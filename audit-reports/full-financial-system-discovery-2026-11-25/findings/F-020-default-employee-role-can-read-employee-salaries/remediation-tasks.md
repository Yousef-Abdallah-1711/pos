# F-020 — Implementation phases and tasks (all NOT STARTED)

Finding: Default Employee role can read employee salaries but cannot use the POS; no built-in cashier role · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-020-01
- **Finding:** F-020
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-020
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-020-02
- **Finding:** F-020
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-020 at code level
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go, backend/internal/domain/hr.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-020-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-020-03
- **Finding:** F-020
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-020
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go, backend/internal/domain/hr.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** Proposed: define built-in roles and stop leaking salary fields. (1) Seed a Cashier role (POS write, orders read/write, no HR, no accounting) next to the existing Employee role in auth_repo.go seedDefaultRolesAndPermissions (:380); (2) serialise base_salary and other pay fields in GET /hr/employees only when the caller has SALARIES READ (domain/hr.go:44 field becomes conditional via a response DTO); (3) change the UI guards so the POS entry follows the Cashier role.
- **Dependencies:** REM-F-020-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A user with the Employee role receives no base_salary or bank fields from GET /hr/employees. | A Cashier can create and pay POS orders and gets 403 on finance and HR endpoints. | RBAC matrix test: module x level for Owner, Manager, Cashier, Employee matches the approved matrix.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-020-04
- **Finding:** F-020
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-020
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** API contract: GET /api/v1/hr/employees: pay fields omitted unless SALARIES READ (breaking for clients that relied on them; coordinate with dashboard). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-020-05
- **Finding:** F-020
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-020
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go, backend/internal/domain/hr.go, backend/internal/repository/auth_repo.go
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-020-06
- **Finding:** F-020
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-020 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** RBAC matrix tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-020-07
- **Finding:** F-020
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-020
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/employee_dto.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
