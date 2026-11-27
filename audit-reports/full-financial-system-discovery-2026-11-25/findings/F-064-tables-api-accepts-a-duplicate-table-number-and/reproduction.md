# F-064 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /tables {table_number:1,name:dup,capacity:4} when table 1 exists; POST /tables with a 1,800-character name

## Expected result
409 duplicate table number; 400 name length

## Actual result (measured)
duplicate -> 201 (two tables numbered 1, ids 1 and 2); long name -> 500 "Failed to create table"; capacity (1-50), negative number and delete-with-active-orders (409) are validated correctly | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph16_tables.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
