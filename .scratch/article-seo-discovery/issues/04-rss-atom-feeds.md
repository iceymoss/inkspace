# 04: RSS 与 Atom 订阅

**What to build:** 提供 RSS 2.0 和 Atom 1.0 订阅端点，让订阅工具稳定获取最近公开文章的摘要、规范链接、作者和发布时间。

**Blocked by:** 02: 文章发布时间与公开文章集合

**Status:** ready-for-agent

- [ ] `/feed.xml` 返回合法 RSS 2.0，`/atom.xml` 返回合法 Atom 1.0；两者使用正确的 UTF-8 MIME 类型。
- [ ] RSS 和 Atom 共享同一份最近公开文章集合，最多 50 篇，只包含 published 且未软删除文章，并使用稳定的 `publish_at DESC, id DESC` 排序。
- [ ] 文章标题、摘要、作者、首次发布时间、更新时间和 slug canonical URL 在两个格式中语义一致；summary 为空时使用安全的纯文本回退，不输出全文。
- [ ] RSS `guid` 和 Atom `id` 使用文章 canonical URL；RSS 和 Atom 都声明自身链接，文章链接不使用旧数字 URL。
- [ ] 端点支持 GET 和 HEAD；数据库查询失败返回 5xx 和 no-store，不返回或缓存伪造的空 feed。
- [ ] 成功响应使用约 5 分钟公共缓存，并支持 ETag 和/或 Last-Modified 条件请求；无变化时返回 304。
- [ ] 站点通用 HTML head 声明 RSS 和 Atom alternate link，且链接使用规范站点地址。
- [ ] 覆盖 XML 解析、特殊字符、中文、公开过滤、排序、时间字段、self link、GET/HEAD、错误响应和缓存条件请求测试。
