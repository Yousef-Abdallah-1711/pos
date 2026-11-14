# Repository, PR, and CI Evidence

## Git state

- Repository: `Yousef-Abdallah-1711/pos` (`origin` = `https://github.com/Yousef-Abdallah-1711/pos.git`).
- Fresh `git fetch origin main` completed during audit; final `origin/main`: `018682beb9448eea6dd72fd2418876b439a94dbf`.
- Active report worktree: `C:\Users\yosea\.codex\worktrees\financial-audit-execution\pos`, detached at `018682beb9448eea6dd72fd2418876b439a94dbf`; only audit reports/evidence were added here.
- Protected checkout: `H:\Projects\pos`, branch `codex/phase2a-continuation`, HEAD `50380eaa1b95768ec40b1cb695035a601d26ab7c`.
- The protected checkout had substantial pre-existing untracked Playwright screenshots/snapshots and frontend test artifacts. It was not cleaned, reset, checked out to main, or edited for application code. New requested reports are isolated under `audit-reports/financial-remediation-master-audit-2026-11-14/`.
- Clean worktree created from the fetched main SHA at `C:\Users\yosea\.codex\worktrees\cfr011-entry-type-width\pos`.
- Candidate fix commit: `092ff33f4ac68d77e13480aad2afdd74ea56c511` (`fix(accounting): widen supplier journal entry type`).
- No local heavy Go/frontend test suite was run. `gofmt` and `git diff --check` were run for the new regression source.

## Current mainline PR reconciliation

| PR | Area | State | Evidence boundary |
|---:|---|---|---|
| #29 | CFR-001 auto-post imbalance validation | Merged | CI green; code and focused tests present. |
| #33 | Atomic journal/GL posting | Merged | CI green; disposable-PostgreSQL integration source present. |
| #34 | PostgreSQL journal integrity | Merged | CI green; migration 180 and direct-SQL/integration coverage. |
| #35 | Durable paid-order event/outbox recovery | Merged | CI green; outbox/lease/replay/reconciliation PostgreSQL tests. |
| #37 | Failed-event recovery operations | Open | Head `decdacfa6cc6d2052d2a64ed3e4ed3dd815fd166`; stale base `ea23d95a213f46fd8cdc59648b56e8474fc679f2`; checks green, no review, merge state `DIRTY`; tracker conflict requires refreshed base/resolution. |
| #39 | Payment/refund transition guards | Merged | CI green; rejects unsupported transitions, does not implement complete refund reversal. |
| #40 | Paid payroll immutability and atomic posting | Merged | CI green; PostgreSQL integration coverage present. |
| #41 | Financial numbering serialization | Merged | CI green; numbering integration test present. |
| #42 | Balance Sheet sign correction | Merged | CI green; current-period income is still not calculated into equity. |
| #44 | CFR-011 supplier entry-type width | Open | Head `6efe8ce9d3125162b518c3394fbe94815c183f35`, based on current main; all required checks green (run 43575466251), merge state `CLEAN`, review decision blank. |
| #45 | Dashboard FontPicker directive/BOM fix | Open | Head `cbd406aa5d33c3da7f73b88e2d723edf3a3de127`; required checks including application-images green (run 43580017753), merge state `CLEAN`, review decision blank. Unrelated to the financial remediation. |

Relevant merged PR statuses were read from GitHub PR checks. `production-ci.yml` defines a disposable PostgreSQL and Redis service, applies migrations, runs `go vet ./...`, and runs `go test -tags integration -count=1 -timeout 15m -json ./...`. It also runs four frontend app test/typecheck jobs, security/config checks, and a Chromium smoke against built images. This report does not invent suite counts where a result artifact was not inspected.

## PR #44 change/evidence

