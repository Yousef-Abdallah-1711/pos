//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ThemeHistoryRepositoryV2 defines the interface for theme history data operations
type ThemeHistoryRepositoryV2 interface {
	Create(history *models.ThemeHistory) (*models.ThemeHistory, error)
	GetByTheme(themeID int64) ([]*models.ThemeHistory, error)
	GetByThemeAndVersion(themeID int64, version int) (*models.ThemeHistory, error)
	GetLatestVersion(themeID int64) (int, error)
	DeleteByTheme(themeID int64) error
}

// ThemeHistoryRepositoryV2Impl implements ThemeHistoryRepositoryV2 using PostgreSQL
type ThemeHistoryRepositoryV2Impl struct {
	db *sql.DB
}

// NewThemeHistoryRepositoryV2 creates a new theme history repository
func NewThemeHistoryRepositoryV2(db *sql.DB) ThemeHistoryRepositoryV2 {
	return &ThemeHistoryRepositoryV2Impl{db: db}
}

// Create inserts a new history entry
func (r *ThemeHistoryRepositoryV2Impl) Create(history *models.ThemeHistory) (*models.ThemeHistory, error) {
	if history == nil {
		return nil, errors.New("history cannot be nil")
	}

	if history.ThemeID == 0 {
		return nil, errors.New("history must have a valid theme ID")
	}

	// Get the next version number
	nextVersion, err := r.GetLatestVersion(history.ThemeID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}
	nextVersion++

	query := `
		INSERT INTO theme_history (
			theme_id, version, changes, change_summary, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, created_at
	`

	now := time.Now()
	var createdID int64
	var createdAt time.Time

	err = r.db.QueryRow(
		query,
		history.ThemeID, nextVersion, history.Changes, history.ChangeSummary,
		history.CreatedBy, now,
	).Scan(&createdID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create theme history: %w", err)
	}

	history.ID = createdID
	history.VersionNumber = nextVersion
	history.CreatedAt = createdAt

	return history, nil
}

// GetByTheme retrieves all history entries for a specific theme
func (r *ThemeHistoryRepositoryV2Impl) GetByTheme(themeID int64) ([]*models.ThemeHistory, error) {
	query := `
		SELECT
			id, theme_id, version, changes, change_summary, created_by, created_at
		FROM theme_history
		WHERE theme_id = $1
		ORDER BY version DESC
	`

	rows, err := r.db.Query(query, themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query theme history: %w", err)
	}
	defer rows.Close()

	var histories []*models.ThemeHistory
	for rows.Next() {
		history := &models.ThemeHistory{}
		var createdBy sql.NullInt64

		err := rows.Scan(
			&history.ID, &history.ThemeID, &history.VersionNumber,
			&history.Changes, &history.ChangeSummary, &createdBy,
			&history.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan theme history: %w", err)
		}

		if createdBy.Valid {
			history.CreatedBy = &createdBy.Int64
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating theme history: %w", err)
	}

	if histories == nil {
		histories = []*models.ThemeHistory{}
	}

	return histories, nil
}

// GetByThemeAndVersion retrieves a specific version of theme history
func (r *ThemeHistoryRepositoryV2Impl) GetByThemeAndVersion(themeID int64, version int) (*models.ThemeHistory, error) {
	query := `
		SELECT
			id, theme_id, version, changes, change_summary, created_by, created_at
		FROM theme_history
		WHERE theme_id = $1 AND version = $2
	`

	history := &models.ThemeHistory{}
	var createdBy sql.NullInt64

	err := r.db.QueryRow(query, themeID, version).Scan(
		&history.ID, &history.ThemeID, &history.VersionNumber,
		&history.Changes, &history.ChangeSummary, &createdBy,
		&history.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch theme history: %w", err)
	}

	if createdBy.Valid {
		history.CreatedBy = &createdBy.Int64
	}

	return history, nil
}

// GetLatestVersion retrieves the latest version number for a theme
func (r *ThemeHistoryRepositoryV2Impl) GetLatestVersion(themeID int64) (int, error) {
	query := `
		SELECT COALESCE(MAX(version), 0)
		FROM theme_history
		WHERE theme_id = $1
	`

	var latestVersion int

	err := r.db.QueryRow(query, themeID).Scan(&latestVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	return latestVersion, nil
}

// DeleteByTheme deletes all history entries for a specific theme
func (r *ThemeHistoryRepositoryV2Impl) DeleteByTheme(themeID int64) error {
	query := `DELETE FROM theme_history WHERE theme_id = $1`

	result, err := r.db.Exec(query, themeID)
	if err != nil {
		return fmt.Errorf("failed to delete theme history: %w", err)
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
