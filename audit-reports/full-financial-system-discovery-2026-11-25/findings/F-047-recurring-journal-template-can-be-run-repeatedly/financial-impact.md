# F-047 — Financial impact

- **Quantified impact (register):** 700 duplicated per extra run in test
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, concurrency

## What is wrong in the books or reports
Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 2026-12-26 | RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_advanced_repo.go:396-420). 

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting - recurring journals

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

