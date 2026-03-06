package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				errMsg := fmt.Sprintf("%v", err)
				path := c.Request.URL.Path
				method := c.Request.Method

				slog.Error("服务器 Panic",
					"error", errMsg,
					"path", path,
					"method", method,
					"stack", stack,
				)

				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
