# F-077 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome tenant 4: save #ff0000 on theme 13 in the editor and activate it; open Settings > Appearance and click Theme Colors; read themes_v2.config

## Expected result
A user preference changes only that user settings

## Actual result (measured)
The click sends PUT /users/4/settings/theme-colors and then PUT /admin/themes/{active id} with the merged colours (primary, secondary, accent from the user profile); the active theme colours and updated_at change (observed on themes 4 and 13), so edits made in the Theme Builder are overwritten whenever any user with settings access saves Appearance

## Re-runnable scripts / evidence
- `tests/runtime/ph05_r4_appearance_theme.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
