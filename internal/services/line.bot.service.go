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
		notifyRecordRepo: notifyRecordRepo,
	}
	s.Start()
	return &s
}

// 主流程
func (s *LineBotService) Start() {
	for i := 0; i < consts.MessageWorker; i++ {
		go func() {
			for params := range s.pushNotifyChan {
				s.pushNotifyChanHandler(params)
			}
		}()
	}
}

func (s *LineBotService) pushNotifyChanHandler(params *types.PushMessageParams) {
	if params.Retry >= 3 {
		log.Printf("訊息id %s 已達重試上限\n", params.NotifyRecordID)
		return
	}

	// 在這裡處理等待，不要進 pushMessage 才 sleep，保持 pushMessage 純粹
	if params.Retry > 0 {
		time.Sleep(consts.RetryInterval)
	}

	ctx, cancel := context.WithTimeout(context.Background(), consts.Timeout)
	defer cancel()

	// 不需要再額外開 go func，因為 Start 已經是在背景執行 handler 了
	// 這樣可以避免 context 提前失效與 goroutine 洩漏
	err := s.pushMessage(ctx, params)

	if err != nil {
		log.Printf("發送失敗，準備重試: %v, 目前重試次數: %d\n", err, params.Retry)
		// 這裡可以直接再丟回 channel
		s.PushNotifyMessage(params)
	}
}

// 供外部呼叫，寫入發送紀錄並排進推播訊息的channel
func (s *LineBotService) PushNotifyMessage(params *types.PushMessageParams) {
	s.pushNotifyChan <- params
}

// 透過line bot client發送訊息
func (s *LineBotService) pushMessage(ctx context.Context, params *types.PushMessageParams) error {
	// 沒有訊息紀錄id，先寫入一筆並取回id
	if utils.IsEmpty(params.NotifyRecordID) {
		ID, err := s.notifyRecordRepo.InsertNotifyRecord(ctx, models.InsertNotifyRecord{
			UserID: params.UserId,
			Content: map[string]any{
				"messages": utils.ToJson(params.Messages),
			},
			Status: false, // 第一次寫入失敗
			Retry:  0,
		})
		params.NotifyRecordID = ID
		if err != nil {
			log.Println("寫入訊息紀錄失敗：", err)
			return err
		}
	}

	status := true
	_, err := s.Client.PushMessage(params.UserId, params.Messages...).Do()
	if err != nil {
		status = false
		params.Retry += 1
		log.Printf("訊息id(%v)發送失敗：%v\n", params.NotifyRecordID, err)
	} else {
		log.Printf("訊息id(%v)發送成功\n", params.NotifyRecordID)
	}

	// update紀錄失敗也不處理
	updateErr := s.notifyRecordRepo.UpdateNotifyRecord(ctx, models.UpdateNotifyRecord{
		ID:     params.NotifyRecordID,
		UserID: params.UserId,
		Status: status,
		Retry:  params.Retry,
	})
	if updateErr != nil {
		log.Printf("訊息id(%v)更新發送紀錄失敗：%v\n", params.NotifyRecordID, err)
	}

	return err
}
