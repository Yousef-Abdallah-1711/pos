# 56b - Policy decision detail (RUN-005, all decisions OPEN)

Companion to [56](56-accounting-policy-decision-package.md). Built by `tests/build_policy_detail.py`. **No decision is approved**; the recommendation is the auditor's proposal with rationale, kept separate from any approved decision. Dependent tasks stay NOT STARTED and policy-dependent calculations in the reports are conditional until a decision is recorded.

## POL-01 - Revenue recognition basis for Sales, Profitability, KPI, order stats and the GL
- **Status:** OPEN
- **Current product behaviour / evidence:** Revenue bridge `evidence/database/RUN-003/PH-12/AUD-REV-TRACE.txt`: Sales 517,000 (tax-incl, unpaid incl) → Profitability 470,000 (net, unpaid incl) → GL 290,000 (net, paid only); KPI 319,660 (tax-incl, paid); order stats 529,100 (all incl cancelled)
- **Business question and alternatives:** (a) order-taking (all non-cancelled orders); (b) payment (paid only); (c) fulfilment (delivered/served)
- **Recommended option (proposal):** (b) payment basis for the GL and every report; show unpaid orders as a separate "pipeline" figure
- **Rationale:** Payment basis matches the GL, which already posts revenue at payment; the other reports then reconcile to it with a labelled bridge.
- **Financial example:** Sales 517,000 (tax-inclusive, unpaid included) vs Profitability 470,000 (net, unpaid included) vs GL 290,000: tax 47,000 + unpaid 180,000 explain the gap; under option (b) all three show 290,000 net.
- **Tax and report implications:** Sales, Profitability, KPI and order statistics change on deploy; every revenue figure gets a basis label; VAT figures are unaffected (they come from the GL).
- **Historical data:** No stored data changes; communicate the change of displayed totals.
- **Migration or backfill:** None (queries and labels).
- **Dependent findings:** F-012, F-042, F-043
- **Dependent future tasks:** 19 (P0 policy tasks: 3); **future regression tests:** 10

## POL-02 - Tax base and tax display
- **Status:** OPEN
- **Current product behaviour / evidence:** F-007/F-041: 10,000 gross, 10% discount, 14% VAT → system 1,400 vs discount-first 1,260; POS cart shows 14% vs stored 10% (F-004)
- **Business question and alternatives:** (a) tax on net-of-discount amount; (b) tax on gross; (c) tax-inclusive prices
- **Recommended option (proposal):** (a) discount-first in both invoice triggers; single configurable tax rate per tenant read by POS and server
- **Rationale:** Discount-first is the usual VAT base (tax on the amount actually charged) and matches the UI preview.
- **Financial example:** Gross 10,000, 10% discount, 14% VAT: discount-first tax 1,260 (total 10,260) vs current 1,400 (total 10,400): 140 over-declared VAT per such invoice; RUN-005 real-UI: UI 630.00 vs stored 700.00 on 5,000.
- **Tax and report implications:** VAT payable, invoice totals, AR and revenue split change for discounted lines; VAT return totals must be reviewed per period.
- **Historical data:** Posted invoices keep their tax; produce an exception report and correct through approved reversing journals.
- **Migration or backfill:** Replace the line-total triggers (migrations 141/142 logic) with a new migration; recalculate drafts only.
- **Dependent findings:** F-004, F-007, F-010, F-041
- **Dependent future tasks:** 30 (P0 policy tasks: 4); **future regression tests:** 13

## POL-03 - Inventory accounting
- **Status:** OPEN
- **Current product behaviour / evidence:** OBS-PUR-001: supplier invoice debits 5110 not 1130; F-011: valuation 132,840 vs GL 1130 −96,120; F-034 waste posts nothing; F-040 untracked product posts no COGS
- **Business question and alternatives:** (a) capitalise purchases to 1130 and relieve on sale; (b) expense purchases to COGS; (c) periodic count adjustment
- **Recommended option (proposal):** (a) perpetual inventory with opening-stock journal and waste/shrinkage expense account
- **Rationale:** Perpetual inventory gives a true COGS and a balance-sheet inventory that equals the stock ledger.
- **Financial example:** Inventory valuation 132,840 vs GL 1130 -96,120 (difference 228,960) because purchases are expensed to 5110 and waste posts nothing.
- **Tax and report implications:** COGS, gross profit, balance sheet inventory, valuation and turnover reports change.
- **Historical data:** Opening-stock journal and a catch-up of purchases already expensed; accountant approval.
- **Migration or backfill:** Posting-rule changes and a data-correction script; no new columns required.
- **Dependent findings:** F-011, F-034, F-040, F-046
- **Dependent future tasks:** 26 (P0 policy tasks: 4); **future regression tests:** 12

