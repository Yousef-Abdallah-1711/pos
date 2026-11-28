# 03 - Complete System Architecture (verified, layer by layer)

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Static, read-only analysis. Legend: VERIFIED = read in source with file:line; UNVERIFIED = not proven in this phase. Companion files: `08-api-inventory.csv` (374 routes), `06-complete-page-inventory.md`, `13-accounting-event-architecture.md`, `diagrams/`.

## 0. Component inventory

| Component | Tech | Evidence |
|---|---|---|
| Backend API | Go, stdlib `net/http` ServeMux (Go 1.22 method+path patterns), single binary `backend/cmd/api` | main.go (2066 lines), `backend/go.mod` |
| Apps | 4 Next.js apps: `dashboard` (staff: POS, accounting, HR, reports, 79 pages), `landing` (marketing + registration/login, 20), `restaurant-website` (customer storefront, 15 incl. route handler), `super-admin` (platform ops, 12) | `frontend/apps/*` |
| Shared packages | `frontend/packages/{form-foundation,i18n,theme,ui,ui-components}` | directory listing |
| DB | PostgreSQL; migrations `backend/migrations/*.sql` (145 .sql files, 163 entries incl. rollbacks; duplicate numeric prefixes at 180, 181, 182) | `ls backend/migrations` |
| Cache/session | Redis (`REDIS_URL`) mandatory in production for revocable sessions and auth rate limits; in-memory fallback in non-production | main.go:144-160, 548-555 |
| Realtime | KDS websocket hub | main.go:164-165, 1457 |
| Payments (platform billing) | Paymob live/mock client; webhook retry worker | main.go:259-276, 489-504 |
| Workers (in-process goroutines) | accounting outbox worker, subscription expiry worker, webhook retry worker, DB health monitors | main.go:501-518, 106, 132 |

## 1. Layer-by-layer

### L1 Frontend UI (pages/components)
- Responsibility: forms, validation hints, RBAC-driven UI hiding (`useUserPermissions().checkPermission(MODULE_IDS.X, PERMISSION_LEVELS.Y)`, `dashboard/src/lib/rbac.ts:11-45`). 
- Input/output: user events -> API client calls; renders `data` envelopes.
- Security boundary: UI checks are advisory only. Many financial pages have NO client permission check (all `/dashboard/reports/financial/*`, `accounting/advanced`, `inventory`; see 06) and rely entirely on backend module checks.
- Financial side effects: computes invoice line `tax_amount`/`line_total` client-side (InvoiceForm.tsx) but the DB trigger recomputes (141_create_invoices_customer.sql:76-84, customer table VERIFIED; supplier trigger UNVERIFIED). POS shows 14% tax (`stores/posStore.ts:107`) while the server prices at 10% (`backend/internal/usecase/order_pricing.go:11,46-51`) - display/ledger mismatch.

### L2 Frontend API client
- Dashboard: axios instance `src/lib/api.ts` (`baseURL` from `lib/config.ts`, `withCredentials: true`, default timeout 6000 ms, GET de-duplication keyed by identity scope, response interceptor at :194). Resource modules in `src/lib/api/*.ts` (`accountingApi.ts`, `posApi.ts`, `financialReportsApi.ts`, ...). Hooks in `src/hooks`.
- Auth: host-only HttpOnly session cookie (browser cannot read it); no bearer token in JS (backend auth.go:28-50). 
- Retry: no retry on write calls found; POS offline queue replays failed orders (`usePOSOffline.ts:165-200`) with NO idempotency key; client timeout 6s can abort a request the server still completes (UNVERIFIED consequence). 
- Restaurant-website: fetch/react-query hooks `lib/hooks/use-api-queries.ts`; sends `expected_total` from a server quote (checkout-page-content.tsx:170).

### L3 HTTP routing
- `mux := http.NewServeMux()` main.go:532; 374 registrations (parity-tested by `routes_manifest_test.go` against `routes.manifest`); registered in `main.go`, `pos_routes.go:16`, `health_routes.go`, `translation_routes.go`. 4 commented-out routes excluded (main.go:618, 1544-1546).
- Global chain (outermost first): RequestID -> Recovery -> RequestMetrics -> SecurityHeaders -> Gzip -> CORS -> NoStoreAPIResponses -> mux (main.go:1982-1989). Server timeouts: ReadHeader 5s, Read 30s, Write 60s, Idle 90s (main.go:1999-2005).
- 46 public (no auth) routes, 302 tenant-session routes, 26 super-admin routes (computed in `tests/gen_api_inventory.py`).

