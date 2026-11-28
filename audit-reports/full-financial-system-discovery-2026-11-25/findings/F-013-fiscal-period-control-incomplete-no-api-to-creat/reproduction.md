# F-013 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Insert closed period (SQL fixture) then post / create / void

## Expected result
Closed period blocks all postings including voids

## Actual result (measured)
Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
