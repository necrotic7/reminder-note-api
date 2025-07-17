package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
)

// Router
func RootRouter(
	lineWebhookHandler *handlers.LineWebhookHandler,
	reminderHandler *handlers.ReminderHandler,
) *gin.Engine {
	r := gin.Default()
	RegisterRouterLineWebhook(r.Group("/line/webhook"), lineWebhookHandler)
	RegisterRouterReminder(r.Group("/reminder"), reminderHandler)
	return r
}
