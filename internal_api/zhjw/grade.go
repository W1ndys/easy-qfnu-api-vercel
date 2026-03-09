package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-lite/common/request"
	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	"github.com/W1ndys/easy-qfnu-api-lite/model"
	zhjwService "github.com/W1ndys/easy-qfnu-api-lite/services/zhjw"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetGradeList 是给 Gin 用的处理函数
func GetGradeList(c *gin.Context) {
	log := requestLogger(c)

	// 获取参数，能放行到这里，说明已经通过鉴权中间件检查
	Authorization := request.GetCurrentUserAuthorization(c)

	// 绑定查询参数到结构体
	var req model.GradeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warn("成绩查询参数绑定失败", zap.Error(err))
		response.FailWithCode(c, response.CodeInvalidParam, "查询参数错误，请检查后重试")
		return
	}

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchGrades 首字母是大写，所以能被跨包调用
	data, err := zhjwService.FetchGrades(Authorization, req.Term, req.CourseType, req.CourseName, req.DisplayType)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, zhjwService.ErrCookieExpired) {
		log.Warn("成绩查询失败，Cookie 已过期")
		response.CookieExpired(c)
		return
	} else if errors.Is(err, zhjwService.ErrResourceNotFound) {
		log.Info("成绩查询无数据",
			zap.String("term", req.Term),
			zap.String("course_name", req.CourseName),
		)
		response.ResourceNotFound(c)
		return
	} else if err != nil {
		log.Error("成绩查询失败", zap.Error(err))
		response.FailWithCode(c, 1, "获取成绩失败: "+err.Error())
		return
	}
	log.Info("成绩查询成功",
		zap.String("term", req.Term),
		zap.Int("record_count", len(data.Grades)),
	)
	response.Success(c, data)

}
