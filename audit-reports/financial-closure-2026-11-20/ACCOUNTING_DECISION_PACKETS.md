# SERVIO accounting decision packets

**Prepared:** 2026-11-20
**Code baseline:** `origin/main` `102c804df820e4e047853eb3d3ea51c0ce11472e`; the engineering candidate source/oracle test is `5c1113c`.
**Approval status:** all decisions below are **PENDING**. Proposals are engineering recommendations only, not accounting approvals. Examples are illustrative and do not authorize postings.

## Approval record

For each decision, the approver must reply with: decision ID; selected option or edited policy; effective date; applicable entities/periods; treatment of existing data; name, title and qualification; and explicit approval. Store the approval in the private governance record. A coding agent, test result, or product owner cannot substitute for a qualified financial approver where one is required.

The current requester has not supplied an approver identity or an authorized contact channel. No approval request has been sent and no decision has been recorded.

## OPEN-01 — Opening balances and offsets

- **Question:** What source records establish opening balances, the cutover date, and the offset account for each tenant/restaurant, and how are incomplete or late balances corrected?
- **Observed behavior:** balances are derived from posted journals; imported/bank balances without a complete opening journal can leave a report residual. Existing evidence does not establish the source or offset of the known fixture.
- **Options:** (A) approved trial balance by entity and cutover date, with an explicitly identified opening-equity/clearing account; (B) reconstruct each opening from source statements and post only documented components; (C) defer opening-balance imports until source support is complete.
- **Proposal:** prefer B when source records are available; otherwise use A with an approved, separately identified opening offset and auditable source reference. Never plug an unexplained difference.
- **Illustrative entry:** Dr cash 10,000; Dr receivables 2,000; Cr payables 3,000; Cr opening equity/clearing 9,000. Amounts are examples only.
- **Impact / implementation:** historical statements depend on cutover date, account mapping and source completeness. Requires an import contract, idempotent source IDs, reviewable preview, scoped posting, and possibly migration/backfill tooling. Do not rewrite posted history.
- **Acceptance:** independently tie each opening account to source evidence; opening debits equal credits; verify date boundary, duplicate import rejection, tenant/restaurant isolation, report cutoffs, and export parity.
- **Approver / status / effective date:** qualified accountant/controller; **PENDING**; none recorded.

## EARN-01 — Current earnings and retained earnings

- **Question:** Should an open-period Balance Sheet include cumulative current-period earnings in equity, and when/how are earnings transferred to retained earnings?
- **Observed behavior:** Balance Sheet reports asset, liability and equity account balances but excludes revenue/expense accounts; the known synthetic fixture shows assets 5,500¢, liabilities 500¢, equity 0¢ and unclosed net income 5,000¢, leaving a 5,000¢ residual. This alone does not establish whether opening offsets or other ledger defects also exist.
- **Options:** (A) present current-period net income within equity as a report calculation without a journal; (B) transfer net income/loss to retained earnings in an approved close journal; (C) use A during open periods and B at close.
- **Proposal:** C is the cleanest operational model if the close lifecycle and reopening rules are approved; do not activate it until OPEN-01 and CLOSE-01 are settled.
- **Illustrative entry for an approved close only:** Dr income summary 5,000; Cr retained earnings 5,000, with account polarity and sign conventions confirmed against the chart of accounts. No entry is authorized now.
- **Impact / implementation:** A changes report aggregation only; B changes ledger history and needs close-period state, audit, reversal/reopen and migration compatibility. Preserve closed-period reports and effective-date semantics.
- **Acceptance:** open and closed period reports tie to independently summed posted journal lines; no duplicate income; multiple fiscal years, loss periods, prior-period adjustments, reopening and exports reconcile; residual equals zero only when underlying statements justify it.
- **Approver / status / effective date:** qualified accountant/controller; **PENDING**; none recorded.

## CLOSE-01 — Fiscal close, reopening and corrections

- **Question:** What closes a period, who may close/reopen it, and how are late or erroneous transactions recorded after close?
- **Observed behavior:** period controls exist, but full close/earnings transfer/reopen acceptance is absent.
- **Options:** (A) hard lock with linked reversal and next-period correction; (B) controlled reopen with reason, approver and immutable audit; (C) adjustment period with explicit reporting rules.
- **Proposal:** require a controlled, audited close; use linked reversals and next-open-period corrections by default, and permit reopen only through an authorized approval flow.
- **Illustrative entries:** original expense Dr 100 / Cr payable 100; correction after close Dr payable 100 / Cr expense or approved prior-period adjustment account 100. Classification/date require approval.
- **Impact / implementation:** period-state schema, role checks, audit trail, APIs, reports and migration of existing period records may change. No silent mutation of posted history.
- **Acceptance:** unauthorized close/reopen denied; duplicate close idempotent; close totals tie; corrections retain original references; closed-period reproducibility and multi-year report tests pass.
- **Approvers / status / effective date:** qualified accountant/controller and operational owner for close schedule; **PENDING**; none recorded.

