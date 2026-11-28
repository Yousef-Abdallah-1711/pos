# F-059 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /admin/themes with colors.primary "not-a-color" and colors.text "url(javascript:alert(1))", activate, GET /public/restaurants/{slug}/theme

## Expected result
Colors validated (hex/rgb/hsl allow-list); unsafe values rejected

## Actual result (measured)
201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline not tested (PROBABLE CSS injection) | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url(javascript:) stored and served by public theme endpoint (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph08_theme_driver.mjs`
- `tests/runtime/ph08_theme2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
