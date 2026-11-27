# F-039 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
POST /accounting/journal-entries with attachments = {"evil":true}, "<script>alert(1)</script>", [{"url":"javascript:alert(1)"}], and a 1 MB blob

## Expected result
400/422 for non-array, unsafe URL scheme or oversize attachments

## Actual result (measured)
All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE

## Re-runnable scripts / evidence
- `tests/runtime/ph19_regressions.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
