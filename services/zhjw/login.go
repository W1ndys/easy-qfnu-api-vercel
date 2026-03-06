package zhjw

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

const maxRetries = 3

// Login 使用用户名和密码完成教务系统登录，返回 cookie 字符串。
// 内部自动处理验证码获取与 OCR 识别，验证码错误时自动重试（最多 3 次）。
func Login(username, password string) (string, error) {
	jar, _ := cookiejar.New(nil)

	client := resty.New()
	client.SetTimeout(loginTimeout)
	client.SetHeader("User-Agent", defaultUA)
	client.GetClient().Jar = jar
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

	// 初始化会话
	_, err := client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return "", fmt.Errorf("初始化会话失败: %w", err)
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		cookie, err := attemptLogin(client, username, password)
		if err == nil {
			return cookie, nil
		}
		if !strings.Contains(err.Error(), "验证码错误") {
			return "", err
		}
		lastErr = err
		slog.Info("验证码错误，正在重试", "attempt", i+1)
	}
	return "", fmt.Errorf("登录失败，已重试 %d 次: %w", maxRetries, lastErr)
}

func attemptLogin(client *resty.Client, username, password string) (string, error) {
	// 获取验证码图片
	resp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %w", err)
	}

	captchaBytes := resp.Body()
	if len(captchaBytes) == 0 {
		return "", fmt.Errorf("验证码图片为空")
	}

	// OCR 识别验证码
	captcha, err := recognizeCaptcha(captchaBytes)
	if err != nil {
		return "", fmt.Errorf("识别验证码失败: %w", err)
	}
	slog.Info("OCR 识别验证码", "captcha", captcha)

	// 编码凭据
	encoded := base64.StdEncoding.EncodeToString([]byte(username)) +
		"%%%" +
		base64.StdEncoding.EncodeToString([]byte(password))

	// 提交登录
	loginResp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"encoded":    encoded,
			"RANDOMCODE": captcha,
		}).
		Post(baseURL + "/jsxsd/xk/LoginToXkLdap")

	if err != nil {
		return "", fmt.Errorf("提交登录请求失败: %w", err)
	}

	body := loginResp.String()

	if strings.Contains(body, "验证码错误") {
		return "", fmt.Errorf("验证码错误")
	}
	if strings.Contains(body, "密码错误") || strings.Contains(body, "用户名或密码错误") {
		return "", fmt.Errorf("用户名或密码错误")
	}

	// 二次验证：访问主页确认登录状态
	verifyResp, err := client.R().Get(baseURL + "/jsxsd/framework/xsMain.jsp")
	if err != nil {
		return "", fmt.Errorf("登录验证请求失败: %w", err)
	}

	if strings.Contains(verifyResp.String(), "用户登录") {
		return "", fmt.Errorf("登录验证失败，请检查账号密码")
	}

	// 提取 cookies
	cookies := client.GetClient().Jar.Cookies(&url.URL{
		Scheme: "http",
		Host:   "zhjw.qfnu.edu.cn",
	})
	if len(cookies) == 0 {
		return "", fmt.Errorf("登录后未获取到 cookie")
	}

	return formatCookies(cookies), nil
}
