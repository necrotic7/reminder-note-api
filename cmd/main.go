package main

import (
	"context"
	"log"

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
		fx.Invoke(StartServer),
		db.Module,
		routes.Module,
		handlers.Module,
		repositories.Module,
		services.Module,
	)

	app.Run()
}

// Server 啟動
func StartServer(lc fx.Lifecycle, r *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("🟢 Starting Gin server...")
			go func() {
				if err := r.Run(":8080"); err != nil {
					log.Fatalf("🛑 Gin server 啟動失敗：%v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Gin server shutdown")
			return nil
		},
	})
}
