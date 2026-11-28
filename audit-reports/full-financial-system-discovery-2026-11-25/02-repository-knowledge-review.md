# 02 - Repository knowledge review

Baseline: git SHA 133260db66b1143c254d79476daa7650b41a33f0 (master, 939 commits). All statements about remediation in historical documents are labelled DOCUMENTED CLAIM - NOT REVALIDATED unless a code path is cited. Documents were read in full for the financial plan, policy package and specs; large reports (audit-reports/**) were sampled by index and the financial-* folders' headline files. Coverage gaps are listed at the end.

## 1. Document inventory

| Document | Size / depth read | What it claims |
|---|---|---|
| `AGENTS.md` | full (527 B) | Agents must read `CODING_STANDARDS.md` before any work. No financial content. |
| `CODING_STANDARDS.md` | NOT read in full (72 KB; only referenced) | Repository's code-quality/security/testing policy (per AGENTS.md). Gap: financial-specific rules not extracted. |
| `PROJECT_KNOWLEDGE_BASE.md` | headings + API/accounting sections (40 KB) | Go + Next.js multi-tenant POS SaaS; RBAC modules 1-20 (Accounting module 15, Financial Reports + Tax module 17, Suppliers 16); accounting API route list; migration history. Predates the CFR program (not authoritative on current behaviour). |
| `README.md` | skim (34 KB) | Product overview/run instructions. Not financial-specific. |
| `FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md` (143 KB) | sections 1-24 read; issue sections for all 21 findings | Canonical tracker: 21 findings (15 CFR + 6 FIN), verdict NO-GO, per-issue root cause, acceptance, and dated continuation sections (13-24). Sections are append-only snapshots; later sections supersede earlier ones. Two sections are both numbered 23. |
| `PHASE4_ACCOUNTING_POLICY_DECISION_PACKAGE.md` | full | Ten pending decisions: CFR-011 supplier contract, OPEN-01, EARN-01, CLOSE-01, TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01, CF-01; all DECISION REQUIRED, none approved. Fixture: Assets 5,500c, Liabilities 500c, Equity 0c, residual 5,000c UNRECONCILED. |
| `SERVIO-COMPLETE-DEVELOPMENT-AND-FINANCIAL-HANDOFF.md` | headings only | 2026-11-12/13 handoff: repo structure, tracker integrity, PR/worktree map, preservation, owner/accountant decisions, NO-GO gates. |
| `REMAINING-TASKS-CHECKLIST.md`, `PROJECT_RECONCILIATION_2026-11-04.md`, `SERVIO_PRODUCTION_READINESS.md`, `PLAYWRIGHT_FULL_ACCEPTANCE_REPORT.md`, AWS/GitHub-Actions reports | NOT read (listed only) | Non-financial or ancillary; gap. |
| `specs/financial-remediation/program-overview.md` | full | Provisional recovery draft, not a completed Spec Kit program; 15 CFR + 6 FIN ids; production financial NO-GO. |
| `specs/financial-remediation/roadmap.md` | full | 6-phase dependency roadmap (Recovery, Integrity, Reports/policy, Transactions, Security/perf/exports, Acceptance); current phase "recovery". |
| `specs/financial-remediation/worktree-recovery.md` | table skimmed | Snapshot 2026-11-22 of ~55 local worktrees; counts of dirty entries and commits absent from cached remotes; inventory, not proof of preservation. |
| `specs/financial-remediation/issue-register.json` / `.md` | full | 21 entries (id, title, severity, current_state, source_paths, PRs, classification). Statuses "provisional". Snapshot main `966a5004`, original commit `fa2e301`. |
| `specs/financial-remediation/issues/{CFR-001..015,FIN-016..021}/{spec,plan,tasks}.md` | NOT read individually (dir listing only); CFR-010 has `continuation.md` | Per-issue recovery drafts. Gap. |
| `specs/financial-remediation/original-tracker-2026-11-12.md`, `pr108-ci-investigation.md`, `spec-kit-setup.md` | NOT read | Immutable original tracker and CI investigation. Gap. |
| `docs-organized/**` | directory listing and keyword search only | Product/feature documentation (themes, API, KDS, tables, offline, Paymob, RBAC). Keyword hits for accounting/tax are incidental; no financial policy located. |
| `audit-reports/**` | index listing; read `financial-closure-2026-11-21/CURRENT_WORK_CLOSURE.md` fully and `cfr010-balance-sheet-fix-2026-11-22/*` headline files | See sections 2 and 3. Other financial folders (phase2/3/3a/3b/4, closure-2026-11-20, master-audit-2026-11-14, cfr010-remediation-2026-11-21, financial-benchmark) listed, not read. |
| `.github/workflows/financial-acceptance.yml`, `financial-benchmark.yml` | acceptance: first ~150 lines; benchmark: first 60 lines | See section 4. |

## 2. Timeline of the financial remediation program

All dates and PR numbers below are DOCUMENTED CLAIM - NOT REVALIDATED (taken from the plan, register and closure docs; recent merge commits are visible in `git log`, e.g. #132-#134 at the baseline).

| Date | Event (source) |
|---|---|
| 2026-11-10/11 | Initial forensic audit artifacts; Phase 00 safety candidate (`fix/financial-phase00-safety`, 5064daa): report warning + export guard + posting diagnostics. |
| 2026-11-12 | Assessment baseline `a6a2c67`; tracker created at `fa2e301`; 21 findings (7 P1, 9 P2, 5 P3), verdict NO-GO. Test tally then: 726 Go passed, 2 skipped; 3 DB defects reproduced (plan s1). |
| 2026-11-13 | PR #29 (CFR-001 reject unbalanced), #33 (CFR-004 atomic post), #34 (CFR-002 DB integrity, migration 180). PR #35 durable outbox. |
| 2026-11-14 | CFR-005 foundations (#39), CFR-009 paid payroll (#40), CFR-003 numbering (#41), CFR-010 signed-balance (#42), CFR-011 width (#44, migration 187), reconciliation of #36/#37/#46/#48 into main. |
| 2026-11-15 | Phase 3 acceptance snapshot: no runtime scenarios run (Docker stopped); NO-GO. Migration 188 redrive audit. |
| 2026-11-16 | Phase 4 engineering and policy decision package. |
| 2026-11-20/21 | Financial execution refresh; dated GL (#93, migration 190); authenticated financial browser acceptance (#72); benchmark (#94: 10,091 journals, 50,207 lines); closure doc with 21-finding disposition. Private acceptance runners offline. |
| 2026-11-21/22 | CFR-010 remediation: journal-GL reconciliation (#104), opening/cross-report diagnostics (#108), current-earnings diagnostic (#109), reconciliation benchmark (#110), decision packets (#112), void history/cache refresh (#114). Synthetic 5,000c residual explained as unclosed earnings; original operating discrepancy NOT ESTABLISHED. |
| 2026-11-22/24 | Recovery/specs drafts (#117), final recovery and preservation (#131). |
| 2026-11-25 | Recovery evidence batches (#132-#134); baseline for this audit (133260d). |

Latest documented release status (closure doc): engineering implementation PARTIAL; financial reconciliation UNRECONCILED; accounting decisions PENDING; performance SLO not accepted; production financial NO-GO.

Internal inconsistencies noticed between documents (relevant for later phases):
- Register classifies CFR-001 as FIXED_IN_MAIN_UNVERIFIED while the plan's CFR-001 section still calls it PARTIAL.
- Plan s13 lists CFR-003/005/009/011 as "NOT STARTED", later sections record them merged; sections are snapshots, not rewritten.
- Plan has two sections numbered 23.
- Register states "main observed 2026-11-23 3fa3702"; repo baseline here is 133260d (later).

## 3. Issue status per historical document

Every remediation status in the plan, register, handoff and closure docs is DOCUMENTED CLAIM - NOT REVALIDATED. Where this pass checked current code, evidence is in `21-existing-issue-regression-matrix.csv`. Summary of code-backed checks (static only): plug removal (CFR-001) present; DB-integrity and outbox migrations (180-190) present in `backend/migrations`; paid-payroll 409 guards present; hard-coded 10% VAT, float64 money, void-without-period-check, heuristic cash flow and balance sheet without current earnings are all still present in source.

## 4. Test suites that exist

Counts below are from file listing and `grep` of test names; none were executed in this pass.

| Suite | Location | Covers |
|---|---|---|
| Go unit (usecase) | `backend/internal/usecase/auto_posting_validation_test.go`, `auto_posting_engine_test.go`, `phase6_accounting_cycle_test.go`, `financial_report_usecase_test.go`, `accounting_posting_diagnostic_test.go` | Balance validation, posting, accounting cycle, report math with mocks |
| Go integration (tag `integration`, PostgreSQL) - repository | `backend/internal/repository/`: `journal_integrity_postgres_integration_test.go`, `general_ledger_atomic_posting_integration_test.go`, `cfr003_numbering_integration_test.go`, `journal_gl_reconciliation_integration_test.go`, `opening_balance_diagnostics_integration_test.go`, `accounting_advanced_rls_integration_test.go`, `accounting_writer_role_integration_test.go`, `accounting_recovery_roles_integration_test.go` | CFR-002/003/004, journal-GL reconciliation, opening-balance diagnostic, RLS and restricted-role grants |
| Go integration - usecase | `backend/internal/usecase/`: `accounting_event_worker_*`, `accounting_event_admin_*`, `cfr010_*_integration_test.go` (8 files), `financial_report_postgres_integration_test.go`, `receivables_aging_integration_test.go` | Outbox worker/redrive, Balance Sheet/earnings/cash-flow/void scenarios, report-to-ledger oracle, concurrency |
| Go integration - handler | `backend/internal/handler/http/`: `salary_financial_integrity_integration_test.go`, `cfr010_cash_flow_void_integration_test.go`, `report_export_rls_integration_test.go`, `accounting_advanced_einvoice_http_integration_test.go` | Payroll immutability and double-pay, exports under RLS, e-invoice |
| Go live-API acceptance | `backend/cmd/api/financial_live_*_integration_test.go` (order-to-ledger, accountant journals, AR lifecycle, payroll, reports, identity), `financial_supplier_*_integration_test.go`, `financial_browser_fixture_integration_test.go`, benchmark tests | Real API process with restricted roles; supplier rollback and concurrent payments; fixtures; benchmark |
| Browser (Playwright) | `frontend/apps/dashboard/tests/financial-browser/acceptance.browser.js` with `frontend/apps/dashboard/playwright.financial.config.ts` | 11 tests: accountant balanced transfer and ledger reconciliation; cashier denied; Balance Sheet vs SQL oracle (EN/AR, desktop/mobile); PDF/CSV export; AR/AP registers and aging; AR invoice partial payment; AP invoice partial payment; payroll mark-paid; second-tenant isolation |
| Frontend unit | `frontend/apps/dashboard/src/lib/financial-integrity.test.ts` | Financial integrity helper (content not reviewed) |
| Other Playwright | `frontend/apps/dashboard/tests/phase2a-e2e`, `visual-regression` | Not financial |
| CI | `.github/workflows/financial-acceptance.yml` (ubuntu-24.04, disposable Postgres and Redis, restricted roles via `scripts/db/create_runtime_roles.sql`, runs migrations, `TestReconcileJournalGL_*`, builds real API, live API tests, seeds browser fixtures, runs dashboard + Playwright); `.github/workflows/financial-benchmark.yml` (synthetic journals 1k/10k/50k) | Workflow content beyond the read prefix not reviewed |
| Oracles in output folder | `audit-reports/full-financial-system-discovery-2026-11-25/oracle/deterministic_control.py` (+ expected JSON) | Pre-existing in the output folder; not reviewed by this pass |

Not found: property tests for monetary precision, a tax-policy test suite, an order-refund reversal test, a closed-period void test, tests for journal attachment validation, a year-end close test.

## 5. Coverage gaps in this pass

`CODING_STANDARDS.md`, `README.md`, handoff body, per-issue `spec/plan/tasks`, original tracker, most financial audit-report folders, benchmark JSON, workflow tails, and `docs-organized` were not read in full. These should be sampled in a later phase if a claim depends on them.
