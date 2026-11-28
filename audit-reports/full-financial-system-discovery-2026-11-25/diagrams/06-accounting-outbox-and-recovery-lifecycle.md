# Diagram 06 - Accounting outbox and recovery lifecycle

All transitions verified in `backend/internal/repository/accounting_event_repo.go` and `accounting_event_admin_repo.go`.

```mermaid
stateDiagram-v2
  [*] --> pending: EnqueuePaidOrderTx in same tx as paid, or ReconcilePaidOrders, or migration 181 backfill
  pending --> processing: ClaimNext (SKIP LOCKED, attempts+1, lease_generation+1, lease 60s)
  retryable --> processing: ClaimNext when next_attempt_at reached
  processing --> completed: Complete (fenced by lease_generation)
  processing --> retryable: Fail while attempts below 8 (backoff 1,2,4,8,16,32,64,64 min)
  processing --> failed: Fail at attempts 8
  processing --> processing: lease expired (worker crash), re-claim after 60s
  processing --> failed: lease expired and attempts at max
  failed --> retryable: POST paid-order-events id redrive (ACCOUNTING WRITE, reason required, audit row, order still paid)
  completed --> [*]
  note right of completed
    Completed does NOT prove a journal exists.
    Inactive posting rule returns nil and the event completes.
    Reconciler only counts CompletedWithoutSalesJournal.
  end note
```

```mermaid
sequenceDiagram
  participant API as UpdatePaymentStatusAs (pos_app)
  participant OB as accounting_event_outbox
  participant WK as Worker (every replica, 15s)
  participant UC as OrderUseCase
  participant EN as AutoPostingEngine
  participant GL as Ledger (pos_accounting)
  API->>OB: BEGIN; lock order; UPDATE paid; INSERT intent ON CONFLICT DO NOTHING; COMMIT
  WK->>OB: ReconcilePaidOrders (insert missing intents, count states)
  WK->>OB: ClaimNext
  WK->>UC: ProcessNextAccountingEvent
  UC->>EN: PostSalesEntry
  EN->>GL: advisory lock; existing posted returns id, else create and post
  UC->>EN: PostCOGSEntry
  EN->>GL: create and post COGS journal
  EN->>OB: (pos_app) inventory movement + movement-event key
  UC->>OB: Complete or Fail (lease fenced)
```
