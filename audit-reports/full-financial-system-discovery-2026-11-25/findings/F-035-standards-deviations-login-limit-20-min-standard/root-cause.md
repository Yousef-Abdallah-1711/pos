# F-035 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/middleware/cache.go:312-314 builds a CSP with script-src unsafe-inline and unsafe-eval; login limit is configured in middleware/rate_limit.go (observed 20/min, 14 failed attempts then 429 in SEC-012c)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/middleware/cache.go | 312-314 | backend/internal/middleware/cache.go | exact | yes |
| backend/internal/middleware/rate_limit.go | 32 | backend/internal/middleware/rate_limit.go | exact | yes |
| middleware/rate_limit.go | - | backend/internal/middleware/rate_limit.go | prefix | yes |

