# Reconciliation Performance Decision

## Evidence

Measurements below come from the exact-SHA benchmark evidence linked in PR #110 (run `44979879216`) and its current PR description. The baseline and candidate runs reported the same Intel Xeon 6975P-C CPU model, PostgreSQL 16.4, and dataset row counts, but they were separate hosted runs rather than a controlled paired comparison.

| Measurement | Baseline | Candidate | Interpretation |
|---|---:|---:|---|
| Warm large-scope reconciliation p50 | 174.044 ms | 180.726 ms | Candidate +6.682 ms / +3.8%; end-to-end regression |
| Warm large-scope reconciliation p95 | 175.441 ms | 182.565 ms | Candidate +7.124 ms |
| Warm large-scope reconciliation max | 175.524 ms | 182.660 ms | Candidate +7.136 ms |
| Isolated journal-line-to-GL query | 160.660 ms | 148.041 ms | Candidate query plan faster in its separate run; not evidence of end-to-end improvement |
| Temporary blocks for line-to-GL query | 3,544 reads / 713 writes | 0 reads / 0 writes | Candidate removes observed temp I/O |

The candidate uses `idx_general_ledger_journal_line` for 50,022 index probes. The prior plan used a broad sort and merge path after a sequential GL scan. The benchmark dataset contained 10,102 journals, 50,229 journal lines, and 50,229 GL rows across two tenants; the large restaurant scope contained 10,000 journals and 50,000 lines/GL rows. PostgreSQL caches were not restarted/flushed, so the first candidate operation is not a verified cold-cache measurement.

## Decision

- Performance CI passing proves the benchmark ran; it does not establish operational acceptance.
- The query-plan and temporary-I/O change is promising, but the available large-scope end-to-end measurements show a 3.8% p50 regression and are not controlled paired runs.
- No approved reconciliation latency budget or acceptance threshold was supplied. Therefore reconciliation performance acceptance is **PENDING**, and this evidence does not approve the candidate for operational financial use.
- PR #110 is benchmark instrumentation plus a query aggregation change. Merging measurement infrastructure, if approved by an independent reviewer and required checks, must remain separate from approval of reconciliation performance for production use.
- Repeat a controlled comparison on the same runner class and environment with a pre-approved performance budget before declaring performance acceptance.

Source: [PR #110](https://github.com/Yousef-Abdallah-1711/pos/pull/110), exact-SHA benchmark run [44979879216](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44979879216), baseline `14fa8640191ae08c2b0789b4d05e370dd3023834`, candidate code `0a1aca0374ce8b357e1709f43ac6e1e242b7f575`.
