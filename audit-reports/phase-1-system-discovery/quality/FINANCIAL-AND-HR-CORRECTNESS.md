# Financial and HR Correctness (Phase 1)

---

# Accounting and HR correctness notes (Phase 1 security/finance auditor)

Read-only. Evidence `file:line` relative to repo root. Items mirrored in `evidence/issues-security.json` (FIN-001..FIN-008).

## What is done well
- Double-entry schema with CHECK constraints (non-negative, one-sided lines) and header-account protection: `backend/migrations/137_create_journal_entries.sql:49-62`, `journal_entry_repo.go:642-660`.
- Application balance check on create and again on post (tolerance 0.009): `journal_entry_repo.go:196,257,389-395,706-735`.
- Posting states `draft|posted|void` with CHECK (`137:38`); only drafts may be edited (`journal_entry_repo.go:281-283`); posted entries cannot be changed through the API and void appends reversing GL rows rather than deleting (`journal_entry_repo.go:459-480`).
- Fiscal period closed check on post (`journal_entry_repo.go:754-775`; no row means open).
- Invoice line and header totals computed by DB triggers, so client-supplied totals are overridden: `141_create_invoices_customer.sql:78-131`; invoice creation posts a linked AR/revenue/VAT journal in the same transaction (`invoice_customer_repo.go:336-413,438-500`). Edits blocked for void or part-paid invoices (`invoice_customer_repo.go:146-150`, `existing.AmountPaid > 0`).
- Payroll gross/net/total_deductions computed by DB trigger (`024:129-148`), with CHECK constraints for non-negative amounts and valid pay period.
- Order to accounting automation exists: `OrderUseCase.AutoPostPaidOrder` posts sales and COGS when payment becomes paid (admin handler `admin_order_handler.go:496-499`), with rule-based accounts (`auto_posting_rules`, migration 152) and idempotency by reference lookup (`auto_posting_engine.go:61-72`). Payroll posts on mark-paid (`hr_salary_handler.go:423-436`). HR to accounting link therefore exists via `PostPayrollEntry`.
- Tax periods, VAT report, fiscal periods, budgets, bank reconciliation, recurring journals, e-invoice submission endpoints exist (main.go:1028-1051,1353-1359).

## Where invoice totals and tax are calculated
- Orders: Go, `usecase/order_usecase.go:127-186` (float64, 10 percent hardcoded).
- Customer/supplier invoices: SQL triggers (authoritative); frontend sends line inputs only (the backend ignores client `line_total`, `tax_amount`).
- Payroll: SQL trigger for gross/net; deductions are client inputs.
- Frontend calculation of financial totals for display was not audited.

## Findings

### FIN-001 [P1] Auto-posting silently plugs any journal imbalance into the last line

- Category: accounting
- Status: OPEN
- Files: `backend/internal/usecase/auto_posting_engine.go`
- Evidence:
  - auto_posting_engine.go:633-663 normalizeBalancedLines adds diff to last credit/debit when |diff|>0.009 regardless of size
  - used by createAndPostEntry (522-566) for sales, COGS, payroll, supplier, depreciation
- Root cause: Defensive rounding fix grew into unconditional balancing.
- Recommended fix: Only absorb rounding diffs up to a few cents with explicit rounding account; otherwise return error and surface failed posting.
- Required regression tests:
  - Order with tax inconsistent by 5.00 fails posting and alerts
  - Rounding diff <=0.02 posts to rounding account

### FIN-002 [P2] No database-level debit=credit enforcement; balance trigger only recomputes totals; float tolerance

- Category: accounting
- Status: OPEN
- Files: `backend/migrations/137_create_journal_entries.sql`, `backend/internal/repository/journal_entry_repo.go`
- Evidence:
  - 137:79-108 validate_journal_balance() does not RAISE; no CHECK total_debit=total_credit; no immutability trigger for posted lines
  - journal_entry_repo.go:393,730 float64 tolerance 0.009 per entry
- Root cause: Integrity rules live only in Go.
- Recommended fix: Add deferred constraint trigger (sum debit = sum credit on commit for posted/post transition), block UPDATE/DELETE of posted lines and entries, use NUMERIC/decimal in Go and exact compare.
- Required regression tests:
  - Direct SQL insert of unbalanced posted entry fails
  - UPDATE of posted line fails
  - Imbalance of 0.005 rejected

### FIN-003 [P1] No row locking: concurrent post/void, running-balance updates and number generation race

- Category: accounting
- Status: OPEN
- Files: `backend/internal/repository/journal_entry_repo.go`, `backend/internal/repository/payment_acct_repo.go`, `backend/internal/repository/invoice_customer_repo.go`
- Evidence:
  - journal_entry_repo.go:363-375 status read without FOR UPDATE; 544-596 running_balance read-modify-write per account; 784-798 MAX()+1 numbering; payment_acct_repo.go:563-575 same; grep shows no FOR UPDATE in repositories
- Root cause: Optimistic code paths assume single writer.
- Recommended fix: SELECT ... FOR UPDATE on entry and account rows (ordered), or compute balances from SUM; use sequences/counter rows per tenant-year; unique constraint retry.
- Required regression tests:
  - Parallel double-post yields one GL set
  - Parallel posts keep running_balance equal to SUM
  - Parallel create yields distinct numbers

