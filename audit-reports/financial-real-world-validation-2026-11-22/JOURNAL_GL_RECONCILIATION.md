# Journal and GL reconciliation

**Synthetic result:** PASS for exercised posted journals; no claim for every journal in the fixture database.

Evidence: authenticated live API acceptance log `financial-live-api-tests.jsonl`, accountant journal and payroll/supplier scenarios.

- Order journal: cash debit 5,500 cents; revenue credit 5,000; tax payable credit 500; three GL rows. Trial balance debit/credit totals both 5,500 cents.
- Manual accountant journal: 10,000-cent debit and credit, two GL rows; cashier posting denied, accountant posting accepted; duplicate posting rejected; validation failures produced no partial writes.
- Payroll journal: two GL rows and 150,000-cent expense/cash effects matching source. A forced failed post rolled back; retry succeeded once.
- Supplier invoice: expense debit 20,000 and input-VAT debit 2,000 against AP credit 22,000. Payment debited AP 22,000 and credited cash 22,000. The separate concurrent-payments scenario left a 2,000-cent payable.

No unexplained duplicate was observed in these exercised scenarios. The broad test database contains unrelated fixtures, so aggregate database totals are not offered as a full reconciliation. Reversal/void paths and orphan detection across all synthetic records were not established by this run.
