package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/services"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
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
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	err := h.svc.CreateReminderFlow(c.Request.Context(), &req)
	if err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("fail:", err),
		})
		return
	}
	utils.Resp(c, utils.RespParams{Status: http.StatusOK})
}

func (h *ReminderHandler) GetUserReminders(c *gin.Context) {
	var query types.ReqGetUserRemindersQuery
	if err := c.BindQuery(&query); err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	result, err := h.svc.GetUserReminders(c.Request.Context(), query)
	if err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("fail: ", err),
		})
		return
	}
	utils.Resp(c, utils.RespParams{
		Status: http.StatusOK,
		Data:   result,
	})
}

func (h *ReminderHandler) UpdateReminder(c *gin.Context) {
	var body types.ReqUpdateReminderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	err := h.svc.UpdateReminderFlow(c.Request.Context(), &body)
	if err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("fail: ", err),
		})
		return
	}
	utils.Resp(c, utils.RespParams{
		Status: http.StatusOK,
	})
}

func (h *ReminderHandler) DeleteReminder(c *gin.Context) {
	var body types.ReqDeleteReminderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	err := h.svc.DeleteReminderFlow(c.Request.Context(), &body)
	if err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("fail: ", err),
		})
		return
	}

	utils.Resp(c, utils.RespParams{
		Status: http.StatusOK,
	})
}
