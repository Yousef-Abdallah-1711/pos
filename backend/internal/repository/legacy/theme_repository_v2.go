//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ThemeRepositoryV2 defines the interface for theme data operations (using new models package)
type ThemeRepositoryV2 interface {
	Create(theme *models.Theme) (*models.Theme, error)
	GetByID(id int64) (*models.Theme, error)
	GetBySlug(restaurantID int64, slug string) (*models.Theme, error)
	GetActive(restaurantID int64) (*models.Theme, error)
	GetActiveByRestaurantSlug(restaurantSlug string) (*models.Theme, error)
	List(restaurantID int64) ([]*models.Theme, error)
	Update(theme *models.Theme) (*models.Theme, error)
	Delete(id int64) error
	Activate(restaurantID, themeID int64) error
	Publish(themeID int64) error
	GetByRestaurant(restaurantID int64, filters map[string]interface{}) ([]*models.Theme, error)
	CheckSlugExists(restaurantID int64, slug string, excludeID *int64) (bool, error)
}

// ThemeRepositoryV2Impl implements ThemeRepositoryV2 using PostgreSQL
type ThemeRepositoryV2Impl struct {
	db *sql.DB
}

// NewThemeRepositoryV2 creates a new theme repository for the new models
func NewThemeRepositoryV2(db *sql.DB) ThemeRepositoryV2 {
	return &ThemeRepositoryV2Impl{db: db}
}

