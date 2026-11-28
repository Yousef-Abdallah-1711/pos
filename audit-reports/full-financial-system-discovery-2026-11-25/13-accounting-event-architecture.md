# 13 - Accounting Event Architecture (durable paid-order outbox)

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Read-only static analysis; nothing was executed against a database. Items not proven from source are labelled UNVERIFIED.

## 1. Summary

There is exactly ONE durable event type: `paid_order` (CHECK constraint `event_type = 'paid_order'`, `backend/migrations/181_durable_accounting_event_outbox.sql:10`). It is an intent row ("post sales + COGS for this order") written in the same transaction that flips `orders.payment_status` to `paid`. A background worker later turns the intent into two journal entries (`auto_sales`, `auto_cogs`). Every other posting path in the system (invoices, payments, payroll, manual/recurring/imported journals) is synchronous and does NOT use the outbox (see 03 section 6).

## 2. Table: `accounting_event_outbox` (migration 181, 183)

| Column | Notes (181:7-26, 183:2-3) |
|---|---|
| `id BIGSERIAL PK`, `event_id VARCHAR(160) UNIQUE` | `event_id = "paid-order:{tenant}:{restaurant}:{order}"` (accounting_event_repo.go:60) |
| `event_type` | only `paid_order` |
| `tenant_id, restaurant_id, order_id` | `order_id` FK `orders(id) ON DELETE RESTRICT`; tenant/restaurant are plain BIGINT (no FK to tenants) |
| `posted_by` | actor user id captured at flip time (nullable) |
| `status` | `pending, processing, retryable, completed, failed` (CHECK) |
| `attempts`, `max_attempts` (default 8) | |
| `next_attempt_at`, `lease_expires_at`, `lease_generation` (183) | fencing token |
| `last_error` (<=1000 chars), `created_at`, `updated_at`, `completed_at` | |
| `UNIQUE (event_type, tenant_id, restaurant_id, order_id)` | dedup key #2 |

Indexes: ready (`next_attempt_at,id` where pending/retryable), stale (`lease_expires_at` where processing), scope (`tenant,restaurant,order,status`) (181:28-35); failed-scope (188:4-6); orders partial index for the recovery scan `idx_orders_paid_accounting_recovery` (182_index_paid_orders_for_accounting_recovery.sql:2-4).

Related tables: `accounting_inventory_movement_events` (184; PK `(tenant,restaurant,order,product)` - replay guard for COGS inventory movements) and `accounting_event_redrive_audit` (188; append-only audit of operator redrives, `reason` 1-1000 chars, `actor_user_id > 0`).

No RLS policy is created on any of these tables (grep of migrations 181/183/185/188 finds none); tenant scoping is by application predicates only.

## 3. Producers

| Producer | Where | Atomicity |
|---|---|---|
| Payment-status flip to `paid` | `OrderRepository.UpdatePaymentStatusAs` order_repo.go:545-603; enqueue at 593-597 via `EnqueuePaidOrderTx` (accounting_event_repo.go:55-71, `ON CONFLICT DO NOTHING`) | SAME tx as `UPDATE orders` and the `FOR UPDATE` order lock (order_repo.go:558-564). Atomic. |
| Migration 181 backfill | 181:37-42 inserts an intent for every already-paid order | one-off |
| Reconciler | `ReconcilePaidOrders` accounting_event_repo.go:74-122 inserts missing intents for any `paid` order each cycle | own short tx (10s timeout) |
| Redrive | `RedriveFailed` accounting_event_admin_repo.go:98-122 | re-queues an existing row (see section 6) |

Only API route that can set `payment_status='paid'`: `PUT /api/v1/admin/orders/{id}/payment-status` (main.go:1437; module ORDERS(10) WRITE) -> `AdminOrderHandler.UpdatePaymentStatus` admin_order_handler.go:466-497 -> `OrderUseCase.UpdatePaymentStatusAs` order_usecase.go:455. Order creation hard-codes `PaymentStatus: "pending"` (order_usecase.go:101). In the frontend only `posApi.updateOrderPaymentStatus` calls it (posApi.ts:266), and only for `payment_method === 'cash'` (pos/page.tsx:302-308; usePOSOffline.ts:179-181). Consequences (verified from source): POS card sales, storefront cash orders and dine-in cash orders have no UI path to `paid`, so they never create an outbox event or revenue journal. A POS cash auto-mark failure is swallowed (`debugError` only, pos/page.tsx:305-307).

## 4. Worker / consumer

