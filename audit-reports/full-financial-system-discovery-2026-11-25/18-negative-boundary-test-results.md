# 18 — Negative and boundary test results

Journal create (tenant C, owner; evidence `runtime/13`): unbalanced → 400 PASS; negative amounts → 400 PASS; zero amounts → 400 PASS; header account → 400 PASS; other-tenant account → 400 PASS; single line → 400 PASS; 1e15 → 400 but leaks `pq: numeric field overflow` (F-025); amounts 0.005 and 10.005 accepted and stored rounded to 0.01 / 10.01 (silent rounding; policy ROUND-01); future date (2027-01-01) and backdated (2020-01-01) accepted as drafts (no date policy).
Payments (runtime/15): customer overpay, zero, negative → 400 PASS; supplier overpay → 400 PASS (sequential).
Fiscal period (runtime/14): post blocked PASS; draft creation allowed; invoice creation blocked (it posts); **void allowed (F-013)**.
Employee form with default values → HTTP 500 (F-014).
Not executed: negative stock/oversell, unicode/XSS/SQL strings in names, duplicate references, inactive account posting, inactive supplier/customer, 3-decimal prices in POS, extreme quantities, all field-level boundaries from the 627-row field inventory (rows remain NOT EXECUTED).
