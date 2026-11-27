# F-053 — Implementation phases and tasks (all NOT STARTED)

Finding: Manual attendance creation silently drops non-RFC3339 clock times and accepts the record with zero hours, so overtime and absence metrics miss the shift · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-053-01
- **Finding:** F-053
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-053 at code level
- **Existing source paths:** backend/internal/handler/http/hr_attendance_handler.go, backend/internal/handler/http/hr_attendance_handler.go, backend/internal/handler/http/hr_attendance_handler.go
- **Proposed paths:** backend/migrations/<next>_attendance_clock_check.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph10_hr_matrix.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-053-02
- **Finding:** F-053
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-053 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/internal/handler/http/hr_attendance_handler.go, backend/internal/handler/http/hr_attendance_handler.go, backend/internal/handler/http/hr_attendance_handler.go
- **Proposed paths:** backend/migrations/<next>_attendance_clock_check.sql (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-053-specific cases. Design: Proposed: strict clock-time handling in backend/internal/handler/http/hr_attendance_handler.go (CreateAttendance, parse block at :235-243): accept RFC3339 and HH:MM (combined with attendance_date in the tenant time zone), return 400 with a field error when a supplied value cannot be parsed (never ignore the error), require clock_out > clock_in, compute and return total hours and overtime server-side; add a CHECK clock_out > clock_in on the table.
- **Dependencies:** F-062 P3 task; REM-F-053-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** POST with clock_in "09:00" and clock_out "19:30" stores both times and 10.5 hours. | POST with clock_in "abc" returns 400 and stores nothing. | clock_out before clock_in returns 400.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-053-03
- **Finding:** F-053
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-053
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_attendance_clock_check.sql (new)
- **Required changes:** Rows with NULL times and 0 hours (evidence) are listed for HR correction.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Rows with NULL times and 0 hours (evidence) are listed for HR correction.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-053-04
- **Finding:** F-053
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-053 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_attendance_clock_check.sql (new)
- **Required changes:** Simple and RFC3339 inputs; reversed times; hours equal independent calculation
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-053-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Reject invalid time format; compute 8 regular + 2.5 overtime hours; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-053-05
- **Finding:** F-053
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-053
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_attendance_clock_check.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-053 no longer reproduce (expected: Reject invalid time format; compute 8 regular + 2.5 overtime hours); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
