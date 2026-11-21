# Reconcile chart-of-accounts current-balance cache

`chart_of_accounts.current_balance` is a derived cache; posted journals, journal lines, and `general_ledger` rows are the source history. Backdated postings made before the chronological balance fix may have left account caches stale. The new posting code repairs accounts it touches, but it cannot repair untouched accounts by itself.

## Preflight

1. Schedule a maintenance window. Stop and drain all API instances, workers, and administrative jobs that can post journals or write chart-of-accounts records. Confirm no posting transaction is in flight. The script takes `SHARE` locks on the two source/cache tables; do not rely on those locks as a substitute for quiescing writers because a posting transaction can already hold an account row lock before it writes the GL.
2. Take the ordinary database backup/snapshot required by the deployment policy. Run `preflight_chart_of_accounts_current_balance.sql` and inspect `stored_balance`, `expected_balance`, and `difference` in currency units (two decimal places). Verify the independent sum formula against the GL movement rows and normal-balance direction, then approve the affected accounts/total before running the repair.
3. Confirm the exact script checksum and target database. Never run this against public, staging, or production from an engineering acceptance session.

## Execution and validation

Run `reconcile_chart_of_accounts_current_balance.sql` using `psql -v ON_ERROR_STOP=1` during the approved maintenance window. Its transaction locks GL/account writes, emits changed account IDs and differences, recomputes the cache for every scoped account—including accounts with no movements, which become zero—and raises an exception if post-update values do not match the independent movement sum. Keep service writers paused until the transaction commits and validation output is saved.

The correction changes only `chart_of_accounts.current_balance` and `updated_at`. It does not edit or recreate posted journal history or stored posting-time `general_ledger.running_balance` snapshots. Post-run checks should compare each account cache to an independent SQL aggregation using tenant, restaurant, account, and normal-balance direction; sample debit-normal and credit-normal accounts, zero-movement accounts, and accounts with negative balances; verify GL row counts and posted journal totals are unchanged; then resume writers.

## Rollback

Before commit, any SQL error or unexpected reconciliation result aborts the transaction; issue `ROLLBACK`, keep writers paused, and investigate. After a successful commit, do not restore stale cache values as a rollback. The repair is idempotent: correct the underlying issue and rerun the same reconciliation procedure. If a committed correction itself is suspected to be wrong, stop writes and reconcile against source GL evidence before any further update.

## Locking and performance

This is an explicit maintenance operation, not an automatically applied schema migration. The transaction scans the scoped ledger and updates only accounts whose cached value differs. `SHARE` locks block concurrent writes to `general_ledger` and `chart_of_accounts` until commit. The preflight is intended to estimate scan time and discrepancy volume; schedule enough downtime for the full aggregation and verification. Do not run concurrently with journal posting or account editing.
