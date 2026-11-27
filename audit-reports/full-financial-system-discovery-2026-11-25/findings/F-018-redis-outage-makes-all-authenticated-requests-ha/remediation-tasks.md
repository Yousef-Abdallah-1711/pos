# F-018 — Implementation phases and tasks (all NOT STARTED)

Finding: Redis outage makes all authenticated requests hang while /health reports OK · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-018-01
- **Finding:** F-018
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-018 at code level
- **Existing source paths:** backend/internal/service/cache.go, backend/cmd/api/main.go, backend/internal/middleware/auth.go, backend/internal/service/auth_session.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph18_r5_redis_timing.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-018-02
- **Finding:** F-018
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-018
- **Existing source paths:** backend/internal/service/cache.go, backend/cmd/api/main.go, backend/internal/middleware/auth.go, backend/internal/service/auth_session.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed (root cause CONFIRMED in RUN-005: middleware/auth.go:60 -> auth_session.go:166 -> cache.go:248 with no deadline, default go-redis options, server WriteTimeout 60 s): 0) add a Redis-down fault-injection test that records per-call durations (ph18_r5_redis_timing.mjs is the reference); 1) set explicit DialTimeout, ReadTimeout, WriteTimeout, PoolTimeout and MaxRetries on the go-redis options built at cache.go:232-237 (target < 500 ms per call); 2) wrap session and rate-limit calls in a context with a 1 s deadline and fail closed with 503 and a Retry-After header instead of hanging; 3) make /health report degraded when Redis is unreachable (add a bounded PING) so orchestration can act; 4) public read endpoints (menu) keep working without Redis.
- **Dependencies:** REM-F-018-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With Redis stopped, an authenticated request returns 503 within 2 s (not a 30 s hang). | /health returns a degraded status within 1 s while Redis is down and recovers after restart. | After Redis returns, requests succeed without restarting the API.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-018-03
- **Finding:** F-018
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-018
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/reports
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: /health response gains a redis component (additive). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** With Redis stopped, an authenticated request returns 503 within 2 s (not a 30 s hang).
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-018-04
- **Finding:** F-018
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-018 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Dependency-failure tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-018-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Fast 503 or degraded mode; health reflects dependency; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-018-05
- **Finding:** F-018
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-018
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-018 no longer reproduce (expected: Fast 503 or degraded mode; health reflects dependency); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
