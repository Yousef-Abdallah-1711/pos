# F-081 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Invoice, payment and related forms load only the first 20 customers or suppliers, so customers beyond the first page cannot be selected.

## 2. Current behavior (verified)
The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20

## 3. Required behavior
Every active customer can be selected (search or paging)

## 4. Reproduction conditions
Headed Chrome tenant 4 (91 active customers): open New customer invoice and count the options in the customer dropdown

## 5. Verified evidence
- evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-RESULT.json
- tests/runtime/ph09_r5_invoice_form.mjs

## 6. Root cause and confidence
CONFIRMED — the customer list endpoint defaults to page_size 20 (backend/internal/handler/http/customer_acct_handler.go:29; suppliers: supplier_handler.go:29) and the invoice pages call getCustomerAccounts() / getSuppliers() without paging or search (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:120; supplier-invoices/page.tsx:107)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/customer_acct_handler.go`:29, `backend/internal/handler/http/supplier_handler.go`:29, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx`:120, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx`:107, `frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx`:107
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: unchanged (search and page_size parameters already supported or additive)
- Accounting logic / posting rules: none named
- Reports affected: Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns)

## 8. Impact
- Financial: 0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI
- Operational: Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: replace the plain customer and supplier selects in the invoice, receipt and payment forms with an async searchable combobox that queries the list endpoints with a search term and page size (backend/internal/handler/http/customer_acct_handler.go:29 and supplier_handler.go:29 keep the default of 20 for list screens); alternatively load all pages before rendering; show the total count and a hint when results are truncated.

## 10. Alternative designs
A: async searchable combobox (recommended, scales). B: raise the default page size (hides the problem until the next threshold).

## 11. Required policy decisions
None.

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
F-080 and shared-file findings F-021, F-062, F-080

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/customer_acct_handler.go, backend/internal/handler/http/supplier_handler.go, frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
With more than 20 customers every customer can be found and selected in the invoice and receipt forms (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With 91 active customers every customer can be found by typing part of the name and selected.
- The form never silently shows a truncated list.
- The same holds for supplier selection in supplier invoices and payments.
