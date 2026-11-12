# Go build & test baseline (2A-003 / 2A-004)
Toolchain: Go is not installed on the host; verified in `golang:1.24-alpine` (go1.24.13, matches `go 1.24.0` in go.mod). Backend mounted read-only; module/build caches in throwaway volumes `pos2a_gomod`, `pos2a_gobuild`.
Commands (from repo root, `MSYS_NO_PATHCONV=1`):
`docker run --rm -v <repo>\backend:/src:ro -v pos2a_gomod:/go/pkg/mod -v pos2a_gobuild:/root/.cache/go-build -w /src golang:1.24-alpine sh -c 'go build -o /dev/null ./...'` and `go test -count=1 $(go list ./... | grep -v "/scripts$")`. No database reachable → DB-dependent tests are untested, not passed.

## Build
`go build ./...`: **1 failure** — `pos-saas/scripts`: `main redeclared` (`scripts/sql_exec.go` vs `scripts/embedded_pg_daemon.go`). All other packages (incl. cmd/api, migrate) compile. Pre-existing; not touched by Phase 2A. `go vet` (run implicitly by `go test`) also reports `internal/sharding/migration.go:174` log.Printf `%d` with a float64.

## Tests (21 packages with tests excluding `scripts`)
PASS (4): `cmd/api`, `internal/config`, `internal/middleware`, `internal/service`.
FAIL at runtime (5), all pre-existing, none security-related except as noted:
- `internal/domain`: `TestSlugValidation/single_char` (isValidSlug("a") false).
- `internal/handler/http`: KDS flow/bump tests (500; scan expects 20 dest args, query returns 21 — KDS column drift), report-export tenant-context tests noisy.
- `internal/handler/ws`: `TestKDSHubBroadcastRoutesMessageToMatchingRoom` timeout.
- `internal/repository`: KDS queue SQL expectation mismatch; Phase6 sql-injection/tenant-isolation list tests (`sql: no rows`); theme repository null/slug tests.
- `internal/usecase`: KDS use-case scan-arity failures (same root cause as above).
BUILD FAILED (9) — test code stale vs production code: `internal/sharding` (vet), `internal/validation` (undefined NewThemeValidator/models), `tests/integration` (GenerateToken signature), `tests/security` (vet: bad Example name), `tests/unit/{auth,middleware,security,sharding,scaling}` (unused vars/imports, signature drift, undefined helper).
No test files: auth, observability, pkg/{database,jwt}, scaling, security, util.
**Implication:** the existing suite cannot serve as a regression gate; most "security" tests under `tests/` don't compile. New regression tests for Phase 2A go in packages that compile (`internal/middleware` etc.).