## POL-04 - Cash-flow statement method and classification
- **Status:** OPEN
- **Current product behaviour / evidence:** F-003: report 69,290 vs cash change 74,480 on clean tenant; asset purchase raised investing
- **Business question and alternatives:** direct (cash accounts) vs indirect (from net income with working-capital terms)
- **Recommended option (proposal):** indirect method derived from GL cash/bank accounts with classification table by account range
- **Rationale:** The indirect method derived from GL cash and bank movements is checkable against the cash accounts.
- **Financial example:** Clean tenant: report 69,290 vs cash change 74,480 (residual 5,190 = working-capital terms); RUN-004: a 20,000 asset purchase paid in cash shows investing +20,000 instead of -20,000.
- **Tax and report implications:** Cash-flow statement and the predictive cash flow (F-045) must use the same method.
- **Historical data:** None stored; historical statements differ when regenerated.
- **Migration or backfill:** None.
- **Dependent findings:** F-003, F-045
- **Dependent future tasks:** 14 (P0 policy tasks: 2); **future regression tests:** 7

## POL-05 - Fiscal period control
- **Status:** OPEN
- **Current product behaviour / evidence:** JE-13/JE-14 posted for 2019 and 2030 with no period; GET /accounting/fiscal-periods is 404
- **Business question and alternatives:** required open period per posting date? backdating window? who may close/reopen?
- **Recommended option (proposal):** require an open period for the entry date; closing requires accountant role; lock after VAT filing
- **Rationale:** An open-period rule is the baseline control for any ledger; it prevents 2019 or 2030 postings with no period.
- **Financial example:** Journals dated 2019-01-15 and 2030-06-01 were posted with no fiscal period.
- **Tax and report implications:** Period close, trial balance cut-offs and VAT filing windows depend on it.
- **Historical data:** List entries outside any period; decide per entry.
- **Migration or backfill:** Fiscal-period tables and a guard in the posting path.
- **Dependent findings:** F-013
- **Dependent future tasks:** 7 (P0 policy tasks: 1); **future regression tests:** 3

## POL-06 - Refund / void / return policy for paid orders
- **Status:** OPEN
- **Current product behaviour / evidence:** F-049: every endpoint rejects paid-order cancel/refund; revenue, VAT, COGS stay posted
- **Business question and alternatives:** full vs partial; tender-specific reversal; restock; approval
- **Recommended option (proposal):** reversing journals (revenue, VAT, COGS) in the original period or the current period per policy POL-05; restock by flag
- **Rationale:** Reversing journals preserve the audit trail and restock keeps inventory true.
- **Financial example:** Refund of a 110.00 paid order (100 net + 10 VAT, COGS 40): Dr revenue 100, Dr VAT payable 10, Cr cash 110; Dr inventory 40, Cr COGS 40.
- **Tax and report implications:** Sales, profitability, VAT and cash reports must exclude refunded amounts consistently.
- **Historical data:** None (new capability).
- **Migration or backfill:** Refund tables and status values.
- **Dependent findings:** F-049, F-050
- **Dependent future tasks:** 11 (P0 policy tasks: 1); **future regression tests:** 4

## POL-07 - Non-cash tender settlement
- **Status:** OPEN
- **Current product behaviour / evidence:** F-005 card/wallet orders never become paid; F-026 manual paid debits clearing 1113; checkout disabled in environment
- **Business question and alternatives:** card/wallet clearing, settlement timing, fees
- **Recommended option (proposal):** provider-confirmed paid status; clearing account settled to bank on statement import
- **Rationale:** A clearing account settled by an explicit bank deposit shows timing differences and fees.
- **Financial example:** Card sale 110.00 debits 1113; a settlement deposit of 107.80 (2% fee) must move 110.00 out of 1113 with 2.20 to fees.
- **Tax and report implications:** Balance sheet (1113 balance), cash reports and bank reconciliation.
- **Historical data:** Opening reconciliation of existing 1113 balances.
- **Migration or backfill:** Settlement tables.
- **Dependent findings:** F-005, F-026
- **Dependent future tasks:** 16 (P0 policy tasks: 2); **future regression tests:** 6

