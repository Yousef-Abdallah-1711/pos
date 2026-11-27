# F-033 — Employee role can export all orders and read billing/subscription data and HR records

- **Severity:** Medium  ·  **Category:** Security (privilege)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders; Subscription; HR
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-SEC-009
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Least privilege

Observed: GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200 RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR employee records with salary fields; billing/subscription route GET /payments/subscription returned 404 for all roles in this build, so the billing part of the original claim is not re-verified.

## Severity rationale
Rated **Medium** (category: Security (privilege); evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Orders; Subscription; HR. n/a

## Technical impact
CONFIRMED (source + RUN-003): same seeding/guard pattern as F-020 (backend/internal/repository/auth_repo.go:380 default Employee permissions; backend/cmd/api/main.go route wrappers) lets the Employee role read HR and export orders

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
API sweep as Employee user

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
