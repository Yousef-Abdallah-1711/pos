# Financial Reconciliation Results
> **Phase 2 execution update (2026-11-14):** `origin/main` was fetched and verified at `165837c7ad2e506dbf20b1a45c9a123476f2f8e9`. Focused CFR-010 A–D PostgreSQL fixtures pass against a newly provisioned schema-only disposable database; independent SQL confirms Fixture A balanced opening (assets/equity 1,000/1,000), Fixture B off-ledger opening reproduces the historical -100 EGP residual, Fixture C expense and Fixture D two-period/close cases reconcile to their independent figures. Migrations `186_accounting_event_redrive_audit.sql` (PR #37), `186_servio_trial_offer.sql` (main), and `187_expand_journal_entry_type_to_32.sql` (PR #44) applied together in fresh and upgraded disposable schemas; long `auto_supplier_*` values persisted, and unsafe rollback refused. This does **not** execute authenticated supplier auto-posting: the only production references to those methods/constants are definitions, with tests as the only callers. Manual browser acceptance remains blocked; no credentialed synthetic browser session is available. PR #46 adds the regression and its CI is green on all required jobs. Financial production remains **NO-GO**.

> **Execution update (2026-11-14):** A later audit pass ran a verified local disposable stack: 1 authenticated supplier invoice/payment API workflow, 6 automated API-backed report views, 1 unauthenticated authorization-negative request, and 0 manual financial browser workflows. Independent SQL reconciliation found source/payment/GL/TB/P&L/AP/AR differences of 0.00 for the narrow fixture, but a -100.00 EGP Balance Sheet residual because opening cash is absent from GL. Local isolation PASS for observed controls; existing staging isolation UNKNOWN. Financial release remains **NO-GO**. Detailed scope, values, commands, and limitations: [disposable execution evidence](evidence/disposable-execution-2026-11-14.md). This update supersedes earlier statements of no local fixture where applicable.

**Baseline:** main `018682beb9448eea6dd72fd2418876b439a94dbf`.  
**Independent-oracle execution:** one synthetic supplier invoice/payment workflow ran against the disposable candidate stack, not an exact mainline release. Expected amounts were recomputed from synthetic source and posted GL rows using read-only SQL, not report functions. This is a narrow fixture and includes a Balance Sheet failure.

## Result summary

| Reconciliation boundary | Expected | Actual | Difference | Expected-value source | Result |
|---|---|---|---|---|---|
| Supplier invoice/payment → accounting event | Invoice/payment fixture persisted through real authenticated API | 100.00 invoice; 100.00 completed payment | 0.00 | Synthetic source lines/payment rows independently summed with read-only SQL | PASS for this workflow; not an order/payment scenario |
| Event → journal | One balanced journal per canonical event | No current-main end-to-end event fixture run in this audit | N/E | None | Blocked |
| Journal → general ledger | Every posted line appears exactly once with same scope/amount | Two posted journals at 100.00 debit/credit each; raw GL 200.00/200.00 | 0.00 | Read-only SQL over persisted synthetic journal/GL rows | PASS for fixture |
| GL → Trial Balance | Debit total equals credit total from independently recomputed posted lines | Raw GL 200.00/200.00; net TB 100.00/100.00 | 0.00 | Read-only sums of synthetic posted GL rows | PASS for fixture |
| P&L → Balance Sheet equity | Current-period earnings/approved close roll-forward tie | P&L net -100.00; BS assets -100.00, liabilities 0.00, equity 0.00; report unreconciled | -100.00 residual | Net P&L checked against net GL; opening cash 1,000.00 is absent from GL | FAIL; opening-balance treatment unresolved |
| Cash movements → Cash Flow | Opening cash 1,000.00 + net movement -100.00 = closing cash 900.00 | Cash Flow showed -100.00 | 0.00 | Independent arithmetic over synthetic opening/closing bank values | Arithmetic PASS; accountant approval needed for classification |
| AR/AP control accounts → invoice aging | Invoice paid 100.00; open balance and aging 0.00 | AP aging 0.00; AR aging 0.00 | 0.00 | Synthetic invoice/payment rows and read-only SQL | PASS for this supplier fixture |
| Payroll → GL | Paid salary obligations/payments agree to posted entries | PR #40 automated PostgreSQL tests pass; no latest-main independent source/GL aggregate | N/E | CI fixture, not independent full payroll basis | Partial |
| Supplier invoice/payment → AP/GL | 100.00 invoice/payment ties to posted journals and ledger | Candidate API flow used `supplier_invoice` / `payment_made`; 2 balanced journals, raw GL 200.00/200.00, net TB 100.00/100.00, AP 0.00 | 0.00 through AP/GL for this fixture | Source invoice/payment plus read-only SQL over journal/GL rows | Fixture PASS; CFR-011 auto_supplier_* path not executed and remains failed on main |
| Report screen → CSV/PDF/XLSX | Same filter, period, totals and exact monetary values | Six report views loaded; Balance Sheet marked unreconciled and PDF export disabled; no file downloaded | N/E for screen-to-file parity | Automated report views; no export artifact | Screen view only; export acceptance blocked |

`N/E` means not executed. It is not zero and must not be interpreted as a passing difference.

## Synthetic oracle required for acceptance

For each future scenario, independently compute expected values from the synthetic business source records before invoking the app. Record currency and approved precision, policy assumptions, and source inputs. Then compare source totals against persisted payment/order state, event, journal, ledger, TB, P&L, BS, Cash Flow, AR/AP controls, and downloaded exports. For every amount record `expected`, `actual`, `difference`, source of expected value, policy approval, and pass/fail.

A suitable minimum neutral double-entry fixture (only after accounting policy selects recognition/tax treatment) should separately identify cash/card clearing, revenue, payable tax, discounts, inventory/COGS, payroll, AP/AR and reversals. Do not create a balancing plug or infer tax, cash-flow, inventory, retained-earnings, or refund policy to make an equation pass.

## Existing evidence boundary

- The historical tracker recorded direct PostgreSQL defect probes for unbalanced posted entries, mutation of posted lines, and supplier `entry_type` length; those are historical findings, not fresh latest-main results.
- PR #34 recorded extensive direct SQL and disposable PostgreSQL integration checks for journal/GL invariants. PR #35 added disposable PostgreSQL outbox/recovery coverage. PR #40 added payroll integrity tests. PR #41 added financial numbering concurrency tests. PR #42 added Balance Sheet sign/equation tests.
- Current PR checks passed for those merged changes. They do not substitute for an independently calculated cross-module statement and export reconciliation.
- No accounting-policy approval or accountant sign-off was found in the examined current evidence.

**Financial reconciliation verdict: FAIL / PARTIAL.** Source, payment, journals, net Trial Balance, P&L, AP and AR in the fixture differ by 0.00; Balance Sheet retains an unexplained -100.00 EGP residual. Cash Flow arithmetic is -100.00 EGP, but its classification requires accountant approval. No balancing entry was added.


## Phase 2 independent Balance Sheet oracle

The new latest-main CFR-010 fixtures use a disposable PostgreSQL database and independent SQL sums. Expected and actual values agree:

| Fixture | Expected | Actual | Difference | Verdict |
|---|---:|---:|---:|---|
| A opening posted journal: assets / liabilities / equity; TB debit/credit | 1,000 / 0 / 1,000; 1,000 / 1,000 | 1,000 / 0 / 1,000; 1,000 / 1,000 | 0 | PASS |
| B bank opening/current and GL cash; BS residual | 1,000 / 900; -100; -100 | 1,000 / 900; -100; -100 | 0 vs expected fixture; residual remains -100 | Reproduced setup discrepancy |
| C revenue / expense / net; post-expense assets/equity | 500 / 300 / 200; 1,200 / 1,000 | 500 / 300 / 200; 1,200 / 1,000 | 0 | PASS; unclosed residual +200 |
| D prior period net; after test close assets/equity; current net and BS residual | +200; 1,150 / 1,150; -20 / -20 | +200; 1,150 / 1,150; -20 / -20 | 0 | PASS; close/current earnings policy unresolved |

The previous -100 EGP report residual results from a bank opening amount stored in `bank_accounts_ledger` without a corresponding opening GL journal. A correctly posted opening journal balances independently. Unclosed earnings explain the separate +200 and -20 period residuals. No balancing plugs were entered and no production policy was implemented. Details and SQL method are in [Phase 2 execution evidence](evidence/phase2-execution-2026-11-14.md).

