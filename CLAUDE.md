# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

InkSpace 是一个 **Gin + GORM 单体后端 + Vue 3 双前端** 的多用户博客系统。注意：本仓库**不是**微服务 / go-zero / RPC / Kafka / WebSocket 架构 —— 如果看到任何提到这些的旧文档，以本文件和实际代码为准。

## 常用命令

```bash
# 后端（三个独立入口，共享 internal/ 代码，各自连同一个 MySQL/Redis）
go run cmd/server/main.go      # 博客/用户服务  :8081（config/config.yaml）
go run cmd/admin/main.go       # 管理后台服务   :8083（config/admin.yaml）
go run cmd/scheduler/main.go   # 定时任务调度器（无 HTTP 端口，可选）

go build ./...                 # 编译全部
go vet ./...                   # 静态检查
gofmt -w <file>                # 格式化

# 依赖数据库（本地开发用 Docker 起 MySQL + Redis）
docker-compose up -d mysql redis

# 前端（两个独立 Vue 应用，包管理用 pnpm）
cd web/blog  && pnpm install && pnpm dev     # 博客前端  :3001
cd web/admin && pnpm install && pnpm dev     # 管理前端  :3002
pnpm build                                   # 生产构建（在各自目录下）
pnpm lint                                    # ESLint --fix

# 测试：当前仓库尚无 Go 单测（无 *_test.go）。新增测试后：
go test ./... -count=1                        # 跑全部
go test ./internal/service/ -run TestXxx -v   # 跑单个测试
```

首次运行需先建库：`CREATE DATABASE inkspace CHARACTER SET = 'utf8mb4';`，表由 GORM `AutoMigrate` 自动创建（`internal/database/migrate.go`）。

## 架构

### 三个二进制、一套 internal/
`cmd/server`、`cmd/admin`、`cmd/scheduler` 是三个 `main`，它们**共享 `internal/` 全部代码**，区别只在于：加载不同配置文件、组装不同路由（`router.SetupUserRouter` vs `router.SetupAdminRouter`）、或不启 HTTP 只跑调度器。新增业务逻辑通常写在 `internal/`，三个入口按需调用。

### 请求分层：handler → service → models
```
HTTP 请求 → router（分组 + 中间件）→ handler（Gin，参数绑定 + 鉴权取值 + 调 service）
          → service（业务逻辑 + 事务 + 缓存）→ models（GORM，操作全局 database.DB）
```
- **handler**（`internal/handler/*_handler.go`）：`NewXxxHandler()` 内部自己 `new` 出对应 service；只做 `c.ShouldBindJSON`、从 `c.Get("user_id")` 取登录态、调 service、用 `utils.*` 返回。**不写业务逻辑**。
- **service**（`internal/service/*_service.go`）：`type XxxService struct{}` + `NewXxxService()`，无状态；直接用全局 `database.DB`（GORM）和 `database.RDB`（Redis）。跨表写操作用 `database.DB.Transaction(func(tx *gorm.DB) error {...})`。
- **models**（`internal/models/*.go`）：GORM 模型 + 内嵌的 `XxxRequest`（`binding:"required"` 绑定结构）+ `ToResponse()`（转对外 DTO，隐藏敏感字段）。同一文件里放模型、请求体、响应转换。

### 关键全局与约定
- **数据库单例**：`database.DB`（`*gorm.DB`，`internal/database/mysql.go`）、`database.RDB`（`*redis.Client`，`redis.go`）。不做依赖注入，全局访问。
- **缓存**：用 `database.SetCache/GetCache/DeleteCache/DeleteCachePattern`（`internal/database/cache.go`），key 形如 `article:%d`、`work:%d`。写数据后主动失效相关缓存。
- **统一响应**：一律用 `internal/utils/response.go` 的 `Success / SuccessWithMessage / PageResponse / BadRequest / Unauthorized / Forbidden / NotFound / Error`。响应体固定 `{code, message, data}`，**成功 code=0**；分页用 `PageResponse` 返回 `{list,total,page,page_size}`。
- **鉴权**：JWT 中间件（`internal/middleware/auth.go`）解析后 `c.Set("user_id"/"username"/"role")`。handler 里 `userID, exists := c.Get("user_id")`，类型是 `uint`。`AuthMiddleware` 强制登录，可选登录版本不设值继续。管理后台额外走 `admin_auth.go`。
- **配置**：Viper 加载，优先级 环境变量 > `.env` > `config/*.yaml`（`internal/config`）。全局 `config.AppConfig`。新增配置项需同步改 `config/*.example.yaml`。
- **日志**：Zap，`utils.InitLogger()` 在每个 `main` 开头调用。
- **定时任务**：实现 `scheduler.Task` 接口（`Run(ctx) error` + `Name() string`），在 `cmd/scheduler/main.go` 用 `sched.RegisterTask(name, task, interval)` 注册。任务文件在 `internal/scheduler/*_task.go`。

### 前端
`web/blog` 与 `web/admin` 是两个独立 Vue 3 应用（Composition API），各自 `views/ components/ layouts/ router/ stores(Pinia)/ utils/`。UI 用 **Element Plus**，HTTP 用 **axios**，Markdown 编辑/渲染用 **vditor**。界面为**中文单语**，无 i18n 框架 —— 不要引入 `vue-i18n`，文案直接写中文。

## 数据库变更须知
- schema 由 GORM struct tag 定义，`AutoMigrate` 自动同步。改字段=改 model 的 `gorm:` tag。
- **只支持 MySQL**（`utf8mb4`）。不要写 SQLite/PostgreSQL 兼容分支。
- 破坏性变更（删列、改类型、加唯一约束）`AutoMigrate` 不会自动处理，需人工确认并提供迁移 SQL（`scripts/`）；变更前先与用户对齐。
- 主键用 `gorm:"primarykey"` 交给 GORM；软删除用内嵌 `gorm.DeletedAt`。
