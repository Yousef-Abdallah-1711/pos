# Phase 2C Recovery Measurements

Date: 2026-11-04 (Africa/Cairo)

## Continuation recovery addendum — 2026-11-04

This addendum records the follow-up runtime recovery after the measurements below. It supersedes the earlier open-status statements for the Dashboard React error, Storefront readiness, and Redis tenant isolation. The CSV measurements below remain the original pre-recovery data unless explicitly stated otherwise.

### Dashboard React #300 recovery

- Root cause: `Sidebar` returned early for blocked/loading states before calling its `useMemo` hook. When permissions or tenant status changed, React observed fewer hooks on the next render and raised error `#300` (“Rendered fewer hooks than expected”). The fix moves both conditional returns below the memo so hook order stays stable.
- Change: `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` (SHA-256 `f22ba0c8bbdc799611a51f19ab3cc1a7d6d210688c9331088652d127833cebae`); regression coverage: `frontend/apps/dashboard/src/components/layout/Sidebar.test.tsx` (SHA-256 `ebfbd86babccdd67c91de54233f8a08aee470ed4696da764da968be46641ce48`) for loading-to-ready and active-to-blocked transitions.
- Focused Vitest result: **2/2 passed**. Production Dashboard image build succeeded: `pos-p2c-dashboard-sidebar-fix:20261104`, image digest `sha256:68808c9b6dd42f71279c11eecb722d857cb88775b9337375d7ea4effb056eee6`. The build emitted existing dynamic-dependency warnings from the component registry.
- Fixed-build authenticated browser replay: **8/8 previously affected routes** (`/en/dashboard`, `/en/dashboard/orders`, `/en/dashboard/inventory`, `/en/dashboard/products`, `/en/dashboard/tables`, `/en/dashboard/hr/employees`, `/en/dashboard/accounting/advanced`, `/en/dashboard/settings`) rendered main content on two samples each, returned HTTP 200, had no uncaught JavaScript exceptions, and did not redirect unexpectedly. POS and KDS were also checked. Authenticated API requests in this replay returned HTTP 200. The synthetic owner session was tenant `3`, restaurant `3`; no credentials are stored in this report.
- The existing `dashboard-page-measurements-final.csv` was recorded from commit `b62ffca` before the Sidebar fix (fingerprint `1fb1789712416e89abd211047f3ee4b9a99bde62`); its `javascript_errors=1` entries are historical and are not post-fix results. A separate fixed-build replay verified recovery, but its TTFB/FCP/LCP/CLS and transfer-size values were not saved as a comparable replacement dataset. Do not use the old timings as post-fix performance evidence or claim a page-load improvement from this recovery.

### Storefront and Redis follow-up

- Storefront routing required a tenant subdomain. On `p2c-restaurant-a.localhost:23203`, the homepage rendered and `/en/menu` showed the expected empty state (`0 items` / `No items found`), with no console errors. The earlier generic `localhost` navigation did not resolve a tenant. This verified routing/readiness for the synthetic tenant; it did not change publication rules or establish populated-menu performance.
- Redis isolation was verified with synthetic tenants `3` and `6`: authenticated status responses returned their respective tenant IDs and Redis held separate `tenant_access:3` and `tenant_access:6` keys, each with a 30-second TTL. No cross-tenant value leakage was observed. The earlier expiry probe also confirmed TTL expiration. No broader response caching was added.

### Admin Orders and runtime status

- Admin Orders remains blocked by active, uncommitted Phase 2A ownership in `H:\ws02` of the overlapping handler, repository, use-case, and integration-test paths. No Admin Orders source was changed and no before/after improvement is claimed. The recorded baseline remains 105 SQL statements per page-size-50 request, with warm median 146.87 ms and p95 215.77 ms. Batch hydration must wait until ownership is released or an integration candidate can be reviewed safely.
- The patched verification runtime was stopped after validation. The original Compose-owned container `pos-p2c-recovery-20261104-dashboard-1` was restarted and returned HTTP 200 at `http://localhost:23202/en/dashboard`. The separate verification container remains stopped. API, PostgreSQL, Redis, Landing, Storefront, and Super Admin containers were not restarted.
- Branch remains `codex/phase2c-notification-dedup`; commits `9e3614b`, `2d86ca3`, `7a62f40`, and `b62ffca` remain ancestors of HEAD. Current worktree changes include the Sidebar fix and focused test; the untracked `scripts/performance/` package is preserved. Browser artifacts under `.playwright-mcp/` are also present and were left intact.

## Outcome

This run collected authenticated API, browser, SQL, PostgreSQL, and Redis evidence on a disposable Phase 2C stack. It verified the Admin Orders N+1 pattern and delivered one safe frontend request-coalescing improvement with before/after measurements. No push, merge, deploy, or production/development database action occurred.

