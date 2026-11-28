# F-057 — Implementation phases and tasks (all NOT STARTED)

Finding: Delivered cash-on-delivery orders stay unpaid and no settlement step records the cash collected by the driver · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-057-01
- **Finding:** F-057
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-057
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-057-02
- **Finding:** F-057
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-057 at code level
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/handler/http/admin_driver_handler.go, backend/internal/usecase/order_usecase.go
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_theme_driver.mjs).
- **Dependencies:** REM-F-057-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-057-03
- **Finding:** F-057
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-057
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-057-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-057-04
- **Finding:** F-057
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-057
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/handler/http/admin_driver_handler.go, backend/internal/usecase/order_usecase.go
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Proposed: cash-on-delivery settlement. Extend OrderUseCase.CompleteDelivery (backend/internal/usecase/order_usecase.go:433-447) so completing a COD delivery records the amount the driver collected (field, default order total) and creates a driver cash liability (Dr driver cash receivable / Cr revenue-and-VAT through the normal sales posting); add a driver settlement operation that records the cash handed in (Dr cash / Cr driver receivable) and reports variances per driver and day.
- **Dependencies:** REM-F-057-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Completing a COD order records the collected amount and posts the sales journal once. | Settling a driver for a day moves exactly the handed-in cash and shows any shortage. | Unsettled driver cash is visible in a report.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-057-05
- **Finding:** F-057
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-057
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-057-06
- **Finding:** F-057
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-057
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Existing delivered-unpaid COD orders are listed for reconciliation with drivers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing delivered-unpaid COD orders are listed for reconciliation with drivers.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-057-07
- **Finding:** F-057
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-057 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Delivered COD order can be settled once; cash and revenue reconcile
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-057-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-057-08
- **Finding:** F-057
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-057
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/driver_settlement_handler.go (new); backend/migrations/<next>_driver_cash_settlements.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-057 no longer reproduce (expected: A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
