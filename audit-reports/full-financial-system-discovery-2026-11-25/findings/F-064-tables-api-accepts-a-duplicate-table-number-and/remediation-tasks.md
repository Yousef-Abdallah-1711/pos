# F-064 — Implementation phases and tasks (all NOT STARTED)

Finding: Tables API accepts a duplicate table number and returns 500 for an over-long name · Severity Low · Root-cause confidence PROBABLE

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-064-01
- **Finding:** F-064
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-064 at code level
- **Existing source paths:** backend/internal/handler/http/table_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph16_tables.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-064-02
- **Finding:** F-064
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-064
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-064-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-064-03
- **Finding:** F-064
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-064
- **Existing source paths:** backend/internal/handler/http/table_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Unique constraint (tenant_id, restaurant_id, table_number) with 409 mapping; length validation
- **Dependencies:** REM-F-064-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Failing test from REM-F-064-01 passes.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-064-04
- **Finding:** F-064
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI/API contract for F-064
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Display/validate values returned by the corrected backend; no hard-coded rates/symbols; keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-064-05
- **Finding:** F-064
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-064 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Duplicate number returns 409; over-long name returns 400
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-064-06
- **Finding:** F-064
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-064
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
