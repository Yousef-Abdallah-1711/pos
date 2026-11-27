# F-069 — Implementation phases and tasks (all NOT STARTED)

Finding: Completing a resignation never terminates the employee: turnover stays 0 and the person stays active and payable until the employee is deleted · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-069-01
- **Finding:** F-069
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-069
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

### REM-F-069-02
- **Finding:** F-069
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-069 at code level
- **Existing source paths:** backend/internal/repository/resignation_repo.go, backend/internal/repository/employee_repo.go, backend/internal/handler/http/hr_report_handler.go, backend/internal/repository/resignation_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-069-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-069-03
- **Finding:** F-069
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-069
- **Existing source paths:** backend/internal/repository/resignation_repo.go, backend/internal/repository/employee_repo.go, backend/internal/handler/http/hr_report_handler.go, backend/internal/repository/resignation_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: terminate the employee when a resignation completes. In backend/internal/repository/resignation_repo.go Complete (:201-215) update employees in the same transaction: employment_status = terminated, termination_date = last_working_date, is_active = false from that date; make DeleteEmployee (employee_repo.go:370-380) use the contract/resignation date when one exists; exclude terminated employees from future payroll and attendance creation; the turnover report then counts the departure.
- **Dependencies:** REM-F-069-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Completing a resignation with last working day 2026-11-25 sets termination_date 2026-11-25 and status terminated. | The turnover report shows departures 1 and the matching rate for the period. | No salary or attendance can be created for the employee after the termination date.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-069-04
- **Finding:** F-069
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-069
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

### REM-F-069-05
- **Finding:** F-069
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-069
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Employees with completed resignations still active are listed for HR review.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-069-06
- **Finding:** F-069
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-069 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Complete resignation -> employee terminated, departure counted, excluded from active payroll
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-069-07
- **Finding:** F-069
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-069
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
