# F-034 — Implementation phases and tasks (all NOT STARTED)

Finding: Waste/shrinkage adjustments change stock but post nothing and the turnover report shows waste 0 · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-034-01
- **Finding:** F-034
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-034
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-034-02
- **Finding:** F-034
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-034 at code level
- **Existing source paths:** backend/internal/handler/http/inventory_adjustment_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-034-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-034-03
- **Finding:** F-034
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-034
- **Existing source paths:** backend/internal/handler/http/inventory_adjustment_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: post inventory adjustments. In backend/internal/handler/http/inventory_adjustment_handler.go Adjust (:29-105) call the auto-posting engine inside the same transaction: waste/shrinkage Dr waste expense (or COGS adjustment) / Cr inventory 1130 at the product cost; count gains the reverse; include the amounts in the turnover report waste column (today 0) by reading the movement type.
- **Dependencies:** REM-F-034-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Adjusting 3 units of a product costing 40 posts Dr waste 120 / Cr 1130 120 and stock and ledger agree. | The turnover report shows waste 120 for the period. | Replaying the same adjustment request (idempotency key) posts once.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-034-04
- **Finding:** F-034
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-034
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Past adjustments are listed with cost impact; accountant approves a catch-up journal.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Past adjustments are listed with cost impact; accountant approves a catch-up journal.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-034-05
- **Finding:** F-034
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-034 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Adjustment/waste scenario tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-034-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Expense posted and waste reported; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-034-06
- **Finding:** F-034
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-034
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-034 no longer reproduce (expected: Expense posted and waste reported); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
