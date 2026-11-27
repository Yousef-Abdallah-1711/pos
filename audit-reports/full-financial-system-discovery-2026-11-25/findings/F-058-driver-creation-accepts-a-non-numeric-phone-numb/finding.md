# F-058 — Driver creation accepts a non-numeric phone number

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Drivers
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-DRV-002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 phone format error

Observed: 201 driver id 2 stored with phone_number "abc" | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): driver phone abc accepted (201) (REPRODUCED). 

## Monetary / operational impact
0

## Where it was reproduced
POST /admin/drivers with phone_number "abc"

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
