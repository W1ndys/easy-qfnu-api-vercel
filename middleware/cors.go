package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Cors 处理跨域请求,支持 options 访问
func Cors() gin.HandlerFunc {
	originConfig := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))

	return func(c *gin.Context) {
		if originConfig == "" {
			c.Next()
			return
		}

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		switch {
		case originConfig == "*":
			allowOrigin = "*"
		case origin != "":
			for _, allowed := range strings.Split(originConfig, ",") {
				if strings.TrimSpace(allowed) == origin {
					allowOrigin = origin
					break
				}
			}
		}

		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, Authorization")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			if allowOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			logger.L().Debug("CORS 请求通过",
				zap.String("request_id", c.GetString(RequestIDKey)),
				zap.String("origin", origin),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", method),
			)
		} else if origin != "" {
			logger.L().Warn("CORS 请求被拒绝",
				zap.String("request_id", c.GetString(RequestIDKey)),
				zap.String("origin", origin),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", method),
			)
		}

		// 放行所有 OPTIONS 方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}
