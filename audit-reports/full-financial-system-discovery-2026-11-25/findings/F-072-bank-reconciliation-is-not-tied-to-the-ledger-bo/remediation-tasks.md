# F-072 — Implementation phases and tasks (all NOT STARTED)

Finding: Bank reconciliation is not tied to the ledger: book balance is the sum of client-supplied line amounts, so a session can be "completed" with invented numbers and stamps the bank account as reconciled · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-072-01
- **Finding:** F-072
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-072 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph11_r4_bank_recon.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-072-02
- **Finding:** F-072
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-072
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-072-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-072-03
- **Finding:** F-072
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-072
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Proposed: rebuild reconciliation on ledger data. (1) Session creation reads the book balance from the GL for the bank account (general_ledger rows of bank_accounts_ledger.gl_account_id up to the statement date) and stores it; the client no longer supplies book_balance; (2) AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) takes the journal line id, verifies it is a line on the bank GL account in the same tenant/restaurant, takes statement/book amounts from that line (the client supplies only the statement amount to match) and rejects a line already matched in another session; (3) refreshReconciliationSummaryTx (:2126-2149) computes difference = statement balance - (ledger balance adjusted by outstanding items) from those matched lines; (4) CompleteBankReconciliationSession (:169-230) completes only when the difference is zero (or an explicit adju
- **Dependencies:** REM-F-072-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A session with a line whose journal has no bank-account line is rejected (400). | Completing a session whose difference is not zero returns 409, leaves status in_progress and completed_at null. | A session matched to real ledger lines completes and sets last_reconciled_balance equal to the ledger balance at the statement date. | The book balance shown equals the GL balance of the bank account (see F-029).
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-072-04
- **Finding:** F-072
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-072
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts ) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-072-05
- **Finding:** F-072
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-072
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Existing completed sessions must be re-verified against the ledger; list sessions whose book balance differs from the GL at the statement date.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing completed sessions must be re-verified against the ledger; list sessions whose book balance differs f); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-072-06
- **Finding:** F-072
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-072 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Session with invented amounts cannot complete; matched real lines complete and update last_reconciled values; non-zero difference returns 409
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-072-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts ; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-072-07
- **Finding:** F-072
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-072
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_bank_reconciliation_line_constraints.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-072 no longer reproduce (expected: Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts ); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
