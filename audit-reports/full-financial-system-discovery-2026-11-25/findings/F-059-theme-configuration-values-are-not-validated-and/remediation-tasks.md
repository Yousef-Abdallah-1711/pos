# F-059 — Implementation phases and tasks (all NOT STARTED)

Finding: Theme configuration values are not validated and are served publicly (invalid colors and url(javascript:) accepted) · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-059-01
- **Finding:** F-059
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-059 at code level
- **Existing source paths:** backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_theme_driver.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-059-02
- **Finding:** F-059
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-059 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-059-specific cases. Design: Proposed: schema-validate theme JSON. In the theme save path (backend/internal/handler/http/admin_theme_handler_v2.go:293 and the theme service) validate every colour against a hex/rgb(a)/hsl pattern, fonts against an allow-list, URLs against http/https and the upload domain only (reject javascript:, data: and url() with other schemes), sizes within ranges, and reject unknown keys; sanitise again when rendering on the storefront.
- **Dependencies:** F-062 P3 task; REM-F-059-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** A colour "not-a-color" returns 400. | A value containing url(javascript:...) returns 400 and never reaches the public theme endpoint. | Valid themes round-trip unchanged.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-059-03
- **Finding:** F-059
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-059
- **Existing source paths:** frontend/apps/restaurant-website/app/[locale]
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** API contract: POST/PUT /api/v1/admin/themes: 400 with field errors. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** A colour "not-a-color" returns 400. (expected behaviour: Colors validated (hex/rgb/hsl allow-list); unsafe values rejected)
- **Risks:** low
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-059-04
- **Finding:** F-059
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-059
- **Existing source paths:** backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-059; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-059-05
- **Finding:** F-059
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-059 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** Invalid/unsafe color or CSS rejected 400; public theme contains only validated values
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-059-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Colors validated (hex/rgb/hsl allow-list); unsafe values rejected; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-059-06
- **Finding:** F-059
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-059
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/service/theme_validator.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-059 no longer reproduce (expected: Colors validated (hex/rgb/hsl allow-list); unsafe values rejected); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
