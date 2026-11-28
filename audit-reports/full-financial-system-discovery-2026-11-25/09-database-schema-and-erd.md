# 09 - Database schema, constraints, security and integrity probes

Baseline git SHA 133260db66b1143c254d79476daa7650b41a33f0. Database: dev container `pos-postgres`, db `pos_saas` (PostgreSQL 15). Read-only catalog queries only. Raw evidence: `evidence/database/*.txt`; SQL: `sql/*.sql`.

## 1. Inventory of the running schema

| Metric | Value | Evidence |
|---|---|---|
| Base tables in `public` | 102 | `SELECT count(*) FROM pg_class ... relkind='r'` |
| Foreign keys | 287 | pg_constraint contype='f' |
| Non-internal triggers | 47 | `evidence/database/triggers.txt` |
| Tables with RLS enabled (and FORCEd) | 20 of 102 | `evidence/database/rls.txt` |
| RLS policies | 23 | pg_policies |
| Float / real / money columns | 0 monetary (4 non-monetary `double precision`: `restaurant_tables.position_x/y`, `table_layouts.canvas_width/height`) | `evidence/database/float_columns.txt` |
| Roles (non-system) | `pos_app`, `pos_migrator`, `pos_accounting` (all NOSUPER, NOBYPASSRLS), `postgres` (super, bypassrls) | `evidence/database/grants_roles.txt` |
| Migrations dir | 163 files, latest 192 (190 = normal-balance history guard) | `backend/migrations` |

### Dev DB data volume (important for interpreting probes)
Row counts in the dev DB: `orders`=0, `journal_entries`=0, `journal_entry_lines`=0, `general_ledger`=0, `invoices_customer`=0, `invoices_supplier`=0, `payments_made`=0, `payments_received`=0, `salaries`=0, `inventory_movements`=0, `accounting_event_outbox`=0, `fiscal_periods`=0, `payments`=0; `chart_of_accounts`=62 (tenant 1 / restaurant 1 only); 2 tenants, 2 restaurants (`evidence/database/integrity_probes.txt`, `probe_details.txt`). **All data-level integrity probes are therefore vacuous (0 rows tested); they prove the SQL runs, not that production data is clean.** Structural findings (section 6) are the valuable part.

## 2. Monetary-type analysis

Monetary columns in business tables are all `numeric` (scales 2; journal/invoice/payment 15,2; products/orders/salaries 10,2; inventory quantity 12,3; exchange rate 18,8; fixed_assets 14,2). Listing: `evidence/database/monetary_types.txt` (108 numeric money-like columns: 31x(10,2), 50x(15,2), 4x(12,2), 2x(14,2), 9x(5,2), 2x(8,2), 1x(18,8)).

Anomalies:
1. **Integer money (cents) in SaaS billing**: `payments.amount_cents integer`, `subscription_plans.price_cents integer`, `tenants.last_payment_amount integer` (int32 caps at 21,474,836.47 major units; inconsistent with numeric(,2) used everywhere else). Evidence: `money_integer_columns.txt`.
2. **No float/real/money DB types for money.** BUT the Go application models and arithmetic use `float64` for every amount (e.g. `usecase/auto_posting_engine.go:645-647 roundedAmount = math.Round(v*100)/100`, `usecase/financial_report_usecase.go:528-530 round2` uses `int(value*100+0.5)` which mis-rounds negatives). Rounding drift is mitigated by cent-validation in `usecase/auto_posting_validation.go:21-103` for auto-posting only; manual/invoice/payment paths use `roundAccountingAmount` (float) and a 0.009 tolerance (`repository/journal_entry_repo.go` PostEntryTx, ~line 395).
3. **Money precision mismatches**: orders/order_items/products/salaries are numeric(10,2) (max 99,999,999.99) while journal/GL are numeric(15,2); `usecase/order_pricing.go:11` caps orders at 9,999,999,999 cents to match.
4. **Type mismatches between ownership columns**: `orders/order_items/accounting_event_outbox/order_item_cancellations/order_stock_reservations` use `bigint` tenant/restaurant ids; journal/GL/invoice/payment/COA/salaries/inventory_movements use `integer`; `journal_entries.reference_id integer` vs `orders.id bigint` (overflow risk after 2^31 orders; `inventory_movements.order_id bigint` vs `reference_id integer`).
5. **Currency**: 16 currency columns, all defaulting to `'EGP'` (`evidence/database/currency_columns.txt`). **`orders`, `order_items`, `journal_entries`, `journal_entry_lines`, `general_ledger`, `payments_made`, `payments_received` have NO currency column.** Multi-currency is not implemented: `exchange_rates` table + `ConvertAmount` exist (repository/accounting_advanced_repo.go:670) but no posting or report uses them. Order quote hardcodes `Currency: "EGP"` (usecase/order_pricing.go:99); Financial Hub frontend hardcodes `currency: 'EGP'` (frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/page.tsx:~33).

