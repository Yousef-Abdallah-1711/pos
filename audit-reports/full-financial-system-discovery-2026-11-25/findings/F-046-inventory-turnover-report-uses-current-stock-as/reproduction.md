# F-046 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
GET /reports/inventory after recording a waste adjustment and sales (tenant 4)

## Expected result
Opening quantity from stock history; waste_qty reflects recorded waste

## Actual result (measured)
opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has waste 1 for products 2 and 4; Audit Burger original stock 50 | RUN-003: turnover opening_stock_qty == current stock (33/21); waste_qty 0 vs movements waste 1 (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph12_batch_apis2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
