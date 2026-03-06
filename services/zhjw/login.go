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

	"github.com/W1ndys/easy-qfnu-api-vercel/services/ocr"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrCaptchaError       = errors.New("验证码识别失败")
	ErrLoginFailed        = errors.New("登录验证失败")
	ErrMaxRetriesExceeded = errors.New("超过最大重试次数")
)

func LoginWithOCR(username, password string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		slog.Info("登录尝试", "attempt", attempt, "max_retries", maxRetries, "username", username)

		cookie, err := loginOnce(username, password)
		if err == nil {
			slog.Info("登录成功", "username", username, "attempt", attempt)
			return cookie, nil
		}

		lastErr = err

		if errors.Is(err, ErrInvalidCredentials) {
			slog.Error("凭证错误，不再重试", "username", username)
			return "", err
		}

		slog.Warn("登录失败，准备重试", "attempt", attempt, "error", err)
	}

	slog.Error("超过最大重试次数", "username", username, "max_retries", maxRetries)
	return "", fmt.Errorf("%w: %v", ErrMaxRetriesExceeded, lastErr)
}

func loginOnce(username, password string) (string, error) {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: loginTimeout,
	}

	client := resty.NewWithClient(httpClient)
	client.SetHeader("User-Agent", defaultUA)

	_, err := client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return "", fmt.Errorf("访问首页失败: %w", err)
	}

	captchaResp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		return "", fmt.Errorf("获取验证码失败: %w", err)
	}

	captchaImage := captchaResp.Body()
	slog.Info("验证码图片大小", "size", len(captchaImage))

	captcha, err := ocr.RecognizeCaptcha(captchaImage)
	if err != nil {
		slog.Error("OCR识别失败", "error", err)
		return "", fmt.Errorf("OCR识别失败: %w", err)
	}

	slog.Info("验证码识别结果", "captcha", captcha)

	encoded := base64.StdEncoding.EncodeToString([]byte(username)) +
		"%%%" +
		base64.StdEncoding.EncodeToString([]byte(password))

	slog.Info("尝试登录", "username", username)

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
		slog.Warn("验证码错误")
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
