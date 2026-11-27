# F-031 — Implementation phases and tasks (all NOT STARTED)

Finding: Business-rule and dependency failures surface as wrong status codes and leak internals · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-031-01
- **Finding:** F-031
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-031 at code level
- **Existing source paths:** backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/errors.go (shared with F-025)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-031-02
- **Finding:** F-031
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared error-mapping layer to cover F-031 (core component delivered once by the P3 task of F-025)
- **Existing source paths:** backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/errors.go (shared with F-025)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-031-specific cases. Design: Proposed: status-code mapping rules, delivered together with the F-025 error mapper. Rules: business-rule violations 409/422 (paid-order cancel, over-settlement, closed period), dependency outages 503 with Retry-After (accounting pool or Redis down), missing resources 404, validation 400; the audited cases (paid-order cancel 500, receipt during outage 400 with a role error, DB down returning 403 TENANT_ACCESS_BLOCKED) each get an explicit mapping and a regression test; /health reports degraded with the failing component.
- **Dependencies:** F-025 P3 task; REM-F-031-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Each audited case returns the documented status and a generic message. | During a database outage authenticated calls return 503, not 403. | No response leaks pq or role names.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-031-03
- **Finding:** F-031
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-031
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/orders
- **Proposed paths:** backend/internal/handler/http/errors.go (shared with F-025)
- **Required changes:** API contract: status codes change for the listed cases (documented in 08-api-inventory.csv). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-031-04
- **Finding:** F-031
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-031 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/errors.go (shared with F-025)
- **Required changes:** Error-contract tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-031-05
- **Finding:** F-031
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-031
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/errors.go (shared with F-025)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
