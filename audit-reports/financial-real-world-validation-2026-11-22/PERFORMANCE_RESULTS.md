# Performance results

**Acceptance: PENDING.** No approved performance budget was found. Benchmark CI passing is not operational performance approval.

## Existing #110 10,000-journal comparison

Source: `../financial-cfr010-remediation-2026-11-21/PERFORMANCE_RESULTS.md`.

| Large restaurant warm reconciliation | Baseline | Candidate | Change |
| --- | ---: | ---: | ---: |
| p50 | 174.044 ms | 180.726 ms | +6.682 ms (+3.8%) |
| p95 | 175.441 ms | 182.565 ms | +7.124 ms |
| max | 175.524 ms | 182.660 ms | +7.136 ms |

End-to-end warm latency regressed. The line-to-GL query-plan execution improved from 160.660 to 148.041 ms. Plan shape switched to indexed lookups and removed temporary I/O (baseline 3,544 temp blocks read / 713 written; candidate zero). Plan improvement does not negate the end-to-end regression. The baseline and candidate each used 10,102 journals overall, with a 10,000-journal large scope.

## Disposable 1,000-order benchmark

`financial-benchmark-results-1000.json` reports PostgreSQL 15.19 on a disposable Linux/AMD Ryzen 7 5700 runner, 16 CPUs, 128 MB shared buffers, 4 MB work memory, 2 tenants, 1,102 journals, 5,229 lines/GL rows, and 1,000 orders. Posting seed took 57,909 ms (17.27 journals/s). Large-scope warm reconciliation was p50/p95/max 70.915/92.638/98.042 ms. Reads during posting were 94.909/141.876/151.989 ms. The 1 ms SQL timeout canceled with SQLSTATE `57014` as configured. The integration test completed in 94.13 seconds; its seed invariant logged 11,000,000 cents per side. The structured benchmark reports 11,032,700 cents on each side.

## Decision

The prior 10,000-journal candidate had better query-plan execution and eliminated temp I/O, but slower full-operation latency. No approved p50/p95/max target exists. Treat measurements as descriptive, synthetic evidence only. The first 10,000-journal attempt that failed before measurement is not a benchmark pass. Performance acceptance and production capacity remain PENDING.