// Create inserts a new theme into the database
func (r *ThemeRepositoryV2Impl) Create(theme *models.Theme) (*models.Theme, error) {
	if theme == nil {
		return nil, errors.New("theme cannot be nil")
	}

	query := `
		INSERT INTO themes_v2 (
			restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23,
			$24, $25
		)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	var createdID int64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(
		query,
		theme.RestaurantID, theme.TenantID, theme.Name, theme.Slug, theme.Description,
		theme.Version, theme.IsActive, theme.IsPublished,
		theme.Colors.Primary, theme.Colors.Secondary, theme.Colors.Accent,
		theme.Colors.Background, theme.Colors.Text, theme.Colors.Border, theme.Colors.Shadow,
		theme.Typography.FontFamily, theme.Typography.BaseFontSize, theme.Typography.BorderRadius, theme.Typography.LineHeight,
		theme.Identity.SiteTitle, theme.Identity.LogoURL, theme.Identity.FaviconURL, theme.Identity.Domain,
		now, now,
	).Scan(&createdID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create theme: %w", err)
	}

	// Set the returned ID and timestamps
	theme.ID = createdID
	theme.CreatedAt = createdAt
	theme.UpdatedAt = updatedAt

	return theme, nil
}

// GetByID retrieves a theme by its ID
func (r *ThemeRepositoryV2Impl) GetByID(id int64) (*models.Theme, error) {
	query := `
		SELECT
			id, restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			header_config, footer_config,
			created_at, updated_at, published_at
		FROM themes_v2
		WHERE id = $1
	`

	theme := &models.Theme{}
	var publishedAt sql.NullTime
	var siteTitle, logoURL, faviconURL, domain sql.NullString
	var headerConfigJSON, footerConfigJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&theme.ID, &theme.RestaurantID, &theme.TenantID,
		&theme.Name, &theme.Slug, &theme.Description,
		&theme.Version, &theme.IsActive, &theme.IsPublished,
		&theme.Colors.Primary, &theme.Colors.Secondary, &theme.Colors.Accent,
		&theme.Colors.Background, &theme.Colors.Text, &theme.Colors.Border, &theme.Colors.Shadow,
		&theme.Typography.FontFamily, &theme.Typography.BaseFontSize, &theme.Typography.BorderRadius, &theme.Typography.LineHeight,
		&siteTitle, &logoURL, &faviconURL, &domain,
		&headerConfigJSON, &footerConfigJSON,
		&theme.CreatedAt, &theme.UpdatedAt, &publishedAt,
	)

	// Handle NULL values - convert to empty strings
	if siteTitle.Valid {
		theme.Identity.SiteTitle = siteTitle.String
	}
	if logoURL.Valid {
		theme.Identity.LogoURL = logoURL.String
	}
	if faviconURL.Valid {
		theme.Identity.FaviconURL = faviconURL.String
	}
	if domain.Valid {
		theme.Identity.Domain = domain.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch theme by id: %w", err)
	}

	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	// Unmarshal header and footer configs
	if len(headerConfigJSON) > 0 && string(headerConfigJSON) != "{}" {
		var header models.HeaderConfig
		if err := json.Unmarshal(headerConfigJSON, &header); err == nil {
			theme.Header = &header
		}
	}

	if len(footerConfigJSON) > 0 && string(footerConfigJSON) != "{}" {
		var footer models.FooterConfig
		if err := json.Unmarshal(footerConfigJSON, &footer); err == nil {
			theme.Footer = &footer
		}
	}

	// Load components for this theme
	components, err := r.loadThemeComponents(theme.ID)
	if err == nil && len(components) > 0 {
		theme.Components = components
	}

	return theme, nil
}

// GetBySlug retrieves a theme by restaurant ID and slug
func (r *ThemeRepositoryV2Impl) GetBySlug(restaurantID int64, slug string) (*models.Theme, error) {
	query := `
		SELECT
			id, restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			header_config, footer_config,
			created_at, updated_at, published_at
		FROM themes_v2
		WHERE restaurant_id = $1 AND slug = $2
	`

	theme := &models.Theme{}
	var publishedAt sql.NullTime
	var siteTitle, logoURL, faviconURL, domain sql.NullString
	var headerConfigJSON, footerConfigJSON []byte

	err := r.db.QueryRow(query, restaurantID, slug).Scan(
		&theme.ID, &theme.RestaurantID, &theme.TenantID,
		&theme.Name, &theme.Slug, &theme.Description,
		&theme.Version, &theme.IsActive, &theme.IsPublished,
		&theme.Colors.Primary, &theme.Colors.Secondary, &theme.Colors.Accent,
		&theme.Colors.Background, &theme.Colors.Text, &theme.Colors.Border, &theme.Colors.Shadow,
		&theme.Typography.FontFamily, &theme.Typography.BaseFontSize, &theme.Typography.BorderRadius, &theme.Typography.LineHeight,
		&siteTitle, &logoURL, &faviconURL, &domain,
		&headerConfigJSON, &footerConfigJSON,
		&theme.CreatedAt, &theme.UpdatedAt, &publishedAt,
	)

	// Handle NULL values - convert to empty strings
	if siteTitle.Valid {
		theme.Identity.SiteTitle = siteTitle.String
	}
	if logoURL.Valid {
		theme.Identity.LogoURL = logoURL.String
	}
	if faviconURL.Valid {
		theme.Identity.FaviconURL = faviconURL.String
	}
	if domain.Valid {
		theme.Identity.Domain = domain.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch theme by slug: %w", err)
	}

	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	// Unmarshal header and footer configs
	if len(headerConfigJSON) > 0 && string(headerConfigJSON) != "{}" {
		var header models.HeaderConfig
		if err := json.Unmarshal(headerConfigJSON, &header); err == nil {
			theme.Header = &header
		}
	}

	if len(footerConfigJSON) > 0 && string(footerConfigJSON) != "{}" {
		var footer models.FooterConfig
		if err := json.Unmarshal(footerConfigJSON, &footer); err == nil {
			theme.Footer = &footer
		}
	}

	// Load components for this theme
	components, err := r.loadThemeComponents(theme.ID)
	if err == nil && len(components) > 0 {
		theme.Components = components
	}

	return theme, nil
}

// GetActive retrieves the active theme for a restaurant
func (r *ThemeRepositoryV2Impl) GetActive(restaurantID int64) (*models.Theme, error) {
	query := `
		SELECT
			id, restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			header_config, footer_config,
			created_at, updated_at, published_at
		FROM themes_v2
		WHERE restaurant_id = $1 AND is_active = true
		LIMIT 1
	`

	theme := &models.Theme{}
	var publishedAt sql.NullTime
	var siteTitle, logoURL, faviconURL, domain sql.NullString
	var headerConfigJSON, footerConfigJSON []byte

	err := r.db.QueryRow(query, restaurantID).Scan(
		&theme.ID, &theme.RestaurantID, &theme.TenantID,
		&theme.Name, &theme.Slug, &theme.Description,
		&theme.Version, &theme.IsActive, &theme.IsPublished,
		&theme.Colors.Primary, &theme.Colors.Secondary, &theme.Colors.Accent,
		&theme.Colors.Background, &theme.Colors.Text, &theme.Colors.Border, &theme.Colors.Shadow,
		&theme.Typography.FontFamily, &theme.Typography.BaseFontSize, &theme.Typography.BorderRadius, &theme.Typography.LineHeight,
		&siteTitle, &logoURL, &faviconURL, &domain,
		&headerConfigJSON, &footerConfigJSON,
		&theme.CreatedAt, &theme.UpdatedAt, &publishedAt,
	)

	// Handle NULL values - convert to empty strings
	if siteTitle.Valid {
		theme.Identity.SiteTitle = siteTitle.String
	}
	if logoURL.Valid {
		theme.Identity.LogoURL = logoURL.String
	}
	if faviconURL.Valid {
		theme.Identity.FaviconURL = faviconURL.String
	}
	if domain.Valid {
		theme.Identity.Domain = domain.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch active theme: %w", err)
	}

	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	// Unmarshal header and footer configs
	if len(headerConfigJSON) > 0 && string(headerConfigJSON) != "{}" {
		var header models.HeaderConfig
		if err := json.Unmarshal(headerConfigJSON, &header); err == nil {
			theme.Header = &header
		}
	}

	if len(footerConfigJSON) > 0 && string(footerConfigJSON) != "{}" {
		var footer models.FooterConfig
		if err := json.Unmarshal(footerConfigJSON, &footer); err == nil {
			theme.Footer = &footer
		}
	}

	// Load components for this theme
	components, err := r.loadThemeComponents(theme.ID)
	if err == nil && len(components) > 0 {
		theme.Components = components
	}

	return theme, nil
}

// GetActiveByRestaurantSlug retrieves the active theme for a restaurant by its slug
func (r *ThemeRepositoryV2Impl) GetActiveByRestaurantSlug(restaurantSlug string) (*models.Theme, error) {
	// First, look up the restaurant ID by slug
	var restaurantID int64
	restaurantQuery := `SELECT id FROM restaurants WHERE slug = $1 LIMIT 1`
	fmt.Printf("DEBUG: Looking up restaurant with slug='%s'\n", restaurantSlug)
	err := r.db.QueryRow(restaurantQuery, restaurantSlug).Scan(&restaurantID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("DEBUG: Restaurant not found with slug='%s'\n", restaurantSlug)
			return nil, sql.ErrNoRows
		}
		fmt.Printf("DEBUG: Error looking up restaurant: %v\n", err)
		return nil, fmt.Errorf("failed to find restaurant by slug: %w", err)
	}

	fmt.Printf("DEBUG: Found restaurant id=%d for slug='%s'\n", restaurantID, restaurantSlug)

	// Now use the existing GetActive method with the restaurant ID
	theme, err := r.GetActive(restaurantID)
	if err != nil {
		fmt.Printf("DEBUG: GetActive returned error for restaurant_id=%d: %v\n", restaurantID, err)
	} else {
		fmt.Printf("DEBUG: GetActive succeeded for restaurant_id=%d, theme_id=%d\n", restaurantID, theme.ID)
	}
	return theme, err
}

// List retrieves all themes for a restaurant
func (r *ThemeRepositoryV2Impl) List(restaurantID int64) ([]*models.Theme, error) {
	query := `
		SELECT
			id, restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			header_config, footer_config,
			created_at, updated_at, published_at
		FROM themes_v2
		WHERE restaurant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list themes: %w", err)
	}
	defer rows.Close()

	var themes []*models.Theme
	for rows.Next() {
		theme := &models.Theme{}
		var publishedAt sql.NullTime
		var siteTitle, logoURL, faviconURL, domain sql.NullString
		var headerConfigJSON, footerConfigJSON []byte

		err := rows.Scan(
			&theme.ID, &theme.RestaurantID, &theme.TenantID,
			&theme.Name, &theme.Slug, &theme.Description,
			&theme.Version, &theme.IsActive, &theme.IsPublished,
			&theme.Colors.Primary, &theme.Colors.Secondary, &theme.Colors.Accent,
			&theme.Colors.Background, &theme.Colors.Text, &theme.Colors.Border, &theme.Colors.Shadow,
			&theme.Typography.FontFamily, &theme.Typography.BaseFontSize, &theme.Typography.BorderRadius, &theme.Typography.LineHeight,
			&siteTitle, &logoURL, &faviconURL, &domain,
			&headerConfigJSON, &footerConfigJSON,
			&theme.CreatedAt, &theme.UpdatedAt, &publishedAt,
		)

		// Handle NULL values - convert to empty strings
		if siteTitle.Valid {
			theme.Identity.SiteTitle = siteTitle.String
		}
		if logoURL.Valid {
			theme.Identity.LogoURL = logoURL.String
		}
		if faviconURL.Valid {
			theme.Identity.FaviconURL = faviconURL.String
		}
		if domain.Valid {
			theme.Identity.Domain = domain.String
		}

		if err != nil {
			return nil, fmt.Errorf("failed to scan theme: %w", err)
		}

		if publishedAt.Valid {
			theme.PublishedAt = &publishedAt.Time
		}

		// Unmarshal header and footer configs
		if len(headerConfigJSON) > 0 && string(headerConfigJSON) != "{}" {
			var header models.HeaderConfig
			if err := json.Unmarshal(headerConfigJSON, &header); err == nil {
				theme.Header = &header
			}
		}

		if len(footerConfigJSON) > 0 && string(footerConfigJSON) != "{}" {
			var footer models.FooterConfig
			if err := json.Unmarshal(footerConfigJSON, &footer); err == nil {
				theme.Footer = &footer
			}
		}

		// Load components for this theme
		components, err := r.loadThemeComponents(theme.ID)
		if err == nil && len(components) > 0 {
			theme.Components = components
		}

		themes = append(themes, theme)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating themes: %w", err)
	}

	if themes == nil {
		themes = []*models.Theme{}
	}

	return themes, nil
}

