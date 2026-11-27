# F-018 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Stop Redis then call login/session/orders

## Expected result
Fast 503 or degraded mode; health reflects dependency

## Actual result (measured)
Requests time out (no response); health 200 | RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 11 ms; after Redis returns the previous session is rejected (401) (REPRODUCED). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
