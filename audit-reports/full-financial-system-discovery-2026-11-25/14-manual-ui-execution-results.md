# 14 — Manual UI execution results

Browser: Playwright Chromium (desktop 1440×900), isolated stack, English UI plus Arabic spot-check. Scripts: `tests/runtime/*.mjs`. Evidence: `evidence/screenshots`, `evidence/network`, `evidence/reports`.

## Executed UI workflows (logical forms: 13)
| # | Workflow | Tenant | Result |
|---|---|---|---|
| 1 | Tenant registration (3×) | A,B,C | PASS (201); onboarding gap F-001 |
| 2 | Product create | A,C | PASS (201); POS shows `$` (F-023) |
| 3 | POS cash order | A,C | PASS create+paid; UI total ≠ stored total (F-004) |
| 4 | POS card / wallet order | C | FAIL: remain pending (F-005) |
| 5 | Manual journal create + post (capital, rent) | C | PASS |
| 6 | Supplier create | C | PASS (requires choosing payable account) |
| 7 | Supplier invoice create + approve | C | PASS with explicit account 1130 (default account unverified at runtime) |
| 8 | Supplier payment | C | PASS (posts JE) |
| 9 | Employee create | C | **FAIL HTTP 500** (F-014) — API fixture used |
| 10 | Salary create + Mark Paid | C | PASS functionally; no approval/accrual (F-015); list filter hides record (F-016) |
| 11 | Customer create; customer invoice; receipt | C | invoice FAIL tax/total mismatch (F-007); receipt PASS |
| 12 | User create (Employee role) | C | PASS |
| 13 | Report pages: P&L, BS, CF, VAT, inventory valuation, AP/AR aging, advanced, overview, sales, profitability, inventory, HR, TB, ledger, orders | C | loaded; see 43 |

## Page coverage
≈45 distinct page loads of 119 inventoried routes (≈38%), of which ≈15 had functional interaction; the rest were loads only. Not visited: super-admin app (0 of 17), most HR sub-modules (attendance, leaves, contracts, insurance, performance, training, resignations), KDS, tables/drivers/categories, theme builder, settings other than users, RBAC role editor, accounting advanced console, bank accounts, chart-of-accounts create, budgets, recurring journals, FX rates, tax periods, imports/exports.

## Arabic/responsive
Arabic (`/ar/...`) on 6 pages: `dir=rtl lang=ar`, 0 page errors (`evidence/runtime/19`). Tablet/mobile responsive checks: NOT EXECUTED.

## Not executed
Buttons/filters/search/pagination/sorting/bulk actions/import-export as individual tests; variants and modifiers in POS; table/takeaway/delivery selection; kitchen status flow; order cancellation (item/order); refunds; storefront ordering (blocked, F-022). All corresponding rows of `34-complete-manual-ui-test-plan.csv` remain NOT EXECUTED.
