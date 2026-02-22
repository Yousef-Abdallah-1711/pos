package domain

import "time"

// Payment statuses.
const (
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

// Subscription history actions related to payments.
const (
	ActionPaymentInitiated      = "PAYMENT_INITIATED"
	ActionPaymentSuccess        = "PAYMENT_SUCCESS"
	ActionPaymentFailed         = "PAYMENT_FAILED"
	ActionSubscriptionActivated = "SUBSCRIPTION_ACTIVATED"
	ActionSubscriptionExpired   = "SUBSCRIPTION_EXPIRED"
	ActionPaymentRefunded       = "PAYMENT_REFUNDED"
	ActionSubscriptionCancelled = "SUBSCRIPTION_CANCELLED"
)

// SubscriptionPlan represents an available billing plan.
type SubscriptionPlan struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	NameAr     string   `json:"name_ar"`
	Price      int      `json:"price"`
	PriceCents int      `json:"price_cents"`
	Currency   string   `json:"currency"`
	Period     string   `json:"period"`
	Features   []string `json:"features"`
	FeaturesAr []string `json:"features_ar"`
}

// SubscriptionPlanSummary is used in current subscription responses.
type SubscriptionPlanSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	PriceCents int    `json:"price_cents"`
	Currency   string `json:"currency"`
	Period     string `json:"period"`
}

// Payment represents one payment record.
type Payment struct {
	ID                  int64                  `json:"id"`
	TenantID            int64                  `json:"tenant_id"`
	UserID              *int64                 `json:"user_id,omitempty"`
	AmountCents         int                    `json:"amount_cents"`
	Currency            string                 `json:"currency"`
	Status              string                 `json:"status"`
	PaymobOrderID       *int64                 `json:"paymob_order_id,omitempty"`
	PaymobTransactionID *int64                 `json:"paymob_transaction_id,omitempty"`
	PaymentMethod       *string                `json:"payment_method,omitempty"`
	CardLastFour        *string                `json:"card_last_four,omitempty"`
	CardBrand           *string                `json:"card_brand,omitempty"`
	BillingName         *string                `json:"billing_name,omitempty"`
	BillingEmail        *string                `json:"billing_email,omitempty"`
	BillingPhone        *string                `json:"billing_phone,omitempty"`
	PlanID              string                 `json:"plan_id"`
	PlanPeriodStart     *time.Time             `json:"plan_period_start,omitempty"`
	PlanPeriodEnd       *time.Time             `json:"plan_period_end,omitempty"`
	MerchantOrderID     string                 `json:"merchant_order_id"`
	ErrorMessage        *string                `json:"error_message,omitempty"`
	ErrorCode           *string                `json:"error_code,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	PaidAt              *time.Time             `json:"paid_at,omitempty"`
}

// TenantBillingInfo stores payment-related fields from tenants table.
type TenantBillingInfo struct {
	SubscriptionPlan        string
	SubscriptionExpiresAt   *time.Time
	SubscriptionCancelledAt *time.Time
	LastPaymentAt           *time.Time
	LastPaymentAmount       *int
	LastPaymentCurrency     string
}

// CurrentSubscriptionResponse is returned by GET /subscription/current.
type CurrentSubscriptionResponse struct {
	SubscriptionStatus      string                    `json:"subscription_status"`
	TrialEndDate            *time.Time                `json:"trial_end_date,omitempty"`
	SubscriptionExpiresAt   *time.Time                `json:"subscription_expires_at,omitempty"`
	SubscriptionCancelledAt *time.Time                `json:"subscription_cancelled_at,omitempty"`
	CancelAtPeriodEnd       bool                      `json:"cancel_at_period_end"`
	LastPaymentAt           *time.Time                `json:"last_payment_at,omitempty"`
	DaysRemaining           int                       `json:"days_remaining"`
	CanAccess               bool                      `json:"can_access"`
	BlockReason             *string                   `json:"block_reason,omitempty"`
	Plan                    SubscriptionPlanSummary   `json:"plan"`
	CurrentPlan             SubscriptionPlanSummary   `json:"current_plan"`
	AvailablePlans          []SubscriptionPlanSummary `json:"available_plans"`
}

// PaymobAuthResponse maps Paymob auth response payload.
type PaymobAuthResponse struct {
	Token string `json:"token"`
}

// PaymobOrderResponse maps Paymob order creation response payload.
type PaymobOrderResponse struct {
	ID int64 `json:"id"`
}

// PaymobPaymentKeyResponse maps Paymob payment key response payload.
type PaymobPaymentKeyResponse struct {
	Token string `json:"token"`
}

// PaymobWebhookPayload maps the transaction callback payload.
type PaymobWebhookPayload struct {
	Obj struct {
		ID                   int64  `json:"id"`
		AmountCents          int    `json:"amount_cents"`
		Success              bool   `json:"success"`
		Currency             string `json:"currency"`
		IsAuth               bool   `json:"is_auth"`
		IsCapture            bool   `json:"is_capture"`
		IsRefunded           bool   `json:"is_refunded"`
		IsVoided             bool   `json:"is_voided"`
		Is3DSecure           bool   `json:"is_3d_secure"`
		IsStandalonePayment  bool   `json:"is_standalone_payment"`
		HasParentTransaction bool   `json:"has_parent_transaction"`
		ErrorOccured         bool   `json:"error_occured"`
		Pending              bool   `json:"pending"`
		IntegrationID        int64  `json:"integration_id"`
		Owner                int64  `json:"owner"`
		CreatedAt            string `json:"created_at"`
		Order                struct {
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

// InitiatePaymentRequest represents POST /subscription/initiate body.
type InitiatePaymentRequest struct {
	PlanID string `json:"plan_id"`
}

// InitiatePaymentResponse represents initiate payment response.
type InitiatePaymentResponse struct {
	PaymentID       int64  `json:"payment_id"`
	IframeURL       string `json:"iframe_url"`
	PaymobOrderID   int64  `json:"paymob_order_id"`
	MerchantOrderID string `json:"merchant_order_id"`
}

// CancelSubscriptionRequest represents POST /subscription/cancel body.
type CancelSubscriptionRequest struct {
	Reason string `json:"reason,omitempty"`
}

// CancelSubscriptionResponse represents cancellation scheduling response.
type CancelSubscriptionResponse struct {
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	CancelledAt           time.Time  `json:"cancelled_at"`
	Message               string     `json:"message"`
}

// RefundPaymentResponse represents super-admin refund response.
type RefundPaymentResponse struct {
	PaymentID          int64     `json:"payment_id"`
	TenantID           int64     `json:"tenant_id"`
	Status             string    `json:"status"`
	SubscriptionStatus string    `json:"subscription_status"`
	RefundedAt         time.Time `json:"refunded_at"`
	AmountCents        int       `json:"amount_cents"`
	Currency           string    `json:"currency"`
}
