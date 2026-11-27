# F-068 — Financial impact

- **Quantified impact (register):** Proportional: each extra active policy repeats the employee base salary and deductions (here 5,000.00 base salary listed twice)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, report

## What is wrong in the books or reports
1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total over the file double-counts

## Which reports and decisions are affected
Financial reports that read the affected data: HR reports > Social insurance export

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

