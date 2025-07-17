package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
)

func RegisterRouterLineWebhook(r *gin.RouterGroup, h *handlers.LineWebhookHandler) {
	r.POST("", h.WebhookHandler)
}
