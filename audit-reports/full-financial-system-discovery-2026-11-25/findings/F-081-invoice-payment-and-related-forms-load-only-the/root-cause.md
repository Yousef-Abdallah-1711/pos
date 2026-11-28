# F-081 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): the customer list endpoint defaults to page_size 20 (backend/internal/handler/http/customer_acct_handler.go:29; suppliers: supplier_handler.go:29) and the invoice pages call getCustomerAccounts() / getSuppliers() without paging or search (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:120; supplier-invoices/page.tsx:107)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/customer_acct_handler.go`:29.
- Propagation: affected modules Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns); related findings F-080; observed consequence: The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20
- Cause versus symptom: the symptom is what users or reports show (The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-INVOICE-UI; scripts: tests/runtime/ph09_r5_invoice_form.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/customer_acct_handler.go | 29 | backend/internal/handler/http/customer_acct_handler.go | exact | yes |
| backend/internal/handler/http/supplier_handler.go | 29 | backend/internal/handler/http/supplier_handler.go | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx | 120 | frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx | 107 | frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx | exact | yes |
| supplier-invoices/page.tsx | 107 | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | RANGE BEYOND FILE |

