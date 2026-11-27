# F-033 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Employee role can export all orders and read billing/subscription data and HR records.

## 2. Current behavior (verified)
GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200 RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR employee records with salary fields; billing/subscription route GET /payments/subscription returned 404 for all roles in this build, so the billing part of the original claim is not re-verified.

## 3. Required behavior
Least privilege

## 4. Reproduction conditions
API sweep as Employee user

## 5. Verified evidence
- evidence/api-sweep-results.csv
- evidence/network/RUN-005/PH-17/AUD-R5-RBAC-MATRIX-ENFORCE.json

## 6. Root cause and confidence
CONFIRMED — same seeding/guard pattern as F-020 (backend/internal/repository/auth_repo.go:380 default Employee permissions; backend/cmd/api/main.go route wrappers) lets the Employee role read HR and export orders

## 7. Affected paths
- Backend/database: `backend/internal/repository/auth_repo.go`:380, `backend/cmd/api/main.go`:877
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`, `frontend/apps/super-admin`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: several GET endpoints start returning 403 for the Employee role (breaking, coordinated with the dashboard)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: Orders; Subscription; HR workflows behave as in section 2.
- Security: Yes — see category: Security (privilege)
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: tighten the Employee role and per-endpoint levels together with F-020. Remove billing/subscription read and bulk order export from the default Employee permission set seeded at backend/internal/repository/auth_repo.go:380; guard the export routes in backend/cmd/api/main.go with the EXPORT level instead of READ; add a role-matrix test covering orders export, subscription/billing reads and HR records for each default role.

## 10. Alternative designs
A: role seed + route levels + matrix test (recommended). B: only hide UI entries (API stays open).

## 11. Required policy decisions
POL-15 (default role matrix) in 56.

## 12. Required data migrations
Data migration updating role_permissions for existing tenants (idempotent).

## 13. Historical data implications
Assess access logs for past Employee exports.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
This finding is the authorization change.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-020 and shared-file findings F-001, F-018, F-020, F-025, F-030, F-031, F-036, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go, backend/internal/repository/auth_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
RBAC matrix tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Employee gets 403 on /orders export, billing and subscription endpoints and on HR records beyond their own.
- Owner/Manager behaviour is unchanged.
- The matrix test in CI fails when a route lacks an explicit level.
