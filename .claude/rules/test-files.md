---
description: 测试文件规则
globs: ["**/*_test.go"]
---

- 当前仓库尚无 Go 单测；新增测试用标准库 `testing`，放在被测包内 `xxx_test.go`
- 测试函数命名 `TestXxx`；多用例用 table-driven（`tests := []struct{...}{...}` + 子测试 `t.Run`）
- service/handler 测试依赖真实 MySQL/Redis 时，用独立测试库并在结束后清理数据；不要污染开发库
- 优先测 service 层业务逻辑（分支、边界、错误路径）；handler 层可用 `httptest` + `gin` 引擎测路由与响应码
- 断言成功响应的 `code==0`，失败按 `utils` 里的语义码