## 3. Per-table catalogue (financial and operational)

Legend: PK = primary key; FK(ondelete); own = ownership columns; st = status column + allowed values (CHECK); money = numeric(p,s).

### 3.1 Sales / orders
| Table | PK | own | Money columns | Status/CHECK | Unique / notable constraints | FKs (delete behaviour) | Key indexes |
|---|---|---|---|---|---|---|---|
| `orders` | id bigint | tenant_id, restaurant_id (bigint) | subtotal(10,2) NOT NULL, tax_amount, discount_amount, delivery_fee (10,2, nullable default 0), total_amount(10,2) NOT NULL | `status` chk_order_status {pending,confirmed,preparing,ready,out_for_delivery,delivered,cancelled}; `payment_status` chk_payment_status {pending,paid,failed,refunded}; `payment_method` free text; `order_source` free text (default 'website') | `orders_order_number_key UNIQUE(order_number)` is GLOBAL (not per tenant); chk_order_totals (total>=0, subtotal>=0); **no CHECK total = subtotal+tax+fee-discount**; no currency | fk_orders_tenant CASCADE, fk_orders_restaurant CASCADE, fk_orders_table SET NULL | idx_orders_paid_accounting_recovery (partial paid), (tenant,restaurant,created_at), (tenant,restaurant,status,created_at), public_token_hash unique partial |
| `order_items` | id bigint | tenant_id, restaurant_id | unit_price, discount_amount, total_price (10,2) | `cancelled_quantity` CHECK 0..quantity; chk_prices_non_negative; chk_quantity_positive | no unique | order CASCADE, product RESTRICT, restaurant/tenant CASCADE | - |
| `order_status_history` | id bigint | via order_id only | - | old/new status text (unconstrained) | - | order_id | - |
| `order_item_cancellations` | (request_key uuid scoped) | tenant,restaurant | - | quantity, reason, actor | idempotency receipt (migration 181) | - | - |
| `order_stock_reservations` | order_id+product_id | tenant,restaurant (bigint) | - | released bool, quantity>0, released_quantity<=quantity | stock hold ledger (migration 180) | - | - |
| `payments` | id bigint | tenant_id | `amount_cents integer` | status; `refund_state text` ('none'...); RLS policy payments_tenant_isolation | **SaaS subscription charges (Paymob) - NOT order payments** | - | - |

DELETE/UPDATE behaviour: orders/items can be deleted by cascade from tenant/restaurant; trigger `order_stock_transition` (BEFORE UPDATE/DELETE on orders, fn `maintain_order_stock`) blocks DELETE when `order_stock_reservations` exist and blocks delivered->cancelled; trigger `order_stock_items` guards `order_items`. **No trigger prevents editing `total_amount`/`payment_status` on a paid order or deleting a paid order that never had stock reservations (e.g. non-tracked products).** No history-immutability on orders.

