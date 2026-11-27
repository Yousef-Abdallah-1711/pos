# F-005 — Implementation phases and tasks (all NOT STARTED)

Finding: Card and wallet POS orders never become paid; stock is consumed anyway · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-005-01
- **Finding:** F-005
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-005
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-005-02
- **Finding:** F-005
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-005 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-005-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-005-03
- **Finding:** F-005
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-005
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-005-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-005-04
- **Finding:** F-005
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-005
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Proposed: define the tender lifecycle for non-cash POS orders. (1) In pos/page.tsx:302-308 keep cash settling immediately, but for card/wallet call a new "record payment" action that stores the provider reference and marks the order paid only after a confirmation step (manual confirmation by the cashier for terminal card payments, or a webhook for online wallets); (2) in the backend add POST /api/v1/admin/orders/{id}/payments (proposed) accepting method, amount and reference, validated against the order balance in one transaction (see F-006 locking); (3) release or reserve stock according to policy POL-05: stock is consumed at order creation today even if payment never arrives, so define an expiry that releases stock for unpaid orders (this reservation rule is not yet a numbered decision: proposed as POL-16).
- **Dependencies:** REM-F-005-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A card order that is confirmed shows payment_status paid, produces one outbox event and one sales journal debiting Card Payment Clearing (1113). | A card order that is never confirmed does not produce revenue and releases its stock after the configured expiry. | Recording the same payment twice (same reference) is idempotent.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-005-05
- **Finding:** F-005
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-005
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-005-06
- **Finding:** F-005
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-005
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Existing pending card/wallet orders (e.g. 8 orders, 88,000 in the evidence) need a one-off review: confirm or cancel each; do not auto-pay them.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-005-07
- **Finding:** F-005
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-005 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Card/wallet/split order tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-005-08
- **Finding:** F-005
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-005
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/order_payment_handler.go (new); backend/internal/jobs/expire_unpaid_orders.go (new, proposed)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
