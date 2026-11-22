# Root Cause

Two separate findings must not be conflated.

1. **The original 5,000-cent Balance Sheet residual is an accounting-policy question.** `GenerateBalanceSheet` reports posted asset, liability, and equity accounts. It does not derive unclosed P&L earnings into equity. The posted source journal debits cash 5,500 cents and credits VAT payable 500 cents plus revenue 5,000 cents. Thus the Balance Sheet is 5,000 cents short of liabilities plus equity while P&L reports 5,000 cents net income. The arithmetic is exact. OPEN-01, EARN-01, and CLOSE-01 remain unapproved, so this work does not move earnings into equity or post a close.

2. **A confirmed engineering defect was stale server-side report caching.** `FinancialReportHandler` cached Balance Sheet, P&L, Cash Flow, and Trial Balance responses for five minutes using tenant, restaurant, path, and query string. Ledger writes did not invalidate these arbitrary report keys. An authenticated request after a posting could therefore return old amounts and old reconciliation status. A regression test failed before the correction and passed after it.

The remediation bypasses caching only for those four ledger-derived reports. The separate inventory valuation cache remains. The unreconciled warning and PDF eligibility gate remain unchanged.

3. **A separate historical P&L defect was found during void testing.** A posted journal later marked void disappeared from P&L even for a reporting period before the void, although its original GL movement remained in the Balance Sheet. The P&L query filtered on the journal's current status. The report correction preserves the original entry in its entry-date period and includes its inverse in the void-date period. See `POSTGRESQL_RECONCILIATION.md` for the failing and passing PostgreSQL evidence. This correction does not resolve the separate unapproved earnings presentation.
