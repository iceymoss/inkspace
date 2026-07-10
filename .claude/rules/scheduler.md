---
description: 定时任务规则（internal/scheduler）
globs: ["internal/scheduler/**/*.go", "cmd/scheduler/**/*.go"]
---

- 每个任务实现 `scheduler.Task` 接口：`Run(ctx context.Context) error` + `Name() string`，文件放 `internal/scheduler/<name>_task.go`，构造函数 `NewXxxTask()`
- 在 `cmd/scheduler/main.go` 用 `sched.RegisterTask(name, task, interval)` 注册；间隔用 `time.Duration` 常量，不要写死在任务内部
- 任务里复用 service 层逻辑，不要把业务逻辑重复实现一遍
- `Run` 必须尊重 `ctx.Done()`，长循环里定期检查以支持优雅退出
- 任务应**幂等**：同一时间窗被跑两次不产生错误副作用（热榜统计这类天然幂等；写库计数类要用条件更新/去重）
- 任务内 `error` 要返回并记日志，不要 panic 拖垮整个调度器
- 调度器进程与 `cmd/server`、`cmd/admin` 共享同一套 `internal/` 和同一个 MySQL/Redis，注意不要和在线服务抢写同一批热点数据时产生竞态
