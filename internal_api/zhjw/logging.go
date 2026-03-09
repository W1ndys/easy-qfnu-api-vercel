package zhjw

import (
	"github.com/W1ndys/easy-qfnu-api-lite/middleware"
	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func requestLogger(c *gin.Context) *zap.Logger {
	return logger.L().With(
		zap.String("request_id", c.GetString(middleware.RequestIDKey)),
		zap.String("ip", c.ClientIP()),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
}
