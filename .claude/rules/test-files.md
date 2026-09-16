---
description: 测试文件规则
globs: ["**/*_test.go"]
---

- 用标准库 `testing`，测试放在被测包内 `xxx_test.go`；现有测试分布在 `internal/{models,service,handler,router,middleware}` 与 `pkg/uploader`
- 测试函数命名 `TestXxx`；多用例用 table-driven（`tests := []struct{...}{...}` + 子测试 `t.Run`）
- service/handler 测试当前均为不依赖 MySQL/Redis 的纯单元测试（handler/router 用 `gin.CreateTestContext` + `httptest`），新增测试优先保持这一点；确需真实 MySQL/Redis 时用独立测试库并在结束后清理数据，不要污染开发库
- 优先测 service 层业务逻辑（分支、边界、错误路径）；handler 层可用 `httptest` + `gin` 引擎测路由与响应码
- 断言成功响应的 `code==0`，失败按 `utils` 里的语义码
