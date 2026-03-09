package main

import (
	"embed"
	"os"
	"strings"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/W1ndys/easy-qfnu-api-lite/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	_ = godotenv.Load()
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	ginMode := strings.ToLower(strings.TrimSpace(os.Getenv("GIN_MODE")))
	switch ginMode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
	default:
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := router.InitRouter(frontendFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8141"
	}

	logConfig := logger.LoadConfigFromEnv()
	logger.L().Info("server starting",
		zap.String("port", port),
		zap.String("gin_mode", ginMode),
		zap.String("log_level", logConfig.Level),
		zap.String("log_format", logConfig.Format),
		zap.String("log_file", logConfig.File),
	)

	if err := r.Run("0.0.0.0:" + port); err != nil {
		logger.L().Error("server exited unexpectedly", zap.Error(err))
		logger.Sync()
		os.Exit(1)
	}
}
