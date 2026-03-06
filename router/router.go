package router

import (
	"github.com/W1ndys/easy-qfnu-api-vercel/common/response"
	zhjw "github.com/W1ndys/easy-qfnu-api-vercel/internal_api/zhjw"
	"github.com/W1ndys/easy-qfnu-api-vercel/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
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

	{
		apiV1.GET("/zhjw/captcha", zhjw.GetCaptcha)
		apiV1.POST("/zhjw/login", zhjw.Login)
	}

	zhjwGroup := apiV1.Group("/zhjw")
	zhjwGroup.Use(middleware.AuthRequired())
	{
		zhjwGroup.GET("/grade", zhjw.GetGradeList)
		zhjwGroup.GET("/course-plan", zhjw.GetCoursePlan)
		zhjwGroup.GET("/exam", zhjw.GetExamSchedules)
		zhjwGroup.GET("/selection", zhjw.GetSelectionResults)
		zhjwGroup.GET("/schedule", zhjw.GetClassSchedules)
	}

	return r
}
