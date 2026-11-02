# 05 - Database and Migrations (Phase 1, read-only)

Scope: backend/migrations, migration runner, schema integrity. Auditor: database/multi-tenancy/sharding.

## 0. Method and limits

- **Live DB not available.** `docker ps` listed no containers (pos-postgres not running; see `evidence/database-evidence/db-status.txt`). No catalog queries were run. Everything below is static analysis of `backend/migrations/*.sql`, Go code and scripts. Trigger/constraint counts are from regex over SQL text, not `pg_catalog`.
- Nothing was modified, migrated or committed.
- Companion file: `06-MULTI-TENANCY-AND-SHARDING.md`. Machine-readable findings: `evidence/issues-database.json` (22 issues, DB-001..DB-022; this file holds the DB-0xx items tagged doc 05).

## 1. Migration inventory

| Metric | Value | Evidence |
|---|---|---|
| Files in backend/migrations | 123 | `evidence/database-evidence/migration-files.txt` |
| Active `.sql` | 107 | same |
| `.disabled` (never run) | 16 (075-080 component/theme-inheritance up/down, 099_create_roles, 102, 104_assign_admin, 105_enhance_tenant_isolation, 106_tenant_audit_logs) | same |
| Total SQL lines (active) | ~6,523 | `wc -l` |
| Distinct tables created (regex) | ~91 (+ ALTERs) | script over `CREATE TABLE` |
| Down/rollback scripts | none active; 7 `*.down.sql.disabled`; one out-of-tree `backend/scripts/rollback_phase6_accounting_hr.sql` | DB-016 |
| Duplicate numeric prefixes | 067, 070, 071, 072, 075, 076, 078, 079, 080, 099, 104, 105, 106 | `duplicate-migration-prefixes.txt` |
| Numbering gaps | 006-011, 018-019, 028-039, 041-049, 051-059, 081-089, 091-098 (and 067 "placeholder" file fills a deliberate gap) | file list |
| Runner | `backend/cmd/migrate/main.go` (filename in `schema_migrations`), executed from `backend/docker-entrypoint.sh:11` on every container start | DB-014 |

Ordering problems: lexical sort with same-prefix tie-break by name (e.g. `070_create_drivers_table` before `070_create_themes_v2_table`; `099_add_test_restaurant` inserts into `themes_v2` which is created at 070, so order works only by accident of prefix). 999_theme_migration_legacy creates `restaurant_themes` last, which `restaurant_theme_repo.go` uses.

Destructive statements: none in forward migrations (grep for DROP TABLE/COLUMN, TRUNCATE, DELETE FROM returned nothing). Seeds/DML are mixed into schema migrations (099, 109, 078_seed, 106, 112, 156, 159). Hardcoded role: `050_create_notifications_system.sql:150-151 TO postgres` (DB-016). Hardcoded credentials seed: DB-013.

## 2. Schema summary (static)

- Hierarchy: `tenants` -> `restaurants` (tenant_id FK) -> `users` (tenant_id NOT NULL, restaurant_id nullable) and operational tables carrying (tenant_id, restaurant_id).
- Core tables: orders (BIGSERIAL, 060), order_items (061, tenant/restaurant added in 065), order_status_history, products/categories/inventory, payments (115, subscription payments via Paymob), RBAC (modules, permission_levels, roles, role_permissions, user_roles), HR (employees, attendance, salaries, leaves, performance_reviews...), accounting (chart_of_accounts, journal_entries/lines, invoices_customer/supplier + lines, payments_received/made, bank_accounts_ledger, general_ledger, fiscal_periods, tax_periods, fixed_assets, auto_posting_rules), KDS, tables, themes (3 models).
- Tenant-bearing tables: ~62 of 91 have `tenant_id`; ~29 do not (list in DB-018). Tenant-level indexes exist for all tenant tables except four child-line tables (DB-018).
- Unique constraints are generally (tenant_id, restaurant_id, number) for accounting (137:42, 143:25) - good. Global uniques that matter for isolation: `orders.order_number UNIQUE` (060:10, DB-010), `tenants.slug/domain/name UNIQUE` (001), `payments.merchant_order_id UNIQUE`.
- Triggers/functions (regex over migrations): updated_at triggers, HR calc triggers (attendance hours, salary totals, leave days/overlap, rating), `validate_employee_role_tenant` (022), invoice line/total recalculation (141/142), journal total recompute (137, misnamed, DB-012), KDS/other. No trigger enforces order totals, state transitions, or posted-entry immutability.
- **RLS: none** in migrations (0 hits); see DB-005.
- Roles: no `CREATE ROLE`/least-privilege app role anywhere; app uses `postgres` superuser (DB-005/DB-020).
- Indexes: broadly present with tenant/restaurant leading columns (066, 121, 129, 157, 160 add more). Duplicates listed in DB-022.

