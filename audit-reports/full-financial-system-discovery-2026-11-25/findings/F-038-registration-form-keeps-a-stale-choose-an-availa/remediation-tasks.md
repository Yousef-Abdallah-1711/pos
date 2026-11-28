# F-038 — Implementation phases and tasks (all NOT STARTED)

Finding: Registration form keeps a stale "choose an available subdomain" error after availability check succeeds · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-038-01
- **Finding:** F-038
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-038 at code level
- **Existing source paths:** frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, frontend/apps/landing/src/app/[locale]/auth/register/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph05_subdomain_probe.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-038-02
- **Finding:** F-038
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-038
- **Existing source paths:** frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, frontend/apps/landing/src/app/[locale]/auth/register/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed (root cause CONFIRMED in RUN-005 by source: nothing clears the manual slug error): in the landing register page (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, around lines 59-63) clear the slug error when availability becomes available (useEffect on availability.state) and either disable submit while the check is pending or queue the submit until the check resolves.
- **Dependencies:** REM-F-038-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** After the availability check returns available the error disappears within 200 ms and submit succeeds on the first click. | Submitting before the check completes does not show a stale error.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-038-03
- **Finding:** F-038
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-038 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Component test: submit during check then check resolves => no stale error
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-038-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Error clears (or submit proceeds) once availability is confirmed; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-038-04
- **Finding:** F-038
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-038
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-038 no longer reproduce (expected: Error clears (or submit proceeds) once availability is confirmed); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
