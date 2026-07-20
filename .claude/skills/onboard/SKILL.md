---
name: onboard
description: 帮助新成员上手 InkSpace。当用户说新人上手、了解项目、项目介绍时触发。
---

帮助新成员快速理解 InkSpace 的架构与开发流程。

## 步骤

1. 读根 `CLAUDE.md` 和 `README.md`，提炼项目简介
2. 说明三个二进制共享 `internal/`：`cmd/server`（博客 :8081）、`cmd/admin`（后台 :8083）、`cmd/scheduler`（定时任务）
3. 说明请求分层数据流：router（分组+中间件）→ handler（Gin）→ service（业务+事务+缓存）→ models（GORM，全局 `database.DB`）
4. 说明前端：`web/blog`（:3001）与 `web/admin`（:3002），Vue 3 + Element Plus + Pinia，pnpm
5. 列出最常见任务及入口（见输出格式）
6. 指向最重要的规则文件（不复制内容）

## 输出格式

```
## 项目简介
（2-3 句：Gin+GORM 单体博客 + Vue 3 双前端）

## 架构数据流
（一段话或 ASCII 图：router → handler → service → models；三个 main 共享 internal/）

## 开发环境
- 依赖: MySQL 8 + Redis（docker-compose up -d mysql redis）
- 后端: go run cmd/server/main.go / cmd/admin/main.go / cmd/scheduler/main.go
- 前端: cd web/blog|admin && pnpm install && pnpm dev
- 测试: go test ./... -count=1（当前尚无 *_test.go）

## 日常开发入口
- 新接口: /new-api
- 新表:   /new-model
- 起服务: /run-services
- 跑测试: /gotest
- 审 diff: /goreview

## 必读规则
（每条一行 + 链接到 .claude/rules/<file>.md：go-backend / database-model / frontend / scheduler / test-files）
```

简洁输出，不展开细节，需要时让用户读对应规则文件。