## POL-08 - Payroll recognition
- **Status:** OPEN
- **Current product behaviour / evidence:** F-015: salaries paid with is_approved=false; posting only at payment; JE credits 2110 withholdings
- **Business question and alternatives:** accrual at approval vs expense at payment; liabilities account
- **Recommended option (proposal):** accrue at approval (Dr 5210/5220, Cr 2110 + payable), clear at payment
- **Rationale:** Accruing at approval puts the expense in the right period and exposes unpaid payroll as a liability.
- **Financial example:** Salary gross 6,000, deductions 750, net 5,250: accrual Dr 5210 6,000 / Cr 2110 750 and payable 5,250; payment Dr payable / Cr cash 5,250.
- **Tax and report implications:** P&L by period, balance sheet liabilities, HR cost report.
- **Historical data:** Salaries 1 and 2 were paid without approval (list for sign-off).
- **Migration or backfill:** Payment records and accrual posting.
- **Dependent findings:** F-015, F-016
- **Dependent future tasks:** 14 (P0 policy tasks: 1); **future regression tests:** 7

## POL-09 - Accounting effect of draft documents
- **Status:** OPEN
- **Current product behaviour / evidence:** F-008: draft customer invoice posts AR/revenue/VAT immediately
- **Business question and alternatives:** post at creation vs at issue
- **Recommended option (proposal):** post only on issue (status sent); drafts have no ledger effect
- **Rationale:** Drafts are working documents; ledger effects belong to issued documents.
- **Financial example:** Draft invoice INV-2026-00005 (5,200) posted Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700 at creation.
- **Tax and report implications:** AR aging, revenue and VAT are overstated while drafts exist.
- **Historical data:** Reverse journals of existing drafts after listing them.
- **Migration or backfill:** None (posting trigger point moves).
- **Dependent findings:** F-008
- **Dependent future tasks:** 7 (P0 policy tasks: 1); **future regression tests:** 2

## POL-10 - Receipt exceeding invoice balance
- **Status:** OPEN
- **Current product behaviour / evidence:** F-006: 5 parallel receipts on a 3,000 invoice accepted (6,000)
- **Business question and alternatives:** reject, or hold as customer credit/unapplied cash
- **Recommended option (proposal):** reject in the same transaction under row lock; optional explicit "customer credit" flow
- **Rationale:** Rejecting under a row lock is simplest; customer credit can be an explicit later feature.
- **Financial example:** Five parallel receipts of 1,000 on a 3,000 invoice were all accepted: 5,000 received, 2,000 over-collected.
- **Tax and report implications:** AR subledger vs GL 1120 (residual 1,800 in RUN-003).
- **Historical data:** List over-settled invoices; decide refund or credit per case.
- **Migration or backfill:** None (locking and validation).
- **Dependent findings:** F-006, F-052
- **Dependent future tasks:** 16 (P0 policy tasks: 2); **future regression tests:** 6

## POL-11 - Currency policy
- **Status:** OPEN
- **Current product behaviour / evidence:** F-023: POS shows $; storefront shows EGP; F-048 same-currency FX pair accepted
- **Business question and alternatives:** single functional currency EGP vs multi-currency; symbol source
- **Recommended option (proposal):** functional currency per tenant; single formatter reading it
- **Rationale:** One functional currency per tenant read by every formatter prevents symbol drift.
- **Financial example:** POS shows $12.00 and $10,000.00 while the dashboard shows EGP 110.00 for the same tenant.
- **Tax and report implications:** All money displays and exports.
- **Historical data:** None.
- **Migration or backfill:** Tenant currency setting.
- **Dependent findings:** F-023, F-048
- **Dependent future tasks:** 11 (P0 policy tasks: 2); **future regression tests:** 5

## POL-12 - Single source of truth for dashboards/KPIs
- **Status:** OPEN
- **Current product behaviour / evidence:** F-043 KPI net profit 181,540 vs P&L 167,180; cash position 50,000 vs GL cash 276,270; F-042
- **Business question and alternatives:** GL-derived vs operational-table-derived
- **Recommended option (proposal):** derive every financial KPI from the GL; operational counts from order tables
- **Rationale:** A single GL-derived source removes the competing revenue and cash figures.
- **Financial example:** KPI net profit 181,540 vs P&L 167,180 (difference 14,360); KPI cash position 50,000 vs GL cash 276,270.
- **Tax and report implications:** Dashboards and KPI endpoints.
- **Historical data:** None.
- **Migration or backfill:** None.
- **Dependent findings:** F-012, F-042, F-043
- **Dependent future tasks:** 19 (P0 policy tasks: 3); **future regression tests:** 10

