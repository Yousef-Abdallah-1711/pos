//go:build legacy
// +build legacy

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pos-saas/internal/domain"
)

// ComponentRepository handles all database operations for components
type ComponentRepository struct {
	db *sql.DB
}

// NewComponentRepository creates a new component repository instance
func NewComponentRepository(db *sql.DB) *ComponentRepository {
	return &ComponentRepository{db: db}
}

// CreateComponentRegistry creates a new component registration
func (r *ComponentRepository) CreateComponentRegistry(ctx context.Context, component *models.ComponentRegistry) error {
	if component.ID == "" {
		component.ID = uuid.New().String()
	}
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()

	query := `
		INSERT INTO component_registry
		(id, name, description, component_id, version, aliases, category, tags, schema, thumbnail_url, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		component.ID,
		component.Name,
		component.Description,
		component.ComponentID,
		component.Version,
		component.Aliases,
		component.Category,
		component.Tags,
		component.Schema,
		component.ThumbnailURL,
		component.IsActive,
		component.CreatedAt,
		component.UpdatedAt,
	)

	return err
}

// GetComponentByID retrieves a component by its ID
func (r *ComponentRepository) GetComponentByID(ctx context.Context, id string) (*models.ComponentRegistry, error) {
	component := &models.ComponentRegistry{}

	query := `
		SELECT id, name, description, component_id, version, aliases, category, tags, schema, thumbnail_url, is_active, created_at, updated_at
		FROM component_registry
		WHERE id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&component.ID,
		&component.Name,
		&component.Description,
		&component.ComponentID,
		&component.Version,
		&component.Aliases,
		&component.Category,
		&component.Tags,
		&component.Schema,
		&component.ThumbnailURL,
		&component.IsActive,
		&component.CreatedAt,
		&component.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("component not found: %s", id)
		}
		return nil, err
	}

	return component, nil
}

// GetComponentByComponentID retrieves a component by its component ID (e.g., "hero")
func (r *ComponentRepository) GetComponentByComponentID(ctx context.Context, componentID string) (*models.ComponentRegistry, error) {
	component := &models.ComponentRegistry{}

	query := `
		SELECT id, name, description, component_id, version, aliases, category, tags, schema, thumbnail_url, is_active, created_at, updated_at
		FROM component_registry
		WHERE component_id = ? AND is_active = true
		ORDER BY version DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, componentID).Scan(
		&component.ID,
		&component.Name,
		&component.Description,
		&component.ComponentID,
		&component.Version,
		&component.Aliases,
		&component.Category,
		&component.Tags,
		&component.Schema,
		&component.ThumbnailURL,
		&component.IsActive,
		&component.CreatedAt,
		&component.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("component not found: %s", componentID)
		}
		return nil, err
	}

	return component, nil
}

// ListComponents retrieves all active components with optional filters
func (r *ComponentRepository) ListComponents(ctx context.Context, category string, limit int, offset int) ([]*models.ComponentRegistry, error) {
	query := "SELECT id, name, description, component_id, version, aliases, category, tags, schema, thumbnail_url, is_active, created_at, updated_at FROM component_registry WHERE is_active = true"
	args := []interface{}{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []*models.ComponentRegistry
	for rows.Next() {
		component := &models.ComponentRegistry{}
		err := rows.Scan(
			&component.ID,
			&component.Name,
			&component.Description,
			&component.ComponentID,
			&component.Version,
			&component.Aliases,
			&component.Category,
			&component.Tags,
			&component.Schema,
			&component.ThumbnailURL,
			&component.IsActive,
			&component.CreatedAt,
			&component.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, component)
	}

	return components, rows.Err()
}

// UpdateComponentRegistry updates an existing component
func (r *ComponentRepository) UpdateComponentRegistry(ctx context.Context, component *models.ComponentRegistry) error {
	component.UpdatedAt = time.Now()

	query := `
		UPDATE component_registry
		SET name = ?, description = ?, version = ?, aliases = ?, category = ?, tags = ?, schema = ?, thumbnail_url = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		component.Name,
		component.Description,
		component.Version,
		component.Aliases,
		component.Category,
		component.Tags,
		component.Schema,
		component.ThumbnailURL,
		component.IsActive,
		component.UpdatedAt,
		component.ID,
	)

	return err
}

// DeleteComponentRegistry soft-deletes a component
func (r *ComponentRepository) DeleteComponentRegistry(ctx context.Context, id string) error {
	query := "UPDATE component_registry SET is_active = false, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// CreateComponentThemeMapping links a component to a theme
func (r *ComponentRepository) CreateComponentThemeMapping(ctx context.Context, mapping *models.ComponentThemeMapping) error {
	if mapping.ID == "" {
		mapping.ID = uuid.New().String()
	}
	mapping.CreatedAt = time.Now()
	mapping.UpdatedAt = time.Now()

	query := `
		INSERT INTO component_theme_mapping
		(id, theme_id, component_id, position, is_enabled, config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		mapping.ID,
		mapping.ThemeID,
		mapping.ComponentID,
		mapping.Position,
		mapping.IsEnabled,
		mapping.Config,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	)

	return err
}

// GetComponentsByTheme retrieves all components used in a theme
func (r *ComponentRepository) GetComponentsByTheme(ctx context.Context, themeID string) ([]*models.ComponentThemeMapping, error) {
	query := `
		SELECT id, theme_id, component_id, position, is_enabled, config, created_at, updated_at
		FROM component_theme_mapping
		WHERE theme_id = ? AND is_enabled = true
		ORDER BY position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, themeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []*models.ComponentThemeMapping
	for rows.Next() {
		mapping := &models.ComponentThemeMapping{}
		err := rows.Scan(
			&mapping.ID,
			&mapping.ThemeID,
			&mapping.ComponentID,
			&mapping.Position,
			&mapping.IsEnabled,
			&mapping.Config,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, rows.Err()
}
