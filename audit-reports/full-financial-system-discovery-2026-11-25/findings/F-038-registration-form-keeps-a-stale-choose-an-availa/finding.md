# F-038 — Registration form keeps a stale "choose an available subdomain" error after availability check succeeds

- **Severity:** Low  ·  **Category:** UX issue  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome, human-paced typing; 2 of 3 early-submit scenarios + slow network)
- **Affected modules:** Landing / registration
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN003-OBS-ONB-001-S1, EX-RUN003-OBS-ONB-001-S4, EX-RUN003-OBS-ONB-001-S5
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Error clears (or submit proceeds) once availability is confirmed

Observed: No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)

## Monetary / operational impact
none

## Where it was reproduced
Open /en/auth/register, fill fields, type the subdomain at ~90 ms/char, click Create account ~200 ms after typing (or with check-subdomain delayed 3 s)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
