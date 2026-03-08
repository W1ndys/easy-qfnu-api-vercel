# CLAUDE.md

此文件为 Claude Code 在处理本仓库代码时提供指导。

## 构建与运行

- 安装 go-task: `go install github.com/go-task/task/v3/cmd/task@latest`
- 前端开发: `task dev:frontend`
- 后端开发: `task dev:backend`
- 前后端并行开发: `task dev`
- 构建前端: `task build:frontend`
- 完整构建（前端 + 后端单二进制）: `task build`
- 测试: `task test`
- Lint: `task lint`

## 部署模式

项目采用单二进制自托管部署：

- 前端 Vue 3 构建产物位于 `frontend/dist`
- Go 使用 `go:embed` 内嵌 `frontend/dist`
- 运行二进制后同时提供 API 与 SPA 静态文件服务

部署前请确保以下项已配置：

- `.env` 中 `PORT`、`GIN_MODE`、`OCR_API_URL`
- 可选 `CORS_ORIGIN`（为空时不启用 CORS）
- ddddocr-api 服务可访问

## 环境变量

- `PORT`: 服务端口，默认 `8141`
- `GIN_MODE`: Gin 模式，`debug` 或 `release`
- `OCR_API_URL`: ddddocr-api 地址，例如 `http://localhost:5000`
- `CORS_ORIGIN`: 可选，允许跨域来源，支持 `*` 或逗号分隔列表
- `FEISHU_WEBHOOK_URL`: 可选，飞书机器人 webhook
- `FEISHU_WEBHOOK_SECRET`: 可选，飞书机器人签名密钥
