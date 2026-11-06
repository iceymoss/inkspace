# Article SEO implementation acceptance evidence

Date: 2026-11-05. Ticket: `.scratch/article-seo-discovery/issues/06-deployment-and-acceptance.md`.
Worktree: `/Users/iceymoss/project/inkspace-seo-06`; branch `feat/seo-ticket-06`;
integration base `34356ce42dfb926a2cc9162c74302ad53d91d4f9` (tickets 01–05 integrated).
This opening record is historical implementation evidence, **not a production deployment**.
The [2026-11-06 MySQL acceptance follow-up](#real-mysql-acceptance-follow-up--2026-11-06)
below supersedes the earlier skipped-MySQL gate; no business migration or deployment is claimed.

## Delivered deployment boundary

- Shared ticket03 canonical config remains the single source of truth. Both YAML examples and
  `env.example` explain direct-Go vs Vite origins and identical BLOG origin for admin/scheduler.
- All five application services in **both** Compose variants explicitly require/forward
  `PUBLIC_SITE_URL`; HTTPS overlay inherits it. No localhost production fallback or build-time origin.
- All three K8s base Deployments require an environment-owned ConfigMap key. Each development
  overlay generates its own hashed ConfigMap with the same development blog origin. Namespace is
  explicit in overlays so Kustomize updates references correctly; no dev canonical origin in base.
- Vite forwards `/feed.xml`, `/atom.xml`, `/sitemap.xml`, `/robots.txt`, `/sitemaps/`, and existing
  `/api` runtime canonical configuration. No `VITE_PUBLIC_SITE_URL` or Host-derived fallback.
- Admin startup no longer retries a valid config.yaml after an invalid admin.yaml. Missing-file
  fallback remains in shared InitWithFile, where canonical validation also runs.

## Automated checks actually executed

| Check | Result | Evidence / qualification |
| --- | --- | --- |
| `go test ./... -count=1` | PASS | All Go packages; gated persistence tests skip without explicit credentials |
| `go build ./...` | PASS | All binaries, including embedded blog assets built from current source |
| `go vet ./...` | PASS | All packages |
| `go test -race ./internal/config ./internal/router ./internal/seo -count=1` | PASS | Affected config/HTTP/pure rendering packages |
| `go test ./internal/webassets/blog -run TestEmbeddedBlogHeadPlaceholders -v -count=1` | PASS | Rebuilt dist and fallback preserve explicit template head markers |
| Blog `npx --yes pnpm@9 install --frozen-lockfile` | PASS | No lock/dependency changes |
| Blog `npx --yes pnpm@9 run test --run` | PASS | **16 files / 98 tests**, including seven new live Vite proxy tests |
| Blog `npx --yes pnpm@9 run build` | PASS | Rebuilt ignored embedded dist before Go verification |
| Blog `npx --yes pnpm@9 run lint` | PASS | Existing ESLint command, includes new proxy test |
| `python3 scripts/seo_smoke_test.py` | PASS | Disposable HTTP server + actual curl script, happy path and intentional wrong-canonical rejection, output privacy assertions |
| `go test ./internal/config -run 'TestRenderedKubernetes\|TestAdminFallback' -v -count=1` | PASS | All three development overlays rendered by installed kubectl, generated ConfigMap hash/reference + namespace checked; admin fallback validation |
| Compose `config --quiet` with explicit origin | PASS | Both base variants and each combined with HTTPS overlay; rendering only, no daemon needed |
| Compose config without origin / with empty origin | EXPECTED REJECTION | Exit 15, required `PUBLIC_SITE_URL` interpolation error, no localhost substitution |
| `git diff --check` | PASS | Owned changes only, no generated dist staged |

Existing nonblocking warnings: unreadable user `.npmrc` (EACCES), Vite CJS Node API deprecation,
large frontend chunks. No dependency/package-manager/framework migration performed.

## HTTP and startup coverage

`internal/router/seo_acceptance_test.go` exercises the **same assembled production router**, with
public readers injected at existing article/site boundaries (no mocked route implementation):

- One public fixture ties HTML canonical/OG/JSON-LD, plain summary, RSS/Atom entry URLs,
  sitemap/shard locations, robots sitemap pointer, and generic HTML alternate links together.
- Slug 200, legacy numeric GET/HEAD 301, unpublished/missing reader results 404, article HEAD,
  MIME/charset, content-length/ETag parity, weak conditional GET/HEAD 304, no SPA fallback.
- Host/X-Forwarded-Host/Forwarded do not affect absolute URL generation. Hostile title text is
  parsed safely as metadata/JSON-LD; private full body, account email, unsafe image protocol and
  internal read-error details are never emitted in public HTML/XML.
- Injected read failures return 500/no-store, never an empty successful document or matching 304;
  robots stays database-independent. These are **reader failures**, not proof of real DB execution.

`internal/config/deployment_startup_test.go` builds/runs **server/admin/scheduler** in empty temp
directories, all debug/test/release modes. Unset/empty URL, query/fragment including empty `?/#`,
and unsafe protocol terminate before DB/Redis initialization. An invalid admin.yaml paired with a
valid config.yaml demonstrated **red** (database initialization attempted) before retry removal,
then **green**. Subprocess environment strips operator DB/Redis credentials and uses closed local
test endpoints if configuration validation regresses. Loader tests also cover missing admin-file
fallback rejecting missing/query/fragment origins in all modes and successful normalized explicit URLs.

`web/blog/src/utils/discoveryProxy.test.js` uses a real ephemeral Vite server and HTTP backend;
four endpoints and `/sitemaps/2.xml` initially returned SPA HTML (**red**), then forwarded correctly
(**green**), including query strings and `/api/settings/public`. Existing Vue route/metadata tests
continue to cover article A→B, stale async requests, settings race, head cleanup and de-duplication.

Rendered K8s review caught a namespace mismatch before delivery: generator resources initially
landed outside `inkspace-backend`, leaving unresolved/unhashed Deployment references. Explicit
overlay namespace and an offline rendered-manifest regression test now verify all three apps.

## Historical isolated MySQL run: explicitly SKIP, not PASS (2026-11-05)

Executed:

```sh
go test ./internal/database ./internal/service ./internal/router \
  -run 'MySQL|Publication|PublicArticleCollections|Interactions' -v -count=1
```

All **13 SEO MySQL fixtures SKIPPED**, because `INKSPACE_SEO_TEST_MYSQL_DSN` is absent:
slug migration/conflicts; publication backfill/readiness/invalid-created-at; lifecycle/concurrent
first publication/historical republish/interaction timestamps/public collection; public permalink
cross-end HTTP; feeds HTTP; sitemap HTTP. The unrelated theme MySQL fixture also skipped its own gate.
No developer database, live Redis, migration DSN, credentials or existing cluster was accessed.

The extended `TestPublicPermalinksMySQL` compiles and is ready for a dedicated MySQL 8 account
with CREATE/DROP test-schema permission and **no selected database** in its DSN. Existing
`SEOTestMySQL` creates/drops only a random `inkspace_seo_test_*` schema. It now checks the same
production HTML/API/feed/sitemap/robots assembly, public/draft/private/deleted exclusion including
author/admin bypass parameters, hostile Host, runtime canonical overriding a conflicting public DB
setting, safe site logo projection, HEAD/304, fresh withdrawal, and real query errors after dropping
only fixture tables. Existing fixtures cover ordering/max50, unique reservation, migration reruns,
publication marker, interaction timestamps and sitemap limits. **These real-SQL assertions have
not executed here.** SQL `verify.sql` is reviewable but has not been run against MySQL.

## Remaining release gates / deliberate limits

1. **Completed for isolated fixtures on 2026-11-06:** the expanded 15-test SEO suite and model
   slug-allocation fixture ran against authorized random schemas on the existing MySQL instance.
   Deployment-specific data checks and offline migration remain operator gates; see the follow-up below.
2. Docker CLI can render Compose, but the local Docker daemon is unavailable (`docker info` fails).
   Container image builds/startup, real nginx/TLS/DNS/CDN routing and cluster apply were not performed.
3. Operator must run [SEO_ROLLOUT.md](SEO_ROLLOUT.md), including **extended migration rerun even
   after ticket01**, `idx_articles_first_publication_v1` verification, backup/rollback rehearsal,
   canary smoke and monitoring. No live canary endpoint/fixture was supplied, so curl smoke was
   validated on its disposable fixture only, not against a running deployed Go+MySQL service.
4. Vite discovery proxy is tested, but Vite's initial page HTML is not Go metadata HTML. Verify the
   built embedded service for crawler acceptance; the marker test confirms the current build artifact.
5. Public discovery cache max-age=300 allows a proxy's bounded stale withdrawal window; origin
   re-reads on every request before conditional validation. No in-process result cache/hit metrics added.
6. Sitemap compact projection is O(public rows), with ID-ordered shards but no persistent snapshot
   across independent requests. Visibility changes may shift membership; load/capacity testing remains
   environment-specific. Robots is guidance, not access control. No full SSR, slug history or new metrics platform.

All changes are confined to this assigned worktree. Original user worktree documentation is untouched.
No subagents, push, cluster mutation, production migration or destructive application-data action occurred.

## Final-review fixes — 2026-11-05

The preceding sections are historical ticket06 evidence. The final-review fixes are
on `fix/seo-review-findings`, based on `64138eb3`, in
`/Users/iceymoss/project/inkspace-seo-review`. See
[SEO_REVIEW_RESOLUTION.md](SEO_REVIEW_RESOLUTION.md) for the complete four SPEC,
D1–D3 and H1–H2 finding-to-code/test mapping and observed red/green evidence.

Revalidated after the fixes:

- `go test ./... -count=1`, `go build ./...`, `go vet ./...`: **PASS**.
- Race tests for config/database/models/service/handler/router/seo: **PASS**.
- Embedded head-placeholder test and `python3 scripts/seo_smoke_test.py`: **PASS**.
- Blog frozen install, tests (**16 files / 114 tests**), build and lint: **PASS**.
- `git diff --check`: **PASS**. No dependency changes or generated dist committed.

That historical gated database/service/router run reported **15 SEO MySQL SKIPs**, not
the earlier 13: the added fixtures cover legacy automatic-timestamp interactions
and executable guarded migration SQL phases. The unrelated theme gate also skips.
`INKSPACE_SEO_TEST_MYSQL_DSN` is absent; no database connection or DDL was performed.
At that point SQL repetition/validation and legacy `updated_at`/sitemap invariance
were unverified. Their isolated-fixture gate is now resolved below; live deployment
acceptance is still outstanding.

Generic SPA settings outages now retain the normal shell with absolute RSS/Atom
links and canonical bootstrap defaults. Explicit `PUBLIC_SITE_URL` remains required
in every mode; article/feed metadata query errors remain **500/no-store**. SQL prepare
and finalize artifacts are separate, opt-in, repeatable phases; application-driven
pinyin/date repair stays between them and the completion marker is installed last.

## Real MySQL acceptance follow-up — 2026-11-06

Worktree: `/Users/iceymoss/project/inkspace-seo-mysql`; branch `fix/seo-mysql-acceptance`;
base `1d2ecf80`. Server version observed by `SELECT VERSION()`: **MySQL 8.0.33**.
The user authorized the already-running instance and its existing account for random-schema
tests. A local Node subprocess launcher parsed the existing config programmatically, passed a
**no-selected-database** DSN only through `INKSPACE_SEO_TEST_MYSQL_DSN`, and redacted secrets
from output. No credentials/DSN were printed or committed. `SEOTestMySQL` retained its selected-DB
refusal and created/dropped only its own `inkspace_seo_test_*` schemas. No business schema was changed.

The initial targeted acceptance run had **4 failures / 11 passes**, each reproduced narrowly:

- SQL phase fixture: the expected missing-opt-in error poisoned GORM's initialized pinned
  connection handle, so later statements inherited it. Fresh statement clones on the **same
  pinned connection** now isolate expected errors. SQL guards were not changed or relaxed.
- Both interaction fixtures: `SyncCounters` used a correlated self-update of `comments`, rejected
  with MySQL error 1093. A `GROUP BY`-materialized derived-table `LEFT JOIN` now counts undeleted
  direct replies and resets empty counts to zero. New assertions cover approved/rejected/deleted
  replies and stale counts, while article editorial/publication time and sitemap lastmod stay unchanged.
- Permalink fixture: an obsolete assertion still demanded a 500 for the generic SPA `/` after the
  reviewed outage fallback. It now verifies GET/HEAD 200/no-cache, default title, canonical bootstrap
  and absolute feeds. Settings API and article/feed/sitemap query failures still assert 500/no-store.

| Command (SEO DSN injected by local launcher) | Actual result |
| --- | --- |
| `go test ./internal/database ./internal/service ./internal/router -run 'MySQL\|Publication\|PublicArticleCollections\|Interactions' -v -count=1` | **15 SEO PASS, 0 SEO SKIP/FAIL**; 1 unrelated theme SKIP |
| `go test ./... -json -count=1` | **124 top-level PASS (412 including subtests), 28 unrelated SKIP, 0 FAIL; 10 test packages PASS**. All **16 SEO MySQL tests PASS**, including model slug allocation. |
| `go test -race ./... -json -count=1` | Same **124 / 412 PASS, 28 unrelated SKIP, 0 FAIL; 10 test packages PASS**; all **16 SEO MySQL tests PASS**, no race reported. |
| `go build ./...` / `go vet ./...` | **PASS / PASS** |
| `git diff --check` | **PASS** |

Executed SQL coverage includes prepare/finalize repetition, missing offline acknowledgements,
pre-repair rejection, missing/invalid/duplicate slugs, zero/partial/impossible publication dates,
malformed unique/marker indexes, populated slug preservation, and Go orchestration validation.
The legacy fixture verified `ON UPDATE CURRENT_TIMESTAMP(3)` really existed before exercising
changed counters; publication/content time and sitemap lastmod invariance now have real-SQL evidence.

Limitations: the 28 unrelated skipped tests use separate bootstrap/admin, reading-time/MySQL/Redis,
and theme opt-ins; those gates were not enabled or counted as passes. The interaction cache uses the
existing disposable RESP invalidation stub, not live Redis. No frontend source/dependencies changed;
the preceding **16 files / 114 tests** frontend result is historical, not rerun in this follow-up.
`verify.sql` was not run against a business/deployment database. No production data migration,
container/cluster rollout, live gateway/TLS/CDN acceptance, or push was performed.
