# F-045 — Implementation phases and tasks (all NOT STARTED)

Finding: Predictive cash flow forecast ignores POS/order receipts and is built only from customer receipts and supplier payments · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-045-01
- **Finding:** F-045
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-045
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

### REM-F-045-02
- **Finding:** F-045
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-045 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph12_batch_apis.mjs).
- **Dependencies:** REM-F-045-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-045-03
- **Finding:** F-045
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared cash-flow method to cover F-045 (core component delivered once by the P3 task of F-003)
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-045-specific cases. Design: Proposed: build the forecast history from the ledger. In backend/internal/repository/accounting_advanced_repo.go GetPredictiveCashFlow (:1105-1175) read monthly cash inflows and outflows from the GL cash and bank accounts (debits and credits on the cash/bank accounts by month, excluding transfers between them) so POS and order receipts are included; keep customer receipts and supplier payments as forecast drivers for open invoices.
- **Dependencies:** F-003 P3 task; REM-F-045-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** For a month with 100,000 of POS receipts and 20,000 of supplier payments the historical series shows net +80,000 equal to the GL cash movement. | The forecast changes when open invoices change.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-045-04
- **Finding:** F-045
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-045
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: response shape unchanged. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** For a month with 100,000 of POS receipts and 20,000 of supplier payments the historical series shows net +80,000 equal to the GL cash movement.
- **Risks:** low
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-045-05
- **Finding:** F-045
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-045
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (History reflects actual cash movements (POS cash/card receipts included)) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-045-06
- **Finding:** F-045
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-045 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Forecast history equals GL cash movements per month
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-045-* test in regression-test-spec.md passes on the isolated stack, including the expected result: History reflects actual cash movements (POS cash/card receipts included); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-045-07
- **Finding:** F-045
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-045
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-045 no longer reproduce (expected: History reflects actual cash movements (POS cash/card receipts included)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
