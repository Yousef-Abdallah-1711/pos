# 04 - Financial model and policies (from current source)

Baseline: git SHA 133260db66b1143c254d79476daa7650b41a33f0. Method: static reading of `backend/internal/**`, `backend/migrations/**`, selected `frontend/` files, plus one SELECT-only catalog query on `pos_saas`. Nothing here was run end-to-end. "Runtime-observed" is `NOT YET OBSERVED` except for the one catalog query noted in section 1.

Policy decisions listed as "Unresolved" are those the repository docs themselves mark PENDING (OPEN-01, EARN-01, CLOSE-01, TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01, CF-01; source: `PHASE4_ACCOUNTING_POLICY_DECISION_PACKAGE.md`). No approval record exists in the repo; none is invented here.

## 1. Monetary representation

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved policy decision |
|---|---|---|---|---|
| Go domain types | CFR-007 says exact/minor-unit contract is desired; not approved (`FINANCIAL-MODULE-...PLAN.md` CFR-007) | `float64` for money everywhere: `backend/internal/domain/accounting.go:38,83-84,109-110,121-127,161-163,185-186,225-227`; `backend/internal/domain/order.go:39-43,79-81,94`. Grep count: 125 `float64` money-named lines in `backend/internal/domain/*.go`. No decimal library in `backend/go.mod` (shopspring absent). | Catalog query (SELECT on information_schema, pos_saas): money-named numeric columns are `numeric(15,2)` x47, `numeric(10,2)` x24, `numeric(5,2)` x9, `numeric(12,2)` x4, `numeric(8,2)` x2, `numeric(14,2)` x1, `numeric(18,8)` x1. No `double precision` or `money` columns matched the name filter. Other behaviour NOT YET OBSERVED. | ROUND-01: currency exponent, representation (integer minor units vs decimal strings), rounding mode |
| Integer-cents islands | n/a | Order pricing computes in int64 cents then converts back to float64: `backend/internal/usecase/order_pricing.go:31,41,46-53`. Auto-posting validation converts to cents for equality: `backend/internal/usecase/auto_posting_validation.go:70-95`. `math/big` only in read-only diagnostics: `backend/internal/repository/opening_balance_diagnostics.go:187`, `backend/internal/usecase/current_earnings_diagnostic.go:193,221`. | NOT YET OBSERVED | Same as above |
| Reports | n/a | Report totals are float64 sums of NUMERIC values scanned into float64: `backend/internal/usecase/financial_report_usecase.go:276-353`. Inventory valuation casts to `DOUBLE PRECISION` in SQL: `:488-510`. | NOT YET OBSERVED | Same |
| Frontend | n/a | Display uses JS `Number` and `toFixed(2)`, e.g. `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx:378-417`. Frontend integrity helper test exists: `frontend/apps/dashboard/src/lib/financial-integrity.test.ts` (content not reviewed). | NOT YET OBSERVED | Same |
| Currency | FIN-020: multi-currency policy not established | Order quote hard-codes `Currency: "EGP"`: `backend/internal/usecase/order_pricing.go:99`. `Currency` fields exist in `domain/accounting.go:39,62,187`. No FX code found (UNVERIFIED by exhaustive search). | NOT YET OBSERVED | Supported currencies, functional currency, FX, consolidation |

## 2. Rounding rules

