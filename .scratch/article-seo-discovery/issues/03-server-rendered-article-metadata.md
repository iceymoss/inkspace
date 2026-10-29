# 03: 文章详情页服务端 SEO metadata

**What to build:** 为公开文章详情页提供包含文章级 SEO metadata 的服务端初始 HTML，使搜索引擎和社交平台无需等待 Vue 执行即可获得文章标题、摘要、规范 URL、分享信息和 BlogPosting 结构化数据。

**Blocked by:** 01: 稳定文章永久链接

**Status:** ready-for-agent

- [ ] 公开 `/blog/<slug>` 的初始 HTML 包含文章 title、description、canonical、Open Graph、Twitter/X Card 和 `BlogPosting` JSON-LD。
- [ ] metadata 使用文章标题、摘要、封面、作者、首次发布时间、内容更新时间、站点名称和站点 logo；缺失字段安全省略或按规格回退。
- [ ] summary 为空时从正文生成安全的纯文本摘要；Markdown、HTML、脚本、样式、控制字符和非法图片/协议值不会进入 metadata。
- [ ] HTML、属性和 JSON-LD 对引号、尖括号、&、脚本片段、反斜杠和中文等特殊内容正确转义，用户内容不能破坏 head 结构。
- [ ] 已发布数字 URL 的 301、未公开/不存在文章的 404 与文章初始 HTML metadata 均由服务端决定，不依赖 JavaScript 跳转或登录状态。
- [ ] 文章详情 metadata 使用 no-cache 或极短缓存；其他 SPA 路由继续使用普通 fallback，XML/TXT 根端点不被 SPA fallback 接管。
- [ ] 覆盖初始 HTML 解析、canonical/OG/Twitter/JSON-LD 字段、特殊字符、无封面/空摘要、公开过滤、301/404 和缓存行为测试。
