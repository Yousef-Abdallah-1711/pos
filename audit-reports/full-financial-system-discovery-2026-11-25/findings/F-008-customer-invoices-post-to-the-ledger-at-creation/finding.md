# F-008 — Customer invoices post to the ledger at creation while status is draft

- **Severity:** Medium  ·  **Category:** Software defect / policy  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customer invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CINV-001, EX-RUN005-INVOICE-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Draft has no accounting effect

Observed: Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700, i.e. posted at creation; the draft state itself could not be re-observed because the invoice has since been paid - fresh draft re-run scheduled in T-P09-001.  | RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE).  RUN-005 real-UI reconfirmation: draft invoice INV-2026-00005 was posted to journal entry 44787 at creation.

## Severity rationale
Rated **Medium** (category: Software defect / policy; evidence: REPRODUCED). Basis recorded in the register: 5200 AR+revenue recognised at draft. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Customer invoices. 5200 AR+revenue recognised at draft

## Technical impact
CONFIRMED (source): backend/internal/repository/invoice_customer_repo.go:452-560 createInvoiceTx calls syncJournalEntryTx (:539) which builds (buildJournalEntry :561-636) and posts the AR/revenue/VAT journal for every invoice at creation regardless of status; draft invoices therefore hit AR, revenue and VAT

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
5200 AR+revenue recognised at draft

## Where it was reproduced
Create invoice with status draft

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
