package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/W1ndys/easy-qfnu-api-vercel/router"
)

func main() {
	_ = godotenv.Load()

	gin.SetMode(gin.ReleaseMode)

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8141"
	}

	r.Run("0.0.0.0:" + port)
}
