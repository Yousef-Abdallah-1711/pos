# F-071 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): config.IsProductionMode() is true only for production/prod (runtime_env.go:20-23); RBAC default mode is audit unless production (middleware/permission.go:36-42); mock Paymob allowed unless production (config/paymob_mode.go:29); metrics token and Redis/session requirements and secure cookies are gated the same way (health_routes.go:52, main.go:148-156, auth_handler.go:223)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/config/runtime_env.go | 20 | backend/internal/config/runtime_env.go | exact | yes |
| backend/internal/middleware/permission.go | 36 | backend/internal/middleware/permission.go | exact | yes |
| backend/internal/config/paymob_mode.go | 29 | backend/internal/config/paymob_mode.go | exact | yes |
| backend/cmd/api/health_routes.go | 52 | backend/cmd/api/health_routes.go | exact | yes |
| backend/cmd/api/main.go | 705 | backend/cmd/api/main.go | exact | yes |
| backend/cmd/api/signing_secrets.go | 9 | backend/cmd/api/signing_secrets.go | exact | yes |
| runtime_env.go | 20-23 | backend/internal/config/runtime_env.go | basename | yes |
| middleware/permission.go | 36-42 | backend/internal/middleware/permission.go | prefix | yes |
| main.go | 148-156 | backend/cmd/api/main.go | basename-ambiguous(5) | yes |
| auth_handler.go | 223 | backend/internal/handler/http/auth_handler.go | prefix | yes |

