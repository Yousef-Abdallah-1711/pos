# F-080 — Financial impact

- **Quantified impact (register):** 0 direct; blocks invoicing and payments for affected customers and suppliers
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui, validation

## What is wrong in the books or reports
Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer returned 400 {"error":"account 63 is a header account and cannot be posted to"}, so the customer can never be invoiced until the account is edited

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting > Customers, Suppliers, Customer invoices, Payments

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

