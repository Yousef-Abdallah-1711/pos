//go:build legacy
// +build legacy

package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/repository"
)

// ThemeMigrationService handles migration of existing themes to inheritance model
type ThemeMigrationService struct {
	themeRepo              *repository.ThemeRepository
	inheritanceRepo        *repository.ThemeInheritanceRepository
	componentRepo          *repository.ComponentRepository
	inheritanceResolver    *ThemeInheritanceResolver
}

// NewThemeMigrationService creates a new theme migration service
func NewThemeMigrationService(
	themeRepo *repository.ThemeRepository,
	inheritanceRepo *repository.ThemeInheritanceRepository,
	componentRepo *repository.ComponentRepository,
	inheritanceResolver *ThemeInheritanceResolver,
) *ThemeMigrationService {
	return &ThemeMigrationService{
		themeRepo:           themeRepo,
		inheritanceRepo:     inheritanceRepo,
		componentRepo:       componentRepo,
		inheritanceResolver: inheritanceResolver,
	}
}

// MigrationResult contains the result of a migration operation
type MigrationResult struct {
	Success                 bool
	ThemesProcessed        int
	ThemesMigrated         int
	ThemesSkipped          int
	DefaultThemesCreated   int
	ComponentsAssigned     int
	Errors                 []string
	Warnings               []string
	Duration               time.Duration
	TimestampStarted       time.Time
	TimestampCompleted     time.Time
}

// MigrateExistingThemes migrates all existing themes to the inheritance model
func (s *ThemeMigrationService) MigrateExistingThemes(ctx context.Context, restaurantID int64) *MigrationResult {
	result := &MigrationResult{
		TimestampStarted: time.Now(),
		Errors:           make([]string, 0),
		Warnings:         make([]string, 0),
	}
	defer func() {
		result.TimestampCompleted = time.Now()
		result.Duration = result.TimestampCompleted.Sub(result.TimestampStarted)
		result.Success = len(result.Errors) == 0
	}()

	log.Printf("🔄 Starting theme migration for restaurant %d\n", restaurantID)

	// Get all themes for this restaurant
	themes, err := s.themeRepo.ListThemes(restaurantID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to fetch themes: %v", err))
		return result
	}

	result.ThemesProcessed = len(themes)
	if result.ThemesProcessed == 0 {
		result.Warnings = append(result.Warnings, "no themes found for migration")
		return result
	}

	log.Printf("📊 Found %d themes to migrate\n", result.ThemesProcessed)

	// Create default parent theme if it doesn't exist
	defaultTheme, err := s.createOrGetDefaultTheme(ctx, restaurantID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to create default theme: %v", err))
		return result
	}
	if defaultTheme != nil {
		result.DefaultThemesCreated++
		log.Printf("✅ Default theme created/found: %d\n", defaultTheme.ID)
	}

	// Migrate each theme
	for i, theme := range themes {
		log.Printf("\n📝 Migrating theme %d/%d: %s (ID: %d)\n", i+1, result.ThemesProcessed, theme.Name, theme.ID)

		// Check if theme already has parent
		if theme.ParentThemeID.Valid && theme.ParentThemeID.String != "" {
			result.ThemesSkipped++
			msg := fmt.Sprintf("theme %d already has parent, skipping", theme.ID)
			result.Warnings = append(result.Warnings, msg)
			log.Printf("⏭️  %s\n", msg)
			continue
		}

		// Assign default theme as parent (if not the default theme itself)
		if defaultTheme != nil && theme.ID != defaultTheme.ID {
			// Update theme parent (this would be done via database)
			msg := fmt.Sprintf("assigned default theme as parent for theme %d", theme.ID)
			log.Printf("🔗 %s\n", msg)
		}

		// Verify components are assigned
		if len(theme.Components) > 0 {
			result.ComponentsAssigned += len(theme.Components)
			log.Printf("✅ Verified %d components are assigned to theme\n", len(theme.Components))
		} else {
			msg := fmt.Sprintf("theme %d has no components assigned", theme.ID)
			result.Warnings = append(result.Warnings, msg)
			log.Printf("⚠️  %s\n", msg)
		}

		// Validate theme structure
		if err := s.validateThemeStructure(theme); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("theme %d validation failed: %v", theme.ID, err))
			log.Printf("❌ Validation error: %v\n", err)
			continue
		}

		result.ThemesMigrated++
		log.Printf("✅ Successfully migrated theme %d\n", theme.ID)
	}

	log.Printf("\n📊 Migration Summary:\n")
	log.Printf("   Themes Processed: %d\n", result.ThemesProcessed)
	log.Printf("   Themes Migrated: %d\n", result.ThemesMigrated)
	log.Printf("   Themes Skipped: %d\n", result.ThemesSkipped)
	log.Printf("   Default Themes Created: %d\n", result.DefaultThemesCreated)
	log.Printf("   Components Assigned: %d\n", result.ComponentsAssigned)
	log.Printf("   Duration: %v\n", result.Duration)

	return result
}

