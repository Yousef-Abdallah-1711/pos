# F-048 — Exchange-rate API accepts a same-currency pair (USD/USD)

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting - exchange rates
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-011
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Rejected (base must differ from quote)

Observed: 200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). 

## Severity rationale
Rated **Low** (category: Validation gap; evidence: REPRODUCED). Basis recorded in the register: 0 (data quality). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting - exchange rates. 0 (data quality)

## Technical impact
CONFIRMED (RUN-004 source): the handler only requires non-empty codes and a positive rate (accounting_advanced_handler.go:475-482); there is no base != quote check

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 (data quality)

## Where it was reproduced
POST /accounting/exchange-rates base_currency=USD quote_currency=USD rate=1

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
