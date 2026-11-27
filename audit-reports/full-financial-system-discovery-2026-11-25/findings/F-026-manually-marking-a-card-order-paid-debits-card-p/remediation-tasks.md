# F-026 — Implementation phases and tasks (all NOT STARTED)

Finding: Manually marking a card order paid debits Card Payment Clearing with no clearing/settlement flow · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-026-01
- **Finding:** F-026
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-026
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-026-02
- **Finding:** F-026
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-026 at code level
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/domain/auto_posting.go, backend/migrations/156_seed_chart_of_accounts.sql, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph17_r5_observed_repro.mjs).
- **Dependencies:** REM-F-026-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-026-03
- **Finding:** F-026
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-026
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-026-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-026-04
- **Finding:** F-026
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-026
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/domain/auto_posting.go, backend/migrations/156_seed_chart_of_accounts.sql, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Proposed: add a settlement step for Card Payment Clearing (1113). (1) A "card settlement" operation (proposed POST /api/v1/accounting/card-settlements) records a bank deposit for a date range and amount, posting Dr bank account / Cr 1113 and, when fees are configured, Dr card-fee expense; (2) the clearing balance report lists unsettled card sales by order so the owner can match them to bank statements; (3) manual "mark paid" for card orders stays allowed only with a provider reference (see F-005).
- **Dependencies:** REM-F-026-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Settling a batch of card orders moves exactly the settled amount from 1113 to the bank account and leaves the unsettled remainder visible. | Fees reduce the deposit and post to the fee expense account. | A second settlement of the same orders is rejected.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-026-05
- **Finding:** F-026
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-026
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Clearing cleared on settlement) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-026-06
- **Finding:** F-026
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-026
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Existing 1113 balances need an opening reconciliation against bank statements.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing 1113 balances need an opening reconciliation against bank statements.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-026-07
- **Finding:** F-026
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-026 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Clearing reconciliation test
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-026-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Clearing cleared on settlement; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-026-08
- **Finding:** F-026
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-026
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/card_settlement_handler.go (new); backend/migrations/<next>_card_settlements.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-026 no longer reproduce (expected: Clearing cleared on settlement); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
