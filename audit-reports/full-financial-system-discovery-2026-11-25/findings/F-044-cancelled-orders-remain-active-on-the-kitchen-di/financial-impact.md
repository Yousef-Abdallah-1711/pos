# F-044 — Financial impact

- **Quantified impact (register):** Operational: food prepared for cancelled orders (COGS without sale)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** general

## What is wrong in the books or reports
Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

