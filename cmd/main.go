package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/db"
	"github.com/zivwu/reminder-note-api/internal/handlers"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/routes"
	"github.com/zivwu/reminder-note-api/internal/services"

	"go.uber.org/fx"
)

func main() {
	config.InitConfig()
	app := fx.New(
		db.Module,
		routes.Module,
		handlers.Module,
		repositories.Module,
		services.Module,
		fx.Invoke(StartServer),
	)

	app.Run()
}

// Server 啟動
func StartServer(r *gin.Engine) {
	r.Run(":8080")
}
