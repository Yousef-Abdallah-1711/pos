# F-078 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Theme edits and activation do not refresh the public theme: the 5-minute server cache is never invalidated (and the storefront adds a 1-hour page cache).

## 2. Current behavior (verified)
Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300

## 3. Required behavior
The public theme reflects an activation or save immediately (or within seconds)

## 4. Reproduction conditions
Tenant 4: prime GET /public/restaurants/audit-c-diag/theme, save #00ff00 in the editor, activate the theme, poll the public API every 45 s

## 5. Verified evidence
- evidence/network/RUN-004/PH-08/AUD-R4-THEME-CACHE-TTL2.json
- evidence/network/RUN-004/PH-08/AUD-R4-THEME-PAYLOAD.json
- tests/runtime/ph08_r4_theme_cache_ttl2.mjs

## 6. Root cause and confidence
CONFIRMED — invalidatePublicThemeCacheByRestaurant is defined (handler/http/theme_public_cache.go:45) but has no caller; GetThemeOnly caches the payload for publicThemeCacheTTL = 5 min (public_homepage_handler.go:117-144)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/theme_public_cache.go`:45, `backend/internal/handler/http/public_homepage_handler.go`:117, `frontend/apps/restaurant-website/app/[locale]/page.tsx`:75, `backend/internal/handler/http/public_homepage_handler.go`:117-144
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/restaurant-website/app/[locale]`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: Theme Builder, public theme API, storefront

## 8. Impact
- Financial: 0
- Operational: Theme Builder, public theme API, storefront workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: call invalidatePublicThemeCacheByRestaurant(restaurantID) (backend/internal/handler/http/theme_public_cache.go:45) from the theme activate, update, delete and restore handlers after the database commit; replace the storefront fixed revalidate (frontend/apps/restaurant-website/app/[locale]/page.tsx:75,131, 3,600 s) with tag-based revalidation triggered by the same events (or reduce it to 60 s); keep Cache-Control max-age small (<= 60 s) for the theme endpoint.

## 10. Alternative designs
A: invalidate on write + short HTTP cache (recommended). B: shorten the TTL only (still stale for the TTL).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-077

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/public_homepage_handler.go, backend/internal/handler/http/theme_public_cache.go, frontend/apps/restaurant-website/app/[locale]/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
After activate or save the public API returns the new colour within 2 s (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- After activate or save the public theme endpoint returns the new values within 2 s (test polls at 1 s).
- The storefront shows the new theme within 60 s.
- No stale payload survives a deleted theme.
