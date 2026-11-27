# 32 - Accounting policy decision register

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Method: static reading of `PHASE4_ACCOUNTING_POLICY_DECISION_PACKAGE.md` (repo root, "P4"), `04-financial-model-and-policies.md` ("04"), `11-financial-posting-matrix.csv` ("11"), `26-business-workflow-catalog.md` ("26"). No application code was run.

## READ THIS FIRST

- **Every "Recommended option" below is a proposal for discussion. NONE is an approved policy.** No approval record exists in the repository (04, P4 "Owner decision record"). Nothing here authorises implementation.
- Existing decision IDs (`OPEN-01`, `EARN-01`, `CLOSE-01`, `TAX-01`, `ROUND-01`, `REFUND-01`, `INV-01`, `PAY-01`, `CF-01`, `CFR-011`) are kept as-is. Where the prior package already proposes something, it is quoted by reference, not replaced. IDs prefixed `D-` below are NEW and exist only in this audit package.
- "Still open" = no approval record found. All items are **OPEN**.
- Egypt context (EGP currency, 14% VAT as used by the control oracle) is recorded as a **question for a qualified Egyptian accountant/tax adviser**, never asserted as law. Statutory deduction rates, VAT recoverability, e-invoicing and fiscal-year rules are all unverified here.
- The oracle in `oracle/scenarios.py` uses recommendations as PROVISIONAL assumptions only so numbers can be computed; changing a decision means regenerating the expected values.

## Index and crosswalk

| # | Decision | Prior ID | New ID | Status | Approver |
|---|---|---|---|---|---|
| 1 | Opening balances | OPEN-01 | - | OPEN | Accountant |
| 2 | Current-period earnings presentation | EARN-01 | - | OPEN | Accountant |
| 3 | Period closing | CLOSE-01 | - | OPEN | Accountant + owner |
| 4 | Revenue recognition | (REFUND-01/CFR-005 touch) | D-REV-01 | OPEN | Accountant + owner |
| 5 | Cash vs accrual basis | (04 s4 "no basis approved") | D-BASIS-01 | OPEN | Accountant + owner |
| 6 | VAT configuration | TAX-01 | - | OPEN | Tax adviser + owner |
| 7 | Input VAT recoverability | TAX-01 (part) | D-TAX-02 | OPEN | Tax adviser |
| 8 | Inventory purchase capitalisation | INV-01 (part) | D-INV-02 | OPEN | Accountant + inventory owner |
| 9 | Inventory costing | INV-01 | - | OPEN | Accountant + inventory owner |
| 10 | COGS recognition | INV-01 (part) | D-COGS-01 | OPEN | Accountant |
| 11 | Discounts | TAX-01 ordering (part) | D-DISC-01 | OPEN | Accountant + owner |
| 12 | Service charges | - | D-SVC-01 | OPEN | Tax adviser + owner |
| 13 | Delivery charges | - | D-DLV-01 | OPEN | Tax adviser + owner |
| 14 | Tips | - | D-TIP-01 | OPEN | Accountant + payroll adviser |
| 15 | Refunds | REFUND-01 | - | OPEN | Accountant + payments owner + owner |
| 16 | Credit notes | (REFUND-01 / CFR-011 adjacent) | D-CN-01 | OPEN | Accountant |
| 17 | Payroll accrual | PAY-01 | - | OPEN | Payroll lead + local adviser |
| 18 | Payroll settlement | PAY-01 | - | OPEN | Payroll lead + local adviser |
| 19 | Bank transfers | CF-01 (part) | D-BANK-01 | OPEN | Accountant |
| 20 | Capital contributions | CF-01 (part) | D-EQ-01 | OPEN | Owner + accountant |
| 21 | Owner withdrawals | CF-01 (part) | D-EQ-02 | OPEN | Owner + accountant |
| 22 | Retained earnings | CLOSE-01/EARN-01 | D-RE-01 | OPEN | Accountant |
| 23 | Historical corrections | CLOSE-01/REFUND-01 | D-HIST-01 | OPEN | Accountant + owner |
| 24 | Currency precision | ROUND-01 | - | OPEN | Accountant + owner |
| 25 | Rounding | ROUND-01 | - | OPEN | Accountant + owner |
| - | Cash Flow classification (cross-cut) | CF-01 | - | OPEN | Accountant |
| - | Legacy supplier auto-post methods | CFR-011 | - | OPEN | Owner + accountant |

