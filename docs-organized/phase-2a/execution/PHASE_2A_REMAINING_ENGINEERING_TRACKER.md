# Phase 2A Remaining Engineering Tracker

**Started:** 2026-11-02  
**Branch:** `main`  
**Starting HEAD:** `5ab6dba62287ddd404ce00529f5fe99220ccdc19`  
**Scope:** WS02, WS04, WS05, WS07, WS08, WS09, WS10. WS03 payments and WS06 GitHub Actions are explicitly deferred.

Statuses: TODO, IN_PROGRESS, IMPLEMENTED_UNVERIFIED, VERIFIED_COMPLETE, PARTIAL, BLOCKED, DEFERRED.

| ID | Task | Files | Tests | Status | Evidence |
|---|---|---|---|---|---|
| WS01 | Preserve least-privilege database foundation | Existing implementation and WS01 report | Historical isolated role/migration checks; fresh Go baseline passed; restricted-role RLS pilot passed on disposable PostgreSQL | VERIFIED_COMPLETE | Kept WS01 implementation untouched; fresh pilot evidence is in `WS10_TENANT_TRANSACTION_FOUNDATION_REPORT.md`. |
| WS02 | Complete tenant/restaurant ownership and side-effect matrix | Driver repository/handlers and new matrix/report as completed | Three fresh driver PostgreSQL scope tests passed against disposable migrated PostgreSQL | PARTIAL | Driver read/update/delete/insert scoping and unchanged-row effects are verified; other resource families, forged HTTP context, bulk and worker mutations remain open. |
| WS03 | Payment integrity | None in this execution | None (explicitly excluded) | DEFERRED | Do not invoke providers or change payment state/reconciliation code. |
| WS04 | Recover Dashboard tests and test-source TypeScript | Dashboard config, theme store, tests, test configuration as completed | 639/639 Dashboard; 35/35 Landing; 64/64 Restaurant Website; 3/3 Super Admin; all app/test typechecks and builds pass | VERIFIED_COMPLETE | Fresh four-app run on 2026-11-03. Dashboard build generated 139 static pages. |
| WS05 | Disposable seven-service integration stack and E2E | Compose, synthetic seed, browser tests as completed | Existing `pos-dev` has seven healthy/running services; a separate tmpfs PostgreSQL RLS pilot passed; complete isolated stack/E2E not run | PARTIAL | No isolated Redis/four-app services, deterministic comprehensive fixture, browser journeys, or full-stack teardown proof. |
| WS06 | GitHub Actions / hosted CI | None (explicitly excluded) | None | DEFERRED | No workflow, runner, branch-protection, or remote CI changes. |
| WS07 | Backend architecture and route parity | Backend route/composition files as completed | Fresh Go build/vet/default tests passed; no route manifest was captured | TODO | Existing composition root remains large; no extraction made. |
| WS08 | Shared frontend contracts and API/auth boundaries | Frontend shared packages and app contracts as completed | All four app and test-source typechecks/builds pass; all 741 frontend tests pass | PARTIAL | Compatibility regression is green; shared theme/API/auth/tenant/restaurant/permission contract consolidation remains open. |
| WS09 | Typed errors, request IDs, redaction, lifecycle | `backend/internal/middleware/request_context.go`, API server, WS09 report | New middleware tests, full Go regression, HTTP health smoke and SIGTERM exit 0 pass | PARTIAL | Request IDs, panic recovery, HTTP timeouts, signal shutdown and DB-health monitor cancellation implemented; typed error mapping, backend-wide log migration, call deadlines and worker cleanup remain. |
| WS10 | Tenant transaction adoption and pooled-connection verification | Existing `backend/internal/pkg/database/tenant_tx.go` and integration pilot | Fresh Go baseline and tagged PostgreSQL RLS/transaction pilot passed in disposable tmpfs DB | PARTIAL | Helper's pilot is freshly verified; selected business-module adoption and broader concurrent pool checks remain open. |
| DOCS | Reconcile blockers, reports, closure record | Phase 2A execution docs | `git diff --check` PASS | VERIFIED_COMPLETE | Closure report records exact scope, results and blockers. |
| PUB | Review, commit and publish verified changes | Source/docs changed in this execution | WS02 cross-resource isolation and WS05 full isolated-stack/E2E gates remain incomplete | BLOCKED | Required security/integration gates are unmet. Repository `CODING_STANDARDS.md` §10.5 also prohibits a direct push to `main`; publication must use an approved PR path. |