| Layer | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| Backend order tax | Comment: "Round the established 10% tax to cents, as PostgreSQL does" | `tax := (subtotal + 5) / 10` on integer cents, i.e. half-up: `backend/internal/usecase/order_pricing.go:46` | NOT YET OBSERVED | ROUND-01 mode (half-up vs half-even), residual allocation |
| Backend generic | n/a | `math.Round(value*100)/100` (half away from zero on float): `backend/internal/usecase/auto_posting_engine.go:645-646`, `backend/internal/handler/http/ledger_handler.go:166`, `backend/internal/repository/accounting_phase3_helpers.go:15` | NOT YET OBSERVED | Same |
| Backend report equation | n/a | `round2(v) = float64(int(v*100+0.5))/100` (truncating int cast, wrong for negative values): `backend/internal/usecase/financial_report_usecase.go:528-530`; used for the equation check at `:343`. Finding: asymmetric for negatives; effect not tested. | NOT YET OBSERVED | Tolerance policy (docs: none approved) |
| Auto-posting | CFR-001 says no tolerance | Sub-cent amounts are rejected, not rounded; imbalance in cents rejected: `auto_posting_validation.go:70-95` | NOT YET OBSERVED | ROUND-01 |
| Order total check | n/a | Client `ExpectedTotal` accepted within 0.004: `backend/internal/usecase/order_usecase.go:163` | NOT YET OBSERVED | Tolerance policy |
| DB | n/a | NUMERIC scale 2 for money (see section 1 observation); invoice line trigger rounds with `ROUND(...,2)`: `backend/migrations/141_create_invoices_customer.sql:80-82` | Scale observed via catalog query only | ROUND-01 |
| Frontend | n/a | `toFixed(2)` at display (see section 1) | NOT YET OBSERVED | Same |

## 3. VAT / tax handling

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| Order tax | CFR-008: expected tenant/jurisdiction-configured tax | Hard-coded 10% of subtotal for every tenant, not configurable: `backend/internal/usecase/order_pricing.go:11,46`; delivery fee forced to 0 at `:52`. | NOT YET OBSERVED | TAX-01: jurisdiction, rate, inclusive/exclusive, exemptions |
| Order discount | n/a | `Order.DiscountAmount` is subtracted after tax: `order_pricing.go:53`. A grep for assignments to `DiscountAmount =` in `backend/internal` found none, so order-level discounts appear to be always 0 unless set at the repository/API layer (UNVERIFIED). | NOT YET OBSERVED | TAX-01 discount ordering |
| Customer invoice tax | CFR-008: tax computed before discount | Per-line `tax_rate` user supplied; trigger computes tax on gross `qty*unit_price` independent of discount: `backend/migrations/141_create_invoices_customer.sql:80-82`. Function defined only in migration 141 (grep). Dashboard form accepts a per-line rate: `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:39,421`. | NOT YET OBSERVED | TAX-01 |
| Supplier invoice tax | n/a | `backend/internal/repository/invoice_supplier_repo.go:449-476` (not re-read in this pass; UNVERIFIED) | NOT YET OBSERVED | TAX-01 |
| GL treatment | n/a | Sales posting credits tax to account `accounts["tax"]` else `2120`: `backend/internal/usecase/auto_posting_engine.go:71,115-120`. Tax periods table exists: `backend/migrations/146_create_tax_periods.sql` (usage not reviewed). | NOT YET OBSERVED | TAX-01 |

## 4. Accounting basis, revenue and recognition

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| Basis | Docs: posted journal/GL is intended source of truth; no basis approved | Order revenue is posted when `payment_status = paid` (cash-like trigger): `auto_posting_engine.go:54-56`. Entry date is `order.CreatedAt`, not payment date: `:137`. Customer invoices post as `customer_invoice`: `backend/internal/repository/invoice_customer_repo.go:630`; supplier invoice as `supplier_invoice`: `invoice_supplier_repo.go:566`. Mixed accrual (AR/AP invoices) and cash-on-paid (orders) behaviour. | NOT YET OBSERVED | Recognition trigger: capture vs settlement vs fulfilment (REFUND-01 / CFR-005); entry-date policy |
| Sales journal shape | n/a | Dr collection account (total); Dr discount if any; Cr sales (subtotal); Cr delivery revenue if any; Cr tax: `auto_posting_engine.go:73-120`. If `DeliveryFee > 0` the shape would not balance against `TotalAmount` (total ignores delivery, `order_pricing.go:52-53`); validation would reject (CFR-001 path). Currently delivery fee forced to 0. | NOT YET OBSERVED | Delivery/service/tip treatment |
| Service charge / tips on orders | n/a | No service-charge or tip field on `domain.Order` (`backend/internal/domain/order.go:39-43`: Subtotal, TaxAmount, DiscountAmount, DeliveryFee, TotalAmount only). Tips exist only on payroll (`salary.Tips`, `auto_posting_engine.go:309,335-343`). | NOT YET OBSERVED | Policy for order tips/service fees |
| Refunds | CFR-005 | Generic `refunded` assertion rejected by `CanTransitionOrderPaymentStatus`: `backend/internal/domain/order.go:219`; no order-level reversal journal found by grep. | NOT YET OBSERVED | REFUND-01 |
| Journal void | CFR-013 | Reversal GL dated `time.Now()`; no closed-period check on void: `backend/internal/repository/journal_entry_repo.go:441,464-472`. | NOT YET OBSERVED | CLOSE-01 |

