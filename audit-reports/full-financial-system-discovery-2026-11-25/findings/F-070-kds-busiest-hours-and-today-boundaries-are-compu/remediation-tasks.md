# F-070 — Implementation phases and tasks (all NOT STARTED)

Finding: KDS busiest-hours and "today" boundaries are computed in UTC, not the restaurant local time · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-070-01
- **Finding:** F-070
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-070 at code level
- **Existing source paths:** backend/internal/repository/kds_repo.go, backend/internal/repository/kds_repo.go, backend/internal/repository/kds_repo.go
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-070-02
- **Finding:** F-070
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-070
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-070-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-070-03
- **Finding:** F-070
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-070
- **Existing source paths:** backend/internal/repository/kds_repo.go, backend/internal/repository/kds_repo.go, backend/internal/repository/kds_repo.go
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Proposed: introduce a restaurant time zone. Add a tenant/restaurant timezone setting (default Africa/Cairo) and use it for KDS "today" and busiest-hours buckets (backend/internal/repository/kds_repo.go:314-317 and :418-430: convert with AT TIME ZONE before CURRENT_DATE / EXTRACT(HOUR ...)); store timestamps as timestamptz going forward; audit other date-bounded reports (sales, dashboard, orders stats) for the same UTC cut-off (currently a CANDIDATE, not tested).
- **Dependencies:** REM-F-070-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** An order received at 23:30 local appears in the previous local day and in hour 23 (not hour 21). | KDS completed_today uses local midnight. | The audit of other reports lists every query that uses CURRENT_DATE or date_trunc on a UTC timestamp.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-070-04
- **Finding:** F-070
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-070
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Historical aggregates are recomputed on read; no stored data changes.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Historical aggregates are recomputed on read; no stored data changes.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-070-05
- **Finding:** F-070
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-070 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Orders at 23:30 and 00:30 local land in the right day and hour
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-070-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-070-06
- **Finding:** F-070
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-070
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_restaurant_timezone.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-070 no longer reproduce (expected: Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
