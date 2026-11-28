# F-073 — Implementation phases and tasks (all NOT STARTED)

Finding: Bank reconciliation lines accept a journal_entry_id that belongs to another tenant (only a foreign key is checked) · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-073-01
- **Finding:** F-073
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-073 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph11_r4_bank_recon.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-073-02
- **Finding:** F-073
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-073
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-073-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-073-03
- **Finding:** F-073
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-073
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Proposed: in AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) validate that the referenced journal entry belongs to the same tenant and restaurant (and, with F-072, touches the bank account); return 400/404 with a generic message; map foreign-key violations through the shared error mapper (F-025); enable row-level security on bank_reconciliation_sessions and bank_reconciliation_lines (pattern of migrations 164-171, see F-019).
- **Dependencies:** REM-F-073-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A journal id from another tenant returns 4xx and stores nothing. | A nonexistent id returns 4xx with no constraint name in the body. | RLS prevents a tenant GUC for tenant 4 from reading tenant 2 reconciliation rows.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-073-04
- **Finding:** F-073
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-073
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-073; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-073-05
- **Finding:** F-073
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-073 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Foreign-tenant and nonexistent journal ids return 4xx with no stored line and no raw database text
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-073-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400/404 because the journal is not visible to the tenant; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-073-06
- **Finding:** F-073
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-073
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_rls_bank_reconciliation.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-073 no longer reproduce (expected: 400/404 because the journal is not visible to the tenant); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
