# F-071 — Implementation phases and tasks (all NOT STARTED)

Finding: ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches finance, HR and user endpoints), mock payments, open metrics and in-memory sessions are allowed · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-071-01
- **Finding:** F-071
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-071 at code level
- **Existing source paths:** backend/internal/config/runtime_env.go, backend/internal/middleware/permission.go, backend/internal/config/paymob_mode.go, backend/cmd/api/health_routes.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph17_sec013.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-071-02
- **Finding:** F-071
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared role/permission matrix to cover F-071 (core component delivered once by the P3 task of F-020)
- **Existing source paths:** backend/internal/config/runtime_env.go, backend/internal/middleware/permission.go, backend/internal/config/paymob_mode.go, backend/cmd/api/health_routes.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-071-specific cases. Design: Proposed: invert the environment check. Replace config.IsProductionMode() (backend/internal/config/runtime_env.go:20-23) usage with an explicit IsDevelopmentMode() allow-list (development, dev, local, test) so that staging and any unknown value behave like production for: RBAC default and refusing RBAC_ENFORCEMENT=audit (middleware/permission.go:36-52), mock payment refusal (config/paymob_mode.go:29), Redis/session and rate-limit store requirements (cmd/api/main.go:148-160), secure cookies (handler/http/auth_handler.go:223) and the metrics token (cmd/api/health_routes.go:52); fail startup on a
- **Dependencies:** F-020 P3 task; REM-F-071-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Starting with ENV=staging and no RBAC variable enforces RBAC: a cashier gets 403 on finance, HR and user endpoints (re-run tests/runtime/ph17_sec013.mjs). | ENV=staging with PAYMOB_MODE unset or mock fails to start. | /health/metrics returns 401 without a token outside development. | An unrecognised ENV value stops the API with a clear message.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-071-03
- **Finding:** F-071
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-071
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: /health/metrics requires a token outside development. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-071-04
- **Finding:** F-071
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-071
- **Existing source paths:** backend/internal/config/runtime_env.go, backend/internal/middleware/permission.go, backend/internal/config/paymob_mode.go, backend/cmd/api/health_routes.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-071-05
- **Finding:** F-071
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-071
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Review staging data for exposure while RBAC was in audit mode.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-071-06
- **Finding:** F-071
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-071 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Start with ENV=staging and no RBAC variable: cashier requests to finance endpoints are 403; mock payment mode is rejected; metrics require a token
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-071-07
- **Finding:** F-071
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-071
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
