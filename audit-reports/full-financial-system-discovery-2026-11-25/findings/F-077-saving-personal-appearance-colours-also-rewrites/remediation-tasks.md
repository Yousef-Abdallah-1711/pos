# F-077 — Implementation phases and tasks (all NOT STARTED)

Finding: Saving personal Appearance colours also rewrites the active public storefront theme, silently overwriting Theme Builder edits · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-077-01
- **Finding:** F-077
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-077 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph05_r4_appearance_theme.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-077-02
- **Finding:** F-077
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-077
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: remove syncBrandColorsToActiveTheme (frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95) from the personal Appearance save; keep the per-user light/dark preference only. If brand colours must stay in step with the public theme, expose them as one tenant-level setting inside the Theme Builder (permission: theme edit) and have the Appearance page link to it. The backend PUT /admin/themes/{id} should require the theme module WRITE permission (check the wrapWithPermission for the route) and record who changed the theme.
- **Dependencies:** REM-F-077-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Saving Appearance leaves the active themes_v2 row (config and updated_at) unchanged. | A colour saved in the Theme Builder persists after any Appearance save. | A user without theme-edit permission cannot change public colours through any screen.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-077-03
- **Finding:** F-077
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-077
- **Existing source paths:** frontend/apps/restaurant-website/app/[locale], frontend/apps/dashboard/src/app/[locale]/dashboard/settings
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: no API change; PUT /admin/themes/{id} is no longer called from Appearance. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** Saving Appearance leaves the active themes_v2 row (config and updated_at) unchanged.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-077-04
- **Finding:** F-077
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-077 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Saving Appearance leaves the active theme row unchanged; Theme Builder colours persist after an Appearance save
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-077-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A user preference changes only that user settings; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-077-05
- **Finding:** F-077
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-077
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-077 no longer reproduce (expected: A user preference changes only that user settings); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
