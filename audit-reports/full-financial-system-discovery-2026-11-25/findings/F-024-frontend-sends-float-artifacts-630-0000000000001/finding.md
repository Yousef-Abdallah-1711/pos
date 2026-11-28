# F-024 — Frontend sends float artifacts (630.0000000000001) in monetary payloads

- **Severity:** Low  ·  **Category:** Software defect (minor)  ·  **Evidence status:** REPRODUCED (RUN-005 headed Chrome)
- **Affected modules:** Invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-OBSERVED-REPRO, EX-RUN005-INVOICE-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Exact decimals

Observed: tax_amount 630.0000000000001 (server stores rounded) RUN-005 (real invoice form, 1 x 5,000, 10% discount, 14% VAT): the POST payload carries "tax_amount":630.0000000000001 (raw float artefact) while the UI shows Tax EGP 630.00.

## Severity rationale
Rated **Low** (category: Software defect (minor); evidence: REPRODUCED (RUN-005 headed Chrome)). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Invoices. n/a

## Technical impact
CONFIRMED (RUN-004 source): InvoiceForm computes discountAmount, taxAmount and lineTotal with raw floating-point multiplication and sends them unrounded (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:620-650)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Create invoice with 14% tax

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
