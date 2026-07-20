---
name: run-services
description: 启动 / 停止 / 调试 InkSpace 的后端服务与前端。当用户说起服务、跑后端、跑前端、停服务时触发。
---

InkSpace = Gin+GORM 单体（三个入口共享 `internal/`）+ 两个 Vue 前端。本 skill 提供本地启停的标准操作。

## 前置依赖（确认已起）
MySQL 8 + Redis：
```bash
docker-compose up -d mysql redis
docker ps | grep -E "mysql|redis"
```
首次需建库：`CREATE DATABASE inkspace CHARACTER SET = 'utf8mb4';`（表由 AutoMigrate 自动创建）。

## 启动后端（三个独立 main）
```bash
go run cmd/server/main.go      # 博客/用户服务  :8081（config/config.yaml）
go run cmd/admin/main.go       # 管理后台服务   :8083（config/admin.yaml）
go run cmd/scheduler/main.go   # 定时任务调度器（无 HTTP 端口，可选）
```
各进程独立、都连同一个 MySQL/Redis。日志走 Zap 输出到终端。

## 启动前端（两个独立 Vue 应用）
```bash
cd web/blog  && pnpm install && pnpm dev   # 博客前端  :3001
cd web/admin && pnpm install && pnpm dev   # 管理前端  :3002（admin / admin123）
```

## 端口冲突排查
```bash
lsof -i :8081 -i :8083 -i :3001 -i :3002
```
后端端口在 `config/*.yaml`，前端端口在各自 `vite.config.js`。

## 停服务
`Ctrl+C` 停对应进程。查找残留：
```bash
ps aux | grep "cmd/server/main.go" | grep -v grep   # 换成对应入口
kill <pid>
```

## 严格约束
- 不要 `pkill -f go` / `kill -9 -1`，会误伤别的进程
- 不要改 `config/*.yaml` / `.env` 后忘记还原（提交前确认）
- 包管理只用 pnpm，不要 npm/yarn/bun
