# F-034 — Waste/shrinkage adjustments change stock but post nothing and the turnover report shows waste 0

- **Severity:** Medium  ·  **Category:** Missing capability / defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Inventory; Reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Expense posted and waste reported

Observed: Stock 48 to 47; no journal; waste_qty 0; reason count_shortage not accepted (allowed: purchase, sale, adjustment, waste, return)

## Severity rationale
Rated **Medium** (category: Missing capability / defect; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Inventory; Reports. n/a

## Technical impact
CONFIRMED (source): backend/internal/handler/http/inventory_adjustment_handler.go (Adjust :29-105) updates stock and writes an inventory movement but contains no journal/posting call (0 references); waste therefore posts nothing, and turnover reads the legacy table (F-046)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
POST inventory adjust reason waste -1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
