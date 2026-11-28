# 25 - Cross-module dependency matrix

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Method: static reading of `backend/**` and `frontend/apps/dashboard/**` plus SELECT-only queries on `pos_saas` (container `pos-postgres`). No endpoint was invoked, no row written. The runtime DB holds almost no business data (2 tenants, 6 products, 0 orders, 0 journal entries; `pg_stat_user_tables`), so every runtime-labelled statement below is about schema/catalog, not behaviour.

Labels: CONFIRMED-BY-SOURCE (CBS), CONFIRMED-BY-READONLY-RUNTIME (CBR), DOCUMENTED-INTENT (DI), INFERRED (INF), UNKNOWN, REQUIRES-EXECUTION (REX).
Link status: IMPLEMENTED / PARTIAL / NOT IMPLEMENTED.

## 1. How accounting is fed (the only four mechanisms that exist)

| # | Mechanism | Evidence | Label |
|---|---|---|---|
| M1 | Outbox-driven auto-post of a paid order (sales + COGS). Only event type is `paid_order` (DB CHECK). Written in the same tx as `orders.payment_status='paid'`; replayed by a 15 s worker. | `backend/internal/repository/order_repo.go:545-605` (enqueue at the `paid` branch), `backend/internal/repository/accounting_event_repo.go:55-72`, `backend/migrations/181_durable_accounting_event_outbox.sql:11` (`CHECK (event_type = 'paid_order')`), worker `backend/internal/service/accounting_event_worker.go:11-60`, started `backend/cmd/api/main.go:513-515`, consumer `backend/internal/usecase/order_usecase.go:465-480,482-515` | CBS |
| M2 | Synchronous in-transaction posting by the owning repository: customer invoice create/update, supplier invoice approve, customer payment, supplier payment, payroll mark-paid. | `invoice_customer_repo.go:522,539-558`; `invoice_supplier_repo.go:295-348`; `payment_acct_repo.go:129-146,275-293`; `hr_salary_handler.go:508-547` + `auto_posting_engine.go:268` | CBS |
| M3 | Human-authored journals: manual JE, CSV import, recurring templates (manual "run" only, no scheduler). | `journal_entry_handler.go:122`, `accounting_advanced_repo.go:396-484` (`RunRecurringJournalTemplate`), grep for ticker in `cmd/api/main.go` finds only subscription/accounting workers (`main.go:513-516`) | CBS |
| M4 | DB triggers. Triggers exist for: stock reservation (`orders`, `order_items`), journal/GL integrity, invoice line/total maths, salary totals, attendance hours, leave days. **No trigger writes to `journal_entries`, `general_ledger`, or `bank_accounts_ledger`.** | `backend/migrations/180_order_stock_reservations.sql:20-101`; CBR trigger catalog (`information_schema.triggers`, 70 rows: only `cfr002_*`, `trg_journal_balance`, `trg_invoice_*`, `order_stock_*`, `trigger_calculate_*`, audit triggers) | CBS + CBR |

Dead posting code: `AutoPostingEngine.PostSupplierInvoiceEntry`, `PostSupplierPaymentEntry`, `PostDepreciationEntries` have zero callers (`grep -rn` over `backend/internal` and `backend/cmd`, non-test: only the three definitions, `auto_posting_engine.go:462,477,492`). The rules `SUPPLIER_INVOICE`, `SUPPLIER_PAYMENT`, `DEPRECIATION` are seeded (CBR: 6 `auto_posting_rules` rows for tenant 1 incl. these) but are inert. CBS.

## 2. Module-by-module matrix

Modules: ORD = orders (POS + storefront); KDS = kitchen queue; DRV = drivers/delivery; INV = products + stock + movements; SUP = suppliers + supplier invoices (AP); CUS = accounting customers + customer invoices (AR); PAY = accounting payments received/made; BNK = `bank_accounts_ledger`; HR = employees/attendance/leaves; PRL = salaries; GL = journals + general ledger; TAX = tax periods/VAT; RPT = financial reports; OPS = operational reports + dashboard.

