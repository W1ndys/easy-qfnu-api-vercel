package zhjw

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func GetCaptchaImage() (imageData []byte, cookie string, err error) {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: loginTimeout,
	}

	client := resty.NewWithClient(httpClient)
	client.SetHeader("User-Agent", defaultUA)

	_, err = client.R().Get(baseURL + "/jsxsd/")
	if err != nil {
		return nil, "", err
	}

	resp, err := client.R().Get(baseURL + "/jsxsd/verifycode.servlet")
	if err != nil {
		return nil, "", err
	}

	u, _ := url.Parse(baseURL)
	cookies := jar.Cookies(u)
	cookieStr := formatCookies(cookies)

	return resp.Body(), cookieStr, nil
}
