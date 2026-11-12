# 02 Docs vs Implementation

Docs read: README.md (651 lines), PROJECT_KNOWLEDGE_BASE.md (1038), CODING_STANDARDS.md (2303, headings), AGENTS.md, docs-organized index (19 sections).

| Claim (source) | Class | Evidence |
|---|---|---|
| "Go 1.18+ (stdlib net/http, no framework)" README Tech Stack | Partial / stale | go.mod `go 1.24.0`; main.go uses Go 1.22 `mux.Handle("GET /path/{id}")` patterns (needs >=1.22). gin v1.11 and chi v5 in go.mod; gin only used in non-wired code (internal/handler/*.go, internal/handlers) |
| "157+ migrations" | Partial | backend/migrations has 123 files, numbered 001..160 plus 999_theme_migration_legacy.sql; gaps, and duplicate prefixes (067,070-072,075,076,078-080,099,123) exist |
| "120+ API endpoints" | Verified (understated) | 361 `mux.Handle*` registrations in cmd/api/main.go |
| Trial "30 days" (README pricing) vs "trial 15 days" (README line ~297) | Inconsistent | migration 112_seed_existing_tenants_trial.sql uses 15 days default; README refers to nonexistent `112_seed_trial_data.sql` |
| "Production Ready" (PROJECT_KNOWLEDGE_BASE.md line 4) | Not supported | unauthenticated /health/metrics, default DB creds, no server timeouts, no backups (see quality docs) |
| Middleware chain CORS->Auth->Tenant->RBAC->Status | Verified | main.go wrapProtected helpers ~line 600 |
| JWT HS256 24h | Verified | config.go JWT_EXPIRY default 24h; production guard on secret length |
| Docker `docker compose up` dev + prod compose | Verified (uncommitted) | docker-compose.yml, docker-compose.prod.yml |
| Backend `go build ./...` passes via `legacy` build tags | Partial / unverified | `//go:build legacy` found in admin_personalization_handler.go, theme_routes_v2.go; not compiled by me |
| "Redis" not claimed in README tech stack | Implemented-undocumented | internal/service/cache.go RedisCache via REDIS_URL; no Redis in any compose file |
| Observability / perf metrics endpoint | Implemented-undocumented | internal/observability/perf_metrics.go, /api/v1/health/metrics |
| AGENTS.md: read `C:/Users/OPT/Desktop/POS/CODING_STANDARDS.md` | Deprecated/stale path | machine-specific absolute path |
| Playwright E2E (playwright.config.ts testDir ./e2e/tests) | Documented-not-found | no e2e/ directory tracked |
| Rollback script backend/scripts/rollback_phase6_accounting_hr.sql | Verified exists | file listed |
| README mojibake (box-drawing, arrows) | Defect | README bytes show double-encoded UTF-8 when read |
| Super-admin default creds via migration 109 | Not verified | not opened |
