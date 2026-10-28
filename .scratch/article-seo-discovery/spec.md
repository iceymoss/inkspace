# 博客文章 SEO 与公开发现能力

Status: ready-for-agent

## Problem Statement

InkSpace 的博客目前主要依赖客户端渲染，文章公开地址使用数字文章 ID。文章详情页没有稳定的语义 URL，也没有文章级的服务端初始 metadata，因此搜索引擎、RSS 阅读器、微信/知乎等社交平台以及其他知识管理工具难以可靠地发现、理解和长期引用文章。

站点也缺少 RSS、Atom、Sitemap 和 robots.txt 这些标准发现入口。文章标题或内容变化时，外部链接没有稳定的 URL 契约；旧数字地址、文章可见性、发布时间和 sitemap/feed 内容之间也没有明确一致的规则。

## Solution

为博客文章建立一套以公开文章为中心的 SEO 与发现能力：

- 为文章持久化生成稳定、全局唯一且默认不可变的 slug，并以 slug 作为新的文章公开 URL。
- 保留已发布文章的旧数字 URL，并由服务端永久重定向到 slug URL。
- 由 Go 服务端为文章详情初始 HTML 注入文章级 title、description、canonical、Open Graph、Twitter Card 和 BlogPosting JSON-LD；正文继续由现有 Vue 应用客户端渲染。
- 增加 RSS 2.0、Atom 1.0、Sitemap 和 robots.txt 根级端点。
- 让 feed、sitemap、公开文章 API、服务端文章页面和前端文章页面共享一致的公开可见性与 canonical URL 语义。
- 通过显式的 `PUBLIC_SITE_URL` 生成所有规范绝对 URL，并在生产环境缺失或非法时拒绝启动。
- 通过可审查、幂等的数据库迁移和历史数据回填建立 slug 与发布时间数据基础。

## User Stories

