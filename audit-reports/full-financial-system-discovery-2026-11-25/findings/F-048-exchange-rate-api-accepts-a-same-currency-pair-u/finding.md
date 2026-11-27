# F-048 — Exchange-rate API accepts a same-currency pair (USD/USD)

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting - exchange rates
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-011
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Rejected (base must differ from quote)

Observed: 200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). 

## Monetary / operational impact
0 (data quality)

## Where it was reproduced
POST /accounting/exchange-rates base_currency=USD quote_currency=USD rate=1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
