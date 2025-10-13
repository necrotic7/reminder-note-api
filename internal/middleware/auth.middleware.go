package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zivwu/reminder-note-api/internal/utils"
)

var ignoreAuthRoutes = []struct {
	Method string
	Path   string
}{
	{
		Method: http.MethodPost,
		Path:   "/users/",
	},
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.FullPath()
		for _, r := range ignoreAuthRoutes {
			if r.Path == url && r.Method == c.Request.Method {
				c.Next()
				return
			}
		}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			utils.Resp(c, utils.RespParams{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Resp(c, utils.RespParams{
				Status:  http.StatusUnauthorized,
				Message: fmt.Sprintf("Auth error: %s", err),
			})
			c.Abort()
			return
		}
		c.Set("tokenInfo", claims)
		c.Next()
	}
}
