package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{}
}

func (h *ReminderHandler) CreateReminder(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
