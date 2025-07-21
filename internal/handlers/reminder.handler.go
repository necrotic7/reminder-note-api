package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/services"
	"github.com/zivwu/reminder-note-api/internal/types"
)

type ReminderHandler struct {
	svc *services.ReminderService
}

func NewReminderHandler(svc *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{
		svc: svc,
	}
}

func (h *ReminderHandler) CreateReminder(c *gin.Context) {
	var req types.ReqCreateReminderBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	err := h.svc.CreateReminder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": fmt.Sprint("fail:", err),
		})
		return
	}
	c.JSON(http.StatusOK, "")
}