## 3. Data integrity matrix (DB-level vs app-only)

| Domain | DB-level protection (evidence) | App-only | Verdict |
|---|---|---|---|
| Orders | status/payment_status CHECK, totals >=0, FKs with cascade (060:44-52) | totals math (float64), transitions, order number (DB-010/011) | PARTIAL |
| Order items | qty>0, price>=0, FK to orders/products (061) | total=qty*price, tenant/restaurant equals parent (DB-018) | PARTIAL |
| Payments (subscription) | status CHECK, merchant_order_id UNIQUE, partial UNIQUE on success per paymob_order_id (115) | webhook idempotency logic | GOOD for dup-success |
| Accounting payments | NUMERIC(15,2), UNIQUE(tenant,restaurant,number), status CHECK (143) | allocation to invoices | PARTIAL |
| Inventory | `inventory_movements` quantity NUMERIC(12,3) (144); products.quantity_in_stock | stock decrement never done on sale; no CHECK >=0 verified (DB-011) | WEAK |
| Payroll | trigger computes net_salary, CHECK amounts/period, UNIQUE(employee,period) (024) | approval workflow | GOOD |
| Journal entries | line CHECKs (non-negative, one side only, 137:60-63), UNIQUE entry number | **balanced debit=credit only in Go (0.009 tolerance)**, posted immutability (DB-012) | WEAK |
| Invoices | status CHECK, DB-side line/total recalculation triggers (141/142) | payment application | GOOD |
| Subscription status | CHECK on tenants.subscription_status (110:29), history table (111) | transitions/workers | GOOD (duplicate is_active flag, DB-022) |
| Money types | NUMERIC/DECIMAL only; no FLOAT | Go float64 (DB-011/019) | DB sound, app not |

## 4. Findings (database and migrations)

### DB-010 [P1] Order numbers collide and are globally unique; generator ignores tenant and uses unix-time modulo

- **Category:** data-integrity
- **Source:** `backend/internal/domain/order.go`, `backend/internal/usecase/order_usecase.go`, `backend/migrations/060_create_orders_table.sql`
- **Evidence:** order_usecase.go:64 GenerateOrderNumber(tenantID, time.Now().Unix()%1000000); domain/order.go:242-244 ignores tenantID and returns "ORD-<year>-%06d". The value repeats every 1,000,000 s (~11.6 days) and two orders in the same second (any tenant) collide. orders.order_number VARCHAR(50) NOT NULL UNIQUE is GLOBAL (060:10): collisions across tenants cause a unique violation (500) on order creation, and numbers are guessable (see DB-002). Not reproduced dynamically.
- **Suspected root cause:** Placeholder comment at order_usecase.go:62 "In production, this should query the current order count".
- **Recommended fix:** Per-restaurant counter (sequence table with row lock) and UNIQUE(tenant_id, restaurant_id, order_number); add a random component for public tracking.
- **Required regression tests:**
  - Concurrent 100 order creates in the same second: no 500s, all unique
  - Two tenants can hold the same order_number
  - Number does not repeat after 11 days (clock mock)

### DB-011 [P1] Order creation is not transactional; inventory is checked but never decremented; totals computed in float64 with hardcoded 10% tax

- **Category:** data-integrity
- **Source:** `backend/internal/usecase/order_usecase.go`, `backend/internal/repository/order_repo.go`
- **Evidence:** order_usecase.go:178-189 CreateOrder inserts the order, then loops CreateOrderItem as separate statements; order_repo.go has no tx in the create path, so a failure leaves an order with missing items. Inventory: order_usecase.go:130-132 only checks product.QuantityInStock < qty; grep shows quantity_in_stock is written only by UpdateProduct (product_repo.go:364): no decrement on sale, and check-then-act is racy. subtotal := 0.0 float64 (:70), TaxAmount = subtotal*0.10 (:174); discount never validated (:176-178). DB: no CHECK that order_items.total_price = qty*unit_price or orders.total = subtotal+tax-discount+fee (060:46-52 only >=0).
- **Suspected root cause:** Business logic prototyped in the usecase; no DB-level invariants.
- **Recommended fix:** Wrap order+items+inventory decrement in one tx with SELECT ... FOR UPDATE on stock; compute money with integer minor units or decimal; add CHECK constraints; per-restaurant tax rate.
- **Required regression tests:**
  - Fault-injection: item insert failure leaves no order row
  - Concurrent orders for the last unit: exactly one succeeds
  - Total-consistency CHECK rejects tampered rows

