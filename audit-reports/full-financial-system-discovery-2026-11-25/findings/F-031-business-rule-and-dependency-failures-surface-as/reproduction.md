# F-031 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Cancel a paid order; receipt during accounting-role outage; DB outage

## Expected result
409/503 with generic message

## Actual result (measured)
Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
