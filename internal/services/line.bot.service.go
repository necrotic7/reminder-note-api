package services

import (
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/consts"
	"github.com/zivwu/reminder-note-api/internal/types"
)

type LineBotService struct {
	Client          *linebot.Client
	pushNotifyChan  chan *types.PushMessageParams
	retryNotifyChan chan *types.RetryPushMessageParams
}

func NewLineBotService() *LineBotService {
	bot, err := linebot.New(
		config.Env.Line.ChannelSecret,
		config.Env.Line.ChannelAccessToken,
	)
	if err != nil {
		log.Panicln("init line bot fail:", err)
	}

	s := LineBotService{
		Client:          bot,
		pushNotifyChan:  make(chan *types.PushMessageParams, consts.Semaphore),
		retryNotifyChan: make(chan *types.RetryPushMessageParams, consts.Semaphore),
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

func (s *LineBotService) PushNotifyMessage(params *types.PushMessageParams) {
	s.pushNotifyChan <- params
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
