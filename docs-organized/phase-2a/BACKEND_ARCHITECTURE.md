# Backend architecture — current state and target (Phase 2A)

## Current state (verified by reading code and running it)
- One binary, `cmd/api/main.go` (~1,900 lines): configuration, repository/usecase/handler construction, middleware wrappers and **all** route registration (363 registrations). Other commands: `migrate`, `migrate-shard`, `seed-themes`, `hash` — all build.
- Layering is not uniform: roughly 87 endpoints go handler → use case → repository, 123 go handler → repository directly, 65 handlers take `*sql.DB`, 40 are inline closures with SQL in `main.go` (Phase 1 inventory). About 31.6k of 96k non-test lines are behind `//go:build legacy`.
- Request pipeline now (all composed in `main.go`): `AuthMiddleware` (bearer JWT, or a **server-validated** impersonation session) → `NewTenantContextMiddleware` (tenant from the token; restaurant verified against the tenant) → `TenantStatusMiddleware` (blocked tenants stop here, billing/status routes excluded) → `NewPermissionMiddleware(module, level)` (effective level from roles; `enforce` in production) → handler. Public routes use `tenantContext` with header verification; auth endpoints are rate limited per IP; operational endpoints use `RequireInternalToken`.
- Dead/unwired code: `internal/auth`, `internal/handlers`, `internal/scaling`, `internal/security` RLS manager; sharding is only reachable via `cmd/migrate-shard` (its router now routes keys correctly — see REFACTORING_LOG).

## Boundaries (what is now enforced in code)
| Concern | Where | Rule |
|---|---|---|
| Authentication | `middleware/auth.go` | Never from header presence; impersonation only through `ImpersonationAuthenticator` |
| Tenant identity | `middleware/tenant.go` | Token wins; headers need `RestaurantVerifier`; no verifier → 403 |
| Authorization | `middleware/permission.go` | Declared module/level compared with `RepositoryPermissionEvaluator`; fail closed on error |
| Subscription | `middleware/tenant_status.go` | Disabled/expired tenants blocked on business routes |
| Public object access | `domain/order_public_access.go` + `usecase.AuthorizePublicOrderAccess` | Token hash (256-bit) or phone proof; identical 404 |
| Payment config | `config/paymob_mode.go` | Production = explicit live, else startup fails |
| Startup safety | `cmd/api/startup_guards.go` | Production refuses stub mode and the seeded super-admin password |
| DB identity | `scripts/db/create_runtime_roles.sql` | Runtime role cannot do DDL |

## Target structure (decision record; **not implemented**)
Keep `internal/{domain,repository,usecase,handler}` while extracting composition: (1) move route registration from `main.go` into per-module files (`routes/<module>.go`) taking a small `Deps` struct (auth, tenant, permission and status wrappers, repositories); (2) migrate one module at a time to handler → use case → repository, starting with those that already have tests (payments, KDS, themes, RBAC), then orders; (3) introduce tenant-scoped transactions (`database.WithTenantTx`) at the use-case boundary as modules move, which is also the precondition for RLS; (4) delete the `legacy`-tagged tree and the dead packages above after the owner confirms; (5) typed domain errors + a single HTTP error mapper (254 handlers currently return raw `err.Error()`), request IDs and structured logging.
Why not done in Phase 2A: the safety net was the blocker, not the design. The suite was red and the security fixes needed their own tests first; large moves in `main.go` would have mixed refactoring with security changes. The suite is now green (680 passing tests, `go vet` clean), so Phase 2B can start extracting routes safely.
