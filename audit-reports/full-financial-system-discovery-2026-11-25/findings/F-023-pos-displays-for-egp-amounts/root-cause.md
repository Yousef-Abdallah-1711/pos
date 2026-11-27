# F-023 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/lib/rtl.ts:81-91 formatCurrencyWithLocale defaults currency to USD (Intl.NumberFormat style currency), so any caller that omits the tenant currency prints $; the storefront passes EGP and shows "EGP" (AUD-STORE-004) while the POS shows $


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/lib/rtl.ts | 81-91 | frontend/apps/dashboard/src/lib/rtl.ts | exact | yes |

