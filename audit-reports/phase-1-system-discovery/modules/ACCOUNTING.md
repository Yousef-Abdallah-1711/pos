# Accounting (dashboard) - static audit

Status: **Verified implemented** for COA, journal entries, ledger, customers/suppliers, invoices, bank accounts; **Partial** for advanced (reconciliation, recurring, budgets, FX, e-invoice, reminders) - UI is one aggregate page; Excel import/export, invoice PDF and tax-period close have no UI.

## Pages (10)
`accounting/chart-of-accounts`, `journal-entries` (+`[id]`), `general-ledger`, `bank-accounts`, `customers`, `customer-invoices` (539 lines, tabs incl. payments), `suppliers`, `supplier-invoices` (506), `advanced` (298; **no gate**). Gate ACCOUNTING/SUPPLIERS on the others. Components `components/accounting/*` (InvoiceForm 600+ lines, JournalEntryForm). API modules: `accountingApi.ts`, `accountingAdvancedApi.ts`, `bankApi.ts`, `customersApi.ts`, `suppliersApi.ts`.

## Contract
All 70+ accounting paths used by the dashboard match `main.go:982-1057, 1353-1359` (RBAC module 15/16 on every route). Backend routes without UI: `GET /accounting/customer-invoices/{id}/pdf`, `PUT|DELETE /accounting/recurring-journals/{id}`, all four `/accounting/excel/*`, `POST /accounting/tax/periods/{id}/close`.

## Client-side financial logic (risk)
- `InvoiceForm.tsx:601-650`: per-line subtotal, discount, tax, line total with unrounded floats; payload sends client `tax_amount`, `discount_amount`, `line_total` (server must recompute) (FE-DASH-018).
- `JournalEntryForm.tsx:53-62`: debit/credit totals and balanced check with epsilon 0.001; server posting rules not read.
- Outstanding balances summed client-side from loaded page of invoices (`customer-invoices/page.tsx:100`, `supplier-invoices/page.tsx:94`) - partial if paginated.
- Currency: `Intl.NumberFormat(... 'EGP')` in advanced page (`advanced/page.tsx:72`).
- Excel import at `advanced/page.tsx:169` posts multipart via `api.post(url, form)` with explicit multipart header (unused backend excel routes suggest different URLs; call goes through `accountingAdvancedApi` path variable - not verified).

## i18n
Custom translator with `bankAccountMessages.*` keys missing in en.json (8 keys, FE-DASH-016).

## Findings
FE-DASH-018, FE-DASH-014 (advanced ungated), FE-DASH-016.

## Unverified
Double-entry enforcement, period locking, auto-posting from orders/payroll (backend `autoPostingEngine`).
