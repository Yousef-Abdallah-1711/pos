# F-006 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Concurrent customer receipts over-settle an invoice.

## 2. Current behavior (verified)
5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment.  | RUN-003 FRESH-INVOICE RE-REPRODUCTION (AUD-CINV-FLOW): invoice 3,000, 5 concurrent receipts of 1,000 -> all 5 accepted (201), amount_paid 6,000, AR credit -3,000; single sequential overpayment 5,000 is rejected (guard exists but is not serialized). RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment. 

## 3. Required behavior
At most 3 accepted; no over-collection

## 4. Reproduction conditions
5 parallel receipts of 1000 on balance 3200

## 5. Verified evidence
- runtime/16 | evidence/database/RUN-003/PH-09/AUD-AR-RESIDUAL-trace.txt | evidence/database/RUN-003/PH-12/ | evidence/database/RUN-003/PH-09/AUD-CINV-FLOW.json | tests/runtime/ph09_customer_flow.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/payment_acct_repo.go:445-449 customer-invoice balance SELECT has no FOR UPDATE (supplier twin :499-503 has FOR UPDATE OF inv); applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0), hiding the overpayment

## 7. Affected paths
- Backend/database: `backend/internal/repository/payment_acct_repo.go`:445-449, `backend/internal/repository/payment_acct_repo.go`:569-582, `backend/internal/repository/payment_acct_repo.go`:499-503
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1120
- Reports affected: Customer invoices; Receipts

## 8. Impact
- Financial: 1800
- Operational: Customer invoices; Receipts workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: in PaymentAccountingRepository.CreatePaymentReceived (backend/internal/repository/payment_acct_repo.go:29) lock the invoice row with SELECT ... FOR UPDATE inside the same transaction (mirror payment_acct_repo.go:499-503 for suppliers), compare amount with the locked balance_due and reject with a 409 'payment amount exceeds invoice balance'; make applyCustomerInvoicePaymentTx (:569-582) conditional (UPDATE ... WHERE balance_due >= amount) and check RowsAffected; remove the silent GREATEST clamp. Add a CHECK (amount_paid <= total_amount) as a safety net after the data pre-check.

## 10. Alternative designs
A: row lock + conditional update (recommended). B: serializable transaction with retry (heavier). C: accept the overpayment as customer credit (needs policy POL-10 and a credit account).

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
F-008 and shared-file findings F-021, F-052

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/payment_acct_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Concurrent settlement tests both sides (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- 5 parallel receipts of 1,000 on a 2,000 balance: exactly 2 return 201, 3 return 409; invoice amount_paid = 2,000 and balance 0.
- GL 1120 for the invoice equals invoice total minus receipts (no negative AR).
- Sequential overpayment still returns 400/409 and creates no receipt.
