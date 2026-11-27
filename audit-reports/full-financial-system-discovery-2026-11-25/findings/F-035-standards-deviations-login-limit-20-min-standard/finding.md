# F-035 — Standards deviations: login limit 20/min (standard 5/min); CSP allows unsafe-eval and unsafe-inline

- **Severity:** Low  ·  **Category:** Standards deviation  ·  **Evidence status:** OBSERVED
- **Affected modules:** Auth; Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Per CODING_STANDARDS 7.3/7.4

Observed: 429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval

## Monetary / operational impact
n/a

## Where it was reproduced
25 rapid bad logins; response headers

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
