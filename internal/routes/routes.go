package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/handlers"
	"go.uber.org/fx"
)

// Router
func RootRouter(
	lineWebhookHandler *handlers.LineWebhookHandler,
	reminderHandler *handlers.ReminderHandler,
	usersHandler *handlers.UsersHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Env.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterRouterLineWebhook(r.Group("/line/webhook"), lineWebhookHandler)
	RegisterRouterReminder(r.Group("/reminders"), reminderHandler)
	RegisterRouterUsers(r.Group("/users"), usersHandler)
	return r
}

// fx Module
var Module = fx.Module(
	"routes",
	fx.Provide(RootRouter),
)
