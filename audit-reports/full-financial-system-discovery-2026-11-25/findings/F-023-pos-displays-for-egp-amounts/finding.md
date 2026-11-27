# F-023 — POS displays $ for EGP amounts

- **Severity:** Low  ·  **Category:** UX defect  ·  **Evidence status:** OBSERVED
- **Affected modules:** POS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-POS-001, EX-RUN003-AUD-POS-002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: EGP symbol

Observed: Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE). 

## Monetary / operational impact
n/a

## Where it was reproduced
Open POS

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
