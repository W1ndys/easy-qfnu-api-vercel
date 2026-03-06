package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/W1ndys/easy-qfnu-api-vercel/common/logger"
	"github.com/W1ndys/easy-qfnu-api-vercel/router"
)

func main() {
	_ = godotenv.Load()

	gin.SetMode(gin.ReleaseMode)

	logger.InitLogger("./logs", "easy-qfnu-api", "info")

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8141"
	}

	log.Printf("Server starting on port %s", port)
	r.Run("0.0.0.0:" + port)
}
