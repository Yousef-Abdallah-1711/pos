# Journal to GL reconciliation

## Candidate contract

`JournalEntryRepository.ReconcileJournalGL(ctx, tenantID, restaurantID)` is a read-only repository operation. It rejects missing scope, uses explicit tenant and restaurant filters for source journals, and compares stored PostgreSQL `NUMERIC` values without converting money to floating point. It returns deterministic counters, not a repair plan or financial adjustment.

It checks journal-line coverage for posted and void entries, duplicate or invalid effects, the original plus inverse pair used by the current audited void workflow, non-posted entries with GL rows, line/account scope, entry dates and fiscal periods, GL rows whose journal/line/account linkage is inconsistent, and source header totals against exact line sums. It does not reject legacy/manual journals by `entry_type`; unsupported historical source patterns still need accountant/engineering classification.

## Validation

- Compile-only repository package: passed (`go test -run '^$' ./internal/repository`).
- `go vet ./internal/repository`: passed.
- `git diff --check`: passed.
- Targeted disposable PostgreSQL tests: attempted and blocked before setup with `pq: role "test" does not exist`. Docker Desktop engine is unavailable and local PostgreSQL services are stopped. No PostgreSQL test result is claimed.

## Limitations

The candidate method has not been exercised against PostgreSQL, the original source dataset, or large-volume performance data. It is not wired to an API or financial user workflow. Exact accounting acceptance, issue-level exception details, and treatment of legacy/manual rows remain open. No source records were changed.
