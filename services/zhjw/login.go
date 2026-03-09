package zhjw

import (
	"encoding/base64"
	"fmt"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const maxRetries = 3

// Login 使用用户名和密码完成教务系统登录，返回 cookie 字符串。
// 内部自动处理验证码获取与 OCR 识别，验证码错误时自动重试（最多 3 次）。
func Login(username, password string) (string, error) {
	start := time.Now()
	log := logger.L().With(zap.String("username", username))

	jar, _ := cookiejar.New(nil)

	client := resty.New()
	client.SetTimeout(loginTimeout)
	client.SetHeader("User-Agent", defaultUA)
	client.GetClient().Jar = jar
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

	// 通过响应中间件收集所有 Set-Cookie，不受域名和路径限制
	allCookies := make(map[string]string)
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		for _, cookie := range resp.Cookies() {
			allCookies[cookie.Name] = cookie.Value
		}
		return nil
	})

	// 初始化会话
	_, err := client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		log.Error("初始化登录会话失败", zap.Error(err))
		return "", fmt.Errorf("初始化会话失败: %w", err)
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := attemptLogin(client, username, password)
		if err == nil {
			// 登录成功，格式化收集到的所有 cookie
			if len(allCookies) == 0 {
				log.Error("登录成功后未获取到 Cookie")
				return "", fmt.Errorf("登录后未获取到 cookie")
			}
			parts := make([]string, 0, len(allCookies))
			for name, value := range allCookies {
				parts = append(parts, name+"="+value)
			}
			log.Info("教务系统登录成功",
				zap.Int("attempts", i+1),
				zap.Int("cookie_count", len(allCookies)),
				zap.Duration("latency", time.Since(start)),
			)
			return strings.Join(parts, "; "), nil
		}
		if !strings.Contains(err.Error(), "验证码错误") {
			if strings.Contains(err.Error(), "用户名或密码错误") || strings.Contains(err.Error(), "登录验证失败") {
				log.Warn("教务系统登录失败",
					zap.Error(err),
					zap.Int("attempts", i+1),
					zap.Duration("latency", time.Since(start)),
				)
			} else {
				log.Error("教务系统登录失败",
					zap.Error(err),
					zap.Int("attempts", i+1),
					zap.Duration("latency", time.Since(start)),
				)
			}
			return "", err
		}
		lastErr = err
		log.Warn("验证码错误，准备重试",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
		)
	}
	finalErr := fmt.Errorf("登录失败，已重试 %d 次: %w", maxRetries, lastErr)
	log.Warn("教务系统登录失败",
		zap.Error(finalErr),
		zap.Int("attempts", maxRetries),
		zap.Duration("latency", time.Since(start)),
	)
	return "", finalErr
}

func attemptLogin(client *resty.Client, username, password string) error {
	log := logger.L().With(zap.String("username", username))

	// 获取验证码图片
	resp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		log.Error("获取验证码失败", zap.Error(err))
		return fmt.Errorf("获取验证码失败: %w", err)
	}

	captchaBytes := resp.Body()
	if len(captchaBytes) == 0 {
		log.Error("验证码图片为空")
		return fmt.Errorf("验证码图片为空")
	}

	// OCR 识别验证码
	captcha, err := recognizeCaptcha(captchaBytes)
	if err != nil {
		log.Error("OCR 识别验证码失败", zap.Error(err))
		return fmt.Errorf("识别验证码失败: %w", err)
	}
	log.Info("OCR 识别验证码",
		zap.String("captcha", captcha),
	)

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
		log.Error("提交登录请求失败", zap.Error(err))
		return fmt.Errorf("提交登录请求失败: %w", err)
	}

	body := loginResp.String()

	if strings.Contains(body, "验证码错误") {
		return fmt.Errorf("验证码错误")
	}
	if strings.Contains(body, "密码错误") || strings.Contains(body, "用户名或密码错误") {
		log.Warn("用户名或密码错误")
		return fmt.Errorf("用户名或密码错误")
	}

	// 二次验证：访问主页确认登录状态
	verifyResp, err := client.R().Get(baseURL + "/jsxsd/framework/xsMain.jsp")
	if err != nil {
		log.Error("登录验证请求失败", zap.Error(err))
		return fmt.Errorf("登录验证请求失败: %w", err)
	}

	if strings.Contains(verifyResp.String(), "用户登录") {
		log.Warn("登录验证失败")
		return fmt.Errorf("登录验证失败，请检查账号密码")
	}

	return nil
}
