# F-044 — Implementation phases and tasks (all NOT STARTED)

Finding: Cancelled orders remain active on the Kitchen Display queue · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-044-01
- **Finding:** F-044
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-044 at code level
- **Existing source paths:** backend/internal/repository/kds_repo.go, backend/migrations/124_create_kds_queue.sql, backend/internal/repository/kds_repo.go
- **Proposed paths:** backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_kds_cancelled.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-044-02
- **Finding:** F-044
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-044
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-044-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-044-03
- **Finding:** F-044
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-044
- **Existing source paths:** backend/internal/repository/kds_repo.go, backend/migrations/124_create_kds_queue.sql, backend/internal/repository/kds_repo.go
- **Proposed paths:** backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
- **Required changes:** Proposed: propagate cancellation to the kitchen. When an order is cancelled (usecase/order_usecase.go cancellation path and repository/order_repo.go status update) update the matching kds_queue row to status cancelled inside the same transaction (or add a trigger on orders.status that does so); exclude cancelled rows from /kds/queue, /kds/stats active_orders and analytics; broadcast the change over the KDS channel so open screens remove the card.
- **Dependencies:** REM-F-044-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Cancelling an order removes it from the KDS board within 2 s and from active_orders. | kds_queue has no active row whose order is cancelled (SQL assertion over the fixture). | Analytics exclude cancelled orders.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-044-04
- **Finding:** F-044
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-044 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
- **Required changes:** Cancel order and assert KDS queue/board and counts
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-044-05
- **Finding:** F-044
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-044
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
