# F-035 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contains unsafe-inline and unsafe-eval in the isolated build; the API /health response carries no CSP header and the landing response carried none in this run.

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

