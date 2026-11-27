# F-021 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Supplier current_balance is inconsistent with payables.

## 2. Current behavior (verified)
current_balance -22000 vs AP 800

## 3. Required behavior
Balance = AP

## 4. Reproduction conditions
Approve invoice then pay

## 5. Verified evidence
- runtime/17
- evidence/database/RUN-005/PH-06/AUD-R5-F021-SUPPLIER-BALANCE.json
- tests/runtime/ph06_r5_supplier_balance.mjs

## 6. Root cause and confidence
CONFIRMED — suppliers.current_balance has exactly one money-movement writer: a supplier payment subtracts the amount (backend/internal/repository/payment_acct_repo.go:302-306), while supplier invoice creation and approval never add to it (no other UPDATE suppliers exists; supplier_repo.go:172-200 writes the value supplied by the handler, which preserves the stored value, supplier_handler.go:73). Replay on a new supplier: invoice 10,000 approved -> current_balance 0.00 while AP 2101 rose by 10,000; payments 4,000 and 6,000 -> current_balance -4,000.00 then -10,000.00 while the invoice balance_due went to 0 and AP returned to its prior 800.00. The stored balance therefore equals minus the sum of payments 

## 7. Affected paths
- Backend/database: `backend/internal/repository/supplier_repo.go`:44-189, `backend/internal/repository/payment_acct_repo.go`:302, `backend/internal/repository/supplier_repo.go`:172, `backend/internal/handler/http/supplier_handler.go`:73, `backend/internal/repository/payment_acct_repo.go`:302-306, `backend/internal/repository/supplier_repo.go`:172-200
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`
- API contracts: supplier responses keep current_balance (value now derived)
- Accounting logic / posting rules: accounts 2101
- Reports affected: none

## 8. Impact
- Financial: n/a (display)
- Operational: Suppliers workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause CONFIRMED in RUN-005: the payment path subtracts from current_balance at payment_acct_repo.go:302-306 and no invoice path ever adds to it): 0) keep the replay in ph06_r5_supplier_balance.mjs as the failing test; 1) stop writing current_balance from supplier_repo.go (:44, :88, :116, :189) as an independent value and derive it from the AP subledger (sum of balance_due of approved invoices minus unapplied payments) in a view or query; 2) reconcile against GL account 2101 in a nightly check that raises an exception when they differ.

## 10. Alternative designs
A: derived balance (recommended; removes the competing source). B: keep the column but update it inside the same transaction as every invoice/payment posting (still two sources). C: delete the column.

## 11. Required policy decisions
POL-12 (single source of truth) in 56 decides whether the subledger or the GL is authoritative.

## 12. Required data migrations
Yes: view or generated value; backfill/validate existing rows and list differences before switching.

## 13. Historical data implications
Existing suppliers with inconsistent balances are listed; accountant approves the opening correction.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-006, F-052, F-062, F-081

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/supplier_handler.go, backend/internal/repository/payment_acct_repo.go, backend/internal/repository/supplier_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Supplier balance tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- After any sequence of supplier invoices, payments and voids, the displayed supplier balance equals GL 2101 for that supplier.
- The reconciliation check reports zero differences for the audit fixtures.
- The column is no longer writable by application code.
