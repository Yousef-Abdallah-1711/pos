# F-027 — Employee-role users get finance, HR and RBAC page shells in the UI (no UI guard) while the API returns 403

- **Severity:** Low  ·  **Category:** Security/usability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Frontend routing; RBAC
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Page guard consistent with API permissions

Observed: Shells render (P&L, BS, CF, VAT, aging, payroll, RBAC roles) with blank data; only 24 pages gated

## Monetary / operational impact
n/a

## Where it was reproduced
Crawl dashboard as Employee user: 50 of 79 routes render; 24 show Access denied; finance report pages render with 403 data calls

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