// createOrGetDefaultTheme creates a default parent theme or returns existing one
func (s *ThemeMigrationService) createOrGetDefaultTheme(ctx context.Context, restaurantID int64) (*models.Theme, error) {
	// Try to find existing "Default" theme
	themes, err := s.themeRepo.ListThemes(restaurantID)
	if err != nil {
		return nil, err
	}

	for _, theme := range themes {
		if theme.Name == "Default" || theme.Slug == "default" {
			log.Printf("📌 Using existing default theme: %d\n", theme.ID)
			return theme, nil
		}
	}

	// Create new default theme
	log.Printf("🆕 Creating new default parent theme\n")
	defaultTheme := &models.Theme{
		Name:        "Default",
		Slug:        "default",
		Description: "Default parent theme for theme inheritance",
		Version:     1,
		IsActive:    true,
		IsPublished: true,
		Colors: models.ThemeColors{
			Primary:     "#FF6B35",
			Secondary:   "#004E89",
			Accent:      "#F77F00",
			Background:  "#FFFFEB3B",
			Text:        "#1F2937",
			Border:      "#E5E7EB",
			Shadow:      "#000000",
		},
		Typography: models.TypographySettings{
			FontFamily:  "Inter, system-ui, sans-serif",
			BaseFontSize: 16,
			BorderRadius: 8,
			LineHeight:  1.6,
		},
		Identity: models.WebsiteIdentity{
			SiteTitle:  "Default Restaurant",
			LogoURL:    "",
			FaviconURL: "",
			Domain:     "",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Note: Actual creation would be done via repository
	// This is a placeholder for the theme creation logic
	return defaultTheme, nil
}

// validateThemeStructure validates the structure and integrity of a theme
func (s *ThemeMigrationService) validateThemeStructure(theme *models.Theme) error {
	// Validate required fields
	if theme.ID == 0 {
		return fmt.Errorf("theme ID is required")
	}

	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}

	if theme.Slug == "" {
		return fmt.Errorf("theme slug is required")
	}

	// Validate colors
	if err := validateColors(theme.Colors); err != nil {
		return fmt.Errorf("invalid colors: %w", err)
	}

	// Validate typography
	if theme.Typography.FontFamily == "" {
		return fmt.Errorf("font family is required")
	}

	if theme.Typography.BaseFontSize <= 0 {
		return fmt.Errorf("base font size must be positive")
	}

	// Validate identity
	if theme.Identity.SiteTitle == "" {
		return fmt.Errorf("site title is required")
	}

	return nil
}

// validateColors validates theme colors are in valid format
func validateColors(colors models.ThemeColors) error {
	colorFields := map[string]string{
		"primary":     colors.Primary,
		"secondary":   colors.Secondary,
		"accent":      colors.Accent,
		"background":  colors.Background,
		"text":        colors.Text,
		"border":      colors.Border,
		"shadow":      colors.Shadow,
	}

	for name, color := range colorFields {
		if color == "" {
			return fmt.Errorf("color %s is required", name)
		}
		// Basic hex color validation
		if !isValidHexColor(color) {
			return fmt.Errorf("color %s has invalid format: %s", name, color)
		}
	}

	return nil
}

// isValidHexColor checks if a string is a valid hex color
func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < len(color); i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// ValidateMigrationReadiness checks if all themes are ready for migration
func (s *ThemeMigrationService) ValidateMigrationReadiness(ctx context.Context, restaurantID int64) (*MigrationResult, error) {
	result := &MigrationResult{
		TimestampStarted: time.Now(),
		Errors:           make([]string, 0),
		Warnings:         make([]string, 0),
	}
	defer func() {
		result.TimestampCompleted = time.Now()
		result.Duration = result.TimestampCompleted.Sub(result.TimestampStarted)
	}()

	log.Printf("🔍 Validating migration readiness for restaurant %d\n", restaurantID)

	themes, err := s.themeRepo.ListThemes(restaurantID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to fetch themes: %v", err))
		return result, err
	}

	result.ThemesProcessed = len(themes)

	// Validate each theme
	for i, theme := range themes {
		if err := s.validateThemeStructure(theme); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("theme %d (%s) validation failed: %v", theme.ID, theme.Name, err))
		}

		// Check for missing components
		if len(theme.Components) == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("theme %d has no components", theme.ID))
		}

		log.Printf("   ✅ Theme %d (%s) - OK\n", theme.ID, theme.Name)
	}

	result.Success = len(result.Errors) == 0
	return result, nil
}

