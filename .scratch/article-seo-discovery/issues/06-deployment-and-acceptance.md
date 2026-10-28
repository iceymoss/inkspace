# 06: SEO 发现能力上线配置与整体验收

**What to build:** 将稳定链接、文章 metadata、订阅和 sitemap 能力接入实际部署与开发环境，并通过跨端验收确认整套公开发现契约可以安全上线。

**Blocked by:** 01: 稳定文章永久链接; 02: 文章发布时间与公开文章集合; 03: 文章详情页服务端 SEO metadata; 04: RSS 与 Atom 订阅; 05: Sitemap 与 robots 发现能力

**Status:** ready-for-agent

- [ ] `PUBLIC_SITE_URL` 被集中配置、校验并同步到配置示例和部署入口；生产环境缺少或配置非法时拒绝启动，开发/测试环境使用显式地址。
- [ ] 所有 canonical、OG、JSON-LD、feed、sitemap、robots 和文章链接均来自同一个规范站点地址，不信任任意请求 Host 头生成生产 URL。
- [ ] Vite 开发代理能够访问 `/feed.xml`、`/atom.xml`、`/sitemap.xml` 和 `/robots.txt`，本地验证结果与生产路由一致。
- [ ] 跨端 HTTP 验收覆盖 slug 文章 200、旧数字 URL 301、未公开文章 404、初始 HTML metadata、四个发现端点、HEAD、条件请求和数据库故障 5xx。
- [ ] 验收确认 feed、sitemap、robots 和文章页面不会泄露草稿、私密或软删除文章，也不会输出错误的 Host、非法图片协议或未转义用户内容。
- [ ] 运行相关 Go 测试、完整 Go 构建和静态检查、博客前端测试、构建和 lint，并记录隔离 MySQL 验证结果。
- [ ] 上线检查清单包含迁移/回填先行、slug 冲突校验、旧 URL 重定向验证、发现端点状态码/解析验证和上线后的 4xx/5xx、延迟、缓存命中监控。
