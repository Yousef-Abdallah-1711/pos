# F-016 — Implementation phases and tasks (all NOT STARTED)

Finding: Payroll page default date window hides the current month's salary records · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-016-01
- **Finding:** F-016
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-016 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph20_evidence_gaps.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-016-02
- **Finding:** F-016
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-016
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: default the payroll list filter to the current pay period, not the last 30 days. In frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise start_date to the first day of the current month and end_date to the last day, and filter by overlap with the pay period (pay_period_start <= end and pay_period_end >= start) in the list query; show the active window in the page header; the HR cost report filter (pay_period_end < date_to, accounting_advanced_repo.go) should use the same overlap rule.
- **Dependencies:** REM-F-016-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A salary record for the current month is visible on first load without changing filters. | The HR cost report includes a pay period that straddles date_to. | The empty state never says "Create your first salary record" when records exist outside the window.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-016-03
- **Finding:** F-016
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-016
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/hr
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: GET /api/v1/hr/salaries: no change except overlap semantics for date filters. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-016-04
- **Finding:** F-016
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-016
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


## P8 Verification

### REM-F-016-05
- **Finding:** F-016
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-016 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** UI filter test
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-016-06
- **Finding:** F-016
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-016
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
