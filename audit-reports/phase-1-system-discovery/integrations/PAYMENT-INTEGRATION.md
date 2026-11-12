# Integration: Paymob (Payment Gateway) - Phase 1 audit

Read-only. No outbound calls made. State machine detail is in `modules/PAYMOB.md`.

## External API surface used (`backend/internal/service/paymob_client_live.go`)
| Step | Endpoint (base default `https://accept.paymob.com/api`) | Notes |
|---|---|---|
| Auth | `POST /auth/tokens` with `api_key` | token reused for order, key, refund; stored in DB (SEC-016) |
| Order | `POST /ecommerce/orders` (amount_cents, currency, merchant_order_id, one hard-coded item "POS SaaS Monthly Subscription") | item name is "Monthly" even for annual plan (`:99`); booleans sent as strings (`:91-96`) |
| Payment key | `POST /acceptance/payment_keys` (billing_data defaults: "NA", `+200000000000`, Cairo, EG; `lock_order_when_paid: "true"`) | 30 s HTTP timeout, no retries/idempotency key on outbound calls |
| Iframe | `GET {base}/acceptance/iframes/{PAYMOB_IFRAME_ID}?payment_token=...` | |
| Refund | `POST /acceptance/void_refund/refund` (transaction_id, amount_cents) | super-admin only route `POST /super-admin/payments/{id}/refund` |
| Webhook in | `POST /api/v1/webhook/paymob` (transaction processed callback) | HMAC-SHA512 over 20 ordered fields (`:21-42,253-287`) |

Config env: `PAYMOB_MODE`, `PAYMOB_API_KEY`, `PAYMOB_INTEGRATION_ID`, `PAYMOB_IFRAME_ID`, `PAYMOB_HMAC_SECRET`, `PAYMOB_BASE_URL`, `WEBHOOK_BASE_URL`/`APP_BASE_URL` (used by retry worker), `SUBSCRIPTION_*_PRICE_CENTS`, `SUBSCRIPTION_CURRENCY`. `PAYMOB_BASE_URL` is operator supplied (no allowlist), so a bad value would send the API key to an arbitrary host (config-level risk, low).

## Webhook security checklist
| Control | Status | Evidence |
|---|---|---|
| Signature verified before any state read | Yes | `payment_usecase.go:185-187` |
| Constant-time compare | Yes | `paymob_client_live.go:203` |
| Fails closed on missing secret/hmac | Yes (live) | `:187-189` |
| Mock mode bypass | `VerifyHMAC` returns true | `paymob_client_mock.go:51-54`; default when `PAYMOB_MODE` unset (`main.go:197-211`) (SEC-004) |
| Amount and currency re-checked server-side | Yes (against plan price and stored payment) | `payment_usecase.go:217-232` |
| Idempotency | Only `status == success` short-circuit; DB partial unique index on success per `paymob_order_id` | `:202-205`, `115:51-52` |
| Replay/refund-callback protection | No (SEC-015) | `:189-252` ignore `is_refunded/is_voided/type` |
| Source IP allowlist / timestamp window | No | handler |
| Rate limit | 60/min per IP, in-memory | `main.go:575-578` |
| Failure semantics | 401 bad HMAC, 400 validation, 500 + enqueue for transient; body/JSON parse errors return 200 | `payment_handler.go:120-207` |
| Observability | `PaymentWebhookMonitor` (200 events in memory), error tracker, super-admin `/payment-metrics`, `/webhook-health` | `main.go:396-425,1775-1780` |

## Retry queue
`webhook_retry_queue` (migration 117) stores endpoint and raw body (including the original HMAC) of failed deliveries; `WebhookRetryWorker` re-posts to `baseURL` (`WEBHOOK_BASE_URL` or `APP_BASE_URL`, mandatory in production, `main.go:407-416`) with header `X-Webhook-Retry-Attempt`. Retried requests pass through the same HMAC check, so the queue cannot be used to inject forged payloads. In mock mode HMAC is bypassed anyway. Raw payloads (card last four, billing details) persist in the queue table (retention not reviewed).

## Mock endpoints
`GET /api/v1/mock/paymob/iframe`, `POST /api/v1/mock/paymob/process` (`main.go:583-588`). `ProcessPayment` builds the webhook URL from the incoming `Host` header (`mock_paymob_handler.go:243-251`) and uses HMAC literal `mock_hmac`; success vs failure chosen by card number. Exposed in any deployment without `PAYMOB_MODE=live`. Redirect target sanitising (`normalizeBillingRedirectURL`, `:288`) was not fully reviewed (open-redirect unverified).

## Money handling
Amounts as integer cents end to end (`amount_cents`, `last_payment_amount` int), which is correct; currency default EGP. Plan price changes are not versioned on the payment row except `amount_cents` (stored) so historical amounts remain accurate.

## Gaps for Phase 2
- Contract tests with recorded real Paymob payloads (success, failed, pending, refund, void, 3DS) against `HandleWebhook`.
- Reconciliation job comparing `payments` to Paymob transaction list (none exists).
- Confirm Paymob callback includes `type: TRANSACTION` and how refund callbacks look for this account (unverified).
- Frontend: `usePayment.ts` flow (iframe embed, redirect statuses `billing_status` query) not audited in depth.
