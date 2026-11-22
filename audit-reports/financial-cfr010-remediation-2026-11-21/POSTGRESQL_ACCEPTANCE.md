# PostgreSQL acceptance

## Corrected candidate run

- Workflow: [SERVIO financial acceptance run 44953908598](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953908598)
- Candidate SHA: `e68d70d45ccf507ee0474449fbf5e3f19a96f76e`
- Environment: GitHub Actions `ubuntu-24.04`; pinned PostgreSQL and Redis containers; random ephemeral credentials; disposable `servio_test`; runtime roles provisioned inside that disposable database; latest migrations applied.
- Target command: `go test -tags=integration -count=1 -timeout=90s ./internal/repository -run '^TestReconcileJournalGL_'`
- Target result: **PASS**, package output `ok pos-saas/internal/repository 0.276s`; the test step ran 2026-11-21 21:55:31–21:55:46 UTC.
- Same-run result: API build/start, authenticated API workflow, independent Balance Sheet SQL-oracle check, EN/AR desktop/mobile browser journeys, and the broader tagged handler/usecase/repository/API integration suite all passed. Workflow conclusion: **success**.
- Cleanup: API and dashboard stopped; disposable PostgreSQL and Redis containers removed successfully.

The target suite covers expected journal/GL effects and injected missing, duplicate, scope, account, amount, date/period, orphan, non-posted, unbalanced-header, reversal, offsetting-line, tenant/restaurant isolation, and concurrent-read scenarios. All fixture data was synthetic.

## Earlier failing run retained for history

- Run [44953642408](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953642408), SHA `59e5d92ad0dbb97ddc5520c7587f4da86e7b25aa`, reached PostgreSQL and exposed an invalid aggregate `FILTER` expression in the header query.
- The query was corrected to aggregate grouped journal rows using `HAVING` and count the mismatching groups. Run 44953908598 verifies that correction.

No original financial data was used or modified. Full CFR-010 acceptance remains open; this PASS is synthetic engineering evidence only.
