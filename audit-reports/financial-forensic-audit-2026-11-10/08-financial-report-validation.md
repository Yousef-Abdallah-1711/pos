# Financial report validation

**Status: FAILED for one Balance Sheet reconciliation; otherwise PARTIAL.**

Controlled journal on 2026-11-10: Dr Audit Cash 9911 EGP 100; Cr Audit Sales Revenue 9912 EGP 100.

| Endpoint | HTTP/result | Assessment |
|---|---|---|
| /api/v1/accounting/trial-balance | 200; debit 100, credit 100, is_balanced true | Totals reconcile |
| /api/v1/reports/financial/profit-loss?date_from=2026-11-10&date_to=2026-11-10 | 200; revenue 100, expenses 0, net profit 100 | Matches fixture |
| /api/v1/reports/financial/balance-sheet?as_of=2026-11-10 | 200; assets 100, liabilities 0, equity 0, equation_balanced false | Fails equation |
| /api/v1/reports/financial/cash-flow?date_from=2026-11-10&date_to=2026-11-10 | 200; operating/net 100 | Directionally matches cash debit; no broader class tested |

A 200 only establishes handler response. Aging, inventory, VAT and comparative periods were empty or not meaningfully populated. Cutoff, currency, voids, reversals, tax, discounts, refunds and multi-restaurant scopes remain unvalidated.

## Independent fixture extension — 2026-11-11

Four separately dated balanced journal states were posted in the disposable tenant and compared against P&L. Balance Sheet returned:

| As of | Assets | Liabilities | Equity | P&L net profit | Difference (A − L − E) |
|---|---:|---:|---:|---:|---:|
| 2026-11-10 | 100 | 0 | 0 | 100 | 100 |
| 2026-11-11 | 60 | 0 | 0 | 60 | 60 |
| 2026-11-12 | 85 | 25 | 0 | 60 | 60 |
| 2026-11-13 | 135 | 25 | 50 | 60 | 60 |

FIN-001 root cause is confirmed at `backend/internal/usecase/financial_report_usecase.go:261-336`: Balance Sheet query excludes revenue/expense and only sums asset/liability/equity; no current earnings/retained earnings are rolled into equity. The EGP 60 discrepancy persists after subsequent liability, asset, and equity movements. These are controlled fixtures, not production data or an accounting-policy sign-off.

Trial balance was observed balanced for the original EGP 100 journal and the conditional EGP 100/100 normalized probe. P&L was cross-checked for the four fixture dates. Cash flow, receivables, payables, inventory valuation, VAT, comparative periods, and CSV/XLSX export were not independently reconciled. A downloaded PDF exists, but row/amount/filter integrity was not established. See `19-coverage-matrix.md`.
