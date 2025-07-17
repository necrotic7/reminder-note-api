package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/handlers"
	"github.com/zivwu/reminder-note-api/internal/routes"
	"github.com/zivwu/reminder-note-api/internal/services"

	"go.uber.org/fx"
)

func main() {
	config.InitConfig()
	app := fx.New(
		fx.Provide(
			routes.RootRouter,
		),
		fx.Provide(
			services.NewLineWebhookService,
			handlers.NewLineWebhookHandler,
		),
		fx.Provide(
			handlers.NewReminderHandler,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}

// Server 啟動
func StartServer(r *gin.Engine) {
	r.Run(":8080")
}
