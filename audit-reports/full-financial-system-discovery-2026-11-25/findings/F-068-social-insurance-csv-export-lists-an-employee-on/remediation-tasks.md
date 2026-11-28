# F-068 — Implementation phases and tasks (all NOT STARTED)

Finding: Social-insurance CSV export lists an employee once per active insurance policy, repeating base salary and deductions · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-068-01
- **Finding:** F-068
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-068 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-068-02
- **Finding:** F-068
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-068
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: one row per salary record in the social-insurance export. In backend/internal/repository/accounting_advanced_repo.go GetSocialInsuranceReportRows (:1421-1446) replace the LEFT JOIN insurance_tracking at :1437 with a LATERAL subquery that returns one active policy per employee (latest start) or aggregates policy numbers with string_agg, so salary and deductions appear once per employee; add a test with two active policies.
- **Dependencies:** REM-F-068-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With two active policies for one employee the CSV has exactly one data row for that employee. | Totals of base_salary and deductions over the file equal the salaries table.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-068-03
- **Finding:** F-068
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-068
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Previously exported files may contain duplicates; list affected months.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Previously exported files may contain duplicates; list affected months.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-068-04
- **Finding:** F-068
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-068 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Two active policies -> one row
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-068-* test in regression-test-spec.md passes on the isolated stack, including the expected result: One row per employee salary record; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-068-05
- **Finding:** F-068
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-068
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-068 no longer reproduce (expected: One row per employee salary record); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
