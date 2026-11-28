# F-036 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Unauthenticated GET /mock/paymob/iframe and POST /mock/paymob/process

## Expected result
Mock disabled outside test

## Actual result (measured)
Return 400 (reachable) without token RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> 400 "payment_token is required": both execute their handler without any session.

## Re-runnable scripts / evidence
- `tests/runtime/ph17_r5_observed_repro.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
