# F-048 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /accounting/exchange-rates base_currency=USD quote_currency=USD rate=1

## Expected result
Rejected (base must differ from quote)

## Actual result (measured)
200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). 

## Re-runnable scripts / evidence
- `tests/runtime/ph11_accounting_flow2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
