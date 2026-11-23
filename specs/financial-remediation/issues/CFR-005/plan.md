# CFR-005 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/handler/http/admin_order_handler.go`
- `backend/internal/repository/order_repo.go`
- `backend/internal/repository/accounting_event_repo.go`
- `backend/internal/repository/order_payment_transition_integration_test.go`
- `backend/internal/domain/order.go`

## Approach

1. Verify the merged payment guard and paid-order outbox against current main, and rerun transition and worker tests against a clean disposable database.
2. Inventory order channels, provider capture/refund callbacks and event identities; distinguish subscription payments from restaurant orders.
3. Define provider idempotency, uncertain outcome recovery and durable state separately from accounting allocation and reversal policy. Preserve payment and journal history; specify migration and rollback before implementation.
4. Once decisions are authorized, implement and test full/partial refund, duplicate and out-of-order callbacks, concurrent requests, timeout, restart, tenant isolation and linked reversal. Reconcile against an independent oracle and validate API/operator flows on the reviewed candidate and main.

## Gaps

Provider contract, recognition point, refund allocation, tax and closed-period treatment remain unresolved; the current guard is not a refund implementation.
