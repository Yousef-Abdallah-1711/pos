# F-005 — Card and wallet POS orders never become paid; stock is consumed anyway

- **Severity:** High  ·  **Category:** Missing capability / defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Accounting; Inventory
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Payment captured or at least recorded; revenue recognised per policy

Observed: Order stays pending; no outbox; no journal; stock decremented (8 unpaid orders = 88000 not in ledger)

## Monetary / operational impact
88000 gross unrecorded in this run

## Where it was reproduced
POS order with Card or Wallet

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