1. 作为博客读者，我希望文章有可读、稳定的 slug URL，以便理解和分享文章主题。
2. 作为已有文章链接的读者，我希望旧的数字文章地址继续工作，以便历史书签和外部引用不会失效。
3. 作为搜索引擎，我希望访问文章 URL 时无需依赖 JavaScript 执行即可获得标题、摘要、规范 URL 和结构化数据，以便正确建立索引。
4. 作为社交平台抓取器，我希望文章初始 HTML 包含 Open Graph 和 Twitter/X metadata，以便生成正确的分享预览。
5. 作为知识管理工具，我希望文章 URL 在标题修改后保持稳定，以便长期引用不会因为编辑标题而漂移。
6. 作为 RSS 阅读器用户，我希望订阅博客的 RSS feed，以便及时发现新发布的公开文章。
7. 作为 Atom 阅读器用户，我希望订阅博客的 Atom feed，以便使用兼容 Atom 的订阅工具。
8. 作为搜索引擎，我希望从 sitemap 发现全部公开文章，以便不依赖站内链接遍历全部内容。
9. 作为爬虫，我希望通过 robots.txt 得到公开内容、私有应用路径和 sitemap 的明确指引。
10. 作为博客运营者，我希望草稿、私密文章和软删除文章不会进入文章页面、feed、sitemap 或 JSON-LD，以避免意外泄露。
11. 作为博客运营者，我希望已发布文章的数字旧地址服务端返回 301，而不是依赖前端 JavaScript 跳转，以便搜索引擎传递链接权重。
12. 作为博客运营者，我希望不存在、未发布和无权访问的文章公开表现一致为 404，以避免泄露文章存在性。
13. 作为作者，我希望标题包含中文、英文、数字或特殊字符时仍能生成安全、可读的 slug。
14. 作为作者，我希望相同标题的文章能够获得确定性的唯一 slug，而不会因为并发或重复迁移产生不可预测的地址。
15. 作为维护者，我希望软删除文章继续占用原 slug，以避免未来新文章复用旧地址并改变历史链接的语义。
16. 作为维护者，我希望 slug 回填可以重复执行而不重写已有 slug、追加重复后缀或覆盖文章数据。
17. 作为维护者，我希望 slug 冲突和无法自动处理的异常被明确报告，而不是静默跳过或删除数据。
18. 作为读者，我希望 feed 中的文章按照首次公开时间稳定排序，并且同一时间的文章不会随机跳动。
19. 作为搜索引擎，我希望 sitemap 使用文章内容或公开状态的更新时间，而不是浏览量、收藏数等互动计数变化时间。
20. 作为读者，我希望文章重新发布时仍保留首次公开时间，避免 feed 把旧文章伪装成新文章。
21. 作为站点管理员，我希望所有 canonical、OG、JSON-LD、feed、sitemap 和 robots URL 使用同一个显式规范站点地址。
22. 作为站点管理员，我希望生产环境缺少规范站点地址时启动失败，而不是生成错误的绝对 URL。
23. 作为读者，我希望文章没有封面或摘要时页面仍然有效，并安全地省略不可用的图片字段或生成合理的摘要回退。
24. 作为站点管理员，我希望文章标题、摘要、作者名、图片 URL 和站点设置中的特殊字符不会破坏 HTML、JSON-LD 或 XML。
25. 作为博客读者，我希望在 SPA 内从一篇文章切换到另一篇文章时浏览器 title 和 metadata 同步更新。
26. 作为博客读者，我希望离开文章详情页后，文章级 metadata 被清理并恢复为站点级默认值。
27. 作为前端应用，我希望浏览器 URL 使用 slug，而收藏、评论、编辑等内部 API 仍可使用数字文章 ID。
28. 作为公开 API 消费者，我希望通过专用的 slug 查询接口获取已发布文章，而不能通过查询参数绕过公开可见性过滤。
29. 作为 feed 阅读器，我希望 RSS 和 Atom 使用相同的文章集合、排序、摘要、作者和时间语义，仅在协议外层格式上有所区别。
30. 作为爬虫，我希望发现端点在数据库故障时返回明确的 5xx，而不是将故障伪装成空 feed 或空 sitemap。
31. 作为缓存代理，我希望成功的 feed 和 sitemap 响应具有短时公共缓存，并支持 ETag 或 Last-Modified 条件请求。
32. 作为站点管理员，我希望发现端点的成功和失败缓存策略不会把私密内容或错误结果公共缓存。
33. 作为维护者，我希望在文章数量超过单个 sitemap 的协议上限时，sitemap 能够以稳定方式分片而不丢失公开文章。
34. 作为开发者，我希望开发环境能够通过 Vite 代理验证与生产环境一致的根级 feed、sitemap 和 robots 端点。
35. 作为发布负责人，我希望可以先执行数据库迁移和数据校验，再部署依赖 slug 的应用版本，以降低上线风险。
36. 作为发布负责人，我希望上线验收能够验证旧 URL 301、新 slug URL 200、未公开文章 404、四个发现端点和初始 HTML metadata。
37. 作为维护者，我希望自动化测试覆盖公开过滤、URL 重定向、XML 合法性、特殊字符转义、缓存头和客户端 metadata 切换。

## Implementation Decisions

### 公开 URL 与 slug

- 文章公开详情 URL 使用 `/blog/<slug>`，slug 是单一路径段。
- slug 只允许小写 ASCII 字母、数字和连字符，最多 120 个 ASCII 字符，不以连字符开头或结尾，不包含斜杠、反斜杠、点号路径片段、控制字符或其他路径结构。
- 英文、数字和常见标点转换为 kebab-case；中文优先通过本地、固定版本的拼音能力转换为拼音。无法生成有效结果时使用稳定的 `article-<id>` 形式兜底。
- slug 在文章记录中持久化，在全体文章记录范围内唯一，不区分草稿、私密、已发布或软删除状态。
- 软删除文章继续占用原 slug；新文章不得复用该 slug。
- 创建文章时自动生成 slug。普通标题编辑不改变 slug；P0 不提供已发布文章 slug 编辑能力，也不提供 slug 历史表。
- 相同标题或其他冲突使用确定性数字后缀解决，结果必须保持长度限制。已有非空 slug 不被自动回填流程重写。

