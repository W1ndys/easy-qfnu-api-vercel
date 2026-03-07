# 项目重构方案：从 Vercel 到自托管单二进制部署

## 一、重构目标

将项目从当前「前端 Vercel 静态托管 + Go Serverless 函数」的部署模式，迁移为「Go 单二进制文件自托管」模式：

- **前端**：Vue 3 SPA 构建产物通过 `go:embed` 嵌入 Go 二进制文件
- **后端**：Go HTTP 服务直接运行在自有服务器上，同时提供 API 和静态文件服务
- **验证码识别**：从 SiliconFlow 大模型 API 切换为自建的 ddddocr-api 服务
- **构建流程**：先构建前端，再将产物嵌入后端，最终输出单个可执行文件

## 二、当前架构分析

### 2.1 现有部署架构

```
Vercel 平台
├── 前端：Vue 3 SPA (frontend/dist) → 静态文件托管
├── 后端：api/index.go → Serverless Go 函数
└── 路由重写：vercel.json
    ├── /api/* → api/index.go
    └── /* → index.html (SPA fallback)
```

### 2.2 现有文件结构（关键部分）

```
项目根目录/
├── api/index.go              ← Vercel Serverless 入口（重构后移除）
├── main.go                   ← Go 入口点（需改造）
├── router/
│   ├── router.go             ← API 路由定义（需改造）
│   └── template.go           ← 遗留模板加载代码（需改造复用）
├── services/zhjw/
│   ├── ocr.go                ← SiliconFlow OCR（需替换）
│   └── login.go              ← 登录流程（OCR 调用处）
├── middleware/
│   └── cors.go               ← CORS 设置（需调整）
├── frontend/                 ← Vue 3 前端（构建配置需调整）
├── vercel.json               ← Vercel 配置（重构后移除）
└── .env.example              ← 环境变量（需更新）
```

### 2.3 当前 OCR 方案

当前使用 SiliconFlow 的视觉语言模型 API 进行验证码识别：

- 接口：`https://api.siliconflow.cn/v1/chat/completions`
- 方式：将验证码图片转为 base64 data URL，通过多模态对话 API 识别
- 环境变量：`SILICONFLOW_API_KEY`、`OCR_MODEL`
- 文件：`services/zhjw/ocr.go`

### 2.4 目标 OCR 方案（ddddocr-api）

切换为自建 ddddocr-api 服务：

- 接口：`POST /ocr/base64`
- 请求体：`{"image": "<base64字符串>", "png_fix": false}`
- 响应体：`{"success": true, "code": 200, "data": {"text": "AB12"}}`
- 环境变量：仅需 `OCR_API_URL`（如 `http://localhost:5000`）

---

## 三、重构阶段划分

整体遵循「先打地基 → 再补墙 → 最后完善细节」的思路，分为 **5 个阶段**：

---

### 阶段一：打地基 —— 构建系统与项目骨架

**目标**：建立新的构建流程，确保前端产物能正确嵌入 Go 二进制文件。

#### 1.1 创建 Taskfile.yml（构建任务管理）

