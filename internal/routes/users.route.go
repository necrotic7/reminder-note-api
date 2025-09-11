package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/handlers"
)

func RegisterRouterUsers(r *gin.RouterGroup, h *handlers.UsersHandler) {
	r.POST("/", h.Login)
}
