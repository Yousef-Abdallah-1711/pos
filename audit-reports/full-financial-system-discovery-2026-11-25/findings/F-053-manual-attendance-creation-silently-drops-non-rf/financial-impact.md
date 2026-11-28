# F-053 — Financial impact

- **Quantified impact (register):** Overtime cost understated whenever times are entered in the simple form
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report, validation

## What is wrong in the books or reports
"09:00" form -> 201 with clock_in/clock_out NULL and regular_hours 0 / overtime_hours 0; reversed order on the same form also 201; RFC3339 form computes 8.00 + 2.50; reversed RFC3339 -> 500 with raw constraint text valid_clock_times; overtime report for November shows 0 hours | RUN-003: attendance with 09:00/19:30 stored NULL times and 0 hours (201); RFC3339 computes 8+2.5 (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: HR attendance; HR overtime/absenteeism reports; payroll overtime input

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

