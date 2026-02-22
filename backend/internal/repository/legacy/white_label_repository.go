//go:build legacy
// +build legacy

package repository

import (
	"context"
	"time"
)

// WhiteLabelRepository handles white-label data persistence
type WhiteLabelRepository interface {
	// Branding operations
	CreateBranding(ctx context.Context, branding *models.BrandingConfig) (*models.BrandingConfig, error)
	GetBranding(ctx context.Context, restaurantID int64) (*models.BrandingConfig, error)
	UpdateBranding(ctx context.Context, branding *models.BrandingConfig) error
	DeleteBranding(ctx context.Context, restaurantID int64) error

	// Theme presets
	GetThemePresets(ctx context.Context) ([]*models.ThemePreset, error)
	GetThemePreset(ctx context.Context, id int64) (*models.ThemePreset, error)
	CreateThemePreset(ctx context.Context, preset *models.ThemePreset) (*models.ThemePreset, error)
	UpdateThemePreset(ctx context.Context, preset *models.ThemePreset) error

	// Reseller partner operations
	CreateResellerPartner(ctx context.Context, partner *models.ResellerPartner) (*models.ResellerPartner, error)
	GetResellerPartner(ctx context.Context, id int64) (*models.ResellerPartner, error)
	GetResellerPartners(ctx context.Context, status string, limit int, offset int) ([]*models.ResellerPartner, error)
	GetResellerPartnerByEmail(ctx context.Context, email string) (*models.ResellerPartner, error)
	UpdateResellerPartner(ctx context.Context, partner *models.ResellerPartner) error
	DeleteResellerPartner(ctx context.Context, id int64) error

	// Reseller portal user operations
	CreateResellerPortalUser(ctx context.Context, user *models.ResellerPortalUser) (*models.ResellerPortalUser, error)
	GetResellerPortalUser(ctx context.Context, id int64) (*models.ResellerPortalUser, error)
	GetResellerPortalUserByEmail(ctx context.Context, email string) (*models.ResellerPortalUser, error)
	GetResellerPortalUsers(ctx context.Context, resellerID int64) ([]*models.ResellerPortalUser, error)
	UpdateResellerPortalUser(ctx context.Context, user *models.ResellerPortalUser) error
	UpdateResellerPortalUserLogin(ctx context.Context, id int64) error
	DeleteResellerPortalUser(ctx context.Context, id int64) error

	// License operations
	CreateLicenseKey(ctx context.Context, license *models.LicenseKey) (*models.LicenseKey, error)
	GetLicenseKey(ctx context.Context, id int64) (*models.LicenseKey, error)
	GetLicenseKeyByKey(ctx context.Context, key string) (*models.LicenseKey, error)
	GetRestaurantLicenses(ctx context.Context, restaurantID int64) ([]*models.LicenseKey, error)
	UpdateLicenseKey(ctx context.Context, license *models.LicenseKey) error
	UpdateLicenseKeyStatus(ctx context.Context, id int64, status string) error
	DeleteLicenseKey(ctx context.Context, id int64) error

	// Marketplace operations
	CreateMarketplace(ctx context.Context, marketplace *models.Marketplace) (*models.Marketplace, error)
	GetMarketplace(ctx context.Context, id int64) (*models.Marketplace, error)
	GetMarketplaceBySlug(ctx context.Context, slug string) (*models.Marketplace, error)
	GetRestaurantMarketplaces(ctx context.Context, restaurantID int64) ([]*models.Marketplace, error)
	UpdateMarketplace(ctx context.Context, marketplace *models.Marketplace) error
	DeleteMarketplace(ctx context.Context, id int64) error

	// Marketplace product operations
	CreateMarketplaceProduct(ctx context.Context, product *models.MarketplaceProduct) (*models.MarketplaceProduct, error)
	GetMarketplaceProduct(ctx context.Context, id int64) (*models.MarketplaceProduct, error)
	GetMarketplaceProducts(ctx context.Context, marketplaceID int64, status string, limit int, offset int) ([]*models.MarketplaceProduct, error)
	GetSellerProducts(ctx context.Context, sellerID int64, limit int, offset int) ([]*models.MarketplaceProduct, error)
	SearchMarketplaceProducts(ctx context.Context, marketplaceID int64, query string, category string) ([]*models.MarketplaceProduct, error)
	UpdateMarketplaceProduct(ctx context.Context, product *models.MarketplaceProduct) error
	UpdateProductStock(ctx context.Context, productID int64, quantity int64) error
	UpdateProductViewCount(ctx context.Context, productID int64) error
	DeleteMarketplaceProduct(ctx context.Context, id int64) error

	// Marketplace seller operations
	CreateMarketplaceSeller(ctx context.Context, seller *models.MarketplaceSeller) (*models.MarketplaceSeller, error)
	GetMarketplaceSeller(ctx context.Context, id int64) (*models.MarketplaceSeller, error)
	GetMarketplaceSellerByStore(ctx context.Context, marketplaceID, storeID int64) (*models.MarketplaceSeller, error)
	GetMarketplaceSellers(ctx context.Context, marketplaceID int64, status string, limit int, offset int) ([]*models.MarketplaceSeller, error)
	UpdateMarketplaceSeller(ctx context.Context, seller *models.MarketplaceSeller) error
	DeleteMarketplaceSeller(ctx context.Context, id int64) error

	// Customization requests
	CreateCustomizationRequest(ctx context.Context, request *models.CustomizationRequest) (*models.CustomizationRequest, error)
	GetCustomizationRequest(ctx context.Context, id int64) (*models.CustomizationRequest, error)
	GetCustomizationRequests(ctx context.Context, restaurantID int64, status string) ([]*models.CustomizationRequest, error)
	UpdateCustomizationRequest(ctx context.Context, request *models.CustomizationRequest) error

	// White-label license operations
	CreateWhiteLabelLicense(ctx context.Context, license *models.WhiteLabelLicense) (*models.WhiteLabelLicense, error)
	GetWhiteLabelLicense(ctx context.Context, id int64) (*models.WhiteLabelLicense, error)
	GetRestaurantWhiteLabelLicense(ctx context.Context, restaurantID int64) (*models.WhiteLabelLicense, error)
	UpdateWhiteLabelLicense(ctx context.Context, license *models.WhiteLabelLicense) error
	GetExpiringSoonLicenses(ctx context.Context, days int) ([]*models.WhiteLabelLicense, error)

	// Reseller commission operations
	CreateResellerCommission(ctx context.Context, commission *models.ResellerCommission) (*models.ResellerCommission, error)
	GetResellerCommissions(ctx context.Context, resellerID int64, period string) ([]*models.ResellerCommission, error)
	GetPendingCommissions(ctx context.Context, resellerID int64) ([]*models.ResellerCommission, error)
	UpdateCommissionStatus(ctx context.Context, id int64, status string) error
	GetResellerTotalCommission(ctx context.Context, resellerID int64, start, end time.Time) (float64, error)

	// Reseller dashboard metrics
	CreateResellerMetrics(ctx context.Context, metrics *models.ResellerDashboardMetrics) (*models.ResellerDashboardMetrics, error)
	GetResellerMetrics(ctx context.Context, resellerID int64, period string, limit int) ([]*models.ResellerDashboardMetrics, error)
	GetResellerMetricsByDate(ctx context.Context, resellerID int64, date time.Time) (*models.ResellerDashboardMetrics, error)

	// Integration keys
	CreateIntegrationKey(ctx context.Context, key *models.IntegrationKey) (*models.IntegrationKey, error)
	GetIntegrationKey(ctx context.Context, id int64) (*models.IntegrationKey, error)
	GetIntegrationKeyByKey(ctx context.Context, key string) (*models.IntegrationKey, error)
	GetRestaurantIntegrationKeys(ctx context.Context, restaurantID int64) ([]*models.IntegrationKey, error)
	UpdateIntegrationKey(ctx context.Context, key *models.IntegrationKey) error
	DeleteIntegrationKey(ctx context.Context, id int64) error
	UpdateIntegrationKeyLastUsed(ctx context.Context, id int64) error

	// Webhook configuration
	CreateWebhookConfig(ctx context.Context, webhook *models.WebhookConfig) (*models.WebhookConfig, error)
	GetWebhookConfig(ctx context.Context, id int64) (*models.WebhookConfig, error)
	GetRestaurantWebhooks(ctx context.Context, restaurantID int64) ([]*models.WebhookConfig, error)
	UpdateWebhookConfig(ctx context.Context, webhook *models.WebhookConfig) error
	DeleteWebhookConfig(ctx context.Context, id int64) error
	GetWebhooksForEvent(ctx context.Context, eventType string) ([]*models.WebhookConfig, error)
}
