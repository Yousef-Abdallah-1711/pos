# F-071 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): config.IsProductionMode() is true only for production/prod (runtime_env.go:20-23); RBAC default mode is audit unless production (middleware/permission.go:36-42); mock Paymob allowed unless production (config/paymob_mode.go:29); metrics token and Redis/session requirements and secure cookies are gated the same way (health_routes.go:52, main.go:148-156, auth_handler.go:223)


## Origin and propagation
- Origin (first verified source location): `backend/internal/config/runtime_env.go`:20.
- Propagation: affected modules Platform configuration (RBAC, payments, sessions, metrics, cookies); related findings F-020, F-033; observed consequence: Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not
- Cause versus symptom: the symptom is what users or reports show (Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Stagin); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-SEC013, EX-RUN005-RBAC-MATRIX; scripts: tests/runtime/ph17_sec013.mjs, tests/runtime/ph17_r5_rbac_matrix.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

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

