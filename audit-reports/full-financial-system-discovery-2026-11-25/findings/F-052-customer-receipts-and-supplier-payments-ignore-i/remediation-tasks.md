# F-052 — Implementation phases and tasks (all NOT STARTED)

Finding: Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the money twice · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-052-01
- **Finding:** F-052
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-052
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-052-02
- **Finding:** F-052
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-052 at code level
- **Existing source paths:** backend/internal/handler/http/payment_acct_handler.go, backend/internal/repository/payment_acct_repo.go, backend/internal/handler/http/order_cancellation_handler.go, backend/internal/handler/http/order_cancellation_handler.go
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph18_reliability.mjs).
- **Dependencies:** REM-F-052-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-052-03
- **Finding:** F-052
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-052
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-052-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-052-04
- **Finding:** F-052
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-052
- **Existing source paths:** backend/internal/handler/http/payment_acct_handler.go, backend/internal/repository/payment_acct_repo.go, backend/internal/handler/http/order_cancellation_handler.go, backend/internal/handler/http/order_cancellation_handler.go
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Proposed: add an idempotency layer for money-creating POST endpoints (receipts, supplier payments, journal create/post, orders, imports): table idempotency_keys(tenant_id, key, endpoint, request_hash, response_status, response_body, created_at) with a unique (tenant_id, key, endpoint); middleware reads the Idempotency-Key header already advertised by CORS, replays the stored response for an identical key+hash, returns 422 when the key is reused with a different body, and requires the key on the listed endpoints. Reuse the idempotent pattern already present in order_cancellation_handler.go:36.
- **Dependencies:** REM-F-052-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Two identical POST /accounting/payments/received with the same key create exactly one payment and return the same body/ids. | Same key with a different amount returns 422 and creates nothing. | Concurrent duplicates (10 parallel, same key) create exactly one payment.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-052-05
- **Finding:** F-052
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-052
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Second request returns the first result; one receipt) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-052-06
- **Finding:** F-052
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-052
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing records created under the defect keep wrong values. Produce a read-only exception report (query by th); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-052-07
- **Finding:** F-052
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-052 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Replay with same key returns same payment id; different key + same body is a new payment
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-052-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Second request returns the first result; one receipt; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-052-08
- **Finding:** F-052
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-052
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/middleware/idempotency.go (new); backend/migrations/<next>_create_idempotency_keys.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-052 no longer reproduce (expected: Second request returns the first result; one receipt); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