### L4 Middleware (per-route wrappers, main.go:675-760)
Order for `wrapWithPermission` (execution): `authenticateRequest` (session cookie -> claims; unsafe methods require trusted Origin, auth.go:28-50) -> `tenantContext` (tenant/restaurant from claims; `X-Tenant-ID` header ignored if it differs, tenant.go:57-105) -> `statusGuard` (`TenantStatusMiddleware`: blocks expired/disabled tenants, 403 `TENANT_ACCESS_BLOCKED`) -> `requirePermission(moduleID, level)` -> handler.
- Variants: `wrapProtectedNoStatus` (billing, /me), `wrapAuthPermission` (RBAC admin endpoints without tenant context), `wrapKDSWithPermission` (+ per-user rate limits 60/min, bump 5/s), `wrapSuperAdmin` (RequireSuperAdmin + IP whitelist + super-admin JWT + 30/min limit, main.go:1900-1906), public rate limiting by IP (orders write 20/min, read 60/min; main.go:34-35, 598-611).
- RBAC enforcement mode: `enforce` in production; **default `audit` (log-only) when `ENV` is not production** (`middleware/permission.go:36-53`). `CurrentEnv()` defaults to `"development"` when none of ENV/APP_ENV/GO_ENV/NODE_ENV is set (`config/runtime_env.go:11-18`), so a deployment missing `ENV=production` silently runs unenforced RBAC, optional Redis, and the accounting-pool fallback. `docker-compose.prod.yml:13,25` sets `ENV: production` and `RBAC_ENFORCEMENT: enforce` (VERIFIED); that compose file has no `ACCOUNTING_DB_*` keys (supplied elsewhere? UNVERIFIED; `ops/aws/compose.json:153-154` does set them).
- Of the 302 tenant-session routes, the subset with no permission module (identity+tenant only): /me, users/{id}/profile|settings*, billing routes, /rbac/me*, etc. (see CSV `permission=none`).

### L5 HTTP handlers (`backend/internal/handler/http`, 60 files)
- Responsibility: decode JSON (`decodeJSONStrict`), pull `tenantID/restaurantID/userID` from context, call usecase or repository, shape `{success,data|error}`.
- Layering is inconsistent: 10 handler files import `database/sql`; 7 run SQL directly (`inventory_adjustment_handler.go`, `report_handler.go`, `hr_report_handler.go`, `hr_attendance_handler.go`, `hr_leave_handler.go`, `accounting_advanced_handler.go`, `payment_monitoring_handler.go`); `hr_salary_handler.go` owns the payroll DB transaction; 39 routes are anonymous closures in `main.go` (users CRUD, hr/roles, notifications stubs, tenant contact-info, rbac stubs) with inline SQL (e.g. main.go:829).
- Error handling: mixed; `writeAccountingMutationError` maps any error containing "not found" to 404 else 400 (`accounting_phase3_handler_helpers.go:51-60`) - business-rule violations and DB failures both surface as 400 with the raw `err.Error()` text.

### L6 Usecase (`internal/usecase`, 22 files)
Used for: auth, product, notification, order, table, KDS, accounts (`AccountingUseCase`), financial reports, auto-posting engine, super-admin, impersonation, subscription, payment. **Not used** for journals, invoices, payments, bank, suppliers, customers, tax, HR (handlers call repositories directly, main.go:424-440).
- Order flows: `OrderUseCase.CreateOrder` (order_usecase.go:62; server reprices from DB, `order_pricing.go:13-39`; `payment_status` forced `pending`, :101), `UpdatePaymentStatusAs` (:455), `ProcessNextAccountingEvent` (:482).
- `AutoPostingEngine` (`auto_posting_engine.go:26-258`): rule lookup, account resolution, line building for sales, COGS, payroll.

### L7 Services (`internal/service`, 12 files)
Session store, Redis cache, Paymob clients, webhook retry queue/worker, subscription expiry worker, `AccountingEventWorker` (accounting_event_worker.go). No accounting business logic lives here except the worker loop.

