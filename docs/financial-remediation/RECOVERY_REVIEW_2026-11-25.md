# Financial Recovery Branch and Worktree Review — 2026-11-25

## Scope and verified baseline

Repository: Yousef-Abdallah-1711/pos
Primary checkout: H:/Projects/pos
Verified origin/main: fd3636370f853150e0ddf9b676f21eb0084ba9bf
Open GitHub PRs at review time: 0
Private preservation repository: Yousef-Abdallah-1711/pos-phase5-private-acceptance (private)

This report records a continuation review and does not claim the entire recovery assignment is complete. It supplements the [financial tracker](FINANCIAL_TRACKING_PROGRESS.md) and existing [worktree recovery manifest](../../specs/financial-remediation/worktree-recovery.md).

## Current-main verification

- Dedicated synthetic PostgreSQL/browser acceptance: [run 45509343366](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45509343366), conclusion success, exact head fd3636370f853150e0ddf9b676f21eb0084ba9bf. Journal/GL reconciliation, API posting, independent Balance Sheet, authenticated EN/AR desktop/mobile journeys, and payment/auth/order/reporting scenarios passed. Raw evidence upload was skipped by workflow policy.
- Financial benchmark: [run 45509925241](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45509925241), conclusion success, exact same head. Synthetic 2-tenant dataset: 10,102 journals, 50,229 journal lines, 50,229 GL rows, 2,000 sales orders, database 75.37 MB. PostgreSQL 16.4 on Linux/amd64; 4 reported CPUs and 16.46 GB reported memory; 5 warmups and 30 measured samples.
- Large-scope reconciliation p50/p95/max: 182.568/184.974/185.149 ms. Four-reader concurrency: 304.714/340.613/345.769 ms. Reads during posting: 299.302/315.089/343.535 ms with 11 concurrent posts. All reported sample errors were zero.
- Current journal_lines_to_gl plan execution: 149.422 ms, with zero temporary blocks read/written. The same-run ledger-page index A/B changed plan execution 0.703 to 0.638 ms; p50 21.109 to 20.980 ms, p95 21.721 to 22.004 ms, and max 21.913 to 22.508 ms. Both plans reported no temporary blocks.
- Historical #110 comparison in audit-reports/financial-cfr010-remediation-2026-11-21/PERFORMANCE_RESULTS.md: line-to-GL execution 160.660 to 146.178 ms and temporary I/O 3,544 reads/713 writes to zero; full large-scope p50 increased 174.044 to 178.529 ms. No approved performance budget exists; performance acceptance remains NOT APPROVED.
- Local command already executed on current main in this session: go test ./internal/repository ./internal/usecase ./internal/handler/http passed all three packages. Dashboard pnpm exec vitest run src/lib/config.test.ts passed 4 tests. These results do not replace the disposable PostgreSQL acceptance suite.
- Financial production remains NO-GO. CFR-010 remains OPEN. OPEN-01, EARN-01, CLOSE-01 and other applicable accounting decisions remain pending. No production records or deployment were changed.

## Branch disposition from this pass

