package zhjw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type chatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

func recognizeCaptcha(imageBytes []byte) (string, error) {
	apiKey := os.Getenv("SILICONFLOW_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("未配置 SILICONFLOW_API_KEY 环境变量")
	}

	model := os.Getenv("OCR_MODEL")
	if model == "" {
		model = "THUDM/GLM-4.1V-9B-Thinking"
	}

	b64Image := base64.StdEncoding.EncodeToString(imageBytes)
	dataURL := "data:image/jpeg;base64," + b64Image

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role: "user",
				Content: []map[string]interface{}{
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": dataURL,
						},
					},
					{
						"type": "text",
						"text": "请识别这张验证码图片中的字符，只返回验证码文本，不要返回任何其他内容。",
					},
				},
			},
		},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post("https://api.siliconflow.cn/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("调用 OCR 服务失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("OCR 服务返回错误状态码 %d: %s", resp.StatusCode(), resp.String())
	}

	var chatResp chatResponse
	if err := json.Unmarshal(resp.Body(), &chatResp); err != nil {
		return "", fmt.Errorf("解析 OCR 响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("OCR 服务未返回结果")
	}

	code := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return code, nil
}
