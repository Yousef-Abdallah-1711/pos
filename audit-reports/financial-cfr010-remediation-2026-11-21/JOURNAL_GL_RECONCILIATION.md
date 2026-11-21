# Journal to GL reconciliation

## Candidate contract

`JournalEntryRepository.ReconcileJournalGL(ctx, tenantID, restaurantID)` is a read-only repository operation. It rejects missing scope, uses explicit tenant and restaurant filters for source journals, and compares stored PostgreSQL `NUMERIC` values without converting money to floating point. It returns deterministic counters, not a repair plan or financial adjustment.

It checks journal-line coverage for posted and void entries, duplicate or invalid effects, the original plus inverse pair used by the current audited void workflow, non-posted entries with GL rows, line/account scope, entry dates and fiscal periods, GL rows whose journal/line/account linkage is inconsistent, and source header totals against exact line sums. It does not reject legacy/manual journals by `entry_type`; unsupported historical source patterns still need accountant/engineering classification.

## Validation

- Targeted disposable PostgreSQL tests passed on exact candidate SHA `e68d70d45ccf507ee0474449fbf5e3f19a96f76e` in workflow [44953908598](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953908598). The workflow provisioned a disposable PostgreSQL/Redis environment, restricted runtime roles, and applied migrations before executing `go test -tags=integration -count=1 -timeout=90s ./internal/repository -run '^TestReconcileJournalGL_'`.
- The target package reported `ok pos-saas/internal/repository 0.276s`; the complete workflow passed, including live API, independent report oracle, browser, and integration steps.
- Earlier local ordinary Go tests, tagged integration compile-only test, `go vet ./...`, workflow YAML parsing, and `git diff --check` passed before the final exact-SHA run.
- Independent read-only review of code at `e68d70d` found no concrete actionable defects. The reviewer did not run tests.

## Limitations

The candidate method has not been exercised against the original source dataset or large-volume performance data. It is not wired to an API. Exact accounting acceptance, issue-level exception details, and treatment of legacy/manual rows remain open. No source records were changed.

## Continuation update — candidate hardening

The candidate now uses one read-only `REPEATABLE READ` transaction for all reconciliation counts, adds separate scope/account/amount/date counters and deterministic `empty` / `reconciled` / `discrepancies` status, and checks journal and GL row totals. Voiding now stores one application timestamp for both `voided_at` and the reversal effect so the diagnostic can compare the reversal to the exact void event date.

Integration scenarios now include missing and duplicate effects, amount/account/date/period mismatches, non-posted effects, cross-tenant contamination, orphan line references, an unbalanced journal, offsetting line-level errors, valid reversal, sibling restaurant isolation, and concurrent reconciliation reads during a posting. They temporarily disable integrity triggers only in the disposable test database and restore fixture state on cleanup.

The existing financial acceptance workflow now runs these tagged cases after applying migrations to its ephemeral PostgreSQL service. Exact-SHA execution passed as recorded above.

Further review-driven updates: non-posted-effect count reports distinct journal headers; regression coverage includes corrupt reversal date and period, explicit empty-scope status, and cleanup of trigger-disabled fixture anomalies. The hosted PostgreSQL run executed the targeted cases successfully.
