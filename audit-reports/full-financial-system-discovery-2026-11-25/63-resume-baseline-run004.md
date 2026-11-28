# 63 — Resume baseline (RUN-004, started 2026-11-27 ~03:55 local)

Authorization (user message, 2026-11-27): finish the audit (option B) · landing/frontend changes are intentional and must be preserved · no product fixes · no commit/push/merge/deploy · no subagents.

## Baselines (kept separate)
| Baseline | Value |
|---|---|
| 1. Historical audit baseline (RUN-001..003 evidence) | SHA `133260db66b1143c254d79476daa7650b41a33f0` |
| 2. Authorized frontend development changes | **UNCOMMITTED working-tree changes** (see below). The brief says they were committed; verified fact: HEAD, branch (`master`), reflog and `git log --all` show **no new commit** since 133260d. Nothing was committed by the audit. |
| 3. Current audit execution baseline (RUN-004) | HEAD `133260d…` + working tree (34 modified, 6 deleted tracked files; 107 untracked paths). Evidence for RUN-004 is attributed to "133260d + dirty tree". |
| 4. Pre-existing local modifications | 2 stashes preserved (`local edits before reset to origin/master`, `pre-pull local changes`); tsconfig/next-env tooling edits |
| 5. Audit-only files | `audit-reports/full-financial-system-discovery-2026-11-25/` (itself untracked, i.e. NOT in any commit/remote), `evidence/runtime/RUN-004`, scripts under `tests/` |

Raw capture: `evidence/runtime/RUN-004/01-git-baseline.txt`. Remote `origin = https://github.com/yousefabdallah171/POSS.git`.

## Frontend work (preserved, not touched)
Tracked: landing `marketing.css`, `styles/marketing-*.css`, `Navbar`, `Footer`, `RootDocument`, `TrialLink`, `ClosingSections`, `CaptureMedia`, i18n `en/ar` marketing+legacy JSON, landing tests/config/package.json, `frontend/pnpm-lock.yaml`, **`dashboard/src/lib/api.ts` (a dashboard file, outside the "landing" description — flagged for revalidation)**. Deleted tracked: `HomeHero`, `HospitalityStory`, `ManagementSections`, `OperationsSections`, `marketing-hero.css`, `marketing-product.css`. Untracked: new marketing components (`ConversionLanding`, `FeatureShowcase`, `OrderStory`, …), `marketing-coral/interactive/workflows/captures.css`, `frontend/apps/landing/.hallmark/` (previously at repo root `.hallmark/`; no longer present there), `.gsap/`, `docs/landing-*`, `docs/superpowers/`, `.playwright-mcp/*`, `writing-block.md`.
Audit never writes outside `audit-reports/` and the isolated Docker project.

## Impact on previously tested functionality
- Landing visual differences vs historical screenshots are NOT regressions.
- `frontend/apps/dashboard/src/lib/api.ts` changed: it is the dashboard API client, so dashboard-driven evidence (RUN-003 UI tests) may be affected → diff must be reviewed and a smoke re-run of login + one mutation done (task T-R04-BASE-002).
- Landing: registration/login/language switching must be re-verified (T-R04-BASE-003).
- Backend, migrations and accounting code are unchanged vs the audit SHA (`git diff --stat 133260d` has no `backend/` paths) → financial evidence stays valid.

## Port 3001 (resolved safely)
Owner changed from the documented pid 6068 to **pid 9064**: `node.exe C:\nvm4w\nodejs\node.exe … next start -p 3001` from this project's `frontend/apps/landing` (started 2026-11-27 03:39 local; parent `cmd.exe` pid 5516 → 23128). Identity verified (command line contains landing path and `start -p 3001`); no established client connections. Stopped pid 9064 only (authorized: obsolete server of the finished frontend work); parent 23128 and all other processes untouched. Port 3001 now free.

## Runtime
Isolated compose project `pos-p2a-ws05-e2e` resumed on the persistent volume; API recreated with the storefront CORS wildcard (`http://*.localhost:13113`) and all four frontends recreated (shared network namespace). Health: API 200; landing `/en` 200, `/en/auth/login` 200, `/ar` 200; dashboard `/en/dashboard` 200 (35 s first compile); super-admin `/login` 200. Dev containers still stopped (restore at closure).
DB at resume: tenants 9, users 12, orders 74, journals 17,155, GL lines 34,365. The RUN-003 checkpoint file lists users 11 (and a truncated GL figure): one user was added after the checkpoint snapshot — noted as drift, final dump unchanged on disk.

## Update during RUN-004: HEAD advanced (verified 2026-11-27)
While RUN-004 was running, the separate frontend workstream committed its work: HEAD moved from `133260db66b1143c254d79476daa7650b41a33f0` to **`0de839b665982a2072f520c20a1aaf9f502f4bea`** (branch master): `3af62dc` feat(landing): build interactive SERVIO conversion page (2026-11-27 04:14 +0200), `4b62e5a` fix(ci): allow local landing browser QA scripts (04:24), `0de839b` fix(landing): align marketing acceptance with conversion page (05:25).
- Verified diff 133260d..HEAD: 107 files; 69 under `frontend/apps/landing/`, docs under `docs/landing-*`, `frontend/pnpm-lock.yaml`, `scripts/check_runtime_localhost.py`, and one dashboard file `frontend/apps/dashboard/src/lib/api.ts` (401 redirect keeps the locale). **0 files under `backend/`.**
- Working tree now: 6 modified tooling files (tsconfig/next-env written by Next dev servers) and untracked paths (`.playwright-mcp/*`, audit folder); the two stashes are intact (count 2).
- Evidence attribution: RUN-004 evidence captured before the commits (landing smoke, form/control matrices, metrics) ran against 133260d + the then-uncommitted working tree of the same landing work; backend-related evidence is unaffected because backend code is identical. The landing smoke result (PASS) is for the tree as it stood at about 03:55 local time and was not repeated against the final commit; repeat it before relying on it for 0de839b.
- Historical RUN-001..003 evidence stays attributed to 133260d.
