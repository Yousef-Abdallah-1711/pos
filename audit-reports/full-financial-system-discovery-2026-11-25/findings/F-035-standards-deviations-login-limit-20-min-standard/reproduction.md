# F-035 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
25 rapid bad logins; response headers

## Expected result
Per CODING_STANDARDS 7.3/7.4

## Actual result (measured)
429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contains unsafe-inline and unsafe-eval in the isolated build; the API /health response carries no CSP header and the landing response carried none in this run.

## Re-runnable scripts / evidence
- `tests/runtime/ph17_r5_observed_repro.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
