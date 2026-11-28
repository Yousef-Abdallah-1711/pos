# F-078 — Theme edits and activation do not refresh the public theme: the 5-minute server cache is never invalidated (and the storefront adds a 1-hour page cache)

- **Severity:** Low  ·  **Category:** Software defect (cache invalidation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Theme Builder, public theme API, storefront
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-THEMEBUILDER-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: The public theme reflects an activation or save immediately (or within seconds)

Observed: Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300

## Severity rationale
Rated **Low** (category: Software defect (cache invalidation); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Theme Builder, public theme API, storefront. 0

## Technical impact
CONFIRMED (source + runtime): invalidatePublicThemeCacheByRestaurant is defined (handler/http/theme_public_cache.go:45) but has no caller; GetThemeOnly caches the payload for publicThemeCacheTTL = 5 min (public_homepage_handler.go:117-144)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Tenant 4: prime GET /public/restaurants/audit-c-diag/theme, save #00ff00 in the editor, activate the theme, poll the public API every 45 s

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
