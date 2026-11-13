# Scale and reliability assessment

**Status: UNVERIFIED.** No load test, high-cardinality dataset, recovery drill or production capacity test was performed.

- Migration 157 defines composite reporting indexes, but no EXPLAIN ANALYZE or selectivity measurement was captured.
- Some CSV export handlers build row arrays before writing; memory and row limits were not measured.
- Development dashboard container used about 2.46 GiB at one point after navigation; no CPU/memory limits were set. This is not production capacity evidence.
- Auto-posting creates a draft before its posting transaction; failed posting may leave an orphan draft for review. Ledger insertion and posted status are transactional in the inspected path.
- No load, concurrency, queue backlog, pool saturation, Redis loss, provider outage, backup restore or failover checks were run.

Before capacity claims: representative scale fixtures, tenant/aggregate load profile, bounded streaming exports, query plans and pool telemetry, idempotent recovery/reconciliation, backup restore, resource limits and SLOs.


## Reliability evidence status — 2026-11-11

No duplicate-request, webhook replay, concurrent payment/refund, retry, or crash-recovery test was run in this continuation. Repository source shows ledger insertion and posted status commit together in the inspected posting transaction, while automatic draft creation is a preceding separate operation. This source observation is not failure-injection proof. Reliability status remains PARTIAL; see the workflow matrix.
