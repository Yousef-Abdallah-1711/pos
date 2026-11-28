# F-042 — Financial impact

- **Quantified impact (register):** 529,100 shown vs 319,660 paid tax-inclusive vs 305,300 GL revenue
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report

## What is wrong in the books or reports
total_revenue 529,100 = SUM of ALL 41 orders incl 3 cancelled (11,440) and 21 unpaid; delivered_today 30 while delivered_count 0 (counts any order created today) | RUN-003: order stats total_revenue 529,100 = all 41 orders incl 3 cancelled; delivered_today 30 vs delivered_count 0 (CONFIRMED order_repo.go:859-870). 

## Which reports and decisions are affected
Financial reports that read the affected data: Orders dashboard; Order stats API

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