| Preservation branch | Exact head | Technical disposition |
|---|---|---|
| preserve/ws02-dirty-snapshot-20261124 | 922c746247096cf412cd7fab583688b953efa4eb | Private snapshot preserved. Do not merge wholesale: diffs include weakened request decoding, internal error disclosure, and reduced tenant-scope/transaction controls. Some isolation tests already exist on main. Remaining individual test-only differences still require review. |
| preserve/financial-phase00-working-20261124 | 24448e3427ba577a66dae4da4369839dae6e9292 | Private snapshot preserved. Core Balance Sheet integrity diagnostics/UI and several tests are already on main. The extra historical execution package/evidence was not cleared for public publication in this pass; it contains dated synthetic outputs and needs full sensitivity/freshness review. |
| preserve/phase00-final-acceptance-scripts-20261124 | cad7d78d0c5b41b996c77fb45270e0e1634765ea | Private snapshot preserved. Core implementation overlaps main, but 11 changed code/test blobs remain different and were not individually approved for recovery in this pass. Acceptance scripts and associated evidence remain unreviewed for publication. |
| preserve/phase00-pr-handoff-20261124 | 8c49b398cdb7b9781559bb489b277ecfa314f841 | The 19 changed paths examined in this comparison have identical blobs on main; historical handoff remains preserved. |
| preserve/cfr004-ops-recovery-unmerged-20261124 | 86c828f1ca43c51cd049c4c875d9c93fe7ebc4de | Superseded by current migration 188 implementation; branch migration 185 is not to be reintroduced. |
| preserve/financial-master-audit-20261124 | 01562e8d63ef38c5081936fa1b1a93aa8107ff44 | Dated audit documents remain private pending line-by-line freshness and sensitive-evidence review. Related executive assessment is already represented in current public tracker/audit records. |
| preserve/financial-current-audit-20261125 | 2841470d8f721ce92a38ce2b777cb1abbb2dad8d | Preserves the previously untracked 2026-11-16 local financial audit in the private repository. It contains dated operational/staging observations and was not cleared for public publication or revalidated against current main. Restore from private-preservation branch preserve/financial-current-audit-20261125. |
| preserve/financial-forensic-audit-20261124 | e38b6d7be2520bd97e782b90963103b3ce031dee | The 82-file forensic/architecture package is preserved privately. Its findings have not been fully mapped to current source and tests in this pass. Do not publish the package wholesale. |
| preserve/cfr010-realworld-validation-20261124 | 4dd766e3981f921600b05305ec2e05216e231ebb (advanced from reviewed source snapshot c76bd74c29ebffe28e1d7e6dbde94b18f964277a) | Current source implementations for the principal diagnostics/benchmark paths match main. Two previously untracked, credential-free synthetic browser harness files were preserved on this private branch after node --check passed. Dated reports still need freshness review; raw local JSONL, browser artifacts, and screenshots remain outside public GitHub. |
| preserve/cfr010-earnings-period-scope-proposal-20261124 | 8b4831576a664a1ee82df5ac5ed59a2b904f9726 | Diagnostic files are present on main; the proposal loosens closed/overlapping-period handling and must remain unmerged pending a correct policy-neutral design. |
| preserve/localonly-cfr010-performance-evidence-20261124 | 79b6409f3e385a0e654a82c59be18a8279749acd | Compared source/workflow paths match current main. Historical performance report remains evidence only and is already summarized in the tracker. |
| preserve/localonly-financial-benchmark-ci-path-20261124 | 2dca4923c4f6fc52cdaf3a9f191044f3ad4ff2e4 | Benchmark workflow and instrumentation candidate require ancestry-aware review; current main already has the workflow fixture-path correction and benchmark implementation. No new merge identified. |
| preserve/phase4-closure-docs-20261124 | a0c16f0025b84621ef70da4da4470865f82cb8b7 | Historical closure documents remain private; two examined claims require freshness review. |
| preserve/phase4-evidence-update-20261124 | 1f0ad49478b6a1486855cc6d6a5c26c8f59c3946 | One historical evidence report remains private pending current-SHA reconciliation. |
| preserve/phase4b-build-recovery-ts-20261120 | 42988c48b0e7a058f5629fdd92d1daaa7cbebec1 | Dependency-lock, package patch and build files differ from main; full dependency/security/build review remains before any integration. |
| preserve/phase4c-hosted-acceptance-20261124 | 18244e78c427c95b70c9ac2725a6b4051c2ce20f | Includes production CI/deployment workflow and operational scripts; not eligible for blind merge. Review separately after financial recovery, with deployment gates unchanged. |
| preserve/dashboard-test-fixes-20261124 | e43d6a1af5fb6dc8938b2c4fbd00df51fef36b74 | One test-file change is a redundant cleanup only; the 4-test selected Vitest run passes on main. No PR warranted. |

This is a targeted review of the named financial/security branches, not a complete review of all 43 private branches or all 125 worktrees.

## Dirty worktree reconciliation

Current inventory: 125 registered worktrees, 33 with any dirty status, 10 with tracked modifications. The primary checkout has pre-existing untracked evidence and consolidation artifacts; none were deleted or cleaned.

Tracked worktree comparisons against current origin/main found:
- Generated Next.js config/next-env.d.ts changes in two worktrees. They are build-generated and not recovered as source.
- registry.json generated timestamp-only change. Not recovered.
- One dashboard config test cleanup; selected 4-test Vitest suite passes on main and this cleanup adds no behavior.
- A stale H:pbi checkout with numerous absent or obsolete files. Its local tree is incomplete relative to main and is not a merge candidate.
- Financial diagnostics, earnings tests, the current financial report PostgreSQL test, and principal CFR-010 benchmark/code paths that were dirty relative to older worktree heads hash-identically to current main.

Untracked financial and browser outputs require content-sensitive treatment. Local PostgreSQL data directories, Playwright traces/screenshots, JSONL acceptance results, local test results, and synthetic PDFs were left in place and excluded from public GitHub in this pass. The CFR-010 screenshot runner and localhost-only Compose override were committed to private preservation branch preserve/cfr010-realworld-validation-20261124 at 4dd766e3981f921600b05305ec2e05216e231ebb; generated outputs remain local. The previously untracked SERVIO_CURRENT_FINANCIAL_AUDIT_2026-11-16.md was preserved privately at 2841470d8f721ce92a38ce2b777cb1abbb2dad8d because it contains dated operational/staging detail; it is not public main. This report still does not claim that every untracked meaningful source/evidence file has been individually classified.

