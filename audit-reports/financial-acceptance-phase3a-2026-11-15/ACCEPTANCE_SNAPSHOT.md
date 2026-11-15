# SERVIO Financial Acceptance — Phase 3A Snapshot

**Snapshot date:** 2026-11-15 (Africa/Cairo)  
**Repository:** `Yousef-Abdallah-1711/pos`  
**Latest-main baseline:** `0a3b104da33209c4b2cdcd0ec64998375e2c8cef`  
**PR head:** `4de0dad048f679cb669fa86dc68f08c98bfebc76`  
**Pull request:** [#57 — isolated financial acceptance runner](https://github.com/Yousef-Abdallah-1711/pos/pull/57)

## Current decision

**Financial release remains NO-GO.** A disposable PostgreSQL/Redis CI stack was actually provisioned, restricted roles initialized, latest-main migrations applied, and database-backed handler/use-case/repository integration tests executed. This is materially beyond the earlier “NOT RUN” snapshot, but it is not a full authenticated running-API order-to-ledger acceptance run, does not exercise positive supplier invoice/payment lifecycles, and does not establish browser acceptance or accountant-approved policy.

## Executed CI evidence

- Financial workflow: [run 43695500770](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43695500770), job `financial-postgres`, **success**. The workflow ran with the PR merge base at the exact latest-main SHA shown above.
- The job started disposable PostgreSQL and Redis containers, initialized the repository’s restricted `pos_migrator`, `pos_app`, and `pos_accounting` roles, and applied migrations successfully. CI stopped the containers at job completion. No production or staging database was used.
- The test command ran integration-tagged packages `internal/handler/http`, `internal/usecase`, `internal/repository`, and `cmd/api` against CI PostgreSQL. Raw Go JSON events were uploaded as artifact `financial-acceptance-76d91630a9483c757cf3aa345114e36a8fea1fe5` ([artifact](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43695500770/artifacts/14154078867)).
- Artifact counts: **422 top-level tests passed, 2 top-level tests skipped, 0 failed**. Across nested test/subtest events, 956 passed, 2 skipped, 0 failed. The skipped benchmark is not acceptance evidence.
- The separate `SERVIO verified CI` run [43695500767](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43695500767) passed backend and its `go test -tags integration ./...` step; the overall run was still in progress at snapshot time. This broader run is not counted as a second financial-acceptance run.
- PR #57 is open and awaits independent review. No merge, deployment, production write, or staging write was made.

## What the persisted tests established

- Durable paid-order intent, worker retry/replay, expired-lease recovery, idempotent posting, tenant scope, and exact journal/GL row counts passed in PostgreSQL-backed integration tests.
- Payroll mark-paid persistence, duplicate/concurrent payment prevention, tenant scope, and rollback-at-failure-boundary tests passed.
- Journal/GL integrity and restricted-role protections, cross-tenant accounting denial, and financial report/earnings boundary assertions passed.
- PostgreSQL-backed auth/session and handler integration tests passed, but they construct handlers and middleware in test processes. They did **not** log into a separately running full API binary. Therefore authenticated application-process acceptance remains **NOT RUN**.
- The supported supplier invoice/payment positive lifecycle was not exercised by this runner. Existing negative-scope checks and mocked auto-posting examples do not prove live AP reconciliation. `PostSupplierInvoiceEntry` / `PostSupplierPaymentEntry` remain outside the active production contract; do not invent a caller. CFR-011 support versus supersession is still an owner decision.

The dated scenario matrix separates PASS, PARTIAL, and NOT RUN and records exact test names: [Phase 3A financial acceptance matrix](FINANCIAL_ACCEPTANCE_MATRIX.csv).

## Infrastructure defect found and corrected

The first run, [43695391211](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43695391211), failed before checkout because the existing PostgreSQL service image digest could not be resolved by Docker Hub. The financial workflow and existing `production-ci.yml` now pin verified official PostgreSQL 16.4 and Redis 7.4.2 manifest digests. The successful run confirms service initialization and migration startup on the updated definitions.

## Remaining execution and policy gates

1. Add a full API-process integration scenario with synthetic authenticated owner/accountant/cashier users and persisted source fixtures; query independent SQL totals for the complete order → payment → event → worker → journal → GL chain.
2. Execute the positive supported supplier invoice approval and payment routes with AP/GL reconciliation, atomicity, and idempotency. Keep exact CFR-011 `auto_supplier_*` behavior pending the owner’s support/accounting decision.
3. Execute authenticated accountant valid/invalid journal workflows and inspect real financial report CSV/PDF/XLSX downloads; do not infer download results from handler response-shape tests.
4. Run authenticated Playwright English/Arabic desktop and 375px mobile scenarios. Do not call automated checks manual acceptance.
5. Obtain qualified accounting/product decisions for opening-balance offsets/dates, fiscal year and cutoff, current-period earnings, close/reopen and retained earnings, distributions, cash-flow classification, and correction/reversal lifecycle. Do not add balancing plugs.

**Release assessment:** NO-GO pending the remaining workflow execution, exact supported supplier AP reconciliation, independent review, browser acceptance, and policy approvals.
