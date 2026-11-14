# Financial Edge-Case Test Results
> **Phase 2 execution update (2026-11-14):** `origin/main` was fetched and verified at `165837c7ad2e506dbf20b1a45c9a123476f2f8e9`. Focused CFR-010 A–D PostgreSQL fixtures pass against a newly provisioned schema-only disposable database; independent SQL confirms Fixture A balanced opening (assets/equity 1,000/1,000), Fixture B off-ledger opening reproduces the historical -100 EGP residual, Fixture C expense and Fixture D two-period/close cases reconcile to their independent figures. Migrations `186_accounting_event_redrive_audit.sql` (PR #37), `186_servio_trial_offer.sql` (main), and `187_expand_journal_entry_type_to_32.sql` (PR #44) applied together in fresh and upgraded disposable schemas; long `auto_supplier_*` values persisted, and unsafe rollback refused. This does **not** execute authenticated supplier auto-posting: the only production references to those methods/constants are definitions, with tests as the only callers. Manual browser acceptance remains blocked; no credentialed synthetic browser session is available. PR #46 adds the regression and its CI is being tracked. Financial production remains **NO-GO**.

> **Execution update (2026-11-14):** A later audit pass ran a verified local disposable stack: 1 authenticated supplier invoice/payment API workflow, 6 automated API-backed report views, 1 unauthenticated authorization-negative request, and 0 manual financial browser workflows. Independent SQL reconciliation found source/payment/GL/TB/P&L/AP/AR differences of 0.00 for the narrow fixture, but a -100.00 EGP Balance Sheet residual because opening cash is absent from GL. Local isolation PASS for observed controls; existing staging isolation UNKNOWN. Financial release remains **NO-GO**. Detailed scope, values, commands, and limitations: [disposable execution evidence](evidence/disposable-execution-2026-11-14.md). This update supersedes earlier statements of no local fixture where applicable.

**Baseline:** `origin/main` `018682beb9448eea6dd72fd2418876b439a94dbf`.  
**Execution rule:** no heavy local test suite was run. Merged PR CI uses a disposable PostgreSQL and Redis service and runs tagged integration tests; this is automated evidence, not manual E2E. The mainline audit itself did not rerun a disposable DB suite. “Not executed” means no current-main observation was made during this audit.

## Journal and General Ledger

| Scenario | Evidence at audit | Result |
|---|---|---|
| Zero-line journal | Migration 180 / integration checks reject zero-line postings; CI green on #34. | Automated test evidence; not rerun here. |
| Unbalanced direct SQL posting | Historical disposable PostgreSQL probe reproduced acceptance before migration 180; #34 adds deferred invariants and CI. | Fix implemented; current runtime not independently rerun. |
| Negative / non-finite / sub-cent auto-posting amounts | CFR-001 validator rejects invalid values; #29 focused validator tests and CI passed. | Automated source/test evidence; cross-currency precision remains policy dependent. |
| Duplicate journal post and void | CFR-002 integration coverage serializes same-draft post and enforces inverse effects; #34 CI passed. | Automated PostgreSQL evidence. |
| Posted journal/line/GL mutation and duplicate GL row | Migration 180 triggers/constraints and direct SQL tests in #34. | Automated PostgreSQL evidence. |
| Cross-tenant account reference | Migration 180 scope guards and direct SQL integration coverage. | Automated PostgreSQL evidence; full API authorization matrix remains open. |
| Parallel document numbering | PR #41 adds advisory transaction lock and `cfr003_numbering_integration_test.go`. | CI passed; no audit-time concurrency stress. |
| Concurrent distinct postings / running balances | Row locks exist in journal/account paths. No current parallel-account balance oracle was run in this audit. | Partial. |
| Closed-period reversal, explicit linked reversal, retry after rollback | No approved behavior or full current integration fixture. | Not executed / policy blocked. |

## Orders and Payments

| Scenario | Evidence at audit | Result |
|---|---|---|
| Paid order persists accounting intent | PR #35 adds outbox/event recovery and PostgreSQL tests for retries, duplicate delivery, worker lease expiry, reconciliation, and paid-status+intent rollback. CI passed. | Automated integration evidence; not manually rerun. |
| Worker interruption/restart/replay | Worker retry/lease tests are in merged #35. | CI evidence only; no live worker restart in this audit. |
| Duplicate / concurrent event delivery | #35 integration coverage; #33 atomic journal+GL and #34 uniqueness/immutability. | CI evidence only. |
| Payment before order completion / out-of-order provider callback | No actual provider/account workflow exercised. | Not executed. |
| Pay-versus-cancel race | No full current provider/order race fixture in this audit. | Not executed. |
| Paid cancellation, partial/full refund, chargeback | PR #39 rejects unsupported transitions. It does not establish a posted refund/cancellation reversal lifecycle. | Unsupported/rejected transition guard exists; accounting behavior remains open/policy blocked. |
| Provider timeout or success followed by local persistence failure | No provider sandbox was used. | Not executed. |
| Duplicate inventory/COGS movement | #35 adds idempotent order sale inventory movement recovery; policy-level quantity/value reconciliation not run. | Automated scope evidence only. |

## Payroll, Invoices, Supplier, Inventory, and Assets

| Scenario | Evidence at audit | Result |
|---|---|---|
| Paid salary edit/delete | PR #40 adds paid-state guard and dedicated PostgreSQL integration coverage. | CI passed; not manually rerun. |
| Concurrent/repeated mark-paid and failed journal posting | PR #40 implements transactional/idempotent posting and failure injection coverage. | CI passed; reversal/correction policy remains absent. |
| Cross-restaurant salary request | Repository/handler scoping exists; no full synthetic role-by-restaurant attack matrix executed. | Not independently accepted. |
| Supplier invoice/payment posting | Disposable candidate stack applied migration 187 and confirmed width 32. One authenticated API invoice/payment workflow persisted 100.00 EGP using `supplier_invoice` and `payment_made`, tying source, payment, journals, GL, TB, P&L, AP/AR for this fixture. It did not execute CFR-011 `auto_supplier_invoice` / `auto_supplier_payment`. PR #44 full required CI passed; review absent. | **Main fails CFR-011; candidate schema verified, exact path not executed.** |
| Customer/supplier tax and discount math | Source reviewed; expected tax treatment is not approved and was not independently calculated. | Policy blocked; no result. |
| Partial payment, overpayment, duplicate invoice payment | No current PostgreSQL invoice lifecycle fixture executed in this audit. | Not executed. |
| Inventory receipts/sales/returns/waste and COGS | No approved cost method or latest-main independent quantity/value/GL fixture. | Policy blocked. |
| Fixed asset acquisition/depreciation/disposal | No current lifecycle fixture. | Policy blocked / not executed. |

## Reports, Exports, and Security

| Scenario | Evidence at audit | Result |
|---|---|---|
| Trial Balance vs GL | No independent full-ledger recomputation on latest main. | Not executed. |
| P&L and current-period earnings vs Balance Sheet | PR #42 changes sign contribution and adds equation integration tests; source still does not derive current-period net income into equity. | Partial; no independent multi-period oracle. |
| Cash Flow / opening-to-closing cash | No cash-method fixture. | Policy blocked. |
| AR/AP aging and balance tie-out | No invoice fixture. | Not executed. |
| CSV formula injection / PDF / XLSX parity / oversized export | No real export file was downloaded or adversarially parsed in this audit. | Blocked. |
| API RBAC, forged tenant/restaurant IDs, stale session | No authenticated test roles or verified isolated environment. | Blocked. |
| RLS restricted runtime role | Prior tracker recorded one focused restricted-role test; #34 CI adds ledger restrictions. This audit did not rerun a complete negative matrix. | Partial evidence only. |
| Arabic/English document semantics | Automated report views loaded with English root `lang=en`, `dir=ltr`, and Arabic root `lang=ar`, `dir=rtl`. | Automated runtime observation; manual review not executed. |

## Migration and failure testing

- CI workflow definition runs `go run ./cmd/migrate` and `go test -tags integration -count=1 -timeout 15m -json ./...` on a fresh disposable PostgreSQL service. This is the correct remote execution lane for the requested heavy suites.
- PR #34 recorded migration 180 preflight, direct SQL guards, downgrade refusal/behavior and failure injection on a disposable PostgreSQL database. These results are attached to that PR history, not a fresh audit rerun.
- PR #35 recorded recovery/outbox/worker integration tests against disposable PostgreSQL. PR #40 recorded payroll integration coverage. PR #41 recorded numbering integration tests. PR #42 recorded Balance Sheet integration tests.
- Candidate migration 187 in PR #44 was applied in the disposable local candidate stack; `entry_type` width 32 was verified. Migration 186 was not applied. Migration rollback was not exercised. The guarded rollback refuses to truncate values above 20 characters.
- No destructive, migration rollback, provider, tenant, or worker failure test was executed against the live staging service.

This report intentionally records no pass for unexecuted workflows and no accounting result without an independent oracle.



## Phase 2 executed edge cases

- **CFR-010 A balanced opening journal:** independent expected assets/equity 1,000/1,000; report reconciled; trial balance 1,000/1,000; opening-date P&L 0/0/0. PASS.
- **CFR-010 B opening cash outside GL:** bank opening/current 1,000/900; GL cash net -100; Balance Sheet residual -100. Historical discrepancy reproduced from fixture setup. No offset posted. Policy unresolved.
- **CFR-010 C expense:** revenue 500, expense 300, net income 200; after-expense asset/equity figures 1,200/1,000 with residual +200. PASS for the current unreconciled report behavior.
- **CFR-010 D two periods:** prior net +200; synthetic explicit close moves that amount to retained earnings (test-only) and brings equation residual to 0; next period net -20 and residual -20. PASS for the report's current close-dependent behavior; production period/earnings policy unresolved.
- Focused tagged PostgreSQL tests for the cases above passed twice on an isolated database restored from schema only. No heavy suite ran locally.
- Migrations main-186 + PR37-186 + PR44-187 were applied together in fresh and upgraded disposable schemas. Same numeric prefix 186 did not collide because `schema_migrations` records complete filenames. Supplier values persisted at 21 characters; PR44 downgrade refused while those values existed and left width at 32.
- Exact authenticated `auto_supplier_*` API workflow remains unexecuted because production source has no caller/route. Manual browser remains blocked. See [Phase 2 evidence](evidence/phase2-execution-2026-11-14.md).
