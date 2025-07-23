package services

import (
	"context"
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

func (s *LineBotService) Start() {
	go func() {
		for {
			select {
			case params := <-s.pushNotifyChan:
				go func(p *types.PushMessageParams) {
					if err := s.pushMessage(p); err != nil {
						log.Println("push message channel operate error:", err)
						s.RetryNotifyMessage(&types.RetryPushMessageParams{
							PushMessageParams: p,
							Retry:             0,
						})
					}
				}(params)
			case retries := <-s.retryNotifyChan:
				go func(r *types.RetryPushMessageParams) {
					// retry 3 次
					if r.Retry >= 3 {
						log.Println("訊息已達重試上限")
						return
					}
					time.Sleep(consts.RetryInterval)
					if err := s.pushMessage(r.PushMessageParams); err != nil {
						log.Println("retry push message channel operate error:", err)
						r.Retry += 1
						s.RetryNotifyMessage(r)
					}
				}(retries)
			}

		}
	}()
}

func (s *LineBotService) PushNotifyMessage(params *types.PushMessageParams) error {
	// 如果沒有傳入ctx 創建一個預設的
	if utils.IsEmpty(params.Ctx) {
		ctx, cancel := context.WithTimeout(context.Background(), consts.Timeout)
		params.Ctx = ctx
		params.Cancel = cancel
	}
	ID, err := s.notifyRecordRepo.InsertNotifyRecord(params.Ctx, models.InsertNotifyRecord{
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

func (s *LineBotService) RetryNotifyMessage(params *types.RetryPushMessageParams) {
	s.retryNotifyChan <- params
}

func (s *LineBotService) pushMessage(params *types.PushMessageParams) error {
	_, err := s.Client.PushMessage(params.UserId, params.Messages...).Do()
	if err != nil {
		log.Println("Line Bot發送訊息失敗：", err)
		return err
	}
	log.Println("line push message success")
	return nil
}
