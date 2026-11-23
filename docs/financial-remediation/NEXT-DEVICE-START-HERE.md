# Financial remediation: start here

## Current verified handoff — 2026-11-23

Fetch `origin/main` before work. Verified integration snapshot at this handoff update: `4ab49cd0b796df778ffa5257acf9b8017eaeeabb`; the tracker addendum in this same documentation PR may advance `main` again, so use `git rev-parse origin/main` for the final SHA. Financial PRs #104, #108, #109, #110, #112, #114, #117, #121 and #122 are merged. PR #110 was merged at `4ab49cd0b796df778ffa5257acf9b8017eaeeabb`; its protected checks and PostgreSQL acceptance passed, while its then-running benchmark is not performance approval.

The 21 issue folders and Spec Kit drafts are in `main`, but they are not 21 verified or accepted fixes. Start by reconciling the current tracker against original issues and evidence; worktree recovery remains incomplete. Keep local database/browser artifacts private until reviewed. Accounting decisions OPEN-01, EARN-01, CLOSE-01 and related policy items remain pending, CFR-010 is open, and production is NO-GO. Do not deploy or change production data.

Read [current consolidated tracker](FINANCIAL_TRACKING_PROGRESS.md) first, then fetch and inspect the live branch/PR state. Preserved snapshots and their caveats are listed in its latest verified section.

---


This is a **partial recovery draft**. Repository: https://github.com/Yousef-Abdallah-1711/pos. Main at the latest 2026-11-23 snapshot was `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`; always fetch and verify. Recovery is [PR #117](https://github.com/Yousef-Abdallah-1711/pos/pull/117), branch `codex/financial-recovery-specs-20261122`.

Read [program overview](../../specs/financial-remediation/program-overview.md), [register](../../specs/financial-remediation/issue-register.md), [roadmap](../../specs/financial-remediation/roadmap.md), and [CFR-010](../../specs/financial-remediation/issues/CFR-010/continuation.md). The snapshot of the original tracker and the root current tracker preserve historical source text. These drafts are not a completed issue inventory or approved accounting policy.

```sh
git -c core.longpaths=true clone --branch codex/financial-recovery-specs-20261122 https://github.com/Yousef-Abdallah-1711/pos.git pos
cd pos
git fetch origin main
git rev-parse HEAD
```

For this handoff branch, fetch the branch named in its documentation PR; do not assume it has merged. Use Go version from `backend/go.mod`, Node and pnpm from repository configuration, disposable PostgreSQL and Redis, and the setup in `.github/workflows/financial-acceptance.yml`. That workflow lists environment variable names and migration/acceptance commands; never copy production credentials. The documented Windows launch was rejected by execution policy, so provision an explicitly permitted local or isolated CI environment.

Immediate task: continue the individual issue audit at CFR-003 and execute the CFR-002 current-SHA disposable PostgreSQL plan. Inspect [worktree inventory](../../specs/financial-remediation/worktree-recovery.md) and recover/verify missing evidence on remote branches. Spec Kit was initialized; see [setup and limitations](../../specs/financial-remediation/spec-kit-setup.md). The 21 issue documents remain incomplete. Do not merge unreviewed financial PRs, invent accounting approvals, deploy, or modify production data.

On 2026-11-23, PRs #108, #109, #110, #112 and #114 were updated with main and need new-head CI and independent review. At that time the only listed collaborator was the PR author. PR #117 remains draft. OPEN-01, EARN-01 and CLOSE-01 remain pending; CFR-010 is open and financial production is NO-GO.

PR #109 was subsequently updated to `244c44c8f7c53fc81adbd83453edf2875acce9fb` to classify overlapping closed periods and unrecognized period statuses as ambiguous. Its disposable PostgreSQL CI job passed; inspect individual test logs when the workflow finishes. The local attempt could not connect to a test database. Main advanced since the candidate, so refresh it and rerun checks. The original #109 worktree remains dirty and must not be reset. Check the PR's current checks and independent review before considering a merge.

