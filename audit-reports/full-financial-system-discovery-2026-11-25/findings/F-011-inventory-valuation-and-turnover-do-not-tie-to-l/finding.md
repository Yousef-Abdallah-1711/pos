# F-011 — Inventory valuation and turnover do not tie to ledger; opening stock has no journal

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Inventory; Products; Reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-OBS-PUR-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Valuation = GL inventory

Observed: 196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no journal (STILL REPRODUCIBLE). 

## Monetary / operational impact
180000 / 148000

## Where it was reproduced
Create product with stock 50 cost 4000 then valuation

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
