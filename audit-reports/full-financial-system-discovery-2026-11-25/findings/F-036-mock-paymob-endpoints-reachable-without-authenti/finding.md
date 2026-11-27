# F-036 — Mock Paymob endpoints reachable without authentication in non-production environments

- **Severity:** Low  ·  **Category:** Security (environment)  ·  **Evidence status:** REPRODUCED (RUN-005 API-DIAGNOSTIC)
- **Affected modules:** Payments (mock)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-OBSERVED-REPRO
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Mock disabled outside test

Observed: Return 400 (reachable) without token RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> 400 "payment_token is required": both execute their handler without any session.

## Severity rationale
Rated **Low** (category: Security (environment); evidence: REPRODUCED (RUN-005 API-DIAGNOSTIC)). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Payments (mock). n/a

## Technical impact
CONFIRMED (source): backend/cmd/api/main.go:667-668 registers GET /mock/paymob/iframe and POST /mock/paymob/process with mux.HandleFunc and no auth wrapper (non-production only); RUN-003 PAY-01 reached the handler unauthenticated (400 payment_token required)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Unauthenticated GET /mock/paymob/iframe and POST /mock/paymob/process

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
