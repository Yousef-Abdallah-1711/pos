# Refactoring log (Phase 2A)

Scope was kept deliberately small, with no module moves. These changes are separate from the owner's earlier uncommitted work (see BASELINE.md).

1. `middleware.AuthMiddleware` gained an optional variadic `ImpersonationAuthenticator`. Existing call sites are unchanged; main.go builds one slice and passes it to all five wrappers.
2. `middleware.TenantContextMiddleware` stays as the strict default. New `NewTenantContextMiddleware(verifier)`. main.go defines `tenantContext` once and replaces 11 direct calls.
3. `WithRequiredPermission` is kept (annotation only). Enforcement lives in `NewPermissionMiddleware`.
4. The permission fallback helpers moved from `handler/http/rbac_permission_handler.go` into `domain` (`MaxPermissionLevel`, `FallbackPermissionLevel`). The handler delegates. No behaviour change.
5. `AuthUseCase` gained a `loginConfirmStore` seam (an interface over two repo methods) so it can be unit tested. The `NewAuthUseCase` signature is unchanged.
6. `wrapProtected` now includes the tenant status check. `wrapProtectedNoStatus` was added for billing, `/me` and status routes.

Planned but not done: extracting route registration from the roughly 1800-line `main.go`, module boundaries, frontend package consolidation, error-handling foundations (typed errors, request IDs), and removal of about 31.6k lines of `legacy`-tagged code.

## Session 2 additions
7. `main.go` wiring helpers: `tenantContext` (verified tenant context), `requirePermission` (permission middleware in `enforce`/`audit` mode), `statusGuard` (tenant status, only when the repository exists), `wrapAuthPermission`; `wrapProtected` now includes the status check and `wrapProtectedNoStatus` was added. Route set unchanged (338 unique registrations before and after).
8. 74 route registrations rewritten by `scripts/test-stack/_wire_permissions.py` (paren-matching codemod; policy table inside the script). 28 `/hr/*` routes re-mapped to the modules that exist in the seeded `modules` table (see SECURITY_REMEDIATION.md).
9. `middleware/permission.go` rewritten (`NewPermissionMiddleware`, `RepositoryPermissionEvaluator`, `ParsePermissionMode`); `WithRequiredPermission` kept as annotation-only.
10. Public orders: `domain/order_public_access.go`, `OrderRepository.SetPublicTokenHash/GetPublicAccessInfo`, `OrderUseCase.AuthorizePublicOrderAccess`, handler guard `requirePublicOrderAccess`.
11. `internal/handler/ws/kds_hub.go`: `addClient` + `drainRegistrations` extracted from `Run` to fix a register/broadcast ordering race.
12. `internal/sharding/router.go`: `jumpHash` replaced by a jump consistent hash (the old one ignored its key).
13. `repository/kds_repo.go`: 8 `defer observability.RecordDB(..., time.Since(...))` sites wrapped in closures.
14. Test infrastructure: `backend/tests/testdb` (guarded test DB URL), build tag `experimental` on six unwired suites, `//go:build ignore` on the two standalone `scripts/*.go` mains, `internal/validation` test tagged `legacy` like the package it tests.
15. Frontend: per-app `tsconfig.json` (app code) / `tsconfig.test.json` (tests) split and `typecheck` scripts; dashboard Vitest setup loads jest-dom and excludes Playwright specs and the unreferenced root `components/`/`lib/`; shared `form-errors.ts`, `qrcode.d.ts`, `product-types.ts`; codemods `_fix_field_errors.py`, `_dedupe_object_keys.py`, `_dashboard_small_fixes.py` (kept for audit, safe to delete).
16. Storefront: `lib/order-access.ts` + interceptors in `lib/api-client.ts` (store and send the order token), phone field on the track page.

Not refactored (planned): extracting routes from `main.go`, moving use cases/repositories into module packages, typed errors/request IDs, frontend package consolidation, feature-boundary refactors.
