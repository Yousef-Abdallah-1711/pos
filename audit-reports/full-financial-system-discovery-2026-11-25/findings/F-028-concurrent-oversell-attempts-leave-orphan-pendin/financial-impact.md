# F-028 — Financial impact

- **Quantified impact (register):** 15 x 11000 = 165000 gross phantom pending sales
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** concurrency

## What is wrong in the books or reports
5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in Sales Analytics

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

