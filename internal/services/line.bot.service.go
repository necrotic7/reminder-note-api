package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/consts"
	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
)

type LineBotService struct {
	Client           *linebot.Client
	pushNotifyChan   chan *types.PushMessageParams
	retryNotifyChan  chan *types.RetryPushMessageParams
	notifyRecordRepo *repositories.NotifyRecordsRepository
}

func NewLineBotService(notifyRecordRepo *repositories.NotifyRecordsRepository) *LineBotService {
	bot, err := linebot.New(
		config.Env.Line.ChannelSecret,
		config.Env.Line.ChannelAccessToken,
	)
	if err != nil {
		log.Panicln("init line bot fail:", err)
	}

	s := LineBotService{
		Client:           bot,
		pushNotifyChan:   make(chan *types.PushMessageParams, consts.Semaphore),
		retryNotifyChan:  make(chan *types.RetryPushMessageParams, consts.Semaphore),
		notifyRecordRepo: notifyRecordRepo,
	}
	s.Start()
	return &s
}

// 主流程
func (s *LineBotService) Start() {
	go func() {
		for {
			select {
			case params := <-s.pushNotifyChan:
				ctx, cancel := context.WithTimeout(context.Background(), consts.Timeout)
				done := make(chan error, 1)
				go func(p *types.PushMessageParams) {
					done <- s.pushMessage(p)
				}(params)

				select {
				case err := <-done:
					if err != nil {
						log.Println("push message channel operate error:", err)
						s.RetryNotifyMessage(&types.RetryPushMessageParams{
							PushMessageParams: params,
							Retry:             0,
						})
					}
				case <-ctx.Done():
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						log.Println("push message channel operate timeout")
						s.RetryNotifyMessage(&types.RetryPushMessageParams{
							PushMessageParams: params,
							Retry:             0,
						})
					}
				}
				cancel()
			case retries := <-s.retryNotifyChan:
				ctx, cancel := context.WithTimeout(context.Background(), consts.Timeout+consts.RetryInterval)
				// retry 3 次
				if retries.Retry >= 3 {
					log.Printf("訊息id %s 已達重試上限\n", retries.NotifyRecordID)
					continue
				}
				done := make(chan error, 1)
				go func(r *types.RetryPushMessageParams) {
					time.Sleep(consts.RetryInterval)
					done <- s.pushMessage(r.PushMessageParams)
				}(retries)

				select {
				case err := <-done:
					if err != nil {
						retries.Retry += 1
						log.Printf("retry push message channel operate error:%v\n, retry counts: %d", err, retries.Retry)
						s.RetryNotifyMessage(retries)
					}
				case <-ctx.Done():
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						retries.Retry += 1
						log.Printf("retry push message channel operate timeout, retry counts: %d\n", retries.Retry)
						s.RetryNotifyMessage(retries)
					}
				}
				cancel()
			}
		}
	}()
}

// 供外部呼叫，寫入發送紀錄並排進推播訊息的channel
func (s *LineBotService) PushNotifyMessage(params *types.PushMessageParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), consts.Timeout)
	defer cancel()
	ID, err := s.notifyRecordRepo.InsertNotifyRecord(ctx, models.InsertNotifyRecord{
		UserID: params.UserId,
		Content: map[string]any{
			"messages": utils.ToJson(params.Messages),
		},
		Status: false, // 第一次寫入十位
		Retry:  0,
	})
	if err != nil {
		log.Println("Push Notify Message insert record fail:", err)
		return err
	}
	params.NotifyRecordID = ID
	s.pushNotifyChan <- params
	return nil
}

// 供外部呼叫，更新發送紀錄並排進重試發送訊息的channel
func (s *LineBotService) RetryNotifyMessage(params *types.RetryPushMessageParams) {
	s.retryNotifyChan <- params
}

// 透過line bot client發送訊息
func (s *LineBotService) pushMessage(params *types.PushMessageParams) error {
	_, err := s.Client.PushMessage(params.UserId, params.Messages...).Do()
	if err != nil {
		log.Println("Line Bot發送訊息失敗：", err)
		return err
	}
	log.Println("line push message success")
	return nil
}
