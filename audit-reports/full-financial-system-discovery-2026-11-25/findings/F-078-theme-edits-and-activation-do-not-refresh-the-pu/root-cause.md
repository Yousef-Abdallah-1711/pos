# F-078 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): invalidatePublicThemeCacheByRestaurant is defined (handler/http/theme_public_cache.go:45) but has no caller; GetThemeOnly caches the payload for publicThemeCacheTTL = 5 min (public_homepage_handler.go:117-144)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/theme_public_cache.go | 45 | backend/internal/handler/http/theme_public_cache.go | exact | yes |
| backend/internal/handler/http/public_homepage_handler.go | 117 | backend/internal/handler/http/public_homepage_handler.go | exact | yes |
| frontend/apps/restaurant-website/app/[locale]/page.tsx | 75 | frontend/apps/restaurant-website/app/[locale]/page.tsx | exact | yes |
| public_homepage_handler.go | 117-144 | backend/internal/handler/http/public_homepage_handler.go | prefix | yes |