### 3.2 Journal / ledger (strongest controls)
| Table | PK | own | Money | CHECK / status | Unique | FKs | Triggers |
|---|---|---|---|---|---|---|---|
| `chart_of_accounts` | id int | tenant_id, restaurant_id | current_balance(15,2) (denormalised cache), tax_rate(5,2) | account_type {asset,liability,equity,revenue,expense}; normal_balance {debit,credit}; level>=1 | UNIQUE(tenant,restaurant,account_code) | parent_id (no cascade), tenant, restaurant | `cfr002_account_scope_write` (BEFORE UPDATE OF tenant_id, restaurant_id), `cfr003_normal_balance_history_guard` (BEFORE UPDATE OF normal_balance; migration 190) |
| `journal_entries` | id int | tenant_id, restaurant_id | total_debit/total_credit(15,2) | status {draft,posted,void}; entry_type varchar(32); is_approved bool | UNIQUE(tenant,restaurant,entry_number); **no unique on (reference_type, reference_id, entry_type)** | tenant, restaurant | `cfr002_journal_header_write` (BEFORE INS/UPD/DEL), `cfr002_journal_integrity_header` (DEFERRABLE INITIALLY DEFERRED constraint trigger) |
| `journal_entry_lines` | id int | tenant_id, journal_entry_id | debit_amount, credit_amount, tax_amount (15,2) | debit>=0, credit>=0, not both >0, at least one >0 | - | journal_entry_id CASCADE, account_id | `cfr002_journal_line_write`, `cfr002_journal_integrity_lines` (deferred), `trg_journal_balance` (AFTER INS/UPD/DEL, fn `validate_journal_balance`) |
| `general_ledger` | id int | tenant_id, restaurant_id | debit_amount, credit_amount, running_balance (15,2) | - | - | account, journal_entry, journal_line, tenant, restaurant | `cfr002_general_ledger_immutable` (BEFORE UPDATE/DELETE always raises 23514), `cfr002_general_ledger_insert`, `cfr002_journal_integrity_ledger` (deferred), `cfr003_lock_account_after_cfr002_general_ledger_insert` |
| `fiscal_periods` | id | tenant, restaurant | - | period_type, status text (no CHECK found), start/end date | UNIQUE(tenant,restaurant,period_name) | closed_by->users | **none** (no trigger blocks posting into a closed period; enforced in app only) |
| `tax_periods` | id | tenant, restaurant | output_vat, input_vat, net_vat (15,2) | status (open/closed; set by app) | - | closed_by->users | none |
| `accounting_event_outbox` | id bigint | tenant_id, restaurant_id, order_id (bigint) | - | event_type CHECK ='paid_order'; status CHECK {pending,processing,retryable,completed,failed}; attempts>=0; max_attempts>0 (default 8) | UNIQUE(event_id), UNIQUE(event_type,tenant,restaurant,order_id) | order RESTRICT | - |
| `accounting_event_redrive_audit` | id | outbox_id | - | reason length 1..1000, actor>0 | - | outbox RESTRICT | pos_app: INSERT,SELECT only |
| `accounting_inventory_movement_events` | (tenant,restaurant,order,product) | - | - | - | PK is the idempotency claim | order RESTRICT | pos_app: INSERT,SELECT only |
| `auto_posting_rules` | id | tenant, restaurant | default_accounts jsonb | is_active | (rule_code per scope) | - | - |
| `recurring_journal_templates`, `account_budgets` (currency, period_month 1..12), `exchange_rates` (rate numeric(18,8)) | id | tenant, restaurant | yes | frequency CHECK | - | - | - |

Immutability semantics (from `evidence/database/trigger_functions.txt`):
- journal header: INSERT must be `draft` and restaurant must belong to tenant; DELETE only if draft; draft->posted requires posted_at+posted_by; draft->void requires void audit fields; posted->void allowed only if every other header field is unchanged and void audit fields+reason present; any other transition raises 23514 (`cfr002_guard_journal_header_write`).
- lines: any write requires parent journal `draft`; account must belong to the same tenant+restaurant; lines cannot move between journals (`cfr002_guard_journal_line_write`).
- deferred integrity (`cfr002_assert_journal_integrity_deferred`): at commit, posted journals need >=2 lines, sum(debit)=sum(credit)=header totals, and exactly one matching GL row per line; voided-after-posting journals need one original and one inverse GL row per line; draft journals cannot have GL rows.
- `general_ledger` UPDATE/DELETE always rejected ("record corrections through an audited journal reversal").
- migration 190 trigger: normal_balance cannot change once any GL row exists for the account.
- grants: `pos_accounting` has full DML on journal/GL/COA tables (`evidence/database/grants_roles.txt`) - history safety comes from triggers, not privileges (migration 182 comment: "protected by CFR-002 guards for every database role").
- `trg_journal_balance` (validate_journal_balance) pre-dates CFR-002 and is redundant.