- Started at boot when `db != nil`: `service.NewAccountingEventWorker(orderUC).Start(rootCtx)` main.go:513-517. Every API replica runs one (no leader election; safety relies on `FOR UPDATE SKIP LOCKED` + leases).
- Loop: `AccountingEventWorker.Start/runOnce` accounting_event_worker.go:24-58. Interval 15s. Each tick: `ReconcileAccountingEvents` (repair missing intents + metrics + alert log), then drain `ProcessNextAccountingEvent` until idle or first error. A processing error stops the drain for that tick (line 55).
- Claim: `ClaimNext` accounting_event_repo.go:125-181. Marks expired-lease exhausted rows `failed`; selects one row `pending|retryable` with `next_attempt_at <= now` OR `processing` with expired lease, `attempts < max_attempts`, `FOR UPDATE SKIP LOCKED`; sets `processing`, `attempts+1`, `lease_generation+1`, lease 1 minute (const `accountingEventLease`, line 18).
- Process: `OrderUseCase.ProcessNextAccountingEvent` order_usecase.go:482-512 -> loads order (`GetOrder`, operational pool) -> `AutoPostPaidOrder` (469-479) which calls `PostSalesEntry` and `PostCOGSEntry` and `errors.Join`s both -> on success `Complete` (183-193), on failure `Fail` (195-214).
- Posting executes on the accounting pool (pos_accounting) through `GeneralLedgerRepository.CreateAndPostAutoEntry` (general_ledger_repo.go:37, tx variant 59-134); inventory movement rows are written on the operational pool (`inventoryMovementRepo := NewInventoryMovementRepository(db)` main.go:197).
- Outbox and the order read use the OPERATIONAL pool (`NewAccountingEventRepository(db)` main.go:331). Journal/GL writes use the ACCOUNTING pool. The "claim -> post -> complete" sequence therefore spans two connections/roles and cannot be one transaction.

## 5. Retry, backoff, dead-letter

- `Fail`: status `retryable` while `attempts < max_attempts (8)`, else `failed`; backoff `1 minute * 2^min(attempts-1, 6)` = 1,2,4,8,16,32,64,64 min (accounting_event_repo.go:198-206). Error text truncated to 1000 chars.
- Lease expiry (worker crash): row becomes re-claimable after 1 minute; if `attempts >= max_attempts` at expiry it is flipped to `failed` with `worker lease expired after final attempt` (134-142).
- Fencing: `Complete` and `Fail` require `status='processing' AND lease_generation=$n` (190, 211); a stale worker gets "lease was lost" (216-227) and cannot overwrite a newer claim.
- Dead letter = `status='failed'` rows. No automatic alerting beyond `log.Printf` ALERT lines and `observability.RecordAccountingEvents` metrics (order_usecase.go:515-536). Alert delivery destination UNVERIFIED.
- No poison-message classification: a permanent error (closed fiscal period, missing account, disabled rule) consumes all 8 attempts (~2 hours) before `failed`.

## 6. Redrive / recovery (migrations 185, 188; handler)

- `GET /api/v1/accounting/paid-order-events/failed` -> `ListFailed` (module ACCOUNTING 15 READ); `POST /api/v1/accounting/paid-order-events/{id}/redrive` -> `RedriveFailed` (ACCOUNTING 15 WRITE) (main.go:1102-1103; handler accounting_event_admin_handler.go:25-98).
- Redrive requires non-empty `reason` <= 1000 bytes, an authenticated user id, scoped to caller's tenant/restaurant; locks event+order `FOR UPDATE`; only when event is `failed` AND order is still `paid`; inserts an audit row then sets `status='retryable', attempts=0, next_attempt_at=now` (accounting_event_admin_repo.go:98-170). Transactional.
- Authorization: any role with ACCOUNTING WRITE can redrive; there is no separate recovery/approve permission and no two-person control. No frontend screen exists (no reference found under frontend/apps).
- Grants: migration 185 gives `pos_app` SELECT/INSERT/UPDATE (no DELETE) on the outbox and SELECT/INSERT on `accounting_inventory_movement_events`; 188 gives `pos_app` SELECT/INSERT on the audit table and revokes UPDATE/DELETE/TRUNCATE (immutable by grant, not by trigger - no append-only trigger seen UNVERIFIED).
- Migration 182 (`182_accounting_runtime_roles.sql`) separately moves journal/GL write rights to `pos_accounting` and REVOKEs INSERT/UPDATE/DELETE on `journal_entries, journal_entry_lines, general_ledger` from `pos_app` (lines 47). 189 lets `pos_accounting` update only payment columns of `salaries` and insert `hr_audit_logs` for the payroll transaction. 190 adds the normal-balance history guard triggers.

## 7. Idempotency and dedup keys (layered)

