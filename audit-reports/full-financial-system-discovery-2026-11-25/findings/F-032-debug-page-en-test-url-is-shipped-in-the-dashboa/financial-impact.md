# F-032 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** security, ui

## What is wrong in the books or reports
Page renders TEST PAGE LOADED with no auth protection and a hydration error RUN-005: GET http://localhost:13112/en/test-url returned HTTP 200 (curl) on the current isolated dashboard.

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

