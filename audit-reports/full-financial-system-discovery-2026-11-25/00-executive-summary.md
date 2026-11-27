    # 00 — Discovery Completion Report (executive summary)

    **Phase:** Discovery / understanding only. No product code, git state, docker state or DB records were changed. No transactions, form submissions or mutating API calls were made.
    **Baseline:** `Yousef-Abdallah-1711/pos` `main` @ `133260db66b1143c254d79476daa7650b41a33f0` (re-verified: local HEAD == upstream `main`; local `origin` still points at the obsolete `yousefabdallah171/POSS`, unchanged).
    **Evidence labels used throughout:** CONFIRMED-BY-SOURCE · CONFIRMED-BY-READONLY-RUNTIME · DOCUMENTED-INTENT · INFERRED · UNKNOWN · REQUIRES-EXECUTION. *Nothing in this package is runtime-executed financial evidence; every behavioural conclusion below is "statically established" until the execution phase.*
    **Release verdict at this point: NO-GO (default).** Discovery cannot satisfy the acceptance gates; several gates are also contradicted by static findings.

    ## 1. What the system is
    A multi-tenant restaurant SaaS: a **Go** API (`net/http` mux, ~374 routes, layered handler → usecase → repository → raw SQL on PostgreSQL 15) and four **Next.js 15** apps — Landing (auth/marketing, :3001), Dashboard (tenant back-office incl. POS, KDS, HR, accounting, reports, :3002), Restaurant Website (public per-tenant storefront, :3003), Super Admin (platform, :3004). Tenancy comes from the JWT; authorization is module-level RBAC (`wrapWithPermission(module, level)`); subscription/trial gating wraps tenant routes. Docs: [03](03-complete-system-architecture.md), [05](05-module-dependency-map.md), [diagrams/](diagrams/).

    ## 2. How money actually flows (statically established)
    Accounting is **not one engine**. Several mechanisms converge on `JournalEntryRepository.CreateJournalEntryTx` / `PostEntryTx`:
    1. **Paid orders** → durable outbox (only event type `paid_order`, written in the same tx as the payment-status flip) → 15 s worker → auto-posting rules → sales/COGS journals. Eventually consistent; an inactive rule completes the event with **no journal**.
    2. **Payroll** → synchronous, handler-owned tx on the accounting pool.
    3. **Customer invoices / payments** → synchronous inside repositories (invoices post at creation, even drafts), bypassing auto-posting rules.
    4. **Manual / recurring / CSV journals** → create + post as separate, non-atomic steps.
    5. **Supplier invoices** post to expense 5110 (not inventory 1130) on approval; `PostSupplierInvoiceEntry`, `PostSupplierPaymentEntry`, `PostDepreciationEntries` exist but have **no callers**.
    Money is `numeric` in PostgreSQL (some `integer` cents columns) but **`float64` in all Go money maths**; no currency column on orders/journals/GL/payments (EGP assumed). Details: [04](04-financial-model-and-policies.md), [09](09-database-schema-and-erd.md), [11](11-financial-posting-matrix.csv), [13](13-accounting-event-architecture.md).

    ## 3. Discovered counts (actual, from this package)
    | Item | Count | Where |
    |---|---|---|
    | Applications | 4 frontends + 1 Go API + PostgreSQL (+ Redis optional) | 03 |
    | RBAC modules | 20 (IDs 1–20; 1–12 frontend `MODULE_IDS`, 13–20 backend migrations) | phase2-progress-notes |
    | Pages / route handlers | 119 | 06 |
    | Form-field rows | 200 (all `NOT TESTED`; EmployeeForm, HR, super-admin, auth, RBAC forms **not covered**) — forms not separately counted | 07 |
    | API routes | 374 (== `routes.manifest`; security agent counted 376 raw registrations incl. variants; 284 with explicit module+level) | 08 |
    | DB tables / FKs / triggers | 102 / 287 / 47 | 09 |
    | RLS | 20 tables (all forced), 23 policies; 70 tables carry `tenant_id`, **57 of them have no RLS** incl. every core financial table | 09, 28 |
    | Business workflows documented | 25 (+ 25 explicitly NOT IMPLEMENTED capabilities, N1–N25) | 26 |
    | Posting-matrix events | 36 | 11 |
    | Competing sources of truth | 12 (S1–S12) | 27 |
    | Financial / advanced report endpoints | 14 (P&L, TB, BS, CF, ledger, receivables & payables aging, VAT, inventory valuation, KPIs, menu profitability, break-even, predictive CF, HR cost) + operational reports (sales, profitability, inventory turnover, overview, order stats, dashboard home, KDS analytics, HR reports) | 12 |
    | Report **metrics** | **NOT COUNTED** (each report's formulas catalogued in 12; no metric tally made) | — |
    | Historical issues reviewed | 36 (21 register issues + 15 historical concerns) | 21 |
    | Knowledge gaps | 40 (5 BLOCKER, 14 HIGH) + addendum | 29 |
    | Diagram files | 27 | diagrams/ |
    | Pages/fields/APIs *tested* | **0** (discovery only) | — |

    ## 4. Most important static findings (spot-verified by the orchestrator unless marked)
    1. **Balance Sheet cannot include current-period earnings** — query limited to `asset/liability/equity` accounts (`financial_report_usecase.go:304`). ✔ verified in source.
    2. **Cash Flow is an accrual P&L plus account-code-prefix guesses** (`'12%'`, `'22%'`; `financial_report_usecase.go:388-392`). ✔
    3. **Tax mismatch:** POS UI uses 14 % (`posStore.ts:107`), server prices orders at hard-coded 10 % (`order_pricing.go:46`); currency hard-coded `"EGP"`; invoice trigger taxes gross before discount (`141_create_invoices_customer.sql:80-82`). ✔
    4. **Fiscal-period lock is effectively off:** 0 `fiscal_periods` rows (✔ runtime), no API to create/close, post path fails open when no row exists, void skips the check (`journal_entry_repo.go:813-842`, `:386`, `:441-520`) — agent-reported, partly verified.
    5. **Money in `float64`**; `round2` truncates negatives wrongly (`financial_report_usecase.go:528`). ✔
    6. **Large parts of the business are not posted to accounting** (matrix 11): refunds/voids (not implemented), order creation & partial payments, inventory purchase/adjustment/waste, expenses, bank transfers, owner capital/withdrawal, VAT/period/year-end close, depreciation; supplier invoices hit expense not inventory. ✔ for supplier→5110 and dead functions; remainder agent-reported.
    7. **POS marks only cash orders paid**, via a second call whose failure is only logged (`pos/page.tsx:302-308`) → unpaid orders never reach the ledger; no idempotency keys on money POSTs. ✔ for the swallowed failure; "card/storefront never reach ledger" is agent-reported.
    8. **Registration seeds no chart of accounts, posting rules or fiscal periods** — runtime-confirmed: tenant 2 (created by the owner at 19:27 UTC today) has 0 COA / 0 rules; tenant 1 has 62 / 6. ✔ CONFIRMED-BY-READONLY-RUNTIME. A first paid order on a new tenant would likely fail to post (INFERRED).
    9. **`products.cost_price` does not exist** but Menu Profitability / Break-even query it (`accounting_advanced_repo.go:976,1045`) — DB error reproduced read-only. ✔ (HTTP response REQUIRES-EXECUTION.)
    10. **Tenant isolation is application-level for all core financial tables** (57/70 tenant tables without RLS). ✔ (dev superuser bypasses RLS anyway — no runtime proof either way.)
    11. **Environment handling:** unset `ENV` ⇒ "development" ⇒ RBAC audit-only + accounting pool fallback; `IsProductionMode()` only for `production|prod` so a `staging` deploy skips production guards; stub-mode mock login with a fixed password exists in source (`auth_repo.go:255-270`); `docker-compose.prod.yml` lacks accounting credentials the code requires in production. ✔ (compose + env function); stub-mode reachability agent-reported.
    12. **Revenue has ~6 competing formulas** (GL net of tax, sales report incl. tax/unpaid, KPIs paid-only incl. tax, order stats incl. cancelled, dashboard tile = paid revenue over the recent-orders list only); `suppliers.current_balance` is decremented on payment but never incremented on invoice approval; inventory valuation reads current stock and ignores its as-of date; "waste" metric reads an unwritten `inventory` table. Agent-reported (27, 26) — REQUIRES-EXECUTION for magnitudes.
    13. **History protection is partial:** CFR-002/003 triggers block UPDATE/DELETE on GL, journal headers/lines (✔ catalog); no TRUNCATE trigger; payments, invoices, salaries, orders have no immutability triggers (agent-reported).
    14. **Journal idempotency is check-then-insert behind an advisory lock, not a DB unique key**; the outbox accepts only `paid_order` (migration 181 CHECK). Agent-reported.

    ## 5. Missing integrations / capabilities (see 26 N1–N25, 25)
    No shift/cashier session, refunds, cancellation of paid orders, bank transfers, fiscal-period creation/closing, inventory-to-supplier-invoice link, unit cost on stock adjustments, recipes→COGS beyond what the posting rules read (UNVERIFIED depth), theme-builder/`/public/themes` backend routes, UI for outbox redrive.

    ## 6. Unresolved accounting policies (need owner/accountant decision)
    VAT rate source of truth (10 % vs 14 % vs configurable) and recoverability; revenue recognition trigger (order created vs paid vs fulfilled); supplier purchases: inventory vs expense capitalisation; inventory costing method; handling of discounts/service/delivery/tips in revenue; payroll accrual vs cash; fiscal-year close & retained-earnings treatment; currency model; refund/credit-note policy; period-lock ownership. Full list: [04](04-financial-model-and-policies.md), [29](29-knowledge-gap-register.md).

    ## 7. Limits of this discovery (read before relying on it)
    - Dev DB is nearly empty (0 orders, 0 journals): integrity probes returning 0 prove only schema health.
    - Dev backend runs as superuser `postgres` (BYPASSRLS), RBAC in audit mode, accounting-pool fallback: no isolation/role claim is supportable from it.
    - Findings from sub-agents were spot-verified by the orchestrator only where marked ✔; the rest are source-reading results awaiting runtime confirmation.
    - Field inventory incomplete (see §3); report-metric tally not made; the pages/fields were not browser-crawled in this phase.

    ## 8. Recommended next phase
    Authorized separately: bring up the isolated restricted-role stack (`docker-compose.phase2a-e2e.yml`), create a fresh audit tenant, run the deterministic scenario ([oracle/](oracle/) is ready: 30/30 self-consistent), then manual UI testing, reconciliation, negative/concurrency/recovery and security under `pos_app`/`pos_accounting`. First runtime checks to prioritise: BS residual with profit, tax 10 % vs 14 %, new-tenant first paid order posting, POS card/cash/storefront ledger reach, duplicate checkout, fiscal-period fail-open, `cost_price` endpoint response.
