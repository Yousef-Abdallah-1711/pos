# Module: Trial and Subscription Lifecycle

Sources: frontend (landing, restaurant-website, super-admin) plus backend usecases/repos/migrations. Static analysis only. Issue IDs: `evidence/issues-frontend-others.json`.

## 1. Data model
`tenants` columns added by migrations 110-119, 114: `subscription_status` (trial|active|expired|suspended|disabled, CHECK constraint), `trial_start_date`, `trial_end_date`, `trial_days` (default 15), `is_disabled/disabled_at/disabled_by/disabled_reason/enabled_at`, `grace_period_days` (30), `deletion_scheduled_at`, plus payment fields (`subscription_plan`, `subscription_expires_at`, `last_payment_at/amount/currency`) and cancellation fields (119). Tables: `tenant_subscription_history` (111/113), `payments` (115; a second `143_create_payments.sql` also exists - not reconciled), `webhook_retry_queue` (117), `subscription_plans` (118, seeds `monthly` 50000 cents EGP and `annual` 500000 cents).
Go constants: StatusTrial/Active/Expired/Suspended/Disabled (`domain/subscription.go:9-13`).

## 2. Lifecycle as implemented
1. Register (`AuthUseCase.Register` -> `AuthRepository.CreateTenantWithUser`): inserts tenant with `subscription_plan='free'`, `max_restaurants=1`, `max_users=5`, `is_active=true`; **no trial dates** (auth_repo.go:61-63). `subscription_status` defaults to `trial`, `trial_days` to 15, `trial_start/end` stay NULL (FE-OTH-003).
2. Trial check: `CheckTenantAccess` blocks only when `is_disabled`, or `status=='trial' && trial_end_date IS NOT NULL && now > trial_end_date`. NULL end date = never expires. `GetTenantSubscription` (used by `/tenant/status` and public `/tenant-status/{slug}`) additionally blocks `status=='expired'` and reports `DaysRemaining`.
3. Enforcement points: dashboard UI (tenantStatusStore + BlockedOverlay/TrialExpiredScreen), restaurant-website middleware (maintenance redirect, fail-open). **API enforcement: TenantStatusMiddleware is wired only to `GET /tenant/status`, which it skips; business routes use `wrapProtected`/`wrapWithPermission` without it** (FE-OTH-031, P0). Login does not check `is_disabled`.
4. Paying: `GET /subscription/plans` (public) -> `POST /subscription/initiate` (auth) -> Paymob iframe (mock handler in dev) -> `POST /webhook/paymob` (HMAC checked, rate limited 60/min/IP, idempotent via `errWebhookAlreadyProcessed`) -> on success `updateTenantSubscriptionTx` sets status `active`, plan, `subscription_expires_at = paidAt + plan months` (monthly 1, annual per `periodMonthsFromPlan`), payment row success, history row. Failed payments mark payment failed.
5. Cancel: `POST /subscription/cancel` allowed only for `active`; sets cancellation timestamp, access continues to period end.
6. Expiry: `SubscriptionExpiryWorker` (started in main.go:419, hourly + at boot) flips `active -> expired` when `subscription_expires_at <= NOW()` and writes history. It does **not** move `trial -> expired` (trial handled by date comparison only); `GetExpiredTrials` exists in the repo with no caller.
7. Refund: super-admin `POST /super-admin/payments/{id}/refund` -> `clearSubscriptionAfterRefundTx`; `resolveStatusAfterRefund` returns trial or expired depending on trial end date.
8. Super-admin controls: extend trial (adds days to existing end date, or to now when NULL; if the old end is already past the result may still be in the past; sets `trial_days += days`; status expired->trial), disable/enable (POST SubscriptionHandler; duplicate PUT in SuperAdminHandler), delete tenant (`DELETE ...?confirm=true`), history list. Audit rows in `tenant_subscription_history`.

## 3. Trial length inconsistencies (FE-OTH-004)
| Location | Value |
|---|---|
| Landing hero/pricing/how-it-works/CTA copy EN+AR (`translations.ts:12,63,78,79,102,119`, AR 228...335) | 30 days |
| README feature description (lines 415-418) and subscription section (522, 529) | 30 days ("seeded via migration 112_seed_trial_data.sql" - actual file is `112_seed_existing_tenants_trial.sql`) |
| README registration flow (line 297) | 15 days |
| DB default `trial_days` and COMMENT (migration 110:7,35); seed `created_at + INTERVAL '15 days'` (112:12-13) | 15 days |
| Dashboard `BlockedOverlay.tsx:29`, `TrialExpiredScreen.tsx:44` | "15-day free trial" |
| Super-admin | shows `trial_days` from API (whatever DB holds) |
| Registration code | no value applied (NULL) |
Effective behavior today: for new tenants, unlimited (see 2.1); for pre-existing seeded tenants, 15 days.

## 4. Pricing consistency
Landing 500 EGP/month (hardcoded) = backend default `SUBSCRIPTION_MONTHLY_PRICE_CENTS=50000` and migration 118 seed. Backend also offers annual (5000 EGP default) that the landing does not show. Landing says "Mobile Apps", "Unlimited Restaurants/Staff"; registration provisions 1 restaurant / 5 users and no limit-enforcement code was found (grep `max_users|max_restaurants`: only domain struct, registration insert and an unrelated advanced-recommendations handler).

## 5. Observed gaps (all in issues file)
- FE-OTH-003 no trial dates at registration; no expiry job for trials.
- FE-OTH-031 enforcement absent on API; evaluators disagree on `expired`/`suspended`; cancelled subscriptions only reflect via the worker.
- FE-OTH-004 trial length drift; FE-OTH-010 marketing vs limits.
- Extend-trial arithmetic starts from old end date even if past (subscription_repo_impl.go ~252-257) - candidate P3 (not logged separately).
- Cache: backend tenant-status cache TTL 30 s (`tenant_status.go:16`) - after upgrade/extension users may stay blocked up to 30 s; whether cache is invalidated on payment webhook was not checked.
- `/subscription/plans` is dependable but unused by landing; dashboard consumes it.
- Two payments-table migrations (115 and 143) - needs reconciliation by the DB auditor.

## 6. Suggested tests
Registration sets trial window; middleware blocks per status matrix on real business routes; webhook success->active->worker expiry->blocked; extend-trial on past end date; landing copy vs `trial_days` contract.

## 7. Not verified
Dashboard payment UI flow details, Paymob live behavior, whether the 112 seed ran in any environment, `subscription_expires_at` handling for grace period (`grace_period_days` appears only in the domain struct and migration 110; no logic found).
