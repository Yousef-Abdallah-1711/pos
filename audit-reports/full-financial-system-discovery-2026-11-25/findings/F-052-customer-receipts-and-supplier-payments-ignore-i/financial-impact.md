# F-052 — Financial impact

- **Quantified impact (register):** Every retried or double-clicked payment duplicates cash and AR/AP movement
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, concurrency

## What is wrong in the books or reports
Both return 201; payments_received has 2 rows (sum 800); invoice amount_paid 800 / balance 200; AR credited twice | RUN-003 REL-025: identical Idempotency-Key replay creates 2 receipts (sum 800) (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Receipts; Supplier payments; Accounting (cash, AR/AP)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

