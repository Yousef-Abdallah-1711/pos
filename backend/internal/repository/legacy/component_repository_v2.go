//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ComponentRepositoryV2 defines the interface for theme component data operations
type ComponentRepositoryV2 interface {
	Create(component *models.ThemeComponent) (*models.ThemeComponent, error)
	GetByID(id int64) (*models.ThemeComponent, error)
	GetByTheme(themeID int64) ([]*models.ThemeComponent, error)
	Update(component *models.ThemeComponent) (*models.ThemeComponent, error)
	Delete(id int64) error
	Reorder(themeID int64, components []*models.ThemeComponent) error
	GetComponentsByThemeAndType(themeID int64, componentType string) ([]*models.ThemeComponent, error)
}

// ComponentRepositoryV2Impl implements ComponentRepositoryV2 using PostgreSQL
type ComponentRepositoryV2Impl struct {
	db *sql.DB
}

// NewComponentRepositoryV2 creates a new component repository
func NewComponentRepositoryV2(db *sql.DB) ComponentRepositoryV2 {
	return &ComponentRepositoryV2Impl{db: db}
}

// Create inserts a new component into the database
func (r *ComponentRepositoryV2Impl) Create(component *models.ThemeComponent) (*models.ThemeComponent, error) {
	if component == nil {
		return nil, errors.New("component cannot be nil")
	}

	if component.ThemeID == 0 {
		return nil, errors.New("component must have a valid theme ID")
	}

	query := `
		INSERT INTO theme_components (
			theme_id, component_type, title, title_en, title_ar,
			subtitle_en, subtitle_ar,
			description_en, description_ar,
			button_text_en, button_text_ar,
			display_order, enabled, settings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	var createdID int64
	var createdAt, updatedAt time.Time

	// Extract bilingual text values or use nil/empty strings as appropriate
	subtitleEN := ""
	subtitleAR := ""
	if component.Subtitle != nil {
		subtitleEN = component.Subtitle.EN
		subtitleAR = component.Subtitle.AR
	}

	descriptionEN := ""
	descriptionAR := ""
	if component.Description != nil {
		descriptionEN = component.Description.EN
		descriptionAR = component.Description.AR
	}

	buttonTextEN := ""
	buttonTextAR := ""
	if component.ButtonText != nil {
		buttonTextEN = component.ButtonText.EN
		buttonTextAR = component.ButtonText.AR
	}

	err := r.db.QueryRow(
		query,
		component.ThemeID, component.ComponentType,
		component.Title.EN, // Old title column (for backwards compatibility)
		component.Title.EN, component.Title.AR,
		subtitleEN, subtitleAR,
		descriptionEN, descriptionAR,
		buttonTextEN, buttonTextAR,
		component.DisplayOrder, component.Enabled, component.Settings,
		now, now,
	).Scan(&createdID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create component: %w", err)
	}

	component.ID = createdID
	component.CreatedAt = createdAt
	component.UpdatedAt = updatedAt

	return component, nil
}

// GetByID retrieves a component by its ID
func (r *ComponentRepositoryV2Impl) GetByID(id int64) (*models.ThemeComponent, error) {
	query := `
		SELECT
			id, theme_id, component_type,
			title_en, title_ar,
			subtitle_en, subtitle_ar,
			description_en, description_ar,
			button_text_en, button_text_ar,
			display_order, enabled,
			settings, created_at, updated_at
		FROM theme_components
		WHERE id = $1
	`

	component := &models.ThemeComponent{}
	var settings sql.NullString
	var subtitleEN, subtitleAR sql.NullString
	var descriptionEN, descriptionAR sql.NullString
	var buttonTextEN, buttonTextAR sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&component.ID, &component.ThemeID, &component.ComponentType,
		&component.Title.EN, &component.Title.AR,
		&subtitleEN, &subtitleAR,
		&descriptionEN, &descriptionAR,
		&buttonTextEN, &buttonTextAR,
		&component.DisplayOrder, &component.Enabled,
		&settings, &component.CreatedAt, &component.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch component by id: %w", err)
	}

	// Construct optional bilingual fields
	if subtitleEN.Valid || subtitleAR.Valid {
		component.Subtitle = &models.BilingualText{
			EN: subtitleEN.String,
			AR: subtitleAR.String,
		}
	}

	if descriptionEN.Valid || descriptionAR.Valid {
		component.Description = &models.BilingualText{
			EN: descriptionEN.String,
			AR: descriptionAR.String,
		}
	}

	if buttonTextEN.Valid || buttonTextAR.Valid {
		component.ButtonText = &models.BilingualText{
			EN: buttonTextEN.String,
			AR: buttonTextAR.String,
		}
	}

	if settings.Valid {
		component.Settings = []byte(settings.String)
	}

	return component, nil
}

// GetByTheme retrieves all components for a specific theme
func (r *ComponentRepositoryV2Impl) GetByTheme(themeID int64) ([]*models.ThemeComponent, error) {
	query := `
		SELECT
			id, theme_id, component_type,
			title_en, title_ar,
			subtitle_en, subtitle_ar,
			description_en, description_ar,
			button_text_en, button_text_ar,
			display_order, enabled,
			settings, created_at, updated_at
		FROM theme_components
		WHERE theme_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query, themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query components: %w", err)
	}
	defer rows.Close()

	var components []*models.ThemeComponent
	for rows.Next() {
		component := &models.ThemeComponent{}
		var settings sql.NullString
		var subtitleEN, subtitleAR sql.NullString
		var descriptionEN, descriptionAR sql.NullString
		var buttonTextEN, buttonTextAR sql.NullString

		err := rows.Scan(
			&component.ID, &component.ThemeID, &component.ComponentType,
			&component.Title.EN, &component.Title.AR,
			&subtitleEN, &subtitleAR,
			&descriptionEN, &descriptionAR,
			&buttonTextEN, &buttonTextAR,
			&component.DisplayOrder, &component.Enabled,
			&settings, &component.CreatedAt, &component.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan component: %w", err)
		}

		// Handle optional bilingual fields
		if subtitleEN.Valid || subtitleAR.Valid {
			component.Subtitle = &models.BilingualText{
				EN: subtitleEN.String,
				AR: subtitleAR.String,
			}
		}

		if descriptionEN.Valid || descriptionAR.Valid {
			component.Description = &models.BilingualText{
				EN: descriptionEN.String,
				AR: descriptionAR.String,
			}
		}

		if buttonTextEN.Valid || buttonTextAR.Valid {
			component.ButtonText = &models.BilingualText{
				EN: buttonTextEN.String,
				AR: buttonTextAR.String,
			}
		}

		if settings.Valid {
			component.Settings = []byte(settings.String)
		}

		components = append(components, component)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating components: %w", err)
	}

	if components == nil {
		components = []*models.ThemeComponent{}
	}

	return components, nil
}

