# F-007 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Invoice 5000 less 10% discount VAT 14%

## Expected result
Tax 630 total 5130 (as UI shows)

## Actual result (measured)
Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE). 

## Re-runnable scripts / evidence
- `tests/runtime/ph09_custinv_draft_headed.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