## Runtime and source identity

- Compose project: `pos-p2c-recovery-20261104`; all ports were loopback-bound. PostgreSQL and Redis were isolated containers. PostgreSQL was not restarted after fixture restoration.
- Synthetic fixture scope: tenant `2`, restaurant `2`; 6 categories, 62 products, 80 tables, 400 orders, 800 order items, and 400 KDS rows.
- API and the unchanged Landing, Storefront, and Super Admin images carry source commit `9e3614b605b4d1691835c37c0a088636ccf6abe2`.
- Final Dashboard image carries source commit `b62ffcaa9c49ad4e2f25720c9c05b8c3165405e3`, source tree fingerprint `1fb1789712416e89abd211047f3ee4b9a99bde62`, and image digest `sha256:13d5cb3fa5051d443832bd92bebbf99ef8b0ff4d8e1128b8e0b569090aaead92`.
- Frontend build: production `Dockerfile.prod`, Node 20 Alpine, pnpm 9.0.0, Next.js 15.5.9, `NEXT_IGNORE_BUILD_ERRORS=false`. Dashboard build passed. Existing warnings remain for dynamic dependency expressions in the component registry and stale Browserslist data.
- Exact source branches and commits are listed below. The measurement scripts remain untracked in the Phase 2C worktree; evidence CSVs are in this recovery directory.

## Authenticated browser measurements

The browser sampler used the synthetic Playwright storage state, two samples per Dashboard and Super Admin route, and three per Storefront route. It recorded navigation milestones, visible-selector readiness, API waterfalls, transferred resource bytes, layout shift, long tasks, console errors, uncaught JavaScript exceptions, failed requests, and final paths. Query strings are represented only by truncated SHA-256 fingerprints in the updated waterfall.

- Dashboard: `/en/dashboard`, POS, Orders, Inventory, Products, Tables, KDS, HR employees, Accounting advanced, and Settings. All navigation responses were HTTP 200. POS and KDS were ready with no JavaScript exception. Eight routes emitted the same minified React error `#300` in both samples; content readiness was inconsistent or timed out on those routes. The measured root route went from 23/22 API calls to 21/20 after the changes below. Cold/warm FCP and TTFB varied substantially, so this run does not claim a page-load time improvement.
- Super Admin: `/dashboard` and `/dashboard/audit-logs`, two samples each; both reached visible main content, returned HTTP 200, and recorded no JavaScript exceptions. The synthetic browser state was refreshed with the persisted user object required by the app.
- Storefront: `/en` and `/en/menu`, three samples each; both returned HTTP 200 but resolved to `/maintenance` with “We Are Currently Offline.” The storefront API endpoints were healthy; the page-level result is a fixture/configuration limitation, not a successful menu-page measurement.
- The first-load resource-byte value is present; repeat samples report zero transferred bytes for cached resources. Do not interpret those zeros as zero resource size.

### Frontend before/after evidence

Equivalent two-sample, authenticated, production Dashboard runs used the same routes, storage state, synthetic IDs, browser, and readiness selectors.

| Evidence | Before (commit `9e3614b`) | Final (commit `b62ffca`) | Change |
|---|---:|---:|---:|
| Dashboard API calls across 10 routes × 2 samples | 190 | 186 | −4 (−2.1%) |
| `/api/v1/notifications/stats` calls | 18 | 16 | −2 (−11.1%) |
| `/api/v1/tenant/status` calls | 22 | 20 | −2 (−9.1%) |
| Exact repeated GETs in final query-fingerprinted waterfalls | — | 0 | Both observed duplicate classes removed |

Each Dashboard root-page sample shed one notification-stats request and one tenant-status request. The two notification-list calls remain because one requests the default list and the other requests `limit=5`; their query fingerprints differ. Page timings are not used as proof of improvement because of low sample counts, warm-cache differences, and the React exceptions.

## Authenticated business API measurements

Each configured route received one warm-up plus 25 measured sequential GETs. All 15 route configurations returned `200:26`; no unauthorized or failed request was counted as coverage. Owner routes used tenant `2` / restaurant `2`; Super Admin used its synthetic token; storefront endpoints were public and used the fixture slug.

| Route | Records | Warm p50 | Warm p95 |
|---|---:|---:|---:|
| Admin Orders, page size 1 | 1 | 14.88 ms | 17.19 ms |
| Admin Orders, page size 10 | 10 | 50.81 ms | 99.94 ms |
| Admin Orders, page size 25 | 25 | 83.25 ms | 142.78 ms |
| Admin Orders, page size 50 | 50 | 146.87 ms | 215.77 ms |
| Products | 20 | 12.36 ms | 15.44 ms |
| Categories | — | 8.36 ms | 10.05 ms |
| Tables | 80 | 10.29 ms | 14.04 ms |
| KDS queue | 300 | 23.19 ms | 36.70 ms |
| Sales report | — | 7.91 ms | 22.81 ms |
| Inventory report | 62 | 5.29 ms | 6.78 ms |
| Super Admin dashboard stats | — | 5.33 ms | 6.28 ms |
| Super Admin tenants | 2 | 14.00 ms | 31.95 ms |
| Super Admin Audit Logs | 20 | 8.05 ms | 11.36 ms |
| Storefront menu | 62 | 7.33 ms | 10.38 ms |
| Storefront homepage | — | 5.26 ms | 8.86 ms |

