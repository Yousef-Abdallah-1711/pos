//go:build legacy
// +build legacy

package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"pos-saas/internal/repository"
)

// WhiteLabelService handles white-label branding and customization
type WhiteLabelService struct {
	repo repository.WhiteLabelRepository
}

// NewWhiteLabelService creates a new white-label service
func NewWhiteLabelService(repo repository.WhiteLabelRepository) *WhiteLabelService {
	return &WhiteLabelService{
		repo: repo,
	}
}

// CreateBranding creates custom branding for a restaurant
func (s *WhiteLabelService) CreateBranding(ctx context.Context, restaurantID int64, branding *models.BrandingConfig) (*models.BrandingConfig, error) {
	if branding.CompanyName == "" {
		return nil, fmt.Errorf("company name is required")
	}

	// Validate colors are hex format
	colors := []string{branding.PrimaryColor, branding.SecondaryColor, branding.AccentColor, branding.BackgroundColor}
	for _, color := range colors {
		if color != "" && !isValidHexColor(color) {
			return nil, fmt.Errorf("invalid hex color: %s", color)
		}
	}

	branding.RestaurantID = restaurantID
	branding.CreatedAt = time.Now()
	branding.UpdatedAt = time.Now()

	return s.repo.CreateBranding(ctx, branding)
}

