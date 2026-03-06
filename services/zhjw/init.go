package zhjw

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func GetInitCookieAndCaptcha() (cookie string, captchaBase64 string, err error) {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: loginTimeout,
	}

	client := resty.NewWithClient(httpClient)
	client.SetHeader("User-Agent", defaultUA)

	_, err = client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return "", "", fmt.Errorf("访问首页失败: %w", err)
	}

	captchaResp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		return "", "", fmt.Errorf("获取验证码失败: %w", err)
	}

	captchaBase64 = base64.StdEncoding.EncodeToString(captchaResp.Body())

	u, _ := url.Parse(baseURL)
	cookies := jar.Cookies(u)
	cookieStr := formatCookies(cookies)

	return cookieStr, captchaBase64, nil
}
