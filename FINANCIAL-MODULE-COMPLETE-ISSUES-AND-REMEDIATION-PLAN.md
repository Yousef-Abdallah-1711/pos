# SERVIO Financial & Accounting Module — Complete Issues and Remediation Plan

**Purpose:** single working record for financial-module remediation. This is an evidence-led inventory and plan, not a claim of financial correctness or production readiness.

**Assessment date:** 2026-11-12 (Africa/Cairo)<br>
**Reference baseline:** fetched `origin/main` / `FETCH_HEAD` = `a6a2c6700a5fe883256deb1e5a3893108bd9865c`<br>
**Current worktree:** `H:\Projects\pos`, branch `codex/phase2a-continuation`, HEAD `ecadb5ab7665df3eb9f69886a9b186c710eff62d`<br>
**Merge base:** `5ab6dba62287ddd404ce00529f5fe99220ccdc19`; divergence: main +178 commits, branch +3 commits from that base.<br>
**Scope and safety:** no application-source edits, protected/existing database writes, branch switching, reset/rebase/merge/cherry-pick, commit, push, PR, or deployment was performed. A no-volume disposable PostgreSQL 16 environment was created and all **109 `.sql` migration files** present in this worktree were applied successfully through the repository migration runner. The directory also contains 16 `.disabled` files, which are not SQL migrations applied by the runner. Synthetic SQL probes and an RLS integration test used that disposable environment. The database/container later became unavailable when the Docker daemon stopped; it was not restarted to avoid affecting unrelated services. Only this existing Markdown tracker was edited in this continuation; existing source changes and untracked files were preserved.

## 1. Executive result

**Financial release verdict: NO-GO pending resolution of the P1 ledger, posting, report-equation, and payroll lifecycle risks.** Current-worktree testing now includes a full Go suite, static analysis/build, all 109 active `.sql` migration files in disposable PostgreSQL 16, schema probes, a focused RLS integration test, dashboard typecheck/build, financial report UI tests, and browser checks with mocked API/auth. Three database defects were reproduced. Multi-step accounting fixtures, concurrency/failure injection, comprehensive authorization matrices, and real backend-authenticated browser workflows remain incomplete; see §6.

Mainline already has meaningful financial breadth: chart of accounts, journals, general ledger/trial balance, order auto-posting, invoices, payment accounting, payroll posting, bank reconciliation, inventory/COGS paths, fixed-asset/tax-period/fiscal-period schemas, and financial reports. Their presence does not establish coherent lifecycle behavior, end-to-end reconciliation, or legal/accounting compliance.

Target direction: posted journal and ledger effects must be the accounting source of truth; operational orders, payment-provider state, invoice/payroll/inventory states, and bank/marketplace settlements must reconcile to that source. Business state plus accounting event creation must be atomic locally or recoverable through a durable outbox/inbox. Posted entries must be immutable and corrected through linked reversals. No tax, rounding tolerance, inventory method, or recognition policy is approved by this assessment.

### Count and current disposition

| Severity | Issue count | Current issues |
|---|---:|---|
| P0 | 0 | None established as an immediate confirmed catastrophic issue in this limited run. |
| P1 | 7 | CFR-001, 002, 003, 004, 005, 009, 010 |
| P2 | 9 | CFR-006, 007, 008, 011, 012, 013, 014, 015; FIN-018 |
| P3 | 5 | FIN-016, FIN-017, FIN-019, FIN-020, FIN-021 |
| **Total** | **21** | 15 canonical CFR items plus 6 additional findings. |

**Current-run test tally:** full Go suite: 726 passed, 0 failed, 2 skipped; focused opt-in RLS repository integration test: 1 passed; dashboard financial/report Vitest selection: 8 passed; direct SQL probes: 2 expected-behavior checks passed and 3 defect reproductions failed as intended. Thus **737 automated/test checks passed and 3 defect probes failed**, plus 2 Go tests skipped; this roll-up counts the two passing SQL behavior probes but excludes browser assertions because they were recorded as workflow observations, not a formal assertion-counted test suite. Browser observations confirmed Arabic report subtree RTL, no JS/page errors, mobile no-horizontal-overflow, and PDF download under mocked API/auth. Full accounting E2E and rollback/concurrency coverage remain incomplete. These are current-run results; historical evidence is identified separately.

## 2. Git, evidence, and worktree boundaries

- `git fetch origin main` completed; fetched SHA recorded above. No branch was switched.
- `git diff --name-only FETCH_HEAD...HEAD` shows eight committed differences: `backend/internal/repository/{contract_repo.go,driver_repo_integration_test.go,hr_child_resource_scope_integration_test.go,insurance_repo.go,resignation_repo.go,training_repo.go}`, `docker-compose.phase2a-e2e.yml`, and `scripts/phase2a-e2e/seed.ps1`. These are not financial-module implementation changes.
- Prior verified local-only candidate refs remain separate from this checkout: Phase 00 candidate `5064daa89a27715c3eeacb8becc9f3f48de1853a` on `fix/financial-phase00-safety` (warning/export guard and posting diagnostics; not on main and not a root-cause financial fix). No ref was switched or modified here.
- The worktree was already dirty before this assessment: 244 changed/untracked paths at inspection. Financial-specific local-only artifacts include `backend/internal/handler/http/inventory_adjustment_handler.go` and its integration test, `audit-reports/financial-forensic-audit-2026-11-10/`, and numerous unrelated Phase 2A changes/artifacts. No pre-existing file was edited by this assessment.
- Current mainline findings below are based on `git show FETCH_HEAD:<path>` / source lines at fetched SHA, corroborated by repository audit evidence where stated. The untracked financial forensic audit folder contains UI snapshots/screenshots dated 2026-11-10; it is historical evidence, not a newly rerun browser test.
- Repository instruction requested `C:\Users\OPT\Desktop\POS\CODING_STANDARDS.md`, which did not exist. The repository copy `H:\Projects\pos\CODING_STANDARDS.md` was read in full before assessment; it is the applied review standard.
- Deployment state was not independently checked. No live or protected database was queried or changed.

## 3. Current financial feature and workflow inventory

Status meanings: **Implemented** = code/schema/API is present; **Partial** = important behavior/evidence gaps remain; **Broken** = source-confirmed unsafe behavior; **Unverified** = cannot assert operational correctness from current-run evidence. “Present” is not a release sign-off.

| Capability / workflow | Mainline status | Evidence / boundary |
|---|---|---|
| Chart of Accounts | Implemented, policy incomplete | Migrations `136_create_chart_of_accounts.sql`, `156_seed_chart_of_accounts.sql`; account repository and accounting UI. Scope, lifecycle, opening balances and owner-configurable structure need stronger contract. |
| Manual journal draft/edit/post/void | Partial / unsafe edges | `journal_entry_repo.go`; `137_create_journal_entries.sql`. Application validation exists; direct DB invariants and concurrency/closed-period reversal behavior remain inadequate (CFR-002/003/013). |
| General Ledger and Trial Balance | Implemented, not DB-verified | `general_ledger_repo.go`, migration `154_create_general_ledger.sql`; reports use `float64` and stored running balances. Independent recomputation not executed. |
| Sales order creation and payment status | Partial | `order_usecase.go`, `admin_order_handler.go`; paid transition invokes auto-posting as best effort after the operational update. No verified universal order-channel convergence. |
| Order sales, tax, delivery and COGS posting | Broken/partial | `auto_posting_engine.go`; auto-posting can rebalance arbitrary differences and write journal then ledger in separate repository transactions. Cancellation/refund reversal not found. |
| Card/web payment provider lifecycle | Partial | Payment usecase/webhook unit tests pass; PostgreSQL integration tests are present but were not run. Provider authorization/capture/settlement to restaurant books is not demonstrated by subscription-payment webhook tests. |
| Customer invoices / AR / receipts | Implemented, unverified | `141_create_invoices_customer.sql`, `invoice_customer_repo.go`, `payment_acct_repo.go`; DB triggers compute totals, but discount/tax calculation order is suspect and DB fixtures not run. |
| Supplier invoices / AP / supplier payments | Partial / broken posting path | `142_create_invoices_supplier.sql`, `invoice_supplier_repo.go`, payment accounting repository. Automated entry type can exceed schema width; posting and reconciliation were not run against PostgreSQL. |
| Inventory movements and COGS | Partial / policy-dependent | `144_create_inventory_movements.sql`, `inventory_movement_repo.go`, auto-posting COGS; costing policy, movement-to-ledger equivalence and returns/wastage handling not verified. Local-only inventory adjustment route is not mainline. |
| Payroll and payroll journal | Broken lifecycle controls | Salary/payroll schema and `hr_salary_handler.go`; source permits edits/deletes without paid-state protection (CFR-009). No live isolated DB test. |
| VAT/tax periods | Implemented, compliance unverified | `146_create_tax_periods.sql`, VAT/report endpoints; no jurisdiction decision, approved rates or accountant validation. |
| Bank accounts/reconciliation | Implemented, unverified | `140_create_bank_accounts.sql`, `bank_reconciliation_repo.go`/accounting APIs; no imported-statement reconciliation fixture or independent reconciliation test run. |
| Balance Sheet | Broken equation presentation | `financial_report_usecase.go:261-336` reads account balances only for asset/liability/equity accounts; no derived current-period earnings line. It can report `EquationBalanced=false`. Current main lacks Phase 00 warning/export guard. |
| Profit & Loss | Implemented, unverified | Posted revenue/expense report path in `financial_report_usecase.go`; no independent fixture/reconciliation to ledger run this session. |
| Cash Flow | Partial / classification heuristic | `financial_report_usecase.go` classifies by account type/code pattern, not confirmed cash-flow policy or cash-account movement. |
| Receivables/payables aging | Implemented, unverified | `accounting_advanced_repo.go`; aging buckets use invoice due dates; fixtures not run. |
| Fixed assets / depreciation | Schema and paths present, lifecycle unverified | `145_create_fixed_assets.sql`, auto-posting/repository paths. Acquisition, depreciation schedule, disposal and GL tie-out not verified. |
| Gift cards / customer deposits / deferred revenue | Missing or not evidenced as complete | No complete end-to-end financial lifecycle established in examined source/catalog. See FIN-016. |
| Fiscal close / retained earnings roll-forward | Partial / missing workflow | Fiscal-period checks exist; a complete approved close, income-to-equity roll-forward, reopen controls and reports were not established. |
| CSV / PDF / XLSX exports | Partial | CSV/PDF report routes/UI evidence exist. Native XLSX parity, exact values, formula injection and limits were not established. CSV is not native XLSX. |
| Financial RBAC and tenant isolation | Partial / unverified | Main contains tenant-aware repositories and some RLS work; no complete negative matrix was executed for finance routes, jobs, exports, bank/close actions. |
| Consolidated tenant reporting / FX | Unverified | Cross-restaurant consolidation and multi-currency conversion policy and parity are not demonstrated. |

