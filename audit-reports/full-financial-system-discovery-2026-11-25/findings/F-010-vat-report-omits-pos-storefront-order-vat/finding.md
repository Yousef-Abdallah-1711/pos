# F-010 — VAT report omits POS/storefront order VAT

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** VAT report
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Output VAT = GL 2120

Observed: output_vat 0 vs GL 1000 | RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 customer-invoice VAT) and GL VAT Payable (acct 87) credit balance 14,700; residual 14,000 = POS order VAT missing from report (STILL REPRODUCIBLE). 

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: 1000 here. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: VAT report. 1000 here

## Technical impact
CONFIRMED (source): backend/internal/repository/tax_repo.go:260-291 GetVATSummary sums tax_amount of invoices_customer / invoices_supplier only (status <> void, drafts included); order (POS/online) VAT never enters the report

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
1000 here

## Where it was reproduced
Open VAT report after POS sale

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
