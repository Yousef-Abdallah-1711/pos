# Phase 2A Next Session Master Prompt

You are continuing POOS Restaurant POS SaaS Phase 2A. Do not repeat completed discovery. Do not start Phase 2B. Repository: `https://github.com/Yousef-Abdallah-1711/pos`.

## Exact baseline and current publication state

At the last verified fetch, `origin/main` was `5ab6dba62287ddd404ce00529f5fe99220ccdc19`. Do not assume later PRs merged: first fetch and verify. PR #1 (`codex/ws02-final-security-closure`) is open with WS02 HR/Accounting/subscription changes and historically failed `runtime-localhost-guardrail` and `url-domain-config-smoke`. PR #3 (`codex/phase2a-final-consolidation`) contains commit `dcdcbe9`, a narrow fix for those baseline failures; at handover its guardrail passed and URL smoke was pending. Verify PR state/checks before acting.

Two original worktrees remain dirty and preserved: `H:\Projects\pos` (`codex/phase2a-continuation`, HEAD `ecadb5ab7665df3eb9f69886a9b186c710eff62d`, 260 dirty paths) and `H:\ws02` (`codex/ws02-final-security-closure`, HEAD `9cb9a288ff3538c68e11caa2be9f50dd79cf8ee2`, 60 dirty paths). Five dirty overlaps are listed in `PHASE_2A_FINAL_SESSION_HANDOVER.md`. A checksum-verified private backup is at `C:\Users\yosea\.codex\backups\poos-phase2a-consolidation-20261104-213732`. Never reset, clean, overwrite, stash-destructively, or blindly merge these worktrees. Consolidation worktree: `H:\pos-phase2a-integration`, branch `codex/phase2a-final-consolidation`.

## Architecture and status

See the existing Phase 2A reports under `docs-organized/phase-2a/execution/` and the final handover in this directory. WS01 is recorded complete by existing evidence. WS02 is partial; WS03 Payments deferred; WS04 partial/current checks need full regression; WS05 partial; WS06 CI expansion deferred; WS07/08/09/10 partial. Do not describe topic/dirty changes as part of `main` unless confirmed reachable from fetched `origin/main`.

A security-critical limitation is intentional: Super Admin tenant handoff is disabled (UI disabled, API HTTP 503, no tenant session creation). Never re-enable unrestricted impersonation. Any future support feature needs actor-attributed least-privilege sessions and auditable start/end events.

## First priorities

1. Read repository `CODING_STANDARDS.md` before code changes; the external path in `AGENTS.md` was unavailable on the prior host.
2. Fetch and inspect `main`, PR #1/#3 states, checks, branch rules, and all worktrees. Preserve their current dirtiness.
3. Re-run/check the PR #3 URL/domain hosted gate. If merged, record exact new main SHA. Re-evaluate PR #1 against it and rerun all applicable checks.
4. Review PR #1 source and tests, then run WS02 HR/Accounting/subscription integration tests against a newly created isolated PostgreSQL database and verify unauthorized DB side effects.
5. Reconcile dirty work by explicit file/hunk transfer from the CSV manifest. Resolve the five overlaps manually; keep Phase 2B/2C performance investigation and generated browser artifacts excluded.
6. Run consolidated Go build/vet/tests, PostgreSQL and Redis isolation/security tests, four frontend test/typecheck/build suites, then Playwright journeys on the selected changes.
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