## POL-13 - Recurring-journal idempotency
- **Status:** OPEN
- **Current product behaviour / evidence:** F-047: two runs same day → two journals
- **Business question and alternatives:** one run per period vs unlimited manual runs
- **Recommended option (proposal):** unique (template, period)
- **Rationale:** A unique (template, period) key makes the scheduled run idempotent.
- **Financial example:** Two runs on one day created two journals of the same amount.
- **Tax and report implications:** Recurring expenses in P&L double.
- **Historical data:** Reverse duplicate journals.
- **Migration or backfill:** Unique index.
- **Dependent findings:** F-047
- **Dependent future tasks:** 8 (P0 policy tasks: 1); **future regression tests:** 3

## POL-14 - Leave entitlement
- **Status:** OPEN
- **Current product behaviour / evidence:** F-054: 30-day annual leave approved over allowance 21
- **Business question and alternatives:** block, allow as unpaid, or allow negative balance
- **Recommended option (proposal):** allow as unpaid beyond entitlement with payroll deduction
- **Rationale:** Allowing the excess as unpaid leave keeps the request workflow usable and protects payroll.
- **Financial example:** Approved 30 days against an allowance of 21: 9 days unpaid (daily rate = base salary / 30).
- **Tax and report implications:** HR leave balance report (remaining currently clamps to 0) and payroll deductions.
- **Historical data:** List over-allowance approvals.
- **Migration or backfill:** Optional unpaid_days column.
- **Dependent findings:** F-054
- **Dependent future tasks:** 6 (P0 policy tasks: 1); **future regression tests:** 3

## POL-15 - Default role matrix
- **Status:** OPEN
- **Current product behaviour / evidence:** F-020/F-033: default Employee role reads HR records incl. base_salary; no cashier role
- **Business question and alternatives:** which default roles exist; salary visibility
- **Recommended option (proposal):** add Cashier/Accountant/HR roles; field-level masking of salary data
- **Rationale:** Least privilege: salary data only for HR and owners; a Cashier role for POS staff.
- **Financial example:** EMPLOYEE role reads base_salary 5,000 from /hr/employees; MANAGER reads /hr/salaries (RUN-005 matrix).
- **Tax and report implications:** HR, payroll and user screens.
- **Historical data:** Review access logs.
- **Migration or backfill:** Seed new roles for existing tenants.
- **Dependent findings:** F-020, F-027, F-033
- **Dependent future tasks:** 19 (P0 policy tasks: 3); **future regression tests:** 10

## POL-16 - Stock reservation for unpaid orders (PROPOSED in RUN-004)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-005/F-028: card/wallet orders stay unpaid yet stock is consumed; load test created 5,719 unpaid cash orders
- **Business question and alternatives:** consume stock at creation (current) vs reserve with expiry vs consume at payment
- **Recommended option (proposal):** reserve at creation with an expiry (e.g. 30 min) that releases stock
- **Rationale:** A reservation with expiry prevents both phantom stock consumption and oversell.
- **Financial example:** The RUN-004 load test created 5,719 unpaid cash orders; tracked stock 12 was consumed by 12 of them regardless of payment.
- **Tax and report implications:** Inventory, orders and KDS counts.
- **Historical data:** Review stale unpaid orders.
- **Migration or backfill:** Expiry job.
- **Dependent findings:** F-005, F-028
- **Dependent future tasks:** 13 (P0 policy tasks: 2); **future regression tests:** 5

## POL-17 - Bank opening balance treatment (PROPOSED in RUN-004)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-029: opening balance 50,000 not in the GL
- **Business question and alternatives:** post to opening-balance equity at the opening date vs require manual journal
- **Recommended option (proposal):** post Dr bank / Cr opening balance equity at creation
- **Rationale:** Opening equity posting at the opening date keeps the bank module and the GL equal.
- **Financial example:** Bank account opened with 50,000: module balance 50,000, GL contribution 0.
- **Tax and report implications:** Balance sheet, KPI cash.
- **Historical data:** Correction journals for existing accounts.
- **Migration or backfill:** None.
- **Dependent findings:** F-029
- **Dependent future tasks:** 7 (P0 policy tasks: 1); **future regression tests:** 3

