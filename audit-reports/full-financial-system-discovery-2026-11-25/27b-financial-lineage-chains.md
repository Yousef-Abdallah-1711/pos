# 27b — Financial lineage chains (implemented paths)

Chain template: UI field → component → API → middleware → handler → usecase → repository → SQL tx → operational table → accounting event → journal → GL → report → dashboard → export. `✔` = confirmed by source (this or earlier pass), `?` = not traced, `✘` = chain broken (exact break stated). Based on [11](11-financial-posting-matrix.csv), [10](10-financial-transaction-lifecycles.md), [31](31-independent-finding-verification.csv). Page/route names per [06](06-complete-page-inventory.md) and [08](08-api-inventory.csv). No runtime verification.

Common middleware ✔: tenant session auth → `wrapWithPermission(module, level)` → tenant status guard.

| # | Operation | Chain (✔ confirmed, ✘ break) | Break / unverified link |
|---|---|---|---|
| 1 | Cash sale (POS) | POS page `pos/page.tsx` → `posApi.createOrder` → orders handler → `order_usecase.CreateOrder` (pricing 10% tax, `order_pricing.go:46`) → order_repo → `orders`/`order_items` (status pending, payment pending ✔) → second call `PUT /admin/orders/{id}/payment-status` (cash only, `pos/page.tsx:302-308` ✔, failure only logged) → `UpdatePaymentStatusAs` tx → `accounting_event_outbox` row ✔ (same tx) → worker (15 s) → `AutoPostPaidOrder` → PostSalesEntry + PostCOGSEntry ✔ → journals → GL ✔ → P&L/TB ✔ | Link between the two UI calls is not atomic ✘ (order may stay pending). Whether worker posts when COA missing: REQUIRES-EXECUTION. Dashboard/export: see 30. |
| 2 | Card sale (POS) | same as 1 until create; **no step marks paid** (`paymentMethod === 'cash'` only) | ✘ no payment capture, no paid status, no outbox, no revenue (VER-008). Staff may mark paid manually via admin orders UI. |
| 3 | Wallet sale | same as card | ✘ same; no wallet/gateway code for orders (VER-010) |
| 4 | Online (storefront) order | storefront checkout → public order API → `CreateOrder` (pending) | ✘ no payment step; paid only by manual admin marking (KG-SEC-032) |
| 5 | Cash on delivery | storefront order → driver/delivery | ✘ driver completion/COD settlement not implemented (N20) |
| 6 | Customer invoice | accounting customer-invoices page → API → invoice repo → `invoices_customer` (trigger computes tax on gross, mig 141:80-82 ✔) → journal posted at creation incl. drafts (earlier read) → GL 1120 / revenue / 2120 | draft-posting behaviour agent-reported; ? |
| 7 | Customer receipt | invoice payment form → `payments_received` + invoice balance_due → journal (cash/AR) | posting path not re-traced ? |
| 8 | Supplier invoice | supplier invoices page → create (draft, no journal ✔) → approve → `buildJournalEntry` ✔ (Dr 5110 expense; Dr 1150 VAT; Cr payable) | ✘ no inventory/stock link (VER-006, N7) |
| 9 | Supplier payment | payments_made record | ✘ `PostSupplierPaymentEntry` has no caller (VER-007); whether payment path posts through another function ? |
| 10 | Inventory receiving | no purchase order/receiving (N6); `POST /inventory/adjust` only | ✘ no GL |
| 11 | Inventory consumption | order confirm decrements stock (earlier read); COGS journal at payment (cost at that time) | timing mismatch (33b) |
| 12 | Inventory adjustment / waste | `/inventory/adjust` → `inventory_movements` | ✘ no GL call (N13); turnover waste reads unwritten table |
| 13 | Payroll | HR salary form → salaries (status transitions) → `MarkAsPaid` handler tx (`hr_salary_handler.go:470-558` ✔) → payroll journal on accounting pool ✔ | accrual vs cash policy open (PAY-01); approval not enforced (N12) |
| 14 | Expenses | no expense module (N8) | only manual journal or supplier invoice |
| 15 | Bank transactions | bank-accounts page → bank_accounts/transactions; reconciliation sessions | ? whether transactions post to GL; transfers absent (N9) |
| 16 | Refund | order `refunded` status value only | ✘ transition from paid blocked; no journal (N2) |
| 17 | Cancellation | item-level cancel (`cancelled_quantity`, migration 181_partial) ; whole-order cancel guarded | reversal accounting ? |
| 18 | Manual journal | journal-entries page → create (draft) → post (separate call) → `CreateJournalEntryTx`/`PostEntryTx` ✔ → GL (cfr002 triggers ✔) | non-atomic create+post (earlier read) |
| 19 | Opening balances | account `opening`/current_balance fields + diagnostics; no workflow | ✘ (N5) |
| 20 | Financial period closing | fiscal_periods table, checker only; no create/close API (VER-013/016) | ✘ entire chain missing |

Downstream (all journals): GL → Trial Balance ✔, P&L ✔, Balance Sheet (✘ excludes revenue/expense, VER-001), Cash Flow (✘ prefix heuristic, VER-002), Ledger; dashboard KPIs read orders not GL (33b); exports: financial statements have none.
