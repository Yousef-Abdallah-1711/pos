# F-072 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Tenant 4 bank account 2 (ledger 12,000.00, GL 0): create a session with statement balance 777,777.77, add one line referencing a real journal that has no bank line with book_amount 777,777.77, complete the session (API-DIAGNOSTIC)

## Expected result
Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts agree with the statement and the ledger

## Actual result (measured)
Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123,456.78 or 500.00 is also stamped completed_at although its status stays in_progress, and the complete call returns HTTP 200

## Re-runnable scripts / evidence
- `tests/runtime/ph11_r4_bank_recon.mjs`
- `tests/runtime/ph11_r4_bank_recon2.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
