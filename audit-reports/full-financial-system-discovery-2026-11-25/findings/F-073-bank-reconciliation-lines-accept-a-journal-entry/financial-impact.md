# F-073 — Financial impact

- **Quantified impact (register):** 0 direct; cross-tenant reference integrity and error-text leakage (F-025 class)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, security

## What is wrong in the books or reports
HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the foreign key, with the raw constraint text in the error body

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting > Bank reconciliation

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

