# F-024 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting

## What is wrong in the books or reports
tax_amount 630.0000000000001 (server stores rounded) RUN-005 (real invoice form, 1 x 5,000, 10% discount, 14% VAT): the POST payload carries "tax_amount":630.0000000000001 (raw float artefact) while the UI shows Tax EGP 630.00.

## Which reports and decisions are affected
Financial reports that read the affected data: Invoices

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

