# F-028 — Implementation phases and tasks (all NOT STARTED)

Finding: Concurrent oversell attempts leave orphan pending orders with items although the API returned 400 · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-028-01
- **Finding:** F-028
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-028
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-028-02
- **Finding:** F-028
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-028 at code level
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/usecase/order_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-028-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-028-03
- **Finding:** F-028
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-028
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/usecase/order_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: make order creation and stock reservation one transaction. In backend/internal/usecase/order_usecase.go CreateOrderWithItems (:168) move the stock reservation/confirmation (currently :201-205, after the commit) inside the same database transaction that inserts the order and items, locking the product rows (SELECT ... FOR UPDATE ordered by product id to avoid deadlocks); on insufficient stock roll back everything so no pending order remains; return 409 with the available quantity.
- **Dependencies:** REM-F-028-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** N concurrent orders for stock S < N create exactly S orders (or S units) and no pending order without a reservation; every rejected request returns 409/400 and leaves no rows. | Sum of order quantities = stock decrease = ledger movement after the burst. | No deadlock errors under 100 concurrent submissions of mixed multi-item carts.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-028-04
- **Finding:** F-028
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-028 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Concurrent oversell tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-028-05
- **Finding:** F-028
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-028
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
