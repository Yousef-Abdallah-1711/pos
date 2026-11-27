# F-055 — Implementation phases and tasks (all NOT STARTED)

Finding: Super Admin extend-trial accepts an unbounded number of days (100,000 days moves the trial end to the year 2300) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-055-01
- **Finding:** F-055
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-055 at code level
- **Existing source paths:** backend/internal/handler/http/subscription_handler.go, backend/internal/handler/http/subscription_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph13_saas.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-055-02
- **Finding:** F-055
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-055
- **Existing source paths:** backend/internal/handler/http/subscription_handler.go, backend/internal/handler/http/subscription_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: bound and record trial extensions. In backend/internal/handler/http/subscription_handler.go ExtendTenantTrial (:128-138) enforce 1 <= days <= 90 (configurable), require a reason, and write a subscription history row (actor, days, previous end, new end); show the history in Super Admin.
- **Dependencies:** REM-F-055-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** days 100000 returns 400. | A valid extension records actor, reason and dates and appears in the tenant history.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-055-03
- **Finding:** F-055
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-055
- **Existing source paths:** backend/internal/handler/http/subscription_handler.go, backend/internal/handler/http/subscription_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-055-04
- **Finding:** F-055
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-055 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** days above the cap rejected; audit row created
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-055-05
- **Finding:** F-055
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-055
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
