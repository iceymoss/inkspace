# 05: Sitemap 与 robots 发现能力

**What to build:** 为搜索引擎提供包含首页、博客入口和全部公开文章的 sitemap，以及明确公开/私有抓取边界的 robots.txt。

**Blocked by:** 02: 文章发布时间与公开文章集合

**Status:** resolved

- [x] `/sitemap.xml` 返回合法 sitemap，包含 `/`、`/blog` 和全部 published 且未软删除文章的 slug URL。
- [x] 文章 `lastmod` 优先使用内容/公开状态更新时间，再回退到首次发布时间或历史创建时间；首页和博客列表没有可靠更新时间时不伪造 lastmod。
- [x] sitemap 不包含草稿、私密、软删除或其他不可公开文章，互动计数变化不会导致伪造的 lastmod。
- [x] 公开 URL 数量超过单个 sitemap 的约 50,000 条协议上限时，端点返回 sitemap index 和稳定分片，且分片不会遗漏或重复公开文章。
- [x] `/robots.txt` 允许公开文章、feed、sitemap 和静态公开资源，禁止 API、后台和用户私有路径，并指向规范站点 sitemap。
- [x] sitemap 和 robots 支持 GET/HEAD，使用正确的 UTF-8 MIME 类型；成功响应使用短时公共缓存，数据库失败返回 5xx 且 no-store。
- [x] 覆盖 XML 解析、公开过滤、lastmod、分片边界、重复/遗漏、robots 规则、MIME、GET/HEAD、缓存和数据库故障测试。

## Answer

Added standard-library XML URL set/index rendering and root sitemap/robots HTTP
routes. Reused ticket02's single fresh compact ID-ordered public query and
`LastModified`, ticket03's canonical origin/URL helpers, and ticket04's shared
ETag/GET/HEAD/error responder. Canonical URLs never trust request Host. Invalid
persisted slugs fail closed with 500/no-store, including conditional requests.

The exact 50,000 URL limit counts **both** static pages: 49,998 articles still fit
one URL set; 49,999 articles produce an index with `/sitemaps/1.xml` and
`/sitemaps/2.xml`. Shard 1 contains both static pages once. Count and membership
come from one query, not separate count/offset queries. For an unchanged public
dataset, shards cover every URL once in stable ID order. Edits preserve membership;
publication/deletion between independent HTTP requests may shift membership or
shrink an index to a URL set. This is not a frozen cross-request snapshot and uses
O(public article count) compact memory. Out-of-range/malformed shards are empty
404/no-store, never SPA HTML.

Robots publishes only generic route rules: `/api/` (including admin API), actual
private `/dashboard`, `/favorites`, `/profile`, `/login`, and private knowledge
storage are disallowed. `/users/*` and `/wiki/*` stay crawlable. The admin SPA is a
separate origin, not a fabricated blog `/admin` path; robots does not replace auth.
Robots needs validated config only, no private settings or DB query.

## Validation evidence

- Pure XML URLset/index and robots tests were red before implementation, then green.
- Real router tests cover exact 0/49,997/49,998/49,999/50,000/99,998/99,999 article
  boundaries, sparse IDs, union ordering/no duplicate/missing URLs, stable shard
  routes, invalid/overflow shards, lastmod fallbacks, hostile Host, ignored list
  filters, MIME, GET/HEAD, public max-age=300, ETag weak/list/wildcard/304, changed
  content, fresh unpublication, index-to-urlset transition, read failure and bad
  canonical config. Invalid persisted slug test demonstrated red 304 then green 500.
- After merging shared dependencies (ticket03 integration `44ab3186`, ticket04
  integration `f9fdc9e5`; own merge `1349cef4` retains feed and sitemap wiring), passed
  `go test ./... -count=1`, `go build ./...`, `go vet ./...`, and
  `go test -race ./internal/router ./internal/seo -count=1`.
- `TestSitemapHTTPMySQL` compiles but **SKIPPED**: no authorized
  `INKSPACE_SEO_TEST_MYSQL_DSN`. Its random-schema fixture exercises production
  router/service public filtering against draft/private/deleted rows, author/admin
  token and list-filter bypass attempts, interaction lastmod stability, immediate
  unpublish and dropped-table errors. No actual DB validation or developer DB use
  is claimed. No Redis needed.
- No frontend source changes; Vite shard proxy/deployment propagation belongs to 06.
