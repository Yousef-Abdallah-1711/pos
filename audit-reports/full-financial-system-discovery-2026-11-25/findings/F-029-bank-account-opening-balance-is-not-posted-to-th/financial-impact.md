# F-029 — Financial impact

- **Quantified impact (register):** 60000 gap
- **Accounts / statements named in the evidence:** 1111, 1112
- **Impact classes:** accounting, capability, report

## What is wrong in the books or reports
No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal (GL 1112 stays 0) (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Bank accounts; Reports KPI; Accounting

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

