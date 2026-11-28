# F-018 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Redis outage makes all authenticated requests hang while /health reports OK.

## 2. Current behavior (verified)
Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 11 ms; after Redis returns the previous session is rejected (401) (REPRODUCED). 

## 3. Required behavior
Fast 503 or degraded mode; health reflects dependency

## 4. Reproduction conditions
Stop Redis then call login/session/orders

## 5. Verified evidence
- runtime/21 | evidence/network/RUN-003/PH-18/AUD-REL-REDIS.json | tests/runtime/ph18_redis_outage.mjs
- evidence/network/RUN-005/PH-18/AUD-R5-F018-REDIS-TIMING.json
- tests/runtime/ph18_r5_redis_timing.mjs

## 6. Root cause and confidence
CONFIRMED — AuthSessionStore.Resolve calls the Redis cache with the request context, which has no deadline (backend/internal/middleware/auth.go:60 -> backend/internal/service/auth_session.go:166-171 -> RedisCache.Get cache.go:248); the Redis client is built from ParseURL defaults with no explicit dial, read, write or pool timeouts (cache.go:232-237); the middleware intends a 503 on a resolver error (auth.go:61-67) but the HTTP server WriteTimeout of 60 s (cmd/api/main.go:2004) closes the connection first. Measured with the isolated Redis stopped: GET /products with a valid cookie failed after 68,179 ms and POST /auth/login after 68,650 ms, both with NO HTTP response (client "fetch failed"); /health answ

## 7. Affected paths
- Backend/database: `backend/internal/service/cache.go`:231-246, `backend/cmd/api/main.go`:2003, `backend/internal/middleware/auth.go`:60, `backend/internal/service/auth_session.go`:166, `backend/cmd/api/main.go`:2004, `backend/internal/service/auth_session.go`:166-171, `backend/internal/middleware/cache.go`:248, `backend/internal/middleware/cache.go`:232-237, `backend/internal/middleware/auth.go`:61-67
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: /health response gains a redis component (additive)
- Accounting logic / posting rules: none named
- Reports affected: Platform

## 8. Impact
- Financial: n/a
- Operational: Platform workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause CONFIRMED in RUN-005: middleware/auth.go:60 -> auth_session.go:166 -> cache.go:248 with no deadline, default go-redis options, server WriteTimeout 60 s): 0) add a Redis-down fault-injection test that records per-call durations (ph18_r5_redis_timing.mjs is the reference); 1) set explicit DialTimeout, ReadTimeout, WriteTimeout, PoolTimeout and MaxRetries on the go-redis options built at cache.go:232-237 (target < 500 ms per call); 2) wrap session and rate-limit calls in a context with a 1 s deadline and fail closed with 503 and a Retry-After header instead of hanging; 3) make /health report degraded when Redis is unreachable (add a bounded PING) so orchestration can act; 4) public read endpoints (menu) keep working without Redis.

## 10. Alternative designs
A: timeouts + 503 + health probe (recommended). B: fall back to in-memory sessions on Redis failure in production (rejected: breaks session revocation guarantees).

## 11. Required policy decisions
None; availability target is an operations decision.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-020, F-025, F-030, F-031, F-033, F-036, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go, backend/internal/middleware/auth.go, backend/internal/middleware/cache.go, backend/internal/service/auth_session.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Dependency-failure tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With Redis stopped, an authenticated request returns 503 within 2 s (not a 30 s hang).
- /health returns a degraded status within 1 s while Redis is down and recovers after restart.
- After Redis returns, requests succeed without restarting the API.
