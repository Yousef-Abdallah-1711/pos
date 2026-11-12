# Implementation Batches and Dependencies

| Batch | Work | Prerequisite | Validation | Rollback / safety |
|---|---|---|---|---|
| 0. Baseline instrumentation | Route-ready browser checks; request-to-DB/Redis spans; production bundle baseline | Synthetic data and endpoint owner review | repeated cold/warm user journeys, low-cardinality metrics, raw traces redacted | Remove instrumentation or disable sampling; no behavior change |
| 1. Correctness and safety | Verify tenant-scoped query keys; verify product sort allowlist; endpoint page-size contract | Batch 0 and security review | cross-tenant synthetic tests, pagination boundaries, input tests | Keep prior query-key/cache policy behind feature flag if needed |
| 2. Admin orders | Batch hydrate list/export; enforce max page size; inspect tenant/date indexes | Query count baseline and response schema capture | compare records/order/items; query count; realistic fixture load | Feature flag / restore prior repository path |
| 3. Frontend readiness | Dashboard shell/auth state redesign; high-traffic loading skeletons | explicit auth-ready condition, provider wiring validated | authenticated owner/cashier/admin traces, no protected-data flash, CLS/keyboard | revert boundary/provider behavior independently |
| 4. Storefront navigation | Tenant status timing/freshness, theme fetch dedupe, progressive public content | synthetic tenant lifecycle and cache freshness policy | tenant active/expired transitions; SEO output; route timings | restore middleware gate and server fetch policy |
| 5. Cache policy | Redis failure model, key normalization/invalidation, only approved cacheable response headers | Batch 0 hit-rate and staleness evidence, tenant review | stale-data and outage tests, tenant isolation, hit/miss counters | remove cache wrapper / shorten TTL / private no-store |
| 6. DB/index follow-up | Query tuning and indexes from plans | realistic row cardinalities and write-rate baseline | EXPLAIN buffers; load tests; migration rollback validation | reversible index/migration plan; do not deploy speculative indexes |
| 7. Production verification | All journeys, slow network, failures, mobile, staged smoke checks | prior batches complete | p50/p75/p95 with sample counts, build/test/rollout checks | documented rollback by application and schema |

Avoid parallel database load tests with browser measurements; they would distort timing. Keep payment and mutation flows mocked/fixture-only.
