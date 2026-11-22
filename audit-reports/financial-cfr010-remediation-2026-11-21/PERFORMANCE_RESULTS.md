# Performance acceptance results

**Status: NOT RUN for journal/GL reconciliation.** No latency, query plan, timeout, or throughput result is claimed.

The hosted workflow's PostgreSQL reconciliation test validates correctness on synthetic anomaly fixtures but does not create a representative high-volume journal corpus or capture `EXPLAIN (ANALYZE, BUFFERS)`, p50, p95, worst latency, or concurrency performance. The full acceptance job passing is not performance evidence.

An existing opt-in benchmark exists at `backend/cmd/api/financial_benchmark_integration_test.go`; it requires `POS_FINANCIAL_BENCHMARK=1` and an output path on a disposable benchmark database. Run it separately and add a reconciliation-specific benchmark before approving this aggregate scan for broad operational use. Record dataset journal/line/GL counts, PostgreSQL version/configuration, index plan, warm-up/repetition method, p50/p95/max, concurrency, and timeout behavior. Do not change indexes or assert a performance gate until measurement and the existing project gate are established.

## Release effect

Performance acceptance is a remaining engineering gate for full CFR-010 closure. The repository method has no public endpoint in this candidate, limiting exposure, but this does not replace measurement before operational use.
