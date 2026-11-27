# F-080 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome tenant 4: open the customer form; the Receivable account dropdown lists 1000 - Assets, 1100 - Current Assets, Petty Cash, Cash on Hand, Bank accounts and others; choose the first account, save; then create a customer invoice for that customer (the form matrix did this for customers 6-13 and suppliers 6, 7)

## Expected result
Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account

## Actual result (measured)
Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer returned 400 {"error":"account 63 is a header account and cannot be posted to"}, so the customer can never be invoiced until the account is edited

## Re-runnable scripts / evidence
- `tests/runtime/ph09_r5_invoice_form.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
