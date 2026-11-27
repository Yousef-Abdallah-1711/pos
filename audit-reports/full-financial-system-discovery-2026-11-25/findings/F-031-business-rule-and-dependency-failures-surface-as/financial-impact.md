# F-031 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, ui

## What is wrong in the books or reports
Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Orders; Receipts; Platform

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

