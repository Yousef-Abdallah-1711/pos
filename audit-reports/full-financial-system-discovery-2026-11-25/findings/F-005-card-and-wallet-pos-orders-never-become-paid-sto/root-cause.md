# F-005 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:302-308 calls updateOrderPaymentStatus(paid) only when paymentMethod === cash; card and wallet orders stay payment_status pending, and no other flow settles them (provider settlement belongs to subscription code)


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`:302-308.
- Propagation: affected modules POS; Orders; Accounting; Inventory; related findings none recorded; observed consequence: Order stays pending; no outbox; no journal; stock decremented (8 unpaid orders = 88000 not in ledger)
- Cause versus symptom: the symptom is what users or reports show (Order stays pending; no outbox; no journal; stock decremented (8 unpaid orders = 88000 not in ledger)); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | 302-308 | frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exact | yes |

