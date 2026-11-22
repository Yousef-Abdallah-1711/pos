# Reconciliation performance results

**Status: CANDIDATE PLAN VERIFIED; performance budget and production acceptance remain open.** These synthetic results do not establish production capacity or financial release acceptance.

## Baseline and candidate

- Baseline SHA: `14fa8640191ae08c2b0789b4d05e370dd3023834`; benchmark [44962457708](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44962457708). Its machine-readable report, including five `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` plans, is `audit-reports/financial-benchmark/cfr010-reconciliation-baseline-14fa8640191ae08c2b0789b4d05e370dd3023834.json`.
- Candidate code SHA: `0a1aca0374ce8b357e1709f43ac6e1e242b7f575`; benchmark [44979879216](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979879216), machine-readable artifact [financial-query-benchmark-0a1aca0...-10000](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979879216/artifacts/14736679271), and exact-SHA financial PostgreSQL acceptance [44979881780](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979881780).
- PR: [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110). Its exact-head CI, including the benchmark, financial PostgreSQL acceptance, and application image smoke, passed. The PR remains open for maintainer review.

Both benchmark runs used PostgreSQL 16.4 on an Intel Xeon 6975P-C with the same reported 10,102 journals, 50,229 journal lines, and 50,229 GL rows. The large restaurant scope had 10,000 journals and 50,000 lines/GL rows; the small scope had 25/75/75; the empty scope was a seeded restaurant with zero rows. Warm samples used five warmups and 30 measurements; concurrent workloads used 40 measurements at concurrency four.

## Measurements

| Reconciliation scope/workload | Baseline p50 / p95 / max | Candidate p50 / p95 / max |
| --- | ---: | ---: |
| Large restaurant, warm | 174.044 / 175.441 / 175.524 ms | 180.726 / 182.565 / 182.660 ms |
| Small restaurant, warm | 8.880 / 9.456 / 9.473 ms | 9.404 / 9.853 / 10.073 ms |
| Empty restaurant, warm | 8.900 / 9.572 / 9.584 ms | 9.608 / 10.249 / 10.329 ms |
| Large scope, four concurrent readers | 288.892 / 298.698 / 306.897 ms | 281.491 / 335.521 / 354.225 ms |
| Large scope reads during posting | 285.280 / 291.821 / 298.111 ms | 278.025 / 293.865 / 317.084 ms |

The warm large-scope operation p50 was 6.682 ms (3.8%) slower in the candidate run. Therefore these measurements do not demonstrate an end-to-end latency improvement. The line-to-GL query plan itself improved: its execution time was 160.660 ms at baseline and 148.041 ms for the candidate. The baseline plan formed 250,118 intermediate rows, then used an incremental sort and broad join that read 3,544 temp blocks and wrote 713. The candidate plan used `idx_general_ledger_journal_line` for 50,022 index probes and had zero temp reads/writes. Treat the plan and full-operation measurements together; no SLO or production acceptance is claimed.

The first operation after seeding in the candidate run was 200.609 ms. PostgreSQL was not restarted and caches were not flushed, so this is not a guaranteed cold-cache result. The 1 ms statement timeout cancelled as configured with SQLSTATE `57014`; this validates cancellation behavior, not a recommended timeout. Eleven concurrent postings completed during the reads workload.

No approved reconciliation performance budget was found. The measured plan no longer spills sort work, but the end-to-end p50 did not improve in this comparison. Obtain an approved latency budget and review repeated, controlled performance results before declaring performance acceptance or production capacity.

## Release effect

The reconciliation implementation is proposed in PR #110 and has passed exact-candidate CI and disposable PostgreSQL acceptance. It has not been merged. An independent read-only review confirmed SQL semantic equivalence and the benchmark correction; no independent GitHub reviewer is available in the repository besides the PR author. CFR-010 remains OPEN and financial production remains NO-GO pending source reconciliation, qualified accounting approvals, and the approved performance budget.