### 3.3 Receivables / payables / banking
| Table | PK | own | Money | Status/CHECK | Unique | FKs | Triggers |
|---|---|---|---|---|---|---|---|
| `customers_accounting` | id | tenant, restaurant | credit_limit, current_balance (15,2) | - | (customer_code scoped) | receivable_account_id->COA | - |
| `suppliers` | id | tenant, restaurant | credit_limit, current_balance | - | supplier_code | payable_account_id->COA | - |
| `invoices_customer` | id | tenant, restaurant | subtotal, tax_amount, discount_amount, total_amount, amount_paid, balance_due (15,2) | status {draft,sent,partially_paid,paid,overdue,void} | UNIQUE(tenant,restaurant,invoice_number); **no CHECK balance_due = total - amount_paid** | customer_id, order_id->orders (nullable), journal_entry_id | `trg_invoice_customer_line_total`, `trg_invoice_customer_totals_ins` on lines |
| `invoice_customer_lines` | id | tenant_id, invoice_id | quantity(10,3), unit_price, tax_rate(5,2), tax_amount, discount_rate, discount_amount, line_total | - | - | invoice CASCADE, account_id | recalc triggers |
| `invoices_supplier` / `invoice_supplier_lines` | id | tenant, restaurant | same shape | status {draft,approved,partially_paid,paid,overdue,void} | UNIQUE(tenant,restaurant,invoice_number) | supplier_id, journal_entry_id | `trg_invoice_supplier_*` |
| `payments_received` | id | tenant, restaurant | amount(15,2) | status {completed,void} | UNIQUE(tenant,restaurant,payment_number) | customer_id, invoice_id, bank_account_id, journal_entry_id | none; **no void/reversal code path** |
| `payments_made` | id | tenant, restaurant | amount(15,2) | status {completed,void} | UNIQUE(tenant,restaurant,payment_number) | supplier_id, invoice_id, bank_account_id, journal_entry_id | none; **no void/reversal code path** |
| `bank_accounts_ledger` | id | tenant, restaurant | opening_balance, current_balance, last_reconciled_balance (15,2) | account_type {checking,savings,credit_card,petty_cash} | - | gl_account_id->COA | none |
| `bank_reconciliation_sessions/lines` | id | tenant, restaurant | statement/book balances | status {draft,in_progress,completed,cancelled}; line match_status {matched,unmatched,manual} | - | - | none |
| `supplier_payment_schedules`, `e_invoice_submissions` (RLS), `invoice_delivery_logs` (RLS) | id | tenant, restaurant | yes | - | - | - | - |
| `fixed_assets` | id | tenant, restaurant | purchase_cost, salvage_value, accumulated_depreciation (14,2) all >=0 | useful_life_months>0 | - | - | **table exists but NO Go code reads/writes it** (grep `fixed_assets` finds nothing outside migrations); `PostDepreciationEntries` (usecase/auto_posting_engine.go:492) has no caller |

Cross-tenant FK posture: none of invoices/payments/bank FKs carry tenant columns, so e.g. `payments_received.invoice_id` could point at another tenant's invoice at the DB level. Application predicates (`WHERE tenant_id=$ AND restaurant_id=$`) are the only guard. `journal_entry_lines` account scope IS DB-enforced (trigger above); invoice/payment->journal link is NOT.

### 3.4 Inventory / catalogue
| Table | PK | own | Money/qty | Constraints | Notes |
|---|---|---|---|---|---|
| `products` | id int | tenant, restaurant | price(10,2) NOT NULL, `cost`(10,2) nullable, discount_price(10,2), discount_percentage(5,2); `quantity_in_stock integer`, low_stock_threshold, reorder_quantity | UNIQUE(restaurant_id, sku); track_inventory bool; status text | **No `cost_price` column** (see report finding F1). Stock is integer: no fractional units (kg, litres). |
| `inventory_movements` | id int | tenant, restaurant | quantity numeric(12,3), unit_cost, total_cost (15,2) default 0 | **No CHECK on movement_type**; no UNIQUE (idempotency via `accounting_inventory_movement_events`); order_id FK SET NULL | append-only for pos_app (migration 181 `REVOKE UPDATE,DELETE FROM pos_app`); but pos_accounting/owner can mutate. Movement types seen in code: reservation, reservation_release, sale, opening, purchase, adjustment, waste, return. |
| `inventory` (legacy) | id int | none (product_id only) | quantity_change, quantity_after integer | reason | **No writers in Go code**; read only by turnover report (handler/http/report_handler.go:2513,2535). Dead/legacy; no tenant column. |
| `order_stock_reservations` | see above | | | | stock decremented at pending->confirmed (trigger `maintain_order_stock`) |
| `low_stock_alerts`, `categories`, `product_audit_log`, `product_images` | | | | | operational |

