# F-050 — Implementation phases and tasks (all NOT STARTED)

Finding: Business-rule rejection when cancelling a paid order via PUT /admin/orders/{id}/status returns HTTP 500 instead of a 4xx · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-050-01
- **Finding:** F-050
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-050 at code level
- **Existing source paths:** backend/internal/handler/http/admin_order_handler.go, backend/internal/handler/http/admin_order_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph19_paid_cancel.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-050-02
- **Finding:** F-050
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared error-mapping layer to cover F-050 (core component delivered once by the P3 task of F-025)
- **Existing source paths:** backend/internal/handler/http/admin_order_handler.go, backend/internal/handler/http/admin_order_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-050-specific cases. Design: Proposed: map domain guard errors to client errors. Introduce sentinel errors (for example ErrPaidOrderCannotBeCancelled) in the usecase layer and translate them in handler/http/admin_order_handler.go:410 to 409 with a stable code; keep 500 only for unexpected errors; add table-driven tests for every guard.
- **Dependencies:** F-025 P3 task; REM-F-050-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Cancelling a paid order via PUT /admin/orders/{id}/status returns 409 with code ORDER_PAID_CANNOT_CANCEL. | Unexpected failures still return 500 with a request id.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-050-03
- **Finding:** F-050
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-050 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Paid-order cancel through every endpoint returns the same 4xx
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-050-04
- **Finding:** F-050
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-050
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
