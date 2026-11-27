# F-019 — Application DB role sees all tenants on tables without RLS and can alter paid orders/payment amounts

- **Severity:** Medium  ·  **Category:** Security / defense in depth  ·  **Evidence status:** REPRODUCED (runtime catalog)
- **Affected modules:** Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Cross-tenant rows hidden; paid records immutable

Observed: All 4 tenants' orders/journals/GL/COA/users visible; UPDATE orders/payments_received succeeded (rolled back); ledger tables protected by grants

## Monetary / operational impact
n/a

## Where it was reproduced
Connect as pos_app without tenant context

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
