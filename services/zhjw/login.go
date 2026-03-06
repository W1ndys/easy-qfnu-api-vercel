package zhjw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrCaptchaError       = errors.New("验证码错误")
	ErrLoginFailed        = errors.New("登录验证失败")
)

func LoginWithCaptcha(username, password, captcha string) (string, error) {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: loginTimeout,
	}

	client := resty.NewWithClient(httpClient)
	client.SetHeader("User-Agent", defaultUA)

	slog.Info("登录流程开始", "username", username)

	_, err := client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return "", fmt.Errorf("获取初始 Cookie 失败: %w", err)
	}

	_, err = client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(username)) +
		"%%%" +
		base64.StdEncoding.EncodeToString([]byte(password))

	slog.Info("尝试登录", "captcha", captcha)

	loginResp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"encoded":    encoded,
			"RANDOMCODE": captcha,
		}).
		Post(baseURL + "/jsxsd/xk/LoginToXkLdap")
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %w", err)
	}

	body := loginResp.String()

	if strings.Contains(body, "验证码错误") {
		return "", ErrCaptchaError
	}
	if strings.Contains(body, "密码错误") || strings.Contains(body, "用户名或密码错误") {
		return "", ErrInvalidCredentials
	}

	slog.Info("登录请求通过，验证登录状态")
	verifyResp, err := client.R().Get(baseURL + "/jsxsd/framework/xsMain.jsp")
	if err != nil {
		return "", fmt.Errorf("验证登录状态失败: %w", err)
	}

	if strings.Contains(verifyResp.String(), "用户登录") {
		return "", ErrLoginFailed
	}

	u, _ := url.Parse(baseURL)
	cookieStr := formatCookies(jar.Cookies(u))

	slog.Info("登录成功", "username", username, "cookie_len", len(cookieStr))
	return cookieStr, nil
}
