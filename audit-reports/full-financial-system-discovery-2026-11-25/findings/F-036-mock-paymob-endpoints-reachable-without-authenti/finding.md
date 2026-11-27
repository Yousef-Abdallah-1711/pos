# F-036 — Mock Paymob endpoints reachable without authentication in non-production environments

- **Severity:** Low  ·  **Category:** Security (environment)  ·  **Evidence status:** OBSERVED
- **Affected modules:** Payments (mock)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Mock disabled outside test

Observed: Return 400 (reachable) without token

## Monetary / operational impact
n/a

## Where it was reproduced
Unauthenticated GET /mock/paymob/iframe and POST /mock/paymob/process

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
