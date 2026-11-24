# Reproduction

## Original source transaction

In disposable PostgreSQL, tenant 10 / restaurant 12 had paid order 2 (`ORD-2026-9221af93-7b0e-41e0-a710-6aeab26cfd58`): subtotal 50.00, tax 5.00, discount 0.00, delivery 0.00, total 55.00. Posted journal 25 references that order and produces cash debit 55.00, VAT payable credit 5.00, and revenue credit 50.00. The payment status was `paid`.

The authenticated order-to-ledger API test had already passed in the preceding validation run. Independent PostgreSQL queries in this run confirmed the source order, journal, GL amounts, and report figures.

## Reproduction totals

For the same journal, the independent SQL aggregate and authenticated report API on tenant 12 / restaurant 14 (synthetic report fixture, journal 43) both returned:

```json
{
  "balance_sheet": {
    "total_assets": 55.0,
    "total_liabilities": 5.0,
    "total_equity": 0.0,
    "equation_difference": 50.0,
    "integrity_status": "UNRECONCILED",
    "equation_balanced": false
  },
  "profit_loss": {"total_revenue": 50.0, "total_expenses": 0.0, "net_profit": 50.0},
  "trial_balance": {"total_debits": 55.0, "total_credits": 55.0, "is_balanced": true}
}
```

Amounts above are dollars in the database/API currency unit; equivalently: assets 5,500 cents, liabilities 500 cents, equity 0, residual 5,000 cents, revenue/net income 5,000 cents, and Trial Balance 5,500 / 5,500 cents.

## Cache defect reproduction

Before the fix, the browser fixture first requested a zero-activity Balance Sheet and cached `RECONCILED`. After posting synthetic journal 43 (cash debit 55.00, VAT payable credit 5.00, revenue credit 50.00), the same report request still returned zero totals and `RECONCILED`; PostgreSQL contained the new GL lines. The new handler regression test reproduced the stale response before the code change (second response remained all zero). After the fix, it recalculates and returns 55.00 / 5.00 / 0.00 with a 50.00 residual and `UNRECONCILED`.
