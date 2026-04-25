# AGENTS.md

## Project overview

InkSpace: Go + Vue 3 blog platform. Go module `github.com/iceymoss/inkspace`, Go 1.21+.

## Architecture

Three **separate** backend binaries, two frontend apps:

| Service | Entry | Default port |
|---|---|---|
| User/blog API | `cmd/server/main.go` | 8081 |
| Admin API | `cmd/admin/main.go` | 8083 |
| Scheduler | `cmd/scheduler/main.go` | (no HTTP port) |

Frontends live under `web/`:
- `web/blog/` — blog UI, dev on :3001, proxies `/api` → :8081
- `web/admin/` — admin UI, dev on :3002, proxies `/api` → :8083

Backend layering: `cmd/` → `internal/router/` → `internal/handler/` → `internal/service/` → `internal/models/` + `internal/database/`

## Dev startup

```bash
# 1. Infrastructure (MySQL 8 + Redis 7)
docker-compose up -d mysql redis

# 2. Config — YAML files are gitignored, must copy from examples
cp config/config.example.yaml config/config.yaml
cp config/admin.example.yaml config/admin.yaml
# edit config/*.yaml or create .env (see env.example)

# 3. Database — auto-migrates on startup; must exist first
mysql -u root -e "CREATE DATABASE IF NOT EXISTS inkspace CHARACTER SET utf8mb4"
# or use scripts/create_database.sql

# 4. Backend (each in its own terminal)
go run cmd/server/main.go      # user API :8081
go run cmd/admin/main.go       # admin API :8083
go run cmd/scheduler/main.go   # optional: hot articles, rankings

# 5. Frontend (pnpm required, not npm/yarn)
cd web/blog && pnpm install && pnpm dev   # :3001
cd web/admin && pnpm install && pnpm dev  # :3002
```

Default admin login: `admin` / `admin123` (auto-created on first DB init).

## Config gotchas

- Config priority: env vars > `.env` file > YAML
- Admin service loads `admin.yaml` first; falls back to `config.yaml`
- Admin port: uses `admin.port` if set, otherwise `server.port + 1`
- Separate JWT secrets: `jwt.secret` (users) vs `jwt.adminSecret` (admin)
- Upload storage: `local` (default) or Tencent COS (`upload.storageType`)

## Testing

**No Go tests exist** in the repo (`*_test.go` files absent).

Frontend lint only:
```bash
cd web/blog && pnpm lint
cd web/admin && pnpm lint
```

## Build & deploy

- Single Dockerfile builds all three Go binaries; `CMD` selects which to run
- Docker Compose runs 3 user-api replicas behind nginx, plus admin-api + scheduler
- Production config: `config/config.yaml` and `config/admin.yaml` mounted into containers
- GitLab CI: tag-triggered (`pro-<service>-v<YYYYMMDDHHMMSS>` or `dev-<service>-v…>`), builds Docker images, updates ArgoCD (GitOps under `deploy/`)

## Key conventions

- `internal/` is the only Go code directory; no `pkg/` usage
- Models use GORM with soft deletes (`deleted_at`)
- Counter fields (article_count, comment_count, etc.) are maintained via DB hooks, not computed
- Scheduler tasks are registered with named intervals in `cmd/scheduler/main.go`
- Frontend uses Vue 3 Composition API + Pinia + Element Plus + Vite
