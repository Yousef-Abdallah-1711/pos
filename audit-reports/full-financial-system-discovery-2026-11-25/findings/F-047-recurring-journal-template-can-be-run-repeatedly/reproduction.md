# F-047 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create template (auto_post) and POST /accounting/recurring-journals/{id}/run twice on 2026-11-26

## Expected result
Second run for the same period rejected or returns the first journal

## Actual result (measured)
Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 2026-12-26 | RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_advanced_repo.go:396-420). 

## Re-runnable scripts / evidence
- `tests/runtime/ph11_accounting_flow2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
