# Phase 2B Coverage, Blockers, and Backlog

## Coverage
The readiness CSV classifies all 110 source page templates. 33 templates have a runtime observation, including six fixture-resolved Dashboard dynamic paths; 77 are NOT_TESTED. Runtime observation includes functional failures and does not imply a healthy page. No route is classified PERMISSION_DENIED_EXPECTED without directly observing denial. No route is assumed FIXTURE_MISSING without evidence.

## Blockers and unmeasured dimensions
- No full sanitized HAR or Chrome trace; no server, DB, Redis, or serialization spans.
- Cold process startup and first-request timing, per-app numeric page timing, INP, mobile/slow network, failure injection, and resource byte breakdowns remain unmeasured.
- Synthetic DB is small; current plans do not support index advice.
- Runtime query counts for varying order page sizes were not instrumented. N+1 is a static candidate with latency correlation, not a measured DB round-trip curve.
- Super Admin list and menu API response sizes are not retained.
- Landing forgot-password is 404; Super Admin audit logs crashes client-side; Landing super-admin login showed Access Failed.

## Ranked implementation and verification backlog
| Rank | Item | Evidence class | Next step |
|---:|---|---|---|
| 1 | Batch admin-order hydration and measure SQL calls vs page size | Static candidate + runtime latency correlation | Instrument query counts; compare reversible batch loader |
| 2 | Fix nullable audit-log details rendering | Reproducible functional failure | Add null-row regression and page smoke |
| 3 | Reduce Dashboard layout shift | Verified UI metric | CLS 0.1059 in 5/5; section disappears around 647 ms |
| 4 | Replace networkidle wait with visible readiness assertion | Reproducible harness timeout | Login rendered while 8 s wait failed 3/3 |
| 5 | Repair/verify forgot-password route | Reproducible functional failure | Linked destination returns 404 |
| 6 | Attribute repeated Dashboard request families | Measured waterfall; cause unmeasured | Capture initiators and critical-path ordering |
| 7 | Bound admin-order page size | Static candidate | Verify client contract and OFFSET work |
| 8 | Add server spans for handler/DB/Redis/serialization | Unmeasured | Attribute endpoint latency |
| 9 | Profile Dashboard authentication hydration | Static candidate | Measure cold/warm auth-ready time |
| 10 | Measure tenant-status middleware critical path | Static candidate | Preserve access freshness while profiling |
| 11 | Verify order cache-key scoping by tenant/restaurant | Static candidate | Existing isolation pass covers product/status only |
| 12 | Measure endpoint-specific Redis hit rate and DB reduction | Unmeasured | Global counters cannot attribute benefit |
| 13 | Verify tenant-status invalidation and 30 s TTL | Partial source/runtime key evidence | Measure access revocation freshness |
| 14 | Measure route-level loading behavior on hot paths | Partial static | Determine actual perceived benefit |
| 15 | Capture production chunks/resource bytes for all apps | Partial build | Only Dashboard shared JS 103 kB recorded |
| 16 | Profile query behavior at larger order volume | Unmeasured | Current 100 orders/restaurant is small |
| 17 | Profile cancellation/context propagation on hot handlers | Static candidate | Instrument before broad changes |
| 18 | Validate product-sort allowlist and plans | Static candidate | Measure workload and options first |
| 19 | Exercise populated storefront journey with disposable fixtures | Unmeasured | Cart and checkout remained empty; no mutation |
| 20 | Resolve Landing super-admin login Access Failed | Functional failure; cause unresolved | Separate auth/config issue from performance |

No optimization or source edit was implemented, so there are no before/after performance claims. Top 20 are opportunities with evidence labels; they are not 20 verified bottlenecks.

Source boundary: two files differ from the recorded container build input (backend/cmd/api/main.go and frontend/apps/landing/src/app/[locale]/auth/login/page.tsx). Runtime results describe the built images, not later worktree edits; see 14_RUNTIME_SOURCE_BOUNDARY.md.
