# F-040 — Implementation phases and tasks (all NOT STARTED)

Finding: No COGS or inventory movement is posted for products with track_inventory=false even when a cost is set; Profitability report still deducts cost · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-040-01
- **Finding:** F-040
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-040
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

### REM-F-040-02
- **Finding:** F-040
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-040 at code level
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/handler/http/product_handler.go, backend/migrations/180_order_stock_reservations.sql, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_inventory_controlled.mjs).
- **Dependencies:** REM-F-040-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-040-03
- **Finding:** F-040
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-040
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-040-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-040-04
- **Finding:** F-040
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-040
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/handler/http/product_handler.go, backend/migrations/180_order_stock_reservations.sql, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: remove the track_inventory guard at backend/internal/usecase/auto_posting_engine.go:173 for COGS (post cost-of-goods for any sold item with a cost, relieving 1130 only when tracked, otherwise a configurable 'cost of non-stocked items' account), default track_inventory to true in the product form or make its effect explicit, and align the profitability report with the GL.
- **Dependencies:** REM-F-040-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Order of an untracked product with cost 40 x 3 posts COGS 120 to the configured account; Profitability COGS equals GL COGS for the period.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-040-05
- **Finding:** F-040
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-040
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-040-06
- **Finding:** F-040
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-040
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-040-07
- **Finding:** F-040
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-040 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Sell untracked and tracked cost-bearing products; assert GL COGS equals Profitability COGS
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-040-08
- **Finding:** F-040
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-040
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
