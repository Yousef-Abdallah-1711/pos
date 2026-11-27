# F-047 — Implementation phases and tasks (all NOT STARTED)

Finding: Recurring journal template can be run repeatedly for the same period and posts a duplicate journal each time · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-047-01
- **Finding:** F-047
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-047
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-047-02
- **Finding:** F-047
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-047 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph11_accounting_flow2.mjs).
- **Dependencies:** REM-F-047-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-047-03
- **Finding:** F-047
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-047
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-047-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-047-04
- **Finding:** F-047
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared idempotency store to cover F-047 (core component delivered once by the P3 task of F-052)
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-047-specific cases. Design: Proposed: make recurring journals idempotent per period. In backend/internal/repository/accounting_advanced_repo.go RunRecurringJournalTemplate (:396-420) lock the template row (SELECT ... FOR UPDATE), reject a run when runDate < next_run_date unless an explicit force flag with a reason is supplied, store a unique key (template_id, period) on the generated journal and advance next_run_date/last_run_date in the same transaction.
- **Dependencies:** F-052 P3 task; REM-F-047-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Running the same template twice for the same period creates one journal; the second call returns 409 with the existing journal id. | Five concurrent runs create one journal. | A forced run records the reason and creates a second journal only when explicitly requested.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-047-05
- **Finding:** F-047
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-047
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-047-06
- **Finding:** F-047
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-047
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Duplicate recurring journals in the evidence (2 runs, 2 journals) are reversed by approved journals.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-047-07
- **Finding:** F-047
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-047 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Double-click / concurrent run creates exactly one journal per period
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-047-08
- **Finding:** F-047
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-047
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_recurring_journal_unique_period.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
