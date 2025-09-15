package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	AllowOrigins []string
	SecretKey    string
	Line         struct {
		ChannelSecret      string
		ChannelAccessToken string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
	}
}

var Env EnvConfig

func InitConfig() EnvConfig {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Println("沒有找到 .env 檔案，跳過")
	}

	// 設定預設值 or 讀環境變數
	viper.AutomaticEnv() // 自動從系統環境變數讀

	Env.AllowOrigins = GetStringSlice("ALLOW_ORIGINS", ",")

	Env.SecretKey = viper.GetString("SECRET_KEY")

	Env.Line.ChannelSecret = viper.GetString("LINE_CHANNEL_SECRET")
	Env.Line.ChannelAccessToken = viper.GetString("LINE_CHANNEL_ACCESS_TOKEN")

	Env.DB.Host = viper.GetString("MONGO_DB_HOST")
	Env.DB.Port = viper.GetString("MONGO_DB_PORT")
	Env.DB.User = viper.GetString("MONGO_DB_USER")
	Env.DB.Password = viper.GetString("MONGO_DB_PASSWORD")

	log.Println("init config success")

	return Env
}

func GetStringSlice(key string, sep string) []string {
	return strings.Split(viper.GetString(key), sep)
}
