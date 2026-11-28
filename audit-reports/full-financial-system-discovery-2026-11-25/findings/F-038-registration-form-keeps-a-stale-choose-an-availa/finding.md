# F-038 — Registration form keeps a stale "choose an available subdomain" error after availability check succeeds

- **Severity:** Low  ·  **Category:** UX issue  ·  **Evidence status:** REPRODUCED (RUN-003 headed Chrome, human-paced typing; 2 of 3 early-submit scenarios + slow network)
- **Affected modules:** Landing / registration
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-OBS-ONB-001-S1, EX-RUN003-OBS-ONB-001-S4, EX-RUN003-OBS-ONB-001-S5
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Error clears (or submit proceeds) once availability is confirmed

Observed: No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)

## Severity rationale
Rated **Low** (category: UX issue; evidence: REPRODUCED (RUN-003 headed Chrome, human-paced typing; 2 of 3 early-submit scenarios + slow network)). Basis recorded in the register: none. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Landing / registration. none

## Technical impact
CONFIRMED (RUN-005 source + RUN-003 runtime): in the current landing code (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx) onSubmit sets the slug error with setError when availability is not available (lines 59-64) and nothing clears it later: the only useEffect (lines 20-24) handles billing and analytics, availability is read from useDebouncedAvailability (line 39) without an effect that calls clearEr

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
none

## Where it was reproduced
Open /en/auth/register, fill fields, type the subdomain at ~90 ms/char, click Create account ~200 ms after typing (or with check-subdomain delayed 3 s)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