### 数据迁移与发布时间

- 文章模型增加持久化 slug，并建立全局唯一约束。由于唯一约束属于破坏性 schema 变更，必须提供可审查的 MySQL 迁移脚本，不能只依赖自动迁移在已有数据上直接创建约束。
- 迁移先增加字段，再按稳定顺序回填历史文章，处理冲突并校验非空和唯一性，最后建立唯一索引。迁移失败时停止并报告，不删除记录、不覆盖已有非空 slug、不静默跳过异常。
- 回填脚本和应用启动兜底检查必须幂等；多实例启动不得重复执行大批量回填。
- 历史文章缺失或非法的 `publish_at` 使用 `created_at` 回填；有效的历史发布时间保留。回填应提供统计结果以便部署核对。
- `publish_at` 表示文章首次公开时间。文章从非 published 状态首次转为 published 时写入当前时间；取消发布或重新发布不清空、不重置该值。
- `updated_at` 表示文章内容或公开状态的最后更新时间。浏览量、收藏数、评论数等互动计数变化不得影响 sitemap 的 `lastmod`。

### 公开可见性与 API

- 公开 SEO HTML 路由、slug 查询 API、feed 和 sitemap 只选择 published 且未软删除的文章。
- 草稿、私密、软删除和不存在的文章在公开 HTML 路由中统一返回 404；登录状态不改变 SEO HTML 路由的公开语义。
- 已发布文章的纯数字旧 URL `/blog/<id>` 由服务端返回 301，并在 `Location` 中指向当前 slug canonical URL。未公开或不存在的数字 ID 返回 404。
- 非数字的单路径段按 slug 查询；规范 URL 使用小写、无尾斜杠形式。可安全规范化的变体永久重定向到唯一 canonical URL。
- 增加专用公开 slug 查询接口，接口只按 slug 查找已发布且未软删除文章，不接受 `status`、`show_all`、`author_id` 等可绕过公开过滤的参数。
- 文章响应增加 slug，同时保留数字 ID，以支持公开链接使用 slug、内部收藏/评论/编辑等操作继续使用 ID。
- 既有按数字 ID 的内部 API 不因浏览器 URL 改为 slug 而改变。

### 服务端文章 metadata

- 只对博客文章详情 URL 做服务端增强；博客列表页和其他 SPA 路由继续使用普通应用 fallback。
- 服务端返回文章详情初始 HTML 时注入：
  - `<title>`，格式为“文章标题 | 站点名称”；
  - description 和 `og:description`；
  - canonical link；
  - `og:title`、`og:description`、`og:image`、`og:url`；
  - Twitter/X card、title、description 和可用图片；
  - `BlogPosting` JSON-LD，包括 headline、description、url、image（如可用）、datePublished、dateModified、author 和 publisher。
- metadata 只注入文章级字段和 JSON-LD，不在 P0 完整服务端渲染正文；正文、Markdown 渲染和交互继续由 Vue 客户端负责。
- HTML 使用明确的 head 占位区域进行安全替换。标题、摘要、作者、站点设置和图片等用户可影响的数据必须按 HTML、属性和 JSON-LD 上下文分别编码，禁止手工拼接未经转义的内容。
- 摘要优先使用文章 summary；为空时从正文提取纯文本，去除 Markdown、HTML、脚本、样式和控制字符并安全截断到约 160 个字符；仍为空时使用站点描述。
- 封面、站点 logo 和作者 URL 只有在能够解析为公开、安全的绝对 URL 时才输出。站内绝对路径可拼接规范站点地址；非法协议、空值或不可公开资源直接省略。
- 作者只输出公开名称和可确认公开的用户 URL；不输出邮箱、手机号、权限或其他私有资料。publisher 使用站点公开名称，logo 缺失时省略 logo 字段。

### 前端文章路由与 metadata

