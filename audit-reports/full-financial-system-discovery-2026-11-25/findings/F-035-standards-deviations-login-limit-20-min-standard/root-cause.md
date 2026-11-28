# F-035 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/middleware/cache.go:312-314 builds a CSP with script-src unsafe-inline and unsafe-eval; login limit is configured in middleware/rate_limit.go (observed 20/min, 14 failed attempts then 429 in SEC-012c)


## Origin and propagation
- Origin (first verified source location): `backend/internal/middleware/cache.go`:312-314.
- Propagation: affected modules Auth; Platform; related findings none recorded; observed consequence: 429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contai
- Cause versus symptom: the symptom is what users or reports show (429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20t); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-OBSERVED-REPRO; scripts: tests/runtime/ph17_r5_observed_repro.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/middleware/cache.go | 312-314 | backend/internal/middleware/cache.go | exact | yes |
| backend/internal/middleware/rate_limit.go | 32 | backend/internal/middleware/rate_limit.go | exact | yes |
| middleware/rate_limit.go | - | backend/internal/middleware/rate_limit.go | prefix | yes |

