# F-044 — Cancelled orders remain active on the Kitchen Display queue

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** KDS; Orders
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Cancelled order is removed from the active kitchen queue

Observed: Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24). 

## Monetary / operational impact
Operational: food prepared for cancelled orders (COGS without sale)

## Where it was reproduced
Cancel a pending order (DELETE /admin/orders/{id}) and open the KDS board (headed Chrome)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
