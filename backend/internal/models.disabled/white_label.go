//go:build legacy
// +build legacy

package models

import (
	"time"
)

// BrandingConfig holds custom branding for a white-label instance
type BrandingConfig struct {
	ID              int64     `json:"id"`
	RestaurantID    int64     `json:"restaurant_id"`
	CompanyName     string    `json:"company_name"`
	CompanyLogo     string    `json:"company_logo"` // URL to logo
	PrimaryColor    string    `json:"primary_color"` // Hex color
	SecondaryColor  string    `json:"secondary_color"`
	AccentColor     string    `json:"accent_color"`
	FontFamily      string    `json:"font_family"`
	BackgroundColor string    `json:"background_color"`
	CustomDomain    string    `json:"custom_domain"` // e.g., restaurant.mypos.com
	CustomCSS       string    `json:"custom_css"`    // Custom CSS overrides
	FaviconURL      string    `json:"favicon_url"`
	EmailTemplate   string    `json:"email_template"`
	FooterText      string    `json:"footer_text"`
	PrivacyURL      string    `json:"privacy_url"`
	TermsURL        string    `json:"terms_url"`
	SupportEmail    string    `json:"support_email"`
	SupportPhone    string    `json:"support_phone"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ThemePreset represents a pre-designed theme template
type ThemePreset struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PrimaryColor    string    `json:"primary_color"`
	SecondaryColor  string    `json:"secondary_color"`
	AccentColor     string    `json:"accent_color"`
	FontFamily      string    `json:"font_family"`
	BackgroundColor string    `json:"background_color"`
	PreviewImage    string    `json:"preview_image"`
	Category        string    `json:"category"` // "modern", "classic", "minimal", "dark"
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// ResellerPartner represents a partner or reseller
type ResellerPartner struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Email             string          `json:"email"`
	Phone             string          `json:"phone"`
	CompanyName       string          `json:"company_name"`
	Address           string          `json:"address"`
	City              string          `json:"city"`
	State             string          `json:"state"`
	Country           string          `json:"country"`
	PostalCode        string          `json:"postal_code"`
	ContactPerson     string          `json:"contact_person"`
	Website           string          `json:"website"`
	CommissionRate    float64         `json:"commission_rate"` // 0.1 = 10%
	Status            string          `json:"status"`          // "pending", "approved", "active", "suspended", "inactive"
	PartnerType       string          `json:"partner_type"`    // "reseller", "affiliate", "technology", "referral"
	TierLevel         string          `json:"tier_level"`      // "bronze", "silver", "gold", "platinum"
	MonthlyRevenue    float64         `json:"monthly_revenue"`
	CustomersCount    int64           `json:"customers_count"`
	ApprovedAt        *time.Time      `json:"approved_at"`
	SuspendedAt       *time.Time      `json:"suspended_at"`
	LastActivityAt    *time.Time      `json:"last_activity_at"`
	BrandingConfig    *BrandingConfig `json:"branding_config"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ResellerPortalUser represents a user in the reseller portal
type ResellerPortalUser struct {
	ID           int64     `json:"id"`
	ResellerID   int64     `json:"reseller_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"` // "admin", "manager", "analyst", "support"
	Permissions  []string  `json:"permissions"`
	IsActive     bool      `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LicenseKey represents a license for a white-label instance
type LicenseKey struct {
	ID                int64     `json:"id"`
	RestaurantID      int64     `json:"restaurant_id"`
	ResellerID        *int64    `json:"reseller_id"`
	LicenseKey        string    `json:"license_key"`
	LicenseType       string    `json:"license_type"`       // "basic", "professional", "enterprise"
	MaxInstances      int       `json:"max_instances"`
	MaxUsers          int       `json:"max_users"`
	MaxCustomers      int       `json:"max_customers"`
	Features          []string  `json:"features"`            // Enabled features
	ExpiresAt         time.Time `json:"expires_at"`
	IssuedAt          time.Time `json:"issued_at"`
	LastValidatedAt   *time.Time `json:"last_validated_at"`
	Status            string    `json:"status"`              // "valid", "expired", "suspended", "revoked"
	ActivationCode    *string   `json:"activation_code"`
	ActivatedAt       *time.Time `json:"activated_at"`
	RenewalDate       *time.Time `json:"renewal_date"`
	AutoRenewEnabled  bool      `json:"auto_renew_enabled"`
	NotificationCount int       `json:"notification_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Marketplace represents a white-label marketplace
type Marketplace struct {
	ID                 int64     `json:"id"`
	RestaurantID       int64     `json:"restaurant_id"`
	Name               string    `json:"name"`
	Slug               string    `json:"slug"` // URL-friendly identifier
	Description        string    `json:"description"`
	LogoURL            string    `json:"logo_url"`
	BannerURL          string    `json:"banner_url"`
	IsActive           bool      `json:"is_active"`
	AllowUserCreation  bool      `json:"allow_user_creation"`
	RequiresApproval   bool      `json:"requires_approval"`
	CommissionRate     float64   `json:"commission_rate"`
	MinimumPrice       float64   `json:"minimum_price"`
	MaximumPrice       float64   `json:"maximum_price"`
	CurrencyCode       string    `json:"currency_code"` // "USD", "EUR", etc.
	PaymentMethods     []string  `json:"payment_methods"`
	Categories         []string  `json:"categories"`
	FeaturedProductIDs []int64   `json:"featured_product_ids"`
	ProductCount       int64     `json:"product_count"`
	SellerCount        int64     `json:"seller_count"`
	TotalRevenue       float64   `json:"total_revenue"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// MarketplaceProduct represents a product in white-label marketplace
type MarketplaceProduct struct {
	ID               int64      `json:"id"`
	MarketplaceID    int64      `json:"marketplace_id"`
	SellerID         int64      `json:"seller_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	SKU              string     `json:"sku"`
	Category         string     `json:"category"`
	Price            float64    `json:"price"`
	SalePrice        *float64   `json:"sale_price"`
	ImageURL         string     `json:"image_url"`
	Images           []string   `json:"images"`
	Stock            int64      `json:"stock"`
	MinimumStock     int64      `json:"minimum_stock"`
	Rating           float64    `json:"rating"`
	ReviewCount      int64      `json:"review_count"`
	Status           string     `json:"status"` // "draft", "published", "inactive", "archived"
	VisibilityLevel  string     `json:"visibility_level"` // "public", "private", "restricted"
	IsApproved       bool       `json:"is_approved"`
	ApprovedAt       *time.Time `json:"approved_at"`
	RejectionReason  *string    `json:"rejection_reason"`
	Tags             []string   `json:"tags"`
	ViewCount        int64      `json:"view_count"`
	SalesCount       int64      `json:"sales_count"`
	Revenue          float64    `json:"revenue"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// MarketplaceSeller represents a seller in the marketplace
type MarketplaceSeller struct {
	ID               int64     `json:"id"`
	MarketplaceID    int64     `json:"marketplace_id"`
	StoreID          int64     `json:"store_id"`
	StoreName        string    `json:"store_name"`
	StoreDescription string    `json:"store_description"`
	StoreIcon        string    `json:"store_icon"`
	StoreBanner      string    `json:"store_banner"`
	Rating           float64   `json:"rating"`
	ReviewCount      int64     `json:"review_count"`
	Status           string    `json:"status"` // "pending", "approved", "active", "suspended"
	Joined           time.Time `json:"joined"`
	Products         int64     `json:"products"`
	TotalSales       float64   `json:"total_sales"`
	MonthlyRevenue   float64   `json:"monthly_revenue"`
	Commission       float64   `json:"commission"` // This seller's commission rate
	VerificationCode *string   `json:"verification_code"`
	VerifiedAt       *time.Time `json:"verified_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CustomizationRequest represents a request to customize branding
type CustomizationRequest struct {
	ID             int64     `json:"id"`
	RestaurantID   int64     `json:"restaurant_id"`
	RequestedBy    int64     `json:"requested_by"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	BrandingConfig *BrandingConfig `json:"branding_config"`
	Status         string    `json:"status"` // "pending", "in_progress", "completed", "rejected"
	ApprovedBy     *int64    `json:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at"`
	RejectionReason *string   `json:"rejection_reason"`
	Priority       string    `json:"priority"` // "low", "medium", "high"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// WhiteLabelLicense represents licensing agreement details
type WhiteLabelLicense struct {
	ID                  int64     `json:"id"`
	RestaurantID        int64     `json:"restaurant_id"`
	LicenseType         string    `json:"license_type"` // "trial", "starter", "professional", "enterprise"
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	Price               float64   `json:"price"`
	BillingCycle        string    `json:"billing_cycle"` // "monthly", "annual"
	IsActive            bool      `json:"is_active"`
	AutoRenew           bool      `json:"auto_renew"`
	Features            []string  `json:"features"`
	SupportLevel        string    `json:"support_level"` // "basic", "priority", "premium"
	MaxSeats            int       `json:"max_seats"`
	CustomBrandingIncluded bool `json:"custom_branding_included"`
	APIAccessIncluded   bool    `json:"api_access_included"`
	SSOIncluded         bool    `json:"sso_included"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ResellerCommission tracks commission earnings for resellers
type ResellerCommission struct {
	ID              int64     `json:"id"`
	ResellerID      int64     `json:"reseller_id"`
	RestaurantID    int64     `json:"restaurant_id"`
	TransactionID   int64     `json:"transaction_id"`
	CommissionType  string    `json:"commission_type"` // "sale", "referral", "license"
	Amount          float64   `json:"amount"`
	Rate            float64   `json:"rate"`
	Period          string    `json:"period"`     // "monthly", "weekly"
	Status          string    `json:"status"`     // "pending", "processing", "paid"
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ResellerDashboardMetrics tracks metrics for reseller dashboard
type ResellerDashboardMetrics struct {
	ID                int64         `json:"id"`
	ResellerID        int64         `json:"reseller_id"`
	Period            string        `json:"period"` // "daily", "weekly", "monthly"
	Date              time.Time     `json:"date"`
	NewCustomers      int64         `json:"new_customers"`
	TotalCustomers    int64         `json:"total_customers"`
	ActiveCustomers   int64         `json:"active_customers"`
	TotalRevenue      float64       `json:"total_revenue"`
	AverageOrderValue float64       `json:"average_order_value"`
	CommissionEarned  float64       `json:"commission_earned"`
	Conversions       int64         `json:"conversions"`
	ConversionRate    float64       `json:"conversion_rate"`
	CreatedAt         time.Time     `json:"created_at"`
}

// IntegrationKey represents an API key for white-label integrations
type IntegrationKey struct {
	ID              int64     `json:"id"`
	RestaurantID    int64     `json:"restaurant_id"`
	Name            string    `json:"name"`
	Key             string    `json:"key"` // Hashed API key
	KeyPreview      string    `json:"key_preview"`
	Permissions     []string  `json:"permissions"`
	IsActive        bool      `json:"is_active"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// WebhookConfig represents webhook configuration for events
type WebhookConfig struct {
	ID            int64     `json:"id"`
	RestaurantID  int64     `json:"restaurant_id"`
	URL           string    `json:"url"`
	EventTypes    []string  `json:"event_types"` // "product.created", "order.completed", etc.
	IsActive      bool      `json:"is_active"`
	Secret        string    `json:"secret"`
	RetryCount    int       `json:"retry_count"`
	MaxRetries    int       `json:"max_retries"`
	LastTriggered *time.Time `json:"last_triggered"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