Source IDs `N1..N25` refer to the NOT IMPLEMENTED capability list in 26 (lines 262-286).

## Cross-cutting conflicts found while compiling this register

1. **Tax rate**: order pricing hard-codes 10% (`backend/internal/usecase/order_pricing.go:11,46`), the POS UI shows 14% (`POSOrderSummary`, 06 POS section), and the control oracle assumes 14%. Same restaurant could display, charge and book different VAT.
2. **Revenue timing**: orders post at `paid` but dated `order.CreatedAt` (`auto_posting_engine.go:137`); customer invoices post at create (even draft) dated invoice date (11); supplier invoices post at approval (11); payroll posts at mark-paid (`auto_posting_engine.go:276-282`). Four triggers, no stated basis.
3. **Purchases**: supplier invoice debits Expense/5110 Food COGS by default, not Inventory 1130 (11, "Supplier invoice approved"), while COGS on sale credits 1130 (`auto_posting_engine.go:143-260`). 1130 can therefore only fall.
4. **VAT report vs GL**: VAT summary reads invoice tables, not GL 2120/1150, so order VAT posted by `auto_sales` is outside it (11, "VAT settlement").
5. **Discount**: invoice trigger computes tax on gross before discount (`backend/migrations/141_create_invoices_customer.sql:80-82`); order tax is applied before `DiscountAmount` is subtracted (`order_pricing.go:53`); sales journal posts 4150 Sales Discounts (`auto_posting_engine.go:73-120`).
6. **Rounding**: integer half-up cents (`order_pricing.go:46`), float `math.Round(v*100)/100` (`auto_posting_engine.go:645-646`), `round2` truncating cast wrong for negatives (`financial_report_usecase.go:528-530`), client tolerance 0.004 (`order_usecase.go:163`).

---

## Decision details

Format per decision: **Current** (cited) / **Conflicts** / **Options** / **Recommended (NOT APPROVED)** / **Consequences** / **Affected** / **Future acceptance tests** / **Approval**.

### 1. Opening balances - OPEN-01 (still OPEN)
- **Current**: Not synthesised; only read-only diagnostic `backend/internal/repository/opening_balance_diagnostics.go` (04 s7); no guided workflow (N5); new tenants get no chart of accounts (N25).
- **Conflicts**: P4 says missing balances leave a visible residual; 5,000 piaster residual kept `UNRECONCILED` in synthetic fixture (P4 safeguards).
- **Options**: (a) one evidence-backed opening journal offset to an Opening Balance Equity account; (b) import per-account source balances with attachments; (c) defer balances lacking evidence and keep a warning.
- **Recommended**: P4's proposal (evidence plus explicit offset account, preserve omitted-balance warnings) plus (a): a single dated journal, separate cash-flow class "opening".
- **Consequences**: sets day-1 assets/liabilities/equity, inventory base for costing, AR/AP aging start; wrong offset distorts equity forever.
- **Affected**: `journal_entries`, `general_ledger`, `chart_of_accounts` (3xxx), `customers_accounting`/`suppliers` balances, Balance Sheet, Cash Flow (opening cash), inventory valuation.
- **Tests**: SCN-01; omitted-balance warning; no plug; duplicate import rejected; BS residual 0 after opening; opening cash excluded from CF categories.
- **Approval**: Qualified accountant/controller.

