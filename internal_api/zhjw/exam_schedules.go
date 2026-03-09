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

// GetExamSchedules 是给 Gin 用的处理函数
func GetExamSchedules(c *gin.Context) {
	log := requestLogger(c)

	// 获取参数，能放行到这里，说明已经通过鉴权中间件检查
	Authorization := request.GetCurrentUserAuthorization(c)

	// 绑定查询参数到结构体
	var req model.ExamSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warn("考试安排查询参数绑定失败", zap.Error(err))
		response.FailWithCode(c, response.CodeInvalidParam, "查询参数错误，请检查后重试")
		return
	}

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchExamSchedules 首字母是大写，所以能被跨包调用
	data, err := zhjwService.FetchExamSchedules(Authorization, req.Term)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, zhjwService.ErrCookieExpired) {
		log.Warn("考试安排查询失败，Cookie 已过期")
		response.CookieExpired(c)
		return
	} else if errors.Is(err, zhjwService.ErrResourceNotFound) {
		log.Info("考试安排查询无数据", zap.String("term", req.Term))
		response.ResourceNotFound(c)
		return
	} else if err != nil {
		log.Error("考试安排查询失败", zap.Error(err))
		response.FailWithCode(c, 1, "获取考试安排失败: "+err.Error())
		return
	}
	log.Info("考试安排查询成功",
		zap.String("term", req.Term),
		zap.Int("record_count", len(data)),
	)
	response.Success(c, data)

}