CFR-004, CFR-005 and CFR-009 now contain source-traced merged safeguards and issue-specific remaining tasks; they are partial, not verified complete. Selected atomic posting, payment-transition and seven payroll PostgreSQL scenarios passed in [#109's disposable job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) at `244c44c8f7c53fc81adbd83453edf2875acce9fb`. The CFR-005 payment-state unit test also passed eight cases at recovery checkout `8667c036a5c2c3e12b8a7d953d25fd5553f15bf9`; provider refund and final-main acceptance remain open.

CFR-011's migration 187 is merged and both constructed supplier entry-type persistence tests passed in the same PostgreSQL job. Current source search found no production caller for those two auto-post hooks; supplier lifecycle support, accounting mapping and end-to-end acceptance remain pending.

CFR-012 remains policy pending: main's Cash Flow classification uses account types and code prefixes. A #114 synthetic taxed-sale fixture observed cash GL +55/-55 against the existing formula's +50/-50; its void-date tests passed in [candidate PostgreSQL CI](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45077468898/job/137377801133). That result is not an approved cash-flow methodology.

PR #110 proposes a production reconciliation query change alongside its benchmark. Its recorded warm large-scope latency increased between separate runs despite a faster line-query plan and eliminated temporary I/O; there is no approved performance budget. Do not treat passing benchmark CI as operational acceptance.

Main was merged into the five PR branches on 2026-11-23. Heads at that snapshot were #108 `d7712f4f84189a848f6bbdf43513fef179c73536`, #109 `bbd23164f81258190e85f1377cb261568ab631d9`, #110 `250028d0cd47fa90c34b8904de3dbefd8f0ac215`, #112 `8ac12b956b4da4503acc44b0db42c3af553d7c68`, and #114 `641e6b515cfb2419b212d771feab680dc4de36b9`. New CI was running; all had no recorded GitHub reviews. Check live heads and required checks again. The previous #109 candidate's three fiscal-period regression tests passed in disposable PostgreSQL.

The old integrated CFR-010 worktree's committed head `98aa00951756a2c89e15678888eeb67b2d66dd5e` is preserved at remote `codex/preserve-cfr010-integrated-20261123`. Its untracked evidence and browser outputs were not pushed and remain to be classified for sensitive data.


## Consolidation update — 2026-11-23

PR #117 merged into main as `a461d2104ee07002789bf0be32bb4ad2536189ad`. Clone main to obtain the draft issue packages and Spec Kit; the earlier branch-specific clone instructions above are historical. A clean clone of that commit contained all 21 issue folders with spec, plan and tasks, the register, constitution, and this handoff. File presence does not establish issue completion or financial acceptance.

The owner authorized technical review followed by merge when GitHub protection and CI pass, without requiring a separate recorded independent approval. Recheck each PR head, current main, required checks and its financial risks before merging. PRs #108, #109, #110, #112 and #114 were updated with that main and their new checks are pending. Accounting decisions OPEN-01, EARN-01 and CLOSE-01 remain pending; CFR-010 remains open and production is NO-GO.

Historical preservation branches include `codex/preserve-phase4-financial-proposals-20261123` at `1595bf54b4d5782e6b87993c08b63bece36ee38f`, `codex/preserve-phase00-integration-20261123` at `8062e54a559a85dc889ed37f78593dbd8890ccb9`, `codex/preserve-cfr004-dirty-20261123` at `5c6a987ba2ead8c96c8115eb966d50ffa5cf399c`, and `codex/preserve-cfr010-synthetic-reports-20261123` at `2a152faaa01b2501de54dad90894eaa74e62227b`. These are unmerged snapshots, not approved implementations. The local worktree inventory and sensitive evidence review are incomplete; not all work is proven recoverable from GitHub.


Consolidated 21-issue snapshots, PRs, policy gates and recovery checklist: [Financial Tracking Progress](FINANCIAL_TRACKING_PROGRESS.md). This is a dated recovery draft; refresh live status before action.
