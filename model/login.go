package model

type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	InitCookie string `json:"init_cookie" binding:"required"`
}

type LoginResponse struct {
	Cookie string `json:"cookie"`
}

type InitCookieResponse struct {
	Cookie       string `json:"cookie"`
	CaptchaImage string `json:"captcha_image"`
}