### DB-012 [P1] Journal entry balance is not enforced by the database; trigger named validate_journal_balance only sums; posted lines are mutable

- **Category:** data-integrity
- **Source:** `backend/migrations/137_create_journal_entries.sql`, `backend/internal/repository/journal_entry_repo.go`
- **Evidence:** 137:76-114 function validate_journal_balance() computes SUM(debit),SUM(credit) and UPDATEs journal_entries.total_*; it never raises when they differ. No CHECK/constraint trigger on journal_entries(total_debit=total_credit) in 137 or later migrations (grep total_debit in 153-160 returns nothing). Balance is verified only in app code at post time: journal_entry_repo.go:393-395 math.Abs(totalDebit-totalCredit) > 0.009 on float64 (tolerates almost one cent). No trigger prevents UPDATE/DELETE of journal_entry_lines for posted/void entries (only repo-level draft check :272-274 before delete :329). journal_entry_lines.tenant_id (137:56) is not an FK and lines carry no restaurant_id. Row-level trigger recomputes totals per row (O(n^2)).
- **Suspected root cause:** Balance rule implemented only in Go; misleading trigger name.
- **Recommended fix:** Add a DEFERRABLE INITIALLY DEFERRED constraint trigger raising when a posted entry is unbalanced; block writes to lines of posted/void entries; composite FK (journal_entry_id,tenant_id); use decimal in Go and exact equality.
- **Required regression tests:**
  - Direct SQL insert of unbalanced posted entry rejected
  - UPDATE line on posted entry rejected
  - Cross-tenant line insert rejected

### DB-013 [P1] Migration 109 seeds a super admin with a super-admin with a legacy hash and fixed email into every environment; 099 seeds a test tenant/restaurant

- **Category:** security
- **Source:** `backend/migrations/109_seed_super_admin.sql`, `backend/migrations/099_add_test_restaurant.sql`
- **Evidence:** 109_seed_super_admin.sql:5-8 legacy seeded bcrypt hash (literal redacted), email superadmin@pos-saas.com, ip_whitelist NULL = all IPs (:26), is_super_admin true, attached to tenant 1. docker-entrypoint.sh:11 runs ./migrate at every container start including production (docker-compose.prod.yml uses the same image). 099_add_test_restaurant.sql:6-34 inserts "Test Tenant"/"restaurant-1" and a theme in any environment.
- **Suspected root cause:** Seed data placed in schema migrations.
- **Recommended fix:** Remove seeds from migrations; create the super admin via a one-time CLI with an env-provided password; rotate/disable this account on any deployed DB; require IP allowlist/MFA for super admin.
- **Required regression tests:**
  - Fresh prod-profile migrate leaves users table empty
  - Login with the documented default credential fails

### DB-014 [P2] Migration runner lacks locking, checksums, explicit transaction control and rollbacks; runs at every container start

- **Category:** migrations
- **Source:** `backend/cmd/migrate/main.go`, `backend/docker-entrypoint.sh`
- **Evidence:** cmd/migrate/main.go:48-95 sorts filenames lexically (sort.Strings), tracks only the filename in schema_migrations(version), executes each file with one db.Exec and inserts the version row in a separate statement (:85-90): a crash in between re-runs a non-idempotent file. No pg_advisory_lock while docker-entrypoint.sh:11 runs ./migrate on every start, so replicas race. No checksum; no down migrations (7 *.down.sql.disabled never run). Duplicate numeric prefixes: 067,070,071,072,075,076,078,079,080,099,104,105,106 (evidence/database-evidence/duplicate-migration-prefixes.txt); active duplicates: 067x2, 070x2, 071x2, 072x2, 078 (seed + disabled), 099, 104x2 (permission_audit_logs + role_permissions), 105, 106; order relies on lexical tie-break (067_create_themes_table < 067_theme_migrations_placeholder). 123 files, 107 active .sql, 16 .disabled; many numeric gaps (006-011, 018-019, 028-039, 041-049, 051-059, 081-089, 091-098).
- **Suspected root cause:** Home-grown runner; accretive numbering across branches.
- **Recommended fix:** Adopt golang-migrate/goose or add advisory lock + transaction per file + checksum; unique numeric versions; run migrations as a separate job, not in the app entrypoint.
- **Required regression tests:**
  - Two concurrent migrate processes: one waits, schema applied once
  - CI fails on duplicate version numbers
  - Clean DB to head then compare against baseline schema dump

