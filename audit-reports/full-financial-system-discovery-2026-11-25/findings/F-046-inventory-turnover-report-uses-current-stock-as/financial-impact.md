# F-046 — Financial impact

- **Quantified impact (register):** Turnover ratio misstated (1.42 vs 1.01 using true average 41.5 qty)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** report

## What is wrong in the books or reports
opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has waste 1 for products 2 and 4; Audit Burger original stock 50 | RUN-003: turnover opening_stock_qty == current stock (33/21); waste_qty 0 vs movements waste 1 (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Inventory turnover report; Overview report

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

