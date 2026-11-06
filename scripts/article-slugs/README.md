# Article slug and first-publication migration (MySQL only)

Back up the database, review this directory and `internal/database/article_slugs.go`,
then **stop every article writer** (blog/admin/scheduler and old application replicas).
Do not deploy slug-dependent application code before this finishes.

The explicit command performs the reviewed preparation/constraint DDL, validates all
existing nonempty slugs (including deleted/private/draft rows), backfills NULL/empty
slugs in ascending ID order, validates completeness, then creates the global unique
index. Existing nonempty values are never rewritten. Invalid/duplicate existing
values fail explicitly before backfill; resolve those manually rather than deleting
articles. The command also backfills publication dates as described below.
Updates do not change `updated_at`. Re-running preserves previous results.
The command shares the startup advisory lock; it is NOT a live traffic migration.

## Repeatable SQL phases (MySQL 8.0)

The executable SQL artifacts are **separate phases**, not a concatenated migration:

1. [`schema.sql`](schema.sql) prepares only missing columns. It refuses unexpected
   existing types; it never rewrites populated slugs or installs readiness indexes.
2. Run the Go command below for deterministic local-pinyin slug allocation and strict
   textual calendar-date preflight/repair. SQL alone is **not** a pinyin/date repair tool.
3. [`finalize.sql`](finalize.sql) independently validates all slugs/dates and both
   existing index shapes, installs the unique constraint if missing, and installs the
   publication completion marker **last**. Invalid/missing/duplicate slugs, unrepaired
   dates, or malformed existing constraints stop with SQLSTATE 45000, never deletion
   or automatic replacement of populated slugs/constraints. It may also verify a fully
   completed Go migration and can be sourced repeatedly without duplicate DDL.

For SQL phases, use a reviewed operator connection with the intended database selected,
stop all writers and explicitly `SET @inkspace_article_seo_offline = 1;` **on that same
connection**, then `SOURCE scripts/article-slugs/schema.sql` or `finalize.sql` as appropriate.
Do not use mysql `--force`. These scripts require routine privileges and MySQL 8.0
(`REGEXP_LIKE`, descending indexes); DDL implicitly commits, so they are not atomic.
The Go command remains the recommended orchestration (including the shared initialization
lock); SQL opt-in is an acknowledgement, not a traffic/advisory lock. A failed phase
must be investigated with writers stopped and then rerun. No script here was executed
against live/developer data during implementation.

```sh
export INKSPACE_ARTICLE_SLUG_MIGRATION_DSN='operator:password@tcp(host:3306)/inkspace?charset=utf8mb4&parseTime=true'
export INKSPACE_ARTICLE_SLUG_MIGRATION_OFFLINE=1
go run ./cmd/article-slugs
```

Review slug `scanned` / `filled` and publication `scanned` / `filled` / `preserved`,
then rerun (expect both `filled=0`) and verify:

```sql
SELECT COUNT(*) FROM articles WHERE slug IS NULL OR slug = '';
SELECT slug, COUNT(*) FROM articles GROUP BY slug HAVING COUNT(*) > 1;
SHOW INDEX FROM articles WHERE Key_name = 'ux_articles_slug';
SHOW INDEX FROM articles WHERE Key_name = 'idx_articles_first_publication_v1';
SELECT COUNT(*) FROM articles WHERE publish_at IS NULL;
```

New installations create the unique index only on their newly created empty table.
Existing installations refuse startup without the slug index and publication
completion index; startup never performs bulk backfill. The nullable slug column
permits a transactional insert-before-ID-fallback
allocation; committed application-created rows always have a slug. Do not write
articles using ad-hoc SQL. Rollback application code without dropping slugs/index:
dropping either loses the permanent URL reservation contract.

## First-publication repair and rollout

Run this **extended command even if the earlier slug migration already completed**.
If absent, it adds nullable `publish_at datetime(3)`. In ID-ascending batches it
examines **all articles**, including draft, private and deleted history:

- NULL, zero dates, partial-zero dates, years outside MySQL's normal DATETIME
  range, and impossible calendar dates are invalid; repair these from `created_at`.
- Existing calendar-valid values are preserved, even if before creation or in the
  future: this migration does not invent scheduling/chronology policy.
- An invalid `created_at` needed as fallback stops the migration with the article
  ID **before any publication row is repaired**. Correct it explicitly, rerun,
  and check counters. Existing nonempty slugs are likewise never overwritten.
- Textual reads detect legacy MySQL `ALLOW_INVALID_DATES` values without Go
  driver's date normalization. No global SQL-mode change is made by the command.
- Publication repair uses raw UPDATE with `updated_at = updated_at`, retaining
  the existing content/state time even on legacy MySQL ON UPDATE schemas.

After **both** publication/slug repair and successful unique-constraint validation,
the command installs nonunique
`idx_articles_first_publication_v1(status, deleted_at, publish_at DESC, id DESC)`.
This supports the recent-public query and is the migration's final completion marker;
**do not install it manually before repair**. Startup checks index metadata only,
not a table-wide date scan. New installations mark their empty article table ready.
If interrupted, keep writers stopped and rerun; completed valid values remain.
Do not remove the completion index while this application version is running.

New article creation sets `publish_at` only for published status. First transition
from draft/private sets it if absent; unpublishing/republishing retains it.
Historical nonpublic records repaired above intentionally retain their fallback
creation time on later publication, as required by the historical repair policy.
Content/state edits advance `updated_at`; views/likes/favorites/comments do not.
Counter-only writes explicitly assign `updated_at = updated_at`, including comment
moderation and database counter maintenance, because skipping GORM timestamps does
not suppress a legacy MySQL `ON UPDATE CURRENT_TIMESTAMP(3)` column.

## Isolated tests

Set `INKSPACE_SEO_TEST_MYSQL_DSN` to a MySQL account allowed to
CREATE/DROP **test databases** with no selected database (not developer credentials).
Tests create a random `inkspace_seo_test_*` database and drop only that database.
Without this explicit opt-in, integration tests skip. No Redis service is required:
the lifecycle tests use an in-process disposable RESP cache-invalidation stub while
all persistence, locking, visibility, ordering and migration assertions use MySQL.
`TestArticleInteractionsPreserveLegacyAutomaticTimeMySQL` changes `updated_at` to
the legacy automatic shape **only in that random schema** and verifies all interaction
paths plus sitemap lastmod. `TestArticleSQLPhasesMySQL` executes the SQL artifacts
twice, rejects pre-repair constraints, invalid slugs/dates/index shapes and preserves
populated slugs. Without the DSN these are explicit SKIPs, not SQL acceptance.

```sh
go test ./internal/database ./internal/service ./internal/router -run 'MySQL|Publication|PublicArticleCollections|Interactions' -v -count=1
```