Cell code: D = direct code call / SQL in same tx; E = event (outbox); T = DB trigger; R = read-only dependency (query joins/reads); I = indirect only via shared column/table; X = no link; (P) = partial.

Read as: row module -> column module.

| from \ to | ORD | KDS | DRV | INV | SUP | CUS | PAY | BNK | HR | PRL | GL | TAX | RPT | OPS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| ORD | - | D | I | T,D | X | X | X | X | X | X | E | X | I | I |
| KDS | D(sql) | - | X | X | X | X | X | X | X | X | X | X | X | R(P) |
| DRV | X (no callers) | X | - | X | X | X | X | X | X | X | X | X | X | X |
| INV | R | X | X | - | X | X | X | X | X | X | D(COGS only, via ORD) | X | R | R |
| SUP | X | X | X | **X** | - | X | D | X | X | X | D | R | R | X |
| CUS | I (`order_id` col only) | X | X | X | X | - | D | X | X | X | D | R | R | X |
| PAY | X | X | X | X | D (supplier balance) | D (customer balance) | - | D | X | X | D | X | R | X |
| BNK | X | X | X | X | X | X | R | - | X | X | R (txn list from GL) | X | R (KPI cash) | X |
| HR | X | X | X | X | X | X | X | X | - | **X** | X | X | R | X |
| PRL | X | X | X | X | X | X | X | **X** | I (FK) | - | D (tx) | X | R | X |
| GL | X | X | X | X | X | X | X | I (`gl_account_id`) | X | X | - | X | R | X |
| TAX | X | X | X | X | R | R | X | X | X | X | **X** | - | X | X |
| RPT | - | - | - | R | R | R | R | R | - | R | R | - | - | - |
| OPS | R | R | - | R | - | - | - | - | R | R | R(KPI) | - | - | - |

Evidence for each non-trivial cell:

