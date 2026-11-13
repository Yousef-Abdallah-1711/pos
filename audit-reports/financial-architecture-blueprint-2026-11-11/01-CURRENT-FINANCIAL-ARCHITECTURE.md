# 01 — Current financial architecture

**Evidence basis:** primary workspace HEAD ecadb5ab7665df3eb9f69886a9b186c710eff62d and noted uncommitted wiring, compared with fetched origin/main e519f52eb0badfb80c283f79a52c51acc1c663c6. No live DB was inspected. SOURCE VERIFIED means repository inspection. HISTORICAL AUDIT EVIDENCE means the previous isolated run; it is not a fresh test.

## System map

~~~mermaid
flowchart TB
  subgraph channels[Channels]
    POS[Dashboard POS]
    Dash[Dashboard orders accounting reports]
    Store[Restaurant website and public API]
    Kitchen[KDS]
    SA[Super Admin]
  end
  subgraph services[Go API]
    Tenant[Auth tenant context permissions]
    Order[Order use case and repository]
    Pay[Paymob subscription payments]
    Invoice[Customer and supplier invoices]
    Auto[Auto-posting]
    Journal[Journal handler and repository]
    Ledger[General ledger]
    Reports[Report handlers]
    Jobs[Webhook retry and expiry workers]
  end
  subgraph db[PostgreSQL]
    Ops[(orders items status history)]
    SaaS[(subscription payments)]
    Subledger[(invoices receipts payments inventory movements)]
    CoA[(chart of accounts and posting rules)]
    JE[(journals and lines)]
    GL[(general ledger)]
  end
  POS --> Tenant
  Dash --> Tenant
  Store --> Tenant
  Kitchen --> Tenant
  SA --> Tenant
  Tenant --> Order
  Tenant --> Pay
  Tenant --> Invoice
  Tenant --> Journal
  Tenant --> Reports
  Order --> Ops
  Order -. paid event .-> Auto
  Pay --> SaaS
  Pay --> Jobs
  Invoice --> Subledger
  Invoice --> Journal
  Auto --> JE
  Journal --> JE
  Journal --> GL
  Ledger --> GL
  Reports --> GL
  Reports --> JE
  Reports --> Subledger
~~~

## Current-state matrix

