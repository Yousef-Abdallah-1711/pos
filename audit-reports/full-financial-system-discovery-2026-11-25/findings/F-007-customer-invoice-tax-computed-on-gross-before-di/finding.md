# F-007 — Customer invoice tax computed on gross before discount and UI/server totals differ

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customer invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CINV-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Tax 630 total 5130 (as UI shows)

Observed: Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE). 

## Monetary / operational impact
70 per invoice

## Where it was reproduced
Invoice 5000 less 10% discount VAT 14%

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
