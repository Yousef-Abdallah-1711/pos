# F-064 — Implementation phases and tasks (all NOT STARTED)

Finding: Tables API returns 500 and a generic message for an over-long name (the duplicate-table-number claim was RETRACTED in RUN-004) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-064-01
- **Finding:** F-064
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-064 at code level
- **Existing source paths:** backend/internal/repository/table_repo.go, backend/internal/handler/http/table_handler.go, backend/internal/repository/table_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph16_tables.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-064-02
- **Finding:** F-064
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared error-mapping layer to cover F-064 (core component delivered once by the P3 task of F-025)
- **Existing source paths:** backend/internal/repository/table_repo.go, backend/internal/handler/http/table_handler.go, backend/internal/repository/table_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-064-specific cases. Design: Proposed (scope reduced in RUN-004: uniqueness is already enforced by uq_restaurant_tables_number and mapped to 409): validate the table name length (<= 100) and location in backend/internal/handler/http/table_handler.go before the insert and map pq 22001 (value too long) to 400 via the shared error mapper (F-025).
- **Dependencies:** F-025 P3 task; REM-F-064-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** A 1,800-character name returns 400 with a field error. | Duplicate table numbers still return 409.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-064-03
- **Finding:** F-064
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-064 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Duplicate number returns 409; over-long name returns 400
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-064-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 409 duplicate table number; 400 name length; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-064-04
- **Finding:** F-064
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-064
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-064 no longer reproduce (expected: 409 duplicate table number; 400 name length); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