// Update updates an existing component
func (r *ComponentRepositoryV2Impl) Update(component *models.ThemeComponent) (*models.ThemeComponent, error) {
	if component == nil || component.ID == 0 {
		return nil, errors.New("component must have a valid ID")
	}

	query := `
		UPDATE theme_components
		SET
			component_type = $1,
			title = $2,
			title_en = $3,
			title_ar = $4,
			subtitle_en = $5,
			subtitle_ar = $6,
			description_en = $7,
			description_ar = $8,
			button_text_en = $9,
			button_text_ar = $10,
			display_order = $11,
			enabled = $12,
			settings = $13,
			updated_at = $14
		WHERE id = $15
		RETURNING created_at, updated_at
	`

	now := time.Now()
	var createdAt, updatedAt time.Time

	// Extract bilingual text values or use empty strings as appropriate
	subtitleEN := ""
	subtitleAR := ""
	if component.Subtitle != nil {
		subtitleEN = component.Subtitle.EN
		subtitleAR = component.Subtitle.AR
	}

	descriptionEN := ""
	descriptionAR := ""
	if component.Description != nil {
		descriptionEN = component.Description.EN
		descriptionAR = component.Description.AR
	}

	buttonTextEN := ""
	buttonTextAR := ""
	if component.ButtonText != nil {
		buttonTextEN = component.ButtonText.EN
		buttonTextAR = component.ButtonText.AR
	}

	err := r.db.QueryRow(
		query,
		component.ComponentType,
		component.Title.EN, // Old title column (for backwards compatibility)
		component.Title.EN, component.Title.AR,
		subtitleEN, subtitleAR,
		descriptionEN, descriptionAR,
		buttonTextEN, buttonTextAR,
		component.DisplayOrder, component.Enabled, component.Settings,
		now, component.ID,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update component: %w", err)
	}

	component.CreatedAt = createdAt
	component.UpdatedAt = updatedAt

	return component, nil
}

