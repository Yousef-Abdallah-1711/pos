# PostgreSQL Reconciliation Continuation

On 2026-11-22, a local PostgreSQL 16 container and newly migrated `cfr010_final_test` database were used for the continuation. The first migration attempt stopped at migration 164 because application roles were absent. After creating the required disposable roles, migration completed. The database contains synthetic data only and the container was stopped after testing.

`go test -tags=integration -v ./internal/usecase -run '^(TestCFR010Postgres|TestGenerateBalanceSheetPostgres)' -count=1` passed four tests. The backdated test reported expected and actual assets of 7,500 cents at the earlier date and 18,500 cents at the later date, with zero differences.

`go test -tags=integration -v ./internal/usecase -run '^TestCFR010SalesAndReversalAccounting$' -count=1` passed after correcting the new oracle's NUMERIC-to-integer scan. Two posted sales produced assets 7,700, liabilities 700, equity 0, residual 7,000, P&L net 7,000, and Trial Balance 7,700/7,700 cents. A subsequent partial reversal produced 6,600/600/0/6,000/6,000/6,600/6,600 cents respectively. A separate tenant-and-restaurant-scoped GL aggregate agreed. These are synthetic journal scenarios; the customer refund flow was not exercised. Equity treatment remains POLICY_PENDING.

These results do not establish original operating data or complete the full scenario matrix. See `SCENARIO_RESULTS.md`.

## Historical void regression

`TestCFR010VoidPreservesHistoricalProfitLoss` first failed on disposable PostgreSQL: a 2025 posted sale voided later produced historical Balance Sheet assets 50.00 but historical P&L 0.00, instead of 50.00. The P&L query filtered by the journal's current `posted` status and erased the original movement from prior periods.

The correction includes the original movement in its entry-date period for journals that were posted before a later void, and its opposite movement on the void date. The PostgreSQL regression passed: historical assets/P&L 50.00/50.00; void-date P&L -50.00; combined-period P&L 0.00; current assets 0.00. A separate tenant-scoped GL aggregate produced original/reversal revenue +5,000/-5,000 cents and a second-tenant posting did not alter the first report. This is a reporting correction; no journal or accounting policy is changed.

A combined six-test run initially failed the legacy fixed-tenant integrity assertion after reuse of the database with immutable fixture history. Repeating the six tests on a newly created, migrated `cfr010_verify_test` database passed. The earlier database was not reset.

## Cash Flow historical void follow-up — unresolved

On the same disposable PostgreSQL database, the synthetic `historical sale before void` fixture had tenant 6, restaurant 7, entry date 2025-01-01, status `void`, and revenue credit 50.00; its void timestamp was in 2026. Evaluating the operating component's current `je.status = 'posted'` and 2025 date filter returned **0.00**, while an independently scoped General Ledger revenue sum for 2025 returned **50.00**. This confirms that the existing Cash Flow operating query drops the original journal after a later void. The application endpoint and reversal-period behavior were not separately tested in this follow-up, and no Cash Flow correction was implemented. Cash Flow also needs its underlying classification and accounting treatment reviewed before its output can be accepted as a cash flow statement. Do not treat this SQL observation as full report acceptance.
