package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

// ThemeRepository defines the interface for theme data operations
type ThemeRepository interface {
	// Create inserts a new theme into the database
	Create(ctx context.Context, theme *domain.Theme) error

	// GetByID retrieves a theme by its ID
	GetByID(ctx context.Context, tenantID, restaurantID, id int64) (*domain.Theme, error)

	// GetBySlug retrieves a theme by tenant ID and slug
	GetBySlug(ctx context.Context, tenantID, restaurantID int64, slug string) (*domain.Theme, error)

	// List retrieves a paginated list of themes for a tenant
	List(ctx context.Context, tenantID, restaurantID int64, offset, limit int) ([]*domain.Theme, int64, error)

	// Update updates an existing theme
	Update(ctx context.Context, theme *domain.Theme) error

	// Delete deletes a theme by ID
	Delete(ctx context.Context, tenantID, restaurantID, id int64) error

	// GetActive retrieves the active theme for a tenant
	GetActive(ctx context.Context, tenantID, restaurantID int64) (*domain.Theme, error)

	// SetActive activates a theme and deactivates others for the same tenant
	SetActive(ctx context.Context, tenantID, restaurantID, id int64) error

	// ListPresets retrieves a paginated list of theme presets
	ListPresets(ctx context.Context, offset, limit int) ([]*domain.Theme, int64, error)

	// GetPreset retrieves a single theme preset by ID
	GetPreset(ctx context.Context, id int64) (*domain.Theme, error)

	// GetByRestaurantAndActive retrieves the active theme for a restaurant
	GetByRestaurantAndActive(ctx context.Context, restaurantID int64) (*domain.Theme, error)
}

// themeRepositoryImpl implements the ThemeRepository interface
type themeRepositoryImpl struct {
	db *sql.DB
}

// NewThemeRepository creates a new instance of the theme repository
func NewThemeRepository(db *sql.DB) ThemeRepository {
	return &themeRepositoryImpl{
		db: db,
	}
}

// Create inserts a new theme into the database
func (r *themeRepositoryImpl) Create(ctx context.Context, theme *domain.Theme) error {
	// Validate theme data
	if err := theme.Validate(); err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
		INSERT INTO themes_v2 (
			tenant_id, restaurant_id, name, slug, config, description,
			is_active, is_published, version, created_by, created_at, updated_at
		) SELECT
			$1::bigint, $2::bigint, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2::bigint AND tenant_id = $1::bigint)
		RETURNING id, created_at, updated_at
	`

	createdBy := sql.NullInt64{}
	if theme.CreatedBy != nil {
		createdBy = sql.NullInt64{Int64: *theme.CreatedBy, Valid: true}
	}

	// Execute the query
	row := r.db.QueryRowContext(ctx, query,
		theme.TenantID,
		theme.RestaurantID,
		theme.Name,
		theme.Slug,
		theme.Config,
		theme.Description,
		theme.IsActive,
		theme.IsPublished,
		theme.Version,
		createdBy,
		time.Now(),
		time.Now(),
	)

	// Scan the returned ID
	err := row.Scan(&theme.ID, &theme.CreatedAt, &theme.UpdatedAt)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"uk_themes_v2_restaurant_slug\"" {
			return &domain.ThemeValidationError{
				Field:   "slug",
				Message: "slug already exists for this restaurant",
			}
		}
		return fmt.Errorf("failed to create theme: %w", err)
	}

	return nil
}

// GetByID retrieves a theme by its ID
func (r *themeRepositoryImpl) GetByID(ctx context.Context, tenantID, restaurantID, id int64) (*domain.Theme, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, name, slug, config, description,
		       is_active, is_published, version, created_by, updated_by,
		       created_at, updated_at, published_at
		FROM themes_v2
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	theme := &domain.Theme{}
	createdBy := sql.NullInt64{}
	updatedBy := sql.NullInt64{}
	publishedAt := sql.NullTime{}

	err := r.db.QueryRowContext(ctx, query, id, tenantID, restaurantID).Scan(
		&theme.ID,
		&theme.TenantID,
		&theme.RestaurantID,
		&theme.Name,
		&theme.Slug,
		&theme.Config,
		&theme.Description,
		&theme.IsActive,
		&theme.IsPublished,
		&theme.Version,
		&createdBy,
		&updatedBy,
		&theme.CreatedAt,
		&theme.UpdatedAt,
		&publishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("theme not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch theme: %w", err)
	}

	// Convert nullable fields
	if createdBy.Valid {
		theme.CreatedBy = &createdBy.Int64
	}
	if updatedBy.Valid {
		theme.UpdatedBy = &updatedBy.Int64
	}
	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	return theme, nil
}

// GetBySlug retrieves a theme by tenant ID and slug
func (r *themeRepositoryImpl) GetBySlug(ctx context.Context, tenantID, restaurantID int64, slug string) (*domain.Theme, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, name, slug, config, description,
		       is_active, is_published, version, created_by, updated_by,
		       created_at, updated_at, published_at
		FROM themes_v2
		WHERE tenant_id = $1 AND restaurant_id = $2 AND slug = $3
		LIMIT 1
	`

	theme := &domain.Theme{}
	createdBy := sql.NullInt64{}
	updatedBy := sql.NullInt64{}
	publishedAt := sql.NullTime{}

	err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID, slug).Scan(
		&theme.ID,
		&theme.TenantID,
		&theme.RestaurantID,
		&theme.Name,
		&theme.Slug,
		&theme.Config,
		&theme.Description,
		&theme.IsActive,
		&theme.IsPublished,
		&theme.Version,
		&createdBy,
		&updatedBy,
		&theme.CreatedAt,
		&theme.UpdatedAt,
		&publishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("theme not found with slug: %s", slug)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch theme by slug: %w", err)
	}

	// Convert nullable fields
	if createdBy.Valid {
		theme.CreatedBy = &createdBy.Int64
	}
	if updatedBy.Valid {
		theme.UpdatedBy = &updatedBy.Int64
	}
	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	return theme, nil
}

