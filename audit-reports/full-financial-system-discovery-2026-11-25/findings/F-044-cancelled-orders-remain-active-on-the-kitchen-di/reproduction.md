# F-044 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Cancel a pending order (DELETE /admin/orders/{id}) and open the KDS board (headed Chrome)

## Expected result
Cancelled order is removed from the active kitchen queue

## Actual result (measured)
Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24). 

## Re-runnable scripts / evidence
- `tests/runtime/ph08_kds_cancelled.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
