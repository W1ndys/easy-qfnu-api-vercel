package zhjw

import (
	"github.com/W1ndys/easy-qfnu-api-lite/common/response"
	zhjwService "github.com/W1ndys/easy-qfnu-api-lite/services/zhjw"
	"github.com/gin-gonic/gin"
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
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "请提供用户名和密码")
		return
	}

	cookie, err := zhjwService.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, "登录失败: "+err.Error())
		return
	}

	response.Success(c, loginData{Cookie: cookie})
}
