# CFR-011 implementation plan — recovery draft

## Verified candidate paths

- `backend/scripts/rollback_187_journal_entry_type_to_32.sql`
- `backend/migrations/187_expand_journal_entry_type_to_32.sql`
- `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`
- `backend/internal/usecase/auto_posting_engine.go`

## Approach

1. Treat merged migration 187 and the passing constructed-journal test as the narrow schema fix; rerun on final main.
2. Test upgrade with existing journal rows and safe rollback refusal using disposable PostgreSQL; never truncate financial history.
3. Inventory supplier source paths and obtain a support/mapping decision for the currently uncalled auto-post hooks.
4. If supported, test invoice/payment lifecycle, scope, failure recovery, reversal, reporting and export with an independent journal/GL oracle. Implement only authorized behavior and validate the reviewed candidate on main.

## Gaps

Constructed entries do not prove the supplier API invokes these hooks. Source lifecycle and accounting mapping remain unresolved.
