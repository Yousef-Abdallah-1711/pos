# F-071 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Start the isolated API image with ENV=staging and RBAC_ENFORCEMENT unset (same DB/Redis) next to the enforce-mode API; log in as the cashier on both and request the same endpoints; production-mode probe for contrast (API-DIAGNOSTIC)

## Expected result
A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented

## Actual result (measured)
Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production RUN-005 expansion (4 roles x 28 GET endpoints + 8 write probes, same database, same users): enforce-mode API: ADMIN 27/28 GET 200, MANAGER 15/28, EMPLOYEE 10/28, DRIVER 3/28; write probes authorised 8, 2, 0, 0 of 8. ENV=staging API with RBAC unset: EVERY role (including DRIVER) gets 27/28 GET 200 with salary/bank/national-id fields in 4 responses and passes authorization on 8/8 write probes (400 body validation instead of 403). Real data was returned (sensitive-field responses 4 per role), not only status codes.

## Re-runnable scripts / evidence
- `tests/runtime/ph17_sec013.mjs`
- `tests/runtime/ph17_r5_rbac_matrix.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
