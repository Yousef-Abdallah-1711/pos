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

## Continuation update — candidate hardening

The candidate now uses one read-only `REPEATABLE READ` transaction for all reconciliation counts, adds separate scope/account/amount/date counters and deterministic `empty` / `reconciled` / `discrepancies` status, and checks journal and GL row totals. Voiding now stores one application timestamp for both `voided_at` and the reversal effect so the diagnostic can compare the reversal to the exact void event date.

Integration scenarios now include missing and duplicate effects, amount/account/date/period mismatches, non-posted effects, cross-tenant contamination, orphan line references, an unbalanced journal, offsetting line-level errors, valid reversal, sibling restaurant isolation, and concurrent reconciliation reads during a posting. They temporarily disable integrity triggers only in the disposable test database and restore fixture state on cleanup.

The existing financial acceptance workflow now runs these tagged cases after applying migrations to its ephemeral PostgreSQL service. This acceptance has not yet run; only compile-only tagged validation is available locally.