| Domain | Operational source | Financial source today | Ledger integration | Failure handling | Evidence/confidence |
|---|---|---|---|---|---|
| POS/orders | orders, order_items, status history; dashboard POS and public/admin routes | subtotal, tax, discount, total, payment status on order | Paid-order caller separately asks for sales and COGS auto-post | Auto-post errors logged, not returned; order/item writes are separate calls | SOURCE VERIFIED; order runtime NOT TESTED |
| Order tenders | order payment_method/status; no evidence SaaS payments table is order tender ledger | paid flag and order total; cash/card account rules | PostSalesEntry selects collection account; split tender unverified | Payment status and accounting effect are separate | SOURCE VERIFIED |
| SaaS subscription payments | payments rows in integer cents and Paymob IDs | payment row, provider callback, subscription/history | Not restaurant sales accounting | HMAC/amount/currency checks; retry queue; external/local gap remains | SOURCE VERIFIED |
| Accounting receipts/payments | payments_received and payments_made | amount, invoice, bank account and journal link | Accounting repository posts journal | Selected actions use DB transactions; full behavior not runtime verified | migration 143, repository; SOURCE VERIFIED |
| Customer AR | customers_accounting, invoices_customer and lines | invoice total, paid and balance fields | Invoice repository journal path | Multiple writes/transactions vary by operation | migrations 139/141; SOURCE VERIFIED |
| Supplier AP | suppliers, invoices_supplier and lines | invoice and due/paid state | Approval creates/posts journal in a shared repo transaction; payment separate | Auto-post probe blockers historically observed | migrations 138/142/143; SOURCE VERIFIED + HISTORICAL |
| Inventory | product stock/cost and inventory_movements | product fields and movement history; costing unclear | COGS auto-post from order; other paths separate | No atomic order-stock-GL proof; uncommitted workspace adds adjustment handler | migrations 015/144; SOURCE VERIFIED |
| Chart of Accounts | chart_of_accounts, seed and auto_posting_rules | mutable current_balance and account metadata | Rules/default account candidates select accounts | Scope is app-enforced in inspected paths; effective-dated mappings not shown | migrations 136/152/156; SOURCE VERIFIED |
| Journals | journal_entries and journal_entry_lines | draft/posted/void, header totals, refs, attachments | Draft create and auto-post are separate operations | UI nil attachments and auto nil attachments failed historically | migration 137; SOURCE VERIFIED + HISTORICAL |
| GL | general_ledger linked to journal lines | posted rows and stored running_balance | Appended during posting; GL rows + posted status share transaction | Draft creation is separate; backdate/concurrency not tested | migration 154; SOURCE VERIFIED |
| Reports | report API / Dashboard pages | P&L and cash flow query journals; Balance Sheet GL; aging invoices; valuation products | Different SQL paths and source snapshots | No shared watermark shown; Balance Sheet omits earnings | financial_report_usecase.go; HISTORICAL + SOURCE |
| Tax | fixed 10% order calculation; tax fields/accounts/periods | order/invoice/journal tax amount | Tax accounts exist; recognition policy unclear | Hard-coded rate; jurisdiction not approved | order_usecase.go, migrations 136/137/146; SOURCE |
| Export | client PDF; API CSV/PDF; export audit tables | generated output and log | Reads report data | XLSX not established by route name | handlers + ReportExportButton; SOURCE |
| Identity/RBAC | auth/session/JWT, tenant middleware, permission wrappers | tenant/restaurant context and role | Queries scoped in many repositories | main.go has uncommitted route/session wiring; local branch removes RLS migration files; live DB state unknown | SOURCE; tenant matrix NOT TESTED |
| Background jobs | webhook retry queue and subscription expiry worker | attempts and status | Webhook retry only; no accounting outbox | primary uncommitted main.go changes worker lifecycle; posting recovery not durable | SOURCE |

## Competing authorities today

Order totals, payment status, SaaS payment rows, accounting receipts/payments, invoices, product stock, inventory movements, journals, GL, current_balance, running_balance and report results are not one fact. Some duplicate or overlap. P&L and cashflow query journal lines; Balance Sheet uses GL; aging uses invoices; valuation uses product data. That prevents assuming reports share a consistent posting snapshot.

The GL writer loads the latest account running balance ordered by date/id and appends. It does not reflow later entries after backdating and does not serialize the read-modify-write by account in the inspected path. Treat running balances as cache, not authority. SOURCE VERIFIED risk, not runtime reproduced.

## Baseline and local workspace differences

Fetched origin/main remains e519f52. Primary branch ecadb5 differs in order/payment/invoice/report/permission code and deletes migration files relative to e519, including prior RLS/payment-state migrations 167–178. This proves a source-tree difference, not that a live database rolled back. The primary worktree is heavily dirty. Its main.go has uncommitted auth/session, POS route, inventory handler and worker-shutdown wiring. No worktree was reset or altered. Core financial report/auto-posting/journal/GL files were unchanged at primary HEAD during inspection; order/payment/invoice paths must be re-reviewed before applying old findings. Primary RefundPayment holds a row lock during Paymob call then commits local refunded state: provider success followed by local commit failure remains possible. SOURCE VERIFIED, runtime untested.

## Entry points and surfaces

Accounting API routes in backend/cmd/api/main.go use module permissions. Financial report APIs use report-specific read/export wrappers. Public order route has tenant context; current uncommitted primary wiring adds authenticated /api/v1/pos/orders. A wrapper alone does not prove object scope. Dashboard surfaces POS, orders, accounting chart/journals/ledger/suppliers/invoices/bank accounts, and financial report pages. Storefront/theme builder and KDS modify/display operational order state, not accounting authority. Convergence across channels is NOT VERIFIED.
