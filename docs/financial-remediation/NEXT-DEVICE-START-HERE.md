# Financial remediation: start here

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
