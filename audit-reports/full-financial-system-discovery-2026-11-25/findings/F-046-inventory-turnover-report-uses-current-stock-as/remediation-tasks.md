# F-046 — Implementation phases and tasks (all NOT STARTED)

Finding: Inventory turnover report uses current stock as opening stock and reads waste from the legacy inventory table, ignoring inventory movements · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-046-01
- **Finding:** F-046
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-046
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

### REM-F-046-02
- **Finding:** F-046
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-046 at code level
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph12_batch_apis2.mjs).
- **Dependencies:** REM-F-046-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-046-03
- **Finding:** F-046
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-046
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: compute turnover from the inventory movement ledger. In backend/internal/handler/http/report_handler.go (:557-645 and :2480-2609) take the opening quantity from inventory_movements at the period start (not current stock or the legacy inventory table) and waste from movements of type waste/adjustment_out; average inventory = (opening value + closing value)/2 at product cost; document the cost basis.
- **Dependencies:** REM-F-046-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** For a fixture with known opening stock, purchases, sales and waste the turnover equals an independent calculation to 2 decimals. | Waste in the report equals the sum of waste movements (see F-034).
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-046-04
- **Finding:** F-046
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-046 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Turnover report with opening stock, sales and waste equals independent calculation
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-046-05
- **Finding:** F-046
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-046
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
