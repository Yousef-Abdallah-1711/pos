# F-012 — At least six revenue/profit definitions disagree and are unlabeled

- **Severity:** High  ·  **Category:** Documentation/design inconsistency  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Reports; Dashboard
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Reports state basis or reconcile

Observed: Sales 121000; Profitability 110000; KPI 22000; GL revenue 20000; KPI net profit 6000 vs P&L 4000; KPI cash position 0 vs GL cash

## Monetary / operational impact
up to 101000 spread

## Where it was reproduced
Open Sales Profitability KPI P&L after mixed paid/unpaid orders

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
