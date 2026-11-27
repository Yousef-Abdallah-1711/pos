# F-042 — Implementation phases and tasks (all NOT STARTED)

Finding: Order stats "total revenue" includes cancelled and unpaid orders and "delivered today" counts every order created today · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-042-01
- **Finding:** F-042
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-042
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

### REM-F-042-02
- **Finding:** F-042
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-042 at code level
- **Existing source paths:** backend/internal/repository/order_repo.go, backend/internal/repository/order_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph12_batch_apis.mjs).
- **Dependencies:** REM-F-042-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-042-03
- **Finding:** F-042
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared shared revenue definitions to cover F-042 (core component delivered once by the P3 task of F-012)
- **Existing source paths:** backend/internal/repository/order_repo.go, backend/internal/repository/order_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-042-specific cases. Design: Proposed: make order statistics follow the approved revenue basis. In backend/internal/repository/order_repo.go GetOrderStats (:859-870) filter total revenue by status <> cancelled and payment_status = paid (or the basis approved in POL-01) and expose both values (booked vs collected); compute delivered_today from status = delivered and a delivered_at timestamp (add the column if missing) instead of COUNT(created_at >= CURRENT_DATE); use the restaurant time zone for the day boundary (see F-070).
- **Dependencies:** F-012 P3 task; REM-F-042-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** For a dataset with cancelled and unpaid orders the stats revenue equals the approved basis computed independently in SQL. | delivered_today counts only orders delivered today. | The dashboard tiles use the same endpoint and show the same numbers.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-042-04
- **Finding:** F-042
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-042
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Revenue excludes cancelled (and unpaid per chosen policy); delivered_today counts delivered orders) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-042-05
- **Finding:** F-042
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-042 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Stats with cancelled/unpaid/delivered orders must equal independent counts
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-042-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Revenue excludes cancelled (and unpaid per chosen policy); delivered_today counts delivered orders; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-042-06
- **Finding:** F-042
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-042
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-042 no longer reproduce (expected: Revenue excludes cancelled (and unpaid per chosen policy); delivered_today counts delivered orders); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
