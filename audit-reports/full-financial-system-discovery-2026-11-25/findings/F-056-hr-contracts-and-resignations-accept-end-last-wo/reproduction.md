# F-056 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /hr/contracts with start_date 2026-11-01 and end_date 2026-10-01; POST /hr/resignations with resignation_date 2026-11-26 and last_working_date 2026-11-01

## Expected result
400 with a clear message

## Actual result (measured)
both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected only by DB check constraints that surface as 500 | RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph10_hr_modules.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
