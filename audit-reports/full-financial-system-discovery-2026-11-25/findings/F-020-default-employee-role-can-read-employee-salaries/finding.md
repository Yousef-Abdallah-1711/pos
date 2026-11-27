# F-020 — Default Employee role can read employee salaries but cannot use the POS; no built-in cashier role

- **Severity:** Medium  ·  **Category:** Security/usability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** RBAC; HR; POS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-SEC-009, EX-RUN005-RBAC-MATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Least privilege

Observed: GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact data; accounting/report/admin routes 403/401 (REPRODUCED).  RUN-005: in enforce mode the EMPLOYEE (cashier) role gets 200 from GET /hr/employees with base_salary (5000) and email in the payload, and 200 from /hr/attendance and /hr/leaves; MANAGER additionally reads /hr/salaries (4 rows, salary fields) and the user and role lists (policy question POL-15, not asserted as a defect); DRIVER sees 0 orders (scoped) and 3/28 readable endpoints.

## Severity rationale
Rated **Medium** (category: Security/usability; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: RBAC; HR; POS. n/a

## Technical impact
CONFIRMED (source + DB): default Employee role is seeded with HR=READ, SALARIES=NONE (auth_repo.go seedDefaultRolesAndPermissions :380), but GET /hr/employees is guarded only by module 3 READ (cmd/api/main.go:877) and serializes base_salary (domain/hr.go:44), so a cashier reads salaries (RUN-003 SEC-009)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Login as Employee-role user

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