// List retrieves a paginated list of themes for a tenant
func (r *themeRepositoryImpl) List(ctx context.Context, tenantID, restaurantID int64, offset, limit int) ([]*domain.Theme, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM themes_v2 WHERE tenant_id = $1 AND restaurant_id = $2`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, tenantID, restaurantID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count themes: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, restaurant_id, name, slug, config, description,
		       is_active, is_published, version, created_by, updated_by,
		       created_at, updated_at, published_at
		FROM themes_v2
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY created_at DESC
		OFFSET $3 LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list themes: %w", err)
	}
	defer rows.Close()

	var themes []*domain.Theme
	for rows.Next() {
		theme := &domain.Theme{}
		createdBy := sql.NullInt64{}
		updatedBy := sql.NullInt64{}
		publishedAt := sql.NullTime{}

		err := rows.Scan(
			&theme.ID,
			&theme.TenantID,
			&theme.RestaurantID,
			&theme.Name,
			&theme.Slug,
			&theme.Config,
			&theme.Description,
			&theme.IsActive,
			&theme.IsPublished,
			&theme.Version,
			&createdBy,
			&updatedBy,
			&theme.CreatedAt,
			&theme.UpdatedAt,
			&publishedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan theme: %w", err)
		}

		// Convert nullable fields
		if createdBy.Valid {
			theme.CreatedBy = &createdBy.Int64
		}
		if updatedBy.Valid {
			theme.UpdatedBy = &updatedBy.Int64
		}
		if publishedAt.Valid {
			theme.PublishedAt = &publishedAt.Time
		}

		themes = append(themes, theme)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating themes: %w", err)
	}

	return themes, total, nil
}

// Update updates an existing theme
func (r *themeRepositoryImpl) Update(ctx context.Context, theme *domain.Theme) error {
	// Validate theme data
	if err := theme.Validate(); err != nil {
		return err
	}

	if theme.ID <= 0 {
		return errors.New("invalid theme ID")
	}

	// Increment version
	theme.Version++
	theme.UpdatedAt = time.Now()

	updatedBy := sql.NullInt64{}
	if theme.UpdatedBy != nil {
		updatedBy = sql.NullInt64{Int64: *theme.UpdatedBy, Valid: true}
	}

	query := `
		UPDATE themes_v2
		SET name = $1, slug = $2, config = $3, description = $4,
		    is_active = $5, is_published = $6, version = $7, updated_by = $8, updated_at = $9
		WHERE id = $10 AND tenant_id = $11 AND restaurant_id = $12
		RETURNING updated_at, version
	`

	row := r.db.QueryRowContext(ctx, query,
		theme.Name,
		theme.Slug,
		theme.Config,
		theme.Description,
		theme.IsActive,
		theme.IsPublished,
		theme.Version,
		updatedBy,
		theme.UpdatedAt,
		theme.ID,
		theme.TenantID,
		theme.RestaurantID,
	)

	err := row.Scan(&theme.UpdatedAt, &theme.Version)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"uk_themes_v2_restaurant_slug\"" {
			return &domain.ThemeValidationError{
				Field:   "slug",
				Message: "slug already exists for this restaurant",
			}
		}
		return fmt.Errorf("failed to update theme: %w", err)
	}

	return nil
}

// Delete deletes a theme by ID
func (r *themeRepositoryImpl) Delete(ctx context.Context, tenantID, restaurantID, id int64) error {
	if id <= 0 {
		return errors.New("invalid theme ID")
	}

	query := `DELETE FROM themes_v2 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`
	result, err := r.db.ExecContext(ctx, query, id, tenantID, restaurantID)
	if err != nil {
		return fmt.Errorf("failed to delete theme: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("theme not found")
	}

	return nil
}

// GetActive retrieves the active theme for a tenant
func (r *themeRepositoryImpl) GetActive(ctx context.Context, tenantID, restaurantID int64) (*domain.Theme, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, name, slug, config, description,
		       is_active, is_published, version, created_by, updated_by,
		       created_at, updated_at, published_at
		FROM themes_v2
		WHERE tenant_id = $1 AND restaurant_id = $2 AND is_active = true
		LIMIT 1
	`

	theme := &domain.Theme{}
	createdBy := sql.NullInt64{}
	updatedBy := sql.NullInt64{}
	publishedAt := sql.NullTime{}

	err := r.db.QueryRowContext(ctx, query, tenantID, restaurantID).Scan(
		&theme.ID,
		&theme.TenantID,
		&theme.RestaurantID,
		&theme.Name,
		&theme.Slug,
		&theme.Config,
		&theme.Description,
		&theme.IsActive,
		&theme.IsPublished,
		&theme.Version,
		&createdBy,
		&updatedBy,
		&theme.CreatedAt,
		&theme.UpdatedAt,
		&publishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no active theme found for tenant")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active theme: %w", err)
	}

	// Convert nullable fields
	if createdBy.Valid {
		theme.CreatedBy = &createdBy.Int64
	}
	if updatedBy.Valid {
		theme.UpdatedBy = &updatedBy.Int64
	}
	if publishedAt.Valid {
		theme.PublishedAt = &publishedAt.Time
	}

	return theme, nil
}

