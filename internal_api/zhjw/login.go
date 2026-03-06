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
		response.FailWithCode(c, response.CodeInvalidParam, "请提供用户名和密码")
		return
	}

	const maxRetries = 3
	cookie, err := zhjwService.LoginWithOCR(req.Username, req.Password, maxRetries)
	if err != nil {
		switch {
		case errors.Is(err, zhjwService.ErrInvalidCredentials):
			response.Fail(c, "用户名或密码错误")
		case errors.Is(err, zhjwService.ErrMaxRetriesExceeded):
			response.Fail(c, "登录失败，请稍后重试")
		default:
			response.Fail(c, "登录失败: "+err.Error())
		}
		return
	}

	response.Success(c, model.LoginResponse{
		Cookie: cookie,
	})
}
