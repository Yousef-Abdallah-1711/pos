# F-049 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Pay an order, then try DELETE /admin/orders/{id}, PUT payment-status refunded/voided/failed/pending, PUT status cancelled, item cancel

## Expected result
Controlled refund/void producing reversing journals, stock return and VAT reversal

## Actual result (measured)
DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refunded/failed/pending from paid -> 400 invalid transition; voided -> "invalid payment status"; status=cancelled -> 500; paid revenue/COGS stay posted; cash 307,330.01 unchanged | RUN-003: paid order cannot be cancelled/refunded/voided through any endpoint; explicit guard message (REPRODUCED). 

## Re-runnable scripts / evidence
- `tests/runtime/ph19_paid_cancel.mjs`
- `tests/runtime/ph19_regs2.mjs | evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
