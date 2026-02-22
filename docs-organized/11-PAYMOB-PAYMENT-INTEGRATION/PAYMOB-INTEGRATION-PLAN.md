# Paymob Payment Integration — Full Implementation Plan

**Created:** February 21, 2026
**Status:** Ready for Implementation
**Target:** Production-ready payment system with mock Paymob endpoints for testing

---

## Table of Contents

1. [Overview](#overview)
2. [Current System State](#current-system-state)
3. [Paymob API Flow](#paymob-api-flow)
4. [Architecture Design](#architecture-design)
5. [Phase 1 — Database Migrations](#phase-1--database-migrations)
6. [Phase 2 — Backend Payment Service](#phase-2--backend-payment-service)
7. [Phase 3 — Backend API Endpoints](#phase-3--backend-api-endpoints)
8. [Phase 4 — Mock Paymob Server](#phase-4--mock-paymob-server)
9. [Phase 5 — Dashboard Billing Page](#phase-5--dashboard-billing-page)
10. [Phase 6 — Dashboard Payment Flow UI](#phase-6--dashboard-payment-flow-ui)
11. [Phase 7 — BlockedOverlay + Sidebar Upgrade](#phase-7--blockedoverlay--sidebar-upgrade)
12. [Phase 8 — Webhook Handler + Security](#phase-8--webhook-handler--security)
13. [Phase 9 — Testing & Verification](#phase-9--testing--verification)
14. [Environment Variables](#environment-variables)
15. [File Reference](#file-reference)
16. [Task Checklist](#task-checklist)

---

## Overview

### What We're Building

A complete Paymob payment integration that:

1. Shows tenant's current subscription status + remaining trial days in the Dashboard
2. After trial expires, presents a payment page (not just a "contact support" email link)
3. Processes payments via **Paymob Egypt** (card payments via iframe)
4. On successful payment: activates the subscription (500 EGP/month, all features)
5. Records all payment transactions in the database with full audit trail
6. Uses **mock Paymob endpoints** for development/testing until real credentials are obtained
7. Is **100% secure**: HMAC webhook verification, no sensitive data in frontend, idempotent transactions

### Pricing Model

| | Details |
|---|---|
| **Trial** | 30 days free (existing system handles this) |
| **Monthly Subscription** | 500 EGP/month |
| **Plan** | Single plan — all features included |
| **Payment Method** | Card via Paymob iframe |
| **Renewal** | Manual (tenant pays monthly via billing page) |

### Key Principle: Mock-First

All Paymob API calls go through a `PaymobClient` interface. In development, a `MockPaymobClient` returns success responses. In production, a `RealPaymobClient` calls the actual Paymob API. The switch is controlled by the `PAYMOB_MODE` environment variable (`mock` vs `live`).

---

## Current System State

### What Already Exists

| Component | Status | Details |
|-----------|--------|---------|
| Trial system | Done | 30-day trial, auto-expiry, blocking middleware |
| `subscription_status` column | Done | Values: `trial`, `active`, `expired`, `suspended`, `disabled` |
| `tenant_subscription_history` | Done | Full audit trail with JSONB metadata |
| `BlockedOverlay.tsx` | Done | Shows when trial expired — currently only has "Contact Support" mailto link |
| `SidebarSubscriptionBadge` | Done | Shows trial days remaining — currently not clickable |
| `useTenantStatus` hook | Done | Returns `isBlocked`, `blockReason`, `daysRemaining`, `status` |
| `tenantStatusStore` | Done | Zustand store, calls `GET /api/v1/tenant/status` |
| `GET /api/v1/tenant/status` | Done | Returns `can_access`, `block_reason`, `days_remaining`, `subscription_status` |
| Super Admin controls | Done | Extend trial, disable, enable, delete |

### What's Missing

| Component | Needed For |
|-----------|------------|
| `payments` table | Recording transactions |
| Payment columns on `tenants` | Tracking active subscription dates |
| Paymob API integration | Creating payment sessions |
| Webhook endpoint | Receiving payment confirmations from Paymob |
| Billing page in Dashboard | Tenant-facing subscription management |
| Payment button in BlockedOverlay | Converting expired trials to paid |
| Mock Paymob server | Testing without real credentials |

---

## Paymob API Flow

### Standard Accept Integration (5 Steps)

```
Step 1: Authentication
  POST https://accept.paymob.com/api/auth/tokens
  Body: { "api_key": "YOUR_API_KEY" }
  Response: { "token": "auth_token_string" }

Step 2: Order Registration
  POST https://accept.paymob.com/api/ecommerce/orders
  Headers: { "Authorization": "Bearer {auth_token}" }
  Body: {
    "auth_token": "{auth_token}",
    "delivery_needed": "false",
    "amount_cents": "50000",          // 500 EGP = 50000 cents
    "currency": "EGP",
    "merchant_order_id": "tenant_{id}_monthly_{timestamp}",
    "items": [{
      "name": "POS SaaS Monthly Subscription",
      "amount_cents": "50000",
      "quantity": "1"
    }]
  }
  Response: { "id": 123456 }          // paymob_order_id

Step 3: Payment Key Request
  POST https://accept.paymob.com/api/acceptance/payment_keys
  Body: {
    "auth_token": "{auth_token}",
    "amount_cents": "50000",
    "expiration": 3600,
    "order_id": 123456,
    "billing_data": {
      "email": "tenant@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "+201234567890",
      "apartment": "NA", "floor": "NA", "building": "NA",
      "street": "NA", "city": "Cairo", "state": "Cairo",
      "country": "EG", "shipping_method": "NA", "postal_code": "NA"
    },
    "currency": "EGP",
    "integration_id": YOUR_INTEGRATION_ID,
    "lock_order_when_paid": "true"
  }
  Response: { "token": "payment_key_token_string" }

Step 4: Iframe Payment
  URL: https://accept.paymob.com/api/acceptance/iframes/{IFRAME_ID}?payment_token={payment_key}
  - User enters card details in Paymob's hosted iframe
  - Paymob handles 3D Secure if required

Step 5: Webhook Callback (Server-to-Server)
  POST https://your-backend.com/api/v1/webhook/paymob
  Paymob sends transaction result with HMAC signature
  Verify HMAC-SHA512 before processing
```

### HMAC Verification (Critical for Security)

Paymob sends an `hmac` query parameter. To verify:

1. Concatenate these fields from the callback payload **in this exact order**:
   ```
   amount_cents + created_at + currency + error_occured + has_parent_transaction +
   id + integration_id + is_3d_secure + is_auth + is_capture + is_refunded +
   is_standalone_payment + is_voided + order.id + owner + pending +
   source_data.pan + source_data.sub_type + source_data.type + success
   ```

2. Compute `HMAC-SHA512(concatenated_string, PAYMOB_HMAC_SECRET)`

3. Compare with the received `hmac` value

4. Also verify: `success === true`, `amount_cents` matches expected, `order.id` matches

---

## Architecture Design

### System Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                     Dashboard (Next.js :3002)                    │
│                                                                  │
│  ┌─────────────┐   ┌──────────────┐   ┌──────────────────────┐  │
│  │ Billing Page │   │ BlockedOverlay│   │ Sidebar Badge        │  │
│  │ /settings/  │   │ "Pay Now"    │   │ (clickable → billing)│  │
│  │  billing    │   │  button      │   │                      │  │
│  └──────┬──────┘   └──────┬───────┘   └──────────────────────┘  │
│         │                 │                                      │
│         ▼                 ▼                                      │
│  ┌──────────────────────────────┐                                │
│  │ POST /subscription/initiate  │                                │
│  │ → returns iframe_url         │                                │
│  └──────────────┬───────────────┘                                │
│                 ▼                                                 │
│  ┌──────────────────────────────┐                                │
│  │ Paymob Iframe (in modal)     │                                │
│  │ User enters card details     │──── 3D Secure if needed        │
│  └──────────────────────────────┘                                │
│                                                                  │
│         After payment: redirect to /settings/billing?status=...  │
│         Then call: tenantStatusStore.checkStatus() to refresh    │
└──────────────────────────────────────────────────────────────────┘

                              │
                              ▼ (Paymob server-to-server callback)

┌──────────────────────────────────────────────────────────────────┐
│                   Go Backend API (:8080)                         │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐      │
│  │ POST /api/v1/webhook/paymob  (Public, HMAC-verified)   │      │
│  │  1. Verify HMAC-SHA512 signature                       │      │
│  │  2. Verify amount_cents = 50000 & currency = EGP       │      │
│  │  3. Look up payment record by paymob_order_id          │      │
│  │  4. If success: set subscription_status = "active"     │      │
│  │  5. Set subscription_expires_at = NOW() + 30 days      │      │
│  │  6. Record in payments table                           │      │
│  │  7. Record in subscription_history (audit)             │      │
│  │  8. Return 200 OK                                      │      │
│  └────────────────────────────────────────────────────────┘      │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐      │
│  │ POST /api/v1/subscription/initiate  (JWT protected)    │      │
│  │  1. Authenticate with Paymob API                       │      │
│  │  2. Create order on Paymob                             │      │
│  │  3. Get payment key                                    │      │
│  │  4. Save pending payment record                        │      │
│  │  5. Return { iframe_url, payment_id } to frontend      │      │
│  └────────────────────────────────────────────────────────┘      │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐      │
│  │ GET /api/v1/subscription/plans  (Public)               │      │
│  │  Returns: [{ id, name, price, currency, features }]    │      │
│  └────────────────────────────────────────────────────────┘      │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐      │
│  │ GET /api/v1/subscription/payments (JWT protected)      │      │
│  │  Returns: paginated list of tenant's payment records   │      │
│  └────────────────────────────────────────────────────────┘      │
│                                                                  │
│  ┌──────────────────┐   ┌──────────────────┐                    │
│  │ PaymobClient     │   │ MockPaymobClient │                    │
│  │ (interface)       │   │ (for dev/test)   │                    │
│  │                  │   │ Returns success   │                    │
│  │ Authenticate()   │   │ with fake IDs     │                    │
│  │ CreateOrder()    │   │                  │                    │
│  │ GetPaymentKey()  │   │                  │                    │
│  └──────────────────┘   └──────────────────┘                    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

                              │
                              ▼

┌──────────────────────────────────────────────────────────────────┐
│                     PostgreSQL                                   │
│                                                                  │
│  tenants (modified):                                             │
│    + subscription_plan        VARCHAR(20)  DEFAULT 'monthly'     │
│    + subscription_expires_at  TIMESTAMP                          │
│    + last_payment_at          TIMESTAMP                          │
│    + last_payment_amount      INTEGER  (cents)                   │
│    + paymob_customer_id       VARCHAR(100)                       │
│                                                                  │
│  payments (new table):                                           │
│    id, tenant_id, amount_cents, currency, status,                │
│    paymob_order_id, paymob_transaction_id, payment_method,       │
│    card_last_four, card_brand, billing_data, error_message,      │
│    created_at, updated_at                                        │
│                                                                  │
│  tenant_subscription_history (existing):                         │
│    + action = 'PAYMENT_SUCCESS' / 'PAYMENT_FAILED'              │
│    + metadata = { paymob_order_id, paymob_txn_id, amount }      │
└──────────────────────────────────────────────────────────────────┘
```

---

## Phase 1 — Database Migrations

### Task 1.1: Add Payment Fields to Tenants Table

**File:** `backend/migrations/114_add_payment_fields_to_tenants.sql`

```sql
-- Add payment-related fields to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_plan VARCHAR(20) DEFAULT 'monthly';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_at TIMESTAMP;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_amount INTEGER;  -- amount in cents (50000 = 500 EGP)
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS last_payment_currency VARCHAR(3) DEFAULT 'EGP';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS paymob_customer_id VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_tenants_subscription_expires_at ON tenants(subscription_expires_at);
CREATE INDEX IF NOT EXISTS idx_tenants_paymob_customer_id ON tenants(paymob_customer_id);
```

### Task 1.2: Create Payments Table

**File:** `backend/migrations/115_create_payments_table.sql`

```sql
CREATE TABLE IF NOT EXISTS payments (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id         BIGINT REFERENCES users(id) ON DELETE SET NULL,

    -- Amount
    amount_cents    INTEGER NOT NULL,                         -- 50000 = 500 EGP
    currency        VARCHAR(3) NOT NULL DEFAULT 'EGP',

    -- Status
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',   -- pending, success, failed, refunded

    -- Paymob references
    paymob_order_id       BIGINT,
    paymob_transaction_id BIGINT,
    paymob_auth_token     TEXT,                               -- temporary, cleared after use

    -- Payment details (from Paymob callback)
    payment_method        VARCHAR(20),                        -- card, wallet, etc.
    card_last_four        VARCHAR(4),
    card_brand            VARCHAR(20),                        -- visa, mastercard, etc.

    -- Billing
    billing_name          VARCHAR(200),
    billing_email         VARCHAR(200),
    billing_phone         VARCHAR(50),

    -- Plan reference
    plan_id               VARCHAR(20) DEFAULT 'monthly',      -- monthly plan
    plan_period_start     TIMESTAMP,                          -- subscription period start
    plan_period_end       TIMESTAMP,                          -- subscription period end

    -- Error handling
    error_message         TEXT,
    error_code            VARCHAR(50),

    -- Metadata
    merchant_order_id     VARCHAR(100) UNIQUE,                -- our internal reference
    metadata              JSONB DEFAULT '{}',

    -- Timestamps
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at               TIMESTAMP
);

-- Indexes
CREATE INDEX idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_paymob_order_id ON payments(paymob_order_id);
CREATE INDEX idx_payments_merchant_order_id ON payments(merchant_order_id);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);
CREATE INDEX idx_payments_tenant_created ON payments(tenant_id, created_at DESC);

-- Prevent duplicate successful payments for the same Paymob order
CREATE UNIQUE INDEX idx_payments_paymob_order_success
    ON payments(paymob_order_id)
    WHERE status = 'success';
```

### Task 1.3: Add New Subscription History Actions

**File:** `backend/migrations/116_add_payment_history_actions.sql`

```sql
-- Update check constraint to include payment-related actions
-- (Only needed if there's a CHECK constraint on the action column)
-- The existing tenant_subscription_history.action column uses VARCHAR(50)
-- and the code handles any string, so we just need to document the new actions:
-- PAYMENT_INITIATED, PAYMENT_SUCCESS, PAYMENT_FAILED, SUBSCRIPTION_ACTIVATED, SUBSCRIPTION_EXPIRED

-- Add comment for documentation
COMMENT ON TABLE tenant_subscription_history IS
'Audit trail for all subscription changes. Actions: TRIAL_STARTED, TRIAL_EXTENDED, TRIAL_EXPIRED, TENANT_DISABLED, TENANT_ENABLED, DELETION_SCHEDULED, DELETION_CANCELLED, TENANT_DELETED, STATUS_CHANGED, PAYMENT_INITIATED, PAYMENT_SUCCESS, PAYMENT_FAILED, SUBSCRIPTION_ACTIVATED, SUBSCRIPTION_EXPIRED';
```

---

## Phase 2 — Backend Payment Service

### Task 2.1: Create Payment Domain Models

**File:** `backend/internal/domain/payment.go`

Define types:

```go
// Payment statuses
const (
    PaymentStatusPending  = "pending"
    PaymentStatusSuccess  = "success"
    PaymentStatusFailed   = "failed"
    PaymentStatusRefunded = "refunded"
)

// Subscription history actions for payments
const (
    ActionPaymentInitiated      = "PAYMENT_INITIATED"
    ActionPaymentSuccess        = "PAYMENT_SUCCESS"
    ActionPaymentFailed         = "PAYMENT_FAILED"
    ActionSubscriptionActivated = "SUBSCRIPTION_ACTIVATED"
    ActionSubscriptionExpired   = "SUBSCRIPTION_EXPIRED"
)

// Plan info
type SubscriptionPlan struct {
    ID          string   `json:"id"`           // "monthly"
    Name        string   `json:"name"`         // "Pro Plan"
    NameAr      string   `json:"name_ar"`      // "الخطة الاحترافية"
    Price       int      `json:"price"`        // 500 (EGP)
    PriceCents  int      `json:"price_cents"`  // 50000
    Currency    string   `json:"currency"`     // "EGP"
    Period      string   `json:"period"`       // "month"
    Features    []string `json:"features"`
    FeaturesAr  []string `json:"features_ar"`
}

// Payment record
type Payment struct {
    ID                    int64      `json:"id"`
    TenantID              int64      `json:"tenant_id"`
    UserID                *int64     `json:"user_id,omitempty"`
    AmountCents           int        `json:"amount_cents"`
    Currency              string     `json:"currency"`
    Status                string     `json:"status"`
    PaymobOrderID         *int64     `json:"paymob_order_id,omitempty"`
    PaymobTransactionID   *int64     `json:"paymob_transaction_id,omitempty"`
    PaymentMethod         *string    `json:"payment_method,omitempty"`
    CardLastFour          *string    `json:"card_last_four,omitempty"`
    CardBrand             *string    `json:"card_brand,omitempty"`
    BillingName           *string    `json:"billing_name,omitempty"`
    BillingEmail          *string    `json:"billing_email,omitempty"`
    PlanID                string     `json:"plan_id"`
    PlanPeriodStart       *time.Time `json:"plan_period_start,omitempty"`
    PlanPeriodEnd         *time.Time `json:"plan_period_end,omitempty"`
    MerchantOrderID       string     `json:"merchant_order_id"`
    ErrorMessage          *string    `json:"error_message,omitempty"`
    Metadata              map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt             time.Time  `json:"created_at"`
    PaidAt                *time.Time `json:"paid_at,omitempty"`
}

// Paymob API types
type PaymobAuthResponse struct {
    Token string `json:"token"`
}

type PaymobOrderResponse struct {
    ID int64 `json:"id"`
}

type PaymobPaymentKeyResponse struct {
    Token string `json:"token"`
}

type PaymobWebhookPayload struct {
    Obj struct {
        ID                    int64  `json:"id"`
        AmountCents           int    `json:"amount_cents"`
        Success               bool   `json:"success"`
        Currency              string `json:"currency"`
        IsAuth                bool   `json:"is_auth"`
        IsCapture             bool   `json:"is_capture"`
        IsRefunded            bool   `json:"is_refunded"`
        IsVoided              bool   `json:"is_voided"`
        Is3DSecure            bool   `json:"is_3d_secure"`
        IsStandalonePayment   bool   `json:"is_standalone_payment"`
        HasParentTransaction  bool   `json:"has_parent_transaction"`
        ErrorOccured          bool   `json:"error_occured"`
        Pending               bool   `json:"pending"`
        IntegrationID         int64  `json:"integration_id"`
        Owner                 int64  `json:"owner"`
        CreatedAt             string `json:"created_at"`
        Order                 struct {
            ID              int64  `json:"id"`
            MerchantOrderID string `json:"merchant_order_id"`
        } `json:"order"`
        SourceData struct {
            Pan     string `json:"pan"`
            SubType string `json:"sub_type"`
            Type    string `json:"type"`
        } `json:"source_data"`
    } `json:"obj"`
    Type string `json:"type"`
    HMAC string `json:"hmac"`
}

// Initiate payment request/response
type InitiatePaymentRequest struct {
    PlanID string `json:"plan_id"` // "monthly"
}

type InitiatePaymentResponse struct {
    PaymentID       int64  `json:"payment_id"`
    IframeURL       string `json:"iframe_url"`
    PaymobOrderID   int64  `json:"paymob_order_id"`
    MerchantOrderID string `json:"merchant_order_id"`
}
```

### Task 2.2: Create PaymobClient Interface + Mock Implementation

**File:** `backend/internal/service/paymob_client.go`

```go
// PaymobClient is the interface for Paymob API interactions
type PaymobClient interface {
    Authenticate() (string, error)
    CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error)
    GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error)
    GetIframeURL(paymentKey string) string
    VerifyHMAC(payload []byte, receivedHMAC string) bool
}
```

**File:** `backend/internal/service/paymob_client_mock.go`

Mock implementation that:
- `Authenticate()` → returns `"mock_auth_token_" + timestamp`
- `CreateOrder()` → returns random int64 as order ID
- `GetPaymentKey()` → returns `"mock_payment_key_" + timestamp`
- `GetIframeURL()` → returns `/api/v1/mock/paymob/iframe?payment_token={key}`
- `VerifyHMAC()` → returns `true` always in mock mode

**File:** `backend/internal/service/paymob_client_live.go`

Real implementation that calls actual Paymob API endpoints:
- Base URL: `https://accept.paymob.com/api`
- Uses `net/http` client with timeouts
- `VerifyHMAC()` → proper HMAC-SHA512 verification

### Task 2.3: Create Payment Repository

**File:** `backend/internal/repository/payment_repo.go` (interface)

```go
type PaymentRepository interface {
    CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error)
    GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error)
    GetPaymentByMerchantOrderID(ctx context.Context, merchantOrderID string) (*domain.Payment, error)
    GetPaymentByPaymobOrderID(ctx context.Context, paymobOrderID int64) (*domain.Payment, error)
    UpdatePaymentStatus(ctx context.Context, id int64, status string, txnID *int64, cardLastFour *string, cardBrand *string, errorMsg *string, paidAt *time.Time) error
    GetPaymentsByTenantID(ctx context.Context, tenantID int64, limit, offset int) ([]domain.Payment, int, error)
    UpdateTenantSubscription(ctx context.Context, tenantID int64, status string, expiresAt time.Time, lastPaymentAt time.Time, lastPaymentAmount int, currency string) error
}
```

**File:** `backend/internal/repository/payment_repo_impl.go`

PostgreSQL implementation of above interface.

### Task 2.4: Create Payment UseCase

**File:** `backend/internal/usecase/payment_usecase.go`

Business logic:

```go
type PaymentUseCase struct {
    paymentRepo      PaymentRepository
    subscriptionRepo SubscriptionRepository
    paymobClient     PaymobClient
    plans            map[string]domain.SubscriptionPlan
}

// Methods:
// InitiatePayment(ctx, tenantID, userID, planID, billingData) → InitiatePaymentResponse, error
//   1. Validate plan exists
//   2. Generate merchant_order_id: "tenant_{tenantID}_monthly_{unix_timestamp}"
//   3. Create pending payment record in DB
//   4. Call paymobClient.Authenticate()
//   5. Call paymobClient.CreateOrder()
//   6. Call paymobClient.GetPaymentKey()
//   7. Update payment record with paymob_order_id
//   8. Log PAYMENT_INITIATED in subscription_history
//   9. Return { iframe_url, payment_id, paymob_order_id, merchant_order_id }
//
// HandleWebhook(ctx, payload, hmac) → error
//   1. Verify HMAC signature (CRITICAL)
//   2. Find payment by paymob_order_id
//   3. If payment already success → return (idempotent)
//   4. If payload.success:
//      a. Update payment status to "success" + transaction details
//      b. Update tenants: subscription_status="active", subscription_expires_at=NOW()+30 days
//      c. Log PAYMENT_SUCCESS + SUBSCRIPTION_ACTIVATED in history
//   5. If !payload.success:
//      a. Update payment status to "failed" + error details
//      b. Log PAYMENT_FAILED in history
//
// GetPlans() → []SubscriptionPlan
//   Returns the single plan (500 EGP/month)
//
// GetPaymentHistory(ctx, tenantID, limit, offset) → []Payment, total, error
//   Returns paginated payment list for the tenant
```

---

## Phase 3 — Backend API Endpoints

### Task 3.1: Create Payment Handler

**File:** `backend/internal/handler/http/payment_handler.go`

Handlers:

```go
type PaymentHandler struct {
    paymentUseCase *usecase.PaymentUseCase
}

// GET /api/v1/subscription/plans (public)
func (h *PaymentHandler) GetPlans(w http.ResponseWriter, r *http.Request)
// Returns: { "success": true, "data": { "plans": [...] } }

// POST /api/v1/subscription/initiate (JWT protected — even expired tenants allowed)
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request)
// Body: { "plan_id": "monthly" }
// Returns: { "success": true, "data": { "iframe_url": "...", "payment_id": 1, "merchant_order_id": "..." } }
// NOTE: This endpoint must be EXCLUDED from TenantStatusMiddleware
//       so that blocked tenants can still initiate payments!

// POST /api/v1/webhook/paymob (public, HMAC verified)
func (h *PaymentHandler) HandlePaymobWebhook(w http.ResponseWriter, r *http.Request)
// Paymob POSTs here after payment processed
// Verify HMAC before any processing
// Returns: 200 OK (always, even on error — Paymob expects 200)

// GET /api/v1/subscription/payments (JWT protected — even expired tenants allowed)
func (h *PaymentHandler) GetPaymentHistory(w http.ResponseWriter, r *http.Request)
// Query: ?page=1&page_size=10
// Returns: { "success": true, "data": { "payments": [...], "total": 5, "page": 1 } }

// GET /api/v1/subscription/current (JWT protected — even expired tenants allowed)
func (h *PaymentHandler) GetCurrentSubscription(w http.ResponseWriter, r *http.Request)
// Returns: extended tenant status with payment info
// { subscription_status, trial_end_date, subscription_expires_at, last_payment_at, plan, days_remaining }
```

### Task 3.2: Register Routes in main.go

**File:** `backend/cmd/api/main.go`

Add these routes:

```go
// Public endpoints (no auth)
mux.HandleFunc("GET /api/v1/subscription/plans", paymentHandler.GetPlans)
mux.HandleFunc("POST /api/v1/webhook/paymob", paymentHandler.HandlePaymobWebhook)

// Payment endpoints (auth required, but SKIP tenant status check)
// These use a special wrapper that has auth + tenant context but NO status middleware
// so expired/blocked tenants can still pay
wrapPayment := func(handler http.HandlerFunc) http.HandlerFunc {
    return authMiddleware(tenantContextMiddleware(handler))
}
mux.HandleFunc("POST /api/v1/subscription/initiate", wrapPayment(paymentHandler.InitiatePayment))
mux.HandleFunc("GET /api/v1/subscription/payments", wrapPayment(paymentHandler.GetPaymentHistory))
mux.HandleFunc("GET /api/v1/subscription/current", wrapPayment(paymentHandler.GetCurrentSubscription))
```

**Critical:** The payment endpoints must NOT use `TenantStatusMiddleware` — otherwise expired tenants can't pay!

---

## Phase 4 — Mock Paymob Server

### Task 4.1: Create Mock Paymob Iframe Page

**File:** `backend/internal/handler/http/mock_paymob_handler.go`

A handler that serves a fake payment iframe for development:

```go
// GET /api/v1/mock/paymob/iframe?payment_token={key}
// Serves an HTML page that looks like a payment form
// Has card fields (prefilled with test card)
// "Pay" button → POST /api/v1/mock/paymob/process
// "Fail" button → POST /api/v1/mock/paymob/process?fail=true

// POST /api/v1/mock/paymob/process
// Simulates payment processing
// Calls the webhook endpoint: POST /api/v1/webhook/paymob with a mock payload
// Redirects to: dashboard_url/settings/billing?status=success or ?status=failed
```

### Task 4.2: Register Mock Routes (Only in Mock Mode)

**File:** `backend/cmd/api/main.go`

```go
if os.Getenv("PAYMOB_MODE") == "mock" || os.Getenv("PAYMOB_MODE") == "" {
    mux.HandleFunc("GET /api/v1/mock/paymob/iframe", mockPaymobHandler.ServeIframe)
    mux.HandleFunc("POST /api/v1/mock/paymob/process", mockPaymobHandler.ProcessPayment)
    log.Println("[PaymentSystem] Running in MOCK mode — mock Paymob endpoints enabled")
}
```

---

## Phase 5 — Dashboard Billing Page

### Task 5.1: Create Billing Settings Page

**File:** `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/billing/page.tsx`

This is the main billing/subscription management page. Sections:

**Section 1: Current Plan Status**
- Shows current subscription status (Trial / Active / Expired)
- If trial: shows days remaining with a countdown badge, trial end date
- If active: shows "Pro Plan — Active", subscription expires date, last payment date
- If expired: shows "Expired" with red warning, "Renew Now" button

**Section 2: Plan Details**
- Card showing: "Pro Plan — 500 EGP/month"
- Feature list (10 items)
- "Subscribe Now" / "Renew" button (shown for trial/expired)
- "Already Active" badge (shown for active)

**Section 3: Payment History**
- Table with columns: Date, Amount, Status, Payment Method, Invoice
- Paginated (from `GET /api/v1/subscription/payments`)
- Empty state: "No payments yet"

### Task 5.2: Add Billing to Settings Sidebar

**File:** `frontend/apps/dashboard/src/components/settings/SettingsSidebar.tsx`

Add new nav item:
```
{ label: t.settings.billing, href: `/${locale}/dashboard/settings/billing`, icon: CreditCard }
```

### Task 5.3: Add Billing Translations

**File:** `frontend/apps/dashboard/src/i18n/messages/en.json`

Add under `settings`:
```json
"billing": "Billing & Subscription",
"billingPage": {
  "title": "Billing & Subscription",
  "currentPlan": "Current Plan",
  "planName": "Pro Plan",
  "priceMonthly": "500 EGP/month",
  "trialActive": "Free Trial Active",
  "trialDaysRemaining": "{days} days remaining",
  "trialEndsOn": "Trial ends on {date}",
  "subscriptionActive": "Active Subscription",
  "subscriptionExpiresOn": "Expires on {date}",
  "lastPayment": "Last payment: {date}",
  "expired": "Subscription Expired",
  "expiredMessage": "Your trial has ended. Subscribe to continue using all features.",
  "subscribeNow": "Subscribe Now — 500 EGP/month",
  "renewSubscription": "Renew Subscription",
  "processing": "Processing payment...",
  "paymentHistory": "Payment History",
  "noPayments": "No payments yet",
  "date": "Date",
  "amount": "Amount",
  "status": "Status",
  "method": "Method",
  "paymentSuccess": "Payment successful! Your subscription is now active.",
  "paymentFailed": "Payment failed. Please try again.",
  "allFeaturesIncluded": "All Features Included",
  "features": {
    "pos": "Full POS System",
    "kitchen": "Kitchen Display System",
    "online": "Online Ordering Website",
    "analytics": "Analytics & Reports",
    "multiRestaurant": "Multi-Restaurant Management",
    "mobile": "Mobile Apps",
    "unlimitedStaff": "Unlimited Staff",
    "support": "Priority Support",
    "updates": "Free Updates",
    "drivers": "Driver Management"
  }
}
```

**File:** `frontend/apps/dashboard/src/i18n/messages/ar.json`

Full Arabic parity for all billing keys above.

---

## Phase 6 — Dashboard Payment Flow UI

### Task 6.1: Create usePayment Hook

**File:** `frontend/apps/dashboard/src/hooks/usePayment.ts`

```typescript
// Hook that manages the Paymob payment flow:
//
// const { initiatePayment, isLoading, error, paymentStatus } = usePayment()
//
// initiatePayment(planId: string) → Promise<void>
//   1. POST /api/v1/subscription/initiate with { plan_id: planId }
//   2. Opens PaymentModal with iframe_url
//   3. Listens for iframe messages (payment completion)
//
// Also handles:
//   - URL query params ?status=success|failed (redirect callback)
//   - Calls tenantStatusStore.checkStatus() after success
//   - Toast notifications for success/failure
```

### Task 6.2: Create PaymentModal Component

**File:** `frontend/apps/dashboard/src/components/billing/PaymentModal.tsx`

A modal that:
- Renders full-screen overlay with Paymob iframe
- Iframe src = the `iframe_url` from initiate response
- Shows loading spinner while iframe loads
- Has a "Cancel" button to close
- Listens for `postMessage` events from iframe (payment completion)
- On completion: closes modal, shows toast, refreshes status

### Task 6.3: Create PlanCard Component

**File:** `frontend/apps/dashboard/src/components/billing/PlanCard.tsx`

Reusable plan display card showing:
- Plan name, price (500 EGP/month), features checklist
- CTA button ("Subscribe" / "Renew" / "Current Plan")

### Task 6.4: Create PaymentHistoryTable Component

**File:** `frontend/apps/dashboard/src/components/billing/PaymentHistoryTable.tsx`

Table with:
- Columns: Date, Amount, Status (badge), Payment Method, Card
- Pagination
- Empty state
- Status badges: green "Success", red "Failed", yellow "Pending"

---

## Phase 7 — BlockedOverlay + Sidebar Upgrade

### Task 7.1: Update BlockedOverlay with Pay Now Button

**File:** `frontend/apps/dashboard/src/components/blocked/BlockedOverlay.tsx`

Modify the `trial_expired` case:

**Before:** "Contact Support" mailto link
**After:**
- Primary CTA: "Subscribe Now — 500 EGP/month" button (orange gradient)
  - Clicking it either:
    - Directly opens the PaymentModal (with initiate call), OR
    - Navigates to `/settings/billing` where the user can review the plan before paying
  - Recommended: Navigate to billing page (more info, less jarring)
- Secondary: "Contact Support" link (keep, but smaller)
- Add trial info: "Your 30-day free trial has ended"
- Add price callout: "Continue with all features for just 500 EGP/month"

### Task 7.2: Update SidebarSubscriptionBadge

**File:** `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` (inside SidebarSubscriptionBadge)

Changes:
- Make the badge **clickable** → navigates to `/settings/billing`
- On `status === 'expired'`: show "Subscribe" CTA styling (orange/red)
- On `status === 'trial'` with < 3 days: show "Subscribe Now" small text below the badge

### Task 7.3: Update Tenant Status Middleware Skip Path

**File:** `frontend/apps/dashboard/src/components/providers/DashboardAuthProvider.tsx`

Ensure the `/settings/billing` page is accessible even when the tenant is blocked:
- The DashboardAuthProvider or layout should allow the billing page to render without BlockedOverlay
- Option A: Add `/settings/billing` to a "skip blocked" paths list
- Option B: The billing page itself handles the blocked state internally

### Task 7.4: Update Tenant Status Response

**File:** `backend/internal/handler/http/subscription_handler.go`

Modify `GetTenantStatus` response to include payment-related fields:
```json
{
  "subscription_status": "expired",
  "trial_end_date": "2026-02-16T14:45:00Z",
  "subscription_expires_at": null,
  "last_payment_at": null,
  "days_remaining": -5,
  "can_access": false,
  "block_reason": "trial_expired",
  "plan": {
    "id": "monthly",
    "price": 500,
    "currency": "EGP"
  }
}
```

---

## Phase 8 — Webhook Handler + Security

### Task 8.1: Implement HMAC Verification

**File:** `backend/internal/service/paymob_client_live.go`

```go
func (c *RealPaymobClient) VerifyHMAC(payload []byte, receivedHMAC string) bool {
    // 1. Parse the payload JSON to extract required fields
    // 2. Concatenate fields in EXACT order (see Paymob docs):
    //    amount_cents, created_at, currency, error_occured, has_parent_transaction,
    //    id, integration_id, is_3d_secure, is_auth, is_capture, is_refunded,
    //    is_standalone_payment, is_voided, order.id, owner, pending,
    //    source_data.pan, source_data.sub_type, source_data.type, success
    // 3. Compute HMAC-SHA512 with PAYMOB_HMAC_SECRET
    // 4. Compare hex digest with receivedHMAC
    // 5. Return true only if they match
}
```

### Task 8.2: Webhook Handler with Idempotency

**File:** `backend/internal/handler/http/payment_handler.go` (HandlePaymobWebhook)

Security measures:
1. **Always return 200** — Paymob expects 200 even on errors (prevents retries that could cause double-processing)
2. **HMAC verification first** — reject before any DB operations
3. **Idempotency** — check if payment already marked success before processing
4. **Amount verification** — confirm `amount_cents` matches expected plan price
5. **Transaction logging** — record everything in metadata JSONB
6. **No sensitive data in logs** — never log card numbers, only last 4

### Task 8.3: Rate Limiting on Webhook

**File:** `backend/cmd/api/main.go`

Apply a rate limiter to `/api/v1/webhook/paymob`:
- Max 60 requests/minute from any single IP
- Prevents abuse of the webhook endpoint

### Task 8.4: CORS Configuration

**File:** `backend/cmd/api/main.go`

Ensure the Paymob iframe domain is allowed in CORS:
- `https://accept.paymob.com` must be in allowed origins for iframe communication
- The webhook endpoint should accept requests from Paymob's servers (no CORS needed for server-to-server)

---

## Phase 9 — Testing & Verification

### Task 9.1: Mock Flow End-to-End Test

1. Start backend with `PAYMOB_MODE=mock`
2. Login to Dashboard with an expired trial tenant
3. See BlockedOverlay with "Subscribe Now" button
4. Click "Subscribe Now" → navigates to `/settings/billing`
5. On billing page, click "Subscribe Now — 500 EGP/month"
6. Mock iframe opens → click "Pay" button
7. Mock webhook fires → payment processed
8. Redirect to billing page with `?status=success`
9. Tenant status refreshes → no longer blocked
10. Dashboard fully accessible
11. Payment appears in billing page history table
12. Subscription history shows PAYMENT_SUCCESS and SUBSCRIPTION_ACTIVATED

### Task 9.2: API Endpoint Tests

```bash
# Get plans
curl http://localhost:8080/api/v1/subscription/plans
# Expected: { "success": true, "data": { "plans": [{ "id": "monthly", "price": 500 }] } }

# Initiate payment (with JWT)
curl -X POST http://localhost:8080/api/v1/subscription/initiate \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{"plan_id": "monthly"}'
# Expected: { "success": true, "data": { "iframe_url": "...", "payment_id": 1 } }

# Get payment history
curl http://localhost:8080/api/v1/subscription/payments \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: 1"
# Expected: { "success": true, "data": { "payments": [...], "total": 1 } }

# Get current subscription
curl http://localhost:8080/api/v1/subscription/current \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: 1"
# Expected: { "success": true, "data": { "subscription_status": "active", ... } }
```

### Task 9.3: Security Tests

1. **HMAC tampered** — modify webhook payload, verify rejection
2. **Amount mismatch** — send webhook with wrong amount_cents, verify rejection
3. **Duplicate webhook** — send same webhook twice, verify idempotent (only 1 payment record)
4. **Missing auth on protected endpoints** — verify 401
5. **SQL injection attempts** — verify parameterized queries
6. **XSS in billing page** — verify escaped output

### Task 9.4: Dashboard Build Verification

```bash
cd frontend
pnpm build
# Verify: no TypeScript errors, all pages compile
```

---

## Environment Variables

### New Backend Variables

Add to `backend/.env`:

```bash
# Paymob Configuration
PAYMOB_MODE=mock                          # "mock" or "live"
PAYMOB_API_KEY=                           # Paymob API key (from dashboard)
PAYMOB_INTEGRATION_ID=                    # Card integration ID
PAYMOB_IFRAME_ID=                         # Iframe ID (from Paymob dashboard)
PAYMOB_HMAC_SECRET=                       # HMAC secret (from Paymob dashboard)
PAYMOB_BASE_URL=https://accept.paymob.com/api  # Paymob API base URL

# Subscription Plan
SUBSCRIPTION_MONTHLY_PRICE_CENTS=50000    # 500 EGP = 50000 cents
SUBSCRIPTION_CURRENCY=EGP
```

### Frontend Variables

Add to `frontend/apps/dashboard/.env.local`:

```bash
NEXT_PUBLIC_PAYMOB_MODE=mock              # Controls whether to show mock indicators
```

---

## File Reference

### All Files to Create/Modify

```
backend/
├── migrations/
│   ├── 114_add_payment_fields_to_tenants.sql     [CREATE] Phase 1
│   ├── 115_create_payments_table.sql              [CREATE] Phase 1
│   └── 116_add_payment_history_actions.sql        [CREATE] Phase 1
│
├── internal/
│   ├── domain/
│   │   └── payment.go                             [CREATE] Phase 2
│   │
│   ├── service/
│   │   ├── paymob_client.go                       [CREATE] Phase 2 (interface)
│   │   ├── paymob_client_mock.go                  [CREATE] Phase 2 (mock)
│   │   └── paymob_client_live.go                  [CREATE] Phase 2 (real)
│   │
│   ├── repository/
│   │   ├── payment_repo.go                        [CREATE] Phase 2 (interface)
│   │   └── payment_repo_impl.go                   [CREATE] Phase 2 (implementation)
│   │
│   ├── usecase/
│   │   └── payment_usecase.go                     [CREATE] Phase 2
│   │
│   └── handler/http/
│       ├── payment_handler.go                     [CREATE] Phase 3
│       └── mock_paymob_handler.go                 [CREATE] Phase 4
│
├── cmd/api/
│   └── main.go                                    [MODIFY] Phase 3 + 4 (register routes)
│
└── .env.example                                   [MODIFY] Add Paymob vars

frontend/apps/dashboard/
├── src/
│   ├── app/[locale]/dashboard/settings/
│   │   └── billing/
│   │       └── page.tsx                           [CREATE] Phase 5
│   │
│   ├── components/
│   │   ├── billing/
│   │   │   ├── PaymentModal.tsx                   [CREATE] Phase 6
│   │   │   ├── PlanCard.tsx                       [CREATE] Phase 6
│   │   │   └── PaymentHistoryTable.tsx            [CREATE] Phase 6
│   │   │
│   │   ├── blocked/
│   │   │   └── BlockedOverlay.tsx                 [MODIFY] Phase 7
│   │   │
│   │   ├── layout/
│   │   │   └── Sidebar.tsx                        [MODIFY] Phase 7
│   │   │
│   │   ├── settings/
│   │   │   └── SettingsSidebar.tsx                [MODIFY] Phase 5
│   │   │
│   │   └── providers/
│   │       └── DashboardAuthProvider.tsx          [MODIFY] Phase 7
│   │
│   ├── hooks/
│   │   └── usePayment.ts                         [CREATE] Phase 6
│   │
│   └── i18n/messages/
│       ├── en.json                                [MODIFY] Phase 5 (add billing keys)
│       └── ar.json                                [MODIFY] Phase 5 (add billing keys)
│
└── .env.local                                     [MODIFY] Add NEXT_PUBLIC_PAYMOB_MODE
```

**Total new files:** 15
**Total modified files:** 9
**Total:** 24 files

---

## Task Checklist

### Phase 1 — Database Migrations
- [ ] 1.1 Create `114_add_payment_fields_to_tenants.sql`
- [ ] 1.2 Create `115_create_payments_table.sql`
- [ ] 1.3 Create `116_add_payment_history_actions.sql`
- [ ] 1.4 Run migrations and verify schema

### Phase 2 — Backend Payment Service
- [ ] 2.1 Create `domain/payment.go` (types, constants)
- [ ] 2.2 Create `service/paymob_client.go` (interface)
- [ ] 2.3 Create `service/paymob_client_mock.go`
- [ ] 2.4 Create `service/paymob_client_live.go`
- [ ] 2.5 Create `repository/payment_repo.go` (interface)
- [ ] 2.6 Create `repository/payment_repo_impl.go`
- [ ] 2.7 Create `usecase/payment_usecase.go`
- [ ] 2.8 Verify `go build ./...` compiles

### Phase 3 — Backend API Endpoints
- [ ] 3.1 Create `handler/http/payment_handler.go` (5 handlers)
- [ ] 3.2 Register routes in `main.go` (with correct middleware — NO status check on payment routes)
- [ ] 3.3 Test `GET /subscription/plans` returns plan data
- [ ] 3.4 Test `POST /subscription/initiate` returns iframe URL
- [ ] 3.5 Test `GET /subscription/payments` returns empty list
- [ ] 3.6 Test `GET /subscription/current` returns status

### Phase 4 — Mock Paymob Server
- [ ] 4.1 Create `handler/http/mock_paymob_handler.go`
- [ ] 4.2 Register mock routes (only when `PAYMOB_MODE=mock`)
- [ ] 4.3 Test mock iframe serves HTML form
- [ ] 4.4 Test mock payment process triggers webhook callback
- [ ] 4.5 Test mock redirect to billing page with status

### Phase 5 — Dashboard Billing Page
- [ ] 5.1 Create `settings/billing/page.tsx`
- [ ] 5.2 Add "Billing" to SettingsSidebar navigation
- [ ] 5.3 Add billing translations to `en.json`
- [ ] 5.4 Add billing translations to `ar.json`
- [ ] 5.5 Verify billing page renders current plan status
- [ ] 5.6 Verify billing page shows payment history (empty initially)

### Phase 6 — Dashboard Payment Flow UI
- [ ] 6.1 Create `hooks/usePayment.ts`
- [ ] 6.2 Create `components/billing/PaymentModal.tsx`
- [ ] 6.3 Create `components/billing/PlanCard.tsx`
- [ ] 6.4 Create `components/billing/PaymentHistoryTable.tsx`
- [ ] 6.5 Wire billing page to usePayment hook
- [ ] 6.6 Test: click "Subscribe Now" opens modal with iframe
- [ ] 6.7 Test: after mock payment, status refreshes

### Phase 7 — BlockedOverlay + Sidebar Upgrade
- [ ] 7.1 Update `BlockedOverlay.tsx` — add "Subscribe Now" button for `trial_expired`
- [ ] 7.2 Update `SidebarSubscriptionBadge` — make clickable, link to billing
- [ ] 7.3 Ensure billing page is accessible when tenant is blocked
- [ ] 7.4 Update backend `GetTenantStatus` to include plan/payment info
- [ ] 7.5 Test: expired tenant sees "Subscribe Now" in overlay
- [ ] 7.6 Test: clicking badge navigates to billing

### Phase 8 — Webhook Handler + Security
- [ ] 8.1 Implement HMAC-SHA512 verification in live client
- [ ] 8.2 Add idempotency check in webhook handler
- [ ] 8.3 Add amount verification in webhook handler
- [ ] 8.4 Add rate limiting on webhook endpoint
- [ ] 8.5 Update CORS config for Paymob iframe domain
- [ ] 8.6 Security audit: no SQL injection, no XSS, no sensitive data in logs

### Phase 9 — Testing & Verification
- [ ] 9.1 Full mock flow end-to-end test (12 steps above)
- [ ] 9.2 API endpoint tests (4 curl tests above)
- [ ] 9.3 Security tests (6 tests above)
- [ ] 9.4 `go build ./...` succeeds
- [ ] 9.5 `pnpm build` succeeds (all frontend apps)
- [ ] 9.6 Arabic RTL billing page renders correctly
- [ ] 9.7 Mobile responsive billing page
- [ ] 9.8 Update documentation (`99-REMAINING-WORK.md`, this file)

---

## Implementation Order (Recommended)

```
Phase 1 (Migrations)          ▎ ~30 min
    ↓
Phase 2 (Backend Service)     ▎ ~2 hrs
    ↓
Phase 3 (Backend API)         ▎ ~1 hr
    ↓
Phase 4 (Mock Server)         ▎ ~1 hr
    ↓
Phase 8 (Webhook Security)    ▎ ~1 hr  ← do security early
    ↓
Phase 5 (Billing Page)        ▎ ~2 hrs
    ↓
Phase 6 (Payment Flow UI)     ▎ ~2 hrs
    ↓
Phase 7 (Overlay + Sidebar)   ▎ ~1 hr
    ↓
Phase 9 (Testing)             ▎ ~1 hr
```

---

## Switching from Mock to Live

When you receive real Paymob credentials:

1. Update `backend/.env`:
   ```bash
   PAYMOB_MODE=live
   PAYMOB_API_KEY=your_real_api_key
   PAYMOB_INTEGRATION_ID=your_integration_id
   PAYMOB_IFRAME_ID=your_iframe_id
   PAYMOB_HMAC_SECRET=your_hmac_secret
   ```

2. In Paymob Dashboard, set callback URLs:
   - **Transaction processed callback:** `https://your-api-domain.com/api/v1/webhook/paymob`
   - **Transaction response callback:** `https://your-dashboard-domain.com/{locale}/dashboard/settings/billing`

3. Restart backend — it will automatically use `RealPaymobClient` instead of `MockPaymobClient`

4. Test with Paymob's test cards before going fully live

---

## Security Checklist

- [ ] HMAC-SHA512 verification on every webhook call
- [ ] Amount verification (must match plan price exactly)
- [ ] Idempotent webhook processing (duplicate safe)
- [ ] No card data stored server-side (only last 4 digits + brand from Paymob callback)
- [ ] `PAYMOB_API_KEY` and `PAYMOB_HMAC_SECRET` in environment variables only (never in code)
- [ ] Payment endpoints skip status middleware but require JWT auth
- [ ] Webhook endpoint has rate limiting
- [ ] All payment SQL queries use parameterized statements (no string concatenation)
- [ ] No sensitive data in console logs (card numbers, API keys)
- [ ] CORS allows Paymob iframe domain
- [ ] Frontend never sees or handles raw card data
- [ ] Mock mode clearly indicated in UI (dev only)

---

## Phase 10 — Unit & Integration Tests

### Task 10.1: PaymobClient Unit Tests

**File:** `backend/internal/service/paymob_client_test.go`

Test cases:
- `MockPaymobClient.Authenticate()` returns valid token format
- `MockPaymobClient.CreateOrder()` returns positive order ID
- `MockPaymobClient.GetPaymentKey()` returns non-empty key
- `MockPaymobClient.VerifyHMAC()` always returns true
- `RealPaymobClient.VerifyHMAC()` with known test vectors (pre-computed HMAC)
- `RealPaymobClient.VerifyHMAC()` rejects tampered data

### Task 10.2: Payment UseCase Unit Tests

**File:** `backend/internal/usecase/payment_usecase_test.go`

Test cases:
- `InitiatePayment()` creates pending payment record and returns iframe URL
- `HandleWebhook()` with success payload updates status to `active`
- `HandleWebhook()` with failed payload records failure
- `HandleWebhook()` duplicate call is idempotent
- `HandleWebhook()` wrong amount is rejected
- `GetPlans()` returns correct plan data
- `GetPaymentHistory()` returns paginated results

### Task 10.3: Payment Handler Unit Tests

**File:** `backend/internal/handler/http/payment_handler_test.go`

Test cases:
- `GetPlans` returns 200 with JSON plan data
- `InitiatePayment` returns 401 without JWT
- `InitiatePayment` returns 200 with valid JWT
- `HandlePaymobWebhook` returns 200 on valid payload
- `GetPaymentHistory` returns 200 with pagination

### Task 10.4: Postman/Newman Integration Tests

**File:** `backend/tests/paymob_api_tests.json`

Postman collection covering all 5 payment endpoints with:
- Environment variables for base URL, JWT token
- Pre-request scripts for authentication
- Test scripts validating response schemas and status codes
- Run with: `newman run backend/tests/paymob_api_tests.json -e backend/tests/env.json`

---

## Phase 11 — Monitoring, Logging & Webhook Reliability

### Task 11.1: Structured Webhook Logging

Add structured logging to all payment operations with `[PaymentWebhook]` prefix:
- Log every webhook: timestamp, paymob_order_id, success/fail, processing_time_ms
- Log HMAC verification result (pass/fail — never log the secret)
- Log payment status transitions (pending → success, pending → failed)
- All logs go through existing debug wrapper

### Task 11.2: Webhook Retry Queue

**File:** `backend/migrations/117_create_webhook_retry_queue.sql`

```sql
CREATE TABLE webhook_retry_queue (
    id              BIGSERIAL PRIMARY KEY,
    endpoint        VARCHAR(500) NOT NULL,
    payload         JSONB NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',   -- pending, processing, success, failed
    attempts        INTEGER DEFAULT 0,
    max_attempts    INTEGER DEFAULT 5,
    next_retry_at   TIMESTAMP NOT NULL,
    last_error      TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webhook_retry_status ON webhook_retry_queue(status, next_retry_at);
```

Exponential backoff: 1min → 5min → 30min → 2hr → 12hr

### Task 11.3: Webhook Retry Worker

**File:** `backend/internal/service/webhook_retry_worker.go`

Background goroutine that:
- Polls `webhook_retry_queue` every 30 seconds
- Processes due retries (where `next_retry_at <= NOW()` and `status = 'pending'`)
- On success: marks as `success`
- On failure: increments `attempts`, sets next backoff, records error
- After max retries: marks as `failed`
- Starts as part of `main.go` initialization

### Task 11.4: Error Tracking (Sentry — Optional)

**File:** `backend/internal/service/error_tracker.go`

```go
type ErrorTracker interface {
    CaptureError(err error, context map[string]interface{})
    CaptureMessage(msg string, level string)
}

// SentryTracker — uses sentry-go SDK (when SENTRY_DSN is set)
// NoopTracker — does nothing (when no SENTRY_DSN)
```

Track: webhook failures, payment initiation errors, HMAC mismatches, DB errors.

### Task 11.5: Payment Metrics Endpoint (Super Admin)

**Endpoint:** `GET /api/v1/super-admin/payment-metrics`

Returns:
```json
{
  "total_payments": 150,
  "successful_payments": 142,
  "failed_payments": 8,
  "success_rate": 94.6,
  "total_revenue_cents": 7100000,
  "total_revenue_formatted": "71,000 EGP",
  "avg_processing_time_ms": 2350,
  "period": "month"
}
```

Supports query params: `?period=today|week|month|all`

### Task 11.6: Webhook Health Dashboard (Super Admin)

Add to super admin tenant management UI:
- Last 50 webhook events table (time, order_id, status, processing_time)
- Success/fail ratio chart
- Retry queue size indicator
- Failed retries alert count

---

## Phase 12 — Multi-Plan Support (Future-Ready)

### Task 12.1: Subscription Plans Table

**File:** `backend/migrations/118_create_subscription_plans_table.sql`

```sql
CREATE TABLE subscription_plans (
    id              VARCHAR(20) PRIMARY KEY,          -- 'monthly', 'annual'
    name_en         VARCHAR(100) NOT NULL,
    name_ar         VARCHAR(100) NOT NULL,
    price_cents     INTEGER NOT NULL,
    currency        VARCHAR(3) DEFAULT 'EGP',
    period_months   INTEGER NOT NULL,                 -- 1 for monthly, 12 for annual
    features        JSONB NOT NULL DEFAULT '[]',
    features_ar     JSONB NOT NULL DEFAULT '[]',
    is_active       BOOLEAN DEFAULT true,
    display_order   INTEGER DEFAULT 0,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed plans
INSERT INTO subscription_plans (id, name_en, name_ar, price_cents, currency, period_months, display_order, features, features_ar)
VALUES
  ('monthly', 'Pro Plan', 'الخطة الاحترافية', 50000, 'EGP', 1, 1,
   '["Full POS System","Kitchen Display System","Online Ordering Website","Analytics & Reports","Multi-Restaurant Management","Mobile Apps","Unlimited Staff","Priority Support","Free Updates","Driver Management"]',
   '["نظام نقاط البيع الكامل","نظام عرض المطبخ","موقع الطلبات أونلاين","التحليلات والتقارير","إدارة مطاعم متعددة","تطبيقات الموبايل","عدد غير محدود من الموظفين","دعم فني مميز","تحديثات مجانية","إدارة السائقين"]'),
  ('annual', 'Pro Plan (Annual)', 'الخطة الاحترافية (سنوية)', 500000, 'EGP', 12, 2,
   '["Everything in Monthly","Save 17% (2 months free)","Annual billing convenience"]',
   '["كل مميزات الخطة الشهرية","وفر 17% (شهرين مجاناً)","فوترة سنوية مريحة"]');
```

### Task 12.2: Plan Repository

**File:** `backend/internal/repository/plan_repo.go`

Interface: `GetAllActive()`, `GetByCode()`, `Create()`, `Update()`, `Deactivate()`

### Task 12.3: Update UseCase to Load Plans from DB

Replace hardcoded plan in `payment_usecase.go` with DB-driven plans.

### Task 12.4: Update Billing Page — Plan Selector

If multiple plans active: show Monthly / Annual toggle with savings badge.

### Task 12.5: Update Subscription Response

Include `available_plans[]` and `current_plan` details from DB in `GetCurrentSubscription` response.

---

## Phase 13 — Refunds & Cancellations

### Task 13.1: Paymob Refund API Integration

Add to `PaymobClient` interface:
```go
RefundTransaction(authToken string, transactionID int64, amountCents int) error
```
- Mock: always returns success
- Live: `POST https://accept.paymob.com/api/acceptance/void_refund/refund`

### Task 13.2: Refund Endpoint

**Endpoint:** `POST /api/v1/super-admin/payments/{id}/refund` (super admin only)

Flow:
1. Validate payment exists and status is `success`
2. Call Paymob refund API
3. Update payment status to `refunded`
4. Update tenant subscription status back to `expired`
5. Log PAYMENT_REFUNDED in subscription history

### Task 13.3: Auto-Expiry Cron Job

**File:** `backend/internal/service/subscription_expiry_worker.go`

Background goroutine runs every hour:
- Find tenants where `subscription_status = 'active'` AND `subscription_expires_at < NOW()`
- Set `subscription_status = 'expired'`
- Log SUBSCRIPTION_EXPIRED in history
- (Future: send email notification)

### Task 13.4: Cancellation Endpoint

**Endpoint:** `POST /api/v1/subscription/cancel` (JWT protected)

- Tenant can cancel their own subscription
- Subscription remains active until `subscription_expires_at` (no immediate cutoff)
- Sets `cancelled_at` flag on tenant
- Logs SUBSCRIPTION_CANCELLED in history

### Task 13.5: Update Billing Page — Cancellation UI

- "Cancel Subscription" button for active subscriptions
- Confirmation modal: "Your subscription will remain active until {date}"
- After cancellation: "Subscription will expire on {date}" message

### Task 13.6: Super Admin Refund UI

In tenant detail page:
- "Refund Last Payment" button
- Confirmation modal showing refund amount
- Updates subscription history table

---

## Phase 14 — Performance & Caching

### Task 14.1: Redis Caching (Optional)

**File:** `backend/internal/service/cache.go`

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// RedisCache — uses go-redis (when REDIS_URL is set)
// InMemoryCache — map-based with TTL (fallback)
```

Cache tenant status for 30 seconds (existing behavior), plans for 5 minutes.

### Task 14.2: Cache Subscription Plans

Plans rarely change — cache for 5 minutes, invalidate on plan update.

### Task 14.3: Cursor-Based Pagination for Payment History

Add optional cursor-based pagination (more efficient for large datasets):
- `GET /api/v1/subscription/payments?cursor={last_id}&limit=10`
- Keep offset-based for backward compatibility

### Task 14.4: Connection Pool Configuration

**Env vars:**
```bash
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=5m
```

### Task 14.5: Response Compression

Add gzip middleware for API responses > 1KB to reduce payload sizes.

---

## Updated Implementation Order (All 14 Phases)

```
MVP (Phases 1-9):
  Phase 1 (Migrations)          ▎ ~30 min
  Phase 2 (Backend Service)     ▎ ~2 hrs
  Phase 3 (Backend API)         ▎ ~1 hr
  Phase 4 (Mock Server)         ▎ ~1 hr
  Phase 8 (Webhook Security)    ▎ ~1 hr  ← do security early
  Phase 5 (Billing Page)        ▎ ~2 hrs
  Phase 6 (Payment Flow UI)     ▎ ~2 hrs
  Phase 7 (Overlay + Sidebar)   ▎ ~1 hr
  Phase 9 (Testing)             ▎ ~1 hr

Pre-Production (Phase 10-11):
  Phase 10 (Unit/Integration Tests)    ▎ ~3 hrs
  Phase 11 (Monitoring + Retry Queue)  ▎ ~3 hrs

Post-Launch Enhancements (Phases 12-14):
  Phase 12 (Multi-Plan Support)        ▎ ~2 hrs
  Phase 13 (Refunds + Cancellations)   ▎ ~3 hrs
  Phase 14 (Performance + Caching)     ▎ ~2 hrs
```

### Additional Files for Phases 10-14

```
backend/
├── internal/service/
│   ├── paymob_client_test.go           [CREATE] Phase 10
│   ├── webhook_retry_worker.go         [CREATE] Phase 11
│   ├── error_tracker.go                [CREATE] Phase 11
│   ├── subscription_expiry_worker.go   [CREATE] Phase 13
│   └── cache.go                        [CREATE] Phase 14
│
├── internal/usecase/
│   └── payment_usecase_test.go         [CREATE] Phase 10
│
├── internal/handler/http/
│   └── payment_handler_test.go         [CREATE] Phase 10
│
├── internal/repository/
│   └── plan_repo.go                    [CREATE] Phase 12
│
├── migrations/
│   ├── 117_create_webhook_retry_queue.sql  [CREATE] Phase 11
│   └── 118_create_subscription_plans_table.sql  [CREATE] Phase 12
│
└── tests/
    ├── paymob_api_tests.json           [CREATE] Phase 10 (Postman collection)
    └── env.json                        [CREATE] Phase 10 (Postman env)

frontend/apps/dashboard/
└── (billing page updates for multi-plan and cancellation UI)

Total additional files (Phases 10-14): ~12
Grand total all phases: ~36 files
```

---

## APPENDIX A — Existing Code Patterns Reference

> **This section provides the EXACT code patterns from the existing codebase so an implementing agent can write code that matches the established conventions WITHOUT needing to search or read files.**

### A.1: Backend Handler Pattern

Every handler follows this structure (from `subscription_handler.go`):

```go
package http

import (
    "log"
    "net/http"
    "strconv"
    "pos-saas/internal/middleware"
    "pos-saas/internal/usecase"
)

type PaymentHandler struct {
    useCase *usecase.PaymentUseCase
}

func NewPaymentHandler(useCase *usecase.PaymentUseCase) *PaymentHandler {
    return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
    log.Printf("[PaymentHandler] GetPlans called\n")
    // ... business logic ...
    OKResponse(w, data)   // 200 + { "success": true, "data": ... }
}
```

**Key conventions:**
- `[HandlerName]` prefix in all log.Printf calls
- Constructor returns pointer, takes `*usecase.XUseCase`
- Extract tenant ID: `tenantID := middleware.GetTenantID(r)` → returns `int64`, 0 = not found
- Extract super admin ID: `userID := middleware.GetSuperAdminID(r)` → returns `int64`
- Extract path param: `r.PathValue("id")` (Go 1.22 net/http built-in)
- Decode POST body: `if !DecodeJSONBody(w, r, &req) { return }` — returns false = already wrote error
- Get IP: `getIPAddress(r)` helper (checks X-Forwarded-For, X-Real-IP, RemoteAddr)
- Pagination: `page := GetQueryParamInt(r, "page", 1)` / `pageSize := GetQueryParamInt(r, "page_size", 20)`

### A.2: Response Helpers (from `utils.go`)

All response helpers are defined in `backend/internal/handler/http/utils.go`:

```go
// Standard response shape: { "success": bool, "data"?: any, "error"?: string }
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

// Success responses:
OKResponse(w, data)           // 200 + { success: true, data: ... }
CreatedResponse(w, data)      // 201 + { success: true, data: ... }
NoContentResponse(w)          // 204

// Error responses:
BadRequestResponse(w, "message")     // 400 + { success: false, error: "message" }
UnauthorizedResponse(w, "message")   // 401
ForbiddenResponse(w, "message")      // 403
NotFoundResponse(w, "message")       // 404
ConflictResponse(w, "message")       // 409
InternalErrorResponse(w, "message")  // 500

// Query helpers:
GetQueryParam(r, "name", "default")      // string
GetQueryParamInt(r, "page", 1)           // int
GetQueryParamBool(r, "active")           // bool

// JSON body decoder (DisallowUnknownFields is ON — extra fields = 400 error):
DecodeJSONBody(w, r, &req) bool  // false = already responded with error
```

### A.3: Repository Pattern

```go
package repository

import (
    "context"
    "database/sql"
    "pos-saas/internal/domain"
)

// Interface file (payment_repo.go):
type PaymentRepository interface {
    CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error)
    // ...
}

// Implementation file (payment_repo_impl.go):
type PaymentRepositoryImpl struct {
    db *sql.DB
}

// Constructor returns INTERFACE, not concrete type:
func NewPaymentRepository(db *sql.DB) PaymentRepository {
    return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error) {
    log.Printf("[PaymentRepo] CreatePayment - tenant_id: %d\n", payment.TenantID)

    var id int64
    err := r.db.QueryRowContext(ctx, `
        INSERT INTO payments (tenant_id, amount_cents, currency, status, merchant_order_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, payment.TenantID, payment.AmountCents, payment.Currency, payment.Status, payment.MerchantOrderID).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("failed to create payment: %w", err)
    }
    return id, nil
}
```

**Key conventions:**
- Always `context.Context` as first parameter
- PostgreSQL: `$1`, `$2` placeholders (not `?`)
- Use `r.db.QueryRowContext(ctx, ...)` for single row, `r.db.QueryContext(ctx, ...)` for multiple
- Nullable fields: `sql.NullString`, `sql.NullTime`, `sql.NullInt64`
- Logging prefix: `[RepoName]`

### A.4: Route Registration in main.go

```go
// Initialization pattern (nil-guarded):
var paymentRepo repository.PaymentRepository
if db != nil {
    paymentRepo = repository.NewPaymentRepository(db)
}

var paymentUC *usecase.PaymentUseCase
if db != nil && paymentRepo != nil {
    paymentUC = usecase.NewPaymentUseCase(paymentRepo, subscriptionRepo, paymobClient)
}

var paymentHandler *handler.PaymentHandler
if paymentUC != nil {
    paymentHandler = handler.NewPaymentHandler(paymentUC)
}

// Route registration:
if paymentHandler != nil {
    // Public (no auth):
    mux.HandleFunc("GET /api/v1/subscription/plans", paymentHandler.GetPlans)
    mux.HandleFunc("POST /api/v1/webhook/paymob", paymentHandler.HandlePaymobWebhook)

    // Payment routes (auth + tenant context, but NO status check):
    // Use wrapProtected, NOT wrapProtectedWithStatus
    mux.Handle("POST /api/v1/subscription/initiate",
        wrapProtected(http.HandlerFunc(paymentHandler.InitiatePayment)))
    mux.Handle("GET /api/v1/subscription/payments",
        wrapProtected(http.HandlerFunc(paymentHandler.GetPaymentHistory)))
    mux.Handle("GET /api/v1/subscription/current",
        wrapProtected(http.HandlerFunc(paymentHandler.GetCurrentSubscription)))
}
```

**Middleware wrappers available:**

| Wrapper | Middleware Chain | Use For |
|---------|-----------------|---------|
| `wrapProtected` | Auth → TenantContext | Routes that need auth but NOT subscription check (payment routes!) |
| `wrapProtectedWithStatus` | Auth → TenantContext → TenantStatus | Standard dashboard routes |
| `wrapWithPermission(h, moduleID, level)` | Auth → TenantContext → RBAC | RBAC-protected routes |
| `wrapSuperAdmin` | RateLimit → SuperAdminAuth → IPWhitelist → RequireSuperAdmin | Super admin routes |
| none (raw `mux.HandleFunc`) | No middleware | Public endpoints |

### A.5: Domain Model Pattern

```go
package domain

import "time"

// Constants grouped by type:
const (
    PaymentStatusPending  = "pending"
    PaymentStatusSuccess  = "success"
    PaymentStatusFailed   = "failed"
    PaymentStatusRefunded = "refunded"
)

// Main struct with JSON tags:
type Payment struct {
    ID          int64      `json:"id"`
    TenantID    int64      `json:"tenant_id"`
    AmountCents int        `json:"amount_cents"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    PaidAt      *time.Time `json:"paid_at,omitempty"`  // nullable = pointer + omitempty
}

// Param structs (for usecase methods, not JSON — no json tags needed):
type InitiatePaymentParams struct {
    TenantID  int64
    UserID    int64
    PlanID    string
    IPAddress string
    UserAgent string
}
```

### A.6: Tenant Status Middleware — Skip Paths

The `TenantStatusMiddleware` in `backend/internal/middleware/tenant_status.go` skips these paths:

```go
// Exact matches:
"/api/v1/tenant/status"
"/api/v1/auth/logout"
"/api/v1/auth/refresh"
"/api/v1/auth/validate-token"

// Prefix matches:
"/api/v1/public/"
"/api/v1/super-admin/"
"/api/v1/auth/"
"/health"
"/metrics"
```

**For Paymob:** Payment routes using `wrapProtected` (instead of `wrapProtectedWithStatus`) automatically bypass the status middleware. The webhook endpoint is public (no middleware at all). No changes to `tenant_status.go` are needed.

### A.7: Frontend API Client Pattern

```typescript
// frontend/apps/dashboard/src/lib/api.ts
import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor auto-adds: Authorization, X-Tenant-ID, X-Restaurant-ID
// Response interceptor: 401 → logout + redirect to landing/auth/login

// Existing exported function pattern:
export async function getTenantStatus() {
  const response = await api.get('/tenant/status')
  return response.data
}

// Add new functions following same pattern:
export async function getSubscriptionPlans() {
  const response = await api.get('/subscription/plans')
  return response.data
}

export async function initiatePayment(planId: string) {
  const response = await api.post('/subscription/initiate', { plan_id: planId })
  return response.data
}

export async function getPaymentHistory(page = 1, pageSize = 10) {
  const response = await api.get(`/subscription/payments?page=${page}&page_size=${pageSize}`)
  return response.data
}

export async function getCurrentSubscription() {
  const response = await api.get('/subscription/current')
  return response.data
}
```

### A.8: Tenant Status Store Pattern

```typescript
// frontend/apps/dashboard/src/stores/tenantStatusStore.ts
import { create } from 'zustand'
import { api } from '@/lib/api'

// State shape:
{
  status: 'trial' | 'active' | 'expired' | 'suspended' | 'disabled' | null,
  trialEndDate: string | null,
  daysRemaining: number | null,
  canAccess: boolean,        // default: true (optimistic)
  blockReason: string | null, // 'trial_expired' | 'tenant_disabled' | 'tenant_not_found'
  isLoading: boolean,
  error: string | null,
  lastChecked: Date | null,
}

// checkStatus() calls GET /tenant/status
// On success: sets all fields from response.data.data
// On 403 with TENANT_ACCESS_BLOCKED: sets canAccess=false, reads blockReason from error.response.data.error.reason
```

**After successful payment, call:**
```typescript
useTenantStatusStore.getState().checkStatus()
```

### A.9: useTenantStatus Hook

```typescript
// frontend/apps/dashboard/src/hooks/useTenantStatus.ts
export function useTenantStatus() {
  const store = useTenantStatusStore()
  return {
    ...store,
    isBlocked: !store.canAccess,  // computed
    getBlockMessage: () => {
      switch (store.blockReason) {
        case 'trial_expired': return 'Your trial period has expired...'
        case 'tenant_disabled': return 'Your account has been disabled...'
        default: return 'Access restricted...'
      }
    }
  }
}
```

### A.10: BlockedOverlay Current Code

```typescript
// frontend/apps/dashboard/src/components/blocked/BlockedOverlay.tsx
interface BlockedOverlayProps {
  blockReason: 'trial_expired' | 'tenant_disabled' | string
  supportEmail?: string  // default: 'support@pos-saas.com'
}

// Currently renders:
// - trial_expired: Clock icon (orange), "Your Trial Has Expired", mailto:support CTA
// - tenant_disabled: Lock icon (red), "Account Access Disabled", mailto:support CTA
// - default: Lock icon (red), "Access Restricted", mailto:support CTA
// + Sign Out link at bottom

// TO MODIFY: Replace mailto CTA in trial_expired case with:
//   1. Primary button: "Subscribe Now — 500 EGP/month" → navigates to /settings/billing
//   2. Secondary link: "Contact Support" (keep mailto but smaller)
```

### A.11: Settings Sidebar — How to Add a Nav Item

```typescript
// frontend/apps/dashboard/src/components/settings/SettingsSidebar.tsx
// Props: { locale: string }

// Existing links array:
const settingsLinks = [
  { href: `/${locale}/dashboard/settings/account`,        label: 'Account Settings',  icon: User,    description: '...' },
  { href: `/${locale}/dashboard/settings/restaurant`,     label: 'Restaurant',         icon: Store,   description: '...' },
  { href: `/${locale}/dashboard/settings/appearance`,     label: 'Appearance',         icon: Palette, description: '...' },
  { href: `/${locale}/dashboard/settings/access-control`, label: 'Access Control',     icon: Shield,  description: '...' },
]

// ADD this to the array:
{ href: `/${locale}/dashboard/settings/billing`, label: 'Billing & Subscription', icon: CreditCard, description: 'Manage payments and subscription' }

// Import CreditCard from lucide-react
```

### A.12: SidebarSubscriptionBadge — Current Code

```typescript
// Located at: frontend/apps/dashboard/src/components/layout/SidebarSubscriptionBadge.tsx
// OR inline in Sidebar.tsx

// Uses: useTenantStatus() hook
// Renders colored badge at bottom of sidebar:
//   active → green "Pro Plan" with Crown icon
//   trial > 7 days → blue "Pro Trial" + "X days left"
//   trial 4-7 days → yellow
//   trial 1-3 days → red with AlertTriangle
//   trial 0 days → red "Trial Expired"
//   expired → red "Expired"

// TO MODIFY: Wrap in <Link> or add onClick → navigate to billing page
// On expired: add orange "Subscribe" styling
```

### A.13: DashboardLayout — Where Blocked State is Handled

The blocked overlay is rendered in `DashboardLayout.tsx` (NOT in `DashboardAuthProvider.tsx`). The auth provider just calls `checkStatus()`. The layout component checks `isBlocked` and renders either `<BlockedOverlay>` or children.

**To allow billing page when blocked:** In `DashboardLayout.tsx`, add a pathname check:
```typescript
const pathname = usePathname()
const isBillingPage = pathname.includes('/settings/billing')

// Render children (not overlay) if: !isBlocked OR isBillingPage
```

### A.14: i18n Translation File Structure

```json
// frontend/apps/dashboard/src/i18n/messages/en.json
{
  "common": { "dashboard": "Dashboard", "save": "Save", ... },
  "sidebar": { ... },
  "products": { ... },
  "categories": { ... },
  "orders": { ... },
  "hr": { ... },
  "settings": { ... },
  "themeBuilder": { ... },
  "notifications": { ... }
}

// ADD new top-level key:
{
  "billing": {
    "title": "Billing & Subscription",
    "currentPlan": "Current Plan",
    ...
  }
}

// Arabic file: frontend/apps/dashboard/src/i18n/messages/ar.json
// Same structure, same keys, Arabic values
```

### A.15: Existing Subscription Domain Constants

```go
// backend/internal/domain/subscription.go
// Statuses:
StatusTrial     = "trial"
StatusActive    = "active"
StatusExpired   = "expired"
StatusSuspended = "suspended"
StatusDisabled  = "disabled"

// History actions (existing):
ActionTrialStarted      = "TRIAL_STARTED"
ActionTrialExtended     = "TRIAL_EXTENDED"
ActionTrialExpired      = "TRIAL_EXPIRED"
ActionTenantDisabled    = "TENANT_DISABLED"
ActionTenantEnabled     = "TENANT_ENABLED"
ActionDeletionScheduled = "DELETION_SCHEDULED"
ActionDeletionCancelled = "DELETION_CANCELLED"
ActionTenantDeleted     = "TENANT_DELETED"
ActionStatusChanged     = "STATUS_CHANGED"

// NEW constants to add in domain/payment.go:
ActionPaymentInitiated      = "PAYMENT_INITIATED"
ActionPaymentSuccess        = "PAYMENT_SUCCESS"
ActionPaymentFailed         = "PAYMENT_FAILED"
ActionSubscriptionActivated = "SUBSCRIPTION_ACTIVATED"
ActionSubscriptionExpired   = "SUBSCRIPTION_EXPIRED"
ActionPaymentRefunded       = "PAYMENT_REFUNDED"
ActionSubscriptionCancelled = "SUBSCRIPTION_CANCELLED"
```

### A.16: Existing TenantSubscriptionInfo Fields

The `GET /api/v1/tenant/status` response currently returns:

```json
{
  "success": true,
  "data": {
    "tenant_id": 1,
    "tenant_name": "Demo Restaurant",
    "slug": "demo",
    "subscription_status": "trial",
    "trial_start_date": "2026-01-20T00:00:00Z",
    "trial_end_date": "2026-02-19T00:00:00Z",
    "trial_days": 30,
    "days_remaining": 10,
    "is_disabled": false,
    "can_access": true,
    "contact_email": "demo@example.com",
    "contact_phone": "+201234567890"
  }
}
```

**After Paymob integration, extend with:**
```json
{
  "subscription_expires_at": "2026-03-21T00:00:00Z",
  "last_payment_at": "2026-02-21T00:00:00Z",
  "plan": {
    "id": "monthly",
    "name": "Pro Plan",
    "price": 500,
    "currency": "EGP"
  }
}
```

### A.17: Database Connection Access

```go
// In main.go, the DB is created as:
db, err := sql.Open("postgres", dsn)
// DSN from env: "host=... port=... user=... password=... dbname=... sslmode=..."

// All repos receive *sql.DB in their constructor:
paymentRepo := repository.NewPaymentRepository(db)
```

### A.18: Paymob HMAC-SHA512 Field Order (CRITICAL)

When verifying Paymob webhook HMAC, concatenate these 20 fields from the `obj` object **in this EXACT alphabetical order**:

```
1.  amount_cents          (string of int)
2.  created_at            (string)
3.  currency              (string)
4.  error_occured         (string of bool: "true"/"false")
5.  has_parent_transaction(string of bool)
6.  id                    (string of int)
7.  integration_id        (string of int)
8.  is_3d_secure          (string of bool)
9.  is_auth               (string of bool)
10. is_capture             (string of bool)
11. is_refunded            (string of bool)
12. is_standalone_payment  (string of bool)
13. is_voided              (string of bool)
14. order.id               (string of int)
15. owner                  (string of int)
16. pending                (string of bool)
17. source_data.pan        (string — last 4 digits like "xxxx-xxxx-xxxx-1234")
18. source_data.sub_type   (string — "MasterCard", "Visa", etc.)
19. source_data.type       (string — "card", etc.)
20. success                (string of bool)
```

Concatenate all values (no separators), then:
```go
mac := hmac.New(sha512.New, []byte(hmacSecret))
mac.Write([]byte(concatenated))
computed := hex.EncodeToString(mac.Sum(nil))
return computed == receivedHMAC
```
