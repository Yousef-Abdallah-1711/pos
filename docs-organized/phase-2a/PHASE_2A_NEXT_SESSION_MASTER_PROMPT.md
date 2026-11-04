# Phase 2A Next Session Master Prompt

You are continuing POOS Restaurant POS SaaS Phase 2A. Do not repeat completed discovery. Do not start Phase 2B. Repository: `https://github.com/Yousef-Abdallah-1711/pos`.

## Exact baseline and current publication state

At the latest verified fetch on 2026-11-04, `origin/main` was `e400400b671437df275f686a84aa644e5b5a1811`, including merged Phase 2C recovery report PR #6 and config fix `5d32752`. Fetch and verify again before acting. PR #1 (`codex/ws02-final-security-closure`) remains open at `9cb9a288ff3538c68e11caa2be9f50dd79cf8ee2` with WS02 HR/Accounting/subscription changes; its hosted `runtime-localhost-guardrail` and `url-domain-config-smoke` checks are still red. Before this handover refresh, PR #3 (`codex/phase2a-final-consolidation`) was at `1bfabef1d83c55027250ae27e1ccaa2526f570d1`, with both hosted checks passing and no review approval. A docs-only handover refresh has since been pushed and restarted its checks. Its unique diff remains the explicit QA URL safeguard and handover documents; the Dashboard config fix is already in main. Repository standards require one independent review approval; no other repository collaborator was available at the latest check, so do not self-approve or merge without an eligible reviewer.

Two original worktrees remain dirty and preserved: `H:\Projects\pos` (`codex/phase2a-continuation`, HEAD `ecadb5ab7665df3eb9f69886a9b186c710eff62d`, 260 dirty paths) and `H:\ws02` (`codex/ws02-final-security-closure`, HEAD `9cb9a288ff3538c68e11caa2be9f50dd79cf8ee2`, 60 dirty paths). Five dirty overlaps are listed in `PHASE_2A_FINAL_SESSION_HANDOVER.md`. A checksum-verified private backup is at `C:\Users\yosea\.codex\backups\poos-phase2a-consolidation-20261104-213732`. Never reset, clean, overwrite, stash-destructively, or blindly merge these worktrees. Consolidation worktree: `H:\pos-phase2a-integration`, branch `codex/phase2a-final-consolidation`.

## Architecture and status

See the existing Phase 2A reports under `docs-organized/phase-2a/execution/` and the final handover in this directory. WS01 is recorded complete by existing evidence. WS02 is partial; WS03 Payments deferred; WS04 partial/current checks need full regression; WS05 partial; WS06 CI expansion deferred; WS07/08/09/10 partial. Do not describe topic/dirty changes as part of `main` unless confirmed reachable from fetched `origin/main`.

A security-critical limitation is intentional: Super Admin tenant handoff is disabled (UI disabled, API HTTP 503, no tenant session creation). Never re-enable unrestricted impersonation. Any future support feature needs actor-attributed least-privilege sessions and auditable start/end events.

## Latest local authentication evidence

The root-owned Dashboard logout change in `H:\Projects\pos` explicitly includes the cookie on the cross-origin logout request and checks the response status. It remains uncommitted in the mixed dirty primary worktree and is not included in PR #3. Against the labeled disposable WS05 stack, the invalid-login and logout/session-revocation/protected-route Playwright journeys passed 2/2, and Dashboard typecheck passed. The logout test's successful rerun followed a disposable Dashboard restart and route warm-up after a cold compile stalled with Watchpack `ENOMEM`; do not treat it as cold-start or full-suite verification. The Go toolchain was unavailable in the host/API container, so no fresh Go regression was run.

## First priorities

1. Read repository `CODING_STANDARDS.md` before code changes; the external path in `AGENTS.md` was unavailable on the prior host.
2. Fetch and inspect `main`, PR #1/#3 states, checks, branch rules, and all worktrees. Preserve their current dirtiness.
3. PR #3 checks passed; obtain the required independent review approval before merge. Re-evaluate PR #1 against current main and rerun its checks.
4. Review PR #1 source and tests, then run WS02 HR/Accounting/subscription integration tests against a newly created isolated PostgreSQL database and verify unauthorized DB side effects.
5. Reconcile dirty work by explicit file/hunk transfer from the CSV manifest. Resolve the five overlaps manually; keep Phase 2B/2C performance investigation and generated browser artifacts excluded.
6. Safely isolate ownership for the remaining auth/session changes before committing them. Then run consolidated Go build/vet/tests, PostgreSQL and Redis isolation/security tests, four frontend test/typecheck/build suites, and the full Playwright suite on the selected changes.
7. Make separately reviewable commits and PRs. Push topic branches only. Never push directly to `main`, force-push, bypass checks, or merge a red/unreviewed PR.
8. Update the existing reports and final handover with actual current evidence. Do not claim Phase 2A or production readiness without all acceptance gates.

## Deferred scope and safety

Do not expand WS03 Payments or WS06 GitHub Actions. Do not deploy, alter production/development data, impersonate a tenant, or use non-disposable databases. Avoid disturbing unrelated Docker projects. The Go suite passed in disposable containers on the clean consolidation branch, but WS02 PostgreSQL, full frontend regression, and Playwright were not run there.

## Acceptance criteria

- Every transferred file has known ownership/dependencies and a recorded source/destination.
- Both original dirty worktrees and backup remain intact.
- All selected changes pass applicable backend, PostgreSQL, Redis, frontend, and browser tests on the consolidated branch.
- PRs target `main`; all applicable checks pass; merge confirmation and exact SHA are verified before reporting integration.
- Workstream statuses separate implemented, tested, committed, pushed, and merged work; Phase 2A remains partial until all required gaps close.
