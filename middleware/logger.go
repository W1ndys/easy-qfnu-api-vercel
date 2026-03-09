package middleware

import (
	"time"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 记录每次 HTTP 请求的详细信息
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 请求处理完后，记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		requestID := c.GetString(RequestIDKey)
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user_agent", userAgent),
			zap.Int("body_size", bodySize),
		}
		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}
		if errorMessage != "" {
			fields = append(fields, zap.String("error", errorMessage))
		}

		if statusCode >= 500 {
			logger.L().Error("[请求服务器错误]", fields...)
		} else if statusCode >= 400 {
			logger.L().Warn("[请求客户端错误]", fields...)
		} else {
			logger.L().Info("[响应成功]", fields...)
		}
	}
}
