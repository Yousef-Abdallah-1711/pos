# WS03 — Payment Callback Integrity Report

**Status:** PARTIAL — core PostgreSQL callback/refund races are verified; provider/local reconciliation remains open.
**Date:** 2026-11-02
**Provider safety:** Only mock clients and synthetic payloads were used. No Paymob API call or real payment was made.

## Initial state

Paymob HMAC verification and a webhook test file already existed. Existing tests covered invalid/missing signatures, unknown orders, duplicate/late events, refund/void non-activation, and currency mismatch using in-memory repositories. Success processing already grouped payment success and tenant subscription activation in one SQL transaction. However, the transactional status update matched only the payment ID, and the failed callback path also updated by ID without a prior-state condition. Two callbacks could both read a pending payment; a late failed callback could overwrite a success after it committed. Refund handling checked the payment before calling the provider, but concurrent refund requests could both pass that check and invoke the provider.

## Changes implemented

- Transactional success transitions now update only when the stored payment is `pending` or `failed`. A duplicate after another success returns the existing idempotent “already processed” outcome and does not repeat subscription activation.
- Transactional refund transitions now require the stored state to still be `success`.
- Non-transactional failed callback updates now match only `pending`. If a concurrent request already moved the payment to `failed`, `success`, or `refunded`, the stale failure callback is acknowledged without downgrading that state.
- Refund processing now opens a DB transaction, locks the payment row with `SELECT ... FOR UPDATE`, rechecks that the state is still `success`, and only then calls the provider. The row remains locked until payment and tenant refund state commit, so another request cannot call the provider for the same already-refunded payment.
- Refund and void webhooks remain acknowledged without activating or extending subscriptions; they are still marked for manual reconciliation by the existing behavior.

## Verification

- Existing `payment_webhook_security_test.go` cases passed, including HMAC, duplicate/late callbacks, refund/void non-activation and amount/currency integrity.
- Added SQL-mock checks for conditional success processing, stale failure callback handling, and suppressing a provider refund call when a concurrent request has already committed the refund.
- Targeted repository, use case and HTTP handler suites passed.
- PostgreSQL 15.19 parsed guarded success/refund updates and the refund row-lock query using `PREPARE` on the isolated restored DB; the session rolled back.
- Final backend gate after WS02/WS03 changes: `gofmt -d`, `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` all passed in a network-isolated Go 1.24 container with source and module cache read-only during verification.

## Files changed

- `backend/internal/repository/payment_repo_impl.go`
- `backend/internal/repository/payment_repo_test.go` (new)
- `backend/internal/usecase/payment_usecase.go`
- `backend/internal/usecase/payment_usecase_test.go`
- `backend/internal/usecase/payment_webhook_security_test.go`
- `docs-organized/phase-2a/execution/PHASE_2A_EXECUTION_TRACKER.md`

## Remaining risks and limitations

- SQL-mock and PREPARE checks do not prove transaction behavior under real concurrent PostgreSQL sessions. Add an isolated integration test that races duplicate success callbacks and duplicate refund requests and asserts exactly one transition/provider call.
- The refund transaction holds a database row lock during the provider refund request. This serializes duplicate requests but consumes a database connection and lock while waiting for the provider. A durable refund-claim/idempotency design with reconciliation would be preferable for production scale.
- The provider refund may succeed while a later local transaction/commit fails. A database and external provider cannot be committed atomically; this needs an idempotent provider request key or a durable reconciliation workflow before production refund automation is considered fully safe.
- Refund/void callbacks do not reconcile local payment/subscription state; they are acknowledged without activation and require manual reconciliation.
- No hosted CI or live provider verification was performed.

## Closure decision

The callback transition and concurrent refund guards are implemented and unit-verified. WS03 remains **IMPLEMENTED BUT UNVERIFIED** until a disposable PostgreSQL concurrency test and refund reconciliation behavior are verified. Production payment processing remains outside this report's validation scope.

## Closure execution update — 2026-11-02

### Payment state and transaction evidence

| Event | Current accepted behavior | Atomicity / limit |
|---|---|---|
| Pending payment creation | `InitiatePayment` inserts local `pending` row before provider setup | Provider calls are outside a DB transaction; failed setup attempts mark local payment failed. |
| Valid successful callback | Pending or failed may transition to success; tenant subscription becomes active | Payment and tenant update share one PostgreSQL transaction. Conditional update makes duplicate success a no-op. |
| Failed callback | Only a payment still pending may become failed | Conditional SQL prevents late failure from downgrading success/refund. |
| Duplicate success | A stored success is acknowledged; concurrent duplicate losers see the conditional update miss | No duplicate activation transaction; subscription history is written after commit. |
| Delayed success after refund | Refunded payment is terminal; event is ignored | No subscription reactivation. |
| Refund callback / void callback | Acknowledged without local state reconciliation | Manual reconciliation remains an operational gap. |
| Admin refund | Requires local success; locks payment `FOR UPDATE`, calls provider, writes refund+tenant state in one transaction | Provider call occurs while holding DB lock. Provider success followed by local commit failure has no durable recovery path. |

### New PostgreSQL tests and defect remediation

Added `backend/internal/usecase/payment_postgres_integration_test.go`. It uses migrated synthetic tenant/payment rows and an in-process mock provider. No provider network call or production credentials are used.

- Twelve simultaneous success callbacks produce one effective DB transition; payment is `success`, tenant is `active`, and exactly two success/activation history entries are observed.
- Invalid HMAC leaves payment pending and tenant trial.
- Failure→success→delayed failure leaves payment successful and tenant active; a delayed success after refund leaves payment refunded and tenant expired.
- A controlled tenant constraint failure after the payment update rolls back both state changes; payment remains pending and tenant remains trial.
- Two concurrent refund requests synchronize before the DB lock; one provider refund call and one local refund commit occur. Payment ends `refunded`, tenant ends `expired`.
- The refund race first exposed a PostgreSQL parameter type inference failure in `clearSubscriptionAfterRefundTx`, causing the mock provider to be called twice while both DB transactions failed. The update now uses a distinct SQL statement for trial versus expired state, and the focused race test passes 10 consecutive runs.
- Three repeated targeted runs of all 8 database integration cases passed.

The new DB tests require an explicit `POS_TEST_DATABASE_URL` and skip when it is absent.

### Current closure decision

WS03 is **PARTIAL**, not complete. Core success/failure/refund database race and rollback behavior has evidence, but refund/void callbacks do not reconcile provider state, local commit failure after provider success has no durable recovery path, and the wider cancellation/refund ordering matrix remains open. Provider idempotency semantics cannot be assumed or verified by the mock boundary.
