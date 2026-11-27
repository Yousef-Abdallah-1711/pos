# F-024 — Frontend sends float artifacts (630.0000000000001) in monetary payloads

- **Severity:** Low  ·  **Category:** Software defect (minor)  ·  **Evidence status:** OBSERVED
- **Affected modules:** Invoices
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Exact decimals

Observed: tax_amount 630.0000000000001 (server stores rounded)

## Monetary / operational impact
n/a

## Where it was reproduced
Create invoice with 14% tax

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
