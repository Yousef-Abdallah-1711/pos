# Performance forensics

**Status: PARTIAL.** No production-build comparison, repeated timing sample, p50/p95, tracing or query-plan profile was completed.

## Environment snapshot

Windows host, 16 logical processors, 34,220,703,744 bytes RAM. Host Node v24.16.0; Go unavailable on host, API container Go 1.24.13; Docker Engine 29.7.2; PostgreSQL 15.19; Redis 7.4.12; browser user agent Chrome 151. Compose NODE_ENV=development. Containers have no explicit CPU/memory limits. Point-in-time Docker stats after navigation: Dashboard about 2.46 GiB / 15.55 GiB container limit; API 46 MiB; PostgreSQL 115 MiB.

## Observations

- Prior route timing: one cold sample per route, about 2.7–10.6s. No median/p95 claim.
- Duplicate /api/v1/rbac/me/permissions and duplicate report/data GETs appeared on several route loads.
- Some invoice content appeared after an additional 2s wait; advanced financial page initially prompted refresh.
- Dashboard memory is elevated in this development container; production impact is not established.
- Development compilation overhead is not separated from product latency.

Next: at least 10 cold/warm samples in dev and production, navigation type, TTFB/LCP/interaction readiness, HAR/trace, backend spans and EXPLAIN ANALYZE BUFFERS on representative anonymized fixtures.

## Completion status — 2026-11-11

No new repeatable performance run was performed in this evidence-closure pass. Earlier report evidence remains one development-mode sample per route (about 2.7–10.6 seconds) plus duplicate permission/report requests. No production-build comparison, 10-sample distributions, browser traces/network waterfalls, query plans, or React profiling were collected. These timings are historical symptoms, not verified production bottleneck root causes. Performance status: **PARTIAL / required measurement matrix NOT TESTED**.
