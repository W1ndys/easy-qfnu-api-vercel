package zhjw

import (
	"github.com/W1ndys/easy-qfnu-api-vercel/common/response"
	zhjwService "github.com/W1ndys/easy-qfnu-api-vercel/services/zhjw"
	"github.com/gin-gonic/gin"
)

func GetCaptcha(c *gin.Context) {
	imageData, cookie, err := zhjwService.GetCaptchaImage()
	if err != nil {
		response.Fail(c, "获取验证码失败: "+err.Error())
		return
	}

	c.SetCookie("zhjw_session", cookie, 300, "/", "", false, true)
	c.Data(200, "image/jpeg", imageData)
}
