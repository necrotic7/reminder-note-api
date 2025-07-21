package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
)

func RegisterRouterReminder(r *gin.RouterGroup, h *handlers.ReminderHandler) {
	r.POST("", h.CreateReminder)
	r.GET("", h.GetUserReminders)
}
