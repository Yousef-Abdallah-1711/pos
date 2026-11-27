# F-077 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Saving personal Appearance colours also rewrites the active public storefront theme, silently overwriting Theme Builder edits.

## 2. Current behavior (verified)
The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accent from the user profile); the active theme colours and updated_at change (observed on themes 4 and 13), so edits made in the Theme Builder are overwritten whenever any user with settings access saves Appearance

## 3. Required behavior
A user preference changes only that user settings

## 4. Reproduction conditions
Headed Chrome tenant 4: save #ff0000 on theme 13 in the editor and activate it; open Settings > Appearance and click Theme Colors; read themes_v2.config

## 5. Verified evidence
- evidence/network/RUN-004/PH-05/AUD-R4-APPEARANCE-THEME.json
- evidence/network/RUN-004/PH-08/AUD-R4-THEME-SAVE-VS-ACTIVATE.json
- tests/runtime/ph05_r4_appearance_theme.mjs

## 6. Root cause and confidence
CONFIRMED — syncBrandColorsToActiveTheme in frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95 copies the user colours into the active tenant theme via themeApi.updateTheme

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx`:38, `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx`:78, `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx`:38-95
- Frontend: `frontend/apps/restaurant-website/app/[locale]`, `frontend/apps/dashboard/src/app/[locale]/dashboard/settings`
- API contracts: no API change; PUT /admin/themes/{id} is no longer called from Appearance
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Settings > Appearance, Theme Builder, public storefront workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: remove syncBrandColorsToActiveTheme (frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx:38-95) from the personal Appearance save; keep the per-user light/dark preference only. If brand colours must stay in step with the public theme, expose them as one tenant-level setting inside the Theme Builder (permission: theme edit) and have the Appearance page link to it. The backend PUT /admin/themes/{id} should require the theme module WRITE permission (check the wrapWithPermission for the route) and record who changed the theme.

## 10. Alternative designs
A: single owner of brand colours in the Theme Builder (recommended). B: keep the sync but ask for confirmation and show the diff (still two writers).

## 11. Required policy decisions
Who may change public branding (owner decision).

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-059

## 21. Risks
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Saving Appearance leaves the active theme row unchanged; Theme Builder colours persist after an Appearance save (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Saving Appearance leaves the active themes_v2 row (config and updated_at) unchanged.
- A colour saved in the Theme Builder persists after any Appearance save.
- A user without theme-edit permission cannot change public colours through any screen.
