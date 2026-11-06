# 06: SEO 发现能力上线配置与整体验收

**What to build:** 将稳定链接、文章 metadata、订阅和 sitemap 能力接入实际部署与开发环境，并通过跨端验收确认整套公开发现契约可以安全上线。

**Blocked by:** 01: 稳定文章永久链接; 02: 文章发布时间与公开文章集合; 03: 文章详情页服务端 SEO metadata; 04: RSS 与 Atom 订阅; 05: Sitemap 与 robots 发现能力

**Status:** resolved

- [x] `PUBLIC_SITE_URL` 被集中配置、校验并同步到配置示例和部署入口；生产环境缺少或配置非法时拒绝启动，开发/测试环境使用显式地址。
- [x] 所有 canonical、OG、JSON-LD、feed、sitemap、robots 和文章链接均来自同一个规范站点地址，不信任任意请求 Host 头生成生产 URL。
- [x] Vite 开发代理能够访问 `/feed.xml`、`/atom.xml`、`/sitemap.xml` 和 `/robots.txt`，本地验证结果与生产路由一致。
- [x] 跨端 HTTP 验收覆盖 slug 文章 200、旧数字 URL 301、未公开文章 404、初始 HTML metadata、四个发现端点、HEAD、条件请求和数据库故障 5xx。
- [x] 验收确认 feed、sitemap、robots 和文章页面不会泄露草稿、私密或软删除文章，也不会输出错误的 Host、非法图片协议或未转义用户内容。
- [x] 运行相关 Go 测试、完整 Go 构建和静态检查、博客前端测试、构建和 lint，并记录隔离 MySQL 验证结果。
- [x] 上线检查清单包含迁移/回填先行、slug 冲突校验、旧 URL 重定向验证、发现端点状态码/解析验证和上线后的 4xx/5xx、延迟、缓存命中监控。

## Answer

完成部署与验收实现，**resolved 表示代码/配置/测试/文档交付，不表示生产或真实 MySQL 已验收**。

- 继承03集中 canonical 校验；补齐两个 Compose 全部五个应用服务、三个 K8s base 与 development
  overlay、配置示例。base 只引用环境拥有的必需 ConfigMap，不写开发 canonical URL。
  实际渲染发现并修复 namespace 导致的 hash 引用错配，增加三个 overlay 离线回归验证。
- 启动子进程测试覆盖三个业务二进制 debug/test/release、缺失/空 URL、query/fragment（包括空 ?/#）、
  非法协议以及 admin fallback。测试先失败发现 cmd/admin 会在非法配置后重试默认文件并进入 DB，
  已移除这个绕过；文件缺失的合法 fallback 仍由共享 loader 承担。
- Vite 四端点与 `/sitemaps/` 命名空间 + `/api/settings/public` 运行时配置，真实 Vite→HTTP fixture
  测试先 red 后 green；不引入编译期 canonical 地址。
- 同一个完整生产 Gin assembly 注入公开 reader，跨 HTML/feeds/sitemap/robots 校验共同 URL/摘要、
  200/301/404、HEAD/304、MIME/cache、Host 独立性、恶意内容/图片/隐私不泄露和 read-error 500/no-store。
  扩展已有 gated `TestPublicPermalinksMySQL` 对真实 SQL/公开 API/runtime canonical 做同组验证，
  撤下文章立即移除，随机测试 schema 内删表模拟 DB 故障；**该真实 SQL fixture 因 DSN 缺失 SKIP**。
- [SEO_ROLLOUT.md](../../../docs/SEO_ROLLOUT.md)：迁移先行、ticket02 扩展命令重跑、completion index、
  review SQL + 只读 verify.sql、curl smoke/人工清单、监控缓存语义/延迟/错误和安全回滚。
  `scripts/seo-smoke.py` 不打印正文/secret，使用一次性文件；disposable fixture 覆盖成功和错误 canonical 拒绝。

验证（2026-11-05）：完整 `go test ./... -count=1`、`go build ./...`、`go vet ./...`、
config/router/seo race 均 PASS；博客 frozen install、**16 files / 98 tests**、build、lint PASS；
embedded dist head markers、Python curl smoke fixture、Compose 两基础及两个 HTTPS 组合渲染、
三 K8s overlay namespace/hash/ref 回归、diff 检查 PASS。Compose 缺失/空 origin 按预期拒绝。
显式运行数据库/service/router SEO fixtures，**13 个 MySQL fixture SKIP**（无 `INKSPACE_SEO_TEST_MYSQL_DSN`），
未连接开发库或 live Redis。Docker daemon 不可用，未实际构建/启动容器、apply 集群或运行生产 canary。
完整证据与发布门禁见 [SEO_VALIDATION.md](../../../docs/SEO_VALIDATION.md)。
