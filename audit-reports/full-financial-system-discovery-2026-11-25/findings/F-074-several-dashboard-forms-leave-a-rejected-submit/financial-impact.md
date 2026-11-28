# F-074 — Financial impact

- **Quantified impact (register):** 0
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** ui

## What is wrong in the books or reports
POST returns 400 (for example pq: value too long for type character varying(20)); the page shows "Runtime AxiosError" and "Request failed with status code 400" (the unhandled-rejection overlay in this build) with no field message; empty submits do show client-side required messages, so only server-side rejections are unhandled

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

