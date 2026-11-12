# Executive Performance Investigation — POOS

**Status: substantial read-only baseline, not a complete runtime performance certification.** The source inventory is broad and the runtime sample is deliberately narrow. Findings that depend on authenticated usage, production rendering, query plans or realistic tenant data are explicitly unmeasured.

## Counts and direct answers

1. **Frontend page files:** 110 total: Dashboard/POS/Inventory/KDS 79; Landing 10; Restaurant Website/Storefront 12; Super Admin 9. The complete route templates and file paths are in `02_COMPLETE_FRONTEND_PAGE_INVENTORY.csv`. Dynamic and locale patterns can resolve to multiple URLs; counts are source page files.
2. **Rendering dependencies:** shared layouts/providers are summarized in `03_PAGE_COMPONENT_DEPENDENCY_MAP.md` and `04_RENDERING_AND_LOADING_BOUNDARIES.md`. The dashboard auth provider and restaurant tenant-status middleware are the most evident shared page-opening dependencies from source.
3. **Endpoint registrations:** source manifest lists 360 registration occurrences, 335 unique method/path patterns. Runtime count is conditional on dependencies/handler construction. See `08_COMPLETE_BACKEND_ENDPOINT_INVENTORY.csv`.
4. **Page-to-endpoint and endpoint-to-handler mapping:** not complete. Static common request flows and remaining per-page mapping fields are in `05_PAGE_TO_API_DEPENDENCY_MATRIX.csv`, `06_FRONTEND_BACKEND_DEPENDENCY_GRAPH.md`, and `11_ENDPOINT_CONSUMER_MAPPING.md`. Do not treat blank/unmapped consumers as unused.
5. **Endpoint response contracts:** only health runtime size/status is measured; endpoint-specific contracts are not fully extracted. See `09_ENDPOINT_RESPONSE_CONTRACTS.md`.
6. **Endpoint latency:** GET `/api/v1/health` measured 25 times in isolated development E2E: median 2.18 ms, p75 2.46 ms, p95 4.94 ms, max 126.93 ms, 69 response bytes, 25/25 HTTP 200. This health result is not representative of business endpoints or production. All other unique route patterns were not exercised.
7. **Slowest verified page:** Landing `/en` first isolated-context navigation had the highest observed document timings: responseStart 10.259 s, DOMContentLoaded 11.279 s, load 12.100 s. Repeat was 2.277/2.378/3.075 s. Development compilation/HMR is plausible but not proven as root cause. No other page has a comparable valid timing; it is not asserted to be platform-wide slowest.
8. **Slowest verified endpoint:** only health was measured; its maximum sample was 126.93 ms (first call), p95 4.94 ms over 25 sequential calls. No other endpoint comparison exists.
9. **Browser/API/DB/Redis/network split:** unavailable for page/business endpoints. Browser document events were captured for Landing only; health was measured host-to-loopback. No spans break down handler, database, Redis, serialization or server render.
10. **Blocking and duplicate requests:** static candidates include dashboard auth hydration, Restaurant Website tenant-status gate, theme metadata/page fetch, and background polling. No full waterfall or duplicate count is established.
11. **Skeletons:** 3 dashboard `loading.tsx` files; 0 route loading files in Landing, Restaurant Website and Super Admin. This does not exclude component/provider skeletons. Per-page visual quality is not verified.
12. **Caching:** Restaurant server fetch revalidation, client query cache, optional Redis and local caches exist; API cache middleware is not wired in. Key scoping/freshness risks and source evidence are in `21`/`22`.
13. **Redis:** actual consumers are narrow (tenant status and financial reports); runtime keyspace was empty. No claim of hit-rate benefit.
14. **SQL/indexes:** static admin-order N+1 is a strong code-level candidate. Product listing has capped page size and dynamic ORDER BY interpolation. No query plans, production row counts or index recommendations are asserted.
15. **Largest opportunities:** admin-order batch hydration; dashboard auth-ready rendering; tenant-status middleware measurement; tenant-safe cache keys; route-level loading; comprehensive tracing. See `33_TOP_20_PERFORMANCE_WINS.md`.
16. **Safe multi-tenant improvements:** preserve tenant/restaurant scoping; validate cache keys and pagination; never use shared-cache data for user-specific scopes without a reviewed key. Priorities and rollback path are in `32` and `34`.
17. **Implementation order:** instrumentation and tenant/cache correctness → admin orders/paging → dashboard/storefront readiness → cache and query tuning only with measured evidence.
18. **Evidence still missing:** production build/page runs; authenticated synthetic role journeys; all page readiness; all endpoint status/size/latency/response schemas; per-endpoint DB/Redis data; query plans/counts; HAR/traces; full frontend consumer map; source/manifest runtime reconciliation; mobile/slow-network/failure tests.

## Preservation

