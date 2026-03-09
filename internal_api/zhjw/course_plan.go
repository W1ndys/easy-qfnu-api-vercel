package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-lite/common/request"
	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	zhjwService "github.com/W1ndys/easy-qfnu-api-lite/services/zhjw"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetCoursePlan 是给 Gin 用的处理函数
func GetCoursePlan(c *gin.Context) {
	log := requestLogger(c)

	// 获取参数，能放行到这里，说明已经通过鉴权中间件检查
	Authorization := request.GetCurrentUserAuthorization(c)

	// 调用业务逻辑 (Service 层)
	// 这里的 FetchCoursePlan 首字母是大写，所以能被跨包调用
	data, err := zhjwService.FetchCoursePlan(Authorization)
	// 处理业务结果
	// 如果有错误，返回错误信息
	if errors.Is(err, zhjwService.ErrCookieExpired) {
		log.Warn("培养方案查询失败，Cookie 已过期")
		response.CookieExpired(c)
		return
	} else if errors.Is(err, zhjwService.ErrResourceNotFound) {
		log.Info("培养方案查询无数据")
		response.ResourceNotFound(c)
		return
	} else if err != nil {
		log.Error("培养方案查询失败", zap.Error(err))
		response.FailWithCode(c, 1, "获取培养方案失败: "+err.Error())
		return
	}
	log.Info("培养方案查询成功", zap.Int("group_count", len(data.Groups)))
	response.Success(c, data)

}