### 3.5 Payroll / HR
| Table | PK | own | Money | CHECK | Unique | FKs |
|---|---|---|---|---|---|---|
| `salaries` | id | tenant, restaurant | base_salary, overtime_rate, overtime_amount, bonus, commission, allowances, tips, other_earnings, tax, social_insurance, health_insurance, pension, loan_deduction, advance_deduction, other_deductions, gross_salary, total_deductions, net_salary (all 10,2) | positive_amounts; status {pending,processing,paid,cancelled,on_hold}; payment_method {cash,bank_transfer,check,mobile_money}; valid_pay_period; valid_payment_date (payment_date >= pay_period_end) | UNIQUE(employee_id, pay_period_start, pay_period_end) | employee CASCADE (deleting an employee deletes paid salary history!), tenant/restaurant CASCADE |
| `employees` | id | tenant, restaurant | base_salary(10,2) | - | employee_code | audit trigger `trigger_audit_employee_changes` |
| `employment_contracts`, `insurance_tracking` (monthly_premium), `attendance` (hours calc trigger), `leaves`, `resignations`, `hr_audit_logs`, `performance_reviews` | | | | | | |
Triggers: `trigger_calculate_salary_totals` (BEFORE INS/UPD recalculates gross/deductions/net), `trigger_set_salary_period`, audit. Note FK `salaries.employee_id ON DELETE CASCADE` + `salaries.tenant_id ON DELETE CASCADE` mean paid payroll rows (whose journals persist) can vanish; app blocks deleting paid salaries (hr_salary_handler.go DeleteSalary) but DB does not. Not seen: payroll accrual tables, loans/advances tables, tips ledger (see 4).

### 3.6 Audit / settings
`product_audit_log`, `hr_audit_logs`, `permission_audit_logs` (RLS, SELECT-only policy), `user_auth_audit_logs` (RLS), `audit_export_logs` (RLS, FORCE), `report_export_audit_logs` (RLS), `super_admin_audit_logs`, `accounting_event_redrive_audit`, `order_status_history`. Settings: `restaurant_settings` (delivery_fee, min/max order; **delivery_fee is not read by pricing** - `usecase/order_pricing.go:52` forces 0), `restaurant_config`, `website_settings`, `user_settings`. **There is no tenant/restaurant tax-rate, currency or fiscal-calendar settings table**; tax is hardcoded (10%) in `order_pricing.go:46`, while CoA seeds VAT accounts at 14% (`chart_of_accounts.tax_rate` 14.00 on 1150/2120).