### Channel convergence and end-to-end accounting chain

Expected chain is **business event → operational persistence → payment state → accounting event → journal → GL → report → reconciliation**. Mainline evidence establishes only parts: the admin order payment-status path calls `AutoPostPaidOrderByID` after updating payment state; channel parity for cashier POS, dashboard order, storefront, theme-builder site, public API, kitchen, and admin-created order was not demonstrated. Kitchen state changes should not themselves recognize revenue. The accounting trigger should be a policy-defined payment/capture/fulfilment event, with explicit treatment of deposits, tips and tax. Each channel needs one canonical event identity and idempotent posting path.

## 4. Issue register — canonical CFR-001 through CFR-015

Statuses mean: **SOURCE-CONFIRMED** = code/schema establishes the defect or gap; **REPRODUCED** = executed test demonstrated it; **IMPLEMENTED / PARTIAL** = a correction is present but one or more acceptance gates remain unverified; **VERIFIED** = independently verified behavior; **SUSPECTED** = evidence incomplete; **BLOCKED** = reproduction needs unavailable safe fixture/runtime. No issue is marked VERIFIED FIXED by this session.

### CFR-001 — Auto-posting silently forces unbalanced lines to balance

- **Severity/status:** P1 — IMPLEMENTED / PARTIAL. The safe rejection path and regression tests are present on `fix/financial-cfr001-reject-unbalanced-posting`; Go verification passes on the latest fetched main, while PostgreSQL integration remains unverified. Do not mark VERIFIED FIXED yet.
- **Problem / impact:** A material tax, discount, quantity or mapping error can be concealed by changing the final accounting line; this can misstate revenue, VAT, expense or cash while producing a balanced journal.
- **Current / expected:** Before this correction, `AutoPostingEngine.createAndPostEntry` passed lines through `normalizeBalancedLines`, which changed the last line to force equal totals. It now validates copied lines before journal creation, rejects any debit/credit total mismatch, and passes unchanged monetary amounts to the repository. It adds no tolerance and no balancing account.
- **Implementation:** `backend/internal/usecase/auto_posting_validation.go` validates finite, non-negative amounts, one-sided nonzero journal lines, the existing two-decimal `NUMERIC(15,2)` storage range, and exact equality using integer cents. Sub-cent values are rejected rather than rounded. `createAndPostEntry` runs this validation before journal or GL repository calls, and the payroll path no longer normalizes independently.
- **Root cause:** `backend/internal/usecase/auto_posting_engine.go` previously called `normalizeBalancedLines` before journal creation; that helper was removed.
- **Reproduction/evidence:** On the final branch based on `origin/main` SHA `c79bc5e7c3f0b8f32d90c583d31d318c971314ba`, `go test ./internal/usecase -run 'Test(ValidateAutoPostingLines|AutoPostingRejectsUnbalanced)' -count=1`, `go test ./internal/usecase -count=1`, and the full `go test ./...` passed. `go vet ./...`, `go build ./...`, and `git diff --check` passed. An earlier full run on the first fetched base (SHA `208bcaebc48d201ea5f84ac41a12bbc5ce48265a`) exposed duplicate validator declarations and seven repository-test failures; main later advanced, those conflicts were resolved, and the branch was rebased before the successful final run. Disposable PostgreSQL integration was not run because the Docker daemon is unavailable; no existing database or Docker service was started or changed.
- **Baseline scope:** The defect was present on fetched `origin/main` SHA `208bcaebc48d201ea5f84ac41a12bbc5ce48265a`. The implementation branch now contains one focused commit on latest fetched `origin/main` SHA `c79bc5e7c3f0b8f32d90c583d31d318c971314ba`.
- **Recommended solution / components:** Keep validation as the shared gate for order, payroll, supplier, COGS and depreciation auto-posting. Preserve source business amounts and surface validation errors to the caller. Any explicit rounding-account policy remains separate and requires approval.
- **Migration:** Only if adding controlled rounding-account mapping/config; no migration until accountant approves policy. Do not add a balancing plug.
- **Dependencies/policy:** The validator follows the current journal schema's two-decimal storage contract as a conservative rejection boundary. Supported currencies, currency precision, rounding mode, and any permitted residual/account remain unapproved. Do not broaden accepted precision or introduce a residual tolerance until the owner/accountant approves policy.
- **Acceptance:** Execute focused regression tests; verify 100/90 and 90/100 rejection before journal/GL writes, 100/100 success, unchanged values, error propagation and valid sales/payroll flows; run relevant backend tests, `go vet`, build, and disposable PostgreSQL integration. Currency-specific behavior remains policy-dependent.
- **Checklist:** PARTIAL — implementation and Go tests pass; disposable PostgreSQL integration and owner/accountant currency/rounding policy evidence remain pending.

### CFR-002 — Journal integrity relies too heavily on application checks

