# Performance acceptance results

**Status: BASELINE MEASURED; optimization verification pending.** These descriptive synthetic results do not establish production capacity or financial release acceptance.

Hosted disposable PostgreSQL benchmark run [44962457708](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44962457708) completed on exact baseline SHA `14fa8640191ae08c2b0789b4d05e370dd3023834`. The machine-readable report, including all five `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` plans and configuration, is `audit-reports/financial-benchmark/cfr010-reconciliation-baseline-14fa8640191ae08c2b0789b4d05e370dd3023834.json`.

Environment: PostgreSQL 16.4, `shared_buffers=128MB`, `work_mem=4MB`, `effective_cache_size=4GB`, four hosted CPUs, `GOMAXPROCS=2`. The principal scope had 10,000 journals, 50,000 journal lines, and 50,000 GL rows; the secondary scope had 25/75/75, plus a real empty restaurant scope. The final corpus also included benchmark and concurrent writes. Each warm result used five warmups and 30 samples; concurrent runs used 40 samples at concurrency four. The first post-seed run is not represented as a guaranteed cold-cache run.

| Reconciliation scope/workload | p50 | p95 | max |
| --- | ---: | ---: | ---: |
| Large restaurant, warm | 174.044 ms | 175.441 ms | 175.524 ms |
| Small restaurant, warm | 8.880 ms | 9.456 ms | 9.473 ms |
| Empty restaurant, warm | 8.900 ms | 9.572 ms | 9.584 ms |
| Large scope, four concurrent readers | 288.892 ms | 298.698 ms | 306.897 ms |
| Large scope reads during concurrent posting | 285.280 ms | 291.821 ms | 298.111 ms |

The large-scope baseline line-to-GL query plan returned 50,022 lines but formed 250,118 intermediate rows, removed 200,022 through its join filter, and spilled sort work (`work_mem=4MB`, 3,544 temp blocks read / 713 written). This is a measured bottleneck being tested with the existing composite `(journal_entry_id, journal_line_id)` index. Other plan components included full scans of the disposable GL/journal tables. The five plans recorded zero shared-read blocks in this warmed environment. A configured 1 ms statement timeout cancelled the operation with SQLSTATE `57014`; this validates cancellation behavior only, not a recommended timeout.

No approved performance budget was found. The baseline identifies a query-plan bottleneck but does not itself establish a production SLO or acceptance. The targeted query rewrite must pass the same correctness and workload checks before its before/after evidence can be assessed.

## Release effect

Performance acceptance is a remaining engineering gate for full CFR-010 closure. The repository method has no public endpoint in this candidate, limiting exposure, but this does not replace optimized-candidate verification and an approved budget before operational use.
