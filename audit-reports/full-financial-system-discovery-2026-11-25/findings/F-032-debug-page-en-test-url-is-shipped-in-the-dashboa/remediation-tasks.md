# F-032 — Implementation phases and tasks (all NOT STARTED)

Finding: Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not behind the auth provider · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-032-01
- **Finding:** F-032
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-032 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx
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

### REM-F-032-02
- **Finding:** F-032
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-032
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: delete the debug route frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx from the shipped app (or move it behind a development-only guard that is compiled out in production builds), and add a build check that fails when a route named test/debug/dev exists in the production route manifest.
- **Dependencies:** REM-F-032-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** GET /en/test-url returns 404 in the production build. | The CI route-manifest check passes.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-032-03
- **Finding:** F-032
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-032
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-032; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-032-04
- **Finding:** F-032
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-032 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Route inventory check
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-032-* test in regression-test-spec.md passes on the isolated stack, including the expected result: No test pages in production build; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-032-05
- **Finding:** F-032
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-032
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-032 no longer reproduce (expected: No test pages in production build); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
