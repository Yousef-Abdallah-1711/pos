# F-033 — Implementation phases and tasks (all NOT STARTED)

Finding: Employee role can export all orders and read billing/subscription data and HR records · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-033-01
- **Finding:** F-033
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-033
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-033-02
- **Finding:** F-033
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-033 at code level
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-033-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-033-03
- **Finding:** F-033
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared role/permission matrix to cover F-033 (core component delivered once by the P3 task of F-020)
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-033-specific cases. Design: Proposed: tighten the Employee role and per-endpoint levels together with F-020. Remove billing/subscription read and bulk order export from the default Employee permission set seeded at backend/internal/repository/auth_repo.go:380; guard the export routes in backend/cmd/api/main.go with the EXPORT level instead of READ; add a role-matrix test covering orders export, subscription/billing reads and HR records for each default role.
- **Dependencies:** F-020 P3 task; REM-F-033-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Employee gets 403 on /orders export, billing and subscription endpoints and on HR records beyond their own. | Owner/Manager behaviour is unchanged. | The matrix test in CI fails when a route lacks an explicit level.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-033-04
- **Finding:** F-033
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-033
- **Existing source paths:** backend/internal/repository/auth_repo.go, backend/cmd/api/main.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-033; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-033-05
- **Finding:** F-033
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-033 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** RBAC matrix tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-033-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Least privilege; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-033-06
- **Finding:** F-033
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-033
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-033 no longer reproduce (expected: Least privilege); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
