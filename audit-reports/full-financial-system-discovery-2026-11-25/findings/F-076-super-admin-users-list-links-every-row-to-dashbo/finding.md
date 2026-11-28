# F-076 — Super Admin Users list links every row to /dashboard/users/{id}, a route that does not exist (404)

- **Severity:** Low  ·  **Category:** Software defect (UI navigation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Super Admin > Users
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-SUPERADMIN-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A user detail page (or no link)

Observed: Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route

## Severity rationale
Rated **Low** (category: Software defect (UI navigation); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Super Admin > Users. 0

## Technical impact
CONFIRMED (source): super-admin/src/app/dashboard/users/page.tsx:132 renders href /dashboard/users/${user.id}; the directory app/dashboard/users contains only page.tsx (no [id] route)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome, Super Admin app (port 13114), Users page: click View on a user row

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
