# F-043 — Financial impact

- **Quantified impact (register):** 14,360 profit overstatement; 226,270.01 cash understatement in this dataset
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report

## What is wrong in the books or reports
KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_accounts_ledger) vs GL cash+bank 276,270.01 | RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_position 50,000 vs GL cash+bank 276,270.01 (CONFIRMED accounting_advanced_repo.go:838-880). 

## Which reports and decisions are affected
Financial reports that read the affected data: Financial KPI dashboard

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

