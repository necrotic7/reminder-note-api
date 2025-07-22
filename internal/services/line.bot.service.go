package services

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/types"
)

type LineBotService struct {
	Client         *linebot.Client
	pushNotifyChan chan types.PushMessageParams
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
		Client:         bot,
		pushNotifyChan: make(chan types.PushMessageParams),
	}

	errChan := make(chan error)

	go func() {
		for {
			select {
			case params := <-s.pushNotifyChan:
				go func(p types.PushMessageParams) {
					log.Println("line bot chan detect")
					if err := s.pushMessage(p); err != nil {
						errChan <- err
					}
				}(params)
			case err := <-errChan:
				log.Println("line bot channel error:", err)
			}

		}
	}()

	return &s
}

func (s *LineBotService) PushToNotifyChan(params types.PushMessageParams) {
	s.pushNotifyChan <- params
}

func (s *LineBotService) pushMessage(params types.PushMessageParams) error {
	_, err := s.Client.PushMessage(params.UserId, params.Messages...).Do()
	if err != nil {
		log.Println("Line Bot發送訊息失敗：", err)
		return err
	}
	log.Println("line push message success")
	return nil
}
