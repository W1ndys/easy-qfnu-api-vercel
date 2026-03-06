package zhjw

import (
	"errors"

	"github.com/W1ndys/easy-qfnu-api-vercel/common/response"
	"github.com/W1ndys/easy-qfnu-api-vercel/model"
	zhjwService "github.com/W1ndys/easy-qfnu-api-vercel/services/zhjw"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.CodeInvalidParam, "请提供用户名、密码和验证码")
		return
	}

	cookie, err := zhjwService.LoginWithCaptcha(req.Username, req.Password, req.Captcha)
	if err != nil {
		switch {
		case errors.Is(err, zhjwService.ErrInvalidCredentials):
			response.Fail(c, err.Error())
		case errors.Is(err, zhjwService.ErrCaptchaError):
			response.Fail(c, err.Error())
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
