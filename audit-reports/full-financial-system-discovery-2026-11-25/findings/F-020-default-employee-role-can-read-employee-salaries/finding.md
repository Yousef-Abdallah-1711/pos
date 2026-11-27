# F-020 — Default Employee role can read employee salaries but cannot use the POS; no built-in cashier role

- **Severity:** Medium  ·  **Category:** Security/usability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** RBAC; HR; POS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-SEC-009
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Least privilege

Observed: GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact data; accounting/report/admin routes 403/401 (REPRODUCED). 

## Monetary / operational impact
n/a

## Where it was reproduced
Login as Employee-role user

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
