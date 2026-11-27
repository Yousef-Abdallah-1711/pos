# F-012 — Implementation phases and tasks (all NOT STARTED)

Finding: At least six revenue/profit definitions disagree and are unlabeled · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-012-01
- **Finding:** F-012
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-012
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-012-02
- **Finding:** F-012
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-012 at code level
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-012-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-012-03
- **Finding:** F-012
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-012
- **Existing source paths:** backend/internal/handler/http/report_handler.go, backend/internal/handler/http/report_handler.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Proposed: publish one definition per concept and label every report with its basis. (1) Adopt the decisions in POL-01 (revenue recognition basis) and POL-12 (single source of truth for dashboards/KPIs): "Sales" (operational, tax-inclusive or net as approved), "Profitability revenue" (net of tax, paid or all as approved) and "GL revenue" (accrual); (2) add a basis label and tooltip to each tile and CSV/PDF header (frontend reports pages and report_handler.go responses gain a "basis" field additively); (3) move the three queries (report_handler.go:2177, :2367 and financial_report_usecase.go:208) behind one shared revenue-definition package so a basis change happens in one place; (4) provide a reconciliation panel showing the bridge Sales -> Profitability -> GL (tax, unpaid orders, cancelled orders).
- **Dependencies:** REM-F-012-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Every revenue figure in the UI, CSV and PDF shows its basis label. | For the RUN-003 dataset the bridge reproduces 517,000 -> 470,000 -> 290,000 with tax 47,000 and unpaid 180,000 explained line by line. | A contract test fails if two reports that claim the same basis return different totals for the same dataset.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-012-04
- **Finding:** F-012
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-012
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/reports
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** API contract: report responses gain additive "basis" fields; no field renamed. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-012-05
- **Finding:** F-012
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-012
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-012-06
- **Finding:** F-012
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-012 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Cross-report reconciliation suite
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-012-07
- **Finding:** F-012
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-012
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/usecase/revenue_definitions.go (new); frontend/apps/dashboard/src/components/reports/BasisLabel.tsx (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
