# F-035 — Standards deviations: login limit 20/min (standard 5/min); CSP allows unsafe-eval and unsafe-inline

- **Severity:** Low  ·  **Category:** Standards deviation  ·  **Evidence status:** REPRODUCED (RUN-005 API-DIAGNOSTIC)
- **Affected modules:** Auth; Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-OBSERVED-REPRO
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Per CODING_STANDARDS 7.3/7.4

Observed: 429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contains unsafe-inline and unsafe-eval in the isolated build; the API /health response carries no CSP header and the landing response carried none in this run.

## Severity rationale
Rated **Low** (category: Standards deviation; evidence: REPRODUCED (RUN-005 API-DIAGNOSTIC)). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Auth; Platform. n/a

## Technical impact
CONFIRMED (source): backend/internal/middleware/cache.go:312-314 builds a CSP with script-src unsafe-inline and unsafe-eval; login limit is configured in middleware/rate_limit.go (observed 20/min, 14 failed attempts then 429 in SEC-012c)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
25 rapid bad logins; response headers

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
