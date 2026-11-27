# F-076 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome, Super Admin app (port 13114), Users page: click View on a user row

## Expected result
A user detail page (or no link)

## Actual result (measured)
Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route

## Re-runnable scripts / evidence
- `tests/runtime/ph13_r4_superadmin_ui.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
