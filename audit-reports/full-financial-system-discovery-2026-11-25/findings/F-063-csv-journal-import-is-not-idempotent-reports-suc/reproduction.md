# F-063 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /accounting/csv/journal-entries/import with a balanced 2-line file twice; with an unbalanced / unknown-account / bad-date / wrong-header file; with 6,000 journals (617 KB)

## Expected result
Duplicate upload rejected or deduplicated; failed rows reported as an error status; large imports asynchronous or chunked with a job id

## Actual result (measured)
identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:true with failed_entries 1 and created_entries 0; a 6,000-journal import ran ~122 s, the client socket was closed (no HTTP response) while the server completed all 6,000 journals, so a client retry duplicates them (17,000 draft journals created in the audit tests); formula characters are stored verbatim (export escapes them) | RUN-003 import battery: re-upload duplicates, all-failed imports answer 200 success, 6,000-journal import completed server-side after the client socket closed (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph11_import.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