### DB-015 [P3] Dead/misleading and non-idempotent migrations; three coexisting theme models

- **Category:** migrations
- **Source:** `backend/migrations/147_create_performance_reviews.sql`, `backend/migrations/026_create_performance_reviews_table.sql`, `backend/migrations/067_create_themes_table.sql`, `backend/migrations/999_theme_migration_legacy.sql`
- **Evidence:** 147 is CREATE TABLE IF NOT EXISTS performance_reviews with a different column set than 026 (areas_to_improve vs areas_for_improvement; no review_type); once 026 has run it is a silent no-op, and the repository matches 026 (performance_review_repo.go:163-170, columns areas_for_improvement/reviewer_comments), so 147 is dead but misleading (domain/hr_advanced.go:20 uses AreasToImprove). Files without IF NOT EXISTS / ON CONFLICT: 065, 067_placeholder, 073, 078_seed, 105, 106, 112, 120, 122, 126, 127, 132, 134, 135, 155, 156, 159 (DML backfills/seeds mixed into schema migrations). Theme models themes (067), themes_v2 (070) and restaurant_themes (999) coexist; 075-080 theme versions are disabled.
- **Suspected root cause:** Parallel feature branches re-created tables.
- **Recommended fix:** Delete dead 147, consolidate theme models, make backfills idempotent.
- **Required regression tests:**
  - Schema snapshot test after full migrate matches committed expected schema

### DB-016 [P2] Hardcoded role grant and destructive/incomplete rollback script

- **Category:** migrations
- **Source:** `backend/migrations/050_create_notifications_system.sql`, `backend/scripts/rollback_phase6_accounting_hr.sql`
- **Evidence:** 050_create_notifications_system.sql:150-151 GRANT EXECUTE ON FUNCTION ... TO postgres fails if that role does not exist and grants nothing to a least-privilege app role. scripts/rollback_phase6_accounting_hr.sql drops performance_reviews CASCADE (also the table from migration 026, an earlier phase) and does not delete schema_migrations rows, so migrate will skip recreating dropped tables. Forward migrations contain no DROP TABLE/COLUMN/TRUNCATE/DELETE (grep) - good.
- **Suspected root cause:** Rollback written by phase, not by migration.
- **Recommended fix:** Replace with versioned down migrations; remove or parameterise role grants; delete schema_migrations rows in rollback.
- **Required regression tests:**
  - Up/down/up cycle on clean DB succeeds
  - Migration run as a non-postgres role succeeds

### DB-018 [P2] Tenant/restaurant consistency not enforced by composite FKs; denormalised tenant_id can drift; slug uniqueness only per tenant

- **Category:** data-integrity
- **Source:** `backend/migrations/065_add_restaurant_id_to_order_items.sql`, `backend/migrations/137_create_journal_entries.sql`, `backend/migrations/002_create_restaurants.sql`, `backend/internal/repository/restaurant_repo.go`
- **Evidence:** order_items.tenant_id/restaurant_id are separate FKs (065:24-26) with nothing tying them to orders; same pattern in journal_entry_lines.tenant_id (no FK at all, 137:56), invoice_*_lines, bank_reconciliation_lines. Tables with no tenant_id column (script over migrations): order_status_history, product_images, inventory, low_stock_alerts, driver_assignments, driver_location_history, theme_sections, theme_history, user_settings, etc. Four tenant-bearing tables lack a leading tenant/restaurant index: bank_reconciliation_lines, invoice_customer_lines, invoice_supplier_lines, journal_entry_lines (indexed by parent id). restaurants UNIQUE(tenant_id,slug) (002:30) while public lookup is slug-only (restaurant_repo.go:20-33 WHERE slug=$1, no LIMIT/ORDER): global uniqueness currently comes only from tenants.slug UNIQUE at registration (auth_repo.go:337) - tenants.max_restaurants is app-only.
- **Suspected root cause:** tenant_id added as a denormalised column without UNIQUE(id,tenant_id) parent keys.
- **Recommended fix:** Add UNIQUE(id,tenant_id[,restaurant_id]) on parents and composite FKs on children; enforce global slug uniqueness or make lookup tenant-aware.
- **Required regression tests:**
  - INSERT child with mismatched tenant rejected by DB
  - Duplicate restaurant slug across tenants rejected or disambiguated

