# F-077 — Saving personal Appearance colours also rewrites the active public storefront theme, silently overwriting Theme Builder edits

- **Severity:** Medium  ·  **Category:** Software defect (cross-feature coupling)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Settings > Appearance, Theme Builder, public storefront
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-THEMEBUILDER-UI
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A user preference changes only that user settings

Observed: The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accent from the user profile); the active theme colours and updated_at change (observed on themes 4 and 13), so edits made in the Theme Builder are overwritten whenever any user with settings access saves Appearance

## Severity rationale
Rated **Medium** (category: Software defect (cross-feature coupling); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Settings > Appearance, Theme Builder, public storefront. 0

## Technical impact
CONFIRMED (source + runtime): syncBrandColorsToActiveTheme in frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95 copies the user colours into the active tenant theme via themeApi.updateTheme

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome tenant 4: save #ff0000 on theme 13 in the editor and activate it; open Settings > Appearance and click Theme Colors; read themes_v2.config

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
