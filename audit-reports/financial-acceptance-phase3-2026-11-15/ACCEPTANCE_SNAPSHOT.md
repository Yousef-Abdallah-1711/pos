# SERVIO Financial Integration Acceptance — Phase 3 Snapshot

> Historical cutoff retained as recorded. The later Phase 3A CI execution and updated scenario statuses are in [the current Phase 3A snapshot](../financial-acceptance-phase3a-2026-11-15/ACCEPTANCE_SNAPSHOT.md) and [matrix](../financial-acceptance-phase3a-2026-11-15/FINANCIAL_ACCEPTANCE_MATRIX.csv). The statements below describe the earlier Phase 3 cutoff and are not the latest status.

**Snapshot date:** 2026-11-15 (Africa/Cairo)

**Repository:** `Yousef-Abdallah-1711/pos`

**Authoritative baseline:** fetched `origin/main` = `3ed336c5f00d25596da3186176256a629d0d474c`

**Current work:** safety PR #63 (`296390d431d1313d47efdddcc2a1663d9c609953`), worker test PR #64 (`779b45c2f4f2e5b0bda99a1fbfb979e856aed7ca`, rebased and stacked on #63), and AR correction PR #65 (`a47a0e64872eafe4341ff103d76cc06121093b92`, stacked on #63) remain open and unmerged.

## Decision

**Financial acceptance: NO-GO.** No verified staging access or visible manual synthetic session was available. Automated isolated GitHub CI did execute the CI safety browser workflow and, separately, a PostgreSQL-backed worker process restart/replay test. Those results are scoped to the open PR heads and do not establish the complete authenticated order-to-ledger, supplier, payroll, reporting, export, or policy acceptance required for release.

## Current source and delivery facts

- Fetched `origin/main`; latest observed SHA is `3ed336c5f00d25596da3186176256a629d0d474c`. This snapshot and its PR evidence do not modify main.
- PR #57 remains open at head `a77818ee343b6795d3bf937b64c2e3aa958c1597`, base `3ed336c5f00d25596da3186176256a629d0d474c`, and merge state `BLOCKED`; its `financial-postgres` check passed in run [43732923459](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43732923459), while the required aggregate run was canceled/failed. That PR's evidence is not treated as current-main acceptance. PR #63 is the CI target lockdown. PR #64 adds a process-boundary accounting worker recovery test. Neither PR has an independent review request or approval; the available collaborator is the PR author, so self-review/approval was not used.
- PR #63 full required CI [run 43747265459](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43747265459) passed on head `e36cdb4dc1979d9c77c172f56c82023b6b55e2f8`. The refreshed safety report at head `296390d431d1313d47efdddcc2a1663d9c609953` also passed full required CI [run 43750104349](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43750104349), including the Chromium smoke and required aggregate. The snapshot/matrix/register commit being prepared here will trigger a new exact-head run.
- PR #64 prior head full required CI [run 43749067506](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43749067506) passed. Refined head `f85dcae4914c4ef09b4c5f90f2351bd3fb037e74` passed backend CI, but full required run [43750826060](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43750826060) failed the unchanged Lighthouse performance/accessibility/CLS threshold in application-images; Chromium scenario checks reported 0 failures and 0 browser errors. Thresholds were not changed. PR #64 was rebased onto `codex/ci-safety-lockdown`; exact head `779b45c2f4f2e5b0bda99a1fbfb979e856aed7ca` passed fresh full required CI [run 43753533112](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43753533112), including application-images and required aggregate; URL smoke [43753533101](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43753533101) and runtime guardrail [43753533130](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43753533130) passed. Its exact-head artifact records PID 3940 stopped after commit, PID 3946 completing attempt 2, one sales journal, one COGS journal, four GL rows, one -2 inventory movement, and zero source-to-GL differences for 10,000-cent revenue and 1,000-cent COGS. Its synthetic worker test was performed against disposable CI PostgreSQL, not staging.
- PR #65 fixes invoice void SQL parameter typing, uses `bigint` for optional order references, and adjusts scoped customer balances atomically on invoice create, edit, customer reassignment, and void. Its repository lifecycle test reconciles customer balance, non-void invoice balances, and receivable GL at 60.00 after a $100 invoice, $40 payment, and a separate $50-to-$70 invoice edit followed by void (difference 0.00). A real PostgreSQL aging report test places an overdue $85 invoice in the 60-day bucket and excludes a void $700 invoice and a foreign-tenant $900 invoice; independent SQL and report total both equal $85. Full required CI, including backend integration, application-images, and the required aggregate, passed on exact head `a47a0e64872eafe4341ff103d76cc06121093b92` in run [43754007733](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43754007733). These tests exercise repository/usecase paths, not authenticated API calls or role denial. Migration 187 is on main and changes `journal_entries.entry_type` to `VARCHAR(32)`. The former 21-character versus `VARCHAR(20)` schema mismatch is no longer an outstanding mainline schema defect. CFR-011 exact `auto_supplier_*` posting paths remain undecided and untested through an authorized production application route.

## CFR-011 disposition

The two methods are historical/planned auto-posting functionality with no demonstrated production integration, and whether they should remain supported requires an explicit product/accounting decision. The evidence does not justify wiring them into a new route or substituting their journal semantics into a live workflow.

