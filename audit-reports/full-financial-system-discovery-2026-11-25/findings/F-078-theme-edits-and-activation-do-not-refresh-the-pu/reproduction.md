# F-078 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Tenant 4: prime GET /public/restaurants/audit-c-diag/theme, save #00ff00 in the editor, activate the theme, poll the public API every 45 s

## Expected result
The public theme reflects an activation or save immediately (or within seconds)

## Actual result (measured)
Public API kept returning #3b82f6 for 182 s after activation although the database held #00ff00, and flipped only at 227 s when the 5-minute entry expired; the storefront page additionally revalidates every 3,600 s (app/[locale]/page.tsx:75,131) and the response carries Cache-Control max-age=300

## Re-runnable scripts / evidence
- `tests/runtime/ph08_r4_theme_cache_ttl2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
