# CFR-002 implementation plan

## Baseline and scope

Main `afda1f91c67922a92f66507adc208e6e1943d631` includes PR #34 and migration 180. Preserve posted history; do not repair historical records without an approved plan. The primary risk is a database writer bypassing application validation.

## Execution

1. Provision a new isolated PostgreSQL database, apply the repository migrations and runtime-role setup, and record the migration and database versions. Confirm no existing database is targeted.
2. Run `backend/scripts/preflight_180_journal_integrity.sql` against synthetic valid and invalid legacy fixtures. Check that the invalid fixture blocks migration without partial guards or data changes.
3. Run `go test -tags=integration ./internal/repository -run '^TestCFR002' -count=1` from `backend`. Independently check journal header totals against line sums and original/inverse GL effects by tenant and restaurant.
4. Exercise documented downgrade and forward recovery on disposable data, verifying that financial rows are unchanged. Retain the historical evidence separately from the new run.
5. Obtain accounting decisions on historical repair, linked reversal representation and closed periods; reconcile CFR-004, CFR-005 and CFR-013 before proposing changes. Review privileged roles and restrict owner/superuser access.
6. Record exact SHA, commands, results, exceptions and required reviewers in the issue register. Close only after the finding-level and financial release gates are met.

## Sources and boundaries

Migration 180, its preflight and rollback scripts, `journal_entry_repo.go`, and `journal_integrity_postgres_integration_test.go` are the primary implementation and test paths. Manual Dashboard coverage is not a substitute for direct SQL assertions; document any user workflow tested separately.
