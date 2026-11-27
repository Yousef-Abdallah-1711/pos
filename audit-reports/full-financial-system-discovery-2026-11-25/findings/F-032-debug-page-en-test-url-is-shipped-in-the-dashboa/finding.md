# F-032 — Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not behind the auth provider

- **Severity:** Low  ·  **Category:** Security/hygiene  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Dashboard app
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: No test pages in production build

Observed: Page renders TEST PAGE LOADED with no auth protection and a hydration error

## Monetary / operational impact
n/a

## Where it was reproduced
Open /en/test-url

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
