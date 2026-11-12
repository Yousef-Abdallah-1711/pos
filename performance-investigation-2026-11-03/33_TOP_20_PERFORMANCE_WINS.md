# Top 20 Performance Wins (Evidence-Ranked Candidates)

These are opportunities, not implemented changes. “Expected benefit” is qualitative until benchmarked.

| Rank | Candidate | Scope | Root cause evidence | Expected benefit | Effort | Main regression risk |
|---:|---|---|---|---|---|---|
| 1 | Batch admin-order detail hydration | Admin order list/export | Static N+1 in handler | Fewer SQL round trips | M | tenant predicate/order/shape mismatch |
| 2 | Bound admin order page size | Admin order endpoints | No apparent max clamp + OFFSET | Bounded work/memory | S | client paging compatibility |
| 3 | Verify and scope order query keys | Storefront order/tracking | ID-oriented query keys | Prevent stale cross-context reuse | S | lower cache reuse |
| 4 | Instrument server-side tenant-status gate | Restaurant host routes | awaited middleware fetch | Identify / reduce blocking if material | S | freshness of access state |
| 5 | Define dashboard auth-ready condition | 79 dashboard templates | global hydration gating | Earlier meaningful shell if safe | M | auth flash/data exposure |
| 6 | Add provider-level request waterfall tracing | Go API | request totals but sparse DB spans | Attribute work accurately | M | metric overhead/cardinality |
| 7 | Verify React Query provider in dashboard | Dashboard theme hooks | hook found; provider not found by source scan | Ensure dedupe/cache is active | S | behavior change if wrapper missing |
| 8 | Add loading boundaries to top high-traffic routes | all apps | only 3 route-level files found | Better perceived navigation | M | layout shift/duplicate skeleton |
| 9 | Separate public shell from protected dashboard data | Dashboard | provider hides descendants | Reduce perceived wait | M/L | tenant/auth leakage |
| 10 | Measure and split production client chunks | four apps | large dev chunk is not production evidence | Reduce JS if bundle confirms | M | chunk loading regressions |
| 11 | Check Restaurant homepage theme fetch dedupe | Storefront | metadata and page fetch theme | Avoid duplicate backend work if confirmed | S | stale metadata/content mismatch |
| 12 | Make polling visibility-aware and deduplicated | KDS/reports/tracking | multiple polling hooks | Reduce idle/background traffic | M | stale active operational view |
| 13 | Profile public tenant-status cache freshness | Redis/middleware | 30s cache | Lower repeat read cost if endpoint hot | S | access revocation delay |
| 14 | Normalize financial report cache query keys | Reports | raw query in key | Reduce duplicate entries | S | changing query semantics |
| 15 | Define Redis outage/backpressure policy | cache users | process memory fallback and ignored Set errors | Predictable degradation | M | replica divergence/stampede |
| 16 | Trace query context/cancellation for hot paths | repositories | direct calls exceed context variants by grep | Stop work after request timeout | L | incomplete cancellation propagation |
| 17 | Add realistic query-plan fixtures | DB | sparse E2E; pg_stat_statements absent | Evidence-led indexes/SQL | M | fixture skew |
| 18 | Validate product sort allowlist | product list | string interpolation in ORDER BY | Safer predictable plans | S | sort option compatibility |
| 19 | Verify service-worker install/cache behavior | Dashboard | worker files exist; activation not established | Avoid stale chunks or offline regressions | M | stale app shell |
| 20 | Add route-ready browser conditions and repeated trials | all four apps | current load events insufficient | Reliable optimization baseline | M | synthetic data not representative |

No benefit percentage is asserted. See `31_PERFORMANCE_ROOT_CAUSE_REGISTER.csv` and per-domain opportunity CSVs.
