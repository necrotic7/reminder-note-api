package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/services"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
)

type UsersHandler struct {
	svc *services.UsersService
}

func NewUsersHandler(svc *services.UsersService) *UsersHandler {
	return &UsersHandler{
		svc: svc,
	}
}

func (h *UsersHandler) Login(c *gin.Context) {
	var req types.ReqLoginBody
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("invalid parameters:", err),
		})
		return
	}

	result, err := h.svc.Login(c.Request.Context(), &req)

	if err != nil {
		utils.Resp(c, utils.RespParams{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("fail:", err),
		})
		return
	}
	utils.Resp(c, utils.RespParams{
		Status: http.StatusOK,
		Data:   result,
	})
}
