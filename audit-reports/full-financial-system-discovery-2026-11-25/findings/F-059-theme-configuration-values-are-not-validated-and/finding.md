# F-059 — Theme configuration values are not validated and are served publicly (invalid colors and url(javascript:) accepted)

- **Severity:** Medium  ·  **Category:** Security / validation  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Theme builder; Storefront; public theme API
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-THM-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Colors validated (hex/rgb/hsl allow-list); unsafe values rejected

Observed: 201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline not tested (PROBABLE CSS injection) | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url(javascript:) stored and served by public theme endpoint (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Security / validation; evidence: REPRODUCED). Basis recorded in the register: 0 (integrity/security of public pages). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Theme builder; Storefront; public theme API. 0 (integrity/security of public pages)

## Technical impact
CONFIRMED (RUN-004 source absence + RUN-003 runtime): the theme handler/service persist config JSON without colour or value validation (no validate call in the theme service; url(javascript:) was accepted and served)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 (integrity/security of public pages)

## Where it was reproduced
POST /admin/themes with colors.primary "not-a-color" and colors.text "url(javascript:alert(1))", activate, GET /public/restaurants/{slug}/theme

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