## POL-18 - Driver cash handling and shortages (PROPOSED in RUN-004)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-057: COD orders stay unpaid after delivery
- **Business question and alternatives:** driver receivable settled daily; shortage expense vs recoverable
- **Recommended option (proposal):** driver receivable, daily settlement, shortage to a variance account
- **Rationale:** A driver receivable makes cash risk visible.
- **Financial example:** A delivered cash order of 110.00 stays unpaid with no record of the cash the driver holds.
- **Tax and report implications:** Cash, AR/driver receivable.
- **Historical data:** Reconcile delivered unpaid orders.
- **Migration or backfill:** Driver settlement tables.
- **Dependent findings:** F-057
- **Dependent future tasks:** 8 (P0 policy tasks: 1); **future regression tests:** 3

## POL-19 - Termination, notice period and final settlement (PROPOSED in RUN-004)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-069: completed resignation leaves the employee active
- **Business question and alternatives:** terminate at last working day; notice period handling; leave payout
- **Recommended option (proposal):** terminate at last working day; settlement per resignation record
- **Rationale:** Termination at the last working day keeps payroll and turnover consistent.
- **Financial example:** Resignation completed with last working day 2026-11-25: employee stays active, turnover 0 instead of 100% for the period.
- **Tax and report implications:** HR turnover, payroll eligibility.
- **Historical data:** List completed resignations with active employees.
- **Migration or backfill:** None.
- **Dependent findings:** F-056, F-069
- **Dependent future tasks:** 11 (P0 policy tasks: 1); **future regression tests:** 5

## POL-20 - Reporting time zone and day boundary (PROPOSED in RUN-005)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-070: KDS hours are UTC buckets (Cairo buckets differ); F-042 delivered_today and report cut-offs use CURRENT_DATE
- **Business question and alternatives:** UTC day (current) vs tenant time zone (Africa/Cairo default) for all date-bounded reports, KDS and exports
- **Recommended option (proposal):** tenant time zone setting, default Africa/Cairo; store timestamps with time zone
- **Rationale:** Business days are local; UTC buckets mislead by 2-3 hours.
- **Financial example:** Same KDS rows: UTC busiest hours 1:7, 20:18, 22:11 vs Cairo 3:7, 22:18, 0:11.
- **Tax and report implications:** KDS, daily sales and any CURRENT_DATE report; exports.
- **Historical data:** Recomputed on read; no stored change.
- **Migration or backfill:** Tenant time-zone setting.
- **Dependent findings:** F-042, F-070
- **Dependent future tasks:** 12 (P0 policy tasks: 1); **future regression tests:** 6

## POL-21 - Bank reconciliation completion rule (PROPOSED in RUN-005)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-072: session completed with an invented 777,777.77 line; F-073 cross-tenant journal reference
- **Business question and alternatives:** zero tolerance vs 0.01 vs approved adjustment journal; who may complete
- **Recommended option (proposal):** book balance from the GL, completion only at zero difference (or with a recorded adjustment journal), finance permission to complete
- **Rationale:** Completion only at a zero difference against the GL gives the reconciled flag meaning.
- **Financial example:** A session with statement 777,777.77 and an invented line completed while the bank ledger balance was 12,000.00.
- **Tax and report implications:** Bank balance reports and last-reconciled fields.
- **Historical data:** Re-verify completed sessions against the GL.
- **Migration or backfill:** Line constraints.
- **Dependent findings:** F-029, F-072, F-073
- **Dependent future tasks:** 20 (P0 policy tasks: 1); **future regression tests:** 7

## POL-22 - Subscription trial extension limits (PROPOSED in RUN-005)
- **Status:** OPEN
- **Current product behaviour / evidence:** F-055: 100,000 days accepted (end year 2300)
- **Business question and alternatives:** unbounded (current) vs capped with reason and history
- **Recommended option (proposal):** cap 90 days, reason required, history row
- **Rationale:** A cap plus reason and history keeps extensions accountable.
- **Financial example:** 100,000 days moved a trial end to the year 2300.
- **Tax and report implications:** Super Admin tenant status.
- **Historical data:** List tenants beyond the cap.
- **Migration or backfill:** History table.
- **Dependent findings:** F-055
- **Dependent future tasks:** 5 (P0 policy tasks: 0); **future regression tests:** 3