## TAX-01 — Tax rates, jurisdiction and basis

- **Question:** Which jurisdictions, tax types/rates, inclusive or exclusive basis, discount ordering, exemptions and effective dates apply to each sale/invoice?
- **Observed behavior:** hardcoded order tax and unresolved invoice tax/discount sequencing remain identified; the applicable jurisdiction and approved rate table are not established here.
- **Options:** versioned jurisdiction-specific rules; explicit tax-inclusive pricing; explicit tax-exclusive pricing; tax-exempt customer/product cases with documented evidence. These can vary by product/jurisdiction.
- **Proposal:** use effective-dated, versioned rules from a qualified tax adviser, with tax calculated from immutable line-level inputs and retained rule IDs. No default rate is proposed without jurisdiction facts.
- **Illustrative entry:** Dr receivable 115; Cr revenue 100; Cr tax payable 15 for a 15% exclusive example only.
- **Impact / implementation:** invoice/order APIs, tax persistence, historical recalculation policy, frontend display and migrations. Existing posted invoices must not be silently recomputed.
- **Acceptance:** jurisdiction/rate/date matrix; inclusive/exclusive and discount ordering; exemptions; refunds/credit notes; line and document rounding; source-to-report/export tie-out; historical rule reproducibility.
- **Approver / status / effective date:** qualified tax adviser/accountant; **PENDING**; none recorded.

## ROUND-01 — Monetary precision and rounding

- **Question:** What currency precision, calculation precision, rounding mode, point of rounding and permitted tolerance apply at line, document, tax and ledger levels?
- **Observed behavior:** monetary models/APIs still contain `float64`; no approved precision or tolerance policy is recorded.
- **Options:** integer minor units for persisted amounts with higher-precision decimal intermediates; decimal values with explicit scale; currency-specific scales and rounding rules.
- **Proposal:** persist integer minor units where currency supports a fixed minor unit; use decimal arithmetic for rates/intermediates and a named, currency-specific rounding rule at approved boundaries. Reject non-finite inputs.
- **Illustrative entry:** three lines of 0.335 at a two-decimal currency can round per-line to 0.34 each (1.02 total) or after summing to 1.01; policy must choose.
- **Impact / implementation:** schema/API contracts, imports, reports, calculations and historical conversion. Conversion requires a reversible preflight and independently reconciled migration; no blanket float coercion.
- **Acceptance:** property tests for arithmetic, boundary/negative/large values, tax and FX; historical round-trip; API validation; database/report/export equality; cross-currency scale tests.
- **Approver / status / effective date:** qualified accountant/controller; **PENDING**; none recorded.

## REFUND-01 — Refunds, chargebacks, voids and reversals

- **Question:** For each order channel/provider, when is revenue recognized, how are partial refunds/fees/tax/tips allocated, how are chargebacks represented, and what happens after period close?
- **Observed behavior:** subscription refund paths exist, but no complete restaurant-order provider/refund/chargeback flow or linked accounting reversal contract was established. Paid-order transitions are guarded; cancellation is not a refund.
- **Options:** event-specific reversal linked to original journal; credit-note workflow; pending provider claim followed by settlement/reversal; distinct void-before-capture and refund-after-capture paths.
- **Proposal:** require stable provider event identity and durable reconciliation state; reverse only confirmed refunded amounts through a linked immutable event; keep inventory return separate from payment reversal.
- **Illustrative entry:** Dr sales returns 40; Dr tax payable 6; Cr cash/refund payable 46 for an approved 40+6 example only.
- **Impact / implementation:** provider contracts, event/outbox tables, retries, account mapping, audit and closed-period rules; migration and reconciliation needed. No provider semantics are assumed.
- **Acceptance:** full/partial, duplicate, out-of-order, concurrent and unknown-outcome callbacks; provider-success/local-failure recovery; chargeback/fee mismatch; scope; linked journal; independent settlement tie-out.
- **Approvers / status / effective date:** qualified accountant/controller and product/payment owner; **PENDING**; none recorded.

## INV-01 — Inventory valuation and COGS recognition

