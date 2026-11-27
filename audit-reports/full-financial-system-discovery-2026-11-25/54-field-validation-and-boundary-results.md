# 54 — Field validation and boundary results (RUN-003)

Scope and honesty note: the 627-field inventory ([07](07-field-level-inventory.csv)) was **not** exercised field by field in the UI. This file consolidates the validation/boundary checks that were actually executed (API-DIAGNOSTIC scripts against the isolated stack, plus UI checks where noted). Rows are evidence-backed by the scripts and JSON under `evidence/*/RUN-003/PH-xx/`. "Rule" = expected behaviour of a correct system. Status: PASS = rejected/handled correctly, FAIL = accepted or wrong status/raw error, OBS = observation.

| Area | Field / rule | Input | Result | Status | Finding / evidence |
|---|---|---|---|---|---|
| Journals | balance | debit 100 / credit 90 | 400 "journal entry is not balanced" | PASS | AUD-ACC-FLOW JE-01 |
| Journals | line sides | debit and credit both set | 400 | PASS | JE-02 |
| Journals | negative amounts | -100 | 400 | PASS | JE-03 |
| Journals | header account | post to header | 400 | PASS | JE-15 |
| Journals | date range | 2019-01-15, 2030-06-01 | posted, periods auto-created | FAIL | F-013 |
| Journals | edit/re-post | posted journal | 400 | PASS | JE-09/10 |
| Journals | CSV import | unbalanced / unknown account / bad date / wrong header | HTTP 200 success:true with failed_entries | FAIL | F-063 |
| Accounts | duplicate code | 3150 twice | 400 with raw pq text | FAIL | F-031 |
| Accounts | account_type | "nonsense" | 400 clear | PASS | ACC-03 |
| Customer invoice | tax vs discount | 10% discount, 14% VAT | tax on gross | FAIL | F-007 |
| Customer receipt | amount > balance | 5,000 on 3,000 | 400; 5 concurrent 1,000 all accepted | FAIL (concurrency) | F-006 |
| Supplier invoice | tax vs discount | same | tax on gross | FAIL | F-041 |
| Supplier payment | amount > balance; concurrent | 3,400 on 2,400; 4x800 | 400; 3 accepted + 1 rejected | PASS | AUD-SUP-FLOW |
| Payments | Idempotency-Key | same key twice | 2 records | FAIL | F-052 |
| Orders | quantity | 0, -3 | 400 clear | PASS | AUD-ORDER-TYPES OT-05/06 |
| Orders | quantity above stock | 100,000 | 400 insufficient inventory | PASS | OT-07 |
| Orders | payment method | "bitcoin" | 400 | PASS | OT-11 |
| Orders | expected_total mismatch, unknown product | — | 500 "Failed to create order" | FAIL | F-050 |
| Orders | customer_name / notes | script and SQL text | stored verbatim; rendered escaped in orders page and KDS | PASS (render) | AUD-XSS-001 |
| Orders | paid order cancel/refund | all endpoints | rejected (no workflow) | FAIL (capability) | F-049 |
| Products | price | -5, 0, 1e12, 10.12345 | 201, 201, 500 overflow, 201 | FAIL | F-062 |
| Products | cost / stock / discount | -1; cost 50 on price 10; stock -4; discount 15 on 10 | all 201 | FAIL | F-062 |
| Products | name | empty; 10,000 chars | 201; 500 raw pq text | FAIL | F-062 / F-031 |
| Products | SKU unique | duplicate | 500 "already exists" (3rd attempt) | FAIL (status) | F-062 |
| Products | negative-price sale | price -5 | 400 "invalid product price" | PASS | AUD-MD-NEGPRICE |
| Products | zero price paid | price 0 | outbox event cannot post | FAIL | F-061 |
| Categories | empty / duplicate name | "" ; same name | 400 ; 201 duplicate | PASS / FAIL | F-062 |
| Suppliers / customers | email, phone, duplicate, credit limit | invalid | 201 | FAIL | F-062 |
| Drivers | phone | "abc" | 201 | FAIL | F-058 |
| Attendance | clock_in/out strings | "09:00"/"19:30" | 201 with NULL times and 0 hours | FAIL | F-053 |
| Attendance | duplicate day, unknown employee | — | 500 raw pq text | FAIL | F-031 |
| Leave | end before start | — | 500 raw check text | FAIL | F-031 |
| Leave | total_days vs dates; over allowance | 40 for 3 days; 30 vs 21 | accepted | FAIL | F-054 |
| Salary | period order, negative base, duplicate | — | 500 raw constraint text | FAIL | F-031 |
| Salary | arithmetic | base 5000 + OT 10x50 + bonus 300 + allowances 200 - 750 | gross 6,000 net 5,250 | PASS | AUD-SAL-001 |
| Contracts / resignations | date order | end before start | 201 | FAIL | F-056 |
| Contracts | duplicate number | — | 400 raw pq text | FAIL | F-031 |
| Exchange rates | rate, same pair | -1; USD/USD | 400; 200 | PASS / FAIL | F-048 |
| Recurring journals | repeat run same period | 2 runs | 2 journals | FAIL | F-047 |
| Super Admin | extend trial | -30; 100000 | 400; 200 (year 2300) | PASS / FAIL | F-055 |
| Themes | colors, custom values | "not-a-color", url(javascript:) | accepted, served publicly | FAIL | F-059 |
| Registration | subdomain availability race | fast click | stale error | FAIL (UX) | F-038 |
| Settings | contact info format | invalid values | accepted | FAIL | F-037 |
| Auth | tampered / empty cookie | — | 401 | PASS | AUD-SEC-BATTERY |
| Auth | login without/foreign Origin | — | 403 | PASS | SEC-012 |
| Auth | rate limit | 30 failed logins | 14 x 401 then 429 | PASS | SEC-012c |
| Search filters | quote / UNION / DROP payloads | — | empty 200, tables intact | PASS | SEC-SQLI |
| Cross-tenant | IDs of another tenant | read/write | 404 / 400, no change | PASS | SEC-IDOR |

## Not executed (explicit)
Per-field UI validation for the 627 inventoried fields; Arabic-input and RTL-specific validation; file-upload fields (images, translations upload); phone/email normalisation rules across all forms; date-picker boundaries in the UI; currency/decimal formatting in Arabic locale. These remain PLANNED in [51](51-phased-audit-task-register.csv) (T-P04-002).
