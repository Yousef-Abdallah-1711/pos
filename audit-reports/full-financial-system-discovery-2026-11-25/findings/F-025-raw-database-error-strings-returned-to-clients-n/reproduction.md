# F-025 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Submit 1e15 amount

## Expected result
Sanitised message

## Actual result (measured)
pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user receives raw text such as pq: value too long for type character varying(3) (REPRODUCED; the filler used the same long synthetic value in every text field, so the inputs are artificial but the exposure is real).  RUN-005 UI: deleting a customer that has an invoice returns 400 with the raw text "pq: update or delete on table customers_accounting violates foreign key constraint invoices_customer_customer_id_fkey".

## Re-runnable scripts / evidence
- `tests/runtime/ph04_r5_destructive_a.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
