# 03: 文章详情页服务端 SEO metadata

**What to build:** 为公开文章详情页提供包含文章级 SEO metadata 的服务端初始 HTML，使搜索引擎和社交平台无需等待 Vue 执行即可获得文章标题、摘要、规范 URL、分享信息和 BlogPosting 结构化数据。

**Blocked by:** 01: 稳定文章永久链接

**Status:** resolved

- [x] 公开 `/blog/<slug>` 的初始 HTML 包含文章 title、description、canonical、Open Graph、Twitter/X Card 和 `BlogPosting` JSON-LD。
- [x] metadata 使用文章标题、摘要、封面、作者、首次发布时间、内容更新时间、站点名称和站点 logo；缺失字段安全省略或按规格回退。
- [x] summary 为空时从正文生成安全的纯文本摘要；Markdown、HTML、脚本、样式、控制字符和非法图片/协议值不会进入 metadata。
- [x] HTML、属性和 JSON-LD 对引号、尖括号、&、脚本片段、反斜杠和中文等特殊内容正确转义，用户内容不能破坏 head 结构。
- [x] 已发布数字 URL 的 301、未公开/不存在文章的 404 与文章初始 HTML metadata 均由服务端决定，不依赖 JavaScript 跳转或登录状态。
- [x] 文章详情 metadata 使用 no-cache 或极短缓存；其他 SPA 路由继续使用普通 fallback，XML/TXT 根端点不被 SPA fallback 接管。
- [x] 覆盖初始 HTML 解析、canonical/OG/Twitter/JSON-LD 字段、特殊字符、无封面/空摘要、公开过滤、301/404 和缓存行为测试。

## Answer

- `internal/seo` 新增显式 canonical origin 校验、绝对 URL、安全公开资源 URL、统一纯文本摘要与文章 metadata 投影。`PUBLIC_SITE_URL` 在所有环境/二进制必填且校验（含 admin fallback），不使用请求 Host；只支持根级 origin。示例配置提供显式本地值，部署传递由 06 完成。
- `SettingService.GetPublicSiteIdentity` 只查询公开且未软删除的站点身份设置，默认名称 `InkSpace`，数据库失败保留错误。文章作者仅选择公开名称，未确认公开的作者 URL 省略；私有上传、签名/凭证 URL、本地网络资源不进入 metadata。
- 复用 01 的公开文章 resolver/render seam，301 Location 改为绝对 canonical URL；`html/template` 和默认 HTML-safe `encoding/json` 安全替换显式 head 区域，保留 Vite 资源标签。源 HTML、构建产物与 fallback 均支持标记；普通博客 fallback 和 `/index.html` 声明 RSS/Atom alternate，不注入文章内容。
- slug API 返回 `data.seo` 服务端同源投影，并改用公开作者 DTO；Vue helper 使用安全 DOM API 更新 title/canonical/OG/Twitter/JSON-LD，支持同组件文章切换、去重、陈旧响应保护、404/离开清理。异步设置不覆盖文章 title，文章加载不再等待设置；MainLayout 使用同一 metadata owner。

## Evidence

- TDD：URL 配置、纯 metadata/summary、公用身份投影、解析初始 HTML 与真实 Vue 路由分别先观察失败，再实现。
- `go test ./... -count=1`、`go build ./...`、`go vet ./...` 通过。
- `go test -race ./internal/seo ./internal/config ./internal/service ./internal/handler ./internal/router ./internal/webassets/blog -count=1` 通过。
- 博客 `npx --yes pnpm@9 run test --run`：15 文件 / 91 测试通过；`run build`、`run lint` 通过。真实 router 测试覆盖 A→B、旧请求晚到、设置晚到、MainLayout title race、404、离开后请求完成以及 head 清理。
- `TestEmbeddedBlogHeadPlaceholders` 在 Vite build 后通过，实际 dist 和 fallback 标记均保留；HTML parser 检查特殊字符、script 结束片段、字段/日期、缺省字段、安全图片、通用 feed alternate、GET/HEAD、no-cache/失败 no-store；公开 slug API 验证不输出作者 email/role。
- `TestPublicPermalinksMySQL` 已扩展真实 SetupUserRouter metadata/身份公开过滤断言并显式执行，但因没有 `INKSPACE_SEO_TEST_MYSQL_DSN` 跳过；未连接开发数据库，不声称真实 MySQL 验证通过。与并发票 02 的该 fixture 修改需在集成合并时保留两组断言。