## Baseline and constraints

- The worktree was clean at start, on `main`, at `5ab6dba62287ddd404ce00529f5fe99220ccdc19`; local and remote were aligned at initial inspection.
- Docker Desktop's Linux engine was recovered during this continuation. Docker daemon startup also brought unrelated Compose projects online automatically; no further action was taken on them.
- Go is not installed on the Windows host. The verified Go 1.24.13 toolchain in `golang:1.24` was used for fresh isolated backend checks.
- Node.js 24.16.0 and pnpm 11.28.9 are available. Dashboard tests and TypeScript checks are being run locally.
- The `pos-dev` database and named volumes were not modified. A tmpfs-only database was used for the RLS pilot. No payment provider or GitHub workflow was touched.

## Execution log

| Step | Status | Evidence |
|---|---|---|
| Standards and repository baseline | VERIFIED_COMPLETE | Read the full repository `CODING_STANDARDS.md`; initial Git state is above. The separately referenced `C:\Users\OPT\Desktop\POS\CODING_STANDARDS.md` path does not exist on this host. |
| Docker/runtime availability | VERIFIED_COMPLETE | Docker Linux engine recovered; `pos-dev` services are up, PostgreSQL/Redis healthy, and API health returns HTTP 200. Unrelated projects auto-started with Docker Desktop and were not operated on afterwards. |
| Dashboard baseline and targeted remediation | VERIFIED_COMPLETE | Latest full run: 639/639 passing across 41 files; app and test-source typechecks pass with zero diagnostics; production build generated 139 static pages. |

### Earlier local Windows rerun — superseded — 2026-11-02

At this earlier checkpoint, the Dashboard count was **513 passed / 126 failed of 639**. App and test-source typechecks passed with zero diagnostics and the production build passed. The matcher setup was made explicit for the app's Vitest instance; the failure groups and runtime versions are recorded in `WS04_FRONTEND_TEST_RECOVERY_REPORT.md`. The Dashboard recovery completion and current publication status are recorded below.

### Dashboard recovery completion — 2026-11-02

The latest Dashboard run passes **639/639 tests across 41 files**. Dashboard app and test-source typechecks pass with zero diagnostics, and the production build generated 139 static pages. WS04 is verified complete. WS02's full resource/forgery/bulk/worker matrix, WS05's seven-service E2E, WS07/WS08/WS09 architecture/foundation work and selected WS10 application adoption remain incomplete, so publication stays blocked.

### Fresh four-app regression and WS09 subset — 2026-11-03

- Frontend tests: Dashboard 639/639, Landing 35/35, Restaurant Website 64/64, Super Admin 3/3 (741 total). Dashboard Vitest used `--pool=threads --maxWorkers=4` because the default fork pool hit Windows `spawn EPERM`; no tests were skipped.
- All four app typechecks, all four test-source typechecks, and all four production builds passed with explicit loopback-only URL configuration.
- Go 1.24.13 container: build, vet and full default test suite passed after WS09 code changes.
- Restricted-role tenant/RLS integration pilot passed against a uniquely named PostgreSQL 15 tmpfs container; that project was torn down by its exact Compose project name.
- WS09 now includes server-generated request IDs, safe panic recovery, HTTP timeouts, signal-aware shutdown, and cancellation of the DB health monitor. A disposable runtime smoke returned HTTP 200 and exited 0 after SIGTERM; see `WS09_ENGINEERING_FOUNDATION_REPORT.md`.
- WS02 cross-resource isolation, WS05 isolated seven-service browser E2E, WS07 route parity, broader WS08 contracts, and WS10 business-module adoption remain incomplete. Publication remains blocked.
| Other frontend regression | VERIFIED_COMPLETE | Landing 35/35, Restaurant Website 64/64, Super Admin 3/3; app and test-source typechecks pass for each; all four production builds pass using local-only public URLs. |
| Security/integration/backend verification | PARTIAL | Fresh `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed in Go 1.24.13 container. Tagged PostgreSQL RLS/transaction pilot passed against a new tmpfs-only disposable DB. Full WS02 resource matrix and isolated E2E remain open. |
