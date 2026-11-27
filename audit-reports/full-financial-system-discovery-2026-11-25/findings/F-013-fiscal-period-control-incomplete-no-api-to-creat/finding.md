# F-013 — Fiscal-period control incomplete: no API to create/close periods; void and draft creation allowed in a closed period

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-006
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Closed period blocks all postings including voids

Observed: Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). 

## Monetary / operational impact
100000 reversed in closed period

## Where it was reproduced
Insert closed period (SQL fixture) then post / create / void

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
