# F-022 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED for API 500 (public_menu_handler.go:63 / product listing with NULL category); UI zero-items cause CONFIRMED in RUN-003: with CORS_ORIGINS lacking a wildcard (the middleware supports http://*.localhost:13113, middleware/cors.go:66-92) browsers block the storefront subdomain origin; once the isolated API was recreated with that origin allowed, see AUD-STORE-004. The isolated compose default does not include the wildcard (environment configuration)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/public_menu_handler.go`:63.
- Propagation: affected modules Storefront; Products; related findings none recorded; observed consequence: Uncategorised product (category_id NULL, allowed by the product form) makes GET public menu/products return 500 for the whole tenant; after assigning a category the API returns the product but the storefront page still shows 0 items; new-te
- Cause versus symptom: the symptom is what users or reports show (Uncategorised product (category_id NULL, allowed by the product form) makes GET public menu/products return 500 for the whole tenant; after ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-STORE-004, EX-RUN003-AUD-UIFORM-001; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/public_menu_handler.go | 63 | backend/internal/handler/http/public_menu_handler.go | exact | yes |
| backend/internal/middleware/cors.go | 66-92 | backend/internal/middleware/cors.go | exact | yes |
| audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml | 140 | audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml | untracked-audit-file | yes |

