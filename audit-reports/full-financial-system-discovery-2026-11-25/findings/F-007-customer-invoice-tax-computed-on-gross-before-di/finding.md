# F-007 — Customer invoice tax computed on gross before discount and UI/server totals differ

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customer invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CINV-001, EX-RUN005-INVOICE-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Tax 630 total 5130 (as UI shows)

Observed: Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE).  RUN-005 real-UI reconfirmation: same invoice: UI preview tax 630.00 / total 5,130.00, stored tax 700.00 / total 5,200.00 (tax on the pre-discount 5,000), journal Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700.

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: 70 per invoice. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Customer invoices. 70 per invoice

## Technical impact
CONFIRMED (source, RUN-003): DB trigger taxes gross | RUN-003 SOURCE: DB trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) computes tax_amount = ROUND(qty*unit_price*tax_rate/100,2) on the PRE-discount amount, then line_total = gross - discount + tax; the supplier twin (142_create_invoices_supplier.sql:75-79) is identical (see F-041).

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
70 per invoice

## Where it was reproduced
Invoice 5000 less 10% discount VAT 14%

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