- `PostSupplierInvoiceEntry` and `PostSupplierPaymentEntry` are defined only in `backend/internal/usecase/auto_posting_engine.go`; repository-wide production-source search finds no caller.
- The `AutoPostingService` production interface contains sales, COGS, and payroll posting only. The two supplier methods are not part of the application service contract.
- Historical accounting specifications and seeded rules document supplier auto-posting as intended/planned functionality. Current production routes instead call the supplier invoice and payment repositories. The invoice approval path transactionally creates/posts a journal and updates invoice status using `supplier_invoice`; payment creation uses its own live path with `payment_made`. These different paths do not prove the exact auto-posting methods obsolete, nor establish that their inventory/input-VAT account mapping is correct for every supplier invoice.
- **Disposition:** keep the schema-width fix; do not claim CFR-011 exact-path acceptance. Product/accounting owner must decide whether the auto supplier entry types are supported and define their source lifecycle/account mapping, or confirm they are superseded so a separate removal change can be reviewed. No accounting policy or public API was invented.

## Execution environment and evidence boundary

- No production or unverified staging environment was accessed, and no production/staging writes were made. Staging database, Redis, provider credentials, and egress isolation remain unverified.
- PR #63's GitHub Actions browser stack verified the disposable database, internal test network, synthetic accounts, and mocked payment path for that run. This was automated CI; it was not a visible manual acceptance session.
- PR #64 backend evidence records an actual worker subprocess stop after committed financial postings, a restarted PID after lease expiry, and persisted SQL reconciliation: revenue 10,000 cents vs GL 10,000 (difference 0); COGS 1,000 vs GL 1,000 (difference 0); one movement of -2 units. One sales journal, one COGS journal, four GL rows, event attempt 2 completed.
- Automated exact-head browser evidence from PR #65 run [43754007733](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43754007733) passed 58/58 checks; the inspected artifact [14178231002](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43754007733/artifacts/14178231002) records all 13 375px pages at `scrollWidth=375`, synthetic $27.50 cash checkout, and no browser errors. English/Arabic Lighthouse performance was 0.96/0.97, accessibility 1.00/1.00, and CLS 0/0. The previously reported overflow was not reproduced. Accountant financial report journeys, PDF monetary extraction, and authenticated Arabic reports were not covered.
- No manual visible-session acceptance is claimed. Heavy tests were run in GitHub Actions rather than locally. No real customer or payment-provider data was used.

## Financial policy boundary

Preserve the evidence that distinguishes an omitted opening journal from unclosed earnings. A bank opening value held outside posted GL can create a Balance Sheet residual; unclosed current/prior earnings can separately explain a period residual. Do not post an artificial balancing entry. The prior synthetic multi-period test close is test-only evidence, not approved production policy.

Qualified accounting/product approval is still required for: opening balance offsets and dates; fiscal year and report cutoff; current-period earnings roll-forward; close/reopen and retained-earnings transfer; distribution treatment; cash-flow classification; and the period-correction/reversal lifecycle. Policy-dependent changes remain pending.

The decision owner, existing behavior, financial impact, options, engineering proposals, approvers, and blocked acceptance tests are listed in the [Accounting Policy Decision Register](ACCOUNTING_POLICY_DECISION_REGISTER.md). Proposals in that register are not approved policy.

## Phase 3 acceptance matrix

The attached [17-row matrix](FINANCIAL_ACCEPTANCE_MATRIX.csv) separates current PR-CI evidence from unexecuted API/report workflows and historical Phase 2 evidence. Worker interruption and duplicate replay passed at PR #64's backend test level; event creation and aggregate journal checks are partial because the order/payment flow was not initiated through the authenticated API. No previous report or source inspection is upgraded to current runtime acceptance.

Historical Phase 2 values, including the narrow 2026-11-14 supplier invoice/payment fixture and its `-100.00 EGP` Balance Sheet residual, remain in [the historical disposable execution report](../financial-remediation-master-audit-2026-11-14/evidence/disposable-execution-2026-11-14.md). They are not current Phase 3 execution results.

## Remaining release gates

1. Complete and inspect the exact-head full required CI runs for the latest #63 and #64 commits; obtain independent code review (no independent repository collaborator is currently available).
2. Add isolated authenticated API acceptance for invoice creation, positive AR aging, partial payment, accountant journal, payroll, cross-tenant denial, reports, and actual CSV/PDF downloads with independent SQL expected/actual/difference values.
3. Resolve CFR-011 support/supersession before adding a caller or removing methods; test the approved route after the product/accounting decision.
4. Obtain qualified accounting approval for opening balance offsets/dates, period earnings and retained earnings, fiscal close/reopen, tax/discount and rounding, refunds/reversals/voids, inventory valuation, payroll deductions, and Cash Flow classification.
5. Run authenticated EN/AR financial report browser acceptance at desktop and 375px, including lang/dir, integrity warnings, downloads, console/hydration errors, and overflow. Keep manual acceptance distinct; no authorized visible session is currently available.
6. Reconcile historical external CI writes using retained read-only evidence and owner authorization. Do not access or modify the external environment automatically.

**No PR was merged, no deployment was made, and no production or unverified staging write was performed. Financial release remains NO-GO.**
