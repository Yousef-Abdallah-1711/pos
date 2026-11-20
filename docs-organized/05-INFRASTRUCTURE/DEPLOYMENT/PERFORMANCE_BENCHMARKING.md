# SERVIO Performance Benchmarking

**Platform-wide performance status: not validated.** SERVIO has no verified platform-wide capacity baseline or 10K readiness result in the repository. The historical aggregate values previously presented as v1.1.0/v2.0.0 measurements were hardcoded example data, not captured benchmark runs. They are withdrawn and must not be used as observed performance or release evidence.

## Evidence available in the repository

The repository contains a scoped admin-orders query benchmark at:

- `docs-organized/phase-2a/execution/ADMIN_ORDERS_BENCHMARK_BEFORE.json`
- `docs-organized/phase-2a/execution/ADMIN_ORDERS_BENCHMARK_AFTER.json`

Those artifacts record one narrow synthetic test: PostgreSQL 15 on disposable tmpfs with `pg_stat_statements`, Go 1.24.13, 75 synthetic orders with two items each, and 30 HTTP samples per page size. For page size 50, the artifacts record median/p95 HTTP latency of 291.768/504.784 ms before and 4.758/7.68 ms after; SQL statement counts were 102 and 3. This is historical, route-specific evidence. It does not establish concurrent throughput, production performance, a general API improvement, or capacity under a larger dataset. This change does not rerun that benchmark.

## Current load-test measurement workflow

The read-only k6 scenarios and measured-summary parser are documented in [LOAD_TESTING.md](./LOAD_TESTING.md). The scripts exercise the current health, public-menu, session, and Landing routes. `summarize_k6.py` derives aggregate values only from a real k6 summary export; required values are not populated with defaults.

For a useful comparison, keep the source commit, environment sizing, data set, concurrency/arrival rate, duration, resource telemetry, and raw k6 summary with the sanitized report. Run baseline and candidate against the same isolated environment and dataset. Treat each route and workload separately. Do not infer web vitals, cache hit rate, database capacity, or memory use from HTTP duration alone.

## Removed unsupported benchmark artifacts

The prior `backend/scripts/benchmark.ts` generated fixed comparisons from assumed v1/v2 values, referred to absent build directories, and was not wired to a runnable package command. The adjacent `baseline-metrics.json` and `optimized-metrics.json` were timestamped example figures rather than test output. They have been removed so they cannot be mistaken for evidence. The scoped admin-orders artifacts above are preserved.

## Required future measurements

Before publishing capacity or performance claims, run the workload in a dedicated environment and record at least:

- Environment sizing, deployed source/image identity, database and Redis configuration.
- Dataset and synthetic account/session/tenant counts.
- Request mix, arrival rate, VU allocation, duration, warm-up, and repetitions.
- Request count, throughput, failed-request rate, p50, p95, and p99 latency.
- CPU, RAM, network, database pool/query behavior, and Redis memory/latency.
- CDN cache-hit and origin traffic for Landing tests.
- Bottleneck, saturation point, aborts, and recovery behavior.

Run separate measurements for registered-account lookup, idle session storage, authenticated session validation, active protected API traffic, cached Landing traffic, login bursts, Redis outage/recovery, tenant isolation, and KDS/WebSocket traffic where supported. A VU count alone does not define any of these workloads.

No arbitrary performance SLO is marked achieved here. Proposed targets must be reviewed against product needs and measured service capacity before being adopted. Public staging and production must not be used for 10K load without explicit authorization and a capacity-safe plan.
