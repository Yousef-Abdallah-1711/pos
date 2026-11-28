# F-060 — Implementation phases and tasks (all NOT STARTED)

Finding: Theme history endpoint returns 501 and the active theme can be deleted · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-060-01
- **Finding:** F-060
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-060 at code level
- **Existing source paths:** backend/internal/handler/http/admin_theme_handler.go
- **Proposed paths:** backend/migrations/<next>_theme_versions.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_theme_driver.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-060-02
- **Finding:** F-060
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-060
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_theme_versions.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-060-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-060-03
- **Finding:** F-060
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-060
- **Existing source paths:** backend/internal/handler/http/admin_theme_handler.go
- **Proposed paths:** backend/migrations/<next>_theme_versions.sql (new)
- **Required changes:** Proposed: implement theme history and protect the active theme. Replace the 501 stub at backend/internal/handler/http/admin_theme_handler.go:68 with a version list (theme_versions table: id, theme_id, config, saved_by, saved_at) and a restore action; in DeleteTheme (admin_theme_handler_v2.go:371) refuse to delete the active theme (409) or require choosing a replacement.
- **Dependencies:** REM-F-060-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Each save creates a version; restoring version N reproduces its config. | Deleting the active theme returns 409.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-060-04
- **Finding:** F-060
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-060 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_theme_versions.sql (new)
- **Required changes:** History returns versions; deleting the active theme is rejected
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-060-* test in regression-test-spec.md passes on the isolated stack, including the expected result: History list; deletion of the active theme refused or falls back to default; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-060-05
- **Finding:** F-060
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-060
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_theme_versions.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-060 no longer reproduce (expected: History list; deletion of the active theme refused or falls back to default); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
