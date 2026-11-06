# 02: 文章发布时间与公开文章集合

**What to build:** 建立所有公开发现功能共同使用的文章发布时间和公开可见性契约，使 feed、sitemap、公开文章页面和 API 对文章状态、首次发布时间和内容更新时间作出一致判断。

**Blocked by:** 01: 稳定文章永久链接

**Status:** resolved

- [x] 文章首次从非 published 状态转为 published 时写入 `publish_at`；后续取消发布和重新发布不清空或重置首次发布时间。
- [x] 历史文章缺失或非法的 `publish_at` 可回填为 `created_at`；有效历史值保留，回填过程可重复并提供可核对的结果。
- [x] 公开文章集合只包含 published 且未软删除的文章；草稿、私密、软删除和作者可见但公众不可见的文章被排除。
- [x] 公开文章查询不会被 status、author_id、show_all 等管理查询参数绕过，并使用稳定的次级排序避免同一时间的结果随机变化。
- [x] 文章内容或公开状态变化可以反映到 `updated_at`；浏览量、收藏数和评论数等互动计数变化不会伪造内容更新时间。
- [x] 文章 feed 时间、sitemap lastmod 和公开集合排序分别遵循首次发布时间与内容更新时间的既定语义。
- [x] 在隔离 MySQL 中覆盖首次发布、取消/重新发布、历史回填、异常值、公开过滤、软删除和稳定排序测试（已实现 opt-in fixtures；本环境未提供 DSN，实际执行跳过，见下文）。

## Answer

- Create 对公开文章设置首次时间；Update 在事务开始时锁文章行，再读取状态及关联，避免并发发布重置首次时间。响应新增 `publish_at`。
- 扩展显式离线 `cmd/article-slugs`：NULL/零日期/部分零日期/非法日历日期使用 created_at 修复，保留有效历史值；全量预检非法 fallback 后再写，统计 scanned/filled/preserved，可幂等重跑。保留 updated_at，包括旧 ON UPDATE schema。支持索引作为迁移完成标记，启动只检查 metadata，不做大批量回填。
- 专用无管理筛选查询：`GetRecentPublicArticles(ctx)` 固定最多50，publish_at DESC/id DESC，仅预加载作者公开姓名；`GetPublicSitemapArticles(ctx)` 单次 SELECT 全量紧凑投影、ID ASC、同快照 count=len；`LastModified()` updated/publish/created 回退。签名和一致性边界见 `docs/agents/public-article-queries.md`。
- 真实隔离 MySQL fixture 覆盖生命周期、并发、历史/非法日期、幂等、未公开及软删除过滤、稳定排序、计数不改变内容时间；真实 SetupUserRouter 验证登录作者/admin 与查询参数也不能绕过公开过滤，以及 publish_at JSON。

## Evidence / limitations

- TDD：`TestArticleResponseFirstPublicationTime` 先执行失败（缺少 publish_at），实现后通过；集合 seam 先因缺失方法/类型编译失败，随后实现。MySQL 集成测试先写，但无 opt-in DSN，不能声称 DB red/green 已执行。
- `go test ./... -count=1`、`go build ./...`、`go vet ./...` 通过。
- `go test -race ./internal/database ./internal/models ./internal/service ./internal/router -count=1` 通过（MySQL fixtures 同样跳过）。
- 定向 `go test ./internal/database ./internal/service ./internal/router -run 'Publication|PublicArticleCollections|PublicPermalinks|Interactions|PublicSitemap' -v -count=1` 明确显示 MySQL 测试 SKIP：缺少 `INKSPACE_SEO_TEST_MYSQL_DSN`；未连接开发数据库、未运行真实迁移。
- sitemap 全量紧凑快照占用 O(公开文章数) 内存；跨 HTTP 请求不承诺冻结数据集，后续分片须保留该边界。无前端或 internal/seo 改动。