使用 [Task](https://taskfile.dev/)（go-task/task）作为构建任务管理工具。Task 是基于 YAML 的跨平台任务运行器，比 Make 更简洁、对 Windows 友好。

安装方式：`go install github.com/go-task/task/v3/cmd/task@latest`

在项目根目录创建 `Taskfile.yml`，定义以下 task：

| 命令 | 说明 |
|------|------|
| `task dev:frontend` | 启动前端开发服务器 |
| `task dev:backend` | 启动后端开发服务器 |
| `task dev` | 同时启动前后端（开发模式，并行执行） |
| `task build:frontend` | 构建前端（安装依赖 + vite build） |
| `task build` | 完整构建：先构建前端，再 `go build` 输出二进制文件 |
| `task clean` | 清理构建产物 |
| `task test` | 运行 Go 测试 |
| `task lint` | 运行 Go vet |

`Taskfile.yml` 示例结构：

```yaml
version: '3'

vars:
  APP_NAME: easy-qfnu-api-lite

tasks:
  dev:
    desc: 同时启动前后端开发服务器
    deps: [dev:frontend, dev:backend]

  dev:frontend:
    desc: 启动前端开发服务器
    dir: frontend
    cmds:
      - npm run dev

  dev:backend:
    desc: 启动后端开发服务器
    cmds:
      - go run .

  build:frontend:
    desc: 构建前端
    dir: frontend
    cmds:
      - npm install
      - npm run build

  build:
    desc: 完整构建（前端 + 后端）
    deps: [build:frontend]
    cmds:
      - go build -o {{.APP_NAME}} .

  clean:
    desc: 清理构建产物
    cmds:
      - rm -f {{.APP_NAME}} {{.APP_NAME}}.exe
      - rm -rf frontend/dist/*

  test:
    desc: 运行 Go 测试
    cmds:
      - go test ./...

  lint:
    desc: 运行 Go 代码检查
    cmds:
      - go vet ./...
```

构建流程核心逻辑：

```
task build
  → 1. [deps] task build:frontend
       → cd frontend && npm install && npm run build
       → 前端产物输出到 frontend/dist/
  → 2. go build -o easy-qfnu-api-lite .
  → 3. 输出单个二进制文件（内嵌前端资源）
```

#### 1.2 调整前端构建输出

当前 Vite 默认输出到 `frontend/dist/`，这个路径保持不变。Go 的 `embed` 指令将从此路径读取文件。

#### 1.3 建立 embed 机制

在 `main.go`（或新建一个 `embed.go`）中添加 `go:embed` 指令：

- 嵌入路径：`frontend/dist`
- 嵌入变量：`var frontendFS embed.FS`
- 该变量传递给 router 层，用于注册静态文件服务

**关键点**：
- `//go:embed` 指令要求嵌入的目录在编译时必须存在
- 因此 `task build` 必须先执行前端构建（通过 `deps` 自动保证），再执行 `go build`
- 开发模式下可以在 `frontend/dist/` 放一个占位文件（如空的 `index.html`），避免 `go run` 编译失败

#### 1.4 创建 `frontend/dist/.gitkeep`

确保 `frontend/dist/` 目录存在于 Git 仓库中（通过 `.gitkeep` 占位），同时在 `.gitignore` 中忽略该目录下的构建产物：

```gitignore
frontend/dist/*
!frontend/dist/.gitkeep
```

---

### 阶段二：砌墙 —— 静态文件服务与路由改造

**目标**：让 Go 服务器同时提供 API 和前端静态文件服务。

#### 2.1 改造 `router/router.go`

在路由初始化函数中增加前端静态文件服务，核心逻辑：

1. **API 路由**（保持不变）：所有 `/api/*` 路径继续走 API handler
2. **静态文件服务**（新增）：非 API 路径从嵌入的 `frontendFS` 中读取文件
3. **SPA Fallback**（新增）：对于找不到的静态文件路径，统一返回 `index.html`（支持 Vue Router 的 history 模式）

函数签名变更：

- 当前：`func InitRouter() *gin.Engine`
- 改为：`func InitRouter(frontendFS embed.FS) *gin.Engine`（传入嵌入的文件系统）

路由优先级：

```
请求进入
  → 匹配 /api/* → 走 API handler
  → 匹配静态文件（JS/CSS/图片等） → 从 embed.FS 返回
  → 都不匹配 → 返回 index.html（SPA fallback）
```

实现方式：使用 Gin 的 `NoRoute` handler 配合 `io/fs` 包从 `embed.FS` 中读取和服务文件。

#### 2.2 改造 `main.go`

- 声明 `//go:embed frontend/dist` 嵌入变量
- 将嵌入的 `embed.FS` 传递给 `router.InitRouter()`
- 保持原有的端口读取和服务启动逻辑不变

#### 2.3 处理 CORS 中间件

当前 CORS `Access-Control-Allow-Origin` 硬编码为 `"easy-qfnu.top"`。由于前后端合为同源服务，需要调整：

- **生产环境**：前后端同源，不再需要 CORS（可选保留用于外部 API 调用场景）
- **开发环境**：前端 Vite 开发服务器（端口 3000）代理到后端（端口 8141），同样不需要 CORS
- **方案**：将 CORS origin 改为可通过环境变量 `CORS_ORIGIN` 配置，默认为 `*` 或留空（同源不启用）

#### 2.4 清理遗留代码

- **删除** `api/index.go`：Vercel Serverless 入口不再需要
- **删除** `vercel.json`：Vercel 部署配置不再需要
- **改造** `router/template.go`：该文件中的模板加载逻辑不再需要（前端改为 SPA 静态文件），可以删除或保留备用

---

### 阶段三：装门窗 —— OCR 服务切换

**目标**：将验证码识别从 SiliconFlow 大模型 API 切换为自建 ddddocr-api。

#### 3.1 重写 `services/zhjw/ocr.go`

**当前实现**（SiliconFlow 方案）：
- 将图片编码为 base64 data URL
- 构造多模态对话请求（chatRequest）
- 调用 SiliconFlow chat completions API
- 解析对话响应中的验证码文本

**目标实现**（ddddocr-api 方案）：
- 将图片编码为纯 base64 字符串（不带 data URL 前缀）
- 构造简单 JSON 请求：`{"image": "<base64>", "png_fix": false}`
- 调用 `POST <OCR_API_URL>/ocr/base64`
- 解析响应：从 `data.text` 字段提取验证码

**变更点**：
- 删除所有 `chatMessage`、`chatRequest`、`chatChoice`、`chatResponse` 结构体
- 新增 `ocrRequest` 和 `ocrResponse` 结构体
- `recognizeCaptcha` 函数签名保持不变：`func recognizeCaptcha(imageBytes []byte) (string, error)`
- 环境变量变更：
  - 删除 `SILICONFLOW_API_KEY`
  - 删除 `OCR_MODEL`
  - 新增 `OCR_API_URL`（如 `http://localhost:5000`）

#### 3.2 对 `login.go` 的影响

`login.go` 中调用 `recognizeCaptcha()` 的代码不需要修改，因为函数签名保持一致。只要 `ocr.go` 内部实现正确替换即可。

#### 3.3 错误处理适配

ddddocr-api 的错误响应格式：

```json
{
  "success": false,
  "code": 400,
  "message": "错误描述信息",
  "data": null
}
```

在新的 `recognizeCaptcha` 中需要检查：
1. HTTP 状态码（非 200 时报错）
2. 响应体中 `success` 字段（为 false 时使用 `message` 字段报错）
3. `data.text` 为空时报错

---

### 阶段四：通水电 —— 环境变量与配置更新

**目标**：更新环境变量、配置文件和开发文档。

#### 4.1 更新 `.env.example`

```env
# 服务端口
PORT=8141

# Gin 运行模式 (debug/release)
GIN_MODE=release

# OCR 服务地址（ddddocr-api）
OCR_API_URL=http://localhost:5000

# CORS 允许的源（可选，留空则不启用 CORS）
CORS_ORIGIN=

# 飞书机器人（可选）
FEISHU_WEBHOOK_URL=
FEISHU_WEBHOOK_SECRET=
```

**移除的变量**：
- `SILICONFLOW_API_KEY`（不再需要）
- `OCR_MODEL`（不再需要）

**新增的变量**：
- `OCR_API_URL`（ddddocr-api 服务地址，必填）
- `GIN_MODE`（之前硬编码为 release，改为可配置）
- `CORS_ORIGIN`（可选）

#### 4.2 更新 `.gitignore`

新增以下忽略项：

```gitignore
# 构建产物
easy-qfnu-api-lite
easy-qfnu-api-lite.exe
frontend/dist/*
!frontend/dist/.gitkeep

# 可移除的 Vercel 相关忽略（如果有）
.vercel
```

#### 4.3 更新 `CLAUDE.md`

更新构建与运行部分的指令：

- `task build`：完整构建（前端 + 后端）
- `task dev`：开发模式启动（前后端并行）
- `task dev:frontend`：单独启动前端开发服务器
- `task dev:backend`：单独启动后端开发服务器
- 新增部署说明：单二进制文件 + `.env` 配置 + ddddocr-api 服务

---

### 阶段五：精装修 —— 测试、文档与部署优化

**目标**：完善细节，确保可靠部署。

#### 5.1 前端 Vite 配置调整

`frontend/vite.config.js` 需确认以下配置：

- `base: '/'`（保持默认，确保生成的资源路径正确）
- `build.outDir: 'dist'`（保持默认）
- 开发代理配置保持不变（`/api` → `http://localhost:8141`）

无需大改，只需确认构建产物路径与 `go:embed` 指令匹配。

#### 5.2 开发模式体验

开发时前后端分离运行：

```
终端1: task dev:backend    → Go 服务监听 8141 端口（仅 API）
终端2: task dev:frontend   → Vite 服务监听 3000 端口（前端 + 代理 API）
```

或者直接 `task dev` 并行启动两者。

浏览器访问 `http://localhost:3000`，API 请求自动代理到 8141。

#### 5.3 生产模式体验

```
task build                → 输出 easy-qfnu-api-lite 二进制文件
./easy-qfnu-api-lite      → 单服务监听 8141 端口（API + 前端）
```

浏览器访问 `http://服务器IP:8141`，一切由同一个进程提供。

#### 5.4 Dockerfile（可选）

提供一个多阶段构建 Dockerfile：

- 阶段1：Node.js 环境构建前端
- 阶段2：Go 环境构建后端（嵌入前端产物）
- 阶段3：最小运行镜像（仅含二进制文件）

#### 5.5 健康检查

当前已有 `/api/health` 端点，无需修改。可用于部署后的服务健康监测。

#### 5.6 更新 README.md

更新部署文档，说明：

- 环境依赖：Go 1.25+、Node.js 18+、[Task](https://taskfile.dev/)、ddddocr-api 服务
- 构建方式：`task build`
- 运行方式：配置 `.env` 后直接运行二进制文件
- Docker 部署方式（可选）

---

## 四、文件变更清单

### 需要新增的文件

| 文件 | 说明 |
|------|------|
| `Taskfile.yml` | 构建任务管理（基于 go-task/task） |
| `frontend/dist/.gitkeep` | 确保 dist 目录存在于 Git 中 |
| `Dockerfile`（可选） | 多阶段 Docker 构建 |

### 需要修改的文件

| 文件 | 变更内容 |
|------|---------|
| `main.go` | 添加 `go:embed` 指令，传递 `embed.FS` 给 router |
| `router/router.go` | 接收 `embed.FS` 参数，添加静态文件服务和 SPA fallback |
| `services/zhjw/ocr.go` | 整体重写，对接 ddddocr-api |
| `middleware/cors.go` | 改为可配置 origin |
| `.env.example` | 更新环境变量说明 |
| `.gitignore` | 添加构建产物忽略规则 |
| `CLAUDE.md` | 更新构建和运行指令 |
| `go.mod` | 可能需要移除 `gin-contrib/multitemplate` 依赖（如果删除 template.go） |

### 需要删除的文件

| 文件 | 说明 |
|------|------|
| `api/index.go` | Vercel Serverless 入口，不再需要 |
| `vercel.json` | Vercel 部署配置，不再需要 |
| `router/template.go` | 遗留模板加载代码，不再需要 |

---

## 五、风险与注意事项

### 5.1 编译时依赖前端产物

`go:embed` 在编译时要求目标目录及文件存在。如果 `frontend/dist/` 为空或不存在：
- `go build` 会编译失败
- `go run .` 会编译失败
- `go test ./...` 也可能受影响

**应对**：
- 在 `frontend/dist/` 中保留 `.gitkeep` 占位文件
- 使用 `go:embed` 的模式匹配（`all:frontend/dist`）来允许空目录
- 或者在 embed 声明中使用条件逻辑（通过 build tag 区分开发和生产模式）

### 5.2 ddddocr-api 服务可用性

新方案依赖外部的 ddddocr-api 服务。需确保：
- ddddocr-api 服务与主服务部署在同一网络
- 服务端口（默认 5000）可达
- 建议添加超时和重试机制（当前登录已有 3 次重试机制）

### 5.3 CORS 策略变更

当前 CORS 硬编码为 `easy-qfnu.top`。切换到同源部署后：
- 同源请求不需要 CORS
- 但如果有外部系统调用 API，仍需保留 CORS 支持
- 建议改为环境变量可配置

### 5.4 Go module 路径

当前 `go.mod` 中的 module path 为 `github.com/W1ndys/easy-qfnu-api-vercel`，需要同步更新为新名称（如 `github.com/W1ndys/easy-qfnu-api-lite`）。这涉及到所有 import 路径的修改，变动面较大 但有利于项目的长期维护和清晰度。

---

## 六、开发顺序检查表

按照依赖关系排列的执行顺序：

- [ ] **阶段一**
  - [ ] 1.1 创建 `frontend/dist/.gitkeep` 和更新 `.gitignore`
  - [ ] 1.2 创建 `Taskfile.yml`
  - [ ] 1.3 在 `main.go` 添加 `go:embed` 声明
- [ ] **阶段二**
  - [ ] 2.1 改造 `router/router.go`（接收 embed.FS，添加静态文件服务）
  - [ ] 2.2 改造 `main.go`（传递 embed.FS 给 router）
  - [ ] 2.3 验证：`task build` 后运行二进制文件，浏览器能访问前端页面
- [ ] **阶段三**
  - [ ] 3.1 重写 `services/zhjw/ocr.go`（对接 ddddocr-api）
  - [ ] 3.2 验证：登录流程能正常走通
- [ ] **阶段四**
  - [ ] 4.1 更新 `.env.example`
  - [ ] 4.2 调整 `middleware/cors.go`
  - [ ] 4.3 更新 `CLAUDE.md`
- [ ] **阶段五**
  - [ ] 5.1 删除 `api/index.go`、`vercel.json`、`router/template.go`
  - [ ] 5.2 清理 `go.mod` 中不再使用的依赖
  - [ ] 5.3 运行 `go vet ./...` 和 `go test ./...` 确保无问题
  - [ ] 5.4 端到端测试：构建 → 运行 → 登录 → 查成绩/课表等全流程
  - [ ] 5.5 编写/更新 README.md 部署文档
  - [ ] 5.6（可选）创建 Dockerfile
