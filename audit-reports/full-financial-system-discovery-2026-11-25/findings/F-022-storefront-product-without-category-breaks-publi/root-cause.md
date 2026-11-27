# F-022 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED for API 500 (public_menu_handler.go:63 / product listing with NULL category); UI zero-items cause CONFIRMED in RUN-003: with CORS_ORIGINS lacking a wildcard (the middleware supports http://*.localhost:13113, middleware/cors.go:66-92) browsers block the storefront subdomain origin; once the isolated API was recreated with that origin allowed, see AUD-STORE-004. The isolated compose default does not include the wildcard (environment configuration)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/public_menu_handler.go | 63 | backend/internal/handler/http/public_menu_handler.go | exact | yes |
| backend/internal/middleware/cors.go | 66-92 | backend/internal/middleware/cors.go | exact | yes |
| audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml | 140 | audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml | untracked-audit-file | yes |

