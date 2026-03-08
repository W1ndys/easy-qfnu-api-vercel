package main

import (
	"embed"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/W1ndys/easy-qfnu-api-lite/router"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	_ = godotenv.Load()

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

	r.Run("0.0.0.0:" + port)
}
