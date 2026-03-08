package zhjw

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ocrRequest struct {
	Image  string `json:"image"`
	PngFix bool   `json:"png_fix"`
}

type ocrResponseData struct {
	Text string `json:"text"`
}

type ocrResponse struct {
	Success bool             `json:"success"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ocrResponseData `json:"data"`
}

func recognizeCaptcha(imageBytes []byte) (string, error) {
	ocrAPIURL := strings.TrimSpace(os.Getenv("OCR_API_URL"))
	if ocrAPIURL == "" {
		return "", fmt.Errorf("未配置 OCR_API_URL 环境变量")
	}
	ocrAPIURL = strings.TrimRight(ocrAPIURL, "/")

	reqBody := ocrRequest{
		Image:  base64.StdEncoding.EncodeToString(imageBytes),
		PngFix: false,
	}

	var ocrResp ocrResponse
	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&ocrResp).
		SetBody(reqBody).
		Post(ocrAPIURL + "/ocr/base64")

	if err != nil {
		return "", fmt.Errorf("调用 OCR 服务失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("OCR 服务返回错误状态码 %d: %s", resp.StatusCode(), resp.String())
	}

	if !ocrResp.Success {
		msg := strings.TrimSpace(ocrResp.Message)
		if msg == "" {
			msg = "未知错误"
		}
		return "", fmt.Errorf("OCR 服务返回失败(code=%d): %s", ocrResp.Code, msg)
	}

	if ocrResp.Data == nil {
		return "", fmt.Errorf("OCR 服务未返回 data 字段")
	}

	code := strings.TrimSpace(ocrResp.Data.Text)
	if code == "" {
		return "", fmt.Errorf("OCR 服务未返回有效验证码")
	}

	return code, nil
}
