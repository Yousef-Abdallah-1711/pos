# F-076 — Implementation phases and tasks (all NOT STARTED)

Finding: Super Admin Users list links every row to /dashboard/users/{id}, a route that does not exist (404) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-076-01
- **Finding:** F-076
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-076 at code level
- **Existing source paths:** frontend/apps/super-admin/src/app/dashboard/users/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx
- **Proposed paths:** frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph13_r4_superadmin_ui.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-076-02
- **Finding:** F-076
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-076
- **Existing source paths:** frontend/apps/super-admin/src/app/dashboard/users/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx
- **Proposed paths:** frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)
- **Required changes:** Proposed: either add frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (user profile, memberships, status, last login, using the existing super-admin users API) or remove the View link at frontend/apps/super-admin/src/app/dashboard/users/page.tsx:132 until the page exists; add a link-check test that visits every internal href rendered on the page.
- **Dependencies:** REM-F-076-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Every link on the Users page returns 200. | The detail page shows the same email, status and role data as the list row.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-076-03
- **Finding:** F-076
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-076 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)
- **Required changes:** Every link rendered on the Users page resolves with 200
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-076-04
- **Finding:** F-076
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-076
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