- 博客前端文章路由参数改为 slug；首次加载通过公开 slug API 获取文章数据，返回的数字 ID继续用于内部操作。
- 前端使用轻量、可测试的 metadata helper，在客户端路由切换时更新 title、canonical、OG、Twitter 和 JSON-LD，并清理前一篇文章的 metadata。
- 客户端 metadata 规则与服务端 metadata 规则保持一致。客户端 metadata 是 SPA 导航补充，不替代服务端初始 HTML metadata。
- 站点设置异步加载时不得覆盖已经设置的文章级 title；离开文章页后恢复站点级默认 metadata。
- 通用 HTML head 声明 RSS 和 Atom alternate link，使用规范站点地址。

### Feed、Sitemap 与 robots

- 提供以下根级端点：
  - `/feed.xml`：RSS 2.0；
  - `/atom.xml`：Atom 1.0；
  - `/sitemap.xml`：公开页面和文章 sitemap；
  - `/robots.txt`：爬虫规则和 sitemap 指引。
- RSS 与 Atom 共享同一份最近公开文章查询结果：最多 50 篇，按 `publish_at DESC, id DESC` 稳定排序，只包含 published、未软删除文章。
- Feed 只输出安全的纯文本摘要，不输出全文。缺失 summary 时使用与页面一致的正文纯文本回退。标题、摘要、作者和 URL 必须使用标准 XML 编码。
- RSS `guid` 和 Atom `<id>` 使用文章 canonical URL；文章链接使用 slug URL。RSS 使用首次发布时间作为 `pubDate`；Atom 使用首次发布时间作为 `published`、使用内容更新时间作为 `updated`。
- RSS 声明自身链接，Atom 声明自身链接；站点 HTML 同时声明 RSS 和 Atom alternate link。
- Sitemap 包含首页 `/`、博客列表 `/blog` 和全部公开文章详情 URL。首页和博客列表没有可靠更新时间时省略 `lastmod`；文章 `lastmod` 优先使用 `updated_at`，再回退到 `publish_at`，最后回退到 `created_at`。
- Sitemap 不因浏览量等互动计数变化更新 lastmod。文章不再公开时立即从 sitemap 中移除。
- sitemap 在未超过协议限制时返回 URL 集合；超过单个 sitemap 的约 50,000 个 URL 上限时返回 sitemap index 和稳定分片。分片继续应用相同的公开过滤和缓存规则。
- robots.txt 允许公开内容和发现端点，禁止 `/api/`、后台和用户私有路径，并指向规范站点的 `/sitemap.xml`。robots.txt 不替代应用鉴权。
- 发现端点使用固定 MIME 类型：RSS 为 `application/rss+xml`，Atom 为 `application/atom+xml`，sitemap 为 `application/xml`，robots 为 `text/plain`，均包含 UTF-8 charset。
- 发现端点支持 GET 和 HEAD。HEAD 与 GET 使用相同状态和主要响应头但不返回正文；XML/TXT 请求不得落入 SPA fallback。

### 规范站点地址与缓存

- 所有 canonical、OG、JSON-LD、feed、sitemap、robots 和文章链接使用集中配置 `PUBLIC_SITE_URL`。
- 规范站点地址必须是带 host 的 `http` 或 `https` URL，不含 query 或 fragment，并规范化尾斜杠。
- 生产环境缺少或配置非法时拒绝启动；开发和测试环境必须显式提供本地或测试地址，不信任任意请求 Host 头生成生产 URL。
- 成功生成的 feed 和 sitemap 使用约 5 分钟公共缓存，可通过 ETag 和/或 Last-Modified 支持条件请求；无变化时返回 304。
- 发现端点查询失败返回 5xx，并使用 `Cache-Control: no-store`；不得把数据库故障伪装成空 XML，也不得缓存错误响应。
- 文章详情 metadata 使用 no-cache 或极短缓存，避免长期锁定旧 metadata。
- 开发环境 Vite 代理应转发四个根级发现端点，使本地验证与生产路由保持一致。

## Testing Decisions