### L8 Repositories (`internal/repository`, 66 files)
- Own all SQL. Tenant scoping is by explicit `tenant_id/restaurant_id` predicates in each query (no ORM, no global scope). RLS is applied to only 13 migrations' worth of ancillary tables (164-172, 174, 176, 177, 191) and used through `database.WithTenantTx` (`pkg/database/tenant_tx.go:11-30`); NO RLS on `orders`, `journal_entries`, `journal_entry_lines`, `general_ledger`, invoices, payments, salaries, or the outbox (grep of `ENABLE ROW LEVEL SECURITY`). 
- Transaction boundaries are repository-local: `CreateInvoice` (invoice_customer_repo.go:128-158), `ApproveInvoice` (invoice_supplier_repo.go:295), `CreatePaymentReceived/Made` (payment_acct_repo.go:29, 175), `PostEntry` (journal_entry_repo.go:344), `UpdatePaymentStatusAs` (order_repo.go:545).
- Retry: none inside repositories; retry exists only in the outbox worker and webhook queue.

### L9 SQL / DB
- Triggers: CFR-002 guards on `journal_entries`, `journal_entry_lines`, `general_ledger`, `chart_of_accounts`, `restaurants` (180_enforce_journal_integrity.sql:471-496: header/line write guards, GL insert guard, GL immutable); CFR-003 normal-balance guard and account lock (190); invoice line/total recompute triggers (141:76-140).
- Roles: `pos_migrator` (DDL), `pos_app` (operational), `pos_accounting` (restricted ledger writer), plus a separate optional read-only auth-lookup role (auth_lookup_db.go; login path disabled until reviewed per main.go:123-125).

## 2. Database connection pools and fallback (VERIFIED)

| Pool | Role | Config | Used by |
|---|---|---|---|
| `db` | `pos_app` (production guard `requireRestrictedRuntimeRole`, startup_guards.go:50-66 - must be `pos_app`, not superuser/BYPASSRLS) | `DB_MAX_OPEN_CONNS` default 75, idle 25 (config.go:78-79) | all operational repos (auth, products, orders, KDS, HR except payroll tx, notifications, RBAC, subscriptions, payments), outbox repo, inventory movement repo, reports (`reportHandler`, `FinancialReportUseCase(db, ...)` main.go:335,411) |
| `accountingDB` | `pos_accounting`, forced `MaxOpenConns=5`, `MaxIdleConns=2` (accounting_db.go:15-26) | requires `ACCOUNTING_DB_USER == "pos_accounting"` and a password; startup verifies `current_user=pos_accounting` and NOT superuser/createdb/createrole/bypassrls (accounting_db.go:47-55) | repos built at main.go:194-211: accounts, journals, auto-posting rules, GL, suppliers, customers, bank accounts, customer/supplier invoices, payment accounting, tax, accounting-advanced; payroll tx (`SetPayrollTransactionDatabase(accountingDB)` main.go:421) |
| `authLookupDB` | separate read-only role | optional | login path disabled pending review |

Fallback behaviour (accounting_db.go:28-38): when `db == nil` in dev/test returns `nil, nil` (accounting repos then hold a nil pool: stub mode; production refuses to start, startup_guards.go:30-32 and accounting_db.go:30-34); when env is `development|dev|local|test` AND `ACCOUNTING_DB_USER`/`ACCOUNTING_DB_PASSWORD` are both unset, **the accounting pool IS the operational pool** (`return primary, nil`). After migration 182 revoked journal/GL INSERT from `pos_app`, that dev fallback cannot post journals (UNVERIFIED by execution: expect permission-denied), so a developer environment without accounting creds cannot exercise posting; conversely tests that pass in fallback mode do not prove the production role split. In production the fallback is disallowed (non-fallback branch builds the restricted pool).

Capacity consequences (derived): all ledger writes, invoice/payment/payroll transactions and the outbox-driven posting share 5 connections per replica; each `PostEntryTx` holds a journal row lock then account row locks (cfr003 trigger comment, 190:54-55).

Cross-pool reads: `pos_accounting` has SELECT-only on `orders, order_items, products, employees, salaries, tenants, restaurants, insurance_tracking, permission_audit_logs` and column-level SELECT on `users(id,name,email,tenant_id,restaurant_id)` (182:30-39); update on salaries payment columns + `hr_audit_logs` insert (189).