### 2. Current-period earnings presentation - EARN-01 (OPEN)
- **Current**: Balance Sheet selects asset/liability/equity only (`financial_report_usecase.go:304`); `EquationBalanced` via `round2` (`:343-352`); diagnostic only (`current_earnings_diagnostic.go`). Any P&L activity makes BS UNRECONCILED (04 finding 1).
- **Conflicts**: oracle shows equity must include current earnings (`oracle/deterministic_control.py:100`); app would show residual equal to net profit/loss (e.g. -2,000 EGP on the control data).
- **Options**: separate "Current earnings" line (no journal); include in equity total (no journal); close entries (see 3).
- **Recommended**: show as a separate explicit line included in equity total, labelled as unclosed, never as a journal. (Consistent with P4.)
- **Consequences**: BS reconciliation status, export gating (P4 safeguards), board/lender reporting.
- **Affected**: `GenerateBalanceSheet`, `BalanceSheet` component, PDF export.
- **Tests**: SCN-23 pre-close equity 98,000 with booked equity 100,000 + current earnings -2,000; multi-period; no double count after close.
- **Approval**: Qualified accountant.

### 3. Period closing - CLOSE-01 (OPEN)
- **Current**: `fiscal_periods` exists (`migrations/153_create_fiscal_periods.sql`); post path rejects `closed` periods, no period row means allowed (`journal_entry_repo.go:386,813-831`); void has no period check and dates reversal `time.Now()` (`:441-520`); no close routine (N10); CBR 0 rows.
- **Conflicts**: void into a closed period is allowed while posting is not.
- **Options**: soft close; hard close with adjustment period; controlled reopen with dual authorisation (P4).
- **Recommended**: P4: define states/roles/audit events first; PROVISIONAL oracle model = hard close via Income Summary to Retained Earnings with adjustment periods.
- **Consequences**: locks history, determines whether corrections are back-dated or current-period; VAT-return alignment.
- **Affected**: `fiscal_periods`, journal post/void, reports, retained earnings, tax periods (`146_create_tax_periods.sql`).
- **Tests**: SCN-23, SCN-24; duplicate close; post into closed period rejected incl. void; reopen audited; roll-forward.
- **Approval**: Accountant + owner.

### 4. Revenue recognition - D-REV-01 (OPEN; related REFUND-01, CFR-005)
- **Current**: order revenue at `payment_status=paid` via async outbox (11; `auto_posting_engine.go:50-141`), dated order creation; no posting at creation or delivery; customer invoice at create; unpaid delivered order never in GL but present in operational Sales report (11, first rows).
- **Conflicts**: GL vs operational reports; back-dating into creation period when paid later; invoice linked to an order can double count (11).
- **Options**: at capture/payment; at fulfilment/delivery; at invoice date (accrual); mixed per channel.
- **Recommended**: recognise at fulfilment (sale completed) for dine-in/POS (equals payment), at delivery for delivery orders; invoice date for B2B; entry date = event date. Deferred-revenue for prepaid not delivered.
- **Consequences**: period revenue, VAT point-of-supply, cutoff, unearned revenue liability, refunds.
- **Affected**: `orders`, `auto_posting_engine`, `accounting_events` outbox, Sales/P&L reports, VAT report.
- **Tests**: SCN-06..11, 24; paid-later order cutoff; prepaid-undelivered at month end; invoice+order no double count.
- **Approval**: Accountant + owner. Question for accountant: Egyptian tax-point rules for restaurant supplies and online orders.

### 5. Cash vs accrual basis - D-BASIS-01 (OPEN)
- **Current**: mixed (04 s4): orders cash-on-paid, payroll cash-on-paid (no accrual, 2130 unused: 11), supplier/customer invoices accrual.
- **Options**: full accrual; cash basis; hybrid with explicit per-document rules.
- **Recommended**: full accrual for reporting (oracle PROMPT assumption), with cash-basis view only through Cash Flow.
- **Consequences**: month-end payroll and unpaid rent/utilities appear in P&L; AP/AR/accruals required.
- **Affected**: all posting engines, P&L, BS, expense module (N8).
- **Tests**: SCN-18, 19, 24; month-end cutoff.
- **Approval**: Accountant + owner.

