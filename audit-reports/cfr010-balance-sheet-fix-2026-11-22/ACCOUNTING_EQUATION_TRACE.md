# Accounting Equation Trace

Synthetic paid sale, tenant 10 / restaurant 12 / order 2 / journal 25:

| Stage | Debit | Credit | Evidence |
|---|---:|---:|---|
| Source order | — | — | Subtotal 5,000¢; tax 500¢; discount/delivery 0¢; total 5,500¢ |
| Posted journal | 5,500¢ cash | 500¢ VAT payable; 5,000¢ revenue | Journal 25, posted; reference type `order`, reference id 2 |
| General Ledger | 5,500¢ | 5,500¢ | Cash, VAT payable, revenue GL lines; scoped by tenant and restaurant |
| Trial Balance | 5,500¢ | 5,500¢ | Balanced |
| Profit & Loss | Revenue 5,000¢; expenses 0¢ | Net income 5,000¢ | Revenue is nominal activity in the period |
| Balance Sheet | Assets 5,500¢ | Liabilities 500¢; equity 0¢ | Residual = 5,000¢; status `UNRECONCILED` |

Equation: `5,500¢ assets - 500¢ liabilities - 0¢ equity = 5,000¢ residual`.

This trace establishes the synthetic case, not the original operating records and not an approved presentation of current earnings.
