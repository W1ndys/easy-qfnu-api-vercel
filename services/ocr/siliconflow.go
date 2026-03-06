package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	siliconFlowAPIURL = "https://api.siliconflow.cn/v1/chat/completions"
	siliconFlowModel  = "THUDM/GLM-4.1V-9B-Thinking"
	ocrTimeout        = 30 * time.Second
)

type SiliconFlowRequest struct {
	Model    string               `json:"model"`
	Messages []SiliconFlowMessage `json:"messages"`
}

type SiliconFlowMessage struct {
	Role    string               `json:"role"`
	Content []SiliconFlowContent `json:"content"`
}

type SiliconFlowContent struct {
	Type     string               `json:"type"`
	Text     string               `json:"text,omitempty"`
	ImageURL *SiliconFlowImageURL `json:"image_url,omitempty"`
}

type SiliconFlowImageURL struct {
	URL string `json:"url"`
}

type SiliconFlowResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func RecognizeCaptcha(imageData []byte) (string, error) {
	apiKey := os.Getenv("SILICONFLOW_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("SILICONFLOW_API_KEY 环境变量未设置")
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	dataURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)

	request := SiliconFlowRequest{
		Model: siliconFlowModel,
		Messages: []SiliconFlowMessage{
			{
				Role: "user",
				Content: []SiliconFlowContent{
					{
						Type: "text",
						Text: "请识别这个验证码图片中的字符，只返回验证码本身，不要包含其他任何文字或解释。",
					},
					{
						Type: "image_url",
						ImageURL: &SiliconFlowImageURL{
							URL: dataURL,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	client := &http.Client{
		Timeout: ocrTimeout,
	}

	req, err := http.NewRequest("POST", siliconFlowAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var response SiliconFlowResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API 错误: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("API 返回空结果")
	}

	captcha := strings.TrimSpace(response.Choices[0].Message.Content)
	captcha = strings.ReplaceAll(captcha, " ", "")
	captcha = strings.ReplaceAll(captcha, "\n", "")
	captcha = strings.ReplaceAll(captcha, "\r", "")

	return captcha, nil
}
