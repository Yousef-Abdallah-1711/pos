# 01 Repository Inventory (Phase 1, read-only)

Tracked files: 1470 (`git ls-files`). Go source ~116.7k lines, TS/TSX ~125.5k lines (tracked, includes tests).

## Top-level layout
backend/ (Go API, cmd/{api,hash,migrate,migrate-shard,seed-themes}, internal/*, migrations/ 123 files, scripts/, tests/), frontend/ (pnpm+turbo workspace: apps/{dashboard,landing,restaurant-website,super-admin}, packages/{theme,ui,ui-components}, layers/{shared-components,ui-themes}), docs-organized/ (~125 tracked .md repo-wide), scripts/ (4 python), setup/ (2 .bat), .github/workflows (2), docker-compose{,.prod,.test}.yml, playwright.config.ts, root package.json (only @playwright/test).

## Tracked junk / artifacts (verified via git ls-files)
| Path | Size KB | Note |
|---|---|---|
| backend/api.exe | 26348 | binary |
| backend/api-local.exe | 25144 | binary |
| backend/api, main, pos-api, pos-backend, app | 11952/11080/10896/11280/10484 | Linux/other binaries |
| backend/pos-api.exe | 10964 | binary |
| backend/tmpclaude-0d2f-cwd, tmpclaude-cc8d-cwd | small | agent temp files |
| ActivateLicense-B73L9oSM.js, obd_app.js, tmp-beep.mp3 (root) | - | unrelated/stray |
| .playwright-mcp/*.png, ORD-2026-957363.pdf | 6888 (pdf) | MCP artifacts (PDF may contain order data) |
| .claude/settings.json, settings.local.json | - | local tool config tracked |
| frontend/apps/dashboard/src/i18n/messages/ar.json.bak-mojibake | 672 | backup; `*.bak` ignore pattern misses `.bak-mojibake` |
| scripts/__pycache__/*.pyc (2) | - | bytecode tracked |
Approx 150 MB of binaries in history/tree (sum of the above).

## Duplicate lockfiles
Tracked: frontend/pnpm-lock.yaml (authoritative, v9), frontend/apps/restaurant-website/pnpm-lock.yaml, frontend/packages/ui-components/pnpm-lock.yaml, frontend/layers/shared-components/package-lock.json (npm). Nested lockfiles are inside the pnpm workspace and ignored by pnpm; stale risk.

## .gitignore gaps
Ignores `*.log`, `*.bak`, `tmp/`, `tmp_*`, `.env*`, `.next/`. Does NOT ignore: `*.exe`, backend binaries, `tmpclaude-*`, `.playwright-mcp/`, `.claude/settings.local.json`, `__pycache__/`, `*.bak-*`, `test-results/`, `playwright-report/`. `backend/tmp` exists on disk. Also: `.dockerignore` at root exists, frontend/.dockerignore untracked (new).

## Largest source files (lines)
report_handler.go 2630; accounting_advanced_repo.go 2133; cmd/api/main.go 1894 (361 `mux.Handle*` registrations); accounting_advanced_handler.go 1346; payment_usecase.go 1110; white_label_handler.go 1086 (dead, see QA); scaling_api.go 1041 (dead); hr_report_handler.go 938; QRStudio.tsx 911; order_repo.go 904; dashboard/page.tsx 889; themeApi.ts 819.

## Tool versions declared
| Item | Declared | Where |
|---|---|---|
| Go | `go 1.24.0` | backend/go.mod; Docker golang:1.24-alpine; README says "Go 1.18+" (stale) |
| Node | engines >=18; Docker node:20-alpine; CI node 20 | frontend/package.json, Dockerfile.prod |
| pnpm | packageManager pnpm@9.0.0, engines >=9; Dockerfile corepack pnpm@9.0.0 | frontend/package.json |
| Postgres | postgres:15-alpine | docker-compose*.yml |
| Alpine runtime | alpine:3.20 | backend/Dockerfile |
No .nvmrc / .tool-versions found at root. Root package.json has no engines.
