# F-060 — Theme history endpoint returns 501 and the active theme can be deleted

- **Severity:** Low  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Theme builder
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-THM-001, EX-RUN004-THEMEBUILDER-UI, EX-RUN005-DESTR-THEME
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: History list; deletion of the active theme refused or falls back to default

Observed: history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). 

## Severity rationale
Rated **Low** (category: Missing capability; evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Theme builder. 0

## Technical impact
CONFIRMED (RUN-004 source): GetThemeHistory returns 501 Not Implemented (admin_theme_handler.go:68)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
GET /admin/themes/{id}/history; DELETE /admin/themes/{active id}

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
