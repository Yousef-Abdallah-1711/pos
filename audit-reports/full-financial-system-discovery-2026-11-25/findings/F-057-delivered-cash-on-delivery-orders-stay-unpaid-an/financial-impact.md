# F-057 — Financial impact

- **Quantified impact (register):** All unsettled COD sales are missing from GL revenue/VAT/cash
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, capability

## What is wrong in the books or reports
status delivered with payment_status pending and 0 journals; revenue stays unrecognised until someone manually marks paid; no driver cash reconciliation exists | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Delivery; Drivers; Orders; Accounting

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

