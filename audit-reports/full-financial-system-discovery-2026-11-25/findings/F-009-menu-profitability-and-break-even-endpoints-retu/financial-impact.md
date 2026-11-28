# F-009 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** report

## What is wrong in the books or reports
HTTP 500 for both (cost_price column absent) | RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONFIRMED: repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost) -> SQL error mapped to generic 500 (accounting_advanced_handler.go:782). 

## Which reports and decisions are affected
Financial reports that read the affected data: Financial reports advanced

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

