# 06 — Financial event catalog (PROPOSED)

Every event has tenant/restaurant, aggregate and event IDs, effect type/sequence, schema version, occurred/recorded/effective timestamps, currency, actor, correlation/causation, immutable tax/discount/FX snapshots, idempotency identity and correction linkage. Operational state changes alone do not imply accounting recognition; recognition policy requires owner/accountant approval.

| Event | Example debit / credit | Required control |
|---|---|---|
| Sale recognized | AR/cash/clearing Dr; revenue/tax payable Cr | approved fulfillment trigger; order-to-journal reconciliation |
| Capture/settlement | clearing Dr; tender receivable/bank Cr; explicit fees | provider ID uniqueness; settlement statement |
| Split/failed tender | confirmed tender only; failed auth has no revenue by itself | replay-safe captures |
| Discount/refund/chargeback | contra-revenue/tax/refund payable Dr; tender/AR Cr | linked original, cumulative cap |
| Tip/service charge/gift card/deposit | cash Dr; payable/deferred revenue Cr until earned | liability roll-forward |
| Cancellation | reverse recognized effects only if already recognized | linked compensating event |
| Inventory receipt/consumption/waste | inventory Dr/AP Cr; COGS/waste Dr/inventory Cr | cost method and movement reconciliation |
| Supplier/customer invoice/payment | expense/inventory Dr/AP Cr; AP Dr/bank Cr; AR Dr/revenue/tax Cr | subledger control-account match |
| Payroll/depreciation | expense Dr; payable/accumulated depreciation Cr | approved batch/schedule |
| Bank reconciliation/fiscal close | explicit adjustment/close transfer | no balancing plug; authorized audit |

Provider calls remain outside database transactions. Failures have durable pending/failed state, retry and alert; correction is compensating, not deletion. Currency, source, tax and reversal policy attach per effect.
