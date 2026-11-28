# F-019 — Implementation phases and tasks (all NOT STARTED)

Finding: Application DB role sees all tenants on tables without RLS and can alter paid orders/payment amounts · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-019-01
- **Finding:** F-019
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-019 at code level
- **Existing source paths:** backend/migrations/164_rls_user_auth_audit_logs.sql, backend/migrations/171_rls_e_invoice_submissions.sql
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-019-02
- **Finding:** F-019
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-019
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-019-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-019-03
- **Finding:** F-019
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-019
- **Existing source paths:** backend/migrations/164_rls_user_auth_audit_logs.sql, backend/migrations/171_rls_e_invoice_submissions.sql
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Proposed: extend row-level security and add immutability where money is stored. (1) New migrations following the pattern of migrations 164-171 enable and force RLS (tenant_id = current_setting(app.tenant_id)) on orders, order_items, payments, invoices_customer/supplier, payments_received/made, chart_of_accounts, employees, salaries and the other tenant tables (57 of 70 are unprotected), with the application setting the tenant GUC per request; (2) revoke UPDATE/DELETE on payments and paid orders from pos_app and route corrections through reversing records; (3) add a trigger that rejects changes to total_amount/payment fields once an order is paid (same style as cfr002 triggers on the ledger).
- **Dependencies:** REM-F-019-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With app.tenant_id set to tenant A, SELECT on every tenant table returns only A rows; without the setting it returns none. | UPDATE of orders.total_amount or payments.amount for a paid order by pos_app fails; the reversing-record path works. | The full existing test suite passes with RLS enabled (no missing tenant GUC).
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-019-04
- **Finding:** F-019
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-019
- **Existing source paths:** backend/migrations/164_rls_user_auth_audit_logs.sql, backend/migrations/171_rls_e_invoice_submissions.sql
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-019; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-019-05
- **Finding:** F-019
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-019 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Restricted-role isolation suite
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-019-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Cross-tenant rows hidden; paid records immutable; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-019-06
- **Finding:** F-019
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-019
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_rls_core_money_tables.sql (new); backend/internal/middleware/tenant_guc.go (new, proposed)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-019 no longer reproduce (expected: Cross-tenant rows hidden; paid records immutable); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
