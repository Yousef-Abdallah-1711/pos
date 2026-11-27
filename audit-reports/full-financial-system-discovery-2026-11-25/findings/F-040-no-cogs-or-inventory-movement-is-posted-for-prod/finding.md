# F-040 — No COGS or inventory movement is posted for products with track_inventory=false even when a cost is set; Profitability report still deducts cost

- **Severity:** Medium  ·  **Category:** Software defect / policy  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Products; POS/Orders; Accounting; Profitability report
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-OBS-INV-001, EX-RUN003-AUD-INV-COGS-UNTRACKED
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Cost-bearing product sold produces COGS consistent with the Profitability report (or the report excludes untracked cost); policy for untracked items documented

Observed: Order 39 (untracked): JE-2026-00064 posts only Dr 1110 330 / Cr 4100 300 / Cr 2120 30 - no COGS; tracked order 41: JE-2026-00066 Dr 5110 120 / Cr 1130 120. Profitability report shows cogs 120 for both products => report COGS exceeds GL COGS by 120 | RUN-003 controlled experiment: untracked product order 39 posts no COGS while Profitability report shows cogs 120; tracked twin order 41 posts COGS 120 (CONFIRMED, auto_posting_engine.go:173). 

## Monetary / operational impact
120 in test (cost x quantity of every untracked sale)

## Where it was reproduced
Create product with cost 40 and track_inventory=false (API default when the flag is omitted), POS order qty 3, mark paid; compare with identical tracked product

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
