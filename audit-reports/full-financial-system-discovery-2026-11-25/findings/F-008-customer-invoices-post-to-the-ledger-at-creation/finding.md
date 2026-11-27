# F-008 — Customer invoices post to the ledger at creation while status is draft

- **Severity:** Medium  ·  **Category:** Software defect / policy  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customer invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CINV-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Draft has no accounting effect

Observed: Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700, i.e. posted at creation; the draft state itself could not be re-observed because the invoice has since been paid - fresh draft re-run scheduled in T-P09-001.  | RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE). 

## Monetary / operational impact
5200 AR+revenue recognised at draft

## Where it was reproduced
Create invoice with status draft

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
