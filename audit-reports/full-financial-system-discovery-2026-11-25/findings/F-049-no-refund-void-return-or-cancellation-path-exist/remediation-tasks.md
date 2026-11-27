# F-049 — Implementation phases and tasks (all NOT STARTED)

Finding: No refund, void, return or cancellation path exists for a paid order (missing financial capability) · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-049-01
- **Finding:** F-049
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-049
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-049-02
- **Finding:** F-049
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-049 at code level
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/repository/order_repo.go, backend/internal/domain/order.go, backend/migrations/175_durable_payment_refund_state.sql
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph19_paid_cancel.mjs).
- **Dependencies:** REM-F-049-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-049-03
- **Finding:** F-049
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-049
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-049-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-049-04
- **Finding:** F-049
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-049
- **Existing source paths:** backend/internal/usecase/order_usecase.go, backend/internal/repository/order_repo.go, backend/internal/domain/order.go, backend/migrations/175_durable_payment_refund_state.sql
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Proposed: design and build a refund/void/return workflow for paid orders (the largest missing capability). Remove the blanket rejection at usecase/order_usecase.go:341 and repository/order_repo.go:634 only behind the new workflow: (1) a refund request document (order, items/quantities, reason, tender, approver) with partial and full refunds; (2) approval gated by a new RBAC action; (3) one transaction that records the refund, returns stock for tracked products at their original cost, and posts reversing journals (Dr revenue and VAT payable, Cr cash/bank/card clearing; Dr inventory / Cr COGS for returned goods); (4) idempotency key on the refund; (5) the order moves to refunded or partially_refunded, reports exclude refunded amounts per POL-06, and KDS/stock are updated.
- **Dependencies:** REM-F-049-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A full refund of a paid cash order returns the cash, reverses revenue/VAT/COGS exactly (sum of journals for the order = 0) and restores stock. | A partial refund of one item posts proportional reversals. | The same refund request twice posts once; concurrent refunds cannot exceed the paid amount. | The P&L, balance sheet and sales reports reflect the refund per POL-06.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-049-05
- **Finding:** F-049
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-049
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-049-06
- **Finding:** F-049
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-049 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Full and partial refund per tender type reconcile revenue, VAT, COGS, inventory, cash/clearing
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-049-07
- **Finding:** F-049
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-049
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/order_refund_handler.go (new); backend/internal/usecase/order_refund_usecase.go (new); backend/migrations/<next>_order_refunds.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
