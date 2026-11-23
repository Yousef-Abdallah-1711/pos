# Financial remediation: start here

This is a **partial recovery draft**. Repository: https://github.com/Yousef-Abdallah-1711/pos. Main at the 2026-11-23 snapshot was `afda1f91c67922a92f66507adc208e6e1943d631`; always fetch and verify. Recovery is [PR #117](https://github.com/Yousef-Abdallah-1711/pos/pull/117), branch `codex/financial-recovery-specs-20261122`.

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

PR #109 was subsequently updated to `244c44c8f7c53fc81adbd83453edf2875acce9fb` to classify overlapping closed periods and unrecognized period statuses as ambiguous. New PostgreSQL regression tests require exact-head CI: the local attempt could not connect to a test database. The original #109 worktree remains dirty and must not be reset. Check the PR's current checks and independent review before considering a merge.
