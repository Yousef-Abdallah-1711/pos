# F-032 — Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not behind the auth provider

- **Severity:** Low  ·  **Category:** Security/hygiene  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Dashboard app
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: No test pages in production build

Observed: Page renders TEST PAGE LOADED with no auth protection and a hydration error RUN-005: GET http://localhost:13112/en/test-url returned HTTP 200 (curl) on the current isolated dashboard.

## Severity rationale
Rated **Low** (category: Security/hygiene; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Dashboard app. n/a

## Technical impact
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx exists in the shipped dashboard app

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Open /en/test-url

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
