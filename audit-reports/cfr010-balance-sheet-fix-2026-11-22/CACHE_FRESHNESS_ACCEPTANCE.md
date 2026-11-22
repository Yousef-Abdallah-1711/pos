# Cache Freshness Acceptance

## Evidence scope

The report handlers for Balance Sheet, Profit & Loss, Cash Flow, and Trial Balance bypass the server response cache. Inventory valuation retains its existing cache. The handler tests use a five-minute in-memory cache and make identical requests with the same tenant, restaurant, dates, and path while the mocked source results change.

| Report | First response | Second response | Evidence |
|---|---:|---:|---|
| Balance Sheet | assets 0; RECONCILED | assets 55, liabilities 5, equity 0; UNRECONCILED, residual 50 | `TestFinancialReportHandler_GetBalanceSheetReflectsNewPostedLedgerActivity` |
| P&L | revenue 1,000; net 700 | revenue 1,200; net 900 | `TestFinancialReportHandler_GetProfitLoss_ReflectsNewPostedLedgerActivity` |
| Cash Flow | net 0 | net 50 | `TestFinancialReportHandler_CashFlowRefreshesAfterLedgerActivity` |
| Trial Balance | debit/credit 0/0 | debit/credit 55/55 | `TestFinancialReportHandler_TrialBalanceRefreshesAfterLedgerActivity` |

On 2026-11-22, `go test ./internal/handler/http -count=1` and `go vet ./internal/handler/http` passed with these tests. These are handler tests with SQL mocks: they verify repeated responses query the changed source rather than serving the old cache, but do not themselves post a real journal or prove concurrent snapshot consistency. The earlier disposable PostgreSQL and authenticated browser findings are recorded in `POSTGRESQL_ACCEPTANCE.md` and `BROWSER_ACCEPTANCE.md`.

The policy-gated Balance Sheet equation remains unreconciled in the synthetic revenue example. Passing these tests is not accounting approval or production acceptance.
