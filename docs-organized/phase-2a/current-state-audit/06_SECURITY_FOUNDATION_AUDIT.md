# Security foundation audit

Classification distinguishes code from verification. No external payment was executed and no production data was accessed.

| Domain | Current status | Evidence / remaining gap |
|---|---|---|
| Authentication | PARTIALLY VERIFIED | Bearer JWT middleware and server-validated impersonation session are described in `internal/middleware/auth.go` and Phase 2A tests. Public login/register/forgot-password are rate limited per current route/source inventory. Password reset and cross-app token transport remain incomplete review areas. |
| Tenant identity | PARTIALLY VERIFIED | Tenant context comes from validated token; restaurant headers require verifier and failure is intended fail-closed. 73 disposable smoke checks passed; coverage is not all internal call paths/resources. |
| RBAC | IMPLEMENTED / route-level evidence | Current registration counts: 274 tenant permission, 26 super-admin. Disposable `RBAC_ENFORCEMENT=enforce` smoke passed. Production forces enforcement; dev override is explicit `enforce`; parser tests exist. Route-level wrappers do not prove object-level scope. |
| Default mode | ENVIRONMENT-SENSITIVE | Prior review: outside production default is `audit` (log-only) if unset. Dev compose override explicitly enforces; generic/dev and prod configuration differ. Verify every deployment setting. |
| Tenant/resource isolation | PARTIAL | Theme isolation and multiple route negatives have tests. Driver repository historical finding 18/23 unscoped has not been fully recounted/reproved. Product/category/order/inventory/users/roles/settings still require paired tenant+restaurant tests. |
| RLS / DB privilege | DEV FAIL / TEST STACK PASS | Live development API uses superuser `postgres`; no RLS enabled. Disposable stack uses `pos_app` and `pos_migrator`, with 73/73 security smoke. |
| Payment/webhook | PARTIAL | Paymob HMAC verification exists; mock/live mode guard and tests cover invalid signature, unknown order and several replay/state cases. Missing verified transactional success activation, concurrent duplicate delivery, out-of-order callbacks, retry recovery. No live callback sent. |
| Public APIs/storefront | PARTIAL | Public order token design uses 256-bit token, hash at rest, phone proof fallback for reads, uniform 404 per implementation report; migration 161 in source and applied only to disposable DB. Sequential IDs/public serialization and browser compatibility require further verification. |
| Secrets/config | MIXED | Development compose contains insecure local defaults (e.g. DB `postgres`, dev JWT default); local override requires `.env.local` variables and mock Paymob. Production config has startup guards documented; secrets values intentionally omitted. |
| CORS/headers/rate limits | PARTIAL | Explicit localhost origins in local compose, but comprehensive prod-origin/header/cookie review is not completed. Auth rate limits exist; distributed multi-replica enforcement not proven. |

## Security findings (evidence template)

### SEC-001 — Development API superuser

- Category/severity: Database privilege, P0 immediate containment / P1 foundation closure.
- Status: CODE VERIFIED + LIVE DEVELOPMENT VERIFIED.
- Affected: `docker-compose.yml` backend service, `backend/cmd/api/main.go` DB startup.
- Observation: development runtime connects as `postgres`; live role inspection shows superuser/BYPASSRLS; 91 public tables, 0 RLS.
- Expected: runtime identity has only necessary DML; migrations use separate DDL identity.
- Impact: API compromise permits schema/role changes and bypasses any RLS protection.
- Recommended remediation: rehearse ownership/privilege transfer on a disposable copy; then separate migrator/runtime identities and validate migrations and rollback.
- Acceptance: live dev app uses restricted role; DDL and schema migration table access denied; API smoke and migrations succeed through separate migrator.
- Dependencies: role grants/default privileges, existing ownership transfer, migration 050 grants, compose configuration.

### SEC-002 — RLS absent from live tables

- Category/severity: Tenant isolation, P1 (critical if app query isolation regresses).
- Status: LIVE DEVELOPMENT VERIFIED absent; tenant transaction pilot implementation exists.
- Affected: application tables / `internal/pkg/database/tenant_tx.go`.
- Impact: isolation relies on app middleware and every query path; RLS is not a backstop.
- Remediation/acceptance: module migration and negative isolation tests before gradual RLS enablement; no table enabled until all paths use scoped transactions.

### SEC-003 — Incomplete resource and callback verification

- Category/severity: Authorization/payment integrity, P1.
- Status: PARTIAL; no claim of exploitability without a repro.
- Affected: resource repositories, Paymob webhook transaction path, storefront.
- Remediation/acceptance: cross-tenant/restaurant ID tests for each priority repository; transaction/idempotency/concurrency tests for callbacks; browser test with disposable fixtures.

See `07_ROUTE_AUTHORIZATION_MATRIX.md`, `08_RESOURCE_ISOLATION_MATRIX.md`, and `docs-organized/phase-2a/SECURITY_REMEDIATION.md`.
