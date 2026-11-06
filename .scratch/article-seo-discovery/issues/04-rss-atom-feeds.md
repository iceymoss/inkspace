# 04: RSS 与 Atom 订阅

**What to build:** 提供 RSS 2.0 和 Atom 1.0 订阅端点，让订阅工具稳定获取最近公开文章的摘要、规范链接、作者和发布时间。

**Blocked by:** 02: 文章发布时间与公开文章集合

**Status:** resolved

- [x] `/feed.xml` 返回合法 RSS 2.0，`/atom.xml` 返回合法 Atom 1.0；两者使用正确的 UTF-8 MIME 类型。
- [x] RSS 和 Atom 共享同一份最近公开文章集合，最多 50 篇，只包含 published 且未软删除文章，并使用稳定的 `publish_at DESC, id DESC` 排序。
- [x] 文章标题、摘要、作者、首次发布时间、更新时间和 slug canonical URL 在两个格式中语义一致；summary 为空时使用安全的纯文本回退，不输出全文。
- [x] RSS `guid` 和 Atom `id` 使用文章 canonical URL；RSS 和 Atom 都声明自身链接，文章链接不使用旧数字 URL。
- [x] 端点支持 GET 和 HEAD；数据库查询失败返回 5xx 和 no-store，不返回或缓存伪造的空 feed。
- [x] 成功响应使用约 5 分钟公共缓存，并支持 ETag 和/或 Last-Modified 条件请求；无变化时返回 304。
- [x] 站点通用 HTML head 声明 RSS 和 Atom alternate link，且链接使用规范站点地址。
- [x] 覆盖 XML 解析、特殊字符、中文、公开过滤、排序、时间字段、self link、GET/HEAD、错误响应和缓存条件请求测试。

## Answer

- `internal/seo/feed.go` 使用标准 `encoding/xml` 序列化协议中立的公开摘要投影。RSS 使用 Dublin Core `creator` 保存公开姓名，不把姓名错误放入要求邮箱的标准 `author`，也不公开邮箱；RSS Atom `updated` 扩展与 Atom 更新时间保持一致。缺少作者姓名时两个协议均使用公开站点名称；空 Atom 使用稳定 epoch 更新时间并保留站点作者。
- `internal/router/feeds.go` 为两个格式复用 `GetRecentPublicArticles(ctx)`、`GetPublicSiteIdentity(ctx)` 及 ticket03 的 `PublicArticleSEOInput` / `Summary` / canonical helper，不读取请求 Host、列表过滤参数或全文缓存。首次发布时间、修改时间的回退沿用共享投影。HTML alternate links 由已合并 ticket03 提供，无重复实现。
- `internal/router/discovery_response.go` 供 feed/sitemap 共用：完整新鲜查询及序列化后才计算 SHA256 ETag，支持强/弱/list/多 header/wildcard 条件请求；成功 `public, max-age=300`，错误 500/no-store，无正文或错误细节。HEAD 保留表示主要响应头而不返回正文。
- `setupUserRouter` 私有 reader seam 使不依赖数据库的测试经过真实生产路由装配；正式公开 `SetupUserRouter` 接口未变。

## Validation evidence

- TDD red：新增 responder 与 XML 测试分别因缺失函数/类型编译失败；实现后通过。空 Atom 作者、1961 年历史更新时间、ETag 空 list 成员测试先观察失败再修复。
- Passed: `go test ./... -count=1`, `go build ./...`, `go vet ./...`, `go test -race ./internal/router ./internal/seo -count=1`, `git diff --check`。
- Passed: 博客前端 `npx --yes pnpm@9 run test --run`（91 tests）、`run build`、`run lint`；无前端源码改动，验证合并后的 ticket03 alternate/head 不回归。
- `go test ./internal/router -run 'TestFeeds|TestDiscovery' -v -count=1` 通过真实 Gin GET/HEAD、XML 解析、中文/特殊字符、摘要边界、作者隐私、时间/排序保持、条件缓存、非首项编辑/可见性变化与错误测试。
- `TestFeedsMySQLHTTP` 使用 `testutil.SEOTestMySQL` 随机隔离 schema，覆盖真实路由 50 篇限制、首次发布时间/ID 排序、published/软删除过滤、作者/admin + 查询参数绕过尝试、公开设置隐私、状态变更及数据库故障。**当前因未设置 `INKSPACE_SEO_TEST_MYSQL_DSN` 明确 SKIP，不能声称实际 MySQL 验证通过；未连接开发库。**
- 已 merge ticket03 integration `44ab3186` 后执行验证；无 push、无子代理。
