package services

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineWebhookService struct {
	BotClient *linebot.Client
}

func NewLineWebhookService() *LineWebhookService {
	bot := InitLineBot()
	return &LineWebhookService{
		BotClient: bot,
	}
}

func (s *LineWebhookService) WebhookRoot(events []*linebot.Event) (err error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// Echo 回傳
				if _, err = s.BotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你說的是："+message.Text)).Do(); err != nil {
					log.Print(err)
					return err
				}
			}
		}
	}

	return nil
}