// UpdateBranding updates existing branding
func (s *WhiteLabelService) UpdateBranding(ctx context.Context, restaurantID int64, branding *models.BrandingConfig) (*models.BrandingConfig, error) {
	existing, err := s.repo.GetBranding(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return s.CreateBranding(ctx, restaurantID, branding)
	}

	// Update fields
	if branding.CompanyName != "" {
		existing.CompanyName = branding.CompanyName
	}
	if branding.CompanyLogo != "" {
		existing.CompanyLogo = branding.CompanyLogo
	}
	if branding.PrimaryColor != "" {
		existing.PrimaryColor = branding.PrimaryColor
	}
	if branding.SecondaryColor != "" {
		existing.SecondaryColor = branding.SecondaryColor
	}
	if branding.AccentColor != "" {
		existing.AccentColor = branding.AccentColor
	}
	if branding.FontFamily != "" {
		existing.FontFamily = branding.FontFamily
	}
	if branding.BackgroundColor != "" {
		existing.BackgroundColor = branding.BackgroundColor
	}
	if branding.CustomDomain != "" {
		existing.CustomDomain = branding.CustomDomain
	}
	if branding.CustomCSS != "" {
		existing.CustomCSS = branding.CustomCSS
	}
	if branding.FaviconURL != "" {
		existing.FaviconURL = branding.FaviconURL
	}
	if branding.EmailTemplate != "" {
		existing.EmailTemplate = branding.EmailTemplate
	}
	if branding.FooterText != "" {
		existing.FooterText = branding.FooterText
	}
	if branding.SupportEmail != "" {
		existing.SupportEmail = branding.SupportEmail
	}
	if branding.SupportPhone != "" {
		existing.SupportPhone = branding.SupportPhone
	}

	existing.UpdatedAt = time.Now()

	err = s.repo.UpdateBranding(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// GetBranding retrieves branding for a restaurant
func (s *WhiteLabelService) GetBranding(ctx context.Context, restaurantID int64) (*models.BrandingConfig, error) {
	return s.repo.GetBranding(ctx, restaurantID)
}

// GetThemePresets retrieves available theme presets
func (s *WhiteLabelService) GetThemePresets(ctx context.Context) ([]*models.ThemePreset, error) {
	return s.repo.GetThemePresets(ctx)
}

// ApplyThemePreset applies a theme preset to a restaurant
func (s *WhiteLabelService) ApplyThemePreset(ctx context.Context, restaurantID int64, presetID int64) (*models.BrandingConfig, error) {
	preset, err := s.repo.GetThemePreset(ctx, presetID)
	if err != nil {
		return nil, err
	}

	if preset == nil {
		return nil, fmt.Errorf("theme preset not found")
	}

	branding := &models.BrandingConfig{
		PrimaryColor:    preset.PrimaryColor,
		SecondaryColor:  preset.SecondaryColor,
		AccentColor:     preset.AccentColor,
		FontFamily:      preset.FontFamily,
		BackgroundColor: preset.BackgroundColor,
	}

	return s.UpdateBranding(ctx, restaurantID, branding)
}

// LicenseService handles license management
type LicenseService struct {
	repo repository.WhiteLabelRepository
}

// NewLicenseService creates a new license service
func NewLicenseService(repo repository.WhiteLabelRepository) *LicenseService {
	return &LicenseService{
		repo: repo,
	}
}

// GenerateLicenseKey generates a new license key
func (s *LicenseService) GenerateLicenseKey(ctx context.Context, restaurantID int64, licenseType string) (*models.LicenseKey, error) {
	// Validate license type
	validTypes := map[string]bool{
		"basic":        true,
		"professional": true,
		"enterprise":   true,
	}

	if !validTypes[licenseType] {
		return nil, fmt.Errorf("invalid license type: %s", licenseType)
	}

	// Generate unique license key
	licenseKey := generateSecureKey(32)
	activationCode := generateSecureKey(16)

	// Set features based on license type
	features := []string{}
	maxUsers := 5
	maxInstances := 1
	maxCustomers := 100

	switch licenseType {
	case "professional":
		features = []string{"custom_branding", "advanced_analytics", "api_access", "sso"}
		maxUsers = 20
		maxInstances = 3
		maxCustomers = 1000
	case "enterprise":
		features = []string{"custom_branding", "advanced_analytics", "api_access", "sso", "dedicated_support", "white_label_marketplace", "custom_integrations"}
		maxUsers = 100
		maxInstances = 10
		maxCustomers = 10000
	default:
		features = []string{"custom_branding"}
		maxUsers = 5
		maxInstances = 1
		maxCustomers = 100
	}

	license := &models.LicenseKey{
		RestaurantID:   restaurantID,
		LicenseKey:     licenseKey,
		LicenseType:    licenseType,
		MaxInstances:   maxInstances,
		MaxUsers:       maxUsers,
		MaxCustomers:   maxCustomers,
		Features:       features,
		IssuedAt:       time.Now(),
		ExpiresAt:      time.Now().AddDate(1, 0, 0), // 1 year
		Status:         "valid",
		ActivationCode: &activationCode,
	}

	return s.repo.CreateLicenseKey(ctx, license)
}

// ValidateLicense validates a license key
func (s *LicenseService) ValidateLicense(ctx context.Context, licenseKey string) (*models.LicenseKey, error) {
	license, err := s.repo.GetLicenseKeyByKey(ctx, licenseKey)
	if err != nil {
		return nil, err
	}

	if license == nil {
		return nil, fmt.Errorf("license key not found")
	}

	if license.Status != "valid" {
		return nil, fmt.Errorf("license status is %s", license.Status)
	}

	if time.Now().After(license.ExpiresAt) {
		license.Status = "expired"
		s.repo.UpdateLicenseKey(ctx, license)
		return nil, fmt.Errorf("license has expired")
	}

	// Update last validated time
	now := time.Now()
	license.LastValidatedAt = &now
	s.repo.UpdateLicenseKey(ctx, license)

	return license, nil
}

// CheckLicenseExpiration checks for expiring licenses
func (s *LicenseService) CheckLicenseExpiration(ctx context.Context) ([]*models.LicenseKey, error) {
	expiringLicenses, err := s.repo.GetExpiringSoonLicenses(ctx, 30)
	if err != nil {
		log.Printf("Error checking expiring licenses: %v", err)
		return nil, err
	}

	return expiringLicenses, nil
}

// RenewLicense renews an expiring license
func (s *LicenseService) RenewLicense(ctx context.Context, licenseID int64) (*models.LicenseKey, error) {
	license, err := s.repo.GetLicenseKey(ctx, licenseID)
	if err != nil {
		return nil, err
	}

	if license == nil {
		return nil, fmt.Errorf("license not found")
	}

	// Extend expiration by 1 year
	license.ExpiresAt = license.ExpiresAt.AddDate(1, 0, 0)
	license.Status = "valid"
	license.UpdatedAt = time.Now()

	err = s.repo.UpdateLicenseKey(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

// ActivateLicense activates a license with activation code
func (s *LicenseService) ActivateLicense(ctx context.Context, licenseKey, activationCode string) (*models.LicenseKey, error) {
	license, err := s.repo.GetLicenseKeyByKey(ctx, licenseKey)
	if err != nil {
		return nil, err
	}

	if license == nil {
		return nil, fmt.Errorf("license key not found")
	}

	if license.ActivationCode == nil || *license.ActivationCode != activationCode {
		return nil, fmt.Errorf("invalid activation code")
	}

	now := time.Now()
	license.ActivatedAt = &now
	license.ActivationCode = nil // Clear activation code after use
	license.Status = "valid"
	license.UpdatedAt = now

	err = s.repo.UpdateLicenseKey(ctx, license)
	if err != nil {
		return nil, err
	}

	return license, nil
}

// ResellerService handles reseller partner management
type ResellerService struct {
	repo repository.WhiteLabelRepository
}

// NewResellerService creates a new reseller service
func NewResellerService(repo repository.WhiteLabelRepository) *ResellerService {
	return &ResellerService{
		repo: repo,
	}
}

// CreateResellerPartner creates a new reseller partner
func (s *ResellerService) CreateResellerPartner(ctx context.Context, partner *models.ResellerPartner) (*models.ResellerPartner, error) {
	if partner.Name == "" || partner.Email == "" {
		return nil, fmt.Errorf("name and email are required")
	}

	// Check if email already exists
	existing, _ := s.repo.GetResellerPartnerByEmail(ctx, partner.Email)
	if existing != nil {
		return nil, fmt.Errorf("reseller with this email already exists")
	}

	partner.Status = "pending"
	partner.CreatedAt = time.Now()
	partner.UpdatedAt = time.Now()

	return s.repo.CreateResellerPartner(ctx, partner)
}

// ApproveResellerPartner approves a reseller partner
func (s *ResellerService) ApproveResellerPartner(ctx context.Context, resellerID int64) (*models.ResellerPartner, error) {
	partner, err := s.repo.GetResellerPartner(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, fmt.Errorf("reseller not found")
	}

	now := time.Now()
	partner.Status = "approved"
	partner.ApprovedAt = &now
	partner.UpdatedAt = now

	err = s.repo.UpdateResellerPartner(ctx, partner)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

// SuspendResellerPartner suspends a reseller partner
func (s *ResellerService) SuspendResellerPartner(ctx context.Context, resellerID int64, reason string) (*models.ResellerPartner, error) {
	partner, err := s.repo.GetResellerPartner(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, fmt.Errorf("reseller not found")
	}

	now := time.Now()
	partner.Status = "suspended"
	partner.SuspendedAt = &now
	partner.UpdatedAt = now

	err = s.repo.UpdateResellerPartner(ctx, partner)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

// UpdateResellerTierLevel updates a reseller's tier based on performance
func (s *ResellerService) UpdateResellerTierLevel(ctx context.Context, resellerID int64) (*models.ResellerPartner, error) {
	partner, err := s.repo.GetResellerPartner(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	if partner == nil {
		return nil, fmt.Errorf("reseller not found")
	}

	// Determine tier based on monthly revenue and customer count
	// Bronze: < $5k/month
	// Silver: $5k-$15k/month
	// Gold: $15k-$50k/month
	// Platinum: > $50k/month

	if partner.MonthlyRevenue > 50000 {
		partner.TierLevel = "platinum"
	} else if partner.MonthlyRevenue > 15000 {
		partner.TierLevel = "gold"
	} else if partner.MonthlyRevenue > 5000 {
		partner.TierLevel = "silver"
	} else {
		partner.TierLevel = "bronze"
	}

	// Update commission rate based on tier
	switch partner.TierLevel {
	case "platinum":
		partner.CommissionRate = 0.20 // 20%
	case "gold":
		partner.CommissionRate = 0.15 // 15%
	case "silver":
		partner.CommissionRate = 0.12 // 12%
	default:
		partner.CommissionRate = 0.10 // 10%
	}

	partner.UpdatedAt = time.Now()
	err = s.repo.UpdateResellerPartner(ctx, partner)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

// CreateResellerPortalUser creates a user account for reseller portal
func (s *ResellerService) CreateResellerPortalUser(ctx context.Context, resellerID int64, user *models.ResellerPortalUser) (*models.ResellerPortalUser, error) {
	if user.Email == "" || user.FirstName == "" {
		return nil, fmt.Errorf("email and first name are required")
	}

	// Check if email already exists
	existing, _ := s.repo.GetResellerPortalUserByEmail(ctx, user.Email)
	if existing != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	user.ResellerID = resellerID
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Hash password (in real implementation, use proper bcrypt)
	if user.PasswordHash != "" {
		user.PasswordHash = hashPassword(user.PasswordHash)
	}

	return s.repo.CreateResellerPortalUser(ctx, user)
}

// GetResellerDashboardMetrics gets metrics for reseller dashboard
func (s *ResellerService) GetResellerDashboardMetrics(ctx context.Context, resellerID int64, period string) ([]*models.ResellerDashboardMetrics, error) {
	return s.repo.GetResellerMetrics(ctx, resellerID, period, 30)
}

// CalculateResellerCommission calculates pending commission for a reseller
func (s *ResellerService) CalculateResellerCommission(ctx context.Context, resellerID int64) (float64, error) {
	pending, err := s.repo.GetPendingCommissions(ctx, resellerID)
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, commission := range pending {
		total += commission.Amount
	}

	return total, nil
}

// MarketplaceService handles white-label marketplace operations
type MarketplaceService struct {
	repo repository.WhiteLabelRepository
}

// NewMarketplaceService creates a new marketplace service
func NewMarketplaceService(repo repository.WhiteLabelRepository) *MarketplaceService {
	return &MarketplaceService{
		repo: repo,
	}
}

// CreateMarketplace creates a new white-label marketplace
func (s *MarketplaceService) CreateMarketplace(ctx context.Context, restaurantID int64, marketplace *models.Marketplace) (*models.Marketplace, error) {
	if marketplace.Name == "" || marketplace.Slug == "" {
		return nil, fmt.Errorf("name and slug are required")
	}

	// Validate slug format (alphanumeric and hyphens only)
	if !isValidSlug(marketplace.Slug) {
		return nil, fmt.Errorf("slug must contain only alphanumeric characters and hyphens")
	}

	marketplace.RestaurantID = restaurantID
	marketplace.CreatedAt = time.Now()
	marketplace.UpdatedAt = time.Now()

	return s.repo.CreateMarketplace(ctx, marketplace)
}

// ListMarketplaceProducts lists products in a marketplace
func (s *MarketplaceService) ListMarketplaceProducts(ctx context.Context, marketplaceID int64, filters map[string]interface{}, page, pageSize int) ([]*models.MarketplaceProduct, error) {
	status := "published"
	if statusVal, ok := filters["status"]; ok {
		status = statusVal.(string)
	}

	offset := (page - 1) * pageSize
	return s.repo.GetMarketplaceProducts(ctx, marketplaceID, status, pageSize, offset)
}

// SearchMarketplaceProducts searches products in a marketplace
func (s *MarketplaceService) SearchMarketplaceProducts(ctx context.Context, marketplaceID int64, query string, category string) ([]*models.MarketplaceProduct, error) {
	return s.repo.SearchMarketplaceProducts(ctx, marketplaceID, query, category)
}

// ApproveMarketplaceProduct approves a product for sale
func (s *MarketplaceService) ApproveMarketplaceProduct(ctx context.Context, productID int64) (*models.MarketplaceProduct, error) {
	product, err := s.repo.GetMarketplaceProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	now := time.Now()
	product.IsApproved = true
	product.ApprovedAt = &now
	product.Status = "published"
	product.UpdatedAt = now

	err = s.repo.UpdateMarketplaceProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// RejectMarketplaceProduct rejects a product
func (s *MarketplaceService) RejectMarketplaceProduct(ctx context.Context, productID int64, reason string) (*models.MarketplaceProduct, error) {
	product, err := s.repo.GetMarketplaceProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	product.IsApproved = false
	product.RejectionReason = &reason
	product.Status = "inactive"
	product.UpdatedAt = time.Now()

	err = s.repo.UpdateMarketplaceProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProductStock updates product stock level
func (s *MarketplaceService) UpdateProductStock(ctx context.Context, productID int64, quantity int64) error {
	return s.repo.UpdateProductStock(ctx, productID, quantity)
}

// TrackProductView records a product view
func (s *MarketplaceService) TrackProductView(ctx context.Context, productID int64) error {
	return s.repo.UpdateProductViewCount(ctx, productID)
}

// Helper functions

func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	_, err := hex.DecodeString(color[1:])
	return err == nil
}

func isValidSlug(slug string) bool {
	for _, ch := range slug {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	return true
}

func generateSecureKey(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