### 6. VAT configuration - TAX-01 (OPEN)
- **Current**: hard-coded 10% on every order (`order_pricing.go:11,46`); invoice per-line user rate (`InvoiceForm.tsx:39,421`); POS display 14%; sales credit to `2120` (`auto_posting_engine.go:71,115-120`); `tax_periods` unused (04 s3); no configurable/inclusive pricing (N19).
- **Conflicts**: 10% vs 14% vs per-line rates; VAT report excludes order VAT.
- **Options**: per-tenant rate table with effective dates; inclusive vs exclusive; line vs invoice rounding; exemptions.
- **Recommended**: P4: persist basis and rate on each source transaction; compute under approved rules. Oracle PROVISIONAL: 14% exclusive, per-document half-up.
- **Consequences**: invoice totals, VAT payable, revenue net, receipts printed, customer-visible prices.
- **Affected**: `orders`, `order_items`, `invoice_*`, GL 2120/1150, VAT report, storefront checkout, POS.
- **Tests**: SCN-03..12, 22; rate-change effective date; inclusive price extraction; per-line vs per-document rounding.
- **Approval**: Tax adviser + owner. Question: correct rate(s) for restaurant dine-in/delivery/takeaway and any additional levies; e-invoice obligations.

### 7. Input VAT recoverability - D-TAX-02 (OPEN)
- **Current**: supplier invoice debits 1150 Input VAT Receivable (11); order COGS ignores VAT; control oracle assumes recoverable (PROMPT assumption, `oracle/deterministic_control.py:12`). No VAT settlement journal (11).
- **Options**: fully recoverable; partly recoverable (apportionment); non-recoverable (capitalise into cost).
- **Recommended**: recoverable only when supported by a valid supplier tax invoice and business use; otherwise capitalise into inventory/expense. Netting Input vs Output at period end via settlement journal.
- **Consequences**: inventory cost, expenses, VAT payable, cash at filing.
- **Affected**: supplier invoices, expense module, VAT report, 1150/2120, inventory cost.
- **Tests**: SCN-03, 04, 05, 17 (waste), 19, 24; VAT settlement journal; non-recoverable path.
- **Approval**: Tax adviser. Question: recoverability conditions in Egypt and treatment on waste/shrinkage.

### 8. Inventory purchase capitalisation - D-INV-02 (OPEN; part of INV-01)
- **Current**: supplier invoice defaults to expense 5110, no `product_id` link, no stock receipt movement (11; N6, N7); stock receipt has no journal and no unit cost (11).
- **Conflicts**: COGS on sale credits 1130 which purchases never debit.
- **Options**: capitalise to Inventory 1130 with receipt; expense on purchase (periodic); hybrid by item class (perishables expensed).
- **Recommended**: capitalise tracked stock items to Inventory on receipt; expense consumables.
- **Consequences**: inventory asset size, COGS timing, gross margin.
- **Affected**: supplier invoices, `inventory_movements`, `products.cost`, valuation report.
- **Tests**: SCN-03, 04, 05, 15, 24; supplier-invoice-to-stock link; return.
- **Approval**: Accountant + inventory owner.

### 9. Inventory costing - INV-01 (OPEN)
- **Current**: last cost = latest positive purchase/opening/adjustment `unit_cost` else `products.cost` (`inventory_movement_repo.go:158-206`); valuation report is `qty * products.cost` at query time and ignores `asOf` (`financial_report_usecase.go:488-526`).
- **Options**: FIFO; weighted average; specific identification; standard cost with variance.
- **Recommended**: weighted average perpetual (simple, deterministic), snapshot cost per movement (P4: persist valuation basis and cost snapshot). PROVISIONAL in SCN-15.
- **Consequences**: COGS, margins, inventory asset, as-of valuation.
- **Affected**: `inventory_movements`, `GetProductCost`, COGS, valuation report, profitability reports.
- **Tests**: SCN-15..17; negative stock; as-of valuation; cost change after sale.
- **Approval**: Accountant + inventory owner.

