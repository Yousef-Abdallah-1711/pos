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
