# Permissions and tenant isolation

The dashboard sidebar filters entries using `src/hooks/useRbac.ts` and `src/lib/rbac.ts`; this is a presentation control, not proof of API authorization. Server routes use `wrapProtected` or module-specific `wrapWithPermission` in `backend/cmd/api/main.go`. Tenant context is obtained by middleware and repositories can use `database.WithTenantTx` to set transaction-local RLS context. The existing user-settings policies derive tenant ownership through `users`.

Future writes must have a specifically chosen tenant-admin permission, server-side tenant context, forced RLS, actor-aware audit, bounded input, and concurrency protection. Verify user, tenant and restaurant scope at each layer and write negative tests for cross-tenant access, no tenant context, invalid actor and missing permissions. Do not expose tenant IDs as authoritative client input. Actual current role assignments and complete route-level permission coverage remain unverified.
