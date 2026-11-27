# F-080 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): CustomerForm.tsx:61-62 filters the dropdown only by account_type === "asset" (the supplier form uses the same pattern), the customer handler passes ReceivableAccountID straight to the repository (backend/internal/handler/http/customer_acct_handler.go:141,244) without checking that it is a postable receivable account, and the first posting attempt fails at backend/internal/repository/journal_entry_repo.go:708


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx`:61.
- Propagation: affected modules Accounting > Customers, Suppliers, Customer invoices, Payments; related findings F-062, F-075; observed consequence: Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer re
- Cause versus symptom: the symptom is what users or reports show (Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-INVOICE-UI; scripts: tests/runtime/ph09_r5_invoice_form.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 61 | frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | exact | yes |
| backend/internal/handler/http/customer_acct_handler.go | 141 | backend/internal/handler/http/customer_acct_handler.go | exact | yes |
| backend/internal/repository/journal_entry_repo.go | 708 | backend/internal/repository/journal_entry_repo.go | exact | yes |
| CustomerForm.tsx | 61-62 | frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | basename | yes |