| Cell | Status | Evidence | Label |
|---|---|---|---|
| ORD->KDS | IMPLEMENTED | `order_usecase.go:309-322` enqueues `kds_queue` on `confirmed`; enqueue `kds_repo.go:74-85` | CBS |
| KDS->ORD | IMPLEMENTED (bypasses usecase) | `kds_repo.go:700-740` `syncOrderStatusAfterBump` runs raw `UPDATE orders SET status` on bump (new->in_progress => preparing; in_progress->ready => ready). Skips `CanTransitionStatus`, relies on `expectedOrderStatus` predicate. | CBS |
| ORD->INV | IMPLEMENTED | Reservation/decrement of `products.quantity_in_stock` by trigger when `pending->confirmed`: `180_order_stock_reservations.sql:20-60` (replaced by `181_partial_order_cancellations.sql:46-103`); stock check at create `order_usecase.go:144`; sale movement + COGS at paid time `auto_posting_engine.go:143-260` | CBS |
| ORD->GL | IMPLEMENTED (event-driven) | M1 above. Sales entry `auto_posting_engine.go:50-141` (entry date `order.CreatedAt`, :137), COGS `:143-260` | CBS |
| ORD->PAY / ORD->BNK | NOT IMPLEMENTED | Order payment is a status column only (`orders.payment_status`, `domain/order.go:ValidPaymentStatus`); collection account is chosen by method->GL account code in `auto_posting_engine.go:611-626` (`resolvePaymentAccount`: cash 1110, card 1113, online/wallet 1111, else 1120). No `bank_accounts_ledger` update on order payments (`AdjustBalanceTx` callers are only `payment_acct_repo.go:164,310`). | CBS |
| ORD->CUS | NOT IMPLEMENTED | `orders` has no customer FK: CBR columns matching `%customer%`: `customer_name, customer_email, customer_phone` only. `invoices_customer.order_id` is a nullable reference validated for scope only (`invoice_customer_repo.go:475-486`); it neither suppresses the order's own auto-posted revenue nor copies lines. | CBR + CBS |
| ORD->TAX | NOT IMPLEMENTED | `GetVATSummary` reads only `invoices_customer` / `invoices_supplier` (`tax_repo.go:260-291`). The 10% order tax (`order_pricing.go:46`) posted to account 2120 never reaches the VAT report. | CBS |
| KDS->OPS | PARTIAL | KDS analytics are separate queries over `kds_queue`; no financial content (`kds_repo.go:300-390`). | CBS |
| ORD->DRV, DRV->ORD | PARTIAL / NOT IMPLEMENTED | Assignment writes `driver_assignments` + `drivers.availability_status` only (`driver_repo.go:411-470`). `UpdateAssignmentStatus` (`:558`) and `UpdateDeliveryMetrics` (`:701`) have no non-test callers (grep). Order reaches `delivered` only via `PUT /api/v1/admin/orders/{id}/status` (`admin_order_handler.go:357`). Delivery fee is forced to 0 (`order_pricing.go:52`); no driver payout/COD concept. | CBS |
| SUP->INV | NOT IMPLEMENTED | `invoice_supplier_lines` / `invoices_supplier` have no `product_id` and no inventory column (grep over `142_create_invoices_supplier.sql`, `141_*.sql`, `138_*.sql`: no match); approval posts GL only (`invoice_supplier_repo.go:295-348`). No purchase-order or goods-receipt objects (grep `purchase_order|goods_receipt|grn` over `internal`, `migrations`, `cmd`: no match; CBR table-name filter `%purchase%`: none). | CBS + CBR |
| SUP->GL | IMPLEMENTED | Approve: Dr expense lines (default acct 5110 = COGS account), Dr input VAT 1150, Cr payable (supplier's `payable_account_id` or 2101) `invoice_supplier_repo.go:498-572`; void reverses `:350-386` | CBS |
| SUP balance | PARTIAL | `suppliers.current_balance` is decremented by supplier payments (`payment_acct_repo.go:302-308`) but nothing increments it on invoice approval; only other writer is manual edit `supplier_repo.go:172`. INF: balance goes negative for invoiced-and-paid suppliers unless manually maintained. Contrast customers, which are maintained on invoice create/update/void (`invoice_customer_repo.go:151,285-292,357`). | CBS (writers) / INF (effect) |
| PAY->SUP/CUS | IMPLEMENTED | Applies amount to invoice (`payment_acct_repo.go:568-600`), updates counterparty balance, posts JE, adjusts bank ledger balance, all one tx. Unapplied customer payments allowed (no invoice) `:437-490`. | CBS |
| PAY->BNK / GL | IMPLEMENTED | JE posted to the bank's `gl_account_id` or fallback code by method `payment_acct_repo.go:546-566`; `bank_accounts_ledger.current_balance` adjusted `:164,310`, `bank_account_repo.go:282-311` | CBS |
| BNK->GL | PARTIAL | Bank account create stores `opening_balance`/`current_balance` as plain columns, no journal (`bank_account_repo.go:95-135`). "Transactions" view is derived from `general_ledger` by `gl_account_id` (`:208-280`). Two parallel balances: column vs GL. | CBS |
| Bank transfers | NOT IMPLEMENTED | grep `bank_transfer_between|fund_transfer|account_transfer|interbank|inter_account` none; CBR table filter `%transfer%`: none; no route (API inventory `08-api-inventory.csv`). | CBS + CBR |
| HR->PRL (attendance->payroll) | NOT IMPLEMENTED | `CreateSalary` takes `base_salary, overtime_hours, days_worked, total_hours_worked...` from the request body (`hr_salary_handler.go:154-200`); `salary_repo.go` and `hr_salary_handler.go` contain no reference to `attendance` (grep). Salary totals computed by trigger from supplied components (`024_create_salaries_table.sql:calculate_salary_totals`). Leaves likewise not read. | CBS |
| PRL->GL | IMPLEMENTED (cash basis) | `MarkAsPaid` locks salary, flips status, posts payroll JE in the **same tx**, commits together `hr_salary_handler.go:508-556`, `salary_repo.go:344-400`, `auto_posting_engine.go:268-345`. No accrual at calculation time. | CBS |
| PRL->BNK | NOT IMPLEMENTED | Payroll credits the account for `cash`/`bank` rule key (`auto_posting_engine.go:295-300`); never calls `AdjustBalanceTx`; no bank account id on salary at all (`024_*.sql` stores free-text `bank_account_number`). | CBS |
| PRL approval gate | NOT IMPLEMENTED | `salaries.is_approved/approved_by` exist (`024_*.sql`) but `MarkAsPaidTx` does not test them (`salary_repo.go:344-400`). | CBS |
| Expenses (module) | NOT IMPLEMENTED | No expense table/route (CBR `%expense%` none; grep `expense_claim|petty_cash_voucher` none). Expenses reach the ledger only via supplier invoices (AP), payroll, manual/imported/recurring JEs. | CBS + CBR |
| TAX<-CUS/SUP | IMPLEMENTED (invoices only) | `tax_repo.go:260-291`; period snapshot columns `output_vat,input_vat,net_vat` (`:171-185`); `CloseTaxPeriod` only flips status (`:293-318`), posts nothing. | CBS |
| GL->TAX | NOT IMPLEMENTED | No VAT settlement/clearing entry; no read of accounts 2120/1150 in tax repo. | CBS |
| RPT<-GL | IMPLEMENTED | P&L/Cash-flow from `journal_entry_lines` (`financial_report_usecase.go:173-250,366-398`); TB/BS from `general_ledger` (`:252-353`). | CBS |
| RPT<-invoice tables | IMPLEMENTED | Aging from `invoices_customer`/`invoices_supplier.balance_due`, not GL (`:428-486`). | CBS |
| RPT<-products | IMPLEMENTED | Inventory valuation = `quantity_in_stock * cost` live, `asOf` is a label only (`:488-526`). | CBS |
| OPS<-orders | IMPLEMENTED | Sales/profitability/inventory reports over `orders`/`order_items`/`products` with `status <> 'cancelled'` (`report_handler.go:2174-2420,2480-2560`). | CBS |
| OPS<-legacy `inventory` | PARTIAL / dead input | Turnover "waste" and opening qty read table `inventory` (`report_handler.go:2513,2535`). No Go code or migration inserts into `inventory` (grep `INTO inventory` over `internal`, `migrations`, `scripts`: none). CBR `inventory` has 0 rows. Stock adjustments write `inventory_movements` instead (`inventory_adjustment_handler.go:92-99`). INF: waste_qty and opening_qty are structurally always empty/fallback. | CBS + CBR (0 rows) |
| OPS->GL | NOT IMPLEMENTED | Ops reports never read journals; KPIs endpoint mixes orders (revenue) with journals (expenses) `accounting_advanced_repo.go:837-965`. | CBS |

## 3. Target chains: link-by-link verdict

### 3.1 POS -> Orders -> Payments -> Inventory -> Accounting

| Link | Status | Evidence |
|---|---|---|
| POS UI -> `POST /api/v1/pos/orders` | IMPLEMENTED | `frontend/apps/dashboard/src/lib/api/posApi.ts:248`; route `backend/cmd/api/pos_routes.go:16` (module 10 WRITE); server forces `order_source=pos` `order_source_guard.go:11-17` |
| Order create -> auto-confirm (-> KDS enqueue -> stock reservation) | IMPLEMENTED | `order_usecase.go:190-214` (POS `Auto-confirmed POS order`), reservation trigger `180_*.sql`; failure aborts the POS sale `:203-206` |
| Order -> Payment | PARTIAL | Separate second call `PUT /admin/orders/{id}/payment-status` from the browser (`posApi.ts:266`). Create and pay are **two non-atomic client calls**; a POS order can exist confirmed/unpaid with stock reserved if the second call never happens. No payment row, tender, change, or split. CBS |
| Payment -> Accounting | IMPLEMENTED | `order_repo.go:545-605` outbox + worker; M1 |
| Payment -> Inventory (cost movement) | IMPLEMENTED (late) | COGS entry and `sale` movement are produced at **payment** time (`auto_posting_engine.go:143-260`), whereas stock decrement is at **confirmation** (trigger). Two different movement rows per sale (`reservation` -q at confirm; `sale` -q at COGS) both inserted into `inventory_movements`; idempotency table `accounting_inventory_movement_events` dedupes only `sale` (`inventory_movement_repo.go:28-80`, legacy lookup filters `movement_type = sale`). Sum of `inventory_movements.quantity` for a sold item is therefore -2q while `products.quantity_in_stock` moved by -q. CBS (both inserts) / INF (aggregate effect; no order data in DB to confirm: REX) |
| Accounting -> reports | IMPLEMENTED | see 27 |

### 3.2 Storefront -> Orders -> Kitchen -> Delivery -> Accounting

| Link | Status | Evidence |
|---|---|---|
| Storefront -> Orders | IMPLEMENTED | `POST /api/v1/public/orders` `public_order_handler.go:39`; server prices from DB `order_usecase.go:144-163`, `order_pricing.go` |
| Orders -> Kitchen | IMPLEMENTED (manual confirm for non-POS) | Non-POS orders stay `pending` until staff confirm (`ConfirmOrder` `order_usecase.go:354`); KDS enqueue on `confirmed` |
| Kitchen -> Order status | IMPLEMENTED | KDS bump `kds_repo.go:190-262,700-740` |
| Order -> Delivery | PARTIAL | `POST /admin/orders/{orderId}/assign` creates assignment (`admin_driver_handler.go:407`, `driver_repo.go:411-470`); no status propagation either direction |
| Delivery -> Accounting | NOT IMPLEMENTED | No delivery fee (forced 0), no COD cash-collection, no driver settlement. Sales posting has a delivery-revenue line but it is unreachable while fee = 0 (`order_pricing.go:52`, `auto_posting_engine.go:100-111`) |
| Online payment -> Order | NOT IMPLEMENTED | Paymob integration is subscription billing only (`payments` table `115_create_payments_table.sql`, `payment_usecase.go`: `paymob_order_id`, no order link). Storefront orders are paid only when staff set `payment_status`. CBS |

### 3.3 Suppliers -> Purchases -> Inventory -> AP

Suppliers (master) IMPLEMENTED; "Purchases" = supplier invoice only (no PO/GRN) PARTIAL; Purchases->Inventory NOT IMPLEMENTED (see matrix); Purchases->AP IMPLEMENTED (GL Cr payable, `invoice_supplier_repo.go:560-566`); AP payment IMPLEMENTED. Inventory cost therefore never increases through purchasing: the stock-increase path is `POST /api/v1/inventory/adjust` with reason `purchase` and **no unit cost field** (`inventory_adjustment_handler.go:20-26,92-99`), writing no journal. GL account 1130 (Inventory) is credited by COGS but no code path debits it (only manual JE or choosing 1130 as a supplier-invoice line `account_id`). CBS; consequence for valuation vs GL: INF.

### 3.4 Customers -> Invoices -> Receipts -> AR

IMPLEMENTED (accounting-side only). Invoice create posts immediately Dr AR / Cr revenue / Cr VAT (`invoice_customer_repo.go:452-558,561-636`) even while invoice `status=draft`; payment posts Dr bank/cash / Cr AR (`payment_acct_repo.go:101-128`). Customer balance maintained in three places: GL, `invoices_customer.balance_due`, `customers_accounting.current_balance`. No credit note (grep none), no write-off, no customer linkage from storefront/POS orders. CBS.

### 3.5 Employees -> Attendance -> Payroll -> Expenses -> Liabilities

Employees/attendance IMPLEMENTED (HR domain). Attendance->Payroll NOT IMPLEMENTED. Payroll->GL IMPLEMENTED at paid time only. "Expenses" exist only as the Dr side of the payroll JE (accounts 5210/5220); liabilities: tax, social insurance, health, pension, loan, other deductions credited to payable accounts at pay time (`auto_posting_engine.go:347-430`) with fallback chain ending in 2110. No later settlement of those liabilities except manual JE. Employer-side contributions are not modelled (no field in `domain.Salary` for them, `024_*.sql`). CBS / INF for the "not modelled" negative (checked column list only).

### 3.6 Inventory -> Product cost -> COGS -> Profitability

PARTIAL and internally inconsistent: three different cost sources.
1. COGS journal: latest `inventory_movements.unit_cost` for types `opening|purchase|adjustment` with `quantity>0 AND unit_cost>0`, else `products.cost` (`inventory_movement_repo.go:158-198`). The only code paths that write a nonzero `unit_cost` are product create/update stock movements (`product_stock_movement.go:16-38`, uses `product.Cost`) and the COGS sale movement itself. The inventory adjust endpoint writes `unit_cost` default 0 (migration `144_create_inventory_movements.sql`), so it never sets cost.
2. Ops profitability/valuation: live `products.cost` (`report_handler.go:2364-2420`; `financial_report_usecase.go:488-526`).
3. Menu-item profitability and break-even: reference `products.cost_price`, a column that does not exist. CBR: `SELECT p.cost_price FROM products p LIMIT 1` -> `ERROR: column p.cost_price does not exist`; CBR `information_schema` shows only `price, cost, discount_price` for products; sources `accounting_advanced_repo.go:976,1045`. HTTP behaviour of `GET /reports/financial/menu-profitability` and `/break-even` therefore expected to be a 5xx or error payload: REQUIRES-EXECUTION to observe the response; the SQL error itself is CBR.
`order_items` has no cost snapshot (CBR column list: no `%cost%` column), so historical profitability changes whenever `products.cost` is edited (INF from schema).

### 3.7 Bank accounts -> Transfers -> GL -> Cash flow

Bank accounts IMPLEMENTED (CRUD); transfers NOT IMPLEMENTED; GL link PARTIAL (see BNK->GL); cash flow report is a heuristic that does not read cash/bank accounts: operating = sum(credit-debit) on revenue+expense accounts, investing = accounts coded `12%`, financing = equity or `22%` (`financial_report_usecase.go:386-391`). Dashboard KPI "cash_position" reads `SUM(bank_accounts_ledger.current_balance)` (`accounting_advanced_repo.go:872`), a different source. Bank reconciliation `book_balance` is client-supplied on create (`accounting_advanced_handler.go:64-96`) and later recomputed from user-entered line amounts (`accounting_advanced_repo.go:2126-2150`), never from GL or `current_balance`. CBS.

### 3.8 Accounting -> Financial reports -> Dashboard

Reports IMPLEMENTED (see 27). Dashboard home does not call any financial report: it lists 5 recent orders and sums `total` of those with paid status client-side as "Recent Paid Revenue" (`frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:281-300,374-377,502-506`). CBS.

## 4. Shared tables and services

| Shared object | Writers | Readers | Risk |
|---|---|---|---|
| `products.quantity_in_stock` | trigger `maintain_order_stock`; `cancel_order_item_quantity()`; product create/update repos (`product_repo.go:154`, `product_update_repo.go:107`); `inventory_adjustment_handler.go:81-85` | order create, ops/financial reports | Multiple writers, no single stock ledger; movements table is not authoritative (see 3.1) |
| `products.cost` | product update only | ops reports, valuation, COGS fallback | Mutable; no history, no per-sale snapshot |
| `inventory_movements` | trigger (reservation/release), COGS engine, product create/update, adjust endpoint | `GetProductCost`, `ListMovements` | Mixed semantics (`movement_type` free text: adjust endpoint stores the request `reason` as movement type) |
| `journal_entries` / `journal_entry_lines` / `general_ledger` | `JournalEntryRepository` (all M1-M3 funnel here), GL repo | all financial reports | Strongest shared sink; protected by integrity triggers (`180_enforce_journal_integrity.sql`, CBR `cfr002_*`) |
| `chart_of_accounts` + `auto_posting_rules` | seed migrations 152/156 (existing restaurants only), account CRUD | engine `resolveAccount` by **account code** | Hard-coded code defaults (1110,1111,1113,1120,1130,2101,2120,4100,4150,5110,5210...) |
| `bank_accounts_ledger.current_balance` | payment_acct_repo only | KPIs | Parallel to GL |
| `customers_accounting.current_balance`, `suppliers.current_balance` | invoice/payment repos; manual edit | UI only (aging reads invoice tables) | Parallel to GL and invoice balance_due |
| `fiscal_periods` | no Go writer (grep: only reads at `journal_entry_repo.go:817`, `cross_report_diagnostics.go:305`, `current_earnings_diagnostic.go:155`); no route | posting guard | Table is unmaintainable via the product (CBR: 0 rows for both restaurants) |
| `tax_periods` | tax handler | VAT UI | Independent from fiscal periods, posts nothing |
| Service: `AutoPostingEngine` | order usecase, salary handler | - | Only injected into `OrderUseCase` and `HRSalaryHandler` (`main.go:323-330,422`) |
| Service: `AccountingEventWorker` | - | `OrderUseCase` as processor | Single worker; only paid-order events |

## 5. Event-driven links (complete list)

1. `orders.payment_status -> paid` => row in `accounting_event_outbox` (same tx) => worker => `PostSalesEntry` + `PostCOGSEntry` + sale movements. CBS (`order_repo.go:593-598`).
2. Reconciliation loop re-creates missing intents for paid orders and counts failed events (`accounting_event_repo.go:74-122`), admin redrive `POST /accounting/paid-order-events/{id}/redrive` (`accounting_event_admin_handler.go:53`).
3. WebSocket broadcast to KDS on enqueue (`order_usecase.go:316`): non-financial.
4. Webhook retry queue (Paymob): subscription billing only.
No event for: payroll (synchronous), purchases, bank movements, stock adjustments, order cancellation, refunds, VAT, period close.

## 6. Missing integrations (consolidated; each also in 26 section 3)

SUP->INV, PO/GRN, ORD->CUS, ORD->PAY ledger, ORD->BNK, HR attendance->PRL, PRL->BNK, PRL approval gate, stock adjustment->GL (waste/shrinkage/return), ORD->TAX report, DRV->ORD completion, delivery->accounting, bank transfers, fiscal period maintenance/closing, VAT settlement, credit notes/refunds, fixed assets->depreciation (table `fixed_assets` exists: CBR; no Go reader/writer except the dead `PostDepreciationEntries`).

## 7. Duplicate responsibilities / inconsistent sources of truth (cross-reference)

Detailed formulas and ids S1-S12 are in `27-financial-data-lineage.md`. Highlights at module level:
- "Payments" is three unrelated things: `payments` (SaaS subscription), `payments_received/made` (accounting), `orders.payment_status` (no record). CBS/CBR.
- "Customers" is two things: free-text order customer fields vs `customers_accounting`. CBS.
- "Suppliers" has two balances (`suppliers.current_balance` vs `invoices_supplier.balance_due`) that are not maintained symmetrically. CBS.
- Stock quantity has four stores: `products.quantity_in_stock`, `inventory_movements`, `order_stock_reservations`, legacy `inventory` (unwritten). CBS/CBR.
- Tax has three computations: order 10% hard-coded, per-line invoice rate, tax_periods snapshot. CBS.
- Closing has two disconnected mechanisms: `tax_periods.status` and `fiscal_periods.status`. CBS.
- New-tenant bootstrap gap: registration (`auth_repo.go:86-200` inserts tenant, user, restaurant, theme, roles) creates **no** chart of accounts, posting rules, fiscal periods or bank accounts; migrations 152/156 seed only restaurants existing at migration time (`156_seed_chart_of_accounts.sql:1-2`, `152_*.sql` `FROM restaurants r`). CBR: tenant 2/restaurant 2 has 0 `chart_of_accounts`, 0 `auto_posting_rules`, 0 `fiscal_periods`, whereas tenant 1 has 62/6/0. Consequence (INF): first paid order of a new tenant cannot post (`resolveAccount` returns `account not found`, `auto_posting_engine.go:628-645`), the outbox event would go retryable then failed after 8 attempts (`181_*.sql` `max_attempts` default 8). REX to observe.