### DB-019 [P2] Money types: DB numeric is sound but Go uses float64; id types mixed; ON DELETE CASCADE on orders

- **Category:** data-integrity
- **Source:** `backend/internal/usecase/order_usecase.go`, `backend/internal/repository/journal_entry_repo.go`, `backend/migrations/060_create_orders_table.sql`, `backend/migrations/115_create_payments_table.sql`
- **Evidence:** Money columns across migrations (grep): 54 NUMERIC, 24 DECIMAL, 5 INTEGER (payments.amount_cents INTEGER, 115:7) and no FLOAT/REAL/MONEY columns - sound; precision differs (orders DECIMAL(10,2) vs accounting NUMERIC(15,2)). Go uses float64 for subtotal/tax (order_usecase.go:70,174) and ledger deltas (journal_entry_repo.go:777). orders/order_items cascade-delete from restaurants/tenants (060:44-45), so deleting a restaurant erases sales history, whereas accounting tables use plain REFERENCES. Ids: orders BIGSERIAL, accounting SERIAL (INT) - e.g. journal_entries.reference_id INT -> orders.id BIGINT.
- **Suspected root cause:** No money-type convention.
- **Recommended fix:** Use decimal or minor units end to end; BIGINT everywhere; soft-delete tenants/restaurants and RESTRICT FKs on financial records.
- **Required regression tests:**
  - Round-trip decimal tests on totals
  - FK RESTRICT test: cannot delete a restaurant that has orders

### DB-020 [P3] Payment provider auth token stored in plaintext; DB defaults to superuser and sslmode=disable

- **Category:** security
- **Source:** `backend/migrations/115_create_payments_table.sql`, `backend/internal/config/config.go`, `docker-compose.prod.yml`, `.env.prod.example`
- **Evidence:** 115:23 paymob_auth_token TEXT stored as received; only card_last_four/brand are kept (no PAN - good). Positive: partial UNIQUE idx_payments_paymob_order_success (115:48-50) prevents duplicate success per Paymob order; merchant_order_id UNIQUE (:35). docker-compose.prod.yml:63 and .env.prod.example:16 default DB_SSLMODE=disable; DB_USER defaults to the postgres superuser (config.go:54).
- **Suspected root cause:** Defaults optimised for local dev.
- **Recommended fix:** Do not persist short-lived auth tokens (or encrypt); require sslmode=require/verify-full in prod; create a least-privilege app role.
- **Required regression tests:**
  - Prod config validation rejects sslmode=disable and superuser

### DB-022 [P3] Duplicate indexes and weak state-transition protection

- **Category:** database
- **Source:** `backend/migrations/003_create_users.sql`, `backend/migrations/066_add_restaurant_id_indexes.sql`, `backend/migrations/060_create_orders_table.sql`
- **Evidence:** Same-name duplicates (second is a no-op IF NOT EXISTS in 066): idx_users_restaurant, idx_categories_restaurant, idx_products_restaurant, idx_orders_restaurant_created, idx_order_items_restaurant(+_created); distinct-name duplicates: idx_performance_reviews_tenant vs _tenant_restaurant, idx_attendance_date vs _date_range, theme_history created_at indexes. Order status/payment_status only have value CHECKs (060:46-47); no trigger enforces legal transitions or that payment_status=paid has a payments row; transitions are app-only (order_repo.go:471-540). Positives: tenants.subscription_status CHECK (110:29); salaries trigger computes net_salary and CHECKs amounts (024:93-97,148). tenants.is_active coexists with subscription_status (auth_repo.go:63) as a second source of truth.
- **Suspected root cause:** Incremental migrations.
- **Recommended fix:** Drop duplicates; add a transition trigger or state-machine table; unify tenant active flag.
- **Required regression tests:**
  - Illegal status transition via SQL rejected
  - Tenant disabled flag consistent with subscription_status


## 5. Not verified

- No live DB: no pg_catalog counts, no EXPLAIN, no actual trigger/constraint list; actual `schema_migrations` state of any environment unknown.
- No dynamic requests: DB-001/002/003/004 exploitability is from code reading and must be confirmed in Phase 2.
- Sharding/scaling tests were not run (need live shards).
- Authorization (who may call RBAC/driver/theme routes) is cross-checked only for route wrappers; role semantics belong to the auth auditor.
- Upload serving path, payment webhook/worker SQL scoping, hr_report dynamic where clauses: not fully traced.
