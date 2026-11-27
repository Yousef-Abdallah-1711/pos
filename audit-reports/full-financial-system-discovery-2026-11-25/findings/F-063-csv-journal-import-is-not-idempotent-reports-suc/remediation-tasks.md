# F-063 — Implementation phases and tasks (all NOT STARTED)

Finding: CSV journal import is not idempotent, reports success when every entry failed, and long imports outlive the HTTP connection leaving the outcome unknown · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-063-01
- **Finding:** F-063
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-063 at code level
- **Existing source paths:** backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph11_import.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-063-02
- **Finding:** F-063
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-063
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-063-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-063-03
- **Finding:** F-063
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-063
- **Existing source paths:** backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Proposed: make journal CSV import safe. In backend/internal/handler/http/accounting_advanced_handler.go ImportJournalEntriesCSV (:1067-1140) (1) require a client batch id (or hash of the file) with a unique constraint so a re-upload is rejected as a duplicate; (2) validate the whole file first and return 422 (or 207 with per-row errors) when any entry fails instead of 200 success:true; (3) process large files as a background job returning 202 plus a status endpoint, because imports of 1,500 journals take about 60 s and 6,000 outlive the HTTP connection; (4) wrap each file in one transaction when all-or-nothing is chosen.
- **Dependencies:** REM-F-063-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Re-uploading the same file creates no duplicate journals. | A file with one bad entry returns a per-row error report and the documented all-or-nothing behaviour. | A 27,000-journal file returns 202 within 2 s and completes in the background with a progress endpoint.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-063-04
- **Finding:** F-063
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-063
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Duplicate journals created by repeated imports are listed.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-063-05
- **Finding:** F-063
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-063 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Re-upload yields the same batch result; all-failed import returns an error status; 6,000-row import returns a job id within 5 s
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-063-06
- **Finding:** F-063
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-063
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/jobs/journal_import_job.go (new); backend/migrations/<next>_journal_import_batches.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
