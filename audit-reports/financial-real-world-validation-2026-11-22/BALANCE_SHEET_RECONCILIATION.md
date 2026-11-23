# Balance Sheet reconciliation

**Result: FAIL / UNRECONCILED; no adjustment made.**

The synthetic paid order produced a tested report with assets of 5,500 cents, liabilities of 500 cents, and equity of zero. The equation residual is:

`5,500 − 500 − 0 = 5,000 cents ($50.00)`

The report identified itself as unreconciled. This is consistent with revenue being present while current earnings/equity presentation is not included in the tested report; that diagnosis does not authorize a posting. No balancing journal or equity adjustment was fabricated.

The browser suite captured English desktop and Arabic mobile report screenshots; both represent this unreconciled synthetic report and are diagnostic evidence, not acceptance evidence. See `acceptance-en-balance-sheet.png`, `acceptance-ar-mobile-balance-sheet.png`, and the exact API assertion in `financial-live-api-tests.jsonl`.

Current earnings and opening-balance treatment remain policy-dependent. EARN-01 and OPEN-01 are PENDING. Balance Sheet acceptance is therefore NOT PASSED.
