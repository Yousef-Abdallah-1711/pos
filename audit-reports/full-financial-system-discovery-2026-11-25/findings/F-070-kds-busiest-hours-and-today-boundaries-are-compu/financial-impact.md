# F-070 — Financial impact

- **Quantified impact (register):** 0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** report

## What is wrong in the books or reports
API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found

## Which reports and decisions are affected
Financial reports that read the affected data: KDS analytics and stats

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

