# Combined candidate verification

## Integration identity

Detached worktree: `H:\pos-pr-candidate-integration`. Base `a6a2c6700a5fe883256deb1e5a3893108bd9865c`. Applied without conflicts in proposed order: baseline commit `5528c6f4bac5c6575c3a185a45ad81eaffa36e6e`, then Phase 00 commit `5064daa89a27715c3eeacb8becc9f3f48de1853a`. Combined HEAD: `8062e54a559a85dc889ed37f78593dbd8890ccb9`.

## Acceptance matrix

| Check | Result | Evidence / qualification |
|---|---|---|
| `git diff --check` | PASS | Combined worktree; no whitespace errors |
| Full frontend test typecheck | PASS | Combined candidates |
| Focused tests | PASS | 4 files / 9 tests on combined candidate |
| Dashboard Vitest | PASS | 48 files / 656 tests |
| Production typecheck | PASS | Combined candidate |
| Targeted ESLint | PASS | Legacy-config compatibility and `--no-ignore` used |
| Dashboard production build | PASS | Combined candidate |
| Go tests | PASS | `go test ./...`, Go 1.24 container |
| Go vet | PASS | `go vet ./...`, Go 1.24 container |
| Disposable PostgreSQL migration | PASS | Entire migration set completed on PostgreSQL 15, tmpfs-only |
| Phase 00 PostgreSQL integration | PASS | `TestGenerateBalanceSheetPostgresIntegrityStatus`, one run |
| Financial browser assertions | PASS WITH OPEN ISSUE | Seven synthetic scenarios passed their assertions; see below |
| Authenticated tenant report/export scenarios | HISTORICAL PASS; NOT RE-RUN | Prior recorded 12 endpoint/role scenarios; does not establish product-wide isolation |
| Arabic document metadata (`lang` / `dir`) | OPEN / NOT VERIFIED | Prior review marked this as separate architecture issue |

## Browser evidence and caveats

The synthetic production-build scenarios asserted: unreconciled warning, preservation of source figures/no export, reconciled export availability and PDF download, unknown-status fail-closed behavior, mobile Arabic warning without horizontal overflow, loading state, API error state, and unauthorized redirect/figure withholding. Seven assertions passed. Artifacts: `evidence/browser/browser-acceptance.json` and synthetic PDF.

The same run recorded minified React error **#419** on English and Arabic route navigations. It was not diagnosed; therefore browser verification is not clean. The mobile Arabic probe found a translated warning and no overflow, but its `documentElement.dir || body.direction` measurement returned **`ltr`**. This does not substantiate RTL direction and remains open. Synthetic 500/401 console errors correspond to deliberate error fixtures; they are expected test stimuli.

The harness recorded 1,135 API requests, including repeated settings requests. This synthetic route harness is not a reliable production performance sample. No tenant-isolation conclusion is made from it.

## Database safety and failed attempts

The passing database run used a newly named PostgreSQL 15 container with `/var/lib/postgresql/data` mounted as tmpfs, a disposable network, seeded test-only tenant/restaurant/account fixtures, and the repository runtime-role initializer. The container/network were removed on completion. Two earlier disposable attempts failed before test execution (first because roles required by migrations were omitted; second because migration seed data already contained ID 1). Their logs are preserved as evidence and are not counted as product failures. No protected database or pre-existing PostgreSQL process was used.

Successful raw log: `evidence/logs/combined-go-postgres-retry2-2026-11-11.log`. Failed setup logs: `combined-go-postgres-2026-11-11.log`, `combined-go-postgres-retry-2026-11-11.log`.
