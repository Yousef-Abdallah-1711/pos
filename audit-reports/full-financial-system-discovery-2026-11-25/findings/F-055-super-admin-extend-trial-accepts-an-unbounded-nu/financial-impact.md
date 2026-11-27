# F-055 — Financial impact

- **Quantified impact (register):** Unbounded free access; no billing record
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** security, validation

## What is wrong in the books or reports
200 days_added 100000, new_trial_end_date 2300-10-02 (negative days correctly rejected) | RUN-003: extend-trial 100000 days accepted -> 2300-10-02 (REPRODUCED; tenant 8 restored by fixture). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

