# F-062 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /products with price -5, price 0, cost -1, cost 50 on price 10, empty name_en, quantity_in_stock -4 (tracked), discount_price 15 on price 10, price 10.12345; POST suppliers/customers with invalid email, duplicates, negative credit limit; duplicate category name

## Expected result
400 with field errors; sellable products must have price > 0 and 2-decimal money

## Actual result (measured)
all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-length return 500 with raw pq text; the SKU unique violation returns 500; selling the negative-price product is blocked later by the order path ("invalid product price"), price 0 is sellable (F-061), 10.12345 is rounded per unit at order time | RUN-003 master-data battery: negative/zero price, negative cost/stock, empty name, discount above price, 5-decimal price, duplicate/malformed supplier/customer accepted (CONFIRMED product_usecase.go:294-306). 

## Re-runnable scripts / evidence
- `tests/runtime/ph16_masterdata.mjs`
- `tests/runtime/ph16_negprice.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
