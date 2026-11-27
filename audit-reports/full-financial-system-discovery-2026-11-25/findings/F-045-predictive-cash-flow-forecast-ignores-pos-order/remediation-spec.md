# F-045 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Predictive cash flow forecast ignores POS/order receipts and is built only from customer receipts and supplier payments.

## 2. Current behavior (verified)
2026-11 history cash_in 13,000 / cash_out 32,400 -> forecast Dec -18,430 while 319,660 of paid POS orders and GL cash +276,270 exist | RUN-003: PCF Nov cash_in 13,000/cash_out 32,400 vs paid POS orders 319,660; forecast negative (CONFIRMED accounting_advanced_repo.go:1105-1175). 

## 3. Required behavior
History reflects actual cash movements (POS cash/card receipts included)

## 4. Reproduction conditions
GET /reports/financial/predictive-cash-flow on tenant 4

## 5. Verified evidence
- evidence/database/RUN-003/PH-12/AUD-MET-PCF-VAT-INV.txt
- evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json
- tests/runtime/ph12_batch_apis.mjs

## 6. Root cause and confidence
CONFIRMED — accounting_advanced_repo.go:1105-1175 GetPredictiveCashFlow aggregates only payments_received and payments_made

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:1105-1175
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: response shape unchanged
- Accounting logic / posting rules: none named
- Reports affected: Financial hub predictive cash flow

## 8. Impact
- Financial: Forecast understates inflows by 319,660 in this dataset
- Operational: Financial hub predictive cash flow workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: build the forecast history from the ledger. In backend/internal/repository/accounting_advanced_repo.go GetPredictiveCashFlow (:1105-1175) read monthly cash inflows and outflows from the GL cash and bank accounts (debits and credits on the cash/bank accounts by month, excluding transfers between them) so POS and order receipts are included; keep customer receipts and supplier payments as forecast drivers for open invoices.

## 10. Alternative designs
A: GL cash/bank history (recommended). B: add orders to the existing two sources (double-counts receipts already posted to cash).

## 11. Required policy decisions
POL-04 (cash-flow method) in 56 (same method must feed both reports).

## 12. Required data migrations
No.

## 13. Historical data implications
None.

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
F-001, F-003 and shared-file findings F-009, F-016, F-043, F-047, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Forecast history equals GL cash movements per month (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For a month with 100,000 of POS receipts and 20,000 of supplier payments the historical series shows net +80,000 equal to the GL cash movement.
- The forecast changes when open invoices change.
