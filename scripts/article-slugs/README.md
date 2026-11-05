# Article slug migration (MySQL only)

Back up the database, review this directory and `internal/database/article_slugs.go`,
then **stop every article writer** (blog/admin/scheduler and old application replicas).
Do not deploy slug-dependent application code before this finishes.

The explicit command performs the reviewed DDL in `schema.sql`, validates all
existing nonempty slugs (including deleted/private/draft rows), backfills NULL/empty
slugs in ascending ID order, validates completeness, then creates the global unique
index. Existing nonempty values are never rewritten. Invalid/duplicate existing
values fail explicitly before backfill; resolve those manually rather than deleting
articles. Updates do not change `updated_at`. Re-running preserves previous results.
The command shares the startup advisory lock; it is NOT a live traffic migration.

```sh
export INKSPACE_ARTICLE_SLUG_MIGRATION_DSN='operator:password@tcp(host:3306)/inkspace?charset=utf8mb4&parseTime=true'
export INKSPACE_ARTICLE_SLUG_MIGRATION_OFFLINE=1
go run ./cmd/article-slugs
```

Review `scanned` / `filled`, then rerun (expect `filled=0`) and verify:

```sql
SELECT COUNT(*) FROM articles WHERE slug IS NULL OR slug = '';
SELECT slug, COUNT(*) FROM articles GROUP BY slug HAVING COUNT(*) > 1;
SHOW INDEX FROM articles WHERE Key_name = 'ux_articles_slug';
```

New installations create the unique index only on their newly created empty table.
Existing installations refuse startup without it; startup never performs bulk
backfill. The nullable column permits a transactional insert-before-ID-fallback
allocation; committed application-created rows always have a slug. Do not write
articles using ad-hoc SQL. Rollback application code without dropping slugs/index:
dropping either loses the permanent URL reservation contract.

Isolated tests: set `INKSPACE_SEO_TEST_MYSQL_DSN` to a MySQL account allowed to
CREATE/DROP **test databases** with no selected database (not developer credentials).
Tests create a random `inkspace_seo_test_*` database and drop only that database.
Without this explicit opt-in, integration tests skip. No Redis is required.

Publication-time backfill belongs to ticket 02 and must also finish before rollout.
