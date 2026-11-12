# WS02 / WS03 Security Closure Tracker

Started 2026-11-02. Branch `main`, HEAD `1815273924b8ade42f51863262b997115776322b`. Existing owner changes are recorded in the current Phase 2A checkpoint; do not reset, stage, or revert them.

Statuses: TODO, IN_PROGRESS, VERIFIED_COMPLETE, PARTIAL, BLOCKED, FAILED.

| ID | Task | Status | Evidence / next action |
|---|---|---|---|
| SEC-001 | Repository baseline and owner-change preservation | VERIFIED_COMPLETE | Branch/HEAD/status and existing diff recorded; no branch switch or Git publication. |
| SEC-002 | Read standards and previous reports | VERIFIED_COMPLETE | Repository `CODING_STANDARDS.md`, required execution reports and required audit reports reviewed. |
| SEC-003 | Verify WS01 restricted runtime role | VERIFIED_COMPLETE | Live `pos-dev` catalog read-only: `pos_app` and `pos_migrator` are non-superuser/NOBYPASSRLS; `tenants.relrowsecurity=false`. Earlier checkpoint verifies API connects as `pos_app`. |
| SEC-004 | Inspect isolated test infrastructure | VERIFIED_COMPLETE | Added `docker-compose.security-closure.yml`; isolated project `pos-security-closure-20261102`, PostgreSQL 15, loopback port 15434, tmpfs, disposable labels; migrations applied only there. Existing `pos-test` was not reused. |
| SEC-005 | Establish baseline backend tests | VERIFIED_COMPLETE | Pre-change Go 1.24 `go test -count=1 ./...` passed in a network-isolated container with source/module cache read-only. Final regression is running below. |
| SEC-006 | Complete resource ownership inventory | PARTIAL | Existing audited route/resource matrices reviewed; complete verified route→handler→usecase→repository→table coverage across all 20 families remains open. |
| SEC-007 | Recount driver repository methods | VERIFIED_COMPLETE | Source has 24 interface/concrete methods: 16 driver/location/statistics, 8 assignment. No intentional global operation. Current admin driver and admin order handlers use scoped methods. |
| SEC-008 | Design deterministic two-tenant fixtures | VERIFIED_COMPLETE | Isolated fixture creates Tenant A/B and Restaurant A1/A2/B1/B2 using stable synthetic slugs; repeatable, explicit `POS_TEST_DATABASE_URL` required. |
| SEC-009 | Implement two-tenant/two-restaurant fixtures | PARTIAL | PostgreSQL fixture seeds both tenants/restaurants and representative driver per restaurant. Users, roles, products, inventory, orders, themes, HR/accounting and other families are not seeded. |
| SEC-010 | Implement cross-tenant read tests | PARTIAL | Real PostgreSQL driver read denial verifies no object returned; HTTP-level and other resource-family checks remain open. |
| SEC-011 | Implement cross-tenant update tests | PARTIAL | Real PostgreSQL driver status update denial and unchanged target row; other resources remain open. |
| SEC-012 | Implement cross-tenant delete tests | PARTIAL | Real PostgreSQL driver delete denial and target row remains; other resources remain open. |
| SEC-013 | Implement cross-restaurant isolation tests | PARTIAL | Real PostgreSQL driver read/update/delete denials and unchanged A2 row from A1 scope; HTTP and other resources remain open. |
| SEC-014 | Implement bulk-operation isolation tests | TODO | No verified bulk mutation matrix. |
| SEC-015 | Verify unauthorized side effects | PARTIAL | Driver base rows checked unchanged after denied writes; related side effects, forged HTTP context, background work and other resource families remain open. |
| SEC-016 | Remediate confirmed vulnerabilities | VERIFIED_COMPLETE | Existing scoped driver/assignment changes retained. Fixed PostgreSQL refund tenant-update type defect exposed by real refund race. |
| SEC-017 | Run WS02 security regression | PARTIAL | Driver repository PostgreSQL tests pass repeated; `TestKDSWebSocketReadPermissionIsRequired` confirms module 13 READ enforcement, and the 2026-11-09 isolated browser journey confirms an employee without KDS READ receives 403. Full resource/RBAC/HTTP matrix remains open. |
| SEC-018 | Map payment callback state machine | VERIFIED_COMPLETE | State transitions and code boundaries recorded in `WS03_PAYMENT_INTEGRITY_REPORT.md`; refund/void callbacks still require reconciliation. |
| SEC-019 | Inspect PostgreSQL transaction boundaries | VERIFIED_COMPLETE | Success updates payment+tenant in one transaction; refund locks payment, calls mock provider, then commits payment+tenant; history writes remain post-commit. |
| SEC-020 | Implement callback integration fixtures | VERIFIED_COMPLETE | Synthetic tenant/payment/order fixture backed by migrated isolated PostgreSQL; no dev data or provider credentials. |
| SEC-021 | Test duplicate callbacks | VERIFIED_COMPLETE | Twelve concurrent duplicate success callbacks; exactly one DB transition and two expected in-memory history entries. |
| SEC-022 | Test concurrent callbacks | VERIFIED_COMPLETE | Success callback race exercised against real PostgreSQL; repeated targeted run passed. |
| SEC-023 | Test out-of-order callbacks | PARTIAL | Failure→success→delayed failure and refund→delayed success remain covered; pending Paymob callbacks now return without payment/subscription writes, verified by `TestPaymentUseCaseHandleWebhookPendingPaymentDoesNotMutateState` and `TestPaymentPostgresPendingCallbackLeavesPaymentPendingUntilSuccess`. Cancellation→delayed-success remains open. |
| SEC-024 | Test refund and void races | PARTIAL | Two concurrent refund requests with mock provider yield one call and one commit; void/refund webhook reconciliation remains open. |
| SEC-025 | Test transaction rollback | VERIFIED_COMPLETE | Controlled tenant constraint failure after payment update leaves payment pending and tenant trial. |
| SEC-026 | Test provider/local reconciliation | BLOCKED | No durable reconciliation/outbox exists. Provider success before local commit cannot be made atomic; provider idempotency semantics remain unverified. |
| SEC-027 | Remediate confirmed payment defects | VERIFIED_COMPLETE | Fixed refund subscription update SQL shape; regression is the PostgreSQL refund race. Cross-system provider/local atomicity remains an architectural risk. |
| SEC-028 | Run WS03 security regression | PARTIAL | PostgreSQL scenarios passed repeated runs; full regression pending. Provider reconciliation and event matrix remain open. |
| SEC-029 | Run complete backend regression | VERIFIED_COMPLETE | Go 1.24 offline `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed; 13 test-bearing packages passed, 0 failed. New PostgreSQL driver/payment tests ran with explicit isolated DB URL. |
| SEC-030 | Verify Docker and database safety | VERIFIED_COMPLETE | `pos-dev` seven services remain up; PostgreSQL healthy/accepting, Redis PONG, API and all four frontends HTTP 200. `pos_app`/`pos_migrator` non-superuser/NOBYPASSRLS; 0 public tables with RLS enabled. Closure DB is a separately labeled tmpfs container. |
| SEC-031 | Update workstream reports | VERIFIED_COMPLETE | WS02 and WS03 reports updated with evidence and remaining gaps. |
| SEC-032 | Update Phase 2A tracker | VERIFIED_COMPLETE | Main execution tracker and remaining blockers updated; Phase 2A remains incomplete. |
| SEC-033 | Produce final security closure report | VERIFIED_COMPLETE | `WS02_WS03_SECURITY_CLOSURE_REPORT.md` records outcomes, tests, safety, file list and open blockers. |

## Security invariants

1. Tenant identity comes from authenticated server-side identity.
2. Restaurant ownership is checked for restaurant-scoped resources.
3. A resource ID alone never authorizes access.
4. Denied writes do not mutate base or related records.
5. Bulk operations cannot partially mutate unauthorized scope unless the contract explicitly requires it and proves atomicity.
6. Background work retains verified tenant scope.
7. Super-admin cross-tenant access requires explicit verified privilege and audit evidence.
8. Missing or invalid tenant scope fails closed.

No RLS rollout is included; no live payment provider is contacted.
