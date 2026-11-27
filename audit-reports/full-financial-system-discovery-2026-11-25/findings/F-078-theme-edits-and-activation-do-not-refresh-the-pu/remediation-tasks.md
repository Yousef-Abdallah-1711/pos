# F-078 — Implementation phases and tasks (all NOT STARTED)

Finding: Theme edits and activation do not refresh the public theme: the 5-minute server cache is never invalidated (and the storefront adds a 1-hour page cache) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-078-01
- **Finding:** F-078
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-078 at code level
- **Existing source paths:** backend/internal/handler/http/theme_public_cache.go, backend/internal/handler/http/public_homepage_handler.go, frontend/apps/restaurant-website/app/[locale]/page.tsx, backend/internal/handler/http/public_homepage_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_r4_theme_cache_ttl2.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-078-02
- **Finding:** F-078
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-078
- **Existing source paths:** backend/internal/handler/http/theme_public_cache.go, backend/internal/handler/http/public_homepage_handler.go, frontend/apps/restaurant-website/app/[locale]/page.tsx, backend/internal/handler/http/public_homepage_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: call invalidatePublicThemeCacheByRestaurant(restaurantID) (backend/internal/handler/http/theme_public_cache.go:45) from the theme activate, update, delete and restore handlers after the database commit; replace the storefront fixed revalidate (frontend/apps/restaurant-website/app/[locale]/page.tsx:75,131, 3,600 s) with tag-based revalidation triggered by the same events (or reduce it to 60 s); keep Cache-Control max-age small (<= 60 s) for the theme endpoint.
- **Dependencies:** REM-F-078-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** After activate or save the public theme endpoint returns the new values within 2 s (test polls at 1 s). | The storefront shows the new theme within 60 s. | No stale payload survives a deleted theme.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-078-03
- **Finding:** F-078
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-078 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** After activate or save the public API returns the new colour within 2 s
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-078-04
- **Finding:** F-078
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-078
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