- 最高测试 seam 是博客服务的公开 HTTP 契约。测试通过真实路由验证 slug 文章 200、数字 URL 301、未公开文章 404、初始 HTML metadata、四个发现端点、公开过滤、内容类型、HEAD 行为、缓存头、条件请求和故障 5xx。
- HTML 测试解析 title、description、canonical、OG、Twitter metadata 和 JSON-LD，而不是只检查字符串存在；覆盖标题、摘要、作者、站点设置和 URL 中的引号、尖括号、&、脚本片段、中文和无效图片值。
- Feed 和 sitemap 测试使用 XML 解析器验证结构、规范 URL、发布时间、排序、摘要、self link、公开过滤和 sitemap 分片；草稿、私密、软删除和作者可见但公众不可见文章必须被排除。
- 公开 HTTP 测试必须验证现有列表查询参数不能绕过 feed、sitemap 或 slug API 的 published 过滤。
- slug 生成和冲突分配使用表驱动测试，覆盖中文、英文、数字、特殊字符、空结果、超长结果、相同标题、软删除占用和幂等回填。
- 使用隔离 MySQL 验证 schema 迁移、历史 slug 回填、唯一约束、publish_at 回填、已有非空 slug 保留和迁移异常停止；不得使用开发数据库作为测试数据。
- 前端测试通过真实文章详情路由和 mock API 验证 slug 加载、SPA 内文章切换、title/head 更新、JSON-LD 去重、无封面/空摘要边界以及离开文章页后的 metadata 恢复。
- 测试优先验证公共行为，不依赖 handler/service 的内部调用顺序、缓存 key 或具体模板实现。可复用现有 Go HTTP、SPA fallback、文章详情和博客前端测试模式。
- P0 验收前运行相关 Go 测试、完整 Go 构建/静态检查、博客前端测试、构建和 lint；具备隔离服务时再运行真实数据库/缓存集成验证。

## Out of Scope

- 不在 P0 对作品、照片、Wiki、用户页等其他公开内容建立 feed、sitemap 或文章级 JSON-LD；架构应保持未来可扩展。
- 不迁移到 Nuxt、完整 Vue SSR、同构渲染或独立 Node 运行时。
- 不在初始 HTML 中完整渲染 Markdown 正文，不改变现有正文渲染器和文章 API 的正文语义。
- 不在 P0 提供已发布文章 slug 编辑、slug 历史表或旧 slug 多次变更兼容；在历史表上线前禁止修改已发布 slug。
- 不为文章生成新的图片服务或社交分享图片；没有合法封面时省略图片 metadata。
- 不输出 RSS/Atom 全文，只输出安全纯文本摘要。
- 不把 robots.txt 当作权限控制，也不允许通过 robots 规则暴露未公开内容。
- 不为发现端点建立长期缓存或在数据库故障时返回伪造的空结果。
- 不改变现有数字 ID API 的内部使用方式，不把所有 API 一并迁移到 slug。
- 不新增文章互动统计语义，不让浏览、收藏、评论等计数成为 SEO 更新时间。

## Further Notes

- 当前项目是 Go/Gin 后端加 Vue 3 客户端 SPA，生产部署将前端构建产物嵌入 Go 服务；因此本方案选择 Go 注入文章 head metadata，而不是引入新的 SSR 运行时。
- 现有文章模型已有 `publish_at`，但发布流程和对外响应需要补齐首次公开时间语义；历史数据回填必须在应用发布前完成。
- 数据库采用 MySQL，schema 通过 GORM 管理，但新增全局唯一 slug 约束和历史数据回填必须使用可审查迁移脚本，并遵循先迁移、回填、校验，再部署应用的上线顺序。
- 上线验收至少包括：公开文章 slug 200、旧数字地址 301、未公开文章 404、初始 HTML SEO metadata、RSS/Atom/Sitemap/robots 正确响应、条件缓存、Vite 开发代理和生产配置校验。
- 上线后应监控数字 URL 301/404、slug URL 200/404、四个发现端点 4xx/5xx、XML 生成失败、slug 冲突、数据库查询延迟和缓存命中情况。
- 本文是需求规格，不代表实现已经开始；下一步应由维护者审核本 spec，再使用 `/to-tickets` 按可独立验证的纵向行为拆分实施票。
