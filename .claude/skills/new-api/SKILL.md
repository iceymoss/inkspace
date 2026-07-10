---
name: new-api
description: 在 InkSpace 里新增一个 HTTP 接口（Gin 路由 + handler + service）。当用户说加接口、新增 api、加路由时触发。
---

在 InkSpace 单体后端新增一个 REST 接口。分层：路由 → handler → service → models。

## 步骤

### 1. 确认输入
- 属于哪个入口：**博客侧**（`internal/router/blog.go`，`cmd/server`）还是**管理后台**（`internal/router/admin.go`，`cmd/admin`）
- 路径与方法（如 `POST /api/articles/:id/like`）
- 是否需要登录（挂 `middleware.AuthMiddleware()`）/ 是否需要管理员（后台走 `admin_auth`）
- 请求字段、响应字段、是否分页

### 2. 请求 / 响应结构
在对应 model 文件（`internal/models/<x>.go`）里加 `XxxRequest`，必填字段带 `binding:"required"`；对外返回用 model 的 `ToResponse()`，不要直接返回 GORM model。

### 3. service 写业务
在 `internal/service/<x>_service.go` 的 `XxxService` 上加方法：
- 用全局 `database.DB`（跨表写用 `database.DB.Transaction`，事务内用 `tx`）
- 需要缓存的读多写少数据用 `database.GetCache/SetCache`，写后 `DeleteCache`/`DeleteCachePattern` 失效
- 返回 `(结果, error)`，不碰 `gin.Context`

### 4. handler 接线
在 `internal/handler/<x>_handler.go` 的 `XxxHandler` 上加方法：
```go
func (h *XxxHandler) DoIt(c *gin.Context) {
    var req models.XxxRequest
    if err := c.ShouldBindJSON(&req); err != nil { utils.BadRequest(c, err.Error()); return }
    userID, exists := c.Get("user_id")
    if !exists { utils.Unauthorized(c, "未登录"); return }
    res, err := h.service.DoIt(&req, userID.(uint))
    if err != nil { utils.Error(c, 400, err.Error()); return }
    utils.Success(c, res)   // 分页用 utils.PageResponse
}
```

### 5. 注册路由
在 `blog.go` / `admin.go` 对应的分组下挂上（public / 需登录 / 管理员）。若是新资源，先在 `SetupXxxRouter` 顶部 `NewXxxHandler()`。

### 6. 手工验证
`go run cmd/server/main.go`（或 `cmd/admin/main.go`）后 curl；确认成功响应 `code==0`。

## 严格约束
- 业务只在 service；handler 只做绑定 + 鉴权取值 + 调 service + `utils.*` 返回
- 统一走 `utils.*` 响应，不要手写 `c.JSON`
- 改/删他人资源必须校验 owner（`user_id` 一致）
- 不要引入 go-zero / RPC / goctl（本项目没有）
