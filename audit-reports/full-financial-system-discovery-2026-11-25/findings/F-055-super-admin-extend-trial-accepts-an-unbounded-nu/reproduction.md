# F-055 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /super-admin/tenants/{id}/extend-trial {days:100000}

## Expected result
Upper bound / confirmation and audit reason

## Actual result (measured)
200 days_added 100000, new_trial_end_date 2300-10-02 (negative days correctly rejected) | RUN-003: extend-trial 100000 days accepted -> 2300-10-02 (REPRODUCED; tenant 8 restored by fixture). 

## Re-runnable scripts / evidence
- `tests/runtime/ph13_saas.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
