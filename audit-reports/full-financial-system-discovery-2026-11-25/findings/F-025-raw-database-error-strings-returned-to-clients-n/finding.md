# F-025 — Raw database error strings returned to clients (numeric overflow; constraint names)

- **Severity:** Low  ·  **Category:** Security/UX  ·  **Evidence status:** REPRODUCED
- **Affected modules:** API
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-UIFORM-003
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Sanitised message

Observed: pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user receives raw text such as pq: value too long for type character varying(3) (REPRODUCED; the filler used the same long synthetic value in every text field, so the inputs are artificial but the exposure is real). 

## Monetary / operational impact
n/a

## Where it was reproduced
Submit 1e15 amount

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
