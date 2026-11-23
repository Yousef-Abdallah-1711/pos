# PostgreSQL Acceptance

Disposable PostgreSQL 16.4: container `cfr010-bs-rootcause-pg`, mapped port 25433, dedicated fresh database `cfr010_fix_test` for the final Balance Sheet integration tests. Redis container: `cfr010-bs-rootcause-redis`, mapped port 26380. Synthetic-only test records.

Commands and results:

- `go test ./internal/handler/http -count=1` — PASS.
- `go test ./internal/usecase -count=1` — PASS.
- `go test -tags=integration -v ./internal/usecase -run '^(TestCFR010Postgres|TestGenerateBalanceSheetPostgres)' -count=1` with `POS_TEST_DATABASE_URL` pointed at `cfr010_fix_test` and `SERVIO_DISPOSABLE_FINANCIAL_TEST_DB=1` — PASS (four tests, including opening balance, unreconciled residual, and backdated `as_of` movement checks).
- The isolated handler freshness regression failed before the cache bypass (second report remained 0.00), then passed afterward.

An earlier run against the already-used `servio_test` database failed the hard-coded tenant-1 integrity test because that legacy fixture intentionally leaves immutable posted entries. The failure was fixture contamination, not a query failure. The final run used a newly created, migrated `cfr010_fix_test` database and passed; the prior database was not cleaned or reset.

The PostgreSQL integration fixture validates opening, prior-period earnings, explicit close, current-period net loss, zero-activity P&L, repeatability, tenant isolation, restaurant isolation, backdated `as_of` boundaries, and no journal-state mutation by report reads. It does not cover every scenario in `SCENARIO_RESULTS.md`.