1. Outbox `event_id` UNIQUE + `uq_accounting_event_order` UNIQUE (181:9, 25) - one intent per paid order, ever. Note: an order that is paid -> (cannot revert; `CanTransitionOrderPaymentStatus` blocks paid->other, domain/order.go:226) never produces a second intent; refund has no posting path.
2. Posting-level: `CreateAndPostAutoEntryTx` takes `pg_advisory_xact_lock(hashtextextended('auto-posting:{tenant}:{restaurant}:{ref_type}:{ref_id}:{entry_type}'))` (general_ledger_repo.go:76-79), then selects existing entries for that source: `posted` -> returns existing id (idempotent success); `draft` -> error "requires recovery"; `void` -> error; >1 -> error (lines 80-124). There is NO database unique constraint on `(tenant, restaurant, reference_type, reference_id, entry_type)` (only `UNIQUE(tenant,restaurant,entry_number)` in 137_create_journal_entries.sql:42), so uniqueness depends on the advisory lock + application check.
3. Inventory movement: `CreateOrderSaleMovementIfAbsent` claims `(tenant,restaurant,order,product)` in `accounting_inventory_movement_events`, recognises legacy movements, then inserts (inventory_movement_repo.go:28-75). Executed AFTER the COGS journal posts (auto_posting_engine.go:253), in its own tx on the operational pool.
4. Partial-failure replay: Sales and COGS are separate journals with separate dedup keys; if sales posts and COGS fails the event is `retryable`; the replay returns the existing sales entry id and retries COGS.

## 8. Atomic vs eventually consistent

| Operation | Consistency |
|---|---|
| Order `payment_status=paid` + outbox intent | Atomic (single tx, operational pool) |
| Intent -> sales/COGS journals + GL rows | Eventually consistent (worker, >= 0-15s typical; up to ~2h on retries; indefinite when `failed`) |
| Sales journal + COGS journal | Two separate transactions (not atomic with each other) |
| COGS journal + inventory movement | Two separate transactions on two pools; replay-safe, not atomic |
| Journal header + lines + GL + status | Atomic (`PostEntryTx` journal_entry_repo.go:358-425; DB triggers cfr002 enforce balance/immutability, 180) |
| Invoice/payment + its journal | Atomic (single accounting-pool tx) |
| Payroll mark-paid + payroll journal | Atomic (handler-owned tx on accounting pool) |
| Reports (P&L etc.) vs orders just paid | Reports read GL (see 03); a paid order is absent from GL until the worker runs |

## 9. Failure modes found (source-proven unless noted)

1. Silent skip -> permanent "completed": if the auto-posting rule is inactive, `ruleAccounts` returns `enabled=false` and `PostSalesEntry`/`PostCOGSEntry` `return err` (nil) (auto_posting_engine.go:59, 152; rule read 593-608). The event then `Complete`s with no journal. Detected only by the reconciler's `CompletedWithoutSalesJournal` counter + log line (accounting_event_repo.go:107-112; order_usecase.go:528-531); nothing repairs it.
2. Entry date = `order.CreatedAt` (auto_posting_engine.go:137, 233), not payment time. If that fiscal period is closed, `ensureFiscalPeriodOpenTx` rejects (journal_entry_repo.go:830-832) every retry; event ends `failed`; redrive cannot help until the period is reopened.
3. `ensureFiscalPeriodOpenTx` allows posting when NO fiscal period row exists (journal_entry_repo.go:825-826 returns nil on `sql.ErrNoRows`) - fail-open.
4. COGS cost drift: unit cost is re-read on each attempt (`GetProductCost`, auto_posting_engine.go:177) but the journal is idempotent while the inventory movement uses the cost at its own execution time; a retry after a cost change can leave journal cost != movement cost (UNVERIFIED which source GetProductCost reads).
5. Worker runs on every replica, and `ClaimNext`/`ReconcilePaidOrders` are not tenant-scoped; the reconcile `INSERT ... SELECT` scans all paid orders every 15s (partial index added in 182).
6. Outbox is on the operational pool whose role `pos_app` is shared with request handlers; compromise of `pos_app` can flip statuses to `completed` without a journal (no trigger prevents it; UNVERIFIED).
7. `PostSupplierInvoiceEntry`, `PostSupplierPaymentEntry`, `PostDepreciationEntries` exist on the engine (auto_posting_engine.go:462-513) but have no callers outside tests (grep) - dead code, not part of the interface `AutoPostingService` (lines 16-20).
8. No idempotency key on POST endpoints that create money records (payments, invoices, orders). The offline POS replay can create an order twice if the follow-up `payment-status` call fails after `createOrder` succeeded (usePOSOffline.ts:176-190: order created, error thrown, queue item retained).

## 10. Test evidence present in repo (not executed here)

`handler/http/accounting_auto_posting_integration_test.go:53` (`UpdatePaymentStatusPersistsDurableAccountingIntent`), `backend/cmd/api/financial_live_*_acceptance_integration_test.go`. Not run in this phase: tested_status = NOT TESTED.
