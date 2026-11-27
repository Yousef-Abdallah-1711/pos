# F-060 — Theme history endpoint returns 501 and the active theme can be deleted

- **Severity:** Low  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Theme builder
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-THM-001, EX-RUN004-THEMEBUILDER-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: History list; deletion of the active theme refused or falls back to default

Observed: history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). 

## Monetary / operational impact
0

## Where it was reproduced
GET /admin/themes/{id}/history; DELETE /admin/themes/{active id}

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
