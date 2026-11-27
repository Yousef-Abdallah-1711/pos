# F-057 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create delivery order (cash), assign driver, move status preparing -> ready -> out_for_delivery -> delivered, then read payment_status and journals

## Expected result
A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment

## Actual result (measured)
status delivered with payment_status pending and 0 journals; revenue stays unrecognised until someone manually marks paid; no driver cash reconciliation exists | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph08_theme_driver.mjs`
- `tests/runtime/ph08_theme2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
