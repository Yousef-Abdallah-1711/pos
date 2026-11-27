# F-070 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Tenant 4: compare /kds/analytics busiest_hours with SQL bucketed by UTC and by Africa/Cairo for the same kds_queue rows

## Expected result
Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3)

## Actual result (measured)
API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