### 10. COGS recognition - D-COGS-01 (OPEN)
- **Current**: COGS at payment, dated order creation, only tracked items with cost>0 (`auto_posting_engine.go:143-260,233`); stock reservation at confirmation (trigger) so quantity moves before accounting (11); no recipe/BOM depletion (N15); all COGS to 5110.
- **Options**: at sale; at fulfilment; recipe-based theoretical plus variance.
- **Recommended**: same event as revenue; recipe consumption deferred until BOM exists, with variance account.
- **Consequences**: margin by period; untracked/zero-cost items produce revenue without cost.
- **Affected**: `auto_posting_engine`, `inventory_movements`, profitability report.
- **Tests**: SCN-13, 15, 24; zero-cost item warning; COGS reversal on refund.
- **Approval**: Accountant.

### 11. Discounts - D-DISC-01 (OPEN; TAX-01 ordering)
- **Current**: no assignment to `DiscountAmount` (N17); engine would debit 4150 (`auto_posting_engine.go:73-120`); per-item discounts not posted (11); invoice tax before discount (migration 141:80-82).
- **Options**: tax after discount; tax before discount; contra-revenue gross presentation vs net revenue.
- **Recommended**: VAT on post-discount price (SCN-12 PROVISIONAL); gross revenue plus Sales Discounts contra.
- **Consequences**: VAT payable, revenue, customer total.
- **Affected**: pricing, orders, invoices, 4150, reports.
- **Tests**: SCN-12; percentage vs fixed; refund of discounted item; invoice vs order parity.
- **Approval**: Accountant + owner. Question: does Egyptian VAT base follow the net-of-discount price.

### 12. Service charges - D-SVC-01 (OPEN)
- **Current**: no field on `domain.Order` (04 s4; N18).
- **Options**: revenue (VAT-able); pass-through liability distributed to staff; excluded.
- **Recommended**: revenue and VAT-able unless distributed to staff as pass-through (needs accountant ruling). SCN-22 PROVISIONAL revenue.
- **Consequences**: VAT, labour cost, customer total.
- **Affected**: pricing, orders, receipts, payroll.
- **Tests**: SCN-22; refund of service charge; rounding.
- **Approval**: Tax adviser + owner. Question: statutory service-charge rules in Egypt.

### 13. Delivery charges - D-DLV-01 (OPEN)
- **Current**: fee forced to 0 (`order_pricing.go:52`); journal credits 4120 Delivery Sales with fallback 4100 but total ignores delivery so a non-zero fee would fail balance validation (04 s4; N20).
- **Options**: revenue (VAT-able); third-party courier pass-through; free.
- **Recommended**: separate VAT-able revenue (SCN-10 PROVISIONAL); courier cost as expense.
- **Consequences**: totals, VAT, COD reconciliation.
- **Affected**: pricing, `orders.delivery_fee`, driver module, 4120.
- **Tests**: SCN-10; fee with discount; refund of fee.
- **Approval**: Tax adviser + owner.

### 14. Tips - D-TIP-01 (OPEN)
- **Current**: tips only on payroll (`salary.Tips`, `auto_posting_engine.go:309,335-343`, 2140 Tips Payable); none on orders (N18).
- **Options**: liability pass-through; revenue then wages; off-books.
- **Recommended**: liability (2140) with payout journal, excluded from revenue and VAT. SCN-22 PROVISIONAL.
- **Consequences**: cash held for staff, payroll reporting.
- **Affected**: orders, payments, payroll, Cash Flow.
- **Tests**: SCN-22; unpaid tips balance at close; card-tip refund.
- **Approval**: Accountant + payroll adviser.

