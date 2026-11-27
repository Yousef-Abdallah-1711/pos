# 17 — Cross-module integration results

| Chain | Result | Evidence |
|---|---|---|
| POS cash order → order paid → outbox → sales+COGS journals → GL → TB/P&L | PASS when COA exists; FAIL on fresh tenant (F-001); amounts follow server 10% tax (F-004) | runtime/06, 24 |
| POS card/wallet → paid → ledger | FAIL (no capture, F-005) | runtime/10 |
| Product stock → inventory movement → COGS | Partial: stock decrements at order creation (also for unpaid orders); COGS at payment; inventory report ≠ GL (F-011) | 24 |
| Supplier invoice → payable → payment → AP | PASS (explicit account 1130); supplier balance field inconsistent (F-021) | TX-004/005 |
| Customer invoice → AR → receipt | PASS for sequential; FAIL for concurrent (F-006); tax mismatch (F-007) | TX-011..013 |
| Employee → salary → payroll journal | UI employee creation FAIL (F-014); salary→journal PASS but cash-basis, unapproved (F-015) | TX-007/008 |
| Attendance/overtime → payroll | NOT EXECUTED (HR attendance not exercised; static: not linked) | — |
| Journals → reports (TB, BS, CF, P&L, VAT, inventory) | TB PASS; BS/CF/VAT/inventory FAIL (F-002/003/010/011) | 43 |
| Orders → KPI/Sales/Profitability dashboards | FAIL consistency (F-012) | 43 |
| Storefront order → ledger | BLOCKED (F-022) | runtime/18 |
| Fiscal period → posting | partial (F-013) | runtime/14 |
