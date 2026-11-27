# F-017 — Financial impact

- **Quantified impact (register):** duplicate revenue risk
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** concurrency

## What is wrong in the books or reports
5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