## 4. Tables the spec expects that do NOT exist (NOT IMPLEMENTED)
Evidence: table-name regex over `information_schema.tables` (`session|shift|cash|drawer|recipe|ingredient|bom|refund|purchase|goods|receipt|expense|transfer|capital|period_clos|closing|tip|advance|loan|credit_note|stock|warehouse|payroll|accrual|asset|currency`) returns only `low_stock_alerts, impersonation_sessions, fixed_assets, bank_reconciliation_sessions, order_stock_reservations`.
| Expected | Status | Evidence |
|---|---|---|
| POS sessions / shifts / till | NOT IMPLEMENTED (POS is just an `order_source='pos'` flag on `orders`; route `POST /api/v1/pos/orders` -> same handler, cmd/api/pos_routes.go:18) | no table; `impersonation_sessions` is admin impersonation |
| Cash movements / cash drawer / petty cash ledger | NOT IMPLEMENTED (`bank_accounts_ledger.account_type='petty_cash'` is only a type) | no table |
| Refunds (order/customer) | NOT IMPLEMENTED (order cancel refused when paid: repository/order_repo.go:634; `CanTransitionOrderPaymentStatus` refuses `refunded`: domain/order.go:219-223). Only SaaS `payments.refund_state` exists (subscription refunds, super-admin) | no table |
| Payment records per order / split payments / tips / change | NOT IMPLEMENTED (single `payment_method`+`payment_status` text on `orders`) | orders columns |
| Stock movements table | EXISTS as `inventory_movements` (+ legacy unused `inventory`) | |
| Recipes / ingredients / BOM | NOT IMPLEMENTED (COGS = product.cost x qty; stock is per sellable product) | no table; engine usecase/auto_posting_engine.go:177-185 |
| Purchase orders / goods receipts / stock receiving | NOT IMPLEMENTED (supplier invoice posts expense/COGS 5110, no inventory movement; usecase repository/invoice_supplier_repo.go:498-572) | no table |
| Expenses table / petty expenses | NOT IMPLEMENTED (only manual journals / supplier invoices) | |
| Bank transfers / owner capital / owner withdrawals | NOT IMPLEMENTED as documents (only possible as manual journal entries; accounts 3100 Owner Capital exists, no Owner Drawings account in seeded CoA) | `evidence/database/coa_default_accounts.txt` |
| Credit notes / customer returns | NOT IMPLEMENTED | |
| Fiscal-year close / period_closing table | NOT IMPLEMENTED. Migration 182 references `period_closing` but `to_regclass` guard skips it; `fiscal_periods` rows are never created or closed by Go code (only SELECTed: repository/journal_entry_repo.go:813, cross_report_diagnostics.go:305). Current-earnings policy is `POLICY_PENDING` (usecase/current_earnings_diagnostic.go:13). | |
| Payroll accruals, loan/advance ledgers, tip ledger | NOT IMPLEMENTED (payroll posts only at mark-paid; 2130 Salaries Payable unused) | usecase/auto_posting_engine.go:280-460 |
| Fixed-asset register usage / depreciation | Table only; NOT wired | see 3.3 |
| Multi-currency posting | NOT IMPLEMENTED | section 2 |

## 5. RLS, grants and roles
- RLS enabled+forced on only 20 tables: `audit_export_logs, auth_identities, e_invoice_submissions, email_change_requests, identity_email_challenges, invoice_delivery_logs, membership_restaurant_grants, membership_role_grants, notification_preferences, notifications, password_change_history, payments, permission_audit_logs, report_export_audit_logs, role_permissions, tenant_memberships, user_auth_audit_logs, user_roles, user_settings, user_settings_audit_log`.
- **No RLS on `orders, order_items, products, inventory_movements, journal_entries, journal_entry_lines, general_ledger, chart_of_accounts, invoices_*, payments_made/received, salaries, employees, bank_*, fiscal_periods`.** Tenant isolation for all core financial data depends on application-layer `WHERE tenant_id` predicates (+ the CFR-002 scope triggers on journals).
- RLS policy predicate pattern: `tenant_id = NULLIF(current_setting('app.tenant_id', true),'')::bigint [AND restaurant_id = ...]` for `pos_app` (and `pos_accounting` on e_invoice_submissions/permission_audit_logs); `pos_migrator` has `USING true` on auth tables.
- Grants (role_table_grants, non-owner): `pos_accounting` = SELECT/INSERT/UPDATE/DELETE on 25 accounting tables (journal, GL, COA, invoices, payments, bank, tax, budgets, rules, e-invoice, templates) and SELECT on `orders, order_items, products, employees, salaries, insurance_tracking, restaurants, tenants, permission_audit_logs`; column-level UPDATE on `salaries(is_paid,paid_at,paid_by,payment_reference,status,updated_at)` and INSERT on `hr_audit_logs` (migration 189); `pos_app` = INSERT/SELECT/UPDATE on `accounting_event_outbox`, INSERT/SELECT on redrive audit and inventory-movement claims, `payments` rows via RLS, role/user-role tables, webhook_retry_queue; `pos_migrator` = auth/membership tables.
- **UNVERIFIED**: in this dev DB `pos_app` has no table grants on `orders/order_items/journal_*` in `role_table_grants` (only the tables listed). Dev falls back to the primary connection when `ACCOUNTING_DB_USER` is unset (cmd/api/accounting_db.go:30-39); the role actually used by the dev API was not verified. Production grants should be checked separately.
- Delete/update behaviour on financial tables for `pos_accounting`: DELETE and UPDATE are GRANTED on `journal_entries`, `journal_entry_lines`, `general_ledger`, `invoices_*`, `payments_*`; protection is trigger-based only. DELETE/UPDATE on `payments_made/received`, `invoices_*`, `bank_accounts_ledger`, `customers_accounting`, `suppliers` have no history triggers (balance fields like `amount_paid`, `current_balance` are mutable).

