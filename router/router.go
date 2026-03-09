package router

import (
	"embed"
	"io/fs"
	"mime"
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

	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		panic(err)
	}

	r.StaticFS("/assets", http.FS(assetsFS))
	r.StaticFileFS("/favicon.svg", "favicon.svg", http.FS(distFS))

	r.NoRoute(func(c *gin.Context) {
		requestPath := path.Clean("/" + c.Request.URL.Path)
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") {
			response.ResourceNotFound(c)
			return
		}

		if serveEmbeddedFile(c, distFS, "index.html") {
			return
		}

		c.String(http.StatusNotFound, "frontend dist not found")
	})

	return r
}

func serveEmbeddedFile(c *gin.Context, distFS fs.FS, filePath string) bool {
	cleanPath := strings.TrimPrefix(path.Clean("/"+filePath), "/")
	if cleanPath == "" || cleanPath == "." {
		return false
	}

	data, err := fs.ReadFile(distFS, cleanPath)
	if err != nil {
		return false
	}

	contentType := mime.TypeByExtension(path.Ext(cleanPath))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	c.Data(http.StatusOK, contentType, data)
	return true
}
