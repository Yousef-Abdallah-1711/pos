# F-054 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create annual leave 2026-12-01..12-30 total_days 30 for an employee with allowance 21 and approve; create a 3-day range with total_days 40

## Expected result
Reject or flag leave beyond remaining entitlement; total_days derived from dates

## Actual result (measured)
30-day leave approved; report shows used 30, remaining 0 (clamped) vs allowance 21; 3-day range stored total_days 40 | RUN-003: 30-day annual leave approved over allowance 21; total_days 40 stored for a 3-day range (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph10_hr_matrix.mjs`
- `tests/runtime/ph10_hr_matrix2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
