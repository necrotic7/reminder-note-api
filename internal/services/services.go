package services

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/utils"
	"go.uber.org/fx"
)

type BaseService struct {
	TokenInfo *utils.TokenClaims
}

func (s *BaseService) SetTokenInfo(c *gin.Context) {
	tokenInfo, ok := c.Get("tokenInfo")
	if !ok {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
		return
	}

	s.TokenInfo = tokenInfo.(*utils.TokenClaims)
}

// 這裡只注入 singleton service，transient(request-level) service 需在各自的 handler 初始化
var Module = fx.Module(
	"services",
	fx.Provide(
		NewScheduleService,
		NewLineBotService,
		NewLineWebhookService,
	),
	fx.Invoke(func(lc fx.Lifecycle, s *ScheduleService) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Println("🟢 Starting ScheduleService")
				return s.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				log.Println("🛑 ScheduleService shutdown")
				return s.Stop(ctx)
			},
		})
	}),
)