## Remaining work and blockers

1. Complete file-level review of every remaining preservation branch and each meaningful changed/untracked file in all dirty worktrees.
2. Reconcile the 82-file financial forensic package and dated Phase 00 / real-world reports against current code, current tracker, and evidence freshness; publish only a sanitized, non-stale subset if useful.
3. Review the remaining local-only worktree source/evidence and preserve legitimate items on verified remote branches. Avoid publishing local databases, raw browser/auth traces, credentials, customer data, or raw financial records.
4. Repeat relevant CI after any recovered source changes, then open normal PRs and merge only after required checks/conversation resolution pass.
5. Obtain finance-owner decisions and original operating evidence for OPEN-01, EARN-01 and CLOSE-01 (plus other applicable decision items), approve a workload-specific performance budget, and complete independent accounting/operational acceptance.
6. Keep financial production NO-GO until all required accounting and operational gates have passed.

## Latest consolidation verification — 2026-11-25

- **Main and PRs:** documentation PR [#132](https://github.com/Yousef-Abdallah-1711/pos/pull/132) merged with merge commit `13f1aab2c099eb0584ba6c2d78c7802ac05cb327`. Current GitHub query reports no open PRs. The primary checkout `H:\Projects\pos` is on `main` at the same SHA with tracked changes clean; existing untracked local files remain untouched.
- **Clean clone:** cloned `origin/main` at `13f1aab2c099eb0584ba6c2d78c7802ac05cb327`. Verified 21 issue folders and `spec.md`, `plan.md`, and `tasks.md` for every folder, plus the tracker, handoff, recovery review, and issue registers. Result: PASS. This verifies documentation availability, not issue acceptance.
- **CFR-001 probe disposition:** a historical local probe attempted a 100.00/90.00 auto-post and also wrote journal entries without cleanup. It was not run. Current main has `validateAutoPostingLines`, which compares integer cents and rejects both debit-heavy and credit-heavy 10.00 differences before repository writes; focused Go tests passed. Historical probe and warning are privately preserved at `preserve/financial-forensic-probes-20261125`, commit `6c9de4c619d7356e7004d76f42caef7140df2e54`; treat it as forensic evidence, not a runnable test.
- **Private evidence preservation:** credential-pattern scanning found an email in raw authenticated acceptance JSONL, so that file and screenshots were excluded. Reviewed synthetic benchmark/browser summary evidence was stored privately at `preserve/cfr010-realworld-evidence-20261125`, commit `181430d6eec9e9d8e28c61ad2f64c44c564d81ed`. Two screenshot-runner source files remain on `preserve/cfr010-realworld-validation-20261124`, commit `4dd766e3981f921600b05305ec2e05216e231ebb`.
- **Recovery remains partial:** inventory still has 125 worktrees, 33 with dirty status, and 10 with tracked modifications; generated caches and raw browser/financial artifacts remain local unless separately verified. No local-only commits outside public/private remote refs were found. Meaningful files in every dirty worktree and all preservation branches have not yet received complete content and sensitivity review. Financial production remains NO-GO; CFR-010 and accounting approvals remain open.

## Continuation review after PR #134 — 2026-11-25

- **Baseline and PR state:** PR #134 merged normally; merge SHA and current `origin/main` are `133260db66b1143c254d79476daa7650b41a33f0`. GitHub returned no open PRs. The primary H checkout was verified at this SHA before starting this update branch; no untracked files were staged or removed.
- **CI:** main run [45609799164](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45609799164) passed backend, financial PostgreSQL, frontend and security/config jobs. The financial PostgreSQL job created restricted disposable roles, applied migrations, reconciled journal/GL, exercised authenticated order-to-ledger posting, independently checked the Balance Sheet, completed authenticated EN/AR desktop/mobile journeys, and destroyed disposable services. Auth-browser-acceptance was skipped; application-images was still running at the last poll. Do not represent either as passed.
- **Fresh clone:** `H:\Projects\pos-clean-clone-20261125-final` is a shallow clone at exact SHA `133260db66b1143c254d79476daa7650b41a33f0`. Verified 21 issue folders, each with `spec.md`, `plan.md`, and `tasks.md`, plus `.specify`, tracker, and handoff. PASS for documentation/package availability; it does not verify all evidence or accounting acceptance.
- **Local tests:** `go test ./internal/usecase ./internal/handler/http` passed. Dashboard command `pnpm exec vitest run src/components/financial-reports/BalanceSheet.test.tsx src/lib/financial-integrity.test.ts` passed 2 files / 3 tests. Local Docker engine and `psql` were unavailable, so no independent local PostgreSQL run was performed.
- **Phase 4C branch review:** private head `18244e78c427c95b70c9ac2725a6b4051c2ce20f`, remote branch `preserve/phase4c-hosted-acceptance-20261124`. Its changes include exact candidate-SHA checkout, unique Compose project and teardown checks, log-redaction helper, authenticated synthetic tenant/browser acceptance, and pinned-tool SBOM/image scans. The branch also carries dated Phase 4 closure documents and a patched `next-intl@4.15.1` proof-of-concept whose watcher replacement throws when `experimental.extract` runs. This dependency patch is not production-ready. No branch-wide merge is safe. Candidate controls have not yet passed a full independent test/security review; keep private pending further item-level review.
- **Open review items:** WS02 ancestry-aware diff contains 35 differing file blobs beyond current `main`; some direct source changes are known regressions, while every residual test/documentation blob has not been individually disposed. Worktree inventory remains 125 total / 33 dirty / 10 tracked-modified as last counted. This pass does not complete the prompt's all-branches/all-worktrees closure criteria.
- **Safety:** no production data changes or deployments; all accounting decisions remain pending, CFR-010 remains open, and financial production remains NO-GO.

## Latest review batch — 2026-11-25 (after PR #133)

- **Baseline:** current origin/main and H:\Projects\pos are both 9a258399e671b5af7d5ff684bbf817684b2e061f; tracked files are clean. PR #133 merged at this SHA. GitHub currently reports zero open PRs.
- **Clean clone:** an independent shallow clone was advanced to exact SHA 9a258399e671b5af7d5ff684bbf817684b2e061f. Verified 21 issue directories; each has spec.md, plan.md, and tasks.md; tracker, handoff, recovery review, and issue registers exist. PASS for artifact availability only.
- **Worktree/commit scan:** 125 registered worktrees; 33 have any dirty status and 10 have tracked modifications. git rev-list --branches --not --remotes=origin --remotes=private-preservation --count returned 0. This does not account for dirty files; those still need content review.
- **preserve/financial-backdated-ledger-20261120 (a5111b7e...):** not safe to integrate wholesale. Its workflow changes remove financial/browser and production CI checks; test diffs delete coverage. Its P&L/Cash Flow query changes omit void reversal movements and introduce cash-flow classification behavior while relevant policy remains pending. Keep private; no integration.
- **preserve/cfr004-ops-recovery-unmerged-20261124 (86c828f1...):** its migration 185 is superseded by current migration 188. The old main.go also reverts signing-secret validation, separate accounting DB wiring, selection-proof enforcement, protected POS order routing, and no-store API response middleware. Current main already has the redrive handler/repository path. Do not merge this snapshot.
- **preserve/cfr010-earnings-period-scope-proposal-20261124 (8b483157...):** removes ambiguous-period regression tests and treats an overlapping period as open. Current main keeps the ambiguous-scope guard; proposal remains private and unmerged pending a policy-neutral correction.
- **preserve/phase4b-build-recovery-ts-20261120 (42988c48...):** dependency patch replaces next-intl's watcher with a function that throws when experimental.extract attempts file watching. Keep as a private proof-of-concept; not a valid build fix.
- **Already represented in main:** all four changed paths on preserve/localonly-tenant-selector-and-earnings-20261124 match main byte-for-byte; all five code/config changes on codex/preserve-performance-tools-original-20261124 also match main. No new public PR warranted.
- **Forensic CFR-010:** private report preserve/cfr010-forensic-inventory-20261124 says the original operating opening-balance root cause remains unestablished without source statements and original journal/GL evidence. It does not support accounting acceptance.
- **Private recovery commits:** reviewed synthetic benchmark/browser summaries at preserve/cfr010-realworld-evidence-20261125 (181430d6eec9e9d8e28c61ad2f64c44c564d81ed); historical CFR-001 probe at preserve/financial-forensic-probes-20261125 (6c9de4c619d7356e7004d76f42caef7140df2e54); reviewed synthetic Phase 00 evidence at preserve/financial-phase00-evidence-20261125 (5b4040940a13c1697119003eb51747286da0cbcb). Raw acceptance JSONL with an email and screenshots were excluded from upload. Each private evidence set is historical/synthetic and not production acceptance.
- **CFR-001 check:** current main integer-cent auto-post validator rejects both debit-heavy and credit-heavy 10.00 imbalance before writes. Focused command passed: go test ./internal/usecase -run 'TestValidateAutoPostingLines|TestAutoPostingRejectsUnbalancedEntriesBeforeJournalOrLedgerWrites' -count=1.
- **Remaining:** review the other private branches and all meaningful dirty/untracked worktree files individually; preserve only reviewed material remotely. CFR-010 operating root cause, OPEN-01/EARN-01/CLOSE-01 and other accounting approvals remain open. Financial production stays NO-GO.
