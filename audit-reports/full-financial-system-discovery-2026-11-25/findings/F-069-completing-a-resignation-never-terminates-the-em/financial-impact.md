# F-069 — Financial impact

- **Quantified impact (register):** Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report

## What is wrong in the books or reports
Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{id} sets termination_date (to CURRENT_DATE, not the last working date); then turnover = 1/1 = 100% RUN-005 real-UI reproduction: Complete on resignation 9 returned 200 with settlement_amount 0; the employee stayed active and no termination date was set.

## Which reports and decisions are affected
Financial reports that read the affected data: HR resignations, HR turnover report, employees, payroll eligibility

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

