package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type EnvConfig struct {
	Line struct {
		ChannelSecret      string
		ChannelAccessToken string
	}
}

var Env EnvConfig

func InitConfig() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Println("沒有找到 .env 檔案，跳過")
	}

	// 設定預設值 or 讀環境變數
	viper.AutomaticEnv() // 自動從系統環境變數讀

	Env.Line.ChannelSecret = viper.GetString("LINE_CHANNEL_SECRET")
	Env.Line.ChannelAccessToken = viper.GetString("LINE_CHANNEL_ACCESS_TOKEN")
}
