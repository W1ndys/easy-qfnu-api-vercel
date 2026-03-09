package middleware

import (
	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthRequired 鉴权中间件
// 作用：强制要求请求必须带 Authorization，否则直接拦截
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization
		Authorization := c.GetHeader("Authorization")
		requestID := c.GetString(RequestIDKey)

		// 2. 检查是否存在
		if Authorization == "" {
			logger.L().Warn("鉴权失败，缺少 Authorization",
				zap.String("request_id", requestID),
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// 如果没有 Authorization，直接报错返回
			response.CookieExpired(c)

			// 🛑 核心步骤：Abort
			// 这一步非常重要！它告诉 Gin 停止执行后面的 Handler，直接返回响应。
			c.Abort()
			return
		}

		// 3. 将 Authorization 放入上下文 (Context)
		// 这样后续的 Handler 就可以直接取用，不用再读 Header 了
		c.Set("Authorization", Authorization)
		logger.L().Debug("鉴权成功",
			zap.String("request_id", requestID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// 4. 放行，执行下一个 Handler
		c.Next()
	}
}
