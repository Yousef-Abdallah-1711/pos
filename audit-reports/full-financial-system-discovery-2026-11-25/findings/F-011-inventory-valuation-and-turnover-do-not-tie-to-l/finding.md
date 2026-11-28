# F-011 — Inventory valuation and turnover do not tie to ledger; opening stock has no journal

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Inventory; Products; Reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-OBS-PUR-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Valuation = GL inventory

Observed: 196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no journal (STILL REPRODUCIBLE). 

## Severity rationale
Rated **High** (category: Software defect / missing capability; evidence: REPRODUCED). Basis recorded in the register: 180000 / 148000. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Inventory; Products; Reports. 180000 / 148000

## Technical impact
CONFIRMED (source): backend/internal/usecase/financial_report_usecase.go:488-526 GenerateInventoryValuation = SUM(quantity_in_stock x cost) for tracked products; it never reads GL 1130, and purchases/opening stock never post to 1130 (supplier invoice lines post to 5110, OBS-PUR-001)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
180000 / 148000

## Where it was reproduced
Create product with stock 50 cost 4000 then valuation

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
