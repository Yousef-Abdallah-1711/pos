# F-024 — Implementation phases and tasks (all NOT STARTED)

Finding: Frontend sends float artifacts (630.0000000000001) in monetary payloads · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-024-01
- **Finding:** F-024
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-024
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-024-02
- **Finding:** F-024
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-024 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx, frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx, frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-024-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-024-03
- **Finding:** F-024
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared discount-first tax calculation to cover F-024 (core component delivered once by the P3 task of F-007)
- **Existing source paths:** frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx, frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx, frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-024-specific cases. Design: Proposed: compute money in integer minor units (cents) in frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx (calculateLineValues :620-635 and mapLineToPayload :637-650): convert unit price and quantity to cents once, apply discount and tax with an explicit rounding rule (half-up to 2 decimals, the same rule the backend trigger uses), and send values rounded to two decimals; better, send only quantity, unit price, discount rate and tax rate and let the server return the authoritative amounts for display.
- **Dependencies:** F-007 P3 task; REM-F-024-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** The payload for 5,000 x 1 with 10% discount and 14% VAT contains tax_amount 630 and line_total 5130 exactly (no 630.0000000000001). | For 1,000 random lines the UI total equals the server-stored total to the cent.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-024-04
- **Finding:** F-024
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-024
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-024-05
- **Finding:** F-024
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-024
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Existing rows with float artefacts are rounded in the numeric(…,2) columns already; list any that differ from the rounded value.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-024-06
- **Finding:** F-024
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-024 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Precision tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-024-07
- **Finding:** F-024
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-024
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/lib/money.ts (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
