---
description: Go 后端代码规则（Gin + GORM 单体）
globs: ["**/*.go"]
---

- 分层职责：handler 只做参数绑定 + 鉴权取值 + 调 service + 返回；业务逻辑放 service；数据操作放 models/GORM。handler 里不写 SQL、不写事务
- 统一响应用 `internal/utils/response.go`（`Success/PageResponse/BadRequest/Unauthorized/...`），不要手写 `c.JSON`；成功 code=0
- 数据库用全局 `database.DB`（`*gorm.DB`），Redis 用全局 `database.RDB`；不要在函数里另建连接
- 跨表写操作用 `database.DB.Transaction(func(tx *gorm.DB) error {...})`，事务内用 `tx` 而不是 `database.DB`
- 缓存读写用 `database.SetCache/GetCache/DeleteCache/DeleteCachePattern`，写数据后主动失效相关 key（如 `article:%d`）
- 错误必须处理，不要用 `_` 忽略 error；返回给前端的错误信息不要泄露内部细节（SQL、堆栈）
- 数据库操作优先用 GORM 方法，避免拼接原始 SQL；确需原始 SQL 时用参数占位符防注入
- 只支持 MySQL，不要写多数据库兼容分支
- 从 `c.Get("user_id")` 取登录态时类型断言为 `uint`，并判断 `exists`
- goroutine 注意并发安全，共享状态加锁；资源（文件句柄、rows）用 defer 关闭
- JSON 序列化用标准库 `encoding/json` 即可（本项目沿用标准库，不引入第三方 JSON 库）
