# F-035 — Implementation phases and tasks (all NOT STARTED)

Finding: Standards deviations: login limit 20/min (standard 5/min); CSP allows unsafe-eval and unsafe-inline · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-035-01
- **Finding:** F-035
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-035 at code level
- **Existing source paths:** backend/internal/middleware/cache.go, backend/internal/middleware/rate_limit.go, backend/internal/middleware/rate_limit.go
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

### REM-F-035-02
- **Finding:** F-035
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-035
- **Existing source paths:** backend/internal/middleware/cache.go, backend/internal/middleware/rate_limit.go, backend/internal/middleware/rate_limit.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: tighten the Content-Security-Policy built at backend/internal/middleware/cache.go:312-314 (remove unsafe-eval; replace unsafe-inline with nonces or hashes generated per response for scripts and styles that need them) and set the login rate limit in middleware/rate_limit.go to the documented standard (5 per minute per IP+account) with a lockout response; keep the Redis-backed store (see F-018).
- **Dependencies:** REM-F-035-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Response CSP contains no unsafe-eval and no unsafe-inline for scripts; the dashboard, landing and storefront still render in Chrome without CSP violations. | The sixth failed login in a minute from one IP returns 429.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-035-03
- **Finding:** F-035
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-035 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Header and limit tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-035-04
- **Finding:** F-035
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-035
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
