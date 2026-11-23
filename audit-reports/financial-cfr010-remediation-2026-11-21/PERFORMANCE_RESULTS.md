# Reconciliation performance results

**Status: CANDIDATE PLAN VERIFIED; performance budget and production acceptance remain open.** These synthetic results do not establish production capacity or financial release acceptance.

## Baseline and candidate

- Baseline SHA: `14fa8640191ae08c2b0789b4d05e370dd3023834`; benchmark [44962457708](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44962457708). Its machine-readable report, including five `EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)` plans, is `audit-reports/financial-benchmark/cfr010-reconciliation-baseline-14fa8640191ae08c2b0789b4d05e370dd3023834.json`.
- Candidate code SHA: `0a1aca0374ce8b357e1709f43ac6e1e242b7f575`; benchmark [44979879216](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979879216), machine-readable artifact [financial-query-benchmark-0a1aca0...-10000](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979879216/artifacts/14736679271), and exact-SHA financial PostgreSQL acceptance [44979881780](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979881780).
- Latest exact-head candidate: source head `05aa99dd5b5e7294a45350d008db380ba1b49729`; benchmark [45217215534](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45217215534), effective code SHA `9c28aa067d207f0a2806f17db46cf39fc86fe94b`, and machine-readable artifact [financial-query-benchmark-9c28aa0...-10000](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45217215534/artifacts/14838919653). Exact-head PostgreSQL acceptance also passed.
- PR: [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110). The latest exact-head benchmark and PostgreSQL acceptance passed. This report evaluates the proposed query change; merging it does not approve operational performance or production use.

The baseline and both candidate runs used PostgreSQL 16.4 on an Intel Xeon 6975P-C with the same reported 10,102 journals, 50,229 journal lines, and 50,229 GL rows. The large restaurant scope had 10,000 journals and 50,000 lines/GL rows; the small scope had 25/75/75; the empty scope was a seeded restaurant with zero rows. Warm samples used five warmups and 30 measurements; concurrent workloads used 40 measurements at concurrency four.

## Measurements

| Reconciliation scope/workload | Baseline p50 / p95 / max | Candidate p50 / p95 / max |
| --- | ---: | ---: |
| Large restaurant, warm | 174.044 / 175.441 / 175.524 ms | 178.529 / 180.294 / 181.966 ms |
| Small restaurant, warm | 8.880 / 9.456 / 9.473 ms | 9.209 / 9.875 / 9.900 ms |
| Empty restaurant, warm | 8.900 / 9.572 / 9.584 ms | 9.392 / 9.989 / 10.157 ms |
| Large scope, four concurrent readers | 288.892 / 298.698 / 306.897 ms | 295.188 / 346.289 / 350.931 ms |
| Large scope reads during posting | 285.280 / 291.821 / 298.111 ms | 294.535 / 307.895 / 330.524 ms |

For the latest exact-head run, the large-scope warm p50/p95/max were 4.485/4.853/6.442 ms (2.6%/2.8%/3.7%) slower than the baseline. An earlier candidate run measured 180.726/182.565/182.660 ms, so repeated candidate runs vary; the latest candidate is faster than that earlier candidate run but still does not demonstrate an end-to-end improvement over baseline. The line-to-GL query plan improved from 160.660 ms at baseline to 146.178 ms for the latest candidate (9.0% faster). The baseline plan formed 250,118 intermediate rows, used an incremental sort and broad join, and read 3,544 temporary blocks while writing 713. The latest candidate used `idx_general_ledger_journal_line` for 50,022 index probes and had zero temporary reads/writes. Treat the plan and full-operation measurements together; no SLO or production acceptance is claimed.

The latest first operation after seeding was 200.512 ms. PostgreSQL was not restarted and caches were not flushed, so this is not a guaranteed cold-cache result. The 1 ms statement timeout cancelled as configured with SQLSTATE `57014`; this validates cancellation behavior, not a recommended timeout. Eleven concurrent postings completed during the reads workload. The latest dataset had 2 tenants, 10,102 journals, 50,229 journal lines, 50,229 GL rows and 2,000 sales orders; PostgreSQL was 16.4 on Linux/amd64, using an Intel Xeon 6975P-C runner, 4 reported CPUs, 16,464,306,176 bytes reported memory, and `GOMAXPROCS=2`.

No approved reconciliation performance budget was found. The measured plan no longer spills sort work, but the end-to-end p50 did not improve in this comparison. Obtain an approved latency budget and review repeated, controlled performance results before declaring performance acceptance or production capacity.

## Release effect

The reconciliation query change has passed the latest exact-head benchmark and disposable PostgreSQL acceptance, and received technical review. The repository owner authorized merge after technical review when required checks and protections pass. No accounting approval or performance budget is expressed by that authorization or by this report. CFR-010 remains OPEN and financial production remains NO-GO pending source reconciliation, qualified accounting approvals, and an approved performance budget.
