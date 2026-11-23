# Transaction-to-ledger trace

## Synthetic paid order

Executable acceptance: `TestFinancialAcceptanceLiveAPIOrderToLedger` in `financial-live-api-tests.jsonl`.

| Layer | Expected | Observed | Difference |
| --- | ---: | ---: | ---: |
| Source subtotal | 5,000 cents | 5,000 cents | 0 |
| Source tax | 500 cents | 500 cents | 0 |
| Discount / delivery | 0 / 0 cents | 0 / 0 cents | 0 |
| Order total | 5,500 cents | 5,500 cents | 0 |
| Paid event | One | One | 0 duplicate effects |
| Cash GL | 5,500 cents | 5,500 cents | 0 |
| Revenue GL | 5,000 cents | 5,000 cents | 0 |
| Tax payable GL | 500 cents | 500 cents | 0 |
| Trial Balance debits / credits | 5,500 / 5,500 cents | 5,500 / 5,500 cents | 0 |

The repeated callback returned HTTP 200 and did not create a second event or journal. The raw log associates this with synthetic tenant 23, restaurant 23, order 3, journal 29, and event key `paid-order:23:23:3`. These identifiers are local test fixtures only.

The authenticated reports/export check showed P&L source/GL revenue of 5,000 cents, sales export and CSV gross of 5,500 cents, and PDF report gross of 5,500 cents for one order. This indicates the tested export values agreed with their respective order totals; it does not make gross receipts equal revenue because tax is separate.

## Other financial lifecycles

Accountant journal: 10,000-cent debit and credit, two GL rows; cashier post denied (403), accountant post succeeded (200), duplicate post rejected (400), no partial write. AR: invoice 10,000, receipt 4,000, open amount 6,000 cents with source/subledger/GL/aging differences all zero. Payroll: 150,000-cent gross/net, matching expense and cash GL; failed post rolled back and retry succeeded. Supplier invoice: 20,000-cent subtotal + 2,000-cent tax = 22,000 cents, balanced invoice and payment journals. Concurrent payment test: four successful requests recorded 20,000 cents against a 22,000-cent payable; 2,000 cents remained.
