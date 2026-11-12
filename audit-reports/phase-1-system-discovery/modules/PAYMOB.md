# Module: Paymob / Subscription Billing (Phase 1 audit)

Read-only code analysis, no live calls. Companion: `integrations/PAYMENT-INTEGRATION.md`.

## 1. Components
| Piece | File |
|---|---|
| Client interface | `backend/internal/service/paymob_client.go` |
| Live client (auth token, order, payment key, refund, iframe URL, HMAC) | `backend/internal/service/paymob_client_live.go` |
| Mock client (`VerifyHMAC` always true) | `backend/internal/service/paymob_client_mock.go:51-54` |
| Mock iframe + processor (calls own webhook) | `backend/internal/handler/http/mock_paymob_handler.go` |
| Handler (initiate, webhook, history, current, cancel, refund) | `backend/internal/handler/http/payment_handler.go` |
| Use case (state machine) | `backend/internal/usecase/payment_usecase.go` |
| Retry queue + worker | `backend/internal/service/webhook_retry_worker.go`, migration `117_create_webhook_retry_queue.sql` |
| Expiry worker | `backend/internal/service/subscription_expiry_worker.go` |
| Tables | `115_create_payments_table.sql`, `110`-`113`, `114`, `116`, `118`, `119` |
| Frontend | `frontend/apps/dashboard/src/hooks/usePayment.ts`, `components/providers/TenantStatusProvider.tsx`, settings/billing page |

## 2. Mode switch
`PAYMOB_MODE` env (`backend/cmd/api/main.go:197-211`): `live` selects the real client; anything else (including unset or typos) selects the mock client. Mock endpoints are registered when mode is empty or `mock` (`main.go:427-431,583-588`). `docker-compose.prod.yml:74` passes `PAYMOB_MODE: ${PAYMOB_MODE:-}` and `.env.prod.example:24` ships it empty, so a production deployment that does not set it silently runs in mock mode with a webhook that accepts any body (SEC-004, P0). There is no startup check that live mode has API key, HMAC secret, iframe id and integration id set (empty HMAC secret in live mode fails closed, `paymob_client_live.go:187-189`; empty API key would only fail at first payment).

## 3. Payment flow
1. `POST /subscription/initiate` (auth + tenant, no RBAC, no tenant-status check by design, `main.go:645`). Plan resolved server-side from DB plans or env defaults (`SUBSCRIPTION_MONTHLY_PRICE_CENTS` default 50000, annual 500000, `payment_usecase.go:907-935`); client only sends `plan_id`. Creates `payments` row `pending` with `merchant_order_id = tenant_<id>_<plan>_<unixnano>` (UNIQUE), then Paymob auth -> order -> payment key; stores `paymob_order_id` and the Paymob auth token on the row (`UpdatePaymentOrder`, `payment_repo_impl.go:122`, plaintext in `payments.paymob_auth_token`).
2. User pays in Paymob iframe (or mock iframe).
3. Paymob posts `POST /webhook/paymob` (public, `RateLimitByIP(60)`, `main.go:575-578`). Handler reads raw body, takes HMAC from `?hmac=` or body (`payment_handler.go:136-140`).
4. `HandleWebhook` (`payment_usecase.go:176-365`): verify HMAC -> require `order.id` -> load payment by `paymob_order_id` -> if status already `success` return (idempotent) -> compare stored amount to plan price, webhook amount to plan price, currency -> on success open a DB transaction: update payment `success` + tenant `subscription_status='active'`, `subscription_plan`, `subscription_expires_at = paidAt + plan months`, `last_payment_*`; commit; write `tenant_subscription_history` rows (history writes are after commit and best-effort). On failure/pending flag: mark payment `failed` (a pending transaction is recorded as failed).
5. Errors that are not HMAC/validation errors return 500 and enqueue to `webhook_retry_queue` (backoff 1m, 5m, 30m, 2h, 12h per comment `main.go:406`); worker re-posts to the own endpoint with `X-Webhook-Retry-Attempt: true`.

## 4. State transitions
Payment: `pending -> success | failed`, `success -> refunded` (super admin only), CHECK constraint `('pending','success','failed','refunded')` (`115:40`). There is no transition guard in SQL: `updatePaymentStatusTx` updates by id only (`payment_usecase.go:750-760`), and a partial unique index limits `success` rows per `paymob_order_id` (`115:51-52`).
Tenant: `trial -> active` (webhook), `active -> expired` (expiry worker when `subscription_expires_at <= NOW()`, `subscription_expiry_worker.go:62-71`), `active` + `cancelled_at` (cancel at period end), refund -> `trial` or `expired` (`resolveStatusAfterRefund`, `clearSubscriptionAfterRefundTx` nulls `expires_at`), super-admin `disable/enable/extend-trial`. Allowed statuses: `trial, active, expired, suspended, disabled` (`110:29`).

## 5. Findings summary (full detail in SECURITY-GAPS.md)
- SEC-004 P0: default mock mode + public webhook + `VerifyHMAC` true + mock endpoints. A tenant can call `initiate` (response returns `paymob_order_id`), then POST a crafted webhook with matching `amount_cents`/`currency` (plan prices are public via `GET /subscription/plans`) to activate a paid subscription for free.
- SEC-015 P1: replay/refund re-activation. Idempotency skip only when stored status is `success`; after a refund (status `refunded`) a replayed valid success webhook, or Paymob's own refund callback (same `order.id`, `is_refunded=true`, amount equal) is processed as a fresh success and re-activates the tenant. The handler/use case ignore `is_refunded`, `is_voided`, `is_capture`, `has_parent_transaction`, transaction `type`, `integration_id` and `owner` (`payment_usecase.go:189-252`). The HMAC covers these fields, but no business check uses them.
- SEC-008 P1: tenant-status middleware is not applied to data routes, so `expired`/`disabled`/`suspended` tenants keep API access (`main.go:610-616,681`).
- SEC-016 P2: Paymob auth token persisted in `payments.paymob_auth_token` and card last four/brand stored; billing email/phone stored. No encryption at rest.
- SEC-017 P3: renewal computes `expires_at = paidAt + period`, discarding unused remaining time on early renewal; cancel flag cleared only by expiry worker; refund calls Paymob before the DB transaction (if the DB step fails, money is refunded but tenant stays active, `payment_usecase.go:571-615`); history rows written after commit may be lost.
- Any tenant user (no RBAC) can initiate payments and cancel the subscription (`main.go:645-648`).
- Trial: default 15 days (`110:8`), `112_seed_existing_tenants_trial.sql` backfills; super-admin `extend-trial`.
- Webhook processing returns 200 for body read/JSON errors (`payment_handler.go:120-134`), so malformed deliveries are not retried by Paymob and are not enqueued.
- HMAC correctness: field order and concatenation match Paymob transaction-callback spec (`paymob_client_live.go:21-42,253-287`); comparison constant-time; hex lowercased. HMAC source may come from query or body.
- No webhook replay protection beyond payment status (no stored transaction ids with unique constraint besides success-per-order; no timestamp window).

## 6. Required tests
See `evidence/issues-security.json` (SEC-004, SEC-008, SEC-015, SEC-016, SEC-017): startup refuses `ENV=production` with mock mode; forged webhook returns 401 in live mode; replay after refund is ignored; refund-callback payload does not activate; expired tenant gets 403 on data routes; amount/currency mismatch rejected; concurrent duplicate webhooks yield one activation; renewal extends from current expiry.
