# F-020 — Financial impact

- **Quantified impact (register):** n/a
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, security, ui

## What is wrong in the books or reports
GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact data; accounting/report/admin routes 403/401 (REPRODUCED).  RUN-005: in enforce mode the EMPLOYEE (cashier) role gets 200 from GET /hr/employees with base_salary (5000) and email in the payload, and 200 from /hr/attendance and /hr/leaves; MANAGER additionally reads /hr/salaries (4 rows, salary fields) and the user and role lists (policy question POL-15, not asserted as a defect); DRIVER sees 0 orders (scoped) and 3/28 readable endpoints.

## Which reports and decisions are affected
Financial reports that read the affected data: RBAC; HR; POS

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

