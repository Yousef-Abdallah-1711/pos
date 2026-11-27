# F-048 — Implementation phases and tasks (all NOT STARTED)

Finding: Exchange-rate API accepts a same-currency pair (USD/USD) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-048-01
- **Finding:** F-048
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-048
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

### REM-F-048-02
- **Finding:** F-048
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-048 at code level
- **Existing source paths:** backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph11_accounting_flow2.mjs).
- **Dependencies:** REM-F-048-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-048-03
- **Finding:** F-048
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-048 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/handler/http/accounting_advanced_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-048-specific cases. Design: Proposed: validate currency pairs. In backend/internal/handler/http/accounting_advanced_handler.go (:475-482) require ISO-4217 uppercase codes from an allow-list, base != quote, and a rate within a sane range (> 0 and < configurable maximum); enforce uniqueness of (tenant, base, quote, effective_date) in the repository (accounting_advanced_repo.go:596).
- **Dependencies:** F-062 P3 task; REM-F-048-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** USD/USD returns 400. | An unknown code returns 400. | Two rates for the same pair and date return 409.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-048-04
- **Finding:** F-048
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-048 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Same-currency and malformed code rejected
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-048-05
- **Finding:** F-048
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-048
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
