# F-018 — Root cause

**Confidence:** PROBABLE

## First break (from the register)
PROBABLE: backend/internal/service/cache.go:231-246 builds the Redis client from ParseURL defaults and the session path calls it without a short per-request deadline; with Redis stopped authenticated requests and login hang beyond 30 s (AUD-REL-REDIS) while /health stays 200 because it does not probe Redis; exact middleware call site not yet identified [RUN-004: kept PROBABLE. Source shows default go-redis options (cache.go:232-237, ping timeout 2 s at startup only); the 30 s hang implies repeated per-request calls but the middleware call site is not yet identified. Required next step: trace session lookup call chain and set a Redis-down fault injection with per-call timing.]


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/service/cache.go | 231-246 | backend/internal/service/cache.go | exact | yes |
| backend/cmd/api/main.go | 2003 | backend/cmd/api/main.go | exact | yes |
| cache.go | 232-237 | backend/internal/middleware/cache.go | basename-ambiguous(2) | yes |


## Missing evidence and required next investigation
- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.
