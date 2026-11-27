# F-060 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
GET /admin/themes/{id}/history; DELETE /admin/themes/{active id}

## Expected result
History list; deletion of the active theme refused or falls back to default

## Actual result (measured)
history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph08_theme_driver.mjs`
- `tests/runtime/ph08_theme2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