## 5. COGS and inventory valuation

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| COGS | INV-01 not approved | `PostCOGSEntry`: for each tracked item `unitCost * quantity` using `inventoryMovementRepo.GetProductCost`, account `food_cogs` else `5110`; skips items with cost <= 0; entry dated `order.CreatedAt`: `backend/internal/usecase/auto_posting_engine.go:143-260,233`. Cost is the current stored product cost, not a layer. | NOT YET OBSERVED | Valuation method (FIFO / weighted average / specific), negative stock, returns/waste |
| Valuation report | n/a | `quantity_in_stock * products.cost` at query time; the `asOf` argument only labels the report, it does not filter the data: `financial_report_usecase.go:488-526`. Finding: not a historical as-of valuation. | NOT YET OBSERVED | INV-01 |
| Movement idempotency | CFR-004 | Migration `backend/migrations/184_idempotent_order_sale_inventory_movements.sql` | NOT YET OBSERVED | - |

## 6. Payroll

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| Recognition | PAY-01 unapproved | Posted only when salary is marked paid (cash basis, no accrual): `auto_posting_engine.go:276-282` (rule lookup :284); entry date = payment date else paid_at else `time.Now()`: `:657-668`. | NOT YET OBSERVED | Accrual vs cash; statutory deductions; employer contributions |
| Journal shape | n/a | Dr salary expense (`gross - overtime - tips`, floored at 0) account `5210`; Dr overtime `5220`; Dr tips payable (`2140`/`5210`) for tips; credits bank/cash and deduction lines (credit side not fully reviewed): `auto_posting_engine.go:309-345`. | NOT YET OBSERVED | Gross vs net treatment |
| Atomicity | CFR-009 / HC-14 | `PostPayrollEntryTx` shares the caller's transaction; refuses to commit if the posting rule is disabled: `auto_posting_engine.go:268,288-292`; migration `189_payroll_accounting_transaction_permissions.sql` (column-scoped grants to `pos_accounting`). | NOT YET OBSERVED | - |
| Immutability | CFR-009 | Paid salary update/delete returns 409: `backend/internal/handler/http/hr_salary_handler.go:320,413,451` | NOT YET OBSERVED | Correction/reversal approval |

## 7. Retained earnings, current-period earnings, fiscal close

