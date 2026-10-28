# 05: Sitemap 与 robots 发现能力

**What to build:** 为搜索引擎提供包含首页、博客入口和全部公开文章的 sitemap，以及明确公开/私有抓取边界的 robots.txt。

**Blocked by:** 02: 文章发布时间与公开文章集合

**Status:** ready-for-agent

- [ ] `/sitemap.xml` 返回合法 sitemap，包含 `/`、`/blog` 和全部 published 且未软删除文章的 slug URL。
- [ ] 文章 `lastmod` 优先使用内容/公开状态更新时间，再回退到首次发布时间或历史创建时间；首页和博客列表没有可靠更新时间时不伪造 lastmod。
- [ ] sitemap 不包含草稿、私密、软删除或其他不可公开文章，互动计数变化不会导致伪造的 lastmod。
- [ ] 公开 URL 数量超过单个 sitemap 的约 50,000 条协议上限时，端点返回 sitemap index 和稳定分片，且分片不会遗漏或重复公开文章。
- [ ] `/robots.txt` 允许公开文章、feed、sitemap 和静态公开资源，禁止 API、后台和用户私有路径，并指向规范站点 sitemap。
- [ ] sitemap 和 robots 支持 GET/HEAD，使用正确的 UTF-8 MIME 类型；成功响应使用短时公共缓存，数据库失败返回 5xx 且 no-store。
- [ ] 覆盖 XML 解析、公开过滤、lastmod、分片边界、重复/遗漏、robots 规则、MIME、GET/HEAD、缓存和数据库故障测试。
