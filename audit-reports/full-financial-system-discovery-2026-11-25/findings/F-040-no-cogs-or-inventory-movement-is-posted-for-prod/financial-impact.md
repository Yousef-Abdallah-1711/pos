# F-040 — Financial impact

- **Quantified impact (register):** 120 in test (cost x quantity of every untracked sale)
- **Accounts / statements named in the evidence:** 1110, 1130, 2120, 4100, 5110
- **Impact classes:** accounting, report

## What is wrong in the books or reports
Order 39 (untracked): JE-2026-00064 posts only Dr 1110 330 / Cr 4100 300 / Cr 2120 30 - no COGS; tracked order 41: JE-2026-00066 Dr 5110 120 / Cr 1130 120. Profitability report shows cogs 120 for both products => report COGS exceeds GL COGS by 120 | RUN-003 controlled experiment: untracked product order 39 posts no COGS while Profitability report shows cogs 120; tracked twin order 41 posts COGS 120 (CONFIRMED, auto_posting_engine.go:173). 

## Which reports and decisions are affected
Financial reports that read the affected data: Products; POS/Orders; Accounting; Profitability report

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

