# F-017 — Implementation phases and tasks (all NOT STARTED)

Finding: Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order creation · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-017-01
- **Finding:** F-017
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-017 at code level
- **Existing source paths:** backend/internal/handler/http/order_cancellation_handler.go, backend/cmd/api/pos_routes.go, backend/internal/handler/http/public_order_handler.go, backend/internal/handler/http/public_order_handler.go
- **Proposed paths:** backend/migrations/<next>_order_idempotency.sql (new); backend/internal/repository/order_idempotency_repo.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-017-02
- **Finding:** F-017
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-017
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_order_idempotency.sql (new); backend/internal/repository/order_idempotency_repo.go (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-017-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-017-03
- **Finding:** F-017
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared idempotency store to cover F-017 (core component delivered once by the P3 task of F-052)
- **Existing source paths:** backend/internal/handler/http/order_cancellation_handler.go, backend/cmd/api/pos_routes.go, backend/internal/handler/http/public_order_handler.go, backend/internal/handler/http/public_order_handler.go
- **Proposed paths:** backend/migrations/<next>_order_idempotency.sql (new); backend/internal/repository/order_idempotency_repo.go (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-017-specific cases. Design: Proposed: make order creation idempotent. (1) Add a table order_idempotency(tenant_id, key, order_id, request_hash, created_at) with a unique (tenant_id, key); (2) in the POS order handler registered at backend/cmd/api/pos_routes.go and public_order_handler.go CreateOrder, read the Idempotency-Key header (the only current consumer is order_cancellation_handler.go:36), insert the key inside the same transaction as the order, and on a repeated key return the original order with 200 and the same body; a key reused with a different payload returns 422; (3) the POS frontend generates one key per ca
- **Dependencies:** F-052 P3 task; REM-F-017-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Posting the same cart twice with the same Idempotency-Key creates exactly one order and decrements stock once. | A different payload with the same key returns 422 and creates nothing. | Five concurrent identical submissions create one order.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-017-04
- **Finding:** F-017
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-017 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_order_idempotency.sql (new); backend/internal/repository/order_idempotency_repo.go (new)
- **Required changes:** Duplicate submit tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-017-05
- **Finding:** F-017
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-017
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_order_idempotency.sql (new); backend/internal/repository/order_idempotency_repo.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