### 15. Refunds - REFUND-01 (OPEN)
- **Current**: none: `refunded` rejected (`domain/order.go:219`), no refund table/endpoint/reversal/stock return (11; N2); only manual journal void, which dates at today with no period check (`journal_entry_repo.go:441,464-472`).
- **Options**: full/partial linked reversal; pre-settlement void; closed-period correction.
- **Recommended**: P4: immutable source references, idempotent linked reversals, never edit posted journals; post to contra Sales Returns dated refund date. SCN-13, 14.
- **Consequences**: cash, revenue, VAT, COGS/inventory, liabilities.
- **Affected**: `orders`, payments, GL, stock, reports, closed periods.
- **Tests**: P4 list plus SCN-13, 14, 24; duplicate request; refund exceeding paid.
- **Approval**: Accountant + payments owner + product owner.

### 16. Credit notes - D-CN-01 (OPEN; N16)
- **Current**: none for customer or supplier; supplier void leaves supplier balance unadjusted (11).
- **Options**: dedicated credit-note document linked to invoice; void-and-reissue only.
- **Recommended**: dedicated linked credit note document (SCN-05 shows supplier side).
- **Consequences**: AR/AP, VAT adjustment in the right period, e-invoice.
- **Affected**: invoices, AP/AR aging, VAT report.
- **Tests**: SCN-05; partial credit; credit over balance; aging.
- **Approval**: Accountant. Question: Egyptian credit-note/e-invoice requirements.

### 17. Payroll accrual - PAY-01 (OPEN)
- **Current**: no accrual; expense only on mark-paid (`auto_posting_engine.go:276-282`); 2130 never used (11); approval not required (N12); attendance not linked (N11).
- **Options**: accrue at period end; accrue at approval; cash basis.
- **Recommended**: accrue gross expense at period end/approval to Payroll Payable (SCN-18).
- **Consequences**: P&L completeness, liabilities at month end.
- **Affected**: salary records, posting engine, P&L, BS.
- **Tests**: SCN-18, 24; duplicate accrual; unpaid at close.
- **Approval**: Payroll lead + local adviser.

### 18. Payroll settlement - PAY-01 (OPEN)
- **Current**: Dr 5210/5220/tips; credits bank/cash and deductions; paid salaries immutable (409: `hr_salary_handler.go:320,413,451`); employer contributions not modelled (11); `PostPayrollEntryTx` atomic (`auto_posting_engine.go:268,288-292`).
- **Options**: gross + liabilities; net settlement with deduction lines.
- **Recommended**: gross to expense, deductions to liabilities, net via bank; remit liabilities separately (SCN-18). Deduction rate in oracle is ILLUSTRATIVE only.
- **Consequences**: statutory liabilities, cash timing.
- **Affected**: payroll page, GL 2110/2140/5210, bank.
- **Tests**: SCN-18; correction after paid; duplicate pay; statutory remittance.
- **Approval**: Payroll lead + local adviser. Question: Egyptian social insurance and tax deductions and employer share.

### 19. Bank transfers - D-BANK-01 (OPEN; N9)
- **Current**: none; bank balance is a mutable cached field adjusted by receipts/payments (11).
- **Options**: dedicated transfer document (two-leg journal); manual journal.
- **Recommended**: dedicated document, Cash and Bank both cash equivalents (CF-01), fee as separate expense (SCN-20).
- **Consequences**: no P&L or CF impact; reconciliation accuracy.
- **Affected**: bank accounts, `bank_accounts_ledger`, Cash Flow, reconciliation.
- **Tests**: SCN-20; cross-currency; transfer failure atomic.
- **Approval**: Accountant.

### 20. Capital contributions - D-EQ-01 (OPEN; N4)
- **Current**: no workflow; account 3100 Owner Capital seeded; Cash Flow financing keys off equity type (11).
- **Options**: dedicated event; manual journal.
- **Recommended**: dedicated event, financing inflow (SCN-02).
- **Consequences**: equity, CF.
- **Affected**: equity accounts, Cash Flow heuristic (`financial_report_usecase.go:389-391`).
- **Tests**: SCN-02; non-cash contribution (assets in kind).
- **Approval**: Owner + accountant.

