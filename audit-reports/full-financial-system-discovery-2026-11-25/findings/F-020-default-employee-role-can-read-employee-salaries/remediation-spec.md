# F-020 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Default Employee role can read employee salaries but cannot use the POS; no built-in cashier role.

## 2. Current behavior (verified)
GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact data; accounting/report/admin routes 403/401 (REPRODUCED).  RUN-005: in enforce mode the EMPLOYEE (cashier) role gets 200 from GET /hr/employees with base_salary (5000) and email in the payload, and 200 from /hr/attendance and /hr/leaves; MANAGER additionally reads /hr/salaries (4 rows, salary fields) and the user and role lists (policy question POL-15, not asserted as a defect); DRIVER sees 0 orders (scoped) and 3/28 readable endpoints.

## 3. Required behavior
Least privilege

## 4. Reproduction conditions
Login as Employee-role user

## 5. Verified evidence
- runtime/09 | evidence/network/RUN-003/PH-17/AUD-SEC-BATTERY.json | tests/runtime/ph17_security.mjs
- evidence/network/RUN-005/PH-17/AUD-R5-RBAC-MATRIX-ENFORCE.json

## 6. Root cause and confidence
CONFIRMED — default Employee role is seeded with HR=READ, SALARIES=NONE (auth_repo.go seedDefaultRolesAndPermissions :380), but GET /hr/employees is guarded only by module 3 READ (cmd/api/main.go:877) and serializes base_salary (domain/hr.go:44), so a cashier reads salaries (RUN-003 SEC-009)

## 7. Affected paths
- Backend/database: `backend/internal/repository/auth_repo.go`:380, `backend/cmd/api/main.go`:877, `backend/internal/domain/hr.go`:44, `backend/internal/repository/auth_repo.go`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: GET /api/v1/hr/employees: pay fields omitted unless SALARIES READ (breaking for clients that relied on them; coordinate with dashboard)
- Accounting logic / posting rules: none named
- Reports affected: RBAC; HR; POS

## 8. Impact
- Financial: n/a
- Operational: RBAC; HR; POS workflows behave as in section 2.
- Security: Yes — see category: Security/usability
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: define built-in roles and stop leaking salary fields. (1) Seed a Cashier role (POS write, orders read/write, no HR, no accounting) next to the existing Employee role in auth_repo.go seedDefaultRolesAndPermissions (:380); (2) serialise base_salary and other pay fields in GET /hr/employees only when the caller has SALARIES READ (domain/hr.go:44 field becomes conditional via a response DTO); (3) change the UI guards so the POS entry follows the Cashier role.

## 10. Alternative designs
A: new Cashier role + DTO that omits pay fields (recommended). B: only hide the pay fields (leaves the Employee role unable to use POS). C: grant Employee POS access (widens an already broad role).

## 11. Required policy decisions
POL-15 (default role matrix) in 56.

## 12. Required data migrations
Seed/upsert of the Cashier role for existing tenants (idempotent data migration).

## 13. Historical data implications
None; existing exposure should be assessed (access logs).

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Core of the fix: role/permission definitions and negative tests.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-001, F-018, F-025, F-030, F-031, F-033, F-036, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go, backend/internal/domain/hr.go, backend/internal/repository/auth_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
RBAC matrix tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A user with the Employee role receives no base_salary or bank fields from GET /hr/employees.
- A Cashier can create and pay POS orders and gets 403 on finance and HR endpoints.
- RBAC matrix test: module x level for Owner, Manager, Cashier, Employee matches the approved matrix.
