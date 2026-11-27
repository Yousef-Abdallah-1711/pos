# F-052 — Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the money twice

- **Severity:** High  ·  **Category:** Software defect (idempotency)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Receipts; Supplier payments; Accounting (cash, AR/AP)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-REL-025
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Second request returns the first result; one receipt

Observed: Both return 201; payments_received has 2 rows (sum 800); invoice amount_paid 800 / balance 200; AR credited twice | RUN-003 REL-025: identical Idempotency-Key replay creates 2 receipts (sum 800) (REPRODUCED). 

## Monetary / operational impact
Every retried or double-clicked payment duplicates cash and AR/AP movement

## Where it was reproduced
POST /accounting/payments/received twice with the same UUID Idempotency-Key and body (400 on a 1,000 invoice)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
