package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				errMsg := fmt.Sprintf("%v", err)
				path := c.Request.URL.Path
				method := c.Request.Method
				requestID := c.GetString(RequestIDKey)

				logger.L().Error("服务器 Panic",
					zap.String("request_id", requestID),
					zap.String("error", errMsg),
					zap.String("path", path),
					zap.String("method", method),
					zap.String("stacktrace", stack),
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
