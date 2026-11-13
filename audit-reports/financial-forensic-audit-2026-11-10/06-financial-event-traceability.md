# Financial event traceability

**Status: PARTIAL.** One manual journal traced end-to-end. Other event paths are source observations only.

| Event | Expected effect | Evidence / gap |
|---|---|---|
| Manual journal | draft, validate, post, ledger | Runtime verified; UI omitted JSONB attachments and failed. |
| Order sale | revenue/tax/discount and tender cash/receivable | PostSalesEntry exists; order_usecase.go:173 uses fixed 10% tax. No order fixture. |
| COGS/inventory | Dr COGS / Cr inventory on sale | PostCOGSEntry exists; no product/stock fixture. |
| Customer invoice/payment | receivable/revenue and settlement | APIs/repositories exist; no runtime invoice. |
| Supplier invoice/payment | expense/inventory payable, then payable/cash | Supplier repository uses CreateJournalEntryTx; no runtime invoice. |
| Refund/chargeback | linked reversal with audit trail | Not executed; reversal and idempotency semantics unknown. |
| Webhook/settlement | idempotent event and reconciliation | Not fully traced or runtime-tested. |

POST journal with attachments [] created a draft; UI post moved it to posted and produced two ledger rows. Trial balance balanced; P&L reported revenue 100; Balance Sheet equation failed. This establishes one case only, not all historical periods or classifications.

## Continuation evidence — 2026-11-11

Source trace shows paid-order accounting is called after payment status changes in `admin_order_handler.go:500-507` and from public order creation in `public_order_handler.go:82-84`. `OrderUseCase.AutoPostPaidOrder` calls sales and COGS posting separately and only logs errors (`order_usecase.go:492-506`); `AutoPostPaidOrderByID` also logs reload errors and returns (`:508-520`). The handlers then return success based on the order/payment operation. No synthetic paid-order API was run, so this establishes silent-error handling by source inspection only, not a paid-order ledger outcome.

Actual supplier auto-post runtime probe failed before persistence on nil JSONB attachments. With valid JSON, supplier entry type width also blocks insertion. Other supplier/customer repository paths remain separate and were not exercised. Workflow status is consolidated in `19-coverage-matrix.md`.