| Aspect | Documented intent | Implemented behaviour (code ref) | Runtime-observed | Unresolved |
|---|---|---|---|---|
| Balance Sheet scope | EARN-01 pending | Selects only `asset`, `liability`, `equity` account types: `financial_report_usecase.go:304`. Balance per account = signed GL sum through `as_of`: `:283-303`; contra-normal balances contribute negatively: `:355-364`. | NOT YET OBSERVED | EARN-01 current-earnings presentation |
| Current-period earnings | Docs: must not be silently added, no plug | Not computed in the report. A separate read-only diagnostic exists: `backend/internal/usecase/current_earnings_diagnostic.go`. | NOT YET OBSERVED | EARN-01 fiscal-year start, no double counting with close journals |
| Equation check | PR #27 guard | `EquationBalanced = round2(assets) == round2(liabilities+equity)`; statuses `RECONCILED`/`UNRECONCILED`: `financial_report_usecase.go:343-352`. Dashboard blocks PDF export when unreconciled (docs/browser test; export code not reviewed). | NOT YET OBSERVED | Tolerance (none approved) |
| Retained earnings account | n/a | Seeded as `3200` (`equity`/`retained_earnings`): `backend/migrations/156_seed_chart_of_accounts.sql:58`. No closing routine posts to it (grep: only read-only diagnostic reads it, `current_earnings_diagnostic.go:123-140`). | NOT YET OBSERVED | CLOSE-01 |
| Opening balances | OPEN-01 | Not synthesized; read-only diagnostic: `backend/internal/repository/opening_balance_diagnostics.go` | NOT YET OBSERVED | OPEN-01 source evidence and offset account |
| Fiscal periods | CLOSE-01 | `fiscal_periods` table with `status` and dates: `backend/migrations/153_create_fiscal_periods.sql`. Post path rejects posting into a `closed` period: `journal_entry_repo.go:386,813-831`; no matching period row means posting is allowed (`:825-827`). No close/reopen/year-end routine found. Void path has no period check (`:441-520`). | NOT YET OBSERVED | CLOSE-01 states, roles, reopen, adjustment periods |
| Cash flow | CF-01 | Heuristic: net income-like operating, `12%` code prefix investing, equity or `22%` financing: `financial_report_usecase.go:389-391`. Not derived from cash/bank accounts. | NOT YET OBSERVED | CF-01 method and cash-equivalent scope |

## 8. Other behaviours relevant to later phases

| Topic | Implemented behaviour (code ref) | Notes |
|---|---|---|
| Journal integrity | DB triggers for balance and immutability: `backend/migrations/180_enforce_journal_integrity.sql`; normal-balance change guard: `190_prevent_historical_normal_balance_change.sql` | Runtime presence NOT YET OBSERVED |
| Numbering | Advisory-lock then MAX()+1: `backend/internal/repository/financial_document_number.go:11` | Not gapless |
| Entry type | `journal_entries.entry_type VARCHAR(20)` default `manual`, no CHECK (`137_create_journal_entries.sql:19`), widened to 32 by `187_expand_journal_entry_type_to_32.sql` | Auto supplier types unreachable (no callers) |
| Attachments | JSONB, raw passthrough, no validation: `137_create_journal_entries.sql:35`; `handler/http/journal_entry_handler.go:156,246` | Finding for later phases |
| Runtime DB roles | `backend/scripts/db/create_runtime_roles.sql` (`pos_migrator`, `pos_app`, NOSUPERUSER/NOBYPASSRLS); `182_accounting_runtime_roles.sql` for `pos_accounting` | Live role/grant state NOT YET OBSERVED |
| Posting rule gating | Auto-posting rules can be disabled; sales/COGS skip silently when disabled (`auto_posting_engine.go:58-61,151-154`); payroll with Tx refuses commit (`:288-292`) | Sales path: disabled rule returns nil without error (potential silent non-posting; confirm at runtime) |

## Findings worth escalating (static, not yet runtime-confirmed)

1. `GenerateBalanceSheet` excludes revenue/expense, so any period with profit/loss shows UNRECONCILED (`financial_report_usecase.go:304`).
2. Order tax is a fixed 10% for all tenants (`order_pricing.go:46`).
3. Journal void bypasses the closed-period guard and dates reversal at `time.Now()` (`journal_entry_repo.go:464-472`).
4. `round2` truncates incorrectly for negatives and is used in the equation check (`financial_report_usecase.go:528`).
5. Inventory valuation report ignores `asOf` (`financial_report_usecase.go:488-526`).
6. Journal attachments are unvalidated JSON (`journal_entry_handler.go:156,246`).
