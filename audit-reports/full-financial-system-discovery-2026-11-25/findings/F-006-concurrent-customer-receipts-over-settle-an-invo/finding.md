# F-006 — Concurrent customer receipts over-settle an invoice

- **Severity:** High  ·  **Category:** Software defect (race)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Customer invoices; Receipts
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-AR-RESIDUAL, EX-RUN003-AUD-SUP-002, EX-RUN003-AUD-CINV-FLOW
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: At most 3 accepted; no over-collection

Observed: 5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment.  | RUN-003 FRESH-INVOICE RE-REPRODUCTION (AUD-CINV-FLOW): invoice 3,000, 5 concurrent receipts of 1,000 -> all 5 accepted (201), amount_paid 6,000, AR credit -3,000; single sequential overpayment 5,000 is rejected (guard exists but is not serialized). RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment. 

## Severity rationale
Rated **High** (category: Software defect (race); evidence: REPRODUCED). Basis recorded in the register: 1800. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Customer invoices; Receipts. 1800

## Technical impact
CONFIRMED (source + RUN-003 reproduction on fresh invoice): backend/internal/repository/payment_acct_repo.go:445-449 customer-invoice balance SELECT has no FOR UPDATE (supplier twin :499-503 has FOR UPDATE OF inv); applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0), hiding the overpayment

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
1800

## Where it was reproduced
5 parallel receipts of 1000 on balance 3200

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
