# F-046 — Inventory turnover report uses current stock as opening stock and reads waste from the legacy inventory table, ignoring inventory movements

- **Severity:** Low  ·  **Category:** Software defect (report)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Inventory turnover report; Overview report
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Opening quantity from stock history; waste_qty reflects recorded waste

Observed: opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has waste 1 for products 2 and 4; Audit Burger original stock 50 | RUN-003: turnover opening_stock_qty == current stock (33/21); waste_qty 0 vs movements waste 1 (REPRODUCED). 

## Severity rationale
Rated **Low** (category: Software defect (report); evidence: REPRODUCED). Basis recorded in the register: Turnover ratio misstated (1.42 vs 1.01 using true average 41.5 qty). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Inventory turnover report; Overview report. Turnover ratio misstated (1.42 vs 1.01 using true average 41.5 qty)

## Technical impact
CONFIRMED (source per inventory 30 rows TURN-001/008): report_handler.go:557-645 / 2480-2609 reads legacy inventory.quantity_after else current stock and SUM(ABS(change)) from the legacy inventory table, while the live ledger is inventory_movements

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Turnover ratio misstated (1.42 vs 1.01 using true average 41.5 qty)

## Where it was reproduced
GET /reports/inventory after recording a waste adjustment and sales (tenant 4)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
