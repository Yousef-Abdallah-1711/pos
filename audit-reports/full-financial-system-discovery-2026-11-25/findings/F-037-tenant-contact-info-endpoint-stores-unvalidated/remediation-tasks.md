# F-037 — Implementation phases and tasks (all NOT STARTED)

Finding: Tenant contact-info endpoint stores unvalidated email/phone · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-037-01
- **Finding:** F-037
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-037 at code level
- **Existing source paths:** backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/tenant_contact_handler.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph05_settings_forms.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-037-02
- **Finding:** F-037
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-037 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/tenant_contact_handler.go (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-037-specific cases. Design: Proposed: move the inline handler at backend/cmd/api/main.go:845-873 to the handler layer, validate email (RFC format) and phone (E.164 using the registration validator), enforce length limits, and return 400 with field errors.
- **Dependencies:** F-062 P3 task; REM-F-037-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Invalid email or phone returns 400 and changes nothing. | Valid values persist and appear in the dashboard.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-037-03
- **Finding:** F-037
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-037
- **Existing source paths:** backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/tenant_contact_handler.go (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-037-04
- **Finding:** F-037
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-037 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/tenant_contact_handler.go (new)
- **Required changes:** API contract test: invalid email/phone => 400, no DB change
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-037-05
- **Finding:** F-037
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-037
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/tenant_contact_handler.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
