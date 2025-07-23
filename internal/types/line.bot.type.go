package types

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
)

type PushMessageParams struct {
	Ctx            context.Context
	Cancel         context.CancelFunc
	NotifyRecordID string
	UserId         string
	Messages       []linebot.SendingMessage
}

type RetryPushMessageParams struct {
	*PushMessageParams
	Retry int
}
