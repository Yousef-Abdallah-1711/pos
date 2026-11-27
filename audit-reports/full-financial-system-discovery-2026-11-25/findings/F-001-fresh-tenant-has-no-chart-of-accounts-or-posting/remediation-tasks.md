# F-001 — Implementation phases and tasks (all NOT STARTED)

Finding: Fresh tenant has no chart of accounts or posting rules; first paid order is never posted · Severity Critical (release-blocking) · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-001-01
- **Finding:** F-001
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-001
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-001-02
- **Finding:** F-001
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-001 at code level
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/internal/handler/http/auth_handler.go, backend/internal/usecase/auth_usecase.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph05_register_headed.mjs).
- **Dependencies:** REM-F-001-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-001-03
- **Finding:** F-001
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-001
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-001-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-001-04
- **Finding:** F-001
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-001
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/internal/handler/http/auth_handler.go, backend/internal/usecase/auth_usecase.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Proposed: perform tenant accounting bootstrap inside the same transaction as registration (backend/internal/repository/auth_repo.go CreateTenantWithUserContext, before commit): insert the default chart of accounts (62 accounts as in backend/migrations/156_seed_chart_of_accounts.sql), the six auto_posting_rules (152_create_auto_posting_rules.sql), the current fiscal_period and a default tax configuration for the new tenant+restaurant. Extract the seed into one reusable function (proposed backend/internal/repository/accounting_bootstrap.go) so a repair command can re-run it idempotently for existing tenants (2 and 3 in the audit database).
- **Dependencies:** REM-F-001-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A tenant registered through the UI (POST /auth/register 201) has chart_of_accounts >= 62, auto_posting_rules = 6, fiscal_periods >= 1 and tax configuration without any manual step. | First POS cash order paid by that tenant produces its sales and COGS journals within 10 s (outbox completed, attempts = 1). | Registering 20 tenants concurrently leaves each with exactly one bootstrap set (no duplicates, no cross-tenant rows). | Re-running the repair function on an already seeded tenant changes nothing.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-001-05
- **Finding:** F-001
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-001
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-001-06
- **Finding:** F-001
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-001
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/internal/handler/http/auth_handler.go, backend/internal/usecase/auth_usecase.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-001-07
- **Finding:** F-001
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-001
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-001-08
- **Finding:** F-001
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-001 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** New-tenant first-order posting test; onboarding smoke
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-001-09
- **Finding:** F-001
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-001
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/repository/accounting_bootstrap.go (new); backend/internal/repository/accounting_bootstrap_integration_test.go (new); backend/migrations/<next>_backfill_tenant_accounting_bootstrap.sql (optional backfill for unseeded tenants)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
