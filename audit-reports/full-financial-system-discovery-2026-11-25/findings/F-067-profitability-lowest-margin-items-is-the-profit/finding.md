# F-067 — Profitability "Lowest Margin Items" is the profit-amount ranking read backwards, so a 100%-margin item ranks second-lowest

- **Severity:** Medium  ·  **Category:** Software defect (report logic)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Reports > Profitability (API, UI, CSV/PDF)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-PROF007
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Items ordered by gross_margin_pct ascending (0%, 60%, 60%, 60%, 100%)

Observed: Order is MD probe 0%, MD probe 100%, Probe A 60%, Probe B 60%, Burger 60%: the list is the gross-PROFIT descending list iterated backwards

## Monetary / operational impact
0 direct (decision quality: menu pricing review points at the wrong items)

## Where it was reproduced
Headed Chrome /en/dashboard/reports/profitability, tenant 4, 2026-10-29..2026-11-27; compare lowest_margin_items with margin percentages from SQL

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