- **Severity/status:** P1 — MERGED TO MAIN via PR #34; financial acceptance remains IMPLEMENTED — ACCEPTANCE PENDING until all §10 acceptance gates are freshly confirmed.
- **Problem / impact:** Alternate writers, maintenance SQL or races could persist posted unbalanced entries or modify posted details because the database did not enforce journal/ledger integrity.
- **Current implementation:** New migration `backend/migrations/180_enforce_journal_integrity.sql` rejects inconsistent historical data before applying guards; requires new journals to start as draft; permits only draft→posted/void and audited posted→void transitions; validates exact debit/credit/header equality and matching line-to-GL effects at transaction commit; locks journal rows during repository post/void and transitions to posted before GL inserts within the same transaction; prevents posted journal/line and GL mutation, duplicate GL effects, account/restaurant scope reassignment, and cross-tenant/restaurant journal lines. The existing void workflow remains an audited inverse GL effect on the original journal. No history is rewritten.
- **Source/test files:** `backend/migrations/180_enforce_journal_integrity.sql`; read-only `backend/scripts/preflight_180_journal_integrity.sql`; emergency downgrade `backend/scripts/rollback_180_journal_integrity.sql`; `backend/internal/repository/journal_entry_repo.go`; new direct-SQL `backend/internal/repository/journal_integrity_postgres_integration_test.go`; adapted CFR-004 atomic-posting and financial-report integration fixtures; updated posting SQL-mock expectations.
- **Migration evidence (2026-11-13):** PostgreSQL 16.15 disposable container. Full repository migration runner applied all migrations through 180 and 999 on a fresh database. Migration 180 reapplied successfully on a populated disposable database with a clean preflight. Read-only preflight returned zero rows on the valid legacy fixture. A synthetic 10.00/9.00 posted journal was surfaced by preflight; applying migration 180 failed in its initial data guard and transaction rollback left no new guard function/index. The documented downgrade was executed on another disposable database and removed only the CFR-002 triggers/functions/index; it did not rewrite rows.
- **PostgreSQL acceptance evidence (2026-11-13):** Direct SQL balanced three-line posting succeeded. Debit-heavy, credit-heavy, and zero-line posting were rejected. Direct UPDATE/DELETE/INSERT against posted lines, posted header amount/status changes, GL UPDATE/DELETE, duplicate GL insertion, cross-tenant account use, and account/restaurant scope moves were rejected. The restricted `pos_app` role has no direct INSERT/UPDATE/DELETE privileges on journal, line, or GL tables. Six concurrent posting attempts against one draft yielded exactly one successful posting and one GL effect per line; audited void then produced exactly two (original+inverse) effects per line, and duplicate void failed. CFR-004 injected line/GL/status failures rolled back. `go test -tags=integration ./... -p 1 -count=1` passed on the disposable PostgreSQL database after setting its timezone to Africa/Cairo; `go test ./... -p 1 -count=1`, `go vet ./...`, `go build ./...`, and `git diff --check` passed. The two unrelated expiry/session integration tests fail with the disposable database timezone set to UTC and pass with Africa/Cairo; their source is unchanged by this PR. No production or development database was used.
- **Historical-data result:** Existing valid rows were unchanged; preflight is read-only. Migration fails closed on invalid posted balance or missing/duplicate ledger history. No automatic correction, balancing plug, or rounding rule was introduced.
- **Remaining limitations:** CFR-002 preserves the existing reversal representation; closed-period reversal eligibility/policy remains under CFR-005/CFR-013. Database superusers/table owners can disable triggers, so production database ownership and migration credentials must remain restricted. The migration does not implement durable paid-business-event intent/recovery, which remains the CFR-004 gap.
- **PR / commit:** PR [#34](https://github.com/Yousef-Abdallah-1711/pos/pull/34) merged 2026-11-13; implementation commit `37d77d56dae15ee52fd916755f294b3a7b04ccac`. Merge and successful CI establish delivery; the acceptance status above remains conservative pending the complete finding-level release gate.
- **Checklist:** IMPLEMENTED — ACCEPTANCE PENDING.

### CFR-003 — Concurrent journal posting, balance update and numbering races

- **Severity/status:** P1 — SOURCE-CONFIRMED structurally; PostgreSQL race reproduction BLOCKED.
- **Problem / impact:** Concurrent duplicate posting can append duplicate GL effects; simultaneous running-balance read-modify-write can diverge; `MAX()+1` numbering can collide.
- **Current / expected:** Status reads are not shown using `FOR UPDATE`; running balance is computed and persisted per line; number generation uses `MAX()+1`. Expected: serialize state transition, guarantee unique numbering, and calculate ledger balances without lost updates.
- **Root cause:** `backend/internal/repository/journal_entry_repo.go:363-395,544-596,784-798`; analogous payment accounting balance paths in `backend/internal/repository/payment_acct_repo.go`.
- **Reproduction/evidence:** Source review only. No DB concurrency run; ordinary sqlmock tests cannot prove PostgreSQL lock/unique behavior.
- **Baseline scope:** Present on main and current branch (no financial delta).
- **Recommended solution / components:** Lock journal row on post/void, lock affected account rows in stable order or derive balances from posted rows, unique sequence/counter per tenant/restaurant/year with conflict retry. Ensure GL uniqueness per journal line.
- **Migration:** Likely unique constraints/indexes and possibly counter table after collision audit.
- **Dependencies/policy:** Numbering scope and reset period need owner policy; historical duplicates require repair decision.
- **Acceptance:** Concurrent same-journal post yields exactly one GL set; parallel journals produce unique numbers and correct sum; race detector plus PostgreSQL isolation tests pass.
- **Checklist:** OPEN.

### CFR-004 — Journal draft creation and GL posting have a recoverable partial-failure window

- **Severity/status at PR #33 baseline:** P1 — PARTIALLY IMPLEMENTED. Local journal creation and GL posting share one PostgreSQL transaction; the continuation update in §13 records the subsequent durable recovery work.
- **Problem / impact:** Previously, a draft could commit before GL posting failed. Separately persisted order/payment state can still exist without durable accounting intent after a process failure.
- **Current implementation:** `CreateAndPostAutoEntry` obtains a transaction-scoped advisory lock keyed by tenant, restaurant, source type/ID, and entry type; checks existing entries within both tenant and restaurant scope; returns an existing posted entry on retry; rejects historical duplicate, draft, void, or unsupported states; creates journal, lines, GL rows, and posted status using the same transaction. No migration was needed. The repository now writes absent JSONB attachments as SQL NULL so a nil attachment payload is valid PostgreSQL input.
- **Source files changed:** `backend/internal/repository/general_ledger_repo.go`, `backend/internal/repository/journal_entry_repo.go`, `backend/internal/usecase/auto_posting_engine.go`, `backend/cmd/api/main.go`, `backend/internal/usecase/auto_posting_engine_test.go`, `backend/internal/usecase/phase6_accounting_cycle_test.go`, and new `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`.
- **Tests and results (2026-11-13):** Dedicated disposable PostgreSQL 16 with the current migration set applied: concurrent duplicate replay returns one posted journal/two GL rows; identical references in separate tenants/restaurants remain isolated; injected failures at journal-line insert, GL insert, and posted-status update leave zero journal/GL rows and unchanged account balances. Focused CFR-004 unit/usecase suite passed. `go test ./... -p 1 -count=1`, `go vet ./...`, and `go build ./...` passed. The two PostgreSQL integration tests `TestImpersonationHTTPAuditSessionLifecycle` and `TestSubscriptionExpiryWorkerPostgresScopesEachTenantByItsOwnExpiry` fail when the disposable database timezone is UTC and pass when it is `Africa/Cairo`; the affected expiry columns are `TIMESTAMP WITHOUT TIME ZONE`, and these tests/source paths are unchanged by PR #33. The complete tagged integration suite passed with the Cairo timezone configured. No production or existing development database was used.
- **Remaining acceptance gaps:** Paid-order persistence and auto-posting are not one shared transaction; `OrderUseCase` currently logs/swallows auto-posting errors. No durable event/outbox or crash/restart recovery is included. There is no automatic recovery for a pre-existing draft or void entry; those cases fail closed for manual reconciliation. Do not claim end-to-end paid-event recovery or VERIFIED FIXED.
- **Required owner/accountant decisions:** Confirm canonical accounting trigger/source identity across payment/order channels and acceptable synchronous posting behavior. Define operational handling for historical orphan drafts/voided sources. No tax, rounding, refund, or reversal policy is changed in this PR.
- **PR / commit / merge:** PR [#33](https://github.com/Yousef-Abdallah-1711/pos/pull/33) merged 2026-11-13 at `c6f320a19ec5a19f6a9ec59e0fe1cf8059a71a8c`; implementation commit `9c6beb8bd60ba848230c3efa2ccee8cea8e91493`.
- **Acceptance still required at PR #33 baseline:** The atomic transaction, duplicate replay, isolation, rollback, CFR-001 imbalance regression, handler/usecase tests, standard Go suite, `go vet`, and build were verified. PR #33 alone did not include durable paid-business-event intent/recovery; §13 records the continuation and current acceptance gaps.
- **Next action at PR #33 baseline:** Continue the separately scoped durable business-event recovery lifecycle after CFR-002 acceptance. This action is now in progress as recorded in §13.

### CFR-005 — Paid orders can suppress posting failures and have no proven cancellation/refund reversal

- **Severity/status:** P1 — SOURCE-CONFIRMED for best-effort call; refund/cancel reversal missing in examined path.
- **Problem / impact:** Paid orders can be absent from financial reports, and refunded/canceled sales may remain as revenue/receivable/cash; period totals can be wrong.
- **Current / expected:** Handler marks payment state and calls `AutoPostPaidOrderByID` without checking a returned error; no reverse call on paid→refunded/canceled was found. Entry date uses order creation date. Expected: policy-date accounting, failure surfaced/retried, and linked reversal/refund entry.
- **Root cause:** `backend/internal/handler/http/admin_order_handler.go` payment-status handler around lines 470-515; `backend/internal/usecase/order_usecase.go:459-474`; `backend/internal/usecase/auto_posting_engine.go:149`.
- **Reproduction/evidence:** Source-confirmed; current session did not run a database end-to-end transition. Historical order auto-posting test proves only happy-path call behavior, not reversal/failure recovery.
- **Baseline scope:** Present on main/current. Phase 00 local warnings/diagnostics do not reverse entries.
- **Recommended solution / components:** Canonical payment/refund events, atomic update/outbox, post-once business event, linked reversing journal, payment-to-GL reconciliation and recovery dashboard.
- **Migration:** Event identity/outbox/reversal-link columns may be needed after design.
- **Dependencies/policy:** Revenue recognition trigger, delivery/fulfilment timing, refund/chargeback fee allocation, backdate policy.
- **Acceptance:** Paid then partial/full refund yields independently calculated net GL; cancellation before/after capture behaves distinctly; failed post visibly retries; paid order vs sales GL reconciliation has zero unexplained difference.
- **Checklist:** OPEN.

### CFR-006 — Historical report performance symptom not remeasured

- **Severity/status:** P2 — BLOCKED (no current performance fixture/measurement).
- **Problem / impact:** Slow or high-memory financial reports can delay close and time out at larger tenant volumes; no current numeric performance claim is made.
- **Current / expected:** Historical audit carried a performance symptom; this session did not benchmark current main. Financial queries scan posted-ledger/report paths; proper bounded pagination, query plans and dataset size have not been validated here.
- **Root cause:** Not established. Candidate query locations: `backend/internal/usecase/financial_report_usecase.go`, `backend/internal/repository/general_ledger_repo.go`, `accounting_advanced_repo.go`.
- **Reproduction/evidence:** No load test or `EXPLAIN (ANALYZE, BUFFERS)` this session. Therefore symptom and root cause remain unverified.
- **Baseline scope:** Historical finding assessed against main; persistence of performance symptom is unknown. Current branch changes are unrelated.
- **Recommended solution / components:** Establish realistic synthetic dataset, baseline p50/p95/query plans and export memory; optimize only measured hotspots.
- **Migration:** Index/partition only after evidence and migration plan.
- **Dependencies/policy:** Workload SLO, data retention and report freshness.
- **Acceptance:** Recorded dataset/command, explain plans, bounded memory, SLO gates; no unsupported scale claim.
- **Checklist:** BLOCKED.

### CFR-007 — Binary floating point is used across monetary domain and APIs

- **Severity/status:** P2 — SOURCE-CONFIRMED; precision behavior not broadly property-tested.
- **Problem / impact:** Repeated decimal arithmetic, allocation and serialization can produce fractional-cent differences and inconsistent totals across Go, PostgreSQL, JavaScript and exports.
- **Current / expected:** Accounting/order/payroll models and repositories use `float64`, while DB uses `NUMERIC/DECIMAL`; JSON is numeric. Expected: one canonical currency-aware exact representation at boundaries, without silent precision loss.
- **Root cause:** `backend/internal/domain/accounting.go`, `domain/order.go`, `backend/internal/usecase/order_usecase.go:127-182`, `backend/internal/handler/http/hr_salary_handler.go:122-141`; numeric schema e.g. migrations `137`, `141`, `024`.
- **Reproduction/evidence:** Source-confirmed type mismatch; no property test or cross-layer precision experiment in this run.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Define integer minor units or arbitrary-precision decimal contract, exact JSON string/decimal policy, currency exponent, supported rounding modes, allocation residual and FX handling. Migrate incrementally with compatibility converters and reconciliation.
- **Migration:** Potential staged type/API migrations; preserve old records and verify sums before any conversion.
- **Dependencies/policy:** Supported currencies and rounding policy require owner/accountant approval.
- **Acceptance:** property tests for fractional line sums, tax inclusive/exclusive, allocations, split/refund/FX, overflow/negative amounts; API→DB→report/export equality at approved precision.
- **Checklist:** OPEN.

### CFR-008 — Order tax hardcoded; invoice tax computed before discount

- **Severity/status:** P2 — SOURCE-CONFIRMED; jurisdictional expected value requires approval.
- **Problem / impact:** Tax can be under/over-collected or discounts misapplied; amounts may not match local policy.
- **Current / expected:** Order creation sets 10% in Go and doesn't set discount from a configured policy; customer invoice trigger calculates tax on gross line amount, then subtracts discount. Expected: tenant/jurisdiction-configured tax and approved inclusive/exclusive tax base/discount ordering.
- **Root cause:** `backend/internal/usecase/order_usecase.go:147-182`; `backend/migrations/141_create_invoices_customer.sql:78-82` (supplier calculation also review `invoice_supplier_repo.go:449-476`).
- **Reproduction/evidence:** Source formula confirms order 100 → tax 10 absent policy, and gross-based invoice tax formula. No approved jurisdiction/rate or DB test is available; numerical result is formula evidence, not legal assessment.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Tax-rule/config service, effective-dated jurisdiction and item tax category, deterministic allocation and approved calculation order; no embedded legal rate.
- **Migration:** Tax config/effective-date schema may be required; do not rewrite historical tax blindly.
- **Dependencies/policy:** Accountant/legal approval of tax registrations, exemptions, inclusive/exclusive, discount and rounding rules.
- **Acceptance:** Approved jurisdiction fixtures match independently calculated invoice/order tax to the minor unit; changing tenant policy cannot leak across tenants; historic records retain rate/basis.
- **Checklist:** OPEN.

### CFR-009 — Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

- **Severity/status:** P1 — SOURCE-CONFIRMED for missing status checks on update/delete; DB reproduction BLOCKED.
- **Problem / impact:** Payroll history can be rewritten after payment/posting, leaving journal, bank and payroll records inconsistent; a hard delete can orphan its journal.
- **Current / expected:** `UpdateSalary` updates fields without paid-state refusal; `DeleteSalary` invokes repository deletion; `MarkAsPaid` behavior has no proven concurrency-safe one-time precondition. Expected: paid record immutable; correction through approval and linked reversal; unpaid-only deletion; mark-paid idempotent.
- **Root cause:** `backend/internal/handler/http/hr_salary_handler.go:237-382,388-438`; `backend/internal/repository/salary_repo.go:216-227`; migration `024_create_salaries_table.sql`.
- **Reproduction/evidence:** Source review; no live salary fixture. Existing sqlmock payroll auto-post test covers a happy path only.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Salary state machine, guarded atomic update, audit trail, reversal/correction workflow, approval and statutory deduction policy.
- **Migration:** Audit/reversal link or state columns only as designed; preserve payroll history.
- **Dependencies/policy:** Payroll statutory rules/roles and correction approval by owner/accountant.
- **Acceptance:** Update/delete paid returns conflict; duplicate concurrent mark-paid produces one payment/journal; correction preserves original and nets correctly; payroll report equals GL.
- **Checklist:** OPEN.

### CFR-010 — Balance Sheet omits current-period earnings; guardrails do not repair the equation

- **Severity/status:** P1 — IMPLEMENTED / PARTIAL. Main now reports the equation status and blocks Balance Sheet export when unreconciled; the missing current-period earnings calculation and accounting policy remain unresolved.
- **Problem / impact:** Balance Sheet can fail Assets = Liabilities + Equity whenever period earnings are not closed into retained earnings; report/export may still appear usable.
- **Current / expected:** The report query still selects balance-sheet account types only and does not derive current earnings from revenue less expenses. Main now returns `integrity_status`, `equation_difference`, and `integrity_check_basis`; the dashboard warns on unreconciled reports and blocks Balance Sheet export unless the equation passes. This is a visibility/export safeguard, not a calculation of missing equity or proof that a balanced equation is correct. Expected: include correctly derived current earnings or use a completed close/retained earnings model, and keep misleading exports blocked.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:261-336`, especially account filter at 288-293 and equation at 336. Mainline is checked at fetched SHA.
- **Reproduction/evidence:** The source query confirms the earnings omission. Phase 00 report integrity fields, warning, export guard, and corresponding unit/integration test files are present on current `main` after PR #27 (merge commit `c79bc5e7c3f0b8f32d90c583d31d318c971314ba`). The recorded Phase 00 disposable-PostgreSQL fixture expects a 100.00 asset / 0.00 equity report to remain unchanged and be marked `UNRECONCILED`; a fresh isolated PostgreSQL run on this current base was not available in this session. No accountant-approved current-earnings/retained-earnings calculation has been implemented.
- **Baseline scope:** At the assessment baseline, the report guard was missing. Current `main` includes the warning/export safeguard through PR #27, while the accounting equation gap remains. The preserved PR #18 contains a WIP copy of earlier evidence; its relevant source and test changes are already present on main. Do not merge the stale preservation branch wholesale.
- **Recommended solution / components:** Compute as-of net income from posted revenue/expense balances respecting close dates, or enforce close journal and report closed/available state. Preserve equity account presentation and expose equation residual.
- **Migration:** Possible close/retained-earnings fields only if lifecycle design requires them.
- **Dependencies/policy:** Accounting basis, close schedule, prior-period adjustment and report cut-off require accountant approval.
- **Acceptance:** Independent fixtures prove Assets = Liabilities + Equity for open and closed periods; reports/export refuse invalid state; P&L, TB and GL agree. Current evidence covers the guard only, not the corrected equation or report-to-ledger reconciliation.
- **Checklist:** PARTIAL — warning/export gate and tests are on main; current-period earnings and independent reconciliation remain open.

### CFR-011 — Automatic supplier entry type exceeds journal schema width

- **Severity/status:** P2 — REPRODUCED in isolated PostgreSQL; end-to-end supplier posting BLOCKED.
- **Problem / impact:** Auto supplier invoice/payment journal insert can fail because `auto_supplier_invoice` / `auto_supplier_payment` exceed `VARCHAR(20)`; AP operation may persist without its journal depending on transaction boundaries.
- **Current / expected:** Auto-posting constants use those strings; migration 137 sets `entry_type VARCHAR(20)`. Expected: valid bounded type persisted consistently with event catalog, no partial AP operation.
- **Root cause:** `backend/internal/domain/auto_posting.go:19-20`; `backend/migrations/137_create_journal_entries.sql:19`; supplier auto-post caller `backend/internal/usecase/auto_posting_engine.go:472-498`.
- **Reproduction/evidence:** Disposable PostgreSQL 16 rejected `auto_supplier_invoice` with `string_data_right_truncation` against `VARCHAR(20)`. No application-level supplier posting fixture ran.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Choose stable shorter enum/type code or widen schema under migration, and enforce schema/domain compatibility in test. Bind supplier invoice/payment state with posting/outbox.
- **Migration:** Required if retaining current 21-character values; otherwise no schema migration if mapping changed safely. Determine historical values first.
- **Dependencies/policy:** Naming/catalog decision; accountant reviews historical/report usage.
- **Acceptance:** Disposable PostgreSQL supplier invoice and payment post, reverse, list/filter and export succeed; migration upgrade and rollback/forward plan validated.
- **Checklist:** OPEN.

### CFR-012 — Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

- **Severity/status:** P2 — SOURCE-CONFIRMED heuristic; report accuracy not independently reproduced.
- **Problem / impact:** Operating/investing/financing classification may misclassify noncash entries and accrual movements as cash flows, yielding misleading statement.
- **Current / expected:** `GenerateCashFlow` classifies using account types/code patterns, including equity or code prefix heuristics. Expected: statement derived from cash/bank movements and approved direct/indirect method with noncash adjustments.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:340 onward` (classification SQL around 346+; code-prefix filter in same function). Exact source at fetched SHA.
- **Reproduction/evidence:** Static SQL review; no independent cash-flow fixture.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Adopt accountant-approved direct or indirect method; identify cash equivalents and map activities; reconcile opening + net movement = closing cash.
- **Migration:** Optional cash-flow classification metadata/account mapping.
- **Dependencies/policy:** Cash-equivalent scope and method need accountant approval.
- **Acceptance:** Independently calculated fixture with noncash depreciation, working-capital movements, asset purchase and financing; closing cash ties to GL/bank; comparative export parity.
- **Checklist:** OPEN.

### CFR-013 — Voiding a journal writes a reversal dated today without adequate closed-period handling

- **Severity/status:** P2 — SOURCE-CONFIRMED.
- **Problem / impact:** Void can alter a closed/current period unexpectedly; original entry's status and GL representation may be confusing and audit trail may not link explicit separate reversal document.
- **Current / expected:** `VoidEntryTx` uses `time.Now()` then appends reversed GL lines; reviewed path has no explicit closed-period validation for the reversal date. Expected: linked immutable reversal journal in an allowed period, preserving original and controlled reason/approval.
- **Root cause:** `backend/internal/repository/journal_entry_repo.go:423-493` (void date 465-477).
- **Reproduction/evidence:** Source review; closed-period void not executed in PostgreSQL.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Replace implicit GL-only reversal with explicit reversal event/journal; use configured open period and uniqueness, authorization, audit.
- **Migration:** Reversal-link metadata or reversal entry as normal journal; preflight existing voided rows.
- **Dependencies/policy:** Backdate and closed-period correction policy.
- **Acceptance:** Closed-period void rejected or linked reversal in approved open period; no original mutation/deletion; duplicate void prevented; reports show both entries consistently.
- **Checklist:** OPEN.

### CFR-014 — Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

- **Severity/status:** P2 — SUSPECTED / partially source-confirmed; method and end-to-end arithmetic BLOCKED.
- **Problem / impact:** Incorrect cost basis or missing return/wastage reversal makes inventory assets and COGS/profit inaccurate.
- **Current / expected:** Inventory movement and COGS auto-post paths exist (`144_create_inventory_movements.sql`, `auto_posting_engine.go`); examined order path calculates item cost using stored unit cost. Costing method and correspondence among stock movement, cost layer and GL were not proven.
- **Root cause:** Not fully established; candidate paths `backend/internal/usecase/auto_posting_engine.go:155-240`, `backend/internal/repository/inventory_movement_repo.go`, migration 144.
- **Reproduction/evidence:** Historical source skim only; no current PostgreSQL or inventory fixture. Local-only adjustment handler is not origin/main.
- **Baseline scope:** Main has partial inventory/COGS path; completeness defect is unverified on both. Local handler is current-branch-only and unverified as integrated.
- **Recommended solution / components:** Approve FIFO/weighted average/specific cost and recognition timing; one atomic stock movement and journal event; explicit returns, comps, waste and negative-stock treatment.
- **Migration:** Cost layer/valuation schema only after method chosen; reconcile opening stock before rollout.
- **Dependencies/policy:** Inventory valuation method and count/cutoff approval.
- **Acceptance:** Independently calculated receipts/sales/returns/waste/adjustments reconcile quantities, valuation, COGS, GL and inventory report; concurrency and rollback injection.
- **Checklist:** OPEN.

### CFR-015 — Product-wide financial tenant/RBAC isolation not independently proven

- **Severity/status:** P2 — BLOCKED for runtime negative authorization matrix; source review is incomplete, not proof of a bypass.
- **Problem / impact:** Any missing tenant/restaurant predicate, weak role gate, export/background-worker/cache scope can expose or mutate another restaurant's financial data.
- **Current / expected:** Main contains tenant-aware SQL/RLS work and finance route/permission wrappers, but this assessment has no full route-to-query matrix. Expected: trusted identity derives tenant/restaurant; every object and job/export is scope-bound; least-privilege role enforcement.
- **Root cause:** Not established as one defect. Review candidates: accounting routes in `backend/cmd/api/main.go`, financial handlers/repositories, middleware and RLS migrations `155+`.
- **Reproduction/evidence:** No cross-tenant PG/API/browser tests run. Historical Phase 2A tenant artifacts are not financial isolation proof.
- **Baseline scope:** Main and current: completeness unverified; no confirmed current cross-tenant exploit asserted. Local branch modifications mainly concern auth/tenant groundwork but are not financial-specific evidence.
- **Recommended solution / components:** Build route × role × tenant × restaurant × job/export matrix; combine application predicates, scoped transactions and RLS where viable; test cache/key/background contexts.
- **Migration:** RLS/policies/constraints only after coverage audit and safe migration plan.
- **Dependencies/policy:** Role matrix and Super Admin/auditor access approval.
- **Acceptance:** Positive and negative tests for each finance operation across tenant/restaurant; no unauthorized read/write/export; background jobs and cache isolation; RLS bypass-role tests.
- **Checklist:** BLOCKED.

## 5. Additional issues beyond the canonical 15

### FIN-016 — No complete gift-card/deposit/deferred-revenue lifecycle established

- **Severity/status:** P3 — SUSPECTED (scope search did not establish a full lifecycle).
- **Problem / impact:** Gift-card sale/redemption/expiry or customer deposits may be treated as immediate sales or omitted from liability/revenue reporting.
- **Current / expected:** No complete financially integrated state machine was evidenced in inspected financial catalog; expected liability on issue and controlled revenue recognition on redemption/expiry per approved policy.
- **Root cause:** Product capability gap not yet proven by exhaustive route/schema search. Candidate areas: order/payment/accounting APIs and `backend/migrations/136_create_chart_of_accounts.sql` onward.
- **Reproduction/evidence:** No gift-card lifecycle tests located/run; absence is not a claim that no partial code exists.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Confirm product scope, define stored value liability, redemption, breakage/expiry, refund and reconciliation; then implement.
- **Migration:** Likely balance/transaction ledger schema if required.
- **Dependencies/policy:** Owner/accountant/legal policy; local stored-value rules.
- **Acceptance:** Issue/redeem/refund/expire flows reconcile liability, cash and revenue; duplicates and cross-tenant access denied.
- **Checklist:** OPEN (scope confirmation).

### FIN-017 — Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

- **Severity/status:** P3 — BLOCKED for isolated DB fixture.
- **Problem / impact:** Asset register, depreciation expense/accumulated depreciation and disposal proceeds may diverge from GL.
- **Current / expected:** Schema exists at migration 145; complete journal lifecycle not verified. Expected approved useful life, method, convention, disposal and reversal behavior.
- **Root cause:** Not established; candidate migration 145 and auto-posting depreciation path.
- **Reproduction/evidence:** No lifecycle test run.
- **Baseline scope:** Feature exists partially on main/current; completeness BLOCKED.
- **Recommended solution / components:** Define asset classes/schedules and connect acquisition/depreciation/disposal to immutable journals.
- **Migration:** Only if lifecycle fields or event identity required.
- **Dependencies/policy:** Accountant approves capitalization threshold, depreciation method/life and disposal policy.
- **Acceptance:** Schedule independently ties book value/expense/GL at period boundaries; disposal and duplicate monthly run safe.
- **Checklist:** BLOCKED.

### FIN-018 — No complete fiscal close and retained-earnings roll-forward established

- **Severity/status:** P2 — SOURCE-CONFIRMED report gap / workflow incompleteness.
- **Problem / impact:** Open periods lack reliable current earnings in equity; closed-period corrections and year-end carryforward may be inconsistent.
- **Current / expected:** Fiscal periods and posting checks exist, but a full close/reopen/earnings transfer/audit workflow was not established. Expected role-gated, repeatable close with TB/report lock and linked corrections.
- **Root cause:** Missing/partial lifecycle across `backend/migrations/153_create_fiscal_periods.sql`, `financial_report_usecase.go`, journal repository.
- **Reproduction/evidence:** Static inventory; no fiscal close test or direct DB run.
- **Baseline scope:** Main/current completeness not established.
- **Recommended solution / components:** Design close state machine, lock, retained earnings policy and reopen/adjustment event. Coordinate CFR-010/013.
- **Migration:** Possible close run/approval and lock metadata after design.
- **Dependencies/policy:** Fiscal year, adjustment periods, retained earnings, reopening authority require accountant/owner approval.
- **Acceptance:** Close twice is idempotent; posting blocked in closed period; current earnings roll forward exactly; reopen/reversal audited; all statements reconcile.
- **Checklist:** OPEN.

### FIN-019 — Export parity and native XLSX / CSV safety are unverified

- **Severity/status:** P3 — BLOCKED for complete export comparison.
- **Problem / impact:** Export totals/filters can disagree with on-screen reports; spreadsheet formulas may execute from user-controlled text; CSV implementation can be mistaken for XLSX support.
- **Current / expected:** CSV/PDF behavior has historical snapshots and targeted API/UI unit coverage; native XLSX parity and formula injection behavior were not demonstrated. Expected schema, exact decimal strings, filters, bounded streaming, audit provenance and safe cell escaping.
- **Root cause:** Export handlers/clients not fully verified; candidates `backend/internal/handler/http/financial_report_handler.go`, `report_handler.go`, frontend `reportsApi.ts`.
- **Reproduction/evidence:** Frontend report tests passed (8) for API parsing/filter use and CSV-export interaction; no export file adversarial fixture/browser download test.
- **Baseline scope:** Feature parity/controls unknown on main/current.
- **Recommended solution / components:** Publish versioned export contract and independent reconciliation fixture; escape formula-leading cells, cap/stream, authorize/audit.
- **Migration:** Export audit schema may exist; verify before adding.
- **Dependencies/policy:** PII/retention and maximum export size.
- **Acceptance:** CSV/PDF/XLSX, if supported, reconcile to same ledger fixture; formula payload safe; role denied; limits/memory proven.
- **Checklist:** OPEN.

### FIN-020 — Multi-currency / FX and consolidated financial reporting policy is not established

- **Severity/status:** P3 — SUSPECTED policy and implementation gap.
- **Problem / impact:** Combining restaurants/currencies without transaction, functional and presentation currency rules can make tenant-level consolidated reports materially wrong.
- **Current / expected:** Accounts and payment flows have currency fields, but conversion rates, dates, realized/unrealized FX and consolidation eliminations were not verified. Expected immutable original currency plus policy-defined base-currency valuation/reporting.
- **Root cause:** Not established; review candidates `domain/accounting.go`, payment/invoice repositories and report usecase.
- **Reproduction/evidence:** No FX test or consolidated report fixture.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Keep FX as explicit policy/config with stored rate source/time and realized/unrealized gains/losses; scope consolidation.
- **Migration:** Potential rate/functional currency metadata after policy approval.
- **Dependencies/policy:** Functional currency, rate source/time, translation and consolidation rules; accountant approval.
- **Acceptance:** Independent multi-currency fixtures tie transaction-currency journal, base-currency ledger, FX gain/loss and consolidated statements.
- **Checklist:** OPEN.

### FIN-021 — Root HTML language/direction does not reflect the active locale

- **Severity/status:** P3 — REPRODUCED in a current-run browser check; localized report content itself did render RTL.
- **Problem / impact:** The browser rendered an Arabic financial report under a nested locale wrapper with `dir="rtl"`, but the document root retained `lang="en"` and no `dir`. Assistive technology, browser text handling, and document-level direction-sensitive behavior can therefore use English/LTR semantics around Arabic content. This is an accessibility/localization defect, not evidence of incorrect financial values.
- **Current / expected:** `frontend/apps/dashboard/src/app/layout.tsx:20-21` hardcodes `<html lang="en">` and leaves `<body>` without a direction. The locale subtree controls its own direction. Expected: root document language and direction correspond to the selected locale, or the application has a documented locale-root strategy that ensures document semantics are correct.
- **Root cause:** Root layout is not locale-aware while localized routes nest under `[locale]`; a nested element cannot change the document `<html>` attributes.
- **Reproduction/evidence:** In a local dashboard browser session on 2026-11-12, Arabic Balance Sheet content appeared, the localized wrapper computed RTL, but `document.documentElement.lang` was `en`, `document.documentElement.dir` was empty, and body direction was LTR. No page/console errors; mobile viewport had no horizontal overflow; PDF download completed. API/auth responses were mocked, so this did not exercise a real authenticated backend workflow.
- **Baseline scope:** Source at fetched `origin/main` and current worktree contains the hardcoded root layout; no independent clean-main browser build was run. Treat as present in both by source comparison, with runtime behavior reproduced on the current worktree.
- **Recommended solution / components:** Make the root layout locale-aware using the framework-supported locale routing pattern and set document-level `lang`/`dir`; verify English and Arabic routes, hydration, nested direction, assistive semantics, and report export.
- **Migration:** None.
- **Dependencies/policy:** Supported locale list and language/direction mapping; frontend routing architecture.
- **Acceptance:** Browser assertions verify root lang/dir for each supported locale; Arabic report and export remain usable; no hydration warnings or horizontal overflow at agreed viewports.
- **Checklist:** OPEN.

## 6. Test execution matrix

The commands below identify the executed test programs/arguments. Go tests, vet, build, migration execution, and the opt-in RLS test ran in an isolated disposable environment; the full Go suite had the disposable PostgreSQL URL/role configured for DB-backed packages. The exact shell wrapper and ephemeral connection string were not retained in the current evidence, so no secret/URL is reconstructed here. The full current worktree was tested, not a clean checkout of `origin/main`. No test used a protected or pre-existing database. The runner applied all 109 active `.sql` migration files in this worktree, including local migration 162. The directory contains 16 additional `.disabled` entries that were not applied. This proves current-worktree migration application on PostgreSQL 16, not main-only migration parity, rollback, or production upgrade safety.

| Suite | Exact command | Environment / result | Count |
|---|---|---|---:|
| Full backend Go suite | `go test ./...` | Disposable PostgreSQL 16 and synthetic test DB configuration; completed before Docker daemon became unavailable. | **726 passed, 0 failed, 2 skipped** |
| Go static analysis | `go vet ./...` | Same isolated Go toolchain environment. | **PASS** (not a test-case count) |
| Go build | `go build ./...` | Same isolated Go toolchain environment. | **PASS** |
| Ordered migration application | Repository migration runner (`go run ./cmd/migrate`, with isolated DB configuration) | Applied **109/109 active `.sql` files** in current worktree to disposable PostgreSQL 16. The directory has 125 entries total: 109 `.sql` files plus 16 `.disabled` files. Includes local-only migration 162. No downgrade/rollback sequence executed. | **109 applied, 0 failed** |
| Opt-in RLS/audit repository integration test | `go test ./internal/repository -run '^TestSuperAdminAuditRepositoryAcceptsDashboardAccessInDisposablePostgres$' -count=1` | Disposable PostgreSQL, required `_test` database and restricted `pos_app` role; initial guarded attempt with wrong URL did not run, corrected invocation passed. | **1 passed** |
| Direct SQL integrity probes | Synthetic disposable PostgreSQL statements (not a Go test suite) | Valid JSONB attachment and negative-amount CHECK rejection behaved as expected; unbalanced posted journal and posted-line update were accepted; supplier entry type exceeded `VARCHAR(20)`. | **2 expected pass; 3 defects reproduced** |
| Dashboard financial/report UI tests | `pnpm exec vitest run 'src/lib/api/reportsApi.test.ts' 'src/components/reports/reports-rtl.test.tsx' 'src/app/[locale]/dashboard/reports/sales/page.test.tsx' --reporter=verbose` | Dashboard workspace, synthetic/mock UI data. | **8 passed, 0 failed** |
| Dashboard typecheck | `pnpm typecheck` | Dashboard workspace; pass. | **PASS** (not a test-case count) |
| Dashboard production build | `pnpm build` (synthetic local public URLs supplied) | Dashboard workspace; pass with repeated webpack “Critical dependency: request of a dependency is an expression” warnings from component registry. | **PASS with warnings** |
| Browser Arabic Balance Sheet / PDF / responsive check | Local dashboard browser workflow (API/auth mocked) | Arabic report rendered; locale wrapper RTL; no page/console errors; mobile viewport no horizontal overflow; PDF downloaded. Root document still `lang=en`, `dir` empty, body LTR (FIN-021). This is not real API/auth E2E. | **Workflow passed; FIN-021 reproduced** |

### Requested test areas and evidence-backed outcome (12 categories)

| Test area | Status | Why / next safe execution |
|---|---|---|
| Full backend Go suite, vet, build | PASS | `go test ./...`: 726 passed, 0 failed, 2 skipped; `go vet ./...` and `go build ./...` passed. The two skips remain unexecuted tests, not failures. |
| PostgreSQL complete migration application / upgrade / rollback | PARTIAL | All 109 active current-worktree `.sql` migrations applied to disposable PostgreSQL 16; 16 `.disabled` entries are not executable SQL migrations and were not applied. Rollback/downgrade, clean-main schema comparison, and upgrade from realistic prior production snapshots were not run. DB is unavailable now because Docker daemon stopped; do not restart unrelated Docker services. |
| Journal direct-SQL constraints, JSONB, supplier type | FAIL — defect reproduced | Five probes: 2 expected behaviors pass, 3 defects reproduced (unbalanced posted row accepted, posted line mutable, supplier type too long for schema). No complete app posting integration fixture. |
| Concurrent posting, running balance, numbering, duplicate capture/refund | BLOCKED | Docker PostgreSQL became unavailable before a concurrency harness could run. No safe standalone DB alternative was available without touching existing services. Resume with fresh disposable DB, duplicate concurrent event IDs, numbering/load and row-lock stress. |
| Transaction rollback/failure injection across order→payment→journal→GL | BLOCKED | Full migrated DB existed, but injection/restart/replay scenarios were not completed before its daemon disappeared. Must run with synthetic records and disposable DB only. |
| JSONB attachment and supplier invoice/payment posting | PARTIAL | JSONB attachment persistence probe passed; supplier type-length constraint defect reproduced. Supplier invoice/payment lifecycle, reversal and reconciliation were not run. |
| GL/TB/Balance Sheet/P&L/Cash Flow/AR/AP numerical reconciliation | PARTIAL | UI tests and mocked browser report rendering passed; no independent seeded ledger fixture was compared across all report formulas. |
| Refund/cancel/chargeback and settlement discrepancy | PARTIAL | Full Go suite passed, but no provider-to-ledger/refund/cancel/settlement fixture was executed. Remains unverified end to end. |
| Supplier, inventory/COGS, payroll, bank reconciliation, depreciation | PARTIAL | Full migration schema loaded and full Go suite passed; supplier column-width defect reproduced. Workflow-level postings, inventory valuation, payroll mutation and bank statement reconciliation lack synthetic accounting fixtures. |
| RBAC, cross-tenant/restaurant, RLS, export/background negative matrix | PARTIAL | One restricted-role RLS/audit repository integration test passed. API route × role × tenant, exports, background jobs and restaurant boundaries were not comprehensively tested. |
| Frontend typecheck, build, financial UI suite | PASS for selected scope; PARTIAL overall | Dashboard typecheck/build passed; 8 financial/report Vitest cases passed. Other frontend apps, full dashboard suite and actual file-content adversarial export checks were not run. Build had component-registry warnings. |
| Browser Arabic RTL, hydration, report export, responsive and authenticated workflows | PARTIAL | Mocked browser workflow verified report subtree RTL, no JS/page errors, mobile overflow, and PDF download; FIN-021 root lang/dir mismatch reproduced. No live API/auth financial workflow; no comprehensive hydration assertion suite. |

**Count convention:** 737 automated test cases/checks passed (726 full Go suite + 1 standalone RLS integration + 8 Vitest + 2 expected SQL probe outcomes); 3 SQL invariant probes failed by demonstrating defects; 2 Go test cases were skipped. The browser workflow is reported separately because its assertion count was not captured as a formal runner result. `go vet`, build, typecheck, migrations, and browser workflow are outcome rows, not added to test-case totals. The 12 requested categories are now classified with evidence, but several are PARTIAL or BLOCKED rather than complete. Docker daemon unavailability is a real external limitation for further disposable PostgreSQL runs after substantial setup and successful full-schema application; it is not being used to excuse unfinished migration preparation.

**Twelve-category disposition:** 1 PASS (full Go/vet/build category), 1 FAIL (SQL invariants reproduced), 8 PARTIAL, and 2 BLOCKED (concurrency and transaction failure injection). “PASS” applies only to the named suite or scope; it does not imply financial-module release readiness.

## 7. Priority implementation roadmap

Phases are proposals only. Do not begin implementation until a new development instruction. Each phase needs independent PR boundary and accounting review where noted.

### Phase 0 — Safety and visibility (engineering-ready, urgent)

- **Issues:** CFR-001, CFR-004, CFR-005, CFR-010.
- **Code areas:** `auto_posting_engine.go`, order payment handler/usecase, financial report usecase/handler/export, alert/reconciliation surfaces.
- **Prerequisites:** Preserve mainline/local branch distinction; identify synthetic PG environment; no invented tolerance or tax value.
- **Tests:** Reject material imbalance; failure-injection for journal/GL write; paid order without posted event visible; Balance Sheet warns and export gate tested.
- **Gate:** No silent posting failure; no material auto-plug; report invalidity is unmistakable; no orphan event can be ignored without reconciliation alert.
- **Migration:** Only if an outbox/event identity is selected; safe forward recovery required.
- **Rollback/correction:** Feature flag or forward correction; never delete posted journals.
- **Independent PRs:** imbalance validation; posting failure visibility; report guard/equation; order/refund lifecycle separate.
- **NO-GO:** Missing alert/retry, or material difference is still plugged.

### Phase 1 — Money, journal and source-of-truth contracts (policy + engineering)

- **Issues:** CFR-002, CFR-003, CFR-007, CFR-011.
- **Code areas:** domain money types, JSON DTOs, journal repository/schema, supplier entry catalog, account/ledger persistence.
- **Prerequisites:** Owner/accountant currency, precision, numbering, rounding and historical-repair decisions.
- **Tests:** Currency property/boundary tests; direct SQL invariant tests; concurrency uniqueness/immutability; migration/rollback in disposable PG.
- **Gate:** exact amounts cross API/DB/report; DB rejects invalid posted journals; supplier post succeeds; no duplicate ledger rows.
- **Migration:** staged money conversion/invariants/entry type only after historical preflight.
- **PR boundaries:** Money contract; ledger invariants; numbering/concurrency; supplier compatibility separate.
- **NO-GO:** silent rounding, irreversible money conversion, unexplained historical imbalance.

### Phase 2 — Reliable order/payment/refund accounting

- **Issues:** CFR-004, CFR-005, CFR-003; cross-channel inventory/workflow gaps.
- **Code areas:** order usecase and all order-channel handlers, provider adapters/webhooks, event/outbox/inbox, auto-posting, reconciliation worker.
- **Prerequisites:** Recognition event, capture vs settlement, deposit/tip/refund/chargeback policy.
- **Tests:** all order channels; duplicate/out-of-order callback; partial/full refund; cancellation pre/post capture; provider and DB failure injection; retry/crash/replay; settlement discrepancy.
- **Gate:** order/payment source totals reconcile to journal/GL with zero unexplained difference; one event has one effect and durable recovery.
- **Migration:** outbox/inbox/event identity if chosen; backfill/reconciliation before enabling.
- **PR boundaries:** canonical event + outbox; sales posting; refunds/reversals; settlement/reconciliation.
- **NO-GO:** operational paid state without durable accounting outcome.

### Phase 3 — Invoices, supplier accounting and payroll lifecycle

- **Issues:** CFR-008, CFR-009, CFR-011; supplier and customer invoice correctness.
- **Code areas:** migrations 141/142, invoice repos, payment accounting repo, salary handler/repo, policy config.
- **Prerequisites:** Tax and payroll accountant decisions, statutory jurisdiction and roles.
- **Tests:** discount/tax independent fixtures; AR/AP payment schedules and overpayment; duplicate supplier post; paid salary immutability, reversal and payroll-to-GL tie out.
- **Gate:** invoice DB totals equal independent calculation; payment changes journal and invoice atomically; paid payroll cannot mutate silently.
- **Migration:** audit/reversal/policy metadata only after data preflight.
- **NO-GO:** tax legal assumption embedded or payroll edit can orphan payment journal.

### Phase 4 — Inventory, COGS and fixed assets

- **Issues:** CFR-014, FIN-017; any stock/COGS mismatch found in fixtures.
- **Code areas:** inventory movement/cost valuation, order COGS, waste/return flows, fixed asset/depreciation repositories and auto-posting.
- **Prerequisites:** Accountant-approved inventory cost method, negative stock, wastage, capitalization and depreciation rules.
- **Tests:** quantity/cost layer/GL reconciliation; return, waste, partial cancellation, failure rollback, duplicate monthly depreciation, disposal.
- **Gate:** inventory asset value and COGS tie to ledger and count; asset schedule ties to GL.
- **Migrations:** cost layers or event identities only after policy/data assessment.
- **NO-GO:** unsupported costing assumption or non-reconciling opening inventory.

### Phase 5 — Reporting, period close and reconciliation

- **Issues:** CFR-006, CFR-010, CFR-012, CFR-013, FIN-018.
- **Code areas:** financial report usecase/repositories, ledger projections, close/period state, exports.
- **Prerequisites:** Accounting basis, fiscal calendar, cash-flow method, adjustment/reopen rules.
- **Tests:** independent TB/BS/P&L/cash-flow/aging fixtures, equations, close/reopen, void/reversal in closed period, exports and report load profile.
- **Gate:** every report reconciles to posted ledger; Assets = Liabilities + Equity; opening + cash movement = closing cash; close locks correctly.
- **Migration:** close-run/approval and classifications as approved.
- **NO-GO:** report differs from ledger or export suppresses balance exception.

### Phase 6 — Financial security, export and operations

- **Issues:** CFR-015, FIN-019; observability/reconciliation cross-cuts.
- **Code areas:** finance routes/RBAC, tenant-scoped queries/RLS, jobs/cache, export code, audit and alerting.
- **Prerequisites:** Approved role matrix, Super Admin/auditor access, export retention.
- **Tests:** complete negative tenant matrix, API/RLS/batch/export tests, formula injection, audit provenance, alert and replay drills.
- **Gate:** no unauthorized cross-tenant view/mutation; all unmatched events/settlements age-alert and have assigned recovery.
- **Migration:** RLS/index/audit changes only with safe backfill and lock analysis.
- **NO-GO:** an unscoped background job or export path.

### Phase 7 — Capacity, multi-currency and commercial release verification

- **Issues:** CFR-006, FIN-020 and any load/security blockers.
- **Code areas:** measured report/query hot paths, partitions/indexes, connection pools/work queues, FX/consolidation.
- **Prerequisites:** Approved SLOs, realistic synthetic dataset, currency and consolidation policy.
- **Tests:** capacity profile (daily volume/peak TPS), 1k-restaurants model, query plans, export memory, backup/restore, disaster recovery, full browser release matrix.
- **Gate:** documented measured SLOs and recovery RPO/RTO; no capacity claim from projection alone; FX/consolidated equations independently tie.
- **NO-GO:** unmeasured capacity, unapproved FX or any unresolved P1.

## 8. Owner and accountant approvals (keep provisional until signed)

### Owner decisions

- Supported currencies, restaurant/tenant consolidation scope, legal entities, roles and approval thresholds.
- Order recognition event, channel behavior, deposits, service charges/tips, gift cards and marketplace settlement treatment.
- Fiscal year, period close schedule, operational downtime and recovery targets.
- Export formats, retention, maximum export size and permitted Super Admin/auditor access.

### Qualified accountant / tax adviser decisions

- Jurisdiction, accounting framework/basis, VAT/tax registration/rates/inclusive versus exclusive rules, discounts/exemptions and historical correction.
- Inventory valuation/COGS timing, spoilage/waste, negative stock, count adjustments and returns.
- Payroll statutory deductions and posting/correction rules.
- FX functional currency, rate source/date, realized/unrealized treatment and consolidation.
- Rounding mode/tolerance/account, void/reversal date, closed-period adjustments, retained earnings and cash-flow method.
- Gift card/deposit deferred revenue, fixed-asset capitalization/depreciation/disposal, audit retention.

No legal, IFRS, tax or jurisdictional compliance is asserted.

## 9. Proposed architecture choices (not accepted)

| Option | Description | Trade-offs |
|---|---|---|
| A. Harden synchronous posting | Update operational row, journal and GL in a single PostgreSQL transaction for local effects. | Lowest moving parts and immediate consistency; provider calls cannot be in that transaction; long transactions and coupling need control. |
| B. Transactional outbox + durable worker | Persist operational state and outbox event atomically; worker validates, posts idempotently, retries and reconciles. | Strong recoverability and decoupled retries; eventual report consistency, outbox ops/alerting and idempotency complexity. |
| C. Hybrid (provisional recommendation) | Synchronously persist critical local business state + event/outbox and journal where feasible; external provider settlement and heavy projections asynchronously reconcile with durable inbox/outbox. | Balances immediate durable accounting intent and external-system realities; needs explicit user-visible pending-accounting state and disciplined reconciliation. |

**Proposed direction:** C for externally initiated payment/settlement flows, A where business event and ledger are entirely local and can share one short DB transaction. No architecture is accepted; first confirm the present repository can safely share transaction boundaries. Exactly-once network delivery is not promised; use at-least-once delivery with idempotent effects.

## 10. Definition of completion and release gate

Do not close any issue based on code presence, a single happy-path mock, or historical screenshot. A closure record must link commit/PR, current-main comparison, migration upgrade/rollback evidence, independently calculated fixtures, targeted regression tests, cross-tenant tests where applicable, operational monitoring/recovery, and owner/accountant approval dependencies. All P1 issues must be resolved or explicitly accepted by accountable owners before commercial release; statutory/accounting decisions must be signed by qualified advisers.

## 11. Prioritized executable tracker

Run one issue per implementation session unless prerequisites couple it. “VERIFIED FIXED” must be changed only with fresh evidence in a future session; no issue is currently fixed.

### P1 — implement in this order

- [ ] **PARTIAL CFR-001** Reject unbalanced auto-posting before persistence. Focused, package, and full Go suites plus vet/build pass; disposable-PostgreSQL acceptance remains pending. Any currency-specific rounding behavior requires owner/accountant policy approval.
- [ ] **PARTIALLY IMPLEMENTED CFR-004** Main contains the atomic journal/GL transaction from PR #33. Durable recovery is in progress in the continuation worktree; see §12 for current evidence and PR state.
- [ ] **OPEN CFR-005** Surface paid-order posting failures and implement linked refund/cancel reversal events. Gate: order/payment/refund totals reconcile to GL.
- [ ] **IMPLEMENTED — ACCEPTANCE PENDING CFR-002** PR #34 and migration 180 are merged. PostgreSQL invariants and CI evidence are recorded above; retain the acceptance-pending label until the remaining §10 release gates are confirmed.
- [ ] **OPEN CFR-003** Serialize journal state/number/balance writes and prevent duplicate GL effects. Gate: concurrent PG tests.
- [ ] **PARTIAL CFR-010** Warning/export gate is on main; derive current earnings or implement approved close/retained earnings and pass independent statement fixtures before closure.
- [ ] **OPEN CFR-009** Lock paid payroll and correct through reversals. Gate: update/delete/duplicate mark-paid tests and payroll-to-GL reconciliation.

### P2 — policy dependencies first where marked

- [ ] **BLOCKED CFR-006** Reproduce historical report performance symptom using realistic disposable data; do not optimize without measurements.
- [ ] **OPEN CFR-007** Approve money representation/precision, then migrate domain, DB and JSON consistently.
- [ ] **BLOCKED CFR-015** Execute financial RBAC, RLS, tenant/restaurant negative matrix; convert any confirmed gap to separate finding.
- [ ] **OPEN CFR-008** Obtain tax policy; replace hardcoded rate and correct discount/tax calculation order.
- [ ] **OPEN CFR-011** Resolve supplier entry type/schema width; prove supplier posting in PostgreSQL.
- [ ] **OPEN CFR-012** Obtain cash-flow method and implement/reconcile cash movements.
- [ ] **OPEN CFR-013** Implement linked reversing journal and closed-period correction rule.
- [ ] **BLOCKED CFR-014** Approve inventory costing; execute quantity/valuation/COGS end-to-end fixtures.
- [ ] **OPEN FIN-018** Design and implement close/retained-earnings lifecycle after accountant approval.

### P3 — confirm scope and close after core ledger/reports stabilize

- [ ] **OPEN FIN-016** Confirm gift-card/deposit scope; define liability/redemption accounting if in product.
- [ ] **BLOCKED FIN-017** Verify fixed asset/depreciation/disposal lifecycle against an isolated DB fixture.
- [ ] **OPEN FIN-019** Define and test CSV/PDF/native XLSX export parity and injection/size protections.
- [ ] **OPEN FIN-020** Confirm multi-currency/consolidation product scope and FX policy.
- [ ] **OPEN FIN-021** Set root document language/direction from active locale; verify Arabic/English browser semantics and hydration.

---

## 12. Final session summary

- **Execution status:** PARTIAL. Full Go suite, vet/build, active SQL migration application, selected frontend checks, limited RLS test and mocked browser checks completed. Several multi-step financial workflows and security matrices remain partial or blocked.
- **Current main SHA:** `a6a2c6700a5fe883256deb1e5a3893108bd9865c` (fetched in this session).
- **Current branch/HEAD:** `codex/phase2a-continuation` / `ecadb5ab7665df3eb9f69886a9b186c710eff62d`.
- **Worktree safety:** existing 244-path dirty/untracked state preserved; no branch change/reset/rebase/merge/cherry-pick. Existing unrelated Docker services were not touched.
- **Tests:** 737 automated test cases/checks passed, 3 direct SQL invariant probes failed by reproducing defects, and 2 Go test cases skipped. Dashboard browser workflow passed its recorded checks with mocked API/auth; this is reported separately, not added to the test-case total. Vet, build, typecheck and 109 active SQL migration applications passed. Twelve major categories are classified; many remain PARTIAL/BLOCKED.
- **Twelve-category outcome:** 1 PASS, 1 FAIL with reproduced defects, 8 PARTIAL, 2 BLOCKED.
- **Issues:** 21 total — P0: 0, P1: 7, P2: 9, P3: 5. None fixed. CFR-002 and CFR-011 have directly reproduced database defects; FIN-021 was reproduced in browser semantics; all others retain their prior status unless updated by evidence above.
- **First implementation phase:** Phase 0 — Safety and visibility, addressing CFR-001, CFR-004, CFR-005 and CFR-010 with explicit imbalance rejection, durable/visible posting failures, and an unmistakable invalid Balance Sheet/export guard. Do not add a rounding tolerance without owner/accountant approval.
- **Recommended sequence after Phase 0:** CFR-002 → CFR-003 → CFR-009; then CFR-007 after monetary policy and CFR-008 after tax policy; then supplier, Cash Flow, void/reversal, inventory/COGS and period close; finally tenant isolation, export/locale, FX and measured scale validation.
- **Source modified:** YES — only this requested Markdown tracker. Application source: NO.
- **Protected database modified:** NO. **Disposable synthetic database used:** YES; 109 active `.sql` migrations applied (16 `.disabled` entries not applied), then the no-volume container became unavailable after Docker stopped.
- **Commits created:** NO. **Pushed:** NO. **PRs opened:** NO. **Merged:** NO. **Deployment action in this session:** NO (live deployment state not checked).
- **Unverified limitations:** No rollback/downgrade test, concurrent posting/load test, transaction failure injection, end-to-end refund/cancel, supplier invoice/payment accounting, inventory/COGS valuation, payroll mutation, bank reconciliation, independent financial statement equation fixture, complete tenant/API/export/background negative matrix, or real backend-authenticated browser financial workflow was run. The RLS result covers one repository integration test only. Browser API/auth were mocked. The disposable Docker-backed DB became unavailable; Docker was not restarted. Historical 2026-11-10 screenshots and audit snapshots are not current test results. No financial correctness, production-readiness, IFRS, tax or scale guarantee is made.

**Immediate blocker to further DB testing:** Docker CLI reported it could not connect to `npipe:////./pipe/dockerDesktopLinuxEngine` because the daemon was unavailable. The full disposable migration/test setup had already succeeded; restarting Docker was intentionally avoided because it may affect unrelated containers/services. Resume PostgreSQL-only suites in a newly authorized disposable environment that does not reuse any existing DB or service.

## 13. Current continuation status — 2026-11-13

This section supersedes §12's earlier Git/test snapshot where they differ. It records the current continuation from fetched `origin/main` `fc8647d6a0b802c3abf7bc3e93dc0b69e1ca9e8d`. Engineering delivery and financial acceptance are separate: merged code is not automatically financially accepted. Findings not changed in this continuation retain their issue-specific evidence and next actions in §§4–5.

### Current 21-finding status snapshot

| Finding | Engineering delivery | Financial acceptance | Next action |
|---|---|---|---|
| CFR-001 | MERGED TO MAIN — PR #29 | PARTIALLY IMPLEMENTED | Add disposable-PostgreSQL acceptance; do not invent rounding policy. |
| CFR-002 | MERGED TO MAIN — PR #34 | IMPLEMENTED — ACCEPTANCE PENDING | Confirm complete finding-level release gate; PR checks passed. |
| CFR-003 | NOT STARTED | BLOCKED BY ENVIRONMENT | Reproduce numbering/balance races in disposable PostgreSQL. |
| CFR-004 | PR #33 MERGED; PR #35 OPEN (ready for review) | IMPLEMENTED — ACCEPTANCE PENDING | Obtain independent approval and operational review. |
| CFR-005 | NOT STARTED | OPEN | Define payment/refund policy and linked reversal workflow after CFR-004. |
| CFR-006 | NOT STARTED | BLOCKED BY ENVIRONMENT | Build realistic disposable report-performance fixture and measure. |
| CFR-007 | NOT STARTED | OPEN | Obtain money precision/representation decisions before migration. |
| CFR-008 | NOT STARTED | BLOCKED BY POLICY | Obtain jurisdictional tax and discount rules before changing calculations. |
| CFR-009 | NOT STARTED | OPEN | Add paid-payroll mutation and duplicate-marked integration coverage. |
| CFR-010 | MERGED TO MAIN — existing report guard | PARTIALLY IMPLEMENTED | Implement or approve current-earnings/close treatment and independent statements. |
| CFR-011 | NOT STARTED | BLOCKED BY ENVIRONMENT | Correct schema compatibility and prove supplier posting in PostgreSQL. |
| CFR-012 | NOT STARTED | OPEN | Define cash movement method and reconcile a cash-flow fixture. |
| CFR-013 | NOT STARTED | OPEN | Define closed-period correction and linked reversal behavior. |
| CFR-014 | NOT STARTED | BLOCKED BY POLICY | Obtain inventory costing/COGS timing policy and test end-to-end. |
| CFR-015 | NOT STARTED | BLOCKED BY ENVIRONMENT | Run financial RBAC/RLS tenant and restaurant negative matrix. |
| FIN-016 | NOT STARTED | BLOCKED BY POLICY | Confirm gift-card/deposit scope and liability treatment. |
| FIN-017 | NOT STARTED | BLOCKED BY ENVIRONMENT | Run fixed-asset lifecycle fixtures in isolated PostgreSQL. |
| FIN-018 | NOT STARTED | OPEN | Obtain close/retained-earnings policy and implement the lifecycle. |
| FIN-019 | NOT STARTED | BLOCKED BY ENVIRONMENT | Compare CSV/PDF/XLSX outputs and exercise safety/size cases. |
| FIN-020 | NOT STARTED | BLOCKED BY POLICY | Confirm currencies, consolidation scope and FX policy. |
| FIN-021 | NOT STARTED | OPEN | Set root language/direction from locale and verify browser behavior. |

### CFR-004 durable recovery continuation

- **Engineering delivery:** PR OPEN (draft) on `fix/cfr004-durable-accounting-recovery` in `H:\Projects\pos-cfr004-durable-recovery`, based directly on the current main SHA above. PR [#35](https://github.com/Yousef-Abdallah-1711/pos/pull/35), implementation commit `34dd55a9d274003ac0f5880f5e5009a695fb7ac9`.
- **Financial acceptance:** IMPLEMENTED — ACCEPTANCE PENDING. PR checks pass; do not mark VERIFIED FIXED until independent approval, migration review and remaining operational limits are accepted.
- **Source and migration changes:** transactional paid-status plus outbox enqueue in `OrderRepository`; scoped admin handler now reports `accounting_intent`; durable repository supports stable event identity, deduplicated enqueue, `SKIP LOCKED` claim, lease expiry/generation fencing, retryable and terminal failure states, bounded exponential backoff and paid-order reconciliation; worker starts with API lifecycle and logs processing/reconciliation failures. `AutoPostingEngine` now writes order/product inventory sale movements through stable idempotency keys and recognizes legacy movements before backfill replay. Migration 181 adds/backfills the outbox, 182 indexes paid-order recovery scans, 183 fences stale worker claims, and 184 protects replayed order-sale inventory movements.
- **Acceptance tests:** On isolated disposable PostgreSQL 16 with database timezone `Africa/Cairo`, migration 184 applied through the repository migration runner. Focused tests passed for paid-status/outbox rollback, duplicate delivery and exactly-once sales journal/GL after retry, concurrent claims, expired lease recovery, stale-lease fencing, missing-intent reconciliation, bounded terminal failure, tenant/restaurant scope, legacy movement recognition, injected movement failure rollback and concurrent movement replay. Full `go test -tags=integration ./... -count=1`, ordinary `go test ./...`, `go vet ./...`, `go build ./...`, and `git diff --check` passed on the implementation before the final tracker/PR-reference-only edits. The full tagged run passed with the disposable database in `Africa/Cairo`; the known session/expiry timezone tests therefore passed in the required environment. No production or existing development database was used.
- **Remaining limits:** The work is not merged or independently reviewed. Migration 184 forward-applied successfully, and the repository's full migration runner passed in GitHub backend CI against its fresh PostgreSQL service; its rollback script was not separately exercised. Provider-side settlement/webhooks are not coupled to this order outbox. Failed events require operational inspection/redrive; no administrative remediation UI or alert integration is included. Refund/reversal policy remains CFR-005, and no tax, inventory valuation, recognition, or rounding policy was added.
- **PR #29:** merged 2026-11-13; merge commit `57c57373ddd1f0b335369ab65717a036695bddcb`.
- **PR #33:** merged 2026-11-13; merge commit `c6f320a19ec5a19f6a9ec59e0fe1cf8059a71a8c`.
- **PR #34:** merged 2026-11-13; merge commit/current main `fc8647d6a0b802c3abf7bc3e93dc0b69e1ca9e8d`.
- **PR #35:** open and ready for review; implementation commit `34dd55a9d274003ac0f5880f5e5009a695fb7ac9`; all GitHub checks pass at the last verified head. No independent review is present; merge pending approval. GitHub reports `main` is not branch-protected.
- **Worktree boundary:** the original protected checkout `H:\Projects\pos` was not changed; implementation is isolated in the dedicated continuation worktree. No merge or deployment occurred in this continuation.

**Single source-of-truth path:** `H:\Projects\pos\FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md`.

## 14. CFR-005 policy-independent safety foundation — 2026-11-14

This update supersedes §13's CFR-004 Git state and CFR-005 “NOT STARTED” status where they differ. Current `origin/main` is `ea23d95a213f46fd8cdc59648b56e8474fc679f2`. PRs #33, #34 and #35 are merged; PR #35 merge commit is `ce65b393fe85fca7ff52e0fc4d88c98861e23044`. PR #37 is open at `decdacfa6cc6d2052d2a64ed3e4ed3dd815fd166`; its reported checks passed, but it has no independent review decision. This CFR-005 change is based on current main and is not included in PR #37.

### CFR-005 status

- **Engineering delivery:** PARTIAL — policy-independent unsafe order-payment transitions are now rejected. The existing CFR-004 paid-order outbox remains the stable accounting-intent mechanism; this change adds no parallel event table or migration.
- **Financial acceptance:** OPEN — no order-level provider capture/refund/chargeback contract or approved reversal policy exists. Do not mark CFR-005 VERIFIED FIXED.
- **Implementation:** `CanTransitionOrderPaymentStatus` rejects a generic `refunded` assertion, any downgrade from `paid`, reactivation from `refunded`, and marking a cancelled order paid. Repository payment updates now lock the tenant/restaurant-scoped order, validate the current state, update payment state and enqueue the existing paid-order intent in one transaction. Order cancellation rechecks payment state under the same row lock and refuses cancellation of paid/refunded orders before stock reservations can be released. The use case also reports the unsupported operation before repository execution; the repository lock is the race-safe gate. Existing paid-order event ID remains `paid-order:{tenant}:{restaurant}:{order}`.
- **Inventory boundary:** no refund implies a physical return. Stock release for cancelled orders remains limited to the existing unpaid/unfulfilled rules; paid cancellation cannot release the reservation without an implemented refund and linked reversal workflow.
- **Tests added/updated:** domain transition matrix; PostgreSQL tests for unsupported refund and paid-state downgrade rejection, cancelled-order payment rejection, stable single paid intent, foreign-tenant and sibling-restaurant denial, and concurrent pay/cancel with exactly one valid winner; reservation assertion for blocked paid cancellation; handler SQL mocks updated for the new locked state read.
- **Disposable PostgreSQL evidence:** fresh PostgreSQL 17 database `pos_cfr005_acceptance_test`, server/database timezone `Africa/Cairo`; all 138 active migration records applied from a clean database by `pos_migrator`. Focused PostgreSQL transition, scope, reservation and race tests passed. The race test ran eight pay-versus-cancel pairs and verified a single safe winner and matching outbox state. No migration was introduced by this change.
- **Regression verification:** `go test -p 1 ./... -count=1 -skip '^TestProductImagesAreContainedRasterAssets$'` passed. `go test -tags=integration -p 1 ./... -count=1 -skip '^TestProductImagesAreContainedRasterAssets$'` passed with both `POS_TEST_DATABASE_URL` and the separate `POS_ACCOUNTING_DATABASE_URL` configured for the same disposable database. The excluded test could not create a Windows symlink because the host lacks the required privilege. An unfiltered run also confirmed this environment constraint. `go vet ./...`, `go build ./...`, and `git diff --check` passed. Package serialization avoids cross-package collisions in tests sharing one disposable database.
- **Test environment note:** an initial integration attempt against the pre-existing disposable cluster database migrated as `postgres` failed partial-cancellation tests because the migration's `pos_migrator` ownership guard was not represented. Those tests pass on the clean database migrated as `pos_migrator`; no application change was made for that environment mismatch. No production or development database was used.
- **Not implemented / remaining acceptance:** no order-provider adapter, durable refund intent, callback event identity/deduplication, out-of-order callback handling, provider-timeout recovery, provider-success/local-failure recovery, refund retry worker, refund concurrency control, chargeback support, or linked accounting reversal exists for restaurant orders. Paymob refund behavior remains subscription-payment functionality and is not reused for restaurant orders. No refund allocation, revenue-recognition, tax, rounding, or closed-period policy was inferred. Next implementation requires an owner/payment-provider contract and accountant-approved reversal/date/period rules; independent idempotency and failure-recovery work can proceed only once a supported provider event contract is defined.
- **PR / commit:** PR [#39](https://github.com/Yousef-Abdallah-1711/pos/pull/39) is open from `fix/cfr005-payment-refund-foundations`, implementation commit `30a7cc3dfc34c0ebf3972de81be1b82db31c63f6`. Independent approval and repository CI are pending. No merge or deployment has occurred.

**CFR-005 next executable step:** obtain and record the restaurant-order provider event contract (capture/settlement/refund IDs and delivery guarantees) and accountant decisions for reversal entry/date and closed periods. Then implement refund intent plus callback idempotency and reconciliation as a separate reviewed PR; retain these transition guards until a tested replacement workflow is deployed.
