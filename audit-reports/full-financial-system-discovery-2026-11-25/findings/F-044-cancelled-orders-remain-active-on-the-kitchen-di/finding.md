# F-044 — Cancelled orders remain active on the Kitchen Display queue

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** KDS; Orders
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-DESTR-ORDER-CANCEL
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Cancelled order is removed from the active kitchen queue

Observed: Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24).  RUN-005 real-UI reproduction: cancelling unpaid order 5796 through the Orders page (confirm dialog, DELETE 200, status cancelled, stock restored 30 -> 32) left its kds_queue row in status "new".

## Severity rationale
Rated **Medium** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: Operational: food prepared for cancelled orders (COGS without sale). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: KDS; Orders. Operational: food prepared for cancelled orders (COGS without sale)

## Technical impact
CONFIRMED (RUN-004 SQL): no trigger or code path updates kds_queue when an order is cancelled (pg_trigger on orders lists only order_stock_transition/order_stock_items); RUN-004 shows 3 cancelled orders still active in the queue (AUD-R4-KDS-after-read.json)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Operational: food prepared for cancelled orders (COGS without sale)

## Where it was reproduced
Cancel a pending order (DELETE /admin/orders/{id}) and open the KDS board (headed Chrome)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
