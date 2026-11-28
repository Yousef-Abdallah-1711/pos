# F-018 — Redis outage makes all authenticated requests hang while /health reports OK

- **Severity:** Medium  ·  **Category:** Reliability defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-REL-REDIS, EX-RUN005-F018-REDIS-TIMING
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Fast 503 or degraded mode; health reflects dependency

Observed: Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 11 ms; after Redis returns the previous session is rejected (401) (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Reliability defect; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Platform. n/a

## Technical impact
CONFIRMED (RUN-005 source + fault injection): AuthSessionStore.Resolve calls the Redis cache with the request context, which has no deadline (backend/internal/middleware/auth.go:60 -> backend/internal/service/auth_session.go:166-171 -> RedisCache.Get cache.go:248); the Redis client is built from ParseURL defaults with no explicit dial, read, write or pool timeouts (cache.go:232-237); the middleware intends a 503 on a

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Stop Redis then call login/session/orders

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
