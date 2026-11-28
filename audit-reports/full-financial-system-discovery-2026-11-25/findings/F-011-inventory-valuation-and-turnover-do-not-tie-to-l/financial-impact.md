# F-011 — Financial impact

- **Quantified impact (register):** 180000 / 148000
- **Accounts / statements named in the evidence:** 1130, 5110
- **Impact classes:** accounting, capability, report

## What is wrong in the books or reports
196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no journal (STILL REPRODUCIBLE). 

## Which reports and decisions are affected
Financial reports that read the affected data: Inventory; Products; Reports

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