- **Question:** Which costing method, negative-stock behavior, wastage/returns treatment and COGS recognition point apply?
- **Observed behavior:** inventory movements and COGS paths exist; valuation and timing are not accepted.
- **Options:** weighted average, FIFO, specific identification, or a documented non-perpetual method where allowed; recognize COGS on fulfillment/consumption or another approved event.
- **Proposal:** choose one valuation policy per supported inventory scope and record immutable cost layers/movements; reject silent negative inventory unless an approved fallback exists.
- **Illustrative entry:** on consumption of inventory costing 30, Dr COGS 30 / Cr inventory 30.
- **Impact / implementation:** movement schema, backfill/cost layers, sales/return/waste flows, reports and migrations. Historical costs may need source-led reconciliation.
- **Acceptance:** purchase, sale, partial return, wastage, transfer, negative stock, backdate and replay tests; inventory subledger equals GL at every cutoff.
- **Approver / status / effective date:** qualified accountant/controller and inventory/product owner; **PENDING**; none recorded.

## PAY-01 — Payroll deductions, posting and correction

- **Question:** Which statutory/voluntary deductions, employer costs, payment timing, rounding and correction method apply for each jurisdiction and pay period?
- **Observed behavior:** paid payroll mutations are guarded, but statutory deductions and correction rules remain unapproved.
- **Options:** jurisdiction-effective deduction configuration; accrual then payment; direct cash-basis posting where approved; correction by linked adjustment rather than editing paid records.
- **Proposal:** keep paid payroll immutable; post approved gross, liability, employer-cost and payment components using effective-dated payroll rules and linked corrections.
- **Illustrative entry:** Dr wage expense 1,000; Cr withholding payable 150; Cr payroll payable 850, followed on payment by Dr payroll payable 850 / Cr cash 850.
- **Impact / implementation:** payroll configuration, journal mapping, reports, privacy/RBAC and historical rule references. Existing paid records are not recalculated automatically.
- **Acceptance:** each deduction and employer cost ties to payroll source; paid immutability; duplicate/concurrent mark-paid; failed posting rollback; correction audit; tenant/restaurant/export privacy.
- **Approver / status / effective date:** qualified payroll adviser/accountant; **PENDING**; none recorded.

## CF-01 — Cash-flow method and classifications

- **Question:** Direct or indirect method, cash-equivalent definition, restricted cash, interest/tax classification, non-cash items and operating/investing/financing mapping?
- **Observed behavior:** no accepted cash-flow classification/method or independently reconciled cash-flow fixture.
- **Options:** direct cash receipts/payments; indirect reconciliation from net income; approved hybrid/statement format.
- **Proposal:** choose the method required by the applicable reporting framework and jurisdiction; make classification rules effective-dated and auditable.
- **Illustrative entry:** customer cash receipt Dr cash 100 / Cr receivable 100; cash-flow classification is a reporting dimension and must not create an extra journal.
- **Impact / implementation:** report query, classification metadata, account mapping and export. Historical restatement depends on approved rules.
- **Acceptance:** cash opening + classified movements = cash closing; non-cash activity excluded or reconciled; each section tied to ledger and independent source; cutoff and export tests.
- **Approver / status / effective date:** qualified accountant/controller; **PENDING**; none recorded.

## CFR-011 — Legacy supplier posting contract

- **Question:** Are `PostSupplierInvoiceEntry` and `PostSupplierPaymentEntry` supported product contracts? If yes, what source lifecycle, event identity, account mapping, authorization and idempotency rules govern them? If no, what supersedes them?
- **Observed behavior:** schema width was corrected by migration 187; repository methods have no production callers and are absent from `AutoPostingService`. This does not establish whether they are supported.
- **Options:** (A) implement as supported source events; (B) confirm supersession by live invoice/payment flows and deprecate/remove through a reviewed migration; (C) retain dormant methods but explicitly mark unsupported.
- **Proposal:** product owner and accountant should choose A or B; preserve compatibility until that decision. Do not add new endpoints or financial postings by inference.
- **Illustrative entry if approved:** invoice Dr expense/inventory 100 / Cr accounts payable 100; payment Dr accounts payable 100 / Cr cash 100.
- **Impact / implementation:** API/use-case contract, mapping, migration/deprecation, replay and reconciliation. Legacy data must be checked before removing any schema or method.
- **Acceptance:** approved source lifecycle; duplicate/replay, rollback/failure, tenant scope, invoice/payable and payment/cash subledger ties; migration upgrade/rollback.
- **Approvers / status / effective date:** product owner and qualified accountant; **PENDING**; none recorded.

## Cross-decision dependencies

Resolve OPEN-01 with EARN-01 before using the known Balance Sheet fixture as proof of a report-only defect. CLOSE-01 controls earnings transfer and closed-period corrections. TAX-01 and ROUND-01 govern invoice/refund and payroll amounts. REFUND-01 depends on provider and product event contracts. INV-01 must reconcile to sales/returns. CF-01 depends on approved account/event classifications. No production policy change is authorized by this packet.
