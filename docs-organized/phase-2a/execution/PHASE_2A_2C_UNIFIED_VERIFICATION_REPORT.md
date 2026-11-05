# POOS Unified Phase 2A + Phase 2C Verification Report

Date: 2026-11-05 (Africa/Cairo)

## Git and source preservation

| Check | Result |
|---|---|
| `git fetch origin`; `git rev-parse origin/main` | PASS; `d84120559db36e8e9e7fdb73872a30cbc9d4e6e7` |
| Clean branch from `origin/main` | PASS; `codex/unified-phase2a-2c-closure` created at `H:\pos-unified-phase2a-2c` |
| Source worktrees | Preserved; primary and WS02 dirty contents not reset/cleaned/overwritten |
| Private snapshot checksums | PASS; 220/220 checksums verified |
| `git diff --check origin/main...HEAD` | PASS before local fixture repair; rerun before publish |

## Go verification

Toolchain: disposable `golang:1.24` container, Go `1.24.13 linux/amd64`. The container uses read-only source and isolated tmpfs PostgreSQL at `127.0.0.1:25436`; no development DB or volume is used.

Initial clean-branch run:

- `go build ./...`: PASS.
- `go vet ./...`: FAIL — `subscription_admin_isolation_integration_test.go` referenced `openProductHTTPTestDB`, `seedProductHTTPFixture`, and `insertSettingsHTTPUser`, which were untracked helper functions in the original WS02 worktree and absent from a clean checkout.
- `go test -count=1 ./...`: FAIL at HTTP package compilation for the same missing helpers.

Remediation: replace the hidden helper dependency with local disposable database and synthetic tenant/restaurant/user fixtures in the subscription integration test. The repair is committed on the unified candidate. A subsequent fresh isolated run completed migrations, build, vet and full tests with exit code 0. The full test suite ran against isolated tmpfs PostgreSQL with `POS_TEST_DATABASE_URL` pointed only to `pos_test`.

The rerun log has `MIGRATE_EXIT:0`, `BUILD_EXIT:0`, `VET_EXIT:0`, and `TEST_EXIT:0`. Explicitly selected Driver, HR child-resource, subscription admin, trial concurrency, tax-period, and recurring-template ownership tests all passed against the disposable PostgreSQL database. The unified Orders continuation also passed `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` under Go 1.24.13 against the disposable PostgreSQL test database.

## Database/security evidence

- A fresh PostgreSQL 15 container uses tmpfs and a unique port. Migration completed successfully.
- Targeted Driver, HR child-resource, subscription trial concurrency/admin, tax period, and recurring-template ownership tests passed, including unchanged-row/side-effect assertions, against disposable PostgreSQL.
- Authentication lifecycle, Redis session revocation, refresh-rotation races, CSRF negatives, audit failure behavior, and all resource-family side-effect checks were not run as one combined acceptance set.
- Do not interpret a response denial alone as isolation proof; unauthorized side effects must remain asserted.

## Frontend/browser/performance

- Frontend source typechecks: PASS for all four apps. Test-source typechecks: PASS for Landing, Restaurant Website and Super Admin; Dashboard fails with existing test fixture and Playwright typing errors.
- Unit tests on the unified candidate: Landing 35/35, Restaurant Website 66/66, Super Admin 16/16 pass. Dashboard: 538/645 pass, 107 fail across 14 files, with one unhandled error. The failures include stale theme/component expectations and timer/error-boundary tests; they remain open.
- Production builds: PASS for all four apps with local example URL configuration. Existing warnings include a dynamic component-registry dependency in Dashboard, `<img>` in Landing, a missing effect dependency in Restaurant Website and missing effect dependencies in Super Admin.
- Authenticated browser login/logout, session, tenant, POS, inventory, storefront cart, Super Admin denial and Admin Orders journeys: NOT RUN against a clean combined candidate stack.
- Admin Orders benchmark: same Go 1.24.13 build, PostgreSQL 15 tmpfs, 75 deterministic synthetic orders with two items each, 30 HTTP samples per page size. Baseline `d841205` → candidate `e30ffe3`: page 1: 4→3 SQL, median 12.188→9.806 ms, p95 16.354→13.012 ms; page 10: 22→3, 57.844→10.692 ms, 79.494→14.647 ms; page 25: 52→3, 142.035→10.610 ms, 211.610→19.043 ms; page 50: 102→3, 291.768→11.784 ms, 504.784→23.033 ms. Payload byte counts and canonical response SHA-256 values match for each page size. Tenant and sibling-restaurant fixture orders are excluded. Full JSON evidence is in `ADMIN_ORDERS_BENCHMARK_BEFORE.json` and `ADMIN_ORDERS_BENCHMARK_AFTER.json`. Historical page-50 figures (105 SQL, 146.87 ms median, 215.77 ms p95) came from a different profile and are not directly comparable.
- Existing PR #7 SQL-mock asserts tenant/restaurant predicates and batched order IDs. It does not establish HTTP behavior, before/after result equality or latency.
- Historical Dashboard and Storefront profiles are in the Phase 2C source branch. Dashboard home cold CLS 0.380, warm POS LCP 2.364 s; Storefront sample had invalid synthetic theme `warm-comfort` 404s. No fresh combined-build profiles were run.

## Hosted CI and review

- PR #1 current head checks failed `runtime-localhost-guardrail` due to stale `frontend/apps/dashboard/scripts/verify-pos-phase4.mjs` localhost fallback and `url-domain-config-smoke` due two stale Dashboard production URL expectations. Failure logs inspected.
- PR #3 latest checks pass; no review approval.
- PR #7 latest checks pass; draft; no review approval; benchmark still missing.
- PR #8 is merged and its storefront optimization is in current main.
- PR #9 head `ec23651ebc42a397b54244554635a0e281122d27`: `runtime-localhost-guardrail` PASS and `url-domain-config-smoke` PASS. No review approval is recorded.

## Final gate

No combined acceptance claim. Phase 2A PARTIAL; Phase 2C PARTIAL; Unified Git Integration PARTIAL; Production Release NO-GO. Exact remote main remains `d84120559db36e8e9e7fdb73872a30cbc9d4e6e7` at report time.
