# F-026 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Mark card order paid

## Expected result
Clearing cleared on settlement

## Actual result (measured)
1113 balance 11000 persists RUN-005: card order 5795 (110.00) marked paid -> journal Dr 1113 Card Payment Clearing 110.00 / Cr 4100 100.00 / Cr 2120 10.00; no settlement route exists.

## Re-runnable scripts / evidence
- `tests/runtime/ph17_r5_observed_repro.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
