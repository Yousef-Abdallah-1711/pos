# F-067 — Implementation phases and tasks (all NOT STARTED)

Finding: Profitability "Lowest Margin Items" is the profit-amount ranking read backwards, so a 100%-margin item ranks second-lowest · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-067-01
- **Finding:** F-067
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-067 at code level
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go
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

### REM-F-067-02
- **Finding:** F-067
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-067
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: sort the lowest-margin list by margin. In backend/internal/handler/http/report_handler.go replace the backwards iteration at :2470-2476 with a query (or an explicit sort) by gross_margin_pct ASC then revenue DESC, keeping the top list ordered by gross profit (:2432); label both lists with their sort key in the UI and exports.
- **Dependencies:** REM-F-067-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** For items with margins 0%, 60%, 60%, 60%, 100% the lowest list is ordered 0%, 60%, 60%, 60%, 100%. | The top list stays ordered by gross profit. | CSV and PDF exports match the screen.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-067-03
- **Finding:** F-067
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-067
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/reports
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: response shape unchanged; ordering of lowest_margin_items changes. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-067-04
- **Finding:** F-067
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-067 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Fixture with items of different margin and profit; assert both lists
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-067-05
- **Finding:** F-067
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-067
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