// SetActive activates a theme and deactivates others for the same tenant
func (r *themeRepositoryImpl) SetActive(ctx context.Context, tenantID, restaurantID, id int64) error {
	if id <= 0 {
		return errors.New("invalid theme ID")
	}

	// Start a transaction to ensure atomicity
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, get the tenant_id of the theme to be activated
	var owningTenantID int64
	var themeRestaurantID int64
	query := `SELECT tenant_id, restaurant_id FROM themes_v2 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, id, tenantID, restaurantID).Scan(&owningTenantID, &themeRestaurantID)
	if err == sql.ErrNoRows {
		return errors.New("theme not found")
	}
	if err != nil {
		return fmt.Errorf("failed to fetch theme tenant: %w", err)
	}

	// Deactivate all themes for this tenant
	deactiveQuery := `UPDATE themes_v2 SET is_active = false WHERE tenant_id = $1 AND restaurant_id = $2`
	_, err = tx.ExecContext(ctx, deactiveQuery, owningTenantID, themeRestaurantID)
	if err != nil {
		return fmt.Errorf("failed to deactivate other themes: %w", err)
	}

	// Activate the specified theme
	activateQuery := `UPDATE themes_v2 SET is_active = true, updated_at = $1 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`
	result, err := tx.ExecContext(ctx, activateQuery, time.Now(), id, owningTenantID, themeRestaurantID)
	if err != nil {
		return fmt.Errorf("failed to activate theme: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("theme not found")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ListPresets retrieves a paginated list of theme presets from theme_presets table
func (r *themeRepositoryImpl) ListPresets(ctx context.Context, offset, limit int) ([]*domain.Theme, int64, error) {
	// Get total count of active presets
	countQuery := `SELECT COUNT(*) FROM theme_presets WHERE is_active = true`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count presets: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, name, slug, description, preset_data, category
		FROM theme_presets
		WHERE is_active = true
		ORDER BY sort_order ASC, created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list presets: %w", err)
	}
	defer rows.Close()

	var presets []*domain.Theme
	for rows.Next() {
		preset := &domain.Theme{}
		var category sql.NullString
		var presetData []byte

		err := rows.Scan(
			&preset.ID,
			&preset.Name,
			&preset.Slug,
			&preset.Description,
			&presetData,
			&category,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan preset: %w", err)
		}

		// Set config from preset_data
		if len(presetData) > 0 {
			preset.Config = presetData
		}

		// Set default values for preset
		preset.IsActive = true
		preset.IsPublished = true
		preset.Version = 1
		preset.TenantID = 0 // Presets are global, not tenant-specific

		presets = append(presets, preset)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating presets: %w", err)
	}

	return presets, total, nil
}

// GetPreset retrieves a single theme preset by ID from theme_presets table
func (r *themeRepositoryImpl) GetPreset(ctx context.Context, id int64) (*domain.Theme, error) {
	query := `
		SELECT id, name, slug, description, preset_data, category
		FROM theme_presets
		WHERE id = $1 AND is_active = true
	`

	preset := &domain.Theme{}
	var category sql.NullString
	var presetData []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&preset.ID,
		&preset.Name,
		&preset.Slug,
		&preset.Description,
		&presetData,
		&category,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("preset not found")
		}
		return nil, fmt.Errorf("failed to get preset: %w", err)
	}

	// Set config from preset_data
	if len(presetData) > 0 {
		preset.Config = presetData
	}

	// Set default values for preset
	preset.IsActive = true
	preset.IsPublished = true
	preset.Version = 1
	preset.TenantID = 0 // Presets are global, not tenant-specific

	return preset, nil
}

// GetByRestaurantAndActive retrieves the active theme for a restaurant
func (r *themeRepositoryImpl) GetByRestaurantAndActive(ctx context.Context, restaurantID int64) (*domain.Theme, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, name, slug, config, description,
		       is_active, is_published, version, created_by, created_at, updated_at
		FROM themes_v2
		WHERE restaurant_id = $1 AND is_active = true
		ORDER BY updated_at DESC
		LIMIT 1
	`

	theme := &domain.Theme{}
	var createdBy sql.NullInt64
	var config []byte

	err := r.db.QueryRowContext(ctx, query, restaurantID).Scan(
		&theme.ID,
		&theme.TenantID,
		&theme.RestaurantID,
		&theme.Name,
		&theme.Slug,
		&config,
		&theme.Description,
		&theme.IsActive,
		&theme.IsPublished,
		&theme.Version,
		&createdBy,
		&theme.CreatedAt,
		&theme.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no active theme found for this restaurant")
		}
		return nil, fmt.Errorf("failed to get active theme: %w", err)
	}

	// Parse config from JSON
	if len(config) > 0 {
		theme.Config = config
	}

	// Set created_by if present
	if createdBy.Valid {
		theme.CreatedBy = &createdBy.Int64
	}

	return theme, nil
}
