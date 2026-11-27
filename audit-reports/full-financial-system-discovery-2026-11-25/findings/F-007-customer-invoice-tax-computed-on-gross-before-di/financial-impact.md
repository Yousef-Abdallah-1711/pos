# F-007 — Financial impact

- **Quantified impact (register):** 70 per invoice
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui

## What is wrong in the books or reports
Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE). 

## Which reports and decisions are affected
Financial reports that read the affected data: Customer invoices

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

