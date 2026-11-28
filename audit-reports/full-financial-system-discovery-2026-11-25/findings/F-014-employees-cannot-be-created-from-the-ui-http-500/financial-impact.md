# F-014 — Financial impact

- **Quantified impact (register):** payroll blocked until API fixture
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
500 positive_working_hours violation; payload keys weekly_hours/monthly_salary vs backend working_hours_per_week/base_salary | RUN-003 (headed Chrome 150, tenant 4 owner): POST /hr/employees -> 500 reproduced; DB before=after (no orphan rows); UI shows only generic 'Request failed with status code 500' - raw DB error is in the API response body (see evidence). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

