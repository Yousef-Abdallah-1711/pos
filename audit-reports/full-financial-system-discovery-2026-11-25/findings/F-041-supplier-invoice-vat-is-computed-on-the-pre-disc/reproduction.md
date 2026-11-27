# F-041 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create supplier invoice line qty 10 x 1,000 with discount 10% and tax 14% (API-DIAGNOSTIC on isolated tenant 4), then approve

## Expected result
Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle)

## Actual result (measured)
API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs oracle 1,260/10,260; JE Dr 5110 9,000 Dr 1150 1,400 Cr 2101 10,400 (CONFIRMED trigger 142:75-79). 

## Re-runnable scripts / evidence
- `tests/runtime/ph06_supplier_flow.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
