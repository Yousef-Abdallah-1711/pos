# F-082 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome, viewport 768x1024, English, tenant 4 owner: open /en/dashboard/reports/financial/profit-loss and compare document scrollWidth with the viewport width

## Expected result
No horizontal page scroll; wide tables scroll inside their own container

## Actual result (measured)
document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1440 px

## Re-runnable scripts / evidence
- `tests/runtime/ph04_r5_lang_viewport.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
