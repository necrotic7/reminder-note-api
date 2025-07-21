package services

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/config"
)

type LineBotService struct {
	Client *linebot.Client
}

func NewLineBotService() *LineBotService {
	bot, err := linebot.New(
		config.Env.Line.ChannelSecret,
		config.Env.Line.ChannelAccessToken,
	)
	if err != nil {
		log.Panicln("init line bot fail:", err)
	}

	return &LineBotService{
		Client: bot,
	}
}
