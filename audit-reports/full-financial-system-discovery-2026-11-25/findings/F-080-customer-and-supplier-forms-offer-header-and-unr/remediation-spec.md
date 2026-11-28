# F-080 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Customer and supplier forms offer header and unrelated accounts as the receivable/payable account, the backend accepts them, and later invoicing fails with "account N is a header account".

## 2. Current behavior (verified)
Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer returned 400 {"error":"account 63 is a header account and cannot be posted to"}, so the customer can never be invoiced until the account is edited

## 3. Required behavior
Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account

## 4. Reproduction conditions
Headed Chrome tenant 4: open the customer form; the Receivable account dropdown lists 1000 - Assets, 1100 - Current Assets, Petty Cash, Cash on Hand, Bank accounts and others; choose the first account, save; then create a customer invoice for that customer (the form matrix did this for customers 6-13 and suppliers 6, 7)

## 5. Verified evidence
- evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-RESULT.json
- evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-DUMP.json
- tests/runtime/ph09_r5_invoice_form.mjs

## 6. Root cause and confidence
CONFIRMED — CustomerForm.tsx:61-62 filters the dropdown only by account_type === "asset" (the supplier form uses the same pattern), the customer handler passes ReceivableAccountID straight to the repository (backend/internal/handler/http/customer_acct_handler.go:141,244) without checking that it is a postable receivable account, and the first posting attempt fails at backend/internal/repository/journal_entry_repo.go:708

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`:61, `backend/internal/handler/http/customer_acct_handler.go`:141, `backend/internal/repository/journal_entry_repo.go`:708, `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`:61-62
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST/PUT /api/v1/accounting/customers and /api/v1/suppliers: 400 for invalid account ids (additive validation)
- Accounting logic / posting rules: none named
- Reports affected: Accounting > Customers, Suppliers, Customer invoices, Payments

## 8. Impact
- Financial: 0 direct; blocks invoicing and payments for affected customers and suppliers
- Operational: Accounting > Customers, Suppliers, Customer invoices, Payments workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: (1) in frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx:61-62 (and the equivalent supplier form) filter the options to active, non-header accounts of the receivable (payable) subtype instead of every account_type "asset" (liability for suppliers); (2) validate server-side in backend/internal/handler/http/customer_acct_handler.go (create :141, update :244) and the supplier handler that the account exists in the tenant, is active, is not a header (the same rule as journal_entry_repo.go:708) and has the expected subtype; return 400 with a field error; (3) default the field to the tenant AR (1120) or AP (2101) account when left empty, so the invoice path never falls back to a header account.

## 10. Alternative designs
A: filter the picker, validate server-side and default to the configured AR/AP account (recommended). B: server validation only (users still see invalid options).

## 11. Required policy decisions
Which account subtype is the receivable/payable control account is an accounting decision (proposed: subtype receivable/payable from the chart).

## 12. Required data migrations
No schema change; data correction for existing rows pointing at header accounts (customers 6-13 and suppliers 6, 7 in the evidence) after listing them.

## 13. Historical data implications
List customers and suppliers whose control account is a header or non-receivable account.

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
F-062, F-075 and shared-file findings F-013, F-074, F-081

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/customer_acct_handler.go, backend/internal/repository/journal_entry_repo.go, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Creating a customer or supplier with a header or non-receivable account returns 400; the dropdown never lists header accounts (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- The receivable dropdown lists no header accounts and no cash, bank or inventory accounts.
- POST/PUT /accounting/customers with a header account id returns 400 with a field error; the same for suppliers and payable accounts.
- A customer saved with the default account can be invoiced without error.