The results are local synthetic-data measurements, not production capacity estimates.

## Admin Orders SQL and PostgreSQL

For each page size, PostgreSQL statistics were reset before a batch of 26 HTTP requests (one warm-up plus 25 measured). The paginated list and count statements each ran 26 times regardless of page size. Order-detail and order-item statements scaled linearly with returned rows:

| Page size | Detail-query calls | Item-query calls | N+1 statements per HTTP request |
|---:|---:|---:|---:|
| 1 | 26 | 26 | 2 |
| 10 | 260 | 260 | 20 |
| 25 | 650 | 650 | 50 |
| 50 | 1,300 | 1,300 | 100 |

This confirms N+1 hydration: one tenant/restaurant-scoped detail fetch and one item fetch per result. The size-50 batch issued 2,730 total SQL calls, or 105 per request including list, count, permission, and role lookups.

On the 400-order synthetic fixture, `EXPLAIN (ANALYZE, BUFFERS)` for the paginated order query used sequential scans and a hash join, returning 50 rows in about 1.03 ms with 30 shared-buffer hits. That is reasonable at this small fixture size and does not justify a speculative index. The order-item plan used `idx_order_items_order_id` and returned two fixture rows in about 0.72 ms. The measured end-to-end growth is consistent with repeated application/database round trips rather than one slow query plan.

The database had 7 backends (1 active, 6 idle) with `max_connections=100` at inspection; there was no pool saturation evidence. A batch-hydration fix is deferred because the active WS02/Phase 2A worktree still owns the overlapping Admin Order handler, repository, use case, and integration-test paths.

## Redis

Redis had measurable application traffic: after resetting disposable Redis statistics, four authenticated Products GETs produced 3 hits, 1 miss, no evictions, and a single `tenant_access:2` key with roughly 30 seconds remaining. The middleware builds this key as `tenant_access:<tenantID>` and sets a 30-second TTL. A later check found it expired, confirming TTL-based invalidation. The key is tenant-scoped in code; only one authenticated tenant identity was available in this fixture, so a two-tenant runtime isolation probe was not performed. No broad Redis response cache was introduced.

## Implemented optimization and verification

- Notification-list cache keys now include tenant, restaurant, and user identity; simultaneous same-scope requests are coalesced. Unit cases verify separation across each identity dimension and separate in-flight scopes.
- Dashboard overview and notification hook now share the same scoped notification-stats request.
- Dashboard overview and tenant-status store now share the same scoped tenant-status request.
- Focused Vitest suite: 4 tests passed. Dashboard TypeScript check passed. The final production Dashboard image built successfully with type checking.
- Audit Logs correctness fix remains a separate local commit and its focused suite previously passed: 16 tests. It is not counted as a performance improvement.

## Git and artifacts

- Current branch: `codex/phase2c-notification-dedup`
- `2d86ca3` — `perf(notifications): coalesce scoped read requests`
- `7a62f40` — `perf(dashboard): share notification stats request`
- `b62ffca` — `perf(dashboard): coalesce tenant status reads`
- `9e3614b` — `fix(super-admin): handle malformed audit details`
- All commits are local. Nothing was pushed or merged. The working tree contains only the preserved untracked `scripts/performance/` measurement package.
- API latency CSVs, per-size SQL TSVs, browser page CSVs, and sanitized waterfalls are in this directory.

## Remaining blockers and next batch

1. After Phase 2A releases the Admin Orders files, implement tenant/restaurant-preserving batch hydration, then rerun the same four page sizes and focused isolation/regression checks. It remains the largest verified backend bottleneck.
2. Capture a comparable post-Sidebar-fix authenticated Dashboard vitals dataset (including TTFB, FCP, LCP, CLS, meaningful-content readiness, transfer sizes, and API waterfalls). The previous CSV belongs to the pre-fix source fingerprint and is not a post-fix timing baseline.
3. The Storefront resolves on the tenant subdomain and reaches the empty-menu state. Seed/publish synthetic content only if a populated-menu journey is needed, then record comparable homepage and menu measurements.
4. Redis isolation across synthetic tenants 3 and 6 passed; keep the existing narrow tenant-status cache and require targeted evidence before expanding shared cache behavior.
