package services

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineWebhookService struct {
	BotSvc *LineBotService
}

func NewLineWebhookService(bot *LineBotService) *LineWebhookService {
	return &LineWebhookService{
		BotSvc: bot,
	}
}

/*
message event:

	{
		"replyToken": "f3ddbf65903349c7a0d5649862237570",
		"type": "message",
		"mode": "active",
		"timestamp": 1752734068865,
		"source": {
		"type": "user",
			"userId": "U5e88c4e102a0fb42114ea7ef7596de3c"
		},
		"message": {
			"id": "570262210472313295",
			"type": "text",
			"text": "test"
		}
	}
*/
func (s *LineWebhookService) RootEventHandler(events []*linebot.Event) (err error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// Echo 回傳
				if _, err = s.BotSvc.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("你說的是："+message.Text)).Do(); err != nil {
					log.Print(err)
					return err
				}
			}
		}
	}

	return nil
}
