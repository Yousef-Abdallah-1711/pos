# FIN-016 implementation plan — recovery draft

## Verified candidate paths

- `backend/migrations/136_create_chart_of_accounts.sql`

## Approach

1. Compare the original issue and latest tracker continuation with current main and related PRs: not yet mapped. Establish whether code is merged, superseded, or still absent.
2. Reproduce in isolated synthetic fixtures with an independent oracle; preserve existing journal and audit history. Record exact environment and baseline SHA.
3. Design the smallest correction after identifying approved accounting decisions, API/UI impact, migration and rollback needs. Do not activate policy-dependent behavior without sign-off.
4. Run applicable PostgreSQL, isolation, concurrent, API, browser and performance checks; record exact SHA and results. Deliver through reviewed PR and validate main.

## Gaps

Architecture, alternatives, contracts, migration, rollback and issue-specific acceptance require owner review. This draft does not assert those decisions.
