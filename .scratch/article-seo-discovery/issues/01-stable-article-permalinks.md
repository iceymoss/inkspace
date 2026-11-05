# 01: 稳定文章永久链接

**What to build:** 为博客文章建立持久化、全局唯一且默认不可变的 slug URL，并让新旧文章访问、公开查询 API 和博客前端路由使用稳定的 slug，同时保留内部数字文章 ID 能力。

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] 文章记录拥有全局唯一、持久化的规范 slug；slug 只使用安全的单路径段格式，软删除、草稿和私密文章也继续占用 slug。
- [x] 新文章创建时自动生成 slug；中文、英文、数字、特殊字符、空结果、超长标题和重复标题都有确定性处理，冲突不会覆盖已有记录。
- [x] 历史文章可以通过可审查、幂等的迁移/回填流程获得 slug；已有非空 slug 不被重写，重复执行不会产生新的重复后缀。
- [x] 数据库在历史数据完成校验后建立全局唯一约束；迁移发现无法处理的冲突时显式失败，不删除或静默跳过数据。
- [x] 公开 slug 查询接口只返回 published 且未软删除的文章，不能通过状态、作者或 show-all 类参数绕过过滤；文章响应携带 slug 并保留数字 ID。
- [x] `/blog/<slug>` 能访问公开文章；已发布的 `/blog/<id>` 由服务端返回 301 到当前 slug；草稿、私密、软删除和不存在文章统一返回 404。
- [x] 浏览器文章路由使用 slug，收藏、评论、编辑等内部操作仍能使用返回的数字 ID。
- [x] 覆盖 slug 生成、冲突、软删除占用、迁移幂等、公开过滤、slug API、301/404 和前端 slug 路由的自动化测试。

## Answer

- 新增 `internal/seo` 固定版本本地拼音 slug 契约、文章持久化字段、事务分配与 MySQL 唯一键冲突/整事务死锁重试。普通 GORM 更新不允许改写 slug；数字 slug 命名空间保留给历史 URL。
- `cmd/article-slugs` 与 `scripts/article-slugs/` 提供离线、幂等的字段准备→全量校验→稳定回填→唯一索引流程；共享初始化锁，失败明确报告，保留原非空 slug 和软删除占用。既有数据库启动只检查迁移结果，不批量回填；仅新建空表自动安装索引。
- 新增 published-only `GetPublicBySlug/GetPublicByID` 与共享 `PublicArticles` scope；`GET /api/articles/slug/:slug` 不绑定可见性查询参数。博客服务显式 GET/HEAD 页面路由执行 301/404、大小写/尾斜杠规范化；metadata 渲染可接收已校验的文章，不重写可见性逻辑。
- 博客前端页面、通知、收藏、评论、编辑跳转和终端公开导航改用 slug；详情在真实 slug 路由切换时加载，旧异步响应不会覆盖新文章，互动操作仍使用返回数字 ID。既有内部数字 API 保留。

## Evidence

- `go test ./... -count=1`、`go build ./...`、`go vet ./...` 全部通过。
- `go test -race ./internal/seo ./internal/models ./internal/service ./internal/handler ./internal/router -count=1` 通过。
- 博客 `pnpm test --run`：14 文件 / 87 测试通过；`pnpm lint`、`pnpm build` 通过。
- 管理前端（仅遗留文章链接改动）`pnpm test --run`：3 文件 / 16 测试通过；`pnpm lint`、`pnpm build` 通过。
- 四个新 MySQL 集成测试已添加并显式执行，但因没有 `INKSPACE_SEO_TEST_MYSQL_DSN` **跳过，未声称真实数据库验证通过**。测试只接受无默认库的显式测试服务器 DSN，随机建/删专属测试库，不使用开发数据库，不依赖 Redis。
- 不执行线上 DDL、不推送。绝对 canonical origin/metadata 由 03 接续，发布时间和公开集合投影由 02 接续；当前 301 `Location` 为规范相对 `/blog/<slug>`。
