# F-078 — Theme edits and activation do not refresh the public theme: the 5-minute server cache is never invalidated (and the storefront adds a 1-hour page cache)

- **Severity:** Low  ·  **Category:** Software defect (cache invalidation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Theme Builder, public theme API, storefront
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-THEMEBUILDER-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: The public theme reflects an activation or save immediately (or within seconds)

Observed: Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300

## Monetary / operational impact
0

## Where it was reproduced
Tenant 4: prime GET /public/restaurants/audit-c-diag/theme, save #00ff00 in the editor, activate the theme, poll the public API every 45 s

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
