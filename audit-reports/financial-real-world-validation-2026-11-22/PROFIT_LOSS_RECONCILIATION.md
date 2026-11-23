# Profit and Loss reconciliation

For the synthetic paid order, the source-derived net revenue expectation was 5,000 cents, tax was 500 cents, and no discount, delivery, or expense was present. The API acceptance observed P&L/source revenue and GL revenue of 5,000 cents and expense of zero. Sales gross exports correctly reported the order gross as 5,500 cents, which includes tax and is not the same measure as revenue.

The exercised scenario reconciled these values exactly. Browser journeys also navigated P&L and compared selected UI/API/report amounts. The evidence does not establish full P&L completeness for all accounts, period boundaries, or COGS. Inventory/COGS workflows and original operating-period results were not reconciled. See `financial-live-api-tests.jsonl` and `financial-browser-tests.txt`.
