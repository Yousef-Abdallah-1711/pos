# Financial remediation: start here

This is a **partial recovery draft**. Repository: https://github.com/Yousef-Abdallah-1711/pos. Main at the snapshot was `966a5004fc594d67f8ef1850909bd7f8e963ecf5`; always fetch and verify.

Read [program overview](../../specs/financial-remediation/program-overview.md), [register](../../specs/financial-remediation/issue-register.md), [roadmap](../../specs/financial-remediation/roadmap.md), and [CFR-010](../../specs/financial-remediation/issues/CFR-010/continuation.md). The snapshot of the original tracker and the root current tracker preserve historical source text. These drafts are not a completed issue inventory or approved accounting policy.

```sh
git clone https://github.com/Yousef-Abdallah-1711/pos.git
cd pos
git fetch origin main
gh pr checkout 114
git rev-parse HEAD
```

For this handoff branch, fetch the branch named in its documentation PR; do not assume it has merged. Use Go version from `backend/go.mod`, Node and pnpm from repository configuration, disposable PostgreSQL and Redis, and the setup in `.github/workflows/financial-acceptance.yml`. That workflow lists environment variable names and migration/acceptance commands; never copy production credentials. The documented Windows launch was rejected by execution policy, so provision an explicitly permitted local or isolated CI environment.

Immediate task: inspect [worktree inventory](../../specs/financial-remediation/worktree-recovery.md) and recover/verify missing evidence on remote branches. Then complete the issue-by-issue verification and official Spec Kit initialization. Do not merge unreviewed financial PRs, invent accounting approvals, deploy, or modify production data.
