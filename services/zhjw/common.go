package zhjw

import (
	"net/http"
	"strings"
	"time"
)

const (
	baseURL      = "http://zhjw.qfnu.edu.cn"
	loginTimeout = 10 * time.Second
	defaultUA    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

func formatCookies(cookies []*http.Cookie) string {
	parts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}
