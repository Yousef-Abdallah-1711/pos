# PostgreSQL acceptance

## Attempt 1

- Workflow: [financial acceptance run 44953642408](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953642408)
- Candidate SHA: `59e5d92ad0dbb97ddc5520c7587f4da86e7b25aa`
- Environment: existing GitHub financial acceptance job; ephemeral PostgreSQL and Redis containers; latest migrations applied; ephemeral credentials; database `servio_test`.
- Target step: `go test -tags=integration -count=1 -timeout=90s ./internal/repository -run '^TestReconcileJournalGL_'`
- Result: **FAILED**. Each scenario reached the repository method, which returned `pq: aggregate functions are not allowed in FILTER` from the journal-header totals query. This is a confirmed implementation defect, not an infrastructure failure.
- Follow-up: replaced the invalid aggregate filter with grouped `HAVING` and an outer count. A fresh exact-SHA run is required before reporting PostgreSQL acceptance.

No original financial data was used or modified. The workflow destroyed its disposable services after failure.
