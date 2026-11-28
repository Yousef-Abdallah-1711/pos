# F-016 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Open payroll page with record period ending after today

## Expected result
Current records visible

## Actual result (measured)
0 records until To-date changed | RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI says 'Create your first salary record', while salaries row id 1 (tenant 4, period 2026-11-01..2026-11-30, paid) exists (STILL REPRODUCIBLE). 

## Re-runnable scripts / evidence
- `tests/runtime/ph20_evidence_gaps.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
