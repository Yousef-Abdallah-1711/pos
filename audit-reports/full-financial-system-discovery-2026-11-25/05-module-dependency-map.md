# 05 - Module Dependency Map

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Import edges computed by grep of `pos-saas/internal/*` imports in non-test Go files (counts = import statements). Read-only.

## 1. Backend package dependency (VERIFIED)

| From \ To | domain | repository | usecase | service | middleware | pkg/jwt | pkg/database | config | observability | commerce | other |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `handler/http` (60 files) | 57 | 62 | 43 | 19 | 83 | 30 | 3 | 4 | - | 1 | auth 2, validation 1, csvutil 1 |
| `usecase` (22) | 36 | 34 | - | 4 | - | 3 | 1 | - | 2 | 2 | |
| `repository` (66) | 92 | - | - | - | - | - | 13 | 1 | 1 | 1 | csvutil 1 |
| `service` (12) | 5 | **1** | - | - | - | - | - | 1 | - | - | |
| `middleware` (20) | 3 | **5** | - | 5 | - | 10 | - | 2 | 2 | 1 | sharding 1 |
| `domain` (31) | - | - | - | - | - | - | - | - | - | - | leaf |

Observations
- No import cycles; `domain` is a leaf; `repository` never imports `usecase`/`handler`/`service` (clean downward flow for SQL).
- Handlers depend on BOTH usecases (43 imports) and repositories (62): the usecase layer is optional, not mandatory. Journal, invoice, payment, bank, supplier, customer, tax, HR handlers go straight to repositories (main.go:423-440).
- `service -> repository` (1) and `middleware -> repository` (5: subscription repo for tenant status, user-role/role-permission repos for RBAC, impersonation) put persistence access into cross-cutting layers.
- Usecases take concrete `*repository.X` types (e.g. `NewOrderUseCase(orderRepo *repository.OrderRepository, ...)`, order_usecase.go:29) - no repository interfaces, so there is no seam to swap a pool at the usecase boundary; the pool choice is entirely a `main.go` wiring decision.
- `OrderUseCase` also owns `*repository.AccountingEventRepository` and an `AutoPostingService` interface (order_usecase.go:24-25) - the only interface seam in the posting path (`AutoPostingService`, auto_posting_engine.go:16-20).

## 2. Financial module dependency map (VERIFIED wiring, main.go)

```
Orders/POS ──► OrderUseCase ──► OrderRepository(db) ──► orders, accounting_event_outbox (same tx)
                    │ (worker) ProcessNextAccountingEvent
                    ▼
             AutoPostingEngine ──► AutoPostingRepository(accountingDB)   [rules]
                    │           ├─► AccountRepository(accountingDB)       [account resolve]
                    │           ├─► ProductRepository(db)                 [product, tracks stock]
                    │           ├─► InventoryMovementRepository(db)       [cost + sale movement]
                    │           └─► GeneralLedgerRepository(accountingDB) ──► JournalEntryRepository(accountingDB)
HR Payroll ──► HRSalaryHandler ─(tx on accountingDB)─► SalaryRepository.MarkAsPaidTx + AutoPostingEngine.PostPayrollEntryTx
Invoices/Payments ──► Customer/SupplierInvoiceRepository, PaymentAccountingRepository (accountingDB) ──► JournalEntryRepository + AccountRepository (+BankAccountRepository)
Journal/Ledger/COA ──► JournalEntryHandler/LedgerHandler/AccountHandler ──► JournalEntryRepository / GeneralLedgerRepository / AccountingUseCase->AccountRepository (accountingDB)
Tax ──► TaxHandler ──► TaxRepository(accountingDB)   (no journal)
Advanced (recurring, CSV, FX, budgets, bank rec, e-invoice, KPIs) ──► AccountingAdvancedHandler(db, repo(accountingDB), customerInvoiceRepo)
Financial reports ──► FinancialReportHandler ──► FinancialReportUseCase(db, GeneralLedgerRepository(accountingDB))   [mixed pools]
Operational reports ──► ReportHandler(db)  (raw SQL in handler)
Inventory adjust ──► InventoryAdjustmentHandler(db)  (raw SQL; no journal)
SaaS billing ──► PaymentUseCase(paymentRepo, subscriptionRepo, planRepo, paymobClient, db) ──► Paymob (external) / WebhookRetryWorker
```

Pool assignment (main.go:194-211, 331, 335, 421): accountingDB for journal, accounts, auto-posting rules, GL, suppliers, customers, bank accounts, invoices, payment accounting, tax, accounting-advanced; `db` for inventory movements, outbox, orders, products, HR (non-payroll-tx), reports, notifications, RBAC, subscriptions.

## 3. Cross-module coupling that crosses the pool boundary

| Coupling | Evidence | Risk |
|---|---|---|
| COGS journal (accountingDB) then inventory movement (db) | auto_posting_engine.go:222-258 | two txs; replay-safe via `accounting_inventory_movement_events` (184) |
| Outbox claim (db) then journal (accountingDB) then outbox complete (db) | order_usecase.go:482-512 | three txs; lease fencing (183) |
| Financial reports read GL through `generalLedgerRepo` (accountingDB) but also raw SQL on `db` (`financial_report_usecase.go:187-194, 377-383`) | | reads may mix sources; grants for `pos_app` on ledger tables UNVERIFIED |
| `pos_accounting` reads orders/products/salaries (182:30-37) | | read-only cross-domain |
| Payroll tx on accountingDB updates `salaries` columns | 189; hr_salary_handler.go:521-557 | `pos_accounting` can set `is_paid`, `paid_by` |

## 4. Frontend dependency map

| App | Shared packages (package.json) | Talks to backend via |
|---|---|---|
| dashboard | `@pos-saas/form-foundation`, `ui-components`, `ui-themes` (+`i18n`, `theme` via tsconfig paths) | axios `src/lib/api.ts` -> `NEXT_PUBLIC_API_URL` (default `http://localhost:8080/api/v1`, config.ts:15); cookie auth |
| landing | `form-foundation` | `src/lib/config.ts` `getApiBaseUrl` (requires `NEXT_PUBLIC_API_URL` in prod) |
| restaurant-website | `ui`, `ui-themes` | fetch + react-query (`lib/hooks/use-api-queries.ts`), server-side theme fetch (`lib/api/get-theme-server`) |
| super-admin | `ui-components`, `ui-themes` | `getSuperAdminApiBaseUrl()` -> `/api/v1/super-admin` |

Inter-app flow: landing registers/logs in tenants then redirects to dashboard; super-admin can mint dashboard login tokens/impersonation (`POST /super-admin/tenants/{id}/dashboard-login`, `/super-admin/impersonate/{userID}`, main.go routes in 08 CSV); restaurant-website resolves tenant by subdomain/slug (`middleware.ts`, public routes `/public/restaurants/{slug}/...`).

## 5. UNVERIFIED / missing connections
- Frontend theme-builder -> backend component/import/export routes: no backend route (06 mismatches list).
- Frontend UI -> outbox failed-events/redrive API: no UI.
- Frontend tax/currency configuration -> backend: no setting exists; both hard-coded (10% backend, 14% POS display).
- Whether any scheduler triggers recurring journal templates automatically: no worker found (only manual `POST .../run`); UNVERIFIED by exhaustive goroutine search - only three workers are started in main.go.
- Email delivery (`send-email`), e-invoice submission external endpoints: handler-level only; external provider dependency UNVERIFIED.
