package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespParams struct {
	Status  int
	Message string
	Data    any
}

func Resp(c *gin.Context, p RespParams) {
	c.JSON(p.Status, gin.H{
		"status":  p.Status == http.StatusOK,
		"message": p.Message,
		"data":    p.Data,
	})
}
