package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/services"
)

type LineWebhookHandler struct {
	svc *services.LineWebhookService
}

func NewLineWebhookHandler(svc *services.LineWebhookService) *LineWebhookHandler {
	return &LineWebhookHandler{svc: svc}
}

func (h *LineWebhookHandler) WebhookHandler(c *gin.Context) {
	events, err := h.svc.BotClient.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	err = h.svc.WebhookRoot(events)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.Status(http.StatusOK)
}
