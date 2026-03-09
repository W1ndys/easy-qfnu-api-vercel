package main

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/W1ndys/easy-qfnu-api-lite/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	log.Printf("Starting server on port %s in %s mode...\n", port, ginMode)
	r.Run("0.0.0.0:" + port)
}
