# SERVIO Financial Acceptance — Phase 3B Current Snapshot

**Snapshot date:** 2026-11-15
**Audited baseline:** `origin/main` `3ed336c5f00d25596da3186176256a629d0d474c`
**Execution commit:** `d843ae6aca6cf4a7d7233233a7373753bd76e5d2`
**Pull request:** [#57](https://github.com/Yousef-Abdallah-1711/pos/pull/57), open and unmerged
**Financial CI:** [run 43730539837](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43730539837), success; evidence artifact `14167589858`
**Acceptance status:** 10 PASS, 6 PARTIAL, 1 NOT RUN, 0 FAIL; policy approvals remain separate and pending.
**Required verified CI:** PR run `43730539790` was canceled before application-image smoke; its required aggregate is therefore not green. Backend/frontend/security passed on earlier head `7d1422f`, before the CSV correction and latest-main merge. Latest-main run `43730274563` was canceled during image build before Chromium smoke.

## Runtime and data safety

The financial workflow ran only in GitHub Actions disposable infrastructure. The workflow created an ephemeral PostgreSQL database named `servio_test`, applied current migrations, initialized restricted `pos_app` and `pos_accounting` roles, and started isolated Redis. The real API used those restricted roles and `PAYMOB_MODE=mock`; test fixtures were newly seeded synthetic tenants, users, suppliers, invoices, orders, accounting journals, and salary records. The runner then destroyed its ephemeral services. The financial acceptance workflow did not access an external production or staging system. The separate verified-CI browser smoke did contact an unverified external API; see the safety finding below.

The API scenarios used the running Go application, authenticated synthetic owner/accountant/cashier sessions, persisted records, a real PostgreSQL backend, and independent SQL readbacks. The full report and scenario-level measurements are in [FINANCIAL_ACCEPTANCE_MATRIX.csv](FINANCIAL_ACCEPTANCE_MATRIX.csv).

## Executed scenarios and results

- Restaurant order, payment confirmation, accounting event worker, journal, GL, P&L, Trial Balance, Balance Sheet, AR/AP aging, Cash Flow, CSV download, and PDF download were exercised through the real API. Order source was 5500¢ total (5000¢ subtotal + 500¢ tax); cash, revenue, and tax GL differences were all 0¢. One event and one journal persisted; event status completed; foreign tenant read returned 404.
- Authenticated accountant journal: create returned 201; cashier post returned 403; accountant post returned 200; duplicate post returned 400. Persisted debit/credit were 10000¢/10000¢, two GL rows, no partial writes, and three validation failures were rejected.
- Payroll source gross/net 150000¢/150000¢ reconciled to expense and cash GL at 150000¢ each. Posting failure returned 500 and rolled back; retry and duplicate mark-paid returned 200; one journal and two GL rows remained. Paid edit/delete and foreign-tenant access were denied.
- Supplier invoice approval posting failure rolled back to draft with no journal; retry succeeded; repeated approval was rejected while retaining exactly one journal. Payment posting failure left no payment, journal, lines, or GL. Five concurrent requests of 5000¢ resulted in exactly four successes; the fifth, the repeated 5000¢ request, and an additional 2100¢ overpayment were rejected. Final source balance, AP net GL credit, and authenticated AP aging report were each 2000¢.
- The sales CSV downloaded through the authenticated API reconciled to 5500¢ gross and one order. A synthetic `=HYPERLINK(...)` product name was apostrophe-prefixed in the downloaded CSV; focused tests cover formula-like names in sales, profitability, and inventory exports.
- Balance Sheet API figures matched independent raw-GL calculations: assets 5500¢, liabilities 500¢, equity 0¢, and residual 5000¢ (`UNRECONCILED`). AR aging was 0¢ vs 0¢ for this order-only fixture; supplier AP aging reconciled 2000¢ vs 2000¢ source and GL. Cash Flow net was 5000¢, matching the SQL classification; classification policy remains open. A valid 17029-byte PDF downloaded, but its monetary content was not extracted.

## Engineering correction

`CreatePaymentMade` previously read `invoices_supplier.balance_due` without locking the invoice row. Competing requests could validate against the same old balance before subsequent writes. `resolveSupplierPaymentContextTx` now obtains `FOR UPDATE OF inv` before checking the balance, serializing decisions for that invoice. The live API concurrency scenario confirms four payments settle 20000¢ against a 22000¢ invoice while preserving the 2000¢ unpaid balance and matching AP GL. The regression test is in the focused PR #57 branch.

## CI scope and safety finding

The repository `production-ci.yml` application-image smoke is not isolated CI: `ops/aws/journal-api-acceptance.cjs` hard-codes `https://api.pos.rakmyat.online/api/v1`. In run 43728374754 it created two synthetic chart-of-account entries, posted a synthetic 1.00 EGP journal, and then successfully created an audited reversal. The journal effects were reversed; the script does not delete the two synthetic accounts. I canceled the run before its remaining external marketing and synthetic-trial steps. This environment was not verified as disposable, so those external API writes are not counted as acceptance evidence. No credentials were included in the report.

On earlier baseline `e35b3c2` run 43726424284, uploaded marketing results show the EN/AR marketing axe checks passed and Lighthouse accessibility was 1.00 for both locales. The run failed its English Lighthouse performance gate at 0.88 against 0.90; Arabic was 0.90. The reported 14 marketing-button contrast failures were not reproduced in that earlier baseline evidence, so no accessibility fix PR was opened.
## Open acceptance work and decisions

- CFR-011 remains a product/accounting contract gap: `auto_supplier_invoice` and `auto_supplier_payment` entry methods still lack an established production caller. This phase tested the supported `supplier_invoice` and `payment_made` API path and does not wire or certify those exact methods.
- A real worker process kill/restart followed by replay was not performed. The observed live worker completed the paid-order event; automated PostgreSQL retries/lease recovery do not replace a process restart scenario.
- The Balance Sheet exposes a 5000¢ earnings residual because the order fixture contains no opening journal and the current-period 5000¢ earnings are not closed into equity. No artificial balancing entry was introduced. The AR aging fixture had no receivable invoices; PDF monetary content and unauthorized download denial remain unmeasured. Cash Flow classification needs policy approval.
- Qualified accounting approval is still needed for current-period earnings, retained earnings, fiscal close/reopen, opening balances, rounding/precision, tax and discount basis, refunds/voids, inventory valuation, payroll deductions, and cash-flow classification.
- No authenticated visible English/Arabic desktop/mobile browser acceptance was run. The 375px horizontal overflow was not reproduced; responsive UI acceptance remains open.
- Independent PR review is outstanding. Do not treat the author’s checks as independent approval.

## Release decision

**Financial release: NO-GO.** The order-to-ledger chain, accountant journal workflow, payroll posting, and supported supplier invoice/payment API lifecycle now have real isolated API evidence with 0¢ differences for the executed reconciliation checks. Release remains blocked by pending independent review, exact CFR-011 contract, worker process restart/replay, Balance Sheet/AR/PDF content acceptance, accounting policy decisions, and authenticated responsive browser acceptance.
