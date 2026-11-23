# CFR-004 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/repository/journal_entry_repo.go`
- `backend/internal/repository/general_ledger_repo.go`
- `backend/internal/usecase/auto_posting_engine.go`
- `backend/cmd/api/main.go`
- `backend/internal/usecase/auto_posting_engine_test.go`
- `backend/internal/usecase/phase6_accounting_cycle_test.go`
- `backend/internal/repository/general_ledger_atomic_posting_integration_test.go`

## Approach

1. Use the merged atomic journal/GL path and paid-order outbox described in the specification as the candidate; verify migration ownership and restricted roles in a disposable PostgreSQL environment.
2. Run `go test -tags=integration -count=1 ./internal/repository -run '^TestCreateAndPostAutoEntryPostgres'` and the paid-order/accounting-event worker PostgreSQL tests. Record exact SHA, database version and each result.
3. Exercise actual API payment, process stop/restart, failed-event inspection and redrive using synthetic orders; compare order, event, journal and GL records independently. Test concurrent delivery and both tenant and restaurant scope.
4. Resolve confirmed engineering failures with focused changes. Define provider/alert contracts and policy decisions separately; do not imply that the outbox covers them. Obtain review and exact-head CI before merge and revalidate on main.

## Gaps

External alerting, provider settlement and operational service-level expectations remain undefined. Existing historical PostgreSQL tests do not establish live restart and operator acceptance on current main.
