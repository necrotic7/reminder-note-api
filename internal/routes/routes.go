package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
)

// Router
func RootRouter(
	lineWebhookHandler *handlers.LineWebhookHandler,
) *gin.Engine {
	r := gin.Default()
	RegisterRouterLineWebhook(r.Group("/line/webhook"), lineWebhookHandler)
	return r
}
