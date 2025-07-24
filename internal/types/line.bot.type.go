package types

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

type PushMessageParams struct {
	NotifyRecordID string
	UserId         string
	Messages       []linebot.SendingMessage
	Retry          int
}
