# F-076 — Super Admin Users list links every row to /dashboard/users/{id}, a route that does not exist (404)

- **Severity:** Low  ·  **Category:** Software defect (UI navigation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Super Admin > Users
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-SUPERADMIN-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A user detail page (or no link)

Observed: Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome, Super Admin app (port 13114), Users page: click View on a user row

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
