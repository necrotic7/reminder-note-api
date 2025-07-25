package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
	"go.uber.org/fx"
)

// Router
func RootRouter(
	lineWebhookHandler *handlers.LineWebhookHandler,
	reminderHandler *handlers.ReminderHandler,
) *gin.Engine {
	r := gin.Default()
	RegisterRouterLineWebhook(r.Group("/line/webhook"), lineWebhookHandler)
	RegisterRouterReminder(r.Group("/reminders"), reminderHandler)
	return r
}

// fx Module
var Module = fx.Module(
	"routes",
	fx.Provide(RootRouter),
)
