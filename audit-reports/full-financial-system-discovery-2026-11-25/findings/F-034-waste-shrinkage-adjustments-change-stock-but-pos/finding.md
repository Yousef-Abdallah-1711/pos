# F-034 — Waste/shrinkage adjustments change stock but post nothing and the turnover report shows waste 0

- **Severity:** Medium  ·  **Category:** Missing capability / defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Inventory; Reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Expense posted and waste reported

Observed: Stock 48 to 47; no journal; waste_qty 0; reason count_shortage not accepted (allowed: purchase, sale, adjustment, waste, return)

## Monetary / operational impact
n/a

## Where it was reproduced
POST inventory adjust reason waste -1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
