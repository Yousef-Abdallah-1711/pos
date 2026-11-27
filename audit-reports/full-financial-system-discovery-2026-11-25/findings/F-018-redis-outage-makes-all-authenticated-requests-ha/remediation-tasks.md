# F-018 — Implementation phases and tasks (all NOT STARTED)

Finding: Redis outage makes all authenticated requests hang while /health reports OK · Severity Medium · Root-cause confidence PROBABLE

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-018-01
- **Finding:** F-018
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-018 at code level
- **Existing source paths:** backend/internal/service/cache.go, backend/cmd/api/main.go, backend/internal/middleware/cache.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
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
- **Existing source paths:** backend/internal/service/cache.go, backend/cmd/api/main.go, backend/internal/middleware/cache.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed (root cause still PROBABLE, so step 0 is investigation): 0) trace the session lookup call chain from the auth middleware to backend/internal/service/cache.go and add a Redis-down fault-injection test that records per-call durations; 1) set explicit DialTimeout, ReadTimeout, WriteTimeout, PoolTimeout and MaxRetries on the go-redis options built at cache.go:232-237 (target < 500 ms per call); 2) wrap session and rate-limit calls in a context with a 1 s deadline and fail closed with 503 and a Retry-After header instead of hanging; 3) make /health report degraded when Redis is unreachable (add a bounded PING) so orchestration can act; 4) public read endpoints (menu) keep working without Redis.
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
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
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
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
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
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
