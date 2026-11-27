# F-004 — Implementation phases and tasks (all NOT STARTED)

Finding: POS shows 14% tax but server stores 10% · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-004-01
- **Finding:** F-004
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-004
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-004-02
- **Finding:** F-004
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-004 at code level
- **Existing source paths:** frontend/apps/dashboard/src/stores/posStore.ts, backend/internal/usecase/order_pricing.go, backend/internal/usecase/order_pricing.go
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-004-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-004-03
- **Finding:** F-004
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared discount-first tax calculation to cover F-004 (core component delivered once by the P3 task of F-007)
- **Existing source paths:** frontend/apps/dashboard/src/stores/posStore.ts, backend/internal/usecase/order_pricing.go, backend/internal/usecase/order_pricing.go
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-004-specific cases. Design: Proposed: make the server the only tax authority. (1) Add a read endpoint (proposed GET /api/v1/pos/pricing-config) returning the effective tax rate and currency from tenant settings; (2) change frontend/apps/dashboard/src/stores/posStore.ts:107 to read that value instead of the hard-coded TAX_RATE = 0.14 and label it with the returned rate; (3) change backend/internal/usecase/order_pricing.go:46 to compute tax from the same tenant setting (default 14% only if policy POL-02 approves) instead of the fixed (subtotal+5)/10 expression; (4) keep rejecting client-supplied tax (already ignored) and r
- **Dependencies:** F-007 P3 task; REM-F-004-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** For subtotal 60,000 with the approved rate r, the POS cart tax, the order response tax_amount and orders.tax_amount are all equal to ROUND(60,000 x r, 2). | Changing the tenant tax rate changes the POS preview and the stored tax without a code change. | A mismatch between the UI preview and the stored total is impossible: contract test posts the cart and compares totals.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-004-04
- **Finding:** F-004
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-004
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-004-05
- **Finding:** F-004
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-004
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** Existing orders keep the stored 10% tax. Produce a read-only exception report of orders whose stored rate differs from the approved rate; any correction is a reversing/adjusting journal approved by the accountant.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-004-06
- **Finding:** F-004
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-004 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** UI vs server total equality tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-004-07
- **Finding:** F-004
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-004
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/pos_pricing_handler.go (new, preview endpoint); frontend/apps/dashboard/src/lib/posPricing.ts (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
