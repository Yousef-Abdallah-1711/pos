# F-052 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the money twice.

## 2. Current behavior (verified)
Both return 201; payments_received has 2 rows (sum 800); invoice amount_paid 800 / balance 200; AR credited twice | RUN-003 REL-025: identical Idempotency-Key replay creates 2 receipts (sum 800) (REPRODUCED). 

## 3. Required behavior
Second request returns the first result; one receipt

## 4. Reproduction conditions
POST /accounting/payments/received twice with the same UUID Idempotency-Key and body (400 on a 1,000 invoice)

## 5. Verified evidence
- evidence/database/RUN-003/PH-18/AUD-REL-BATTERY.json
- tests/runtime/ph18_reliability.mjs

## 6. Root cause and confidence
CONFIRMED — the only Idempotency-Key consumer in the backend is order_cancellation_handler.go (item cancel); payment_acct handler/repository/usecase contain no key handling although the CORS allow-list advertises the header

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/payment_acct_handler.go`, `backend/internal/repository/payment_acct_repo.go`, `backend/internal/handler/http/order_cancellation_handler.go`:36, `backend/internal/handler/http/order_cancellation_handler.go`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: All money-creating POSTs: header Idempotency-Key (UUID) required; documented in 08-api-inventory.csv
- Accounting logic / posting rules: none named
- Reports affected: Receipts; Supplier payments; Accounting (cash, AR/AP)

## 8. Impact
- Financial: Every retried or double-clicked payment duplicates cash and AR/AP movement
- Operational: Receipts; Supplier payments; Accounting (cash, AR/AP) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: add an idempotency layer for money-creating POST endpoints (receipts, supplier payments, journal create/post, orders, imports): table idempotency_keys(tenant_id, key, endpoint, request_hash, response_status, response_body, created_at) with a unique (tenant_id, key, endpoint); middleware reads the Idempotency-Key header already advertised by CORS, replays the stored response for an identical key+hash, returns 422 when the key is reused with a different body, and requires the key on the listed endpoints. Reuse the idempotent pattern already present in order_cancellation_handler.go:36.

## 10. Alternative designs
A: generic middleware + table (recommended). B: per-endpoint natural keys (reference numbers) - weaker for payments with free-text references. C: client-side double-submit prevention only (not sufficient for retries).

## 11. Required policy decisions
Decide and record in 56-accounting-policy-decision-package.md before implementation: the accounting treatment of the affected transaction

## 12. Required data migrations
Yes — a numbered migration after the current highest in backend/migrations (data-correction + constraint/trigger change). Existing rows must be reconciled first.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-017 and shared-file findings F-006, F-017, F-021

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/order_cancellation_handler.go, backend/internal/handler/http/payment_acct_handler.go, backend/internal/repository/payment_acct_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Replay with same key returns same payment id; different key + same body is a new payment (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Two identical POST /accounting/payments/received with the same key create exactly one payment and return the same body/ids.
- Same key with a different amount returns 422 and creates nothing.
- Concurrent duplicates (10 parallel, same key) create exactly one payment.
