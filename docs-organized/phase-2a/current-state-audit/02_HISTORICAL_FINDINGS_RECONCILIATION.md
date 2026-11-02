# Historical findings reconciliation

Audit date: 2026-11-02. Evidence hierarchy: current source and live read-only checks, then current test runs, then existing Phase 2A reports. Historical reports are not treated as proof.

| Claim | Current assessment | Evidence / limitation |
|---|---|---|
| 154 Phase 1 findings, 13 P0, 40 P1 | DOCUMENTATION ONLY pending item-by-item replay | Phase 1 files remain in `audit-reports/phase-1-system-discovery`; not every finding is reproduced in this audit. |
| 363 route registrations / 338 unique method-paths | REPRODUCIBLE from current committed `backend/cmd/api/main.go` inventory | Current Phase 2A route matrix reports 363 registrations and 338 unique method-path pairs; registrations may be duplicated between DB/stub handlers. |
| 274 permission-wrapped routes | IMPLEMENTED; broad route counts reproduced, not equivalent to resource-level authorization | Current matrix: 274 tenant-permission, 38 public, 20 authenticated, 26 super-admin, 3 session-token, 1 internal, 1 webhook. Route registration evidence plus 73/73 disposable security smoke. |
| 91 tables, no RLS | LIVE DEVELOPMENT VERIFIED | Read-only query against `pos-postgres/pos_saas`: PostgreSQL 15.19, 91 public base tables, 0 with `relrowsecurity`; 108 schema migration records. |
| Restricted DB runtime role | VERIFIED only in disposable stack | `docker-compose.pos-test.yml` sets `pos_app` runtime and `pos_migrator` for migration. Current dev compose still configures `DB_USER=postgres`; live role query showed only `postgres` role exists. |
| Backend build/vet/tests | VERIFIED in prior current-checkpoint run | Prior evidence: `go build ./...`, `go vet ./...`, `go test -count=1 ./...`; 680 tests pass. No source changes since that evidence are claimed except preserved pre-existing login-page edit; run details in `09_TEST_INFRASTRUCTURE_AUDIT.md`. |
| Dashboard tests | Reproduced in prior checkpoint | 530/639 pass, 109 fail across 15/41 files; not a passing regression gate. |
| Landing tests | CURRENTLY FAILED | In running landing container: 4 tests pass, 3 suites fail resolving `@/stores/authStore`; host lacks Vitest executable. |
| Storefront tests | CURRENTLY FAILED | In running storefront container: 42 pass / 8 fail, 3 suites fail; theme store API/persistence expectations, color utility behavior/contrast, and alias mapping. |
| Super-admin tests | VERIFIED in prior checkpoint | 3/3 passed; no rerun required for this audit. |
| All four app typechecks and production builds | Prior-run evidence only | Previously reported 0 app-code TypeScript errors and all four Docker production builds pass without bypass. No fresh build run because the dev apps use mounted `.next` volumes and rebuilding could alter active runtime artifacts. |
| Frontend runtime | CURRENTLY LIVE VERIFIED | API health plus ports 13001–13004 each returned HTTP 200 on 2026-11-02. This is a smoke check, not functional flow coverage. |
| Test stack | PARTIAL | Disposable stack is API+Postgres only. It omits Redis and four frontend applications; the 73 checks do not cover browser flows or all tenant resource paths. |

The historical phase report’s former host/branch baseline is superseded by the current baseline in `01_REPOSITORY_BASELINE.md`. Historical task claims were reconciled conservatively in `15_PHASE_2A_TASK_RECONCILIATION.md`.
