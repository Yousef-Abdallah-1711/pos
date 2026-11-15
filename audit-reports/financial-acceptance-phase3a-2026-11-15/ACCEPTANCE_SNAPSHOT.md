# SERVIO Financial Acceptance — Phase 3A Snapshot

**Snapshot date:** 2026-11-15 (Africa/Cairo)  
**Repository:** `Yousef-Abdallah-1711/pos`  
**Latest-main baseline:** `5c00e53434fe96b5ec5e2f73d6d3d853fbfa7bd9`
**Acceptance source commit tested against latest main:** `ab9d09e978fdd6081523d863473f3b3727e1f54b`
**Pull request:** [#57 — isolated financial acceptance runner](https://github.com/Yousef-Abdallah-1711/pos/pull/57)

## Current decision

**Financial release remains NO-GO.** A disposable PostgreSQL/Redis CI stack was provisioned, restricted roles initialized, latest-main migrations applied, and the real API process started with restricted database identities. Authenticated order-to-ledger and supported supplier invoice-to-payment lifecycles passed with independent SQL reconciliation. Supplier retry/rollback guarantees, accountant/cashier API-process workflows, financial report downloads, browser acceptance, and accountant-approved policy remain open.

## Executed CI evidence

- Latest-main financial workflow: [run 43701118387](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43701118387), job `financial-postgres`, **success** on acceptance head `ab9d09e978fdd6081523d863473f3b3727e1f54b`, with base `5c00e53434fe96b5ec5e2f73d6d3d853fbfa7bd9`.
- CI started disposable, loopback-bound PostgreSQL and password-protected Redis containers, initialized `pos_migrator`, `pos_app`, and `pos_accounting`, applied latest-main migrations, built and started the actual API, and removed containers after test completion. The API authenticated a synthetic owner. No production or staging system was used.
- The live test repeated the authenticated `paid` confirmation (HTTP 200) and persisted exactly one accounting event (`paid-order:2:2:1`, event count 1), which the worker completed. It found one posted `auto_sales` journal (ID 1) with three GL rows. Independent source totals were subtotal **5,000¢**, tax **500¢**, discount **0¢**, delivery **0¢**, gross **5,500¢**. Persisted order total was **5,500¢** (difference **0¢**); GL cash was **5,500¢** (difference **0¢**), sales revenue **5,000¢** (difference **0¢**), and tax payable **500¢** (difference **0¢**). The synthetic product was non-inventory with zero cost, so this scenario did not post or assert a COGS journal.
- The foreign synthetic tenant authenticated and received HTTP **404** reading the first tenant’s order. Both live API acceptance scenarios passed. The broader PostgreSQL-backed handler/use-case/repository/route suite recorded **422 top-level passes, 3 skips, 0 failures**; skips are not counted as passed acceptance scenarios.
- The same isolated API process authenticated an owner who created a synthetic supplier, created and approved an invoice, and paid it through the supported payment route. Independent source totals were subtotal **20,000¢**, input VAT **2,000¢**, and total **22,000¢**. The invoice journal credited AP **22,000¢**; the payment journal debited AP **22,000¢** and credited cash **22,000¢**. The invoice persisted as paid with **0¢** due. All checked source-to-GL differences were **0¢**. The test used one payment submission; it did not establish retry idempotency or failure rollback for this lifecycle.
- Supporting artifact: [latest-main financial acceptance JSON evidence](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43701118387/artifacts/14156416799). It contains synthetic tenant/source identifiers and monetary values only.
- Full verified CI [run 43701118404](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43701118404) passed backend, frontend, security/config and image builds, but failed its generic landing-page Chromium accessibility gate on serious CTA color-contrast findings. Lighthouse performance scored 0.92 (EN) / 0.95 (AR), accessibility 0.97 in both locales, and CLS 0. The aggregate required check therefore failed. This is separate from the financial API workflow.
- CI history is preserved: run [43696982964](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43696982964) appeared green because its `go test | tee` pipeline lacked `pipefail`; its uploaded artifact actually recorded a worker retry. Strict pipeline handling exposed the fixture’s missing standard COGS accounts; the next run exposed an incorrect assertion comparing gross order total to revenue. These were corrected and run 43701118387 passes on the latest-main base. No product-side accounting behavior was changed.
- PR #57 remains open. The only repository collaborator available to request as reviewer is the PR author, so independent review is still pending; no self-review, merge, deployment, production write, or staging write was made.

## What the persisted tests established

- The real running API authenticated a synthetic owner, created and paid an order, accepted an identical payment confirmation idempotently, and persisted exactly one paid-order event. The worker completed that event and posted one sales journal; SQL independently reconciled order total and cash/revenue/tax GL with zero-cent differences.
- The real running API authenticated a synthetic owner and completed supplier creation, invoice approval, and supplier payment. SQL verified one posted invoice journal, one posted payment journal, paid invoice state, zero balance due, and matching expense, input VAT, AP, and cash GL amounts.
- Durable paid-order intent, worker retry/replay, expired-lease recovery, idempotent posting, tenant scope, and exact journal/GL row counts also passed in PostgreSQL-backed integration tests.
- Payroll mark-paid persistence, duplicate/concurrent payment prevention, tenant scope, and rollback-at-failure-boundary tests passed.
- Journal/GL integrity and restricted-role protections, cross-tenant accounting denial, and financial report/earnings boundary assertions passed.
- PostgreSQL-backed auth/session and handler integration tests passed, but those tests construct handlers and middleware in test processes. The only authenticated workflows exercised through a separately running API binary are owner order-to-ledger and owner supplier invoice-to-payment. Cashier/accountant, payroll, report, and export workflows through that API remain **NOT RUN**.
- The supported positive supplier lifecycle is now exercised. Supplier retry idempotency, rollback-on-failure, cross-tenant AP denial, and AP subledger-to-GL checks under replay remain unverified. `PostSupplierInvoiceEntry` / `PostSupplierPaymentEntry` remain outside the active production contract; do not invent a caller. CFR-011 support versus supersession is still an owner decision.

The dated scenario matrix records PASS and NOT RUN statuses with exact test names: [Phase 3A financial acceptance matrix](FINANCIAL_ACCEPTANCE_MATRIX.csv).

## Infrastructure defect found and corrected

The first run, [43695391211](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43695391211), failed before checkout because the existing PostgreSQL service image digest could not be resolved by Docker Hub. The financial workflow and existing `production-ci.yml` now pin verified official PostgreSQL 16.4 and Redis 7.4.2 manifest digests. The successful run confirms service initialization and migration startup on the updated definitions.

## Remaining execution and policy gates

1. Extend the real API-process scenarios to cashier and accountant roles and test the complete inventory/COGS branch, including worker interruption/restart and replay through the process boundary.
2. Add supplier route retry/idempotency, rollback-on-failure, and cross-tenant checks. Keep exact CFR-011 `auto_supplier_*` behavior pending the owner’s support/accounting decision.
3. Execute authenticated accountant valid/invalid journal workflows and inspect real financial report CSV/PDF/XLSX downloads; do not infer download results from handler response-shape tests.
4. Run authenticated Playwright English/Arabic desktop and 375px mobile scenarios. Do not call automated checks manual acceptance.
5. Obtain qualified accounting/product decisions for opening-balance offsets/dates, fiscal year and cutoff, current-period earnings, close/reopen and retained earnings, distributions, cash-flow classification, and correction/reversal lifecycle. Do not add balancing plugs.

**Release assessment:** NO-GO pending supplier retry/rollback checks, the remaining role/report/export workflows, independent review, browser acceptance, and policy approvals.
