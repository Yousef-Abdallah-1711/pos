# 13 — Financial scale and performance (planning target)

Ambition: 1,000 × 10,000 = 10,000,000 business events/day, ~116/s average across 86,400 seconds; an illustrative 10× burst is ~1,160/s. At 4–10 writes/event, illustrative peak is ~4,600–11,600 writes/s. These are assumptions for workload design, not achieved capacity.

Measure event sizes, lines, indexes, retention, report concurrency and retry amplification before sizing. Preserve synchronous source/outbox and journal atomicity. Projections/caches may be eventually consistent only with visible watermark/staleness. Evaluate partitioning from plans, avoid hot running-balance locks, use bounded pools, queue backpressure, tenant fairness, streaming exports, eligible replicas, backups and restore exercises. Propose SLOs only with owner approval: p95/p99, queue age, DB headroom, no loss/duplication on restart, storage growth. Existing dev browser timing is single-sample; production baseline NOT TESTED.
