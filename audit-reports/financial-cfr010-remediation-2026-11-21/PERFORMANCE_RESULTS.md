# Reconciliation performance results

**Status: measured before/after; performance acceptance remains OPEN.** These disposable synthetic results do not establish production capacity. The repository has no approved reconciliation latency budget, so no performance target is inferred here.

## Runs and environment

- Baseline: [hosted run 44962457708](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44962457708), exact SHA `14fa8640191ae08c2b0789b4d05e370dd3023834`.
- Optimized candidate: [hosted run 44963630632](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44963630632), exact SHA `b06296730ced62e68a595473f66ff425d15c1223`.
- Machine-readable reports: baseline `audit-reports/financial-benchmark/cfr010-reconciliation-baseline-14fa8640191ae08c2b0789b4d05e370dd3023834.json`; optimized `audit-reports/financial-benchmark/cfr010-reconciliation-optimized-b06296730ced62e68a595473f66ff425d15c1223.json`.
- Both runs used disposable PostgreSQL 16.4 (`shared_buffers=128MB`, `work_mem=4MB`, `effective_cache_size=4GB`), four hosted CPUs, `GOMAXPROCS=2`, five warmups, and 30 warm samples. The large restaurant scope had 10,000 journals, 50,000 lines, and 50,000 GL rows; the small scope had 25/75/75. An empty restaurant scope was also measured. The corpus balanced journal and GL debits/credits at 110,032,700 cents.
- The first post-seed operation is reported as `FIRST_RUN_ONLY`, not as a guaranteed cold-cache measurement; neither database run flushed PostgreSQL or host caches.

## Reconciliation latency

| Scope/workload | Baseline p50 / p95 / max | Optimized p50 / p95 / max |
| --- | ---: | ---: |
| Large restaurant, warm (30 samples) | 174.044 / 175.441 / 175.524 ms | 179.243 / 180.141 / 180.271 ms |
| Small restaurant, warm (30 samples) | 8.880 / 9.456 / 9.473 ms | 9.110 / 10.026 / 10.448 ms |
| Empty restaurant, warm (30 samples) | 8.900 / 9.572 / 9.584 ms | 8.831 / 9.576 / 9.578 ms |
| Large restaurant, four concurrent readers (40 samples) | 288.892 / 298.698 / 306.897 ms | 224.396 / 239.402 / 240.067 ms |
| Four readers during concurrent posting (40 samples) | 285.280 / 291.821 / 298.111 ms | 220.453 / 229.043 / 240.090 ms |

All recorded reconciliation workloads returned zero errors. Eleven legitimate concurrent posts overlapped with each run's posting workload. A 1 ms statement timeout cancelled the operation with SQLSTATE `57014` in both runs; this verifies cancellation behavior only and is not a recommended timeout.

## Bottleneck and observed effect

The baseline large-scope line-to-GL plan took 160.303 ms, built 250,118 intermediate rows for 50,022 lines, discarded 200,022 rows through a join filter, and spilled sort work. The lateral source-link rewrite in `b062967` reduced that plan's execution to 115.720 ms and removed the join-filter discards. The optimized plan still uses sequential scans and spills sort work (`work_mem=4MB`); the JSON artifacts retain all five `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` plans and index statistics.

This plan-level improvement did not lower single-reader whole-operation latency: large-scope p50 increased by about 3.0% (174.044 to 179.243 ms). Four-reader p50 improved by about 22.3% (288.892 to 224.396 ms), and reads during posting improved by about 22.7% (285.280 to 220.453 ms). The before/after comparison is descriptive: it shows a selectivity improvement and better concurrent results alongside a small single-reader regression. It does not prove a production benefit or meet a defined SLO.

## Release effect

The candidate passes its hosted synthetic PostgreSQL benchmark and correctness checks on the exact candidate SHA, but reconciliation performance acceptance remains OPEN pending an approved workload budget and decision about the observed tradeoff. CFR-010 remains OPEN and financial production remains NO-GO. Original operating root cause remains **NOT ESTABLISHED**; OPEN-01, EARN-01, and CLOSE-01 remain pending qualified accounting approval.
