# F-005 — Card and wallet POS orders never become paid; stock is consumed anyway

- **Severity:** High  ·  **Category:** Missing capability / defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Accounting; Inventory
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Payment captured or at least recorded; revenue recognised per policy

Observed: Order stays pending; no outbox; no journal; stock decremented (8 unpaid orders = 88000 not in ledger)

## Severity rationale
Rated **High** (category: Missing capability / defect; evidence: REPRODUCED). Basis recorded in the register: 88000 gross unrecorded in this run. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: POS; Orders; Accounting; Inventory. 88000 gross unrecorded in this run

## Technical impact
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:302-308 calls updateOrderPaymentStatus(paid) only when paymentMethod === cash; card and wallet orders stay payment_status pending, and no other flow settles them (provider settlement belongs to subscription code)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
88000 gross unrecorded in this run

## Where it was reproduced
POS order with Card or Wallet

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
