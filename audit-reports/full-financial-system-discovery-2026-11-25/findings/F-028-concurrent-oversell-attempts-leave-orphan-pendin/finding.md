# F-028 — Concurrent oversell attempts leave orphan pending orders with items although the API returned 400

- **Severity:** High  ·  **Category:** Software defect (race)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Inventory; KDS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 5 orders accepted; 15 rejected leave no order

Observed: 5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in Sales Analytics

## Monetary / operational impact
15 x 11000 = 165000 gross phantom pending sales

## Where it was reproduced
Set stock 5; 20 parallel POS orders of 1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
