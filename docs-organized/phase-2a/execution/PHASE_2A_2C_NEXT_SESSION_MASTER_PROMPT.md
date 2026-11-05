# Next Session Master Prompt — POOS Unified Phase 2A + Phase 2C

Continue the POOS unified engineering closure from the verified state below. This is a self-contained continuation prompt.

## Repository and baseline

- Repository: `https://github.com/Yousef-Abdallah-1711/pos`
- Last verified `origin/main`: `d84120559db36e8e9e7fdb73872a30cbc9d4e6e7`
- Integration branch/worktree: `codex/unified-phase2a-2c-closure` at `H:\pos-unified-phase2a-2c`; PR #9 remote head is `ec23651ebc42a397b54244554635a0e281122d27`; benchmarked code commit is `e30ffe3`.
- PR #9 is open; both hosted checks pass at the current head, while independent human review approval remains outstanding.
- Read `H:\Projects\pos\CODING_STANDARDS.md` before code changes (mandatory project instruction).
- Preserve dirty source worktrees `H:\Projects\pos` and `H:\ws02`; never reset, clean, overwrite or blind-stage their contents.
- A private snapshot at `C:\Users\yosea\.codex\backups\poos-unified-closure-20261105-014347` was checksum-verified 220/220. Keep it private.
- Reconciliation CSV: `docs-organized/phase-2a/execution/PHASE_2A_2C_RECONCILIATION_MANIFEST.csv`.

## Current integration candidate

Five committed WS02 PR #1 commits were cherry-picked to the clean branch: `7075011`, `a268f9b`, `1816970`, `63c92cd`, `b249ac4`. During clean checkout testing, a hidden dependency was found in `subscription_admin_isolation_integration_test.go`; its database and synthetic data fixture helpers existed only in the source WS02 untracked tests. The integration test has been changed to use self-contained fixtures. This correction is committed at `f56f32e`.

A disposable Go 1.24.13 run migrated an isolated tmpfs PostgreSQL database successfully and passed `go build ./...`. Vet/full test rerun is in `H:\pos-unified-verification\backend-isolated.log`; inspect exit markers and finish targeted integration tests before publishing.

## Source/PR facts

- PR #1 open; old head checks fail because stale QA verifier uses `http://localhost:3002` and two stale Dashboard production URL tests no longer match current main. Current main contains fixes; rebase candidate should be checked with hosted CI.
- PR #3 open; checks pass, no approval. It is mostly handover/docs and requires deduplication against the new unified deliverables.
- PR #7 open draft; checks pass, no approval. Its Admin Orders optimization is reconciled into the local unified candidate. Fresh isolated measurements at page sizes 1/10/25/50 show fixed three-query candidate behavior, matching canonical response hashes/payloads and exclusion of tenant/sibling-restaurant fixture orders. See the unified verification report for exact metrics.
- PR #8 merged; Storefront coalescing is in `origin/main`. PRs #2, #4, #5, #6 merged; avoid duplicate patches.
- Current remote main SHA at handover: `d84120559db36e8e9e7fdb73872a30cbc9d4e6e7`.

## Required next actions

1. Inspect `backend-isolated.log`; repair any actual failures without weakening tenant isolation. Run named integration tests using only the disposable PostgreSQL fixture.
2. Update the verification report with exact current results. Re-run `git diff --check`, inspect every candidate hunk, secret scan, and migration consistency before commit.
3. Review the remaining committed and dirty WS02 changes path-by-path; specifically Orders handler/repository/use-case and all tenant/restaurant database side effects. Preserve original trees.
4. Benchmark Admin Orders with comparable synthetic data and build fingerprints at page sizes 1, 10, 25, 50. Record SQL/detail/item counts, query timing, HTTP median/p95, payload, response equality, cross-tenant and cross-restaurant negatives. If instrumentation/fixtures are insufficient, report blocked and do not claim improvement.
5. Dashboard has 107 failing tests out of 645 and test-source TypeScript errors. Its app typecheck and all four production builds pass; other frontend suites pass. Fix the Dashboard test failures and rerun all frontend gates.
6. Run isolated Redis tests and browser login/logout, auth/session, tenant/restaurant, POS, inventory, storefront/cart, Super Admin and Admin Orders journeys against a clean candidate Docker stack. Do not use the existing `pos-dev` stack.
7. Obtain required independent review on PR #9 after addressing the Dashboard test gate; merge only if all gates pass. Never push to main directly or bypass checks/review. Fetch main and record its exact merge SHA.
7. Maintain the five handover deliverables in `docs-organized/phase-2a/execution/`: final handover, verification report, reconciliation manifest, main integration report, next-session prompt. Clearly distinguish implemented, tested, committed, pushed, merged, deferred, blocked and production-ready.

## Known status and safety posture

Phase 2A PARTIAL; Phase 2C PARTIAL; Unified Git Integration PARTIAL; Production Release NO-GO. Super Admin support handoff must remain disabled/fail-closed. Payments (WS03) and broad CI expansion (WS06) remain deferred. Do not submit production orders/payments or run tests against development/customer databases. Do not close PRs before confirming their intended changes are merged or preserved.
