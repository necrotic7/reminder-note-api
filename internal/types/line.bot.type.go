package types

import "github.com/line/line-bot-sdk-go/linebot"

type PushMessageParams struct {
	UserId   string
	Messages []linebot.SendingMessage
}
