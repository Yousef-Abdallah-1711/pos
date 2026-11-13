# 04 — Phase 00 test results

- Go: `go test ./internal/usecase ./internal/handler/http` — PASS in golang:1.24 container.
- Go static: `go vet ./internal/usecase ./internal/handler/http` — PASS.
- PostgreSQL integration: `go test -tags integration ./internal/usecase -run '^TestGenerateBalanceSheetPostgresIntegrityStatus$' -count=1` — PASS on isolated migrated PostgreSQL 15.
- Integration cleanup verification: journal_rows=0 and ledger_rows=0.
- Go formatting: gofmt applied to changed Go files.
- Frontend install: after user approval, lockfile-pinned pnpm 9.0.0 installed selected workspace dependencies with `--frozen-lockfile --ignore-scripts` into a private worktree store. Lockfile SHA-256 remained `E23E2910259F119D67C9BF4B2C88A8D2CC5091E5268EF1218005AA904D7C5E95`. Earlier compatibility probing added pnpm 11.28.9 to the shared pnpm store; that shared store was not modified afterward.
- Frontend typecheck: PASS with pinned pnpm 9.0.0. The earlier pnpm 11.28.9 auto-install/React diagnostics did not reproduce.
- Relevant component/helper tests: PASS — 2 files, 3 tests. `typecheck:test` remains FAIL due existing unrelated diagnostics in `TableList.test.tsx` (duplicate JSX attributes) and `config.test.ts` (`NODE_ENV` readonly assignments).
- Dashboard production build: PASS, including synthetic local URL configuration used by browser checks.
- Targeted ESLint: PASS using the repository's legacy ESLint configuration (`ESLINT_USE_FLAT_CONFIG=false --no-ignore`).
- Browser desktop/mobile/UI acceptance: PASS for seven synthetic-fixture scenarios; see `08-BROWSER-ACCEPTANCE.md`. This does not constitute authenticated end-to-end or tenant-isolation coverage.
- Tenant adversarial tests: NOT RUN.

Therefore tenant authorization and some frontend test-suite gates remain NOT VERIFIED/FAIL; Phase 00 acceptance is PARTIAL.