- Base: `018682beb9448eea6dd72fd2418876b439a94dbf`.
- Initial code head: `092ff33f4ac68d77e13480aad2afdd74ea56c511`; latest head after tracker update: `6efe8ce9d3125162b518c3394fbe94815c183f35`.
- Files: `backend/migrations/187_expand_journal_entry_type_to_32.sql`, `backend/scripts/rollback_187_journal_entry_type_to_32.sql`, `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`.
- `187` was selected because open PR #37 already uses `186` for its recovery-runtime grant migration.
- The tagged integration regression creates synthetic fixture context, posts both supplier entry types, reads the persisted type and verifies one GL row per journal line. It would fail on the current 20-character column.
- Downgrade checks for values longer than 20 and raises an exception instead of truncating financial history.
- Initial backend GitHub check: [run 43574869555](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43574869555) passed for the code-only head. At the final refreshed query, all required checks at `6efe8ce` were successful; the review decision remains blank.
- PR: [#44](https://github.com/Yousef-Abdallah-1711/pos/pull/44).

## Mainline workflow and deployment observations

- Mainline staging release verification workflow is defined in `.github/workflows/production-release.yml`. It waits for successful CI on a main push, checks the exact release SHA/status/environment, then performs public health/page GET checks.
- The most recent exact-main staging verification attempts for SHA `018682b` were skipped/cancelled, not successful. No verified deployment of `018682b` was established.
- Public release identity GET at audit time returned HTTP 200 with `environment=staging`, `status=candidate`, SHA `f733efb5e7f5730fea5a56a65d66d9c339e6c040`.
- Public health GET returned HTTP 200, `status=ok`, `db_healthy=true`.
- These observations confirm an API and database health response only. They do not expose or prove database contents, Redis isolation, integration credentials, production separation, or a disposable synthetic dataset.
- Browser observation is separately recorded in `staging-browser-observation.md`.
- The subsequent local disposable API/browser execution is separately recorded in `disposable-execution-2026-11-14.md`; it does not upgrade the staging isolation result.

## Commands and outcomes

1. `git fetch origin main` — succeeded; fetched `018682b...`.
2. `gh pr list --state all` — PR #37 open, #39–#42 merged, #43 unrelated draft, and #44 created by this audit.
3. `gh pr view <number> --json ...statusCheckRollup` — refreshed PR #37/#44/#45 checks and merge state; exact heads, results and run URLs are recorded below.
4. `gh workflow view production-ci.yml --yaml` — confirmed disposable Postgres/Redis and tagged integration workflow commands.
5. Read-only `GET /_ops/release` and `GET /api/v1/health` — HTTP 200; values recorded above.
6. Open dashboard `/en/dashboard` in fresh browser tab — reached login; no session submitted.
7. Source review confirmed supplier code values in `backend/internal/domain/auto_posting.go` are 21 characters while `backend/migrations/137_create_journal_entries.sql` declares `entry_type VARCHAR(20)`.
8. No production/staging data mutation, payment capture, merge, or deployment. GitHub collaborator count returned 1 (owner only); no independent reviewer was available.



## Refreshed repository state — execution update

- `git fetch origin main` completed again; `origin/main` is still `018682beb9448eea6dd72fd2418876b439a94dbf`. Active audit worktree is detached at this SHA; report-only changes are isolated under this directory.
- PR #37 remains open, stale-base and `DIRTY`, with all listed checks successful but no review decision. Exact head/base values are in the row above. GitHub merge-tree identified the authoritative tracker conflict.
- PR #44 remains open at `6efe8ce9d3125162b518c3394fbe94815c183f35`, based directly on current main. All required CI checks, including required gate and application-images, succeeded in run [43575466251](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43575466251); merge state is `CLEAN`, review decision is blank. No merge or deployment.
- PR #45 is unrelated dashboard FontPicker fix at `cbd406aa5d33c3da7f73b88e2d723edf3a3de127`; all returned CI checks including required gate and application-images succeeded in run [43580017753](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43580017753). It is `CLEAN`, open, with no review decision.
- GitHub collaborator count returned 1 (repository owner only). No independent authorized reviewer was available to request in this audit. Required owner action: add/identify an authorized independent reviewer and request review on #37 and #44 (and #45 if that fix is to be merged). Do not self-approve.
- Disposable execution, financial oracle, exact counts and limits are documented in [disposable-execution-2026-11-14.md](disposable-execution-2026-11-14.md).
