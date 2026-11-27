# F-066 — Financial impact

- **Quantified impact (register):** 0 (wrong operational count; misleads the owner about catalogue size)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

