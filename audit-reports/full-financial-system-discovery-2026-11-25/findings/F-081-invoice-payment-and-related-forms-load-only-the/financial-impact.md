# F-081 — Financial impact

- **Quantified impact (register):** 0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui

## What is wrong in the books or reports
The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

