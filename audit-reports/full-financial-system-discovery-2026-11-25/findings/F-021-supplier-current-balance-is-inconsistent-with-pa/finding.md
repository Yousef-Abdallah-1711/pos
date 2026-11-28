# F-021 — Supplier current_balance is inconsistent with payables

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Suppliers
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-F021-SUPPLIER-BALANCE
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Balance = AP

Observed: current_balance -22000 vs AP 800

## Severity rationale
Rated **Medium** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: n/a (display). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Suppliers. n/a (display)

## Technical impact
CONFIRMED (RUN-005 source + replay): suppliers.current_balance has exactly one money-movement writer: a supplier payment subtracts the amount (backend/internal/repository/payment_acct_repo.go:302-306), while supplier invoice creation and approval never add to it (no other UPDATE suppliers exists; supplier_repo.go:172-200 writes the value supplied by the handler, which preserves the stored value, supplier_handler.go:7

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a (display)

## Where it was reproduced
Approve invoice then pay

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