No application source/config/schema/migration/infrastructure file was intentionally edited. No data mutation, service restart, test/build, cache clearing, payment, message, commit, push or deploy was performed. The only runtime calls were browser GET navigation/resource loads and 25 GET health requests to the verified isolated WS05 disposable E2E stack. Existing dirty changes in both worktrees were preserved. The only new repository content is this dedicated report directory; session-generated Playwright snapshot/console evidence was relocated into its trace-evidence subfolder.

For exact limitations and counts see `35_COVERAGE_AND_EVIDENCE_REPORT.md`; for source and measurements begin with `01_SYSTEM_ARCHITECTURE_AND_RUNTIME.md`, `02_COMPLETE_FRONTEND_PAGE_INVENTORY.csv`, `08_COMPLETE_BACKEND_ENDPOINT_INVENTORY.csv`, and `10_ENDPOINT_RUNTIME_PERFORMANCE.csv`.

## Phase 2B — Production Runtime Baseline (2026-11-03)

Phase 2B adds a substantial, authenticated runtime sample in the isolated production-build stack. It does not certify all routes or represent production-scale behavior. The Dashboard build completed on the recorded source fingerprint A00BD5F7...; Next.js resolved to 15.5.9, and shared first-load Dashboard JavaScript was 103 kB. The synthetic fixture contained two Phase 2B tenants with two restaurants each plus one bootstrap tenant, 24 users, 62 products, 400 orders and 80 tables (full cardinalities are in the Phase 2B baseline).

### Measured results

- Authenticated Dashboard home: 5 hard reloads; median FCP/LCP 472 ms, meaningful content 436 ms, 23 successful API calls and 73 resources. CLS was 0.1059 in all five samples; a three-column section disappeared near 647 ms. No long tasks were observed.
- Admin orders: 25 sequential authenticated GET samples for each page size 1, 5, 25 and 50. Warm median rose from 12.9 ms to 91.2 ms; warm p95 rose from 15.1 ms to 117.4 ms; response sizes were 417 B to 18,055 B. Source shows list/count work plus one order and item hydration per row. This is a strong N+1 candidate with observed latency scaling, but SQL query counts were not captured. The synthetic EXPLAIN plans used existing tenant/order indexes; no index change is supported by this data.
- Tenant isolation: two synthetic tenants each returned their 21 products with no ID overlap. A mixed 20-request probe yielded global Redis delta of 19 hits and 1 miss. This does not establish per-endpoint cache rates.
- Authenticated route sweep: Dashboard 13 static templates and 6 fixture-resolved dynamic templates returned 200. Storefront menu displayed 21 synthetic products; cart and checkout showed empty-cart states. Super Admin dashboard, tenant list/detail, users, payments and settings returned 200.
- Coverage is 34 of 110 page templates observed; 76 remain NOT_TESTED. Runtime observation does not imply a passing page; see the readiness CSV.

### Failures and interpretation

The anonymous Dashboard navigation timed out on an 8-second Playwright networkidle condition in 3/3 runs, but redirected to and rendered the Landing login quickly. A prefetch to the linked forgot-password route returned 404. This is a reproducible harness timeout plus a broken route, not evidence of an authenticated Dashboard hang.

Super Admin Audit Logs returned HTTP 200 and 9 synthetic rows, but 2 rows had null details. The frontend calls Object.keys on that nullable value and crashes. This is a reproducible functional failure, not a performance finding. Landing super-admin login displayed Access Failed; its cause was not established.

### Evidence boundaries and next actions

The Phase 2B report files under phase2b/ include browser and API matrices, the full 110-template readiness CSV, waterfall summary, PostgreSQL/Redis findings, root-cause reports, and a 20-item evidence-labeled backlog. The 20 items are opportunities to verify or fix, not 20 verified performance bottlenecks. No application optimization was implemented, so there are no before/after claims. Per-endpoint query counts, handler/DB/Redis timing splits, cold startup, consistent per-route Web Vitals, INP, mobile/slow-network, full sanitized HAR, realistic large-tenant behavior, and several page families remain unmeasured.

The Phase 2B disposable Compose project was removed after verifying all seven containers carried the exact project label and disposable=true. No named project volumes existed; containers and the project network were removed, while all five uniquely tagged Phase 2B images were retained. Cleanup did not target Phase 2A, pos-dev, or unrelated projects. See phase2b/15_CLEANUP_VERIFICATION.md.

Source boundary: API measurements describe the production image built at fingerprint A00BD5F7.... The final checkout later gained changes in backend/cmd/api/main.go (including route and middleware/lifecycle changes) and the Landing login page. Those Phase 2A edits were not mounted into the containers. Repeat endpoint profiling after rebuilding before comparing to the later worktree state; see phase2b/14_RUNTIME_SOURCE_BOUNDARY.md.
