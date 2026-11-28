# F-041 — Supplier invoice VAT is computed on the pre-discount amount, overstating input VAT and accounts payable

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-SUP-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle)

Observed: API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs oracle 1,260/10,260; JE Dr 5110 9,000 Dr 1150 1,400 Cr 2101 10,400 (CONFIRMED trigger 142:75-79). 

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: 140 per 10,000 gross at 10% discount and 14% VAT (VAT x discount). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports. 140 per 10,000 gross at 10% discount and 14% VAT (VAT x discount)

## Technical impact
CONFIRMED (source): trigger recalculate_invoice_supplier_line_total backend/migrations/142_create_invoices_supplier.sql:75-79 computes tax on quantity*unit_price before subtracting the discount (same formula as customer trigger, F-007)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
140 per 10,000 gross at 10% discount and 14% VAT (VAT x discount)

## Where it was reproduced
Create supplier invoice line qty 10 x 1,000 with discount 10% and tax 14% (API-DIAGNOSTIC on isolated tenant 4), then approve

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
