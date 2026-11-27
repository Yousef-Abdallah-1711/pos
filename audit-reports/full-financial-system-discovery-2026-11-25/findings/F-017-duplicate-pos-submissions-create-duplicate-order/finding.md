# F-017 — Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order creation

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Inventory
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-REL-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: One order

Observed: 5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). 

## Monetary / operational impact
duplicate revenue risk

## Where it was reproduced
5 identical parallel POST and 3 with same key

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
