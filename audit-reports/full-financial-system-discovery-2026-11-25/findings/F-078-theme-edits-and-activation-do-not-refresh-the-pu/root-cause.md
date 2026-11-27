# F-078 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): invalidatePublicThemeCacheByRestaurant is defined (handler/http/theme_public_cache.go:45) but has no caller; GetThemeOnly caches the payload for publicThemeCacheTTL = 5 min (public_homepage_handler.go:117-144)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/theme_public_cache.go`:45.
- Propagation: affected modules Theme Builder, public theme API, storefront; related findings F-077; observed consequence: Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75
- Cause versus symptom: the symptom is what users or reports show (Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-THEMEBUILDER-UI; scripts: tests/runtime/ph08_r4_theme_cache_ttl2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/theme_public_cache.go | 45 | backend/internal/handler/http/theme_public_cache.go | exact | yes |
| backend/internal/handler/http/public_homepage_handler.go | 117 | backend/internal/handler/http/public_homepage_handler.go | exact | yes |
| frontend/apps/restaurant-website/app/[locale]/page.tsx | 75 | frontend/apps/restaurant-website/app/[locale]/page.tsx | exact | yes |
| public_homepage_handler.go | 117-144 | backend/internal/handler/http/public_homepage_handler.go | prefix | yes |