// Update updates an existing theme
func (r *ThemeRepositoryV2Impl) Update(theme *models.Theme) (*models.Theme, error) {
	if theme == nil || theme.ID == 0 {
		return nil, errors.New("theme must have a valid ID")
	}

	// Marshal header and footer configs to JSON
	headerConfigJSON := []byte("{}")
	footerConfigJSON := []byte("{}")

	if theme.Header != nil {
		var err error
		headerConfigJSON, err = json.Marshal(theme.Header)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal header config: %w", err)
		}
	}

	if theme.Footer != nil {
		var err error
		footerConfigJSON, err = json.Marshal(theme.Footer)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal footer config: %w", err)
		}
	}

	query := `
		UPDATE themes_v2
		SET
			name = $1,
			slug = $2,
			description = $3,
			version = $4,
			primary_color = $5,
			secondary_color = $6,
			accent_color = $7,
			background_color = $8,
			text_color = $9,
			border_color = $10,
			shadow_color = $11,
			font_family = $12,
			base_font_size = $13,
			border_radius = $14,
			line_height = $15,
			site_title = $16,
			logo_url = $17,
			favicon_url = $18,
			domain = $19,
			header_config = $20,
			footer_config = $21,
			updated_at = $22
		WHERE id = $23 AND restaurant_id = $24
		RETURNING created_at, updated_at
	`

	now := time.Now()
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(
		query,
		theme.Name, theme.Slug, theme.Description, theme.Version,
		theme.Colors.Primary, theme.Colors.Secondary, theme.Colors.Accent,
		theme.Colors.Background, theme.Colors.Text, theme.Colors.Border, theme.Colors.Shadow,
		theme.Typography.FontFamily, theme.Typography.BaseFontSize, theme.Typography.BorderRadius, theme.Typography.LineHeight,
		theme.Identity.SiteTitle, theme.Identity.LogoURL, theme.Identity.FaviconURL, theme.Identity.Domain,
		headerConfigJSON, footerConfigJSON,
		now, theme.ID, theme.RestaurantID,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update theme: %w", err)
	}

	theme.CreatedAt = createdAt
	theme.UpdatedAt = updatedAt

	return theme, nil
}

