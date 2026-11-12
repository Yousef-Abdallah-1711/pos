# Production Readiness — 28 gates
Rule: PASS only with evidence. Nothing here is PASS.

| Gate | Status | Basis |
|---|---|---|
| G01 Services start reliably | PARTIAL | Dev stack starts and is healthy; prod stack built and ran with test values (earlier this session); startup runs migrations on every boot with no lock |
| G02 Auth enforced | **FAIL** | Clusters A, C; stub mode (SEC-012) |
| G03 RBAC enforced by backend | **FAIL** | Cluster D |
| G04 Tenant isolation | **FAIL** | Clusters B, E, F; no RLS (verified) |
| G05 Restaurant isolation | **FAIL** | Same; driver_repo 18/23 methods unscoped |
| G06 Core order workflows | FAIL | Non-atomic, no inventory decrement, number collisions |
| G07 POS | PARTIAL | Partial implementation; offline duplicates |
| G08 KDS real-time | NOT_TESTED | Code exists; WebSocket auth unverified |
| G09 Inventory integrity | FAIL | adjust route missing; no decrement on order |
| G10 Payments safe | **FAIL** | Cluster G |
| G11 Accounting correct | FAIL | FIN-001/003/005; float columns (4) |
| G12 HR workflows | NOT_TESTED | Static only |
| G13 Migrations safe | FAIL | No lock/checksum/rollback; hardcoded `postgres` role; duplicate prefixes |
| G14 Error handling consistent | FAIL | QA-002 |
| G15 Logs/observability | FAIL | No structured logs, no request IDs |
| G16 Security checks | FAIL | See P0 clusters |
| G17 Frontend builds reliable | PARTIAL | Landing/storefront/super-admin/dashboard build in Docker only with type-gate bypass for dashboard |
| G18 Backend builds reliable | PARTIAL | Docker image builds; `go build ./...` / tests not run (no Go on host) |
| G19 Responsive UI | PARTIAL | Public landing OK at 320/768/1440; authenticated UI NOT_TESTED |
| G20 Accessibility | NOT_TESTED | |
| G21 Arabic & English | PARTIAL | Landing OK; storefront maintenance ignores locale |
| G22 Performance | NOT_TESTED | No representative measurements |
| G23 CI/CD | FAIL | No build/test in CI |
| G24 Backup/restore | FAIL | None found |
| G25 Deploy & rollback | PARTIAL | New Docker prod setup verified once; no rollback procedure |
| G26 Regression coverage | PARTIAL | 62 Go + 59 frontend test files; not run; many behind build tags |
| G27 Dependency risks | NOT_TESTED | Registry/CVE checks not done |
| G28 Monitoring | FAIL | Public unauthenticated metrics only |
