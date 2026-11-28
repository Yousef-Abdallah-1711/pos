# F-079 — Financial impact

- **Quantified impact (register):** 0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400" for every valid submission; the request body carries start_date "2026-11-20T00:00:00.000Z"

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

