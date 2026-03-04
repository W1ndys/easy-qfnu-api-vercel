package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-vercel/common/response"
	"github.com/W1ndys/easy-qfnu-api-vercel/model"
	zhjwService "github.com/W1ndys/easy-qfnu-api-vercel/services/zhjw"
	"github.com/gin-gonic/gin"
)

// Login 教务系统模拟登录
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "请提供用户名和密码")
		return
	}

	cookie, err := zhjwService.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, zhjwService.ErrInvalidCredentials):
			response.Fail(c, err.Error())
		case errors.Is(err, zhjwService.ErrCaptchaFailed):
			response.Fail(c, err.Error())
		case errors.Is(err, zhjwService.ErrOCRServiceUnavailable):
			response.FailWithCode(c, response.CodeTargetError, "OCR 服务不可用，请联系管理员")
		case errors.Is(err, zhjwService.ErrLoginFailed):
			response.Fail(c, "登录验证失败，请重试")
		default:
			response.Fail(c, "登录失败: "+err.Error())
		}
		return
	}

	response.Success(c, model.LoginResponse{
		Cookie: cookie,
	})
}
