# F-013 — Financial impact

- **Quantified impact (register):** 100000 reversed in closed period
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, capability

## What is wrong in the books or reports
Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). 

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

