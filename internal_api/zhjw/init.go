package zhjw

import (
	"github.com/W1ndys/easy-qfnu-api-vercel/common/response"
	"github.com/W1ndys/easy-qfnu-api-vercel/model"
	zhjwService "github.com/W1ndys/easy-qfnu-api-vercel/services/zhjw"
	"github.com/gin-gonic/gin"
)

func GetInitCookie(c *gin.Context) {
	cookie, captchaBase64, err := zhjwService.GetInitCookieAndCaptcha()
	if err != nil {
		response.Fail(c, "获取初始Cookie失败: "+err.Error())
		return
	}

	response.Success(c, model.InitCookieResponse{
		Cookie:       cookie,
		CaptchaImage: captchaBase64,
	})
}
