# Public article discovery read contracts

These service reads accept `context.Context`, preserve database/cancellation
errors, never count views, and never use the internal article/list/Redis cache.
Visibility is always `status=1 AND deleted_at IS NULL`, independent of identity
or `status`, `show_all`, `author_id`, ranking and sort HTTP query parameters.
Do not use `ArticleService.GetList` for discovery: it deliberately supports
internal/admin filtering.

```go
(*service.ArticleService).GetPublicBySlug(ctx context.Context, slug string) (*models.Article, error)
(*service.ArticleService).GetPublicByID(ctx context.Context, id uint) (*models.Article, error)
(*service.ArticleService).GetRecentPublicArticles(ctx context.Context) ([]*models.Article, error)
(*service.ArticleService).GetPublicSitemapArticles(ctx context.Context) ([]service.PublicSitemapArticle, error)
```

- Detail reads return `gorm.ErrRecordNotFound` for draft/private/deleted/missing.
  `ArticleResponse.publish_at` exposes the first publication time (nullable).
- **Recent/feed**: fixed maximum 50; `publish_at DESC, id DESC`; selected
  ID/slug/title/summary/content/author_id/created_at/updated_at/publish_at fields.
  Content is present only for shared plain-text excerpt fallback. `Author` loads
  only ID, username and nickname (no password/email/role). Reuse this same query
  for RSS and Atom. It returns article models as a projection, not complete rows.
- **Sitemap**: `[]PublicSitemapArticle` contains only ID, slug, created_at,
  updated_at, publish_at; ID ASC. One SELECT gives matching entries and count
  (`len(entries)`) from the same statement snapshot. Consumers split this compact
  snapshot into protocol-sized shards; do not add content/author preloads. Memory
  use is O(public article count), not O(total article body size).
- `PublicSitemapArticle.LastModified()` prefers nonzero updated_at, then nonzero
  publish_at, then created_at. Zero means an unavailable time, not today's date.
- Separate HTTP requests intentionally see fresh visibility. ID ordering is stable
  under edits but not an immutable cross-request version: publication/deletion may
  shift offset-based shard membership. Downstream must document this consistency
  boundary rather than claim independent requests form a frozen snapshot.

Before rollout run the extended explicit offline migration in
[`scripts/article-slugs/README.md`](../../scripts/article-slugs/README.md).
First-publication changes are serialized by the Update transaction's article row
lock. Existing valid first-publication values survive content edits, unpublishing,
and republishing. Interaction services' `UpdateColumn` writes preserve updated_at;
sitemap lastmod must not be based on interaction counters.