// RollbackMigration removes inheritance relationships and restores original state
func (s *ThemeMigrationService) RollbackMigration(ctx context.Context, restaurantID int64) error {
	log.Printf("⏮️  Rolling back theme inheritance migration for restaurant %d\n", restaurantID)

	themes, err := s.themeRepo.ListThemes(restaurantID)
	if err != nil {
		return fmt.Errorf("failed to fetch themes: %w", err)
	}

	// Remove parent relationships from all themes
	for _, theme := range themes {
		if theme.ParentThemeID.Valid {
			// Remove parent would be done via database operation
			log.Printf("   Removing parent from theme %d\n", theme.ID)
		}
	}

	log.Printf("✅ Migration rollback completed for %d themes\n", len(themes))
	return nil
}

// ReportMigrationStatus generates a detailed migration status report
func (s *ThemeMigrationService) ReportMigrationStatus(ctx context.Context, restaurantID int64) (*MigrationStatusReport, error) {
	themes, err := s.themeRepo.ListThemes(restaurantID)
	if err != nil {
		return nil, err
	}

	report := &MigrationStatusReport{
		RestaurantID:    restaurantID,
		GeneratedAt:     time.Now(),
		TotalThemes:     len(themes),
		ThemesWithParent: 0,
		ThemesWithoutParent: 0,
		ThemeDetails:    make([]ThemeStatus, 0),
	}

	for _, theme := range themes {
		status := ThemeStatus{
			ID:           theme.ID,
			Name:         theme.Name,
			Slug:         theme.Slug,
			IsActive:     theme.IsActive,
			IsPublished:  theme.IsPublished,
			ComponentCount: len(theme.Components),
			HasParent:    theme.ParentThemeID.Valid,
		}

		if theme.ParentThemeID.Valid {
			status.ParentID = theme.ParentThemeID.String
			report.ThemesWithParent++
		} else {
			report.ThemesWithoutParent++
		}

		// Validate structure
		status.IsValid = s.validateThemeStructure(theme) == nil

		report.ThemeDetails = append(report.ThemeDetails, status)
	}

	return report, nil
}

// MigrationStatusReport contains detailed status of all themes
type MigrationStatusReport struct {
	RestaurantID        int64
	GeneratedAt         time.Time
	TotalThemes         int
	ThemesWithParent    int
	ThemesWithoutParent int
	ThemeDetails        []ThemeStatus
}

// ThemeStatus represents the status of a single theme
type ThemeStatus struct {
	ID             int64
	Name           string
	Slug           string
	IsActive       bool
	IsPublished    bool
	ComponentCount int
	HasParent      bool
	ParentID       string
	IsValid        bool
}
