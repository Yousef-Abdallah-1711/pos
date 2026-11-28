# F-023 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/lib/rtl.ts:81-91 formatCurrencyWithLocale defaults currency to USD (Intl.NumberFormat style currency), so any caller that omits the tenant currency prints $; the storefront passes EGP and shows "EGP" (AUD-STORE-004) while the POS shows $


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/lib/rtl.ts`:81-91.
- Propagation: affected modules POS; related findings none recorded; observed consequence: Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE).  RUN-005: the POS page shows $12.00, $10.00, $100.00, $10,000.00 while the dashboard home shows "EGP 1
- Cause versus symptom: the symptom is what users or reports show (Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE).  R); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-POS-001, EX-RUN003-AUD-POS-002, EX-RUN005-OBSERVED-REPRO; scripts: tests/runtime/ph20_evidence_gaps.mjs, tests/runtime/ph07_r5_pos_currency.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/lib/rtl.ts | 81-91 | frontend/apps/dashboard/src/lib/rtl.ts | exact | yes |

