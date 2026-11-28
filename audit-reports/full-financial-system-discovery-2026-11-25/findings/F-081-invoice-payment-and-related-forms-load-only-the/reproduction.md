# F-081 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome tenant 4 (91 active customers): open New customer invoice and count the options in the customer dropdown

## Expected result
Every active customer can be selected (search or paging)

## Actual result (measured)
The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20

## Re-runnable scripts / evidence
- `tests/runtime/ph09_r5_invoice_form.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
