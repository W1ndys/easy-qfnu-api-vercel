package zhjw

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// 登录相关哨兵错误
var (
	ErrInvalidCredentials    = errors.New("用户名或密码错误")
	ErrCaptchaFailed         = errors.New("验证码识别失败，请重试")
	ErrOCRServiceUnavailable = errors.New("OCR 服务不可用")
	ErrLoginFailed           = errors.New("登录验证失败")
)

const (
	baseURL       = "http://zhjw.qfnu.edu.cn"
	maxRetries    = 3
	loginTimeout  = 10 * time.Second
	defaultUA     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Login 模拟登录教务系统，返回登录后的 Cookie 字符串
func Login(username, password string) (string, error) {
	ocrURL := os.Getenv("OCR_URL")
	if ocrURL == "" {
		return "", ErrOCRServiceUnavailable
	}

	// 创建带 CookieJar 的客户端
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: loginTimeout,
	}

	client := resty.NewWithClient(httpClient)
	client.SetHeader("User-Agent", defaultUA)

	// 1. 访问首页获取初始 Cookie (JSESSIONID)
	slog.Info("登录流程开始", "username", username)
	_, err := client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return "", fmt.Errorf("获取初始 Cookie 失败: %w", err)
	}

	// 2. 循环尝试验证码识别 + 登录
	encoded := base64.StdEncoding.EncodeToString([]byte(username)) +
		"%%%" +
		base64.StdEncoding.EncodeToString([]byte(password))

	for i := range maxRetries {
		slog.Info("尝试登录", "attempt", i+1)

		// 获取验证码图片
		captchaResp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
		if err != nil {
			slog.Warn("获取验证码失败", "attempt", i+1, "error", err)
			continue
		}

		// 调用 OCR 服务识别验证码
		captchaBase64 := base64.StdEncoding.EncodeToString(captchaResp.Body())
		ocrResp, err := client.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(map[string]string{
				"image": captchaBase64,
			}).
			Post(ocrURL + "/ocr/base64")
		if err != nil {
			slog.Warn("OCR 服务请求失败", "attempt", i+1, "error", err)
			continue
		}
		if ocrResp.StatusCode() != 200 {
			slog.Warn("OCR 服务返回异常", "attempt", i+1, "status", ocrResp.StatusCode())
			continue
		}

		captchaCode := strings.TrimSpace(ocrResp.String())
		slog.Info("验证码识别结果", "code", captchaCode, "attempt", i+1)

		// 提交登录表单
		loginResp, err := client.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(map[string]string{
				"encoded":    encoded,
				"RANDOMCODE": captchaCode,
			}).
			Post(baseURL + "/jsxsd/xk/LoginToXkLdap")
		if err != nil {
			slog.Warn("登录请求失败", "attempt", i+1, "error", err)
			continue
		}

		body := loginResp.String()

		// 检查登录结果
		if strings.Contains(body, "验证码错误") {
			slog.Warn("验证码错误，重试", "attempt", i+1)
			continue
		}
		if strings.Contains(body, "密码错误") || strings.Contains(body, "用户名或密码错误") {
			return "", ErrInvalidCredentials
		}

		// 登录成功，验证登录状态
		slog.Info("登录请求通过，验证登录状态")
		verifyResp, err := client.R().Get(baseURL + "/jsxsd/framework/xsMain.jsp")
		if err != nil {
			return "", fmt.Errorf("验证登录状态失败: %w", err)
		}

		if strings.Contains(verifyResp.String(), "用户登录") {
			return "", ErrLoginFailed
		}

		// 从 CookieJar 中提取 Cookie
		u, _ := url.Parse(baseURL)
		cookies := jar.Cookies(u)
		cookieStr := formatCookies(cookies)

		slog.Info("登录成功", "username", username, "cookie_len", len(cookieStr))
		return cookieStr, nil
	}

	return "", ErrCaptchaFailed
}

// formatCookies 将 Cookie 列表拼接为 "name=value; name2=value2" 格式
func formatCookies(cookies []*http.Cookie) string {
	parts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}
