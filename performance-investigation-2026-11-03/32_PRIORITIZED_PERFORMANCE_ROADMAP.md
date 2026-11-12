# Prioritized Performance Roadmap

This roadmap prioritizes investigations and implementation candidates using static confidence, breadth, tenant safety and effort. No percentage or time-saving promise is made without production-like benchmark evidence.

## P0 — Verify before broad optimization

1. **Cache key tenant scope:** prove Restaurant Website order/tracking keys cannot reuse customer/order data after tenant/restaurant switching; then include all required context if absent. This is primarily a correctness/security gate.
2. **Admin order N+1:** confirm query count and response cost with realistic synthetic data; batch detail hydration and cap page size while preserving tenant predicates and response compatibility.

## P1 — High-breadth candidates

3. **Dashboard shell readiness:** split auth/session readiness from unrelated shell rendering only if the result never exposes protected data; measure route usable time across roles.
4. **Restaurant tenant middleware:** instrument the awaited tenant-status check and cache only after verifying disable/expiry freshness semantics.
5. **Backend request tracing:** correlate route duration with middleware/auth, DB acquisition/query, Redis, external calls and serialization, without tenant/user IDs as metric labels.
6. **HTTP API caching policy:** decide explicitly which responses may be cached; middleware present in source is not currently wired. Keep sensitive responses private/no-store unless a reviewed scope allows otherwise.

## P2 — Targeted opportunities

7. Add route loading boundaries for high-traffic uncovered routes based on readiness traces.
8. Run production bundle analysis for all four apps and split only measured heavy client chunks.
9. Audit query cancellation and DB pool wait for endpoints with real latency evidence.
10. Normalize report cache query keys and add invalidation only for demonstrated stale data.
11. Test polling visibility behavior, dedupe and jitter for KDS/reports/tracking.
12. Validate product sort allowlist and query plans; security correctness precedes performance tuning.

## P3 — Optional / evidence-dependent

13. Virtualize lists only where post-pagination client lists still exceed rendering budgets.
14. Optimize public asset caching/CDN only after deployed headers and hit rates are known.
15. Consider indexes only after `EXPLAIN (ANALYZE, BUFFERS)` on safe realistic synthetic SELECTs and write-cost review.

## Release gates

- Capture a production-mode baseline on safe isolated synthetic fixture.
- Include median/p75/p95 only with adequate sample counts; retain raw samples.
- Verify tenant boundary, stale-data, empty/error/loading states and rollback for each batch.
- Run prescribed test/build/deployment checks in the implementation task; they were not run in this investigation.
