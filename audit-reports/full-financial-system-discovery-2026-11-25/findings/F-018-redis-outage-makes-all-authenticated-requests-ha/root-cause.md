# F-018 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-005 source + fault injection): AuthSessionStore.Resolve calls the Redis cache with the request context, which has no deadline (backend/internal/middleware/auth.go:60 -> backend/internal/service/auth_session.go:166-171 -> RedisCache.Get cache.go:248); the Redis client is built from ParseURL defaults with no explicit dial, read, write or pool timeouts (cache.go:232-237); the middleware intends a 503 on a resolver error (auth.go:61-67) but the HTTP server WriteTimeout of 60 s (cmd/api/main.go:2004) closes the connection first. Measured with the isolated Redis stopped: GET /products with a valid cookie failed after 68,179 ms and POST /auth/login after 68,650 ms, both with NO HTTP response (client "fetch failed"); /health answered 200 in 5 ms and the public theme endpoint answered in 45 ms. After Redis restarted the old cookie returned 401 (all sessions lost) and a new login worked. Residual not isolated: the exact split of the 68 s between go-redis retries and pool waits


## Origin and propagation
- Origin (first verified source location): `backend/internal/service/cache.go`:231-246.
- Propagation: affected modules Platform; related findings none recorded; observed consequence: Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 
- Cause versus symptom: the symptom is what users or reports show (Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login bot); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-REL-REDIS, EX-RUN005-F018-REDIS-TIMING; scripts: tests/runtime/ph18_r5_redis_timing.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/service/cache.go | 231-246 | backend/internal/service/cache.go | exact | yes |
| backend/cmd/api/main.go | 2003 | backend/cmd/api/main.go | exact | yes |
| backend/internal/middleware/auth.go | 60 | backend/internal/middleware/auth.go | exact | yes |
| backend/internal/service/auth_session.go | 166 | backend/internal/service/auth_session.go | exact | yes |
| backend/cmd/api/main.go | 2004 | backend/cmd/api/main.go | exact | yes |
| backend/internal/service/auth_session.go | 166-171 | backend/internal/service/auth_session.go | exact | yes |
| cache.go | 248 | backend/internal/middleware/cache.go | basename-ambiguous(2) | yes |
| cache.go | 232-237 | backend/internal/middleware/cache.go | basename-ambiguous(2) | yes |
| auth.go | 61-67 | backend/internal/middleware/auth.go | basename | yes |

