# F-038 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Open /en/auth/register, fill fields, type the subdomain at ~90 ms/char, click Create account ~200 ms after typing (or with check-subdomain delayed 3 s)

## Expected result
Error clears (or submit proceeds) once availability is confirmed

## Actual result (measured)
No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)

## Re-runnable scripts / evidence
- `tests/runtime/ph05_subdomain_probe.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
