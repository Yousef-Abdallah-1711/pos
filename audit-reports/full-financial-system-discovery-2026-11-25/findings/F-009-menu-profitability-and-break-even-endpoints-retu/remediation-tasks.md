# F-009 — Implementation phases and tasks (all NOT STARTED)

Finding: Menu Profitability and Break-even endpoints return HTTP 500 · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-009-01
- **Finding:** F-009
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-009
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-009-02
- **Finding:** F-009
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-009 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-009-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-009-03
- **Finding:** F-009
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-009
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
- **Required changes:** Proposed: change the two SQL statements in backend/internal/repository/accounting_advanced_repo.go at :976 and :1045 from p.cost_price to p.cost (the actual products column, see migrations and report_handler.go:2367 which already uses p.cost); map a missing-column/SQL error to a logged 500 with a generic message instead of leaking it; add a repository integration test that runs both queries against the migrated schema so a column rename fails CI.
- **Dependencies:** REM-F-009-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** GET /reports/financial/menu-profitability and /break-even return 200 with revenue and cost computed from products.cost for a fixture with 3 products. | The totals equal an independent SQL calculation over order_items for the same date range. | A schema-contract test fails if a referenced column does not exist.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-009-04
- **Finding:** F-009
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-009 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
- **Required changes:** Endpoint smoke tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-009-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reports render; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-009-05
- **Finding:** F-009
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-009
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/repository/accounting_advanced_repo_schema_test.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-009 no longer reproduce (expected: Reports render); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
