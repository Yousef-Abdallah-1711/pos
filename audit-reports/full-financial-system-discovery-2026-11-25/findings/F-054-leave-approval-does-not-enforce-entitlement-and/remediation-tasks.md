# F-054 — Implementation phases and tasks (all NOT STARTED)

Finding: Leave approval does not enforce entitlement and total_days is not validated against the requested dates · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-054-01
- **Finding:** F-054
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-054
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-054-02
- **Finding:** F-054
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-054 at code level
- **Existing source paths:** backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph10_hr_matrix.mjs).
- **Dependencies:** REM-F-054-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-054-03
- **Finding:** F-054
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-054
- **Existing source paths:** backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go, backend/internal/handler/http/hr_leave_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: derive and enforce leave entitlement. In backend/internal/handler/http/hr_leave_handler.go (:207-210 create, :328-329 update, approve near :412) ignore client total_days and compute working days from start/end (excluding configured weekends/holidays); at approval compare with the remaining annual allowance and either reject, or approve the excess as unpaid leave according to POL-14, recording the unpaid portion; the leave balances report reads the same ledger (it currently clamps remaining to 0 when used exceeds the allowance, hiding the overdraw).
- **Dependencies:** REM-F-054-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Requesting 3 calendar days stores total_days 3 regardless of the client value. | Approving 30 days against a remaining allowance of 21 follows POL-14 (reject, or 21 paid + 9 unpaid recorded). | The balances report shows the true remaining balance including negative values or the unpaid portion.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-054-04
- **Finding:** F-054
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-054
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Existing over-allowance approvals (30 vs 21 in the evidence) are listed for HR.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-054-05
- **Finding:** F-054
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-054 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Over-allowance and mismatched total_days rejected or converted to unpaid
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-054-06
- **Finding:** F-054
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-054
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
