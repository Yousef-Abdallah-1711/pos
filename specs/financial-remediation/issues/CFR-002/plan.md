# CFR-002 implementation plan — recovery draft

## Verified candidate paths

- `backend/migrations/180_enforce_journal_integrity.sql`
- `backend/scripts/preflight_180_journal_integrity.sql`
- `backend/scripts/rollback_180_journal_integrity.sql`
- `backend/internal/repository/journal_entry_repo.go`
- `backend/internal/repository/journal_integrity_postgres_integration_test.go`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: 34. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
