# F-027 — Implementation phases and tasks (all NOT STARTED)

Finding: Employee-role users get finance, HR and RBAC page shells in the UI (no UI guard) while the API returns 403 · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-027-01
- **Finding:** F-027
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-027
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-027-02
- **Finding:** F-027
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-027 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-027-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-027-03
- **Finding:** F-027
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared role/permission matrix to cover F-027 (core component delivered once by the P3 task of F-020)
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-027-specific cases. Design: Proposed: one route-guard mechanism. (1) Add a shared guard wrapper for dashboard pages that checks the same module and level the API enforces (the permission hook already exists and is used by 6 page files); (2) apply it to every page under accounting, hr, settings/roles, reports and billing, driven by a single route-to-permission table kept next to the navigation definition so the menu and the guard cannot diverge; (3) unauthorised users get an access-denied page, never the page shell.
- **Dependencies:** F-020 P3 task; REM-F-027-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** For the Employee role every finance/HR/RBAC route renders the access-denied page and issues no data requests. | The menu hides exactly the routes the guard would block. | A test enumerates all dashboard routes and fails when a route is missing from the permission table.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-027-04
- **Finding:** F-027
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-027
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-027; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-027-05
- **Finding:** F-027
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-027 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** Page-permission matrix tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-027-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Page guard consistent with API permissions; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-027-06
- **Finding:** F-027
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-027
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/lib/routePermissions.ts (new); frontend/apps/dashboard/src/components/RouteGuard.tsx (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-027 no longer reproduce (expected: Page guard consistent with API permissions); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
