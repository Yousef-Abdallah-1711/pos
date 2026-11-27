# F-005 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:302-308 calls updateOrderPaymentStatus(paid) only when paymentMethod === cash; card and wallet orders stay payment_status pending, and no other flow settles them (provider settlement belongs to subscription code)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | 302-308 | frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exact | yes |

