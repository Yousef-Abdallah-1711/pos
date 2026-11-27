# F-033 — Employee role can export all orders and read billing/subscription data and HR records

- **Severity:** Medium  ·  **Category:** Security (privilege)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders; Subscription; HR
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-SEC-009
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Least privilege

Observed: GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200

## Monetary / operational impact
n/a

## Where it was reproduced
API sweep as Employee user

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
