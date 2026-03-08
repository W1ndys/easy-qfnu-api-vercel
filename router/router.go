package router

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	zhjw "github.com/W1ndys/easy-qfnu-api-lite/internal_api/zhjw"
	"github.com/W1ndys/easy-qfnu-api-lite/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(frontendFS embed.FS) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Cors())

	apiRoot := r.Group("/api")
	{
		apiRoot.GET("/health", func(c *gin.Context) {
			response.Success(c, "API is healthy")
		})
	}

	apiV1 := apiRoot.Group("/v1")

	// 登录接口不需要 AuthRequired 中间件
	apiV1.POST("/zhjw/login", zhjw.Login)

	zhjwGroup := apiV1.Group("/zhjw")
	zhjwGroup.Use(middleware.AuthRequired())
	{
		zhjwGroup.GET("/grade", zhjw.GetGradeList)
		zhjwGroup.GET("/course-plan", zhjw.GetCoursePlan)
		zhjwGroup.GET("/exam", zhjw.GetExamSchedules)
		zhjwGroup.GET("/selection", zhjw.GetSelectionResults)
		zhjwGroup.GET("/schedule", zhjw.GetClassSchedules)
	}

	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		panic(err)
	}

	r.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") {
			response.ResourceNotFound(c)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
		if cleanPath == "" || cleanPath == "." {
			cleanPath = "index.html"
		}

		if serveFromEmbeddedDist(c, distFS, cleanPath) {
			return
		}

		if serveFromEmbeddedDist(c, distFS, "index.html") {
			return
		}

		c.String(http.StatusNotFound, "frontend dist not found")
	})

	return r
}

func serveFromEmbeddedDist(c *gin.Context, distFS fs.FS, filePath string) bool {
	file, err := distFS.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return false
	}

	c.FileFromFS(filePath, http.FS(distFS))
	return true
}
