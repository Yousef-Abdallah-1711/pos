# Complete feature inventory

**Status: PARTIAL.** Routes were opened in a synthetic tenant, mostly empty. API 200 does not establish financial correctness.

| Module | Frontend route | API/source | Runtime evidence | Status |
|---|---|---|---|---|
| Chart of Accounts | /en/dashboard/accounting/chart-of-accounts | accounting/accounts | Opened; two isolated accounts created | Partial |
| Journal Entries | /en/dashboard/accounting/journal-entries and /{id} | accounting/journal-entries GET/POST/PUT, post, void | UI submit failed when attachments omitted; API with [] created and UI posted | Verified with defect |
| General Ledger | /en/dashboard/accounting/general-ledger | accounting/general-ledger | Opened empty; DB verified rows after post | Partial |
| Trial Balance | no Dashboard page found | /api/v1/accounting/trial-balance | API totals 100/100, balanced | API only |
| Bank Accounts | /en/dashboard/accounting/bank-accounts | accounting/bank-accounts | Opened, empty | Partial |
| Advanced Accounting | /en/dashboard/accounting/advanced | accounting export/audit handlers | Opened | Partial |
| Customers | /en/dashboard/accounting/customers | customer/accounting APIs | Opened, empty | Partial |
| Customer invoices/payments | /en/dashboard/accounting/customer-invoices and ?tab=payments | invoice/payment APIs | Both opened, empty | Partial |
| Financial report index | /en/dashboard/reports/financial | financial report APIs | Opened | Partial |
| Advanced financial | /en/dashboard/reports/financial/advanced | predictive cash-flow/audit exports | Opened; initially prompted refresh | Partial |
| Inventory valuation | /en/dashboard/reports/financial/inventory-valuation | report endpoint | Opened/API 200; no fixture | Partial |
| Payables aging | /en/dashboard/reports/financial/payables-aging | report endpoint | Opened/API 200; no fixture | Partial |
| Receivables aging | /en/dashboard/reports/financial/receivables-aging | report endpoint | Opened/API 200; no fixture | Partial |
| Cash flow | /en/dashboard/reports/financial/cash-flow | report endpoint | API 200; controlled scenario net 100 | Partial |
| Balance Sheet | /en/dashboard/reports/financial/balance-sheet | report endpoint | equation_balanced=false after posting | Failed reconciliation |
| Profit & Loss | /en/dashboard/reports/financial/profit-loss | report endpoint | UI/API revenue and net profit 100 | Partial; reconciliation issue |
| Dashboard | /en/dashboard | dashboard aggregates | Opened and screenshot captured | Partial |
| VAT report (additional) | /en/dashboard/reports/financial/vat | accounting/tax/vat-report | Opened/API 200; empty | Partial |
| Suppliers and supplier invoices | supplier routes | supplier/invoice/payment APIs | Opened/API 200; empty | Partial |

See reports 02, 04, 05, and 06 for source paths. Permission-negative tests and all existing tests were not exhaustively reviewed or run.
