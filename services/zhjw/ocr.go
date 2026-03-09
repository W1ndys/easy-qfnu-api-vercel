package zhjw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type ocrLegacyResponseData struct {
	Text string `json:"text"`
}

type ocrLegacyResponse struct {
	Success bool                   `json:"success"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *ocrLegacyResponseData `json:"data"`
}

type ocrGenericResponse struct {
	Success *bool           `json:"success"`
	Code    json.RawMessage `json:"code"`
	Message string          `json:"message"`
	Msg     string          `json:"msg"`
	Result  string          `json:"result"`
	Text    string          `json:"text"`
	Data    json.RawMessage `json:"data"`
}

func recognizeCaptcha(imageBytes []byte) (string, error) {
	ocrAPIURL := strings.TrimSpace(os.Getenv("OCR_API_URL"))
	if ocrAPIURL == "" {
		return "", fmt.Errorf("未配置 OCR_API_URL 环境变量")
	}
	ocrAPIURL = normalizeOCRURL(ocrAPIURL)
	start := time.Now()

	logger.L().Debug("开始调用 OCR 服务",
		zap.String("ocr_url", ocrAPIURL),
	)

	client := resty.New().SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetFormData(map[string]string{
			"image":       base64.StdEncoding.EncodeToString(imageBytes),
			"probability": "false",
			"png_fix":     "false",
		}).
		Post(ocrAPIURL)

	if err != nil {
		logger.L().Error("调用 OCR 服务失败",
			zap.String("ocr_url", ocrAPIURL),
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return "", fmt.Errorf("调用 OCR 服务失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.L().Error("OCR 服务返回非 200 状态码",
			zap.String("ocr_url", ocrAPIURL),
			zap.Int("status", resp.StatusCode()),
			zap.String("response_snippet", truncateText(resp.String(), 200)),
			zap.Duration("latency", time.Since(start)),
		)
		return "", fmt.Errorf("OCR 服务返回错误状态码 %d: %s", resp.StatusCode(), resp.String())
	}

	code, parseErr := parseOCRCode(resp.Body())
	if parseErr != nil {
		logger.L().Error("OCR 响应解析失败",
			zap.String("ocr_url", ocrAPIURL),
			zap.Error(parseErr),
			zap.String("response_snippet", truncateText(resp.String(), 200)),
			zap.Duration("latency", time.Since(start)),
		)
		return "", fmt.Errorf("OCR 响应解析失败: %w，响应体: %s", parseErr, truncateText(resp.String(), 200))
	}

	logger.L().Info("OCR 识别成功",
		zap.String("ocr_url", ocrAPIURL),
		zap.String("captcha", code),
		zap.Duration("latency", time.Since(start)),
	)

	return code, nil
}

func normalizeOCRURL(rawURL string) string {
	rawURL = strings.TrimRight(strings.TrimSpace(rawURL), "/")
	if strings.HasSuffix(rawURL, "/ocr") {
		return rawURL
	}
	return rawURL + "/ocr"
}

func parseOCRCode(respBody []byte) (string, error) {
	bodyText := strings.TrimSpace(string(respBody))
	if bodyText == "" {
		return "", fmt.Errorf("OCR 响应体为空")
	}

	// 兼容纯文本返回，例如直接返回 "ABCD"
	if !strings.HasPrefix(bodyText, "{") && !strings.HasPrefix(bodyText, "[") {
		return bodyText, nil
	}

	// 兼容历史格式：{"success":true,"data":{"text":"ABCD"}}
	var legacyResp ocrLegacyResponse
	if err := json.Unmarshal(respBody, &legacyResp); err == nil {
		if legacyResp.Data != nil {
			if code := strings.TrimSpace(legacyResp.Data.Text); code != "" {
				return code, nil
			}
		}
		if !legacyResp.Success && (legacyResp.Message != "" || legacyResp.Code != 0) {
			msg := strings.TrimSpace(legacyResp.Message)
			if msg == "" {
				msg = "未知错误"
			}
			return "", fmt.Errorf("OCR 服务返回失败(code=%d): %s", legacyResp.Code, msg)
		}
	}

	// 兼容常见格式：{"result":"ABCD"} / {"text":"ABCD"} / {"data":"ABCD"} / {"data":{"text":"ABCD"}}
	var genericResp ocrGenericResponse
	if err := json.Unmarshal(respBody, &genericResp); err != nil {
		return "", fmt.Errorf("OCR 响应不是有效 JSON: %w", err)
	}

	if genericResp.Success != nil && !*genericResp.Success {
		msg := strings.TrimSpace(genericResp.Message)
		if msg == "" {
			msg = strings.TrimSpace(genericResp.Msg)
		}
		if msg == "" {
			msg = "未知错误"
		}
		return "", fmt.Errorf("OCR 服务返回失败: %s", msg)
	}

	for _, value := range []string{genericResp.Result, genericResp.Text} {
		if code := strings.TrimSpace(value); code != "" {
			return code, nil
		}
	}

	if len(genericResp.Data) > 0 && strings.TrimSpace(string(genericResp.Data)) != "null" {
		var dataText string
		if err := json.Unmarshal(genericResp.Data, &dataText); err == nil {
			if code := strings.TrimSpace(dataText); code != "" {
				return code, nil
			}
		}

		var dataObj map[string]any
		if err := json.Unmarshal(genericResp.Data, &dataObj); err == nil {
			for _, key := range []string{"text", "result", "code", "captcha"} {
				if rawValue, ok := dataObj[key].(string); ok {
					if code := strings.TrimSpace(rawValue); code != "" {
						return code, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("OCR 服务未返回有效验证码")
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
