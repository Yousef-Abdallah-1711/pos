# 02 — Financial source-of-truth contract (PROPOSED)

## Authority ownership

| Fact | Authority owner | Copy/projection that is not authority |
|---|---|---|
| Accepted items, quantities, agreed price, discount, tax quote, fees | Immutable accepted order revision | POS cart, receipt preview, analytics |
| Tender authorization/capture/failure/refund/settlement | Payment aggregate plus immutable payment events; provider-signed callback and settlement statement are external evidence | Order payment flag or redirect success |
| Customer receivable / supplier payable | Approved invoice subledger with immutable lines and source references | Aging screen and invoice UI total |
| Physical stock receipt/consumption/waste/return | Append-only inventory movement | Product stock field and valuation screen |
| Recognized account/period/functional amount | Posted balanced journal with source-event identity | Draft or mapping preview |
| Account activity | Posted journal lines; GL is a one-to-one posting projection | current_balance and running_balance |
| Statement at cutoff | Versioned report from a consistent set of posted journals and policy | Page/export cache |
| External cash settlement | Bank/provider statement and reconciliation result | Payment callback alone |

IFRS 15 describes revenue reporting according to transfer of promised goods/services; IAS 2 addresses inventory cost and expense recognition; IAS 1 covers financial-statement presentation including changes in equity; IAS 21 addresses foreign currency. Applicability and local statutory treatment are ACCOUNTANT APPROVAL REQUIRED, not certified here. [IFRS 15](https://www.ifrs.org/issued-standards/list-of-standards/ifrs-15-revenue-from-contracts-with-customers/), [IAS 2](https://www.ifrs.org/issued-standards/list-of-standards/ias-2-inventories/), [IAS 1](https://www.ifrs.org/issued-standards/list-of-standards/ias-1-presentation-of-financial-statements/), [IAS 21](https://www.ifrs.org/issued-standards/list-of-standards/ias-21-the-effects-of-changes-in-foreign-exchange-rates/).

## Event identity and provenance

Each financially relevant transition carries tenant_id, restaurant_id, aggregate type/id, source_event_id, event type/schema version, occurred_at, recorded_at, effective_at, timezone, currency, actor/system principal, correlation and causation IDs, immutable pricing/tax/discount/FX snapshots, and correction linkage. Provider events also retain provider account/event ID, transaction ID, signature-verification result and receive time.

Accounting effect identity is unique per tenant, restaurant, source event, effect type and sequence. Human reference numbers or timestamps alone are not idempotency keys. A single order may create sale, COGS, tender, fee, refund and restock effects.

A canonical request fingerprint is computed after validation. Same key plus same fingerprint returns the existing outcome. Same key with different semantic payload is rejected and audited.

## Mutable versus immutable

- Before acceptance, cart and quote can change.
- After acceptance, order pricing/tax/discount snapshot is immutable; amendment or cancellation appends a new event.
- Provider events, captured tender facts, settlement lines, inventory movements and posted journals are append-only.
- Draft journals can change with version checks; posted entries cannot be edited or deleted.
- Invoice correction uses credit/debit note or pre-posting void per approved policy.
- Reports can be rebuilt, but each artifact retains input watermark, parameters, policy version, generation time and manifest/hash.

## Propagation

1. Quote calculation consumes versioned price/tax/discount rules and returns exact decimal strings and currency.
2. Accepted order stores item snapshots; product price changes do not restate it.
3. Payment records actual tender, not accounting revenue.
4. Inventory movement records quantity and cost method snapshot.
5. Accounting builder consumes source event and approved mapping to emit balanced lines.
6. Posted journal and GL are accounting authority; reports derive from them.
7. Settlement compares captures to processor/bank results; fees, chargebacks and timing differences are explicit events.

## Consistency and reconciliation

Source transition and outbox record commit in the same local transaction. An accounting event is posted, pending/retryable, or an operator exception; no silent disappearance is allowed. Reports distinguish pending accounting and bind to a posted watermark. Daily controls compare order net/tax/tips/discount by tender to journals; provider captures/refunds/settlements to clearing and bank; AR/AP to control accounts; inventory movement cost to inventory control; GL lines to journals; reports to balances.

Corrections preserve original evidence and add linked compensating events. Retention and approvals are OWNER DECISION REQUIRED / ACCOUNTANT APPROVAL REQUIRED.

## Existing gap and acceptance

Current code stores order amounts, order payment status, subscription payments, accounting receipts/payments, invoice balances, stock, journals, GL and reports separately. Prior audit tested one manual journal, not all order-to-ledger paths.

Acceptance: 100 concurrent same-event deliveries yield one effect; changed payload under reused key is rejected; each eligible event has posted/pending/error disposition; reports reconcile at declared watermark; corrections never rewrite source/provider/posting history.