### 21. Owner withdrawals - D-EQ-02 (OPEN; N4)
- **Current**: no Owner Drawings account, no document type (11).
- **Options**: contra-equity drawings; reduce capital; loan to owner (receivable).
- **Recommended**: Owner Drawings contra-equity, financing outflow, closed to retained earnings (SCN-21).
- **Consequences**: equity, tax treatment, CF.
- **Affected**: CoA seed, equity report, Cash Flow.
- **Tests**: SCN-21, 24; withdrawal exceeding equity warning.
- **Approval**: Owner + accountant.

### 22. Retained earnings - D-RE-01 (OPEN; CLOSE-01/EARN-01)
- **Current**: account 3200 seeded (`migrations/156_seed_chart_of_accounts.sql:58`); nothing posts to it (04 s7).
- **Options**: close each period; close only at fiscal year end; never close (derive).
- **Recommended**: close at fiscal year end (months stay open P&L), by journal, per CLOSE-01. SCN-23.
- **Consequences**: BS equity, comparatives.
- **Affected**: 3200, close routine, BS.
- **Tests**: SCN-23, 24; duplicate close; loss carried as debit.
- **Approval**: Accountant.

### 23. Historical corrections - D-HIST-01 (OPEN)
- **Current**: journals immutable (`180_enforce_journal_integrity.sql`); void appends inverse dated today (`journal_entry_repo.go:464-472`); paid salaries cannot be edited; no receipt void (11).
- **Options**: back-dated reversal into the original period; current-period adjustment; adjustment period.
- **Recommended**: if period open reverse in place; if closed, post adjustment in current open period with reference to original; never mutate history.
- **Consequences**: comparatives, VAT filed, audit trail.
- **Affected**: void paths, fiscal periods, reports.
- **Tests**: void in closed period; reversal date; report as-of stability.
- **Approval**: Accountant + owner.

### 24. Currency precision - ROUND-01 (OPEN)
- **Current**: Go `float64` money (`domain/accounting.go`, `order.go`); DB `numeric(15,2)`; hard-coded `"EGP"` (`order_pricing.go:99`); FX tables without posting (N24, FIN-020).
- **Options**: integer minor units; decimal strings; keep float with tolerance.
- **Recommended**: minor-unit integers/decimal in domain; 2 decimals for EGP; single functional currency until FX policy exists.
- **Consequences**: eliminates float drift in totals.
- **Affected**: all money code, reports, exports.
- **Tests**: boundary amounts, large sums, aggregation vs line sums.
- **Approval**: Accountant + owner.

### 25. Rounding - ROUND-01 (OPEN)
- **Current**: see cross-cutting item 6 plus 0.004 client tolerance (`order_usecase.go:163`), invoice trigger `ROUND(...,2)` (migration 141:80-82).
- **Options**: half-up; half-even; per-line vs per-document; residual allocation.
- **Recommended**: half-up at document level for VAT (oracle PROVISIONAL), explicit residual-cent allocation for splits, no tolerance in posting.
- **Consequences**: 1-piaster differences between receipt, invoice, VAT return and GL.
- **Affected**: pricing, invoices, posting validation, reports.
- **Tests**: half-piaster boundary, negative values (`round2` defect), split payments.
- **Approval**: Accountant + owner.

### Cross-cutting: Cash Flow classification - CF-01 (OPEN)
- **Current**: heuristic: `12%` investing, equity or `22%` financing, else operating (`financial_report_usecase.go:389-391`), not derived from cash accounts.
- **Recommended**: explicit mapping, direct method reconciled to cash+bank GL movement (SCN oracle design).
- **Tests**: every scenario CF reconciliation; transfer excluded.
- **Approval**: Accountant.

### Cross-cutting: Legacy supplier auto-post - CFR-011 (OPEN)
- Kept unreachable per P4 until a source lifecycle is approved. Not re-decided here.

## Owner decision record (to be completed by humans)

| ID | Outcome (APPROVED/REJECTED/DEFERRED) | Exact policy text | Effective date | Approver | Evidence |
|---|---|---|---|---|---|
| all rows above | blank | blank | blank | blank | blank |
