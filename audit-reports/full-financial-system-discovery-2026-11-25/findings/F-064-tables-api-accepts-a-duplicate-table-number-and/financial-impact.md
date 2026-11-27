# F-064 — Financial impact

- **Quantified impact (register):** 0 (dine-in orders can attach to an ambiguous table)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** concurrency, validation

## What is wrong in the books or reports
duplicate -> 201 (two tables numbered 1, ids 1 and 2); long name -> 500 "Failed to create table"; capacity (1-50), negative number and delete-with-active-orders (409) are validated correctly | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