## 3. Tenant isolation (VERIFIED parts)
1. Identity: session cookie -> server session store (Redis) -> claims (tenant, restaurant, user, roles).
2. Context: `tenantContext` trusts claims over headers; if claims lack a restaurant, the `X-Restaurant-ID` header is accepted only after `SQLRestaurantVerifier` confirms it belongs to the tenant (tenant.go:76-88; main.go:589-593). Public routes accept `X-Tenant-ID`/`X-Restaurant-ID` headers verified against the DB (tenant.go:92-102).
3. Data: every repository query carries `tenant_id` (+`restaurant_id`) predicates (spot-verified in outbox, journal, GL, invoice, payment repos). DB-level backstops: triggers rejecting cross-tenant journal lines / GL accounts (180:32-44, 190:37-48), RLS on ancillary tables only.
4. Gaps: no RLS on core financial tables; the outbox worker and reconciler are intentionally cross-tenant; `ClaimNext` does not set a tenant context.

## 4. Where accounting posting actually happens (PROOF)

All code that creates journal entries (exhaustive grep for `CreateJournalEntry(Tx)` and `INSERT INTO journal_entries` over `backend/`, excluding tests): the only SQL insert is `journal_entry_repo.go:213` (`INSERT INTO journal_entries`) and lines at :717 for lines; every caller goes through `JournalEntryRepository`. Call sites:

| # | Path | Caller (file:line) | Sync/async | Tx scope | Engine used | Posted state |
|---|---|---|---|---|---|---|
| 1 | Paid order sales + COGS | `AutoPostingEngine.PostSalesEntry/PostCOGSEntry` (auto_posting_engine.go:50,143) <- `OrderUseCase.AutoPostPaidOrder` (order_usecase.go:469-479) <- worker `ProcessNextAccountingEvent` (:482) | ASYNC via durable outbox | own tx per journal (`CreateAndPostAutoEntry`, general_ledger_repo.go:37-134) | AutoPostingEngine + auto_posting_rules | create+post |
| 2 | Payroll payment | `HRSalaryHandler.MarkAsPaid` hr_salary_handler.go:476-565 -> `PostPayrollEntryTx` (auto_posting_engine.go:268) | SYNC | single handler-owned tx incl. salary status update (accounting pool) | AutoPostingEngine (tx variant); disabled rule aborts the payment (auto_posting_engine.go:288-291) | create+post |
| 3 | Customer invoice create | `CustomerInvoiceRepository.createInvoiceTx` -> `syncJournalEntryTx` invoice_customer_repo.go:522, 539-560 | SYNC | same tx as invoice + balance | repo-embedded (builds lines itself, `buildJournalEntry` :561) | create+post immediately |
| 4 | Customer invoice void | `VoidInvoice` invoice_customer_repo.go:323-362 -> `VoidEntryTx` | SYNC | same tx | repo-embedded | void |
| 5 | Supplier invoice approve | `ApproveInvoice` invoice_supplier_repo.go:295-348 (`CreateJournalEntryTx` :317) | SYNC | same tx | repo-embedded (`buildJournalEntry` :498) | create+post |
| 6 | Supplier invoice void | `VoidInvoice` invoice_supplier_repo.go:350 | SYNC | same tx | repo-embedded | void |
| 7 | Payment received / made | `CreatePaymentReceived` payment_acct_repo.go:29 (journal :129); `CreatePaymentMade` :175 (journal :275) | SYNC | same tx as payment + invoice apply | repo-embedded | create+post |
| 8 | Manual journal | `JournalEntryHandler.CreateJournalEntry` journal_entry_handler.go:162 -> `CreateJournalEntry` (draft); `PostJournalEntry` -> `PostEntry` | SYNC, two calls | create tx; post tx | none | draft then post |
| 9 | Recurring journal run | `AccountingAdvancedRepository.RunRecurringJournalTemplate` accounting_advanced_repo.go:396-480 (create :459, `PostEntry` if `AutoPost`) | SYNC | NOT atomic (create then post are separate txs) | none | draft or posted |
| 10 | CSV journal import | `CreateImportedJournalEntry` accounting_advanced_repo.go:1921-1970 (create :1964; post if `autoPost`) | SYNC | NOT atomic | none | draft or posted |
| 11 | Void manual entry | `VoidJournalEntry` -> `VoidEntryTx` | SYNC | tx | none | void |

Dead engine methods (no callers): `PostSupplierInvoiceEntry`, `PostSupplierPaymentEntry`, `PostDepreciationEntries` (auto_posting_engine.go:462-513). Not a posting path.

