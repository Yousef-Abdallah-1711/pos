# Phase 2C Prioritized Backlog and Blockers

## Blockers

1. **Admin Orders query instrumentation/optimization:** P1 paths remain dirty in the Phase 2A `H:\ws02` worktree. Do not edit those paths until ownership is released. Source shows a per-row detail reload, but runtime query counts remain uncollected.
2. **Disposable runtime stack:** Phase 2B containers were removed. The Phase 2A seed/stack is guarded for another workspace and must not be used. No new P2C synthetic PostgreSQL/Redis stack was brought up, so authenticated API and browser measurements were not run.
3. **Authenticated browser state and route fixtures:** Synthetic storage state and stable IDs for Dashboard, POS, Inventory, Orders, Products, and KDS were not available for Phase 2C.

## Next work, in order

1. After the order-path ownership gate clears, implement the list-specific batch hydration in `04_ADMIN_ORDERS_N1_ANALYSIS.md`; instrument and benchmark equivalent fixtures before and after.
2. Build a disposable P2C-only production-mode stack from the recorded source change-set fingerprint. Do not attach Phase 2C measurements to Phase 2B's older fingerprint.
3. Measure authenticated Dashboard, POS, Inventory, Orders, Products, and KDS, including visible-content readiness, cold reload, warm reload, client navigation, API waterfalls, redirects, and errors.
4. Sample safe authenticated GET endpoints for Orders, Products, Inventory, Tables, KDS, restaurant status, menu, tenants, Reports, and Accounting. Add SQL/Redis timing only where instrumentation is available.
5. Use measurement results to select one independent request-deduplication, cache correctness, or loading-boundary change; preserve tenant/restaurant/user/permission dimensions and test invalidation/isolation.
6. Evaluate PostgreSQL indexes and Redis changes only after collecting real query plans and cache hit/miss/TTL/memory evidence.

## Explicitly not completed

No dashboard navigation waterfall, frontend cache isolation experiment, Redis optimization, PostgreSQL index change, page runtime expansion, endpoint runtime expansion, or performance budget was claimed. The Phase 2C coverage CSVs distinguish unmeasured routes from prior Phase 2B evidence.
