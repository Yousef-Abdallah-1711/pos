# F-076 — Financial impact

- **Quantified impact (register):** 0
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