**Conclusion (verified):** posting is NOT a single central engine and NOT a single outbox. It is multiple competing mechanisms:
1. Durable transactional outbox + async worker + `AutoPostingEngine` (orders; sales and COGS only) - eventual consistency.
2. Synchronous, in-transaction posting through the same engine (payroll).
3. Synchronous, in-transaction posting embedded in operational repositories that build their own journal lines and bypass `auto_posting_rules` (customer/supplier invoices, payments).
4. User-driven/batch journal creation with non-atomic create+post (manual, recurring, CSV import).
All converge on `JournalEntryRepository.CreateJournalEntryTx/PostEntryTx` (journal_entry_repo.go:182, 358) which enforces balance (tolerance 0.009), fiscal period open check (fail-open if no period row, :825-826), appends GL rows and is backed by DB triggers (180). No async event bus other than the single outbox; no message broker.

Inconsistency worth noting: account selection differs per path - engine uses rule config + hard-coded fallbacks `4100, 2120, 4150, 4120, 5110, 1130, 5210` (auto_posting_engine.go:~69-250, 304) while invoices/payments use customer/supplier-linked or "system default" accounts (`lookupCustomerReceivableAccountTx` invoice_customer_repo.go:638, `resolveCollectionAccountTx` payment_acct_repo.go:546).

## 5. Transaction boundaries, error handling, retry (summary)

| Concern | Finding |
|---|---|
| Request timeout | server Write 60s; client 6s default; repository calls mostly use `context.Background()` with 5s timeouts for order/outbox paths (order_repo.go:551-553) and no timeout for many accounting repos (`r.db.Begin()`), UNVERIFIED impact |
| Error contract | JSON `{success:false,error:"..."}` with raw error strings in many accounting handlers |
| Idempotency | only the outbox event_id/UNIQUE, auto-posting advisory lock + existence check, inventory event key, and `order_cancellation_handler` (the only handler referencing "idempotency"); no Idempotency-Key for payments/invoices/orders |
| Retry | outbox (8 attempts, exp. backoff), webhook queue (1m,5m,30m,2h,12h; main.go:492); none elsewhere |
| Concurrency | `SELECT ... FOR UPDATE` on orders/journals; advisory lock per auto-posting source; GL insert trigger locks account |
| Observability | `observability.RecordAccountingEvents`; `[AccountingRecovery]`/`[AccountingAutoPost]` log lines |

## 6. Operational vs ledger consistency gaps (derived)
- Revenue reaches the GL only for orders flipped to `paid` via the single admin endpoint, and the frontend only does this for POS cash (13 section 3). Card/online/wallet methods are accepted at creation (order_usecase.go:605) but no gateway integration exists for restaurant orders (Paymob is platform subscription billing only, main.go:259-276).
- Inventory adjustments change stock and movement rows with no unit cost and no journal (`inventory_adjustment_handler.go:81-95`); purchase/waste/return therefore never hit the GL; inventory valuation report reads operational data (`FinancialReportUseCase(db,...)`) - GL vs report divergence is possible. UNVERIFIED by report-level reconciliation in this phase.
- Reports use the operational pool while postings use the ledger pool: reading a just-posted entry works (same DB), but `pos_app` has SELECT on ledger tables only by default grants (UNVERIFIED; reports evidently work in tests).
- Hard-coded tax 10% (order_pricing.go:11) and currency `EGP` (order_pricing.go:~99; invoice defaults) - there is no tenant tax/currency settings table/UI (settings/restaurant form has name/email/phone only).

## 7. Top architectural risks (ranked)
1. Silent "completed" events when the sales/COGS rule is inactive (auto_posting_engine.go:59,152; accounting_event_repo.go:107-112 only counts it).
2. Environment fail-open default: unset `ENV` => RBAC audit-mode, optional Redis, accounting-pool fallback (runtime_env.go:11-18; permission.go:36-42; accounting_db.go:29,36).
3. Revenue recognition depends on a client-side two-call sequence for POS cash only; card/other orders never reach the GL; failure swallowed (pos/page.tsx:302-308) and offline replay can duplicate orders (usePOSOffline.ts:176-190); no idempotency keys.
4. Mutually inconsistent posting mechanisms and account resolution (section 4); non-atomic create+post for recurring/CSV (accounting_advanced_repo.go:459-465, 1964-1970); no DB unique constraint on auto-posting source.
5. Segregation-of-duties gaps: payroll mark-paid and order payment-status need only SALARIES/ORDERS WRITE (main.go:1437, 1096); outbox redrive needs only ACCOUNTING WRITE; tenant isolation of core financial tables is application-only (no RLS).