// Delete removes a component from the database
func (r *ComponentRepositoryV2Impl) Delete(id int64) error {
	query := `DELETE FROM theme_components WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete component: %w", err)
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

// Reorder updates the display order of components for a theme
func (r *ComponentRepositoryV2Impl) Reorder(themeID int64, components []*models.ThemeComponent) error {
	if len(components) == 0 {
		return errors.New("components list cannot be empty")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update display order for each component
	for _, comp := range components {
		if comp.ThemeID != themeID {
			return fmt.Errorf("component %d does not belong to theme %d", comp.ID, themeID)
		}

		query := `
			UPDATE theme_components
			SET display_order = $1, updated_at = $2
			WHERE id = $3 AND theme_id = $4
		`

		result, err := tx.Exec(query, comp.DisplayOrder, time.Now(), comp.ID, themeID)
		if err != nil {
			return fmt.Errorf("failed to update component order: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("component %d not found or doesn't belong to theme %d", comp.ID, themeID)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetComponentsByThemeAndType retrieves all components of a specific type for a theme
func (r *ComponentRepositoryV2Impl) GetComponentsByThemeAndType(themeID int64, componentType string) ([]*models.ThemeComponent, error) {
	query := `
		SELECT
			id, theme_id, component_type,
			title_en, title_ar,
			subtitle_en, subtitle_ar,
			description_en, description_ar,
			button_text_en, button_text_ar,
			display_order, enabled,
			settings, created_at, updated_at
		FROM theme_components
		WHERE theme_id = $1 AND component_type = $2 AND enabled = true
		ORDER BY display_order ASC
	`

	rows, err := r.db.Query(query, themeID, componentType)
	if err != nil {
		return nil, fmt.Errorf("failed to query components: %w", err)
	}
	defer rows.Close()

	var components []*models.ThemeComponent
	for rows.Next() {
		component := &models.ThemeComponent{}
		var settings sql.NullString
		var subtitleEN, subtitleAR sql.NullString
		var descriptionEN, descriptionAR sql.NullString
		var buttonTextEN, buttonTextAR sql.NullString

		err := rows.Scan(
			&component.ID, &component.ThemeID, &component.ComponentType,
			&component.Title.EN, &component.Title.AR,
			&subtitleEN, &subtitleAR,
			&descriptionEN, &descriptionAR,
			&buttonTextEN, &buttonTextAR,
			&component.DisplayOrder, &component.Enabled,
			&settings, &component.CreatedAt, &component.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan component: %w", err)
		}

		// Construct optional bilingual fields
		if subtitleEN.Valid || subtitleAR.Valid {
			component.Subtitle = &models.BilingualText{
				EN: subtitleEN.String,
				AR: subtitleAR.String,
			}
		}

		if descriptionEN.Valid || descriptionAR.Valid {
			component.Description = &models.BilingualText{
				EN: descriptionEN.String,
				AR: descriptionAR.String,
			}
		}

		if buttonTextEN.Valid || buttonTextAR.Valid {
			component.ButtonText = &models.BilingualText{
				EN: buttonTextEN.String,
				AR: buttonTextAR.String,
			}
		}

		if settings.Valid {
			component.Settings = []byte(settings.String)
		}

		components = append(components, component)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating components: %w", err)
	}

	if components == nil {
		components = []*models.ThemeComponent{}
	}

	return components, nil
}
