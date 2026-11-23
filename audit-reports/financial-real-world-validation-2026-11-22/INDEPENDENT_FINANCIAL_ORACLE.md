# Independent financial oracle

The acceptance checks calculate expected amounts from synthetic source transaction inputs using integer minor units (cents), then compare them with persisted/API/report values. These expected values are not obtained by calling the report implementation under test. Raw assertions and observed values are recorded in `financial-live-api-tests.jsonl`; UI cross-check assertions are in the Playwright suite and its text/JUnit output.

## Independent expected values

- Paid order: subtotal 5,000 + tax 500 − discount 0 + delivery 0 = 5,500 cents.
- Expected posting: debit cash 5,500; credit revenue 5,000 and tax payable 500.
- AR: invoice 10,000 − payment 4,000 = open AR 6,000 cents.
- Payroll path: gross/net 150,000 cents; expected expense and cash effects each 150,000 cents.
- Supplier invoice: subtotal 20,000 + input tax 2,000 = AP 22,000 cents; concurrent settlement of 20,000 leaves 2,000 cents.
- Accountant manual journal: debit = credit = 10,000 cents.

The tested order and supporting scenarios matched source, journal/GL, trial balance, and relevant API/export assertions exactly. The Balance Sheet assertion did not pass the accounting equation: 5,500 assets − 500 liabilities − 0 equity = 5,000-cent residual. The result is reported, not concealed by a tolerance or fabricated entry. Cash-flow account classification, opening balances, current-earnings presentation, and close behavior depend on unresolved accounting policy and are not PASS results.

## Scope limitations

The executed oracle is scenario-specific. It does not establish independent results for refunds, COGS/inventory valuation, bank reconciliation, all payment methods, full cash operations, period locks, or original operating transactions. Full decimal formatting is not exercised; all comparisons reported here are exact integer cents.
