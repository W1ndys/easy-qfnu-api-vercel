package handler

import (
	"net/http"

	"github.com/W1ndys/easy-qfnu-api-vercel/router"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	r = router.InitRouter()
}

func Handler(w http.ResponseWriter, req *http.Request) {
	r.ServeHTTP(w, req)
}
