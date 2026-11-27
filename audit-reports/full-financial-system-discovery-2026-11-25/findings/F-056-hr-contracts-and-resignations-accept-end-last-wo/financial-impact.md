# F-056 — Financial impact

- **Quantified impact (register):** 0 directly; wrong notice/tenure calculations downstream
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** validation

## What is wrong in the books or reports
both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected only by DB check constraints that surface as 500 | RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

