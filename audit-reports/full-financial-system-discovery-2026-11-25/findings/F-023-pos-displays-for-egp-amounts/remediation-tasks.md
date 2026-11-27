# F-023 — Implementation phases and tasks (all NOT STARTED)

Finding: POS displays $ for EGP amounts · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-023-01
- **Finding:** F-023
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-023
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-023-02
- **Finding:** F-023
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-023 at code level
- **Existing source paths:** frontend/apps/dashboard/src/lib/rtl.ts
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph20_evidence_gaps.mjs).
- **Dependencies:** REM-F-023-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-023-03
- **Finding:** F-023
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-023
- **Existing source paths:** frontend/apps/dashboard/src/lib/rtl.ts
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** Proposed: stop defaulting to USD. In frontend/apps/dashboard/src/lib/rtl.ts:81-91 (formatCurrencyWithLocale) require an explicit currency argument (no default) and make every POS/cart/receipt caller pass the tenant currency from one source (a tenant settings hook, today missing: no currency setting exists); until a currency setting exists use the single functional currency decided in POL-11 as a named constant, not an implicit default.
- **Dependencies:** REM-F-023-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** The POS, receipts and orders list show the same currency code/symbol as the storefront and reports for the same tenant. | A lint/test fails when formatCurrencyWithLocale is called without a currency. | Switching the tenant currency setting changes all displays.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-023-04
- **Finding:** F-023
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-023
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** API contract: GET tenant settings gains currency (additive) if not present. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-023-05
- **Finding:** F-023
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-023 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** Currency display test
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-023-06
- **Finding:** F-023
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-023
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useTenantCurrency.ts (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
