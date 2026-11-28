# F-045 — Financial impact

- **Quantified impact (register):** Forecast understates inflows by 319,660 in this dataset
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report, ui

## What is wrong in the books or reports
2026-11 history cash_in 13,000 / cash_out 32,400 -> forecast Dec -18,430 while 319,660 of paid POS orders and GL cash +276,270 exist | RUN-003: PCF Nov cash_in 13,000/cash_out 32,400 vs paid POS orders 319,660; forecast negative (CONFIRMED accounting_advanced_repo.go:1105-1175). 

## Which reports and decisions are affected
Financial reports that read the affected data: Financial hub predictive cash flow

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

