//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

// ComponentLibraryRepositoryV2 defines the interface for component library data operations
type ComponentLibraryRepositoryV2 interface {
	GetByType(componentType string) (*models.ComponentLibrary, error)
	GetByID(id int64) (*models.ComponentLibrary, error)
	List() ([]*models.ComponentLibrary, error)
	GetActive() ([]*models.ComponentLibrary, error)
	Create(library *models.ComponentLibrary) (*models.ComponentLibrary, error)
	Update(library *models.ComponentLibrary) (*models.ComponentLibrary, error)
	Delete(id int64) error
}

// ComponentLibraryRepositoryV2Impl implements ComponentLibraryRepositoryV2 using PostgreSQL
type ComponentLibraryRepositoryV2Impl struct {
	db *sql.DB
}

// NewComponentLibraryRepositoryV2 creates a new component library repository
func NewComponentLibraryRepositoryV2(db *sql.DB) ComponentLibraryRepositoryV2 {
	return &ComponentLibraryRepositoryV2Impl{db: db}
}

// GetByType retrieves a component library entry by component type
func (r *ComponentLibraryRepositoryV2Impl) GetByType(componentType string) (*models.ComponentLibrary, error) {
	if componentType == "" {
		return nil, errors.New("component type cannot be empty")
	}

	query := `
		SELECT
			id, component_type, title, description, icon_url,
			default_settings, settings_schema, is_active, created_at, updated_at
		FROM component_library
		WHERE component_type = $1
	`

	library := &models.ComponentLibrary{}
	var defaultSettings, settingsSchema sql.NullString

	err := r.db.QueryRow(query, componentType).Scan(
		&library.ID, &library.ComponentType, &library.Title, &library.Description,
		&library.IconURL, &defaultSettings, &settingsSchema, &library.IsActive,
		&library.CreatedAt, &library.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch component library by type: %w", err)
	}

	if defaultSettings.Valid {
		library.DefaultSettings = []byte(defaultSettings.String)
	}

	if settingsSchema.Valid {
		library.SettingsSchema = []byte(settingsSchema.String)
	}

	return library, nil
}

// GetByID retrieves a component library entry by ID
func (r *ComponentLibraryRepositoryV2Impl) GetByID(id int64) (*models.ComponentLibrary, error) {
	query := `
		SELECT
			id, component_type, title, description, icon_url,
			default_settings, settings_schema, is_active, created_at, updated_at
		FROM component_library
		WHERE id = $1
	`

	library := &models.ComponentLibrary{}
	var defaultSettings, settingsSchema sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&library.ID, &library.ComponentType, &library.Title, &library.Description,
		&library.IconURL, &defaultSettings, &settingsSchema, &library.IsActive,
		&library.CreatedAt, &library.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch component library by id: %w", err)
	}

	if defaultSettings.Valid {
		library.DefaultSettings = []byte(defaultSettings.String)
	}

	if settingsSchema.Valid {
		library.SettingsSchema = []byte(settingsSchema.String)
	}

	return library, nil
}

// List retrieves all component library entries
func (r *ComponentLibraryRepositoryV2Impl) List() ([]*models.ComponentLibrary, error) {
	query := `
		SELECT
			id, component_type, title, description, icon_url,
			default_settings, settings_schema, is_active, created_at, updated_at
		FROM component_library
		ORDER BY component_type ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query component library: %w", err)
	}
	defer rows.Close()

	var libraries []*models.ComponentLibrary
	for rows.Next() {
		library := &models.ComponentLibrary{}
		var defaultSettings, settingsSchema sql.NullString

		err := rows.Scan(
			&library.ID, &library.ComponentType, &library.Title, &library.Description,
			&library.IconURL, &defaultSettings, &settingsSchema, &library.IsActive,
			&library.CreatedAt, &library.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan component library: %w", err)
		}

		if defaultSettings.Valid {
			library.DefaultSettings = []byte(defaultSettings.String)
		}

		if settingsSchema.Valid {
			library.SettingsSchema = []byte(settingsSchema.String)
		}

		libraries = append(libraries, library)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating component library: %w", err)
	}

	if libraries == nil {
		libraries = []*models.ComponentLibrary{}
	}

	return libraries, nil
}

// GetActive retrieves all active component library entries
func (r *ComponentLibraryRepositoryV2Impl) GetActive() ([]*models.ComponentLibrary, error) {
	query := `
		SELECT
			id, component_type, title, description, icon_url,
			default_settings, settings_schema, is_active, created_at, updated_at
		FROM component_library
		WHERE is_active = true
		ORDER BY component_type ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active component library: %w", err)
	}
	defer rows.Close()

	var libraries []*models.ComponentLibrary
	for rows.Next() {
		library := &models.ComponentLibrary{}
		var defaultSettings, settingsSchema sql.NullString

		err := rows.Scan(
			&library.ID, &library.ComponentType, &library.Title, &library.Description,
			&library.IconURL, &defaultSettings, &settingsSchema, &library.IsActive,
			&library.CreatedAt, &library.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan component library: %w", err)
		}

		if defaultSettings.Valid {
			library.DefaultSettings = []byte(defaultSettings.String)
		}

		if settingsSchema.Valid {
			library.SettingsSchema = []byte(settingsSchema.String)
		}

		libraries = append(libraries, library)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating component library: %w", err)
	}

	if libraries == nil {
		libraries = []*models.ComponentLibrary{}
	}

	return libraries, nil
}

// Create inserts a new component library entry
func (r *ComponentLibraryRepositoryV2Impl) Create(library *models.ComponentLibrary) (*models.ComponentLibrary, error) {
	if library == nil {
		return nil, errors.New("library cannot be nil")
	}

	if library.ComponentType == "" {
		return nil, errors.New("component type cannot be empty")
	}

	query := `
		INSERT INTO component_library (
			component_type, title, description, icon_url,
			default_settings, settings_schema, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		RETURNING id, created_at, updated_at
	`

	var createdID int64

	err := r.db.QueryRow(
		query,
		library.ComponentType, library.Title, library.Description, library.IconURL,
		library.DefaultSettings, library.SettingsSchema, library.IsActive,
	).Scan(&createdID, &library.CreatedAt, &library.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create component library: %w", err)
	}

	library.ID = createdID

	return library, nil
}

// Update updates an existing component library entry
func (r *ComponentLibraryRepositoryV2Impl) Update(library *models.ComponentLibrary) (*models.ComponentLibrary, error) {
	if library == nil || library.ID == 0 {
		return nil, errors.New("library must have a valid ID")
	}

	query := `
		UPDATE component_library
		SET
			component_type = $1,
			title = $2,
			description = $3,
			icon_url = $4,
			default_settings = $5,
			settings_schema = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt interface{}

	err := r.db.QueryRow(
		query,
		library.ComponentType, library.Title, library.Description, library.IconURL,
		library.DefaultSettings, library.SettingsSchema, library.IsActive, library.ID,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update component library: %w", err)
	}

	return library, nil
}

// Delete removes a component library entry
func (r *ComponentLibraryRepositoryV2Impl) Delete(id int64) error {
	query := `DELETE FROM component_library WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete component library: %w", err)
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