### FIN-004 [P2] float64 used for money in Go while DB uses NUMERIC/DECIMAL

- Category: accounting
- Status: OPEN
- Files: `backend/internal/domain/accounting.go`, `backend/internal/domain/order.go`, `backend/internal/usecase/order_usecase.go`, `backend/internal/handler/http/hr_salary_handler.go`
- Evidence:
  - ~100 float64 fields in accounting/hr domains; order_usecase.go:142-182 float arithmetic; hr_salary_handler.go:122-141
  - migrations use NUMERIC(15,2) / DECIMAL(10,2)
- Root cause: Convenience typing.
- Recommended fix: Adopt shopspring/decimal or integer minor units at boundaries; round half-even at defined points.
- Required regression tests:
  - Property test: sum of 0.1 x 3 lines equals exact total
  - Totals match SQL recomputation

### FIN-005 [P1] Order accounting not reversed on refund/cancel; posting failures only logged; create and post not atomic

- Category: accounting
- Status: OPEN
- Files: `backend/internal/usecase/order_usecase.go`, `backend/internal/handler/http/admin_order_handler.go`, `backend/internal/usecase/auto_posting_engine.go`
- Evidence:
  - order_usecase.go:459-474 errors only log.Printf
  - admin_order_handler.go:482-499 posts when paid; no reverse on paid->refunded/cancelled (no reversal code in order_usecase)
  - auto_posting_engine.go:522-566 CreateJournalEntry commits, then AppendEntry separately
  - entry_date = order.CreatedAt (auto_posting_engine.go:149) not payment date; closed period would fail silently
- Root cause: Accounting integration added as best effort side effect.
- Recommended fix: Single transaction for order status + posting (or outbox with retry/alert), reversing entry on refund/cancel, post at payment date, reconcile report of paid orders without journal.
- Required regression tests:
  - Paid then refunded order nets to zero in GL
  - Posting failure surfaces and retries
  - Orders-paid total equals sales GL total for period

### FIN-006 [P2] Tax/discount calculation defects: hardcoded 10% order tax; invoice tax ignores discounts

- Category: accounting
- Status: OPEN
- Files: `backend/internal/usecase/order_usecase.go`, `backend/migrations/141_create_invoices_customer.sql`
- Evidence:
  - order_usecase.go:174-175 TaxAmount = subtotal*0.10, DiscountAmount never set
  - 141:80-82 tax = qty*price*rate before discount_amount; discount_rate unused in line_total
- Root cause: Placeholder logic retained.
- Recommended fix: Per-tenant tax config (rate, inclusive/exclusive); tax on net-of-discount base; support discount_rate; document rounding.
- Required regression tests:
  - Invoice with discount computes tax on discounted base
  - Tenant with 14% VAT gets 14% on orders

### FIN-007 [P2] Payroll records editable/deletable after payment and posting; tax/insurance are manual inputs

- Category: hr
- Status: OPEN
- Files: `backend/internal/handler/http/hr_salary_handler.go`, `backend/internal/repository/salary_repo.go`, `backend/migrations/024_create_salaries_table.sql`
- Evidence:
  - hr_salary_handler.go:250-354 UpdateSalary no status guard (can change status/amounts after paid)
  - 356-380 DeleteSalary hard DELETE, journal entry orphaned
  - salary_repo.go:216-227 MarkAsPaid no precondition
  - req fields tax, social_insurance etc. are client-supplied (122-141); gross/net computed by DB trigger (024:129-148)
- Root cause: CRUD handler without lifecycle rules.
- Recommended fix: Lock paid records; reversal journal on correction; soft-delete only for unpaid; compute statutory deductions from configurable rules; require approval workflow.
- Required regression tests:
  - PUT on paid salary returns 409
  - Delete paid salary blocked
  - Mark paid twice posts once
  - Net = gross - deductions matches SQL

### FIN-008 [P3] Void and account deletion edge cases

- Category: accounting
- Status: OPEN
- Files: `backend/internal/repository/journal_entry_repo.go`, `backend/internal/repository/account_repo.go`
- Evidence:
  - journal_entry_repo.go:437-504 void of posted entry reverses in GL dated today with no closed-period check
  - account_repo.go:212 hard DELETE of chart_of_accounts (FK behaviour unverified)
  - order number time%1e6 (order_usecase.go:63) collision-prone
- Root cause: Edge cases not covered.
- Recommended fix: Create explicit reversing journal entry linked to original, respect fiscal period, restrict account delete to unused accounts.
- Required regression tests:
  - Void in closed period blocked or posts to open period with link
  - Delete account with GL rows rejected

## Data-type summary
DB: `NUMERIC(15,2)` for journals and invoices, `DECIMAL(10,2)` for salaries (10 digits: max 99,999,999.99; adequate for most but check high-denomination currencies), `amount_cents` integers for subscription payments. Go: `float64` throughout accounting/HR/orders (FIN-004). Journal tolerance 0.009 per entry (FIN-002).

## Unverified
Supplier invoice and payment-made posting logic (`invoice_supplier_repo.go`, `payment_acct_repo.go`) was only skimmed; inventory movement costing (`PostCOGSEntry`, moving average vs last cost); bank reconciliation completeness; depreciation entries; tax period close effects on posting; frontend total calculations; RLS enforcement in DB.
