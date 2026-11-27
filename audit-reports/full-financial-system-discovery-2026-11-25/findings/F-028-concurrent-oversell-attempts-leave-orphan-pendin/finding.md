# F-028 — Concurrent oversell attempts leave orphan pending orders with items although the API returned 400

- **Severity:** High  ·  **Category:** Software defect (race)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Inventory; KDS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 5 orders accepted; 15 rejected leave no order

Observed: 5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in Sales Analytics

## Severity rationale
Rated **High** (category: Software defect (race); evidence: REPRODUCED). Basis recorded in the register: 15 x 11000 = 165000 gross phantom pending sales. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: POS; Orders; Inventory; KDS. 15 x 11000 = 165000 gross phantom pending sales

## Technical impact
CONFIRMED (source): backend/internal/usecase/order_usecase.go:168 CreateOrderWithItems commits the pending order and items first; POS confirmation/stock reservation runs afterwards (:201-205) and its failure returns an error (HTTP 400) while the draft order stays persisted ("complete draft is preserved for audit" comment), leaving orphan pending orders that stats and sales reports count (F-042)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
15 x 11000 = 165000 gross phantom pending sales

## Where it was reproduced
Set stock 5; 20 parallel POS orders of 1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
