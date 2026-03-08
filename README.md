# easy-qfnu-api-lite

QFNU 综合校园服务网关，采用 Go 单二进制自托管部署模式。

## 运行架构

- 前端：Vue 3 SPA，构建产物输出到 `frontend/dist`
- 后端：Gin API 服务
- 打包方式：`go:embed` 将 `frontend/dist` 内嵌到 Go 二进制
- OCR：对接自建 `ddddocr-api`（`POST /ocr/base64`）

## 环境依赖

- Go 1.25+
- Node.js 18+
- Task（go-task）
- ddddocr-api 服务

安装 Task:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

## 环境变量

参考 `.env.example`：

```env
PORT=8141
GIN_MODE=release
OCR_API_URL=http://localhost:5000
CORS_ORIGIN=
FEISHU_WEBHOOK_URL=
FEISHU_WEBHOOK_SECRET=
```

## 开发

```bash
# 启动前端（3000）
task dev:frontend

# 启动后端（8141）
task dev:backend

# 或并行启动
task dev
```

前端开发服务器会将 `/api` 代理到 `http://localhost:8141`。

## 构建与运行

```bash
# 完整构建（先前端，再后端）
task build

# 运行单二进制
./easy-qfnu-api-lite
```

启动后同一进程提供 API 与前端页面，默认地址为 `http://127.0.0.1:8141`。

## Docker（可选）

仓库内提供多阶段构建 `Dockerfile`：

```bash
docker build -t easy-qfnu-api-lite .
docker run -d --name easy-qfnu-api-lite -p 8141:8141 easy-qfnu-api-lite
```

## 健康检查

- `GET /api/health`

## 免责声明

项目仅供学习交流，请遵守目标系统服务条款与当地法律法规。
