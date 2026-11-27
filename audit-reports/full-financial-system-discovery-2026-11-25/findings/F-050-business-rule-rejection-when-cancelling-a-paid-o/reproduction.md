# F-050 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
PUT /admin/orders/{paid id}/status {status:cancelled}

## Expected result
409/400 with the same message as DELETE

## Actual result (measured)
500 "Failed to update order status" (DELETE returns 400 with the explicit reason) | RUN-003: PUT status cancelled on paid order -> 500.  | RUN-003: PUT status cancelled on paid order -> 500. Same class: unsigned Paymob webhook -> 500 and refund of unknown payment -> 500 (AUD-SAAS-MATRIX PAY-02/PAY-03). 

## Re-runnable scripts / evidence
- `tests/runtime/ph19_paid_cancel.mjs`
- `tests/runtime/ph19_regs2.mjs`
- `tests/runtime/ph07_order_types.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
