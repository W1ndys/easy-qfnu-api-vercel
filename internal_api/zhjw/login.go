package zhjw

import (
	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	zhjwService "github.com/W1ndys/easy-qfnu-api-lite/services/zhjw"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginData struct {
	Cookie string `json:"cookie"`
}

// Login 处理教务系统登录请求
func Login(c *gin.Context) {
	log := requestLogger(c)

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("登录请求参数绑定失败", zap.Error(err))
		response.Fail(c, "请提供用户名和密码")
		return
	}

	log.Info("收到教务系统登录请求", zap.String("username", req.Username))
	cookie, err := zhjwService.Login(req.Username, req.Password)
	if err != nil {
		log.Warn("教务系统登录失败",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		response.Fail(c, "登录失败: "+err.Error())
		return
	}

	log.Info("教务系统登录成功", zap.String("username", req.Username))
	response.Success(c, loginData{Cookie: cookie})
}
