# Article SEO final-review resolution

Date: 2026-11-05. Worktree: `/Users/iceymoss/project/inkspace-seo-review`.
Branch: `fix/seo-review-findings`; reviewed implementation/base: `64138eb3`.
Scope: the four SPEC findings, Standards D1–D3, and minimal H1–H2 cleanup.
The opening record is historical implementation/local verification, **not deployment acceptance**.
Real isolated-MySQL acceptance was completed in the [2026-11-06 follow-up](#real-mysql-acceptance-resolution--2026-11-06)
below; the earlier SKIP statements describe the original review run only.

## Finding → implementation → regression evidence

| Finding | Resolution | Regression / evidence |
| --- | --- | --- |
| SPEC 1 / P1: terminal interaction reload loop | `web/blog/src/views/BlogDetail.vue` watches a stable resolved article ID, independent of temporary UI clearing. Route loading resets that ID; resolving a route does not replay an old refresh signal. | `BlogDetail.test.js`: actual Pinia terminal store + executor confirmed favorite/unfavorite + real Vue router + delayed slug responses. Each interaction adds one request; resolving it adds none. Leaving/revisiting fetches once and retains canonical/title/single JSON-LD. Test demonstrated red (`expected 3 to be 2`) before the fix, then green. |
| SPEC 2 / P2: browser numeric IPv4 resource hosts | `internal/seo/url.go` rejects noncanonical browser numeric-ending hosts rather than partially reimplementing WHATWG IPv4 normalization. Existing IP range checks and canonical public-IP/domain support remain. No network request or DNS lookup. | `TestBrowserNumericResourceHostsAreOmitted` in `internal/seo/metadata_test.go`: hex/mixed/octal/short/single-number/trailing-dot and invalid numeric-ending DNS variants; helper and actual cover/logo/author metadata projections. Red before filtering, green after. Ordinary domains, hex-looking nonfinal labels and canonical public IPv4 remain allowed. |
| SPEC 3 / P2: settings outage breaks generic SPA | `internal/router/blog.go` uses `service.PublicSiteIdentityFromSettings(base, nil)` only for generic-shell settings failures. It still requires valid explicit canonical config and emits absolute feeds/bootstrap. Article/feed/sitemap errors remain strict. | `TestGenericSPASettingsFailureKeepsShellButSEORemainsStrict`: actual assembled router, GET/HEAD for `/`, `/blog`, `/login`, `/dashboard`, `/share/token`, `/index.html`; default shell/absolute feeds/bootstrap and no-cache. Article/RSS/Atom/sitemap retain 500/no-store. Missing canonical config is not silently accepted. Red 500s before fix, green after. |
| SPEC 4 / P2: legacy automatic timestamps change lastmod | `internal/service/article_counter.go` supplies counter delta plus explicit `updated_at = updated_at`. Applied to public/legacy views, favorite/unfavorite, like/unlike, comment create/delete/moderation. `internal/database/migrate.go` also preserves timestamps in comment helpers and both article `SyncCounters` statements. Content/state writes are unchanged. | Connection-free `TestArticleViewStatementsPreserveLegacyAutomaticTimestamp` observes actual service/helper SQL; red without self-assignment, green with it. Fresh and legacy gated interaction fixtures exercise changed counters, publication/content time and sitemap lastmod after every operation, including moderation/helpers/sync. **MySQL assertions SKIPPED, not proven here.** |
| D1: incomplete centralized columns | Article counter columns added to `ArticleColumns`; `User.TableName` and privacy-minimal `UserColumns` projection added. Live slug occupancy uses mapped, unscoped GORM; narrow create-only slug allocation SQL uses trusted mapped identifiers. No repository/DI layer or broader model migration. | Existing slug/public-collection suites and full Go tests pass. Real allocation/uniqueness/projection SQL remains in the gated MySQL suite. |
| D2: nonrepeatable SQL artifact | `scripts/article-slugs/schema.sql` is guarded preparation only; new `finalize.sql` validates repaired rows and both existing index shapes before installing anything, unique index first and completion marker last. Both require an explicit same-session offline acknowledgement. Go command also validates slug constraint before marking publication ready, and checks marker column order/direction. Updated migration README and rollout instructions. | `TestArticleSQLArtifactIsGuardedAndPhased` passes. Gated `TestArticleSQLPhasesMySQL` runs prepare/finalize twice; rejects absent opt-in, pre-repair finalization, missing/invalid/duplicate slugs, zero/partial/impossible dates, nonunique slug index, malformed marker; preserves populated `reserved` slug. It also checks Go rejection of malformed indexes/marker ordering. **Artifact execution SKIPPED.** |
| D3: request log bypasses Zap | `internal/handler/public_article_handler.go` uses configured global Zap with `article_id` and error fields. View-write failure remains nonfatal and private. | `TestPublicSlugAPIViewFailureUsesStructuredNonfatalLog` demonstrated red with no observed Zap entry, then green; asserts 200 response without private error and exactly one structured warning. |
| H1: router middle-man chain (heuristic) | Production `SetupUserRouter` directly invokes the one real assembly function. Convenience reader constructors moved to `internal/router/blog_test_helpers_test.go`. | Existing real-router HTTP acceptance tests continue to pass; no new runtime wrapper/bundle/framework. |
| H2: repeated article-reference parsing (heuristic) | `parseArticleReference` in `virtualFileSystem.js` is shared by navigation, terminal reads, and confirmed writes. Preserves slug versus legacy numeric ID and rejects malformed nested article paths. | Parser table and executor tests cover numeric-prefix slugs, legacy numeric files, safe integer/length/path/extension boundaries and confirmed ID-based favorite. New legacy/malformed cases demonstrated red before shared parsing and green after. |

The existing mandatory explicit `PUBLIC_SITE_URL` contract is unchanged in every mode;
Host headers are not trusted, no localhost default was introduced, and strict SEO
article/feed metadata query errors still produce 500/no-store.

## Verification actually executed

All commands run in this worktree (frontend commands in `web/blog`):

| Command | Result |
| --- | --- |
| `go test ./... -count=1` | PASS; integration gates skip without credentials |
| `go build ./...` | PASS |
| `go vet ./...` | PASS |
| `go test -race ./internal/config ./internal/database ./internal/models ./internal/service ./internal/handler ./internal/router ./internal/seo -count=1` | PASS |
| `go test ./internal/webassets/blog -run TestEmbeddedBlogHeadPlaceholders -v -count=1` | PASS |
| `npx --yes pnpm@9 install --frozen-lockfile` | PASS; no dependency/lock changes |
| `npx --yes pnpm@9 run test --run` | PASS: **16 files / 114 tests** |
| `npx --yes pnpm@9 run build` | PASS; ignored embedded dist rebuilt |
| `npx --yes pnpm@9 run lint` | PASS |
| `python3 scripts/seo_smoke_test.py` | PASS; disposable HTTP/curl happy path and intentional rejection fixture only |
| `git diff --check` | PASS |

Nonblocking existing warnings: unreadable user `.npmrc`, Vite CJS Node API deprecation,
and large frontend chunks. Generated dist and dependencies are not committed.

## Historical MySQL gate — 15 explicit SEO SKIPs, not PASS (2026-11-05)

Executed:

```sh
go test ./internal/database ./internal/service ./internal/router \
  -run 'MySQL|Publication|PublicArticleCollections|Interactions' -v -count=1
```

The 15 SEO fixtures skipped because `INKSPACE_SEO_TEST_MYSQL_DSN` is absent:
three publication migration/readiness/preflight tests; two slug migration/conflict tests;
SQL phase artifact test; lifecycle, concurrent first publication, historical republish;
fresh and legacy interaction timestamps; public collections; permalink HTTP; feed HTTP;
sitemap HTTP. The unrelated theme integration test also skipped its independent gate.
No developer database, live Redis, migration credentials or live DDL was accessed.

At that point, an operator still needed to supply an explicitly authorized isolated **MySQL 8.0**
DSN with CREATE/DROP test-schema privileges and **no selected database**. The fixture
creates/drops only its random `inkspace_seo_test_*` schema. SQL artifact testing also
requires routine privileges. Run the command above and record actual results.

The SQL fixture separates pinned prepare/finalize connections from the Go command's
own advisory-lock connection and opts in separately for each SQL session. The legacy
interaction fixture performs all association AutoMigrate work before installing
`ON UPDATE CURRENT_TIMESTAMP(3)`, asserts that metadata, then performs no more DDL.
These safeguards were then only reviewed/compiled; the follow-up below records their execution.

Container/live gateway/TLS/CDN acceptance, production backup/rollback rehearsal and
the offline migration remain operator gates in [SEO_ROLLOUT.md](SEO_ROLLOUT.md).
No deploy, push, PR, production migration or subagent work was performed.

## Real MySQL acceptance resolution — 2026-11-06

Branch `fix/seo-mysql-acceptance`, worktree `/Users/iceymoss/project/inkspace-seo-mysql`,
base `1d2ecf80`. The authorized existing server reports **MySQL 8.0.33**. Credentials were read
programmatically, never printed/committed; the environment DSN selected no database. Only fixture-owned
random `inkspace_seo_test_*` schemas were created/changed/dropped, not any business schema.

| Acceptance failure | Root cause and minimal resolution | Real regression evidence |
| --- | --- | --- |
| SQL phases: error 1644 after offline opt-in | GORM's initialized `Connection` handle retained the deliberately triggered guard error. Use `Session(NewDB: true)` statement clones on each pinned phase connection. SQL artifacts/guards remain unchanged. | Prepare/finalize reruns, guard and invalid row/index cases all execute and pass. Added explicit unset acknowledgements and no-change checks, including the pooled-session case. |
| Fresh + legacy interactions: `SyncCounters` error 1093 | MySQL rejects the correlated `comments` self-update. Join a grouped/materialized reply-count projection and `COALESCE` missing groups to zero; preserve original undeleted/direct-reply semantics. | Both existing fixtures demonstrated red before the query fix, then green. Added active/rejected/deleted reply assertions plus stale-to-zero reset. Existing changed article counters preserve `updated_at`, first publication and sitemap lastmod, including verified legacy `ON UPDATE`. |
| Permalink `/` expected 500 despite correct shell 200 | Stale gated assertion contradicted SPEC 3's reviewed SPA fallback. Separate strict settings-API failure from generic-shell success. | Real dropped settings table: API 500/no-store; shell GET/HEAD 200/no-cache, empty HEAD, default title, canonical bootstrap and absolute RSS/Atom links. Existing article/feed/sitemap 500/no-store checks remain. |

Verification (SEO DSN injected through a local subprocess launcher):

- Original targeted database/service/router command: **15 SEO PASS, 0 SEO SKIP/FAIL**;
  the unrelated theme test alone skips its independent gate.
- `go test ./... -json -count=1`: **124 top-level PASS / 412 including subtests**, **28 unrelated
  SKIP**, **0 FAIL**, **10 test packages PASS**. Includes **16 SEO MySQL PASS**, adding model slug allocation.
- `go test -race ./... -json -count=1`: identical pass/skip counts; all **16 SEO MySQL PASS**,
  no race reported.
- `go build ./...`, `go vet ./...`, `git diff --check`: **PASS**.

See [SEO_VALIDATION.md](SEO_VALIDATION.md#real-mysql-acceptance-follow-up--2026-11-06) for complete
scope and limits. The 28 unrelated bootstrap/admin, reading-time/MySQL/Redis and theme integration
gates remain skipped, not passing claims. Cache invalidation uses the existing RESP stub, not live
Redis. Frontend code was unaffected; its earlier **114-test** evidence remains historical.
The isolated SEO MySQL gate is resolved; business-data `verify.sql`, offline migration, backup/rollback,
container/gateway/TLS/CDN and live deployment acceptance remain operator responsibilities.
