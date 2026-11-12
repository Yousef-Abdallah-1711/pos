# Architecture decisions (Phase 2A)

ADR-1 to ADR-3 are implemented. ADR-4 to ADR-7 are proposals only.

- **ADR-1 Identity comes only from verified credentials.** Tenant, restaurant and user come from JWT claims or a server-validated impersonation session. Client headers are hints that must be verified (implemented: `RestaurantVerifier`).
- **ADR-2 Authorization is a middleware concern with explicit modes.** Declared route levels are enforced centrally. `audit` mode exists only to validate mappings and is refused in production (implemented).
- **ADR-3 Fail closed.** A missing verifier, evaluator or authenticator, or a lookup error, denies the request with 401, 403 or 503 (implemented).
- **ADR-4 Backend structure (proposal).** Keep `cmd/api` thin. Move route registration out of `main.go` into per-module route files that receive a small `Deps` struct (auth, tenant, permission and status wrappers). Migrate modules one at a time. Defer the `internal/modules/*` rename until the 16 failing tests and 9 stale test packages are repaired: a refactor needs a trustworthy suite first.
- **ADR-5 Frontend structure (proposal).** Keep the four apps. Resolve the undeclared `@pos-saas/ui` dependency. Delete or adopt `packages/theme` and `layers/shared-components` (no consumers). Add one `api-client` package to replace the per-app axios/fetch wrappers. Fix the vitest `@/` alias so frontend tests become a usable gate.
- **ADR-6 Public order access.** Needs an unguessable per-order tracking token (design in SECURITY_REMEDIATION.md, SEC-E).
- **ADR-7 Database.** Add a non-superuser runtime role and a separate migration role. Classify tables (global, tenant, restaurant, platform, audit) before any RLS. Carry tenant context with `SET LOCAL app.tenant_id` per transaction. Not started.