## 6. Integrity probes (SQL: `sql/10_integrity_probes.sql`, `sql/11_probe_details.sql`; results: `evidence/database/integrity_probes.txt`, `probe_details.txt`)
64 read-only probes. All data probes returned 0 because the tables are empty, except two structural anomalies:
| Probe | Result | Detail |
|---|---|---|
| p31 COA normal_balance vs account_type | **2** | Tenant 1/Restaurant 1: `1240 Accumulated Depreciation` (asset, credit) and `4150 Sales Discounts` (revenue, debit) - legitimate contra-accounts, but Balance Sheet uses `balanceSheetContribution` sign flip (usecase/financial_report_usecase.go:355-364) and P&L treats revenue as credit-debit; the discount account IS used by auto-posting (usecase/auto_posting_engine.go:84-93) and appears as negative revenue. Mixed handling = presentation risk, not a data fault |
| p33 restaurants without chart of accounts | **1** | restaurant id 2 / tenant 2 has 0 accounts; any accounting action for that restaurant would fail with "account not found" (auto-posting resolveAccount, engine:628-643). CoA seeding is not guaranteed at restaurant creation |
| p1..p10 (balance, header vs lines, orphans, GL coverage, COA cache vs GL) | 0 | vacuous (0 journals) |
| p11..p14 (order scope, totals, subtotal vs items) | 0 | vacuous |
| p15..p16 (paid orders vs outbox vs journals, duplicate journals per reference) | 0 | vacuous |
| p17..p22 (invoice balance, payments vs invoice paid, cross-tenant invoice/customer/supplier, journal linkage) | 0 | vacuous |
| p23..p26 (salaries, payroll journals, posting after close) | 0 | vacuous |
| p27..p30, p32, p34 (inventory scope, negative stock, cost math, COA parent scope, bank scope) | 0 | vacuous |
| `cost_price` probe (`sql/12_cost_price_column_probe.sql`) | **ERROR 42703 column p.cost_price does not exist** | confirms menu-profitability / break-even queries cannot execute (see 12-report-formula-catalog.md F1) |

Probe design gaps worth re-running on production-like data: p10 (COA.current_balance vs GL), p16 (duplicate journals per reference - no unique index guards this), p18/p18b (invoice amount_paid vs payments - no DB invariant), p15c/p15d (paid orders without journals).

## 7. Structural findings summary
1. No RLS on any core financial table (section 5).
2. Journal idempotency is not DB-enforced: `idx_je_ref (tenant_id, reference_type, reference_id)` is non-unique; app uses advisory lock + SELECT (repository/general_ledger_repo.go:78-117). Manual journals accept arbitrary `reference_type/reference_id/entry_type` from the request (handler/http/journal_entry_handler.go:151-154,241-244), so a manual entry can occupy an auto-posting identity: `CreateJournalEntryTx` only defaults empty entry_type to "manual" (repository/journal_entry_repo.go:189-191) and performs no reserved-type/reference check (no validation found in handler or repo).
3. Fiscal-period lock is app-only and effectively disabled: no fiscal_periods rows are ever created by code; `ensureFiscalPeriodOpenTx` returns nil when no row exists (repository/journal_entry_repo.go:813-842); void reversals are not checked at all (VoidEntryTx 441-521 never calls it).
4. Dual stock ledgers: `inventory` (legacy, no writers, no tenant column) vs `inventory_movements`; the turnover report reads the dead one.
5. Soft polymorphic links from journals to sources; cascade deletes from `employees` -> `salaries` and `tenants/restaurants` -> `orders` can orphan the immutable ledger's source documents.