// Delete removes a theme from the database
func (r *ThemeRepositoryV2Impl) Delete(id int64) error {
	query := `DELETE FROM themes_v2 WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete theme: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Activate sets a theme as active and deactivates all other themes for the same restaurant
func (r *ThemeRepositoryV2Impl) Activate(restaurantID, themeID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Deactivate all other themes
	deactivateQuery := `UPDATE themes_v2 SET is_active = false WHERE restaurant_id = $1 AND id != $2`
	_, err = tx.Exec(deactivateQuery, restaurantID, themeID)
	if err != nil {
		return fmt.Errorf("failed to deactivate other themes: %w", err)
	}

	// Activate the specified theme
	activateQuery := `UPDATE themes_v2 SET is_active = true, updated_at = $1 WHERE id = $2 AND restaurant_id = $3`
	result, err := tx.Exec(activateQuery, time.Now(), themeID, restaurantID)
	if err != nil {
		return fmt.Errorf("failed to activate theme: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Publish marks a theme as published
func (r *ThemeRepositoryV2Impl) Publish(themeID int64) error {
	query := `UPDATE themes_v2 SET is_published = true, published_at = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, time.Now(), time.Now(), themeID)
	if err != nil {
		return fmt.Errorf("failed to publish theme: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByRestaurant retrieves themes for a restaurant with optional filters
func (r *ThemeRepositoryV2Impl) GetByRestaurant(restaurantID int64, filters map[string]interface{}) ([]*models.Theme, error) {
	query := `
		SELECT
			id, restaurant_id, tenant_id, name, slug, description,
			version, is_active, is_published,
			primary_color, secondary_color, accent_color,
			background_color, text_color, border_color, shadow_color,
			font_family, base_font_size, border_radius, line_height,
			site_title, logo_url, favicon_url, domain,
			created_at, updated_at, published_at
		FROM themes_v2
		WHERE restaurant_id = $1
	`

	args := []interface{}{restaurantID}
	paramCount := 2

	// Apply optional filters
	if isActive, ok := filters["isActive"].(bool); ok {
		query += fmt.Sprintf(` AND is_active = $%d`, paramCount)
		args = append(args, isActive)
		paramCount++
	}

	if isPublished, ok := filters["isPublished"].(bool); ok {
		query += fmt.Sprintf(` AND is_published = $%d`, paramCount)
		args = append(args, isPublished)
		paramCount++
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query themes by restaurant: %w", err)
	}
	defer rows.Close()

	var themes []*models.Theme
	for rows.Next() {
		theme := &models.Theme{}
		var publishedAt sql.NullTime
		var siteTitle, logoURL, faviconURL, domain sql.NullString

		err := rows.Scan(
			&theme.ID, &theme.RestaurantID, &theme.TenantID,
			&theme.Name, &theme.Slug, &theme.Description,
			&theme.Version, &theme.IsActive, &theme.IsPublished,
			&theme.Colors.Primary, &theme.Colors.Secondary, &theme.Colors.Accent,
			&theme.Colors.Background, &theme.Colors.Text, &theme.Colors.Border, &theme.Colors.Shadow,
			&theme.Typography.FontFamily, &theme.Typography.BaseFontSize, &theme.Typography.BorderRadius, &theme.Typography.LineHeight,
			&siteTitle, &logoURL, &faviconURL, &domain,
			&theme.CreatedAt, &theme.UpdatedAt, &publishedAt,
		)

		// Handle NULL values - convert to empty strings
		if siteTitle.Valid {
			theme.Identity.SiteTitle = siteTitle.String
		}
		if logoURL.Valid {
			theme.Identity.LogoURL = logoURL.String
		}
		if faviconURL.Valid {
			theme.Identity.FaviconURL = faviconURL.String
		}
		if domain.Valid {
			theme.Identity.Domain = domain.String
		}

		if err != nil {
			return nil, fmt.Errorf("failed to scan theme: %w", err)
		}

		if publishedAt.Valid {
			theme.PublishedAt = &publishedAt.Time
		}

		// Load components for this theme
		components, err := r.loadThemeComponents(theme.ID)
		if err == nil && len(components) > 0 {
			theme.Components = components
		}

		themes = append(themes, theme)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating themes: %w", err)
	}

	if themes == nil {
		themes = []*models.Theme{}
	}

	return themes, nil
}

// CheckSlugExists checks if a slug already exists for a restaurant
func (r *ThemeRepositoryV2Impl) CheckSlugExists(restaurantID int64, slug string, excludeID *int64) (bool, error) {
	query := `SELECT COUNT(*) FROM themes_v2 WHERE restaurant_id = $1 AND slug = $2`
	args := []interface{}{restaurantID, slug}

	// Exclude a specific theme ID if provided (for update checks)
	if excludeID != nil {
		query += ` AND id != $3`
		args = append(args, *excludeID)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return count > 0, nil
}

// loadThemeComponents is a helper function to load components for a theme
func (r *ThemeRepositoryV2Impl) loadThemeComponents(themeID int64) ([]models.ThemeComponent, error) {
	query := `
		SELECT
			id, theme_id, component_type, title,
			display_order, enabled, settings, created_at, updated_at
		FROM theme_components
		WHERE theme_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query, themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query theme components: %w", err)
	}
	defer rows.Close()

	var components []models.ThemeComponent
	for rows.Next() {
		component := models.ThemeComponent{}
		var settings sql.NullString

		err := rows.Scan(
			&component.ID, &component.ThemeID, &component.ComponentType,
			&component.Title, &component.DisplayOrder, &component.Enabled,
			&settings, &component.CreatedAt, &component.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan component: %w", err)
		}

		if settings.Valid {
			component.Settings = []byte(settings.String)
		}

		components = append(components, component)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating components: %w", err)
	}

	return components, nil
}
