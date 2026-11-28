# F-023 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Open POS

## Expected result
EGP symbol

## Actual result (measured)
Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE).  RUN-005: the POS page shows $12.00, $10.00, $100.00, $10,000.00 while the dashboard home shows "EGP 110.00" for the same tenant (currency EGP).

## Re-runnable scripts / evidence
- `tests/runtime/ph20_evidence_gaps.mjs`
- `tests/runtime/ph07_r5_pos_currency.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
