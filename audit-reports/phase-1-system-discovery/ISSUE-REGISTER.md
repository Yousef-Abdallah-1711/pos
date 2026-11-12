# Issue Register
Full machine-readable register: `ISSUE-REGISTER.json` — **154 raw findings** merged from six auditors (P0 13, P1 40, P2 60, P3 39, P4 2). Many are the same defect found independently by several auditors; clustered below. Per-finding detail (file:line, evidence, fix, tests) lives in `evidence/issues-*.json` and the module reports. All status OPEN. "Verified" = reproduced or read first-hand by the lead auditor; everything else is static-reading by a subagent and still needs Phase 2 confirmation.

## P0 clusters (deduplicated)
| Cluster | Raw IDs | Verification |
|---|---|---|
| **A. Auth bypass via `X-Impersonation-Token`** — any non-empty header skips JWT validation (`middleware/auth.go:18-31`) | BE-001, SEC-001 | **Verified live**: 401 → passes auth. Data access then stopped at an in-handler role check on `/products` (403); other routes untested |
| **B. Header-controlled tenant** — `X-Tenant-ID/X-Restaurant-ID` override tenant when JWT `restaurant_id`=0 and on public routes (`middleware/tenant.go:26-63`) | DB-001, SEC-026, BE-005 | Code-read; needs fixture to reproduce |
| **C. Passwordless login confirm** — `POST /auth/login/confirm` mints a JWT from email+tenant_id | SEC-002, FE-OTH-001 | Code-read |
| **D. RBAC middleware is decorative** — `WithRequiredPermission` only sets context (`permission.go:16-24`), 193 routes wrapped; handlers do inconsistent ad-hoc checks; 17 `/rbac/*` routes + categories/users/drivers/themes have no wrapper | BE-002, BE-003, SEC-003, SEC-005, DB-004 | **Middleware behaviour verified in source; "no enforcement anywhere" is overstated** (see evidence/VERIFICATION-NOTES.md) |
| **E. Public order endpoints** trust caller tenant headers, sequential IDs, unauthenticated read/cancel | DB-002, FE-OTH-002, SEC-006, BE-006 | Code-read; BLOCKED (no orders in dev DB) |
| **F. Theme v2 IDOR** — get/update/delete/activate filter by id only | DB-003 | Code-read |
| **G. Mock Paymob is default; webhook public; mock HMAC always true; prod compose leaves `PAYMOB_MODE` empty** | SEC-004, BE-007 | Code-read; **my `docker-compose.prod.yml` has this default and must be fixed** |
| **H. Tenant status/subscription not enforced on business routes**; new tenants get no trial dates | FE-OTH-031, SEC-008, BE-008, FE-OTH-003 | Code-read |

## Other P1 highlights
Seeded super admin with documented password (DB-013/SEC-010); stub mode accepts any email when DB is down (SEC-012); no login rate limit (SEC-009); webhook replay/refund re-activation (SEC-015); unauthenticated translations upload overwrites global state (BE-022/SEC-007); `/health/metrics` public (QA-001, **verified**); no server timeouts/graceful shutdown/body limits (BE-009/QA-006); non-atomic order creation, colliding `order_number`, no inventory decrement (DB-010/011, BE-010); journal balance enforced only in Go, auto-posting plugs imbalances (DB-012, FIN-001); no RLS and app uses superuser (DB-005, **verified live**); shard tooling unsafe/unwired (DB-006..009); raw `err.Error()` to clients, 254 sites (QA-002); CI has no build/test (QA-003); no backup/restore (QA-004); POS offline replay can duplicate orders, SW sync wrong origin (FE-DASH-002/003); `/inventory/adjust` has no backend route (FE-DASH-001); JWT in URLs (FE-OTH-005); backend login logs password/hash prefix (FE-OTH-007).

## Observations confirmed in browser
Bare `/login`,`/register` on landing end at 404; storefront for a non-existent tenant shows a generic "Temporarily Unavailable" page; storefront `/ar` maintenance page renders `lang=en dir=ltr`; super-admin `<html>` lacks `dir`; dashboard `tsc` reports **705 errors** (run in dev container; the dashboard auditor could not run it) — production builds currently need `NEXT_IGNORE_BUILD_ERRORS=true`.
