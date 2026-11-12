# Candidate 1 — frontend test baseline

## Identity

- Branch: `fix/dashboard-test-typecheck-baseline`
- Base: `a6a2c6700a5fe883256deb1e5a3893108bd9865c`
- Commit: `5528c6f4bac5c6575c3a185a45ad81eaffa36e6e`
- Commit message: `test(dashboard): repair baseline test typecheck failures`
- Worktree: `H:\pos-pr-candidate-baseline`

## Scope review

Exactly two test files changed from the base:

1. `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableList.test.tsx` — removes duplicate `onOpenPOS` JSX properties while preserving the test behavior.
2. `frontend/apps/dashboard/src/lib/config.test.ts` — replaces readonly `NODE_ENV` assignment with Vitest environment stubbing and cleanup.

No runtime application code, compiler configuration, or test exclusions are included. Diff summary: 2 files, 5 insertions, 6 deletions. Full diff and commit manifest are recorded under `evidence/git/`.

## Verification

Recorded PASS on this candidate: full `typecheck:test`; focused tests (6); dashboard Vitest (46 files / 653 tests); production typecheck; targeted ESLint; production build. The combined integration run repeats full test typecheck, dashboard Vitest (48 files / 656 tests), production typecheck, targeted lint, and build after applying both commits.

## Readiness

**Local review ready.** Candidate is isolated and based on current main. No push or PR creation has occurred. Before opening a PR, review the two-file diff and attach the verification evidence index. No standalone raw frontend command transcripts were retained in the candidate worktree; this report records the execution results from the ongoing verification session.
