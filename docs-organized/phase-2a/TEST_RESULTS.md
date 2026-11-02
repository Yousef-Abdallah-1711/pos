# Phase 2A test and build results (end of session 2)

All numbers below were produced in this session by the commands shown. Nothing is estimated.

## Backend (Go 1.24.13 in `golang:1.24-alpine`; Go is not installed on the host)
| Check | Command | Result |
|---|---|---|
| Build everything | `go build ./...` | **OK** |
| Build each command | `go build ./cmd/{api,migrate,migrate-shard,seed-themes,hash}` | **all 5 OK** |
| Static analysis | `go vet ./...` | **clean** (found a real metrics bug, fixed) |
| Tests | `go test -count=1 ./...` | **680 passed, 0 failed**; 13 packages with tests pass, 10 packages have no test files |
| Race detector | `-race` | not run (alpine image has no cgo) |
| Stability of the previously flaky KDS hub test | `go test -count=30 ./internal/handler/ws` | passes 30/30 after the fix |

Before this session: 16 failing tests, 9 failing packages, `go build ./...` failing in `scripts`. Root causes and fixes: see PHASE_2A_TASKS.md C013 (stale KDS fixtures after the `order_source` column; stale SQL regexes; pagination COUNT queries; tautological theme tests with inverted logic; one self-contradictory slug test; the real hub race; the real shard-hash defect; `go vet` findings).

**Quarantined, NOT passing and NOT counted** (`-tags experimental`, `backend/tests/EXPERIMENTAL.md`): `tests/unit/{auth,security,sharding,middleware,scaling}`, `tests/integration`. They test subsystems that `cmd/api` does not use and need a schema/DB that does not exist (e.g. `set_tenant_context()`, shard hosts). Running them against the disposable DB showed ~40 failures, e.g. the `internal/auth` RBAC manager targets columns (`roles.name`) the real schema does not have.

New tests added in Phase 2A (all in the 680): middleware 31 (impersonation 6, tenant context 9, permission 6, tenant status 4, internal token 3, + existing), jwt 4, usecase `TestLoginConfirm` 7 + `TestWebhookSecurity` 8, themes 3, config Paymob 6, domain public-order access 3, sharding router 4, database tenant-tx 6, `cmd/api` startup guards 2.
Integration (tag `integration`, disposable DB): `TestRLSPilot_TenantIsolationWithRestrictedRole` **PASS**.

## Live smoke on the disposable stack — `scripts/test-stack/security-smoke.sh`
Stack: `docker-compose.pos-test.yml` (Postgres 15 on tmpfs, API built from `backend/Dockerfile`, `ENV=test`, `PAYMOB_MODE=mock`, `RBAC_ENFORCEMENT=enforce`, API connected as `pos_app`). The script refuses to run unless both containers are labelled `com.pos.disposable=true`, the DB is `pos_test_disposable`, the backend has `POS_DISPOSABLE_TEST_STACK=1`, `ENV=test`, mock payments and a local URL.
**Result: 73 passed, 0 failed.** Sections: impersonation (3), tenant context (4), login confirm (2), theme isolation (6), RBAC incl. employee fixture (4), subscription status (5), public order access (15), formerly ungated routes (14), unsafe defaults (4), restricted DB role (15). Run history this session: 25 → 58 → 73.
Not covered by the smoke: orders created through `/public/orders` (rows were inserted directly), POS flow, KDS WebSocket, impersonation round trip, payment success activation, any frontend.

## Route compatibility
Unique `(method, path)` registrations in `main.go`: **338 at HEAD, 338 now; none removed, none added** (script comparison). Payload changes are intentional and listed in SECURITY_REMEDIATION.md (`tracking_token` on create, `X-Order-Token` on public order reads/cancel, `selection_token` on login-confirm).

## Frontend
| App | `typecheck` (app code) | `typecheck:test` | Unit tests |
|---|---|---|---|
| landing | **0** | 23 errors | cannot run (no Vitest config / alias; jsdom not a dependency) |
| super-admin | **0** | 3 errors | 1 file, 3 tests pass |
| restaurant-website | **0** | 24 errors | Jest, not run |
| dashboard | **0** (was 705 combined: 542 test + 163 app) | 158 errors | 41 files: 26 pass / 15 fail; 639 tests: 530 pass / 109 fail (was 596 tests: 334 / 262; 43 more tests now run because one jest-style file was ported to Vitest) |

Production builds (`frontend/Dockerfile.prod`, `NEXT_IGNORE_BUILD_ERRORS=false`, all `NEXT_PUBLIC_*` supplied): see the table at the bottom (filled in from the build logs). All four apps build without the type-gate bypass. Lint (`next lint`) was not run separately; `next build` runs it.

## Production builds without the type-gate bypass (`frontend/Dockerfile.prod`, `NEXT_IGNORE_BUILD_ERRORS=false`)
| App | Result |
|---|---|
| landing | **PASS** — `docker build` exit 0 (`Compiled successfully`, lint + type check passed, standalone image produced) |
| super-admin | **PASS** — `docker build` exit 0 |
| restaurant-website | **PASS** — `docker build` exit 0 |
| dashboard | **PASS** — `docker build` exit 0 without the bypass (built from the tree at the moment typecheck reached 0; later edits touched only tests, the Vitest setup and one `export` keyword) |
The bypass flag is now `false` in `docker-compose.prod.yml`.
