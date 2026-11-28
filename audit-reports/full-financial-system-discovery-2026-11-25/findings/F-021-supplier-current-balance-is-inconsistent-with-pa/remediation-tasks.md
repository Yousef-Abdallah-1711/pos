# F-021 — Implementation phases and tasks (all NOT STARTED)

Finding: Supplier current_balance is inconsistent with payables · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-021-01
- **Finding:** F-021
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-021
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-021-02
- **Finding:** F-021
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-021 at code level
- **Existing source paths:** backend/internal/repository/supplier_repo.go, backend/internal/repository/payment_acct_repo.go, backend/internal/repository/supplier_repo.go, backend/internal/handler/http/supplier_handler.go
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_r5_supplier_balance.mjs).
- **Dependencies:** REM-F-021-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-021-03
- **Finding:** F-021
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-021
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-021-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-021-04
- **Finding:** F-021
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-021
- **Existing source paths:** backend/internal/repository/supplier_repo.go, backend/internal/repository/payment_acct_repo.go, backend/internal/repository/supplier_repo.go, backend/internal/handler/http/supplier_handler.go
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Proposed (root cause CONFIRMED in RUN-005: the payment path subtracts from current_balance at payment_acct_repo.go:302-306 and no invoice path ever adds to it): 0) keep the replay in ph06_r5_supplier_balance.mjs as the failing test; 1) stop writing current_balance from supplier_repo.go (:44, :88, :116, :189) as an independent value and derive it from the AP subledger (sum of balance_due of approved invoices minus unapplied payments) in a view or query; 2) reconcile against GL account 2101 in a nightly check that raises an exception when they differ.
- **Dependencies:** REM-F-021-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** After any sequence of supplier invoices, payments and voids, the displayed supplier balance equals GL 2101 for that supplier. | The reconciliation check reports zero differences for the audit fixtures. | The column is no longer writable by application code.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-021-05
- **Finding:** F-021
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-021 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Supplier balance tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-021-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Balance = AP; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-021-06
- **Finding:** F-021
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-021
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_supplier_balance_view.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-021 no longer reproduce (expected: Balance = AP); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
