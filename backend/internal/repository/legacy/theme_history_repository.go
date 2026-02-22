//go:build legacy
// +build legacy

/**
 * ThemeHistoryRepository Interface
 * Defines operations for theme version history records
 */

package repository

import (
)

// ThemeHistoryRepository defines operations for theme history
type ThemeHistoryRepository interface {
	// Create saves a new history record
	Create(history *models.ThemeHistory) (*models.ThemeHistory, error)

	// GetVersionHistory retrieves paginated version history for a theme
	GetVersionHistory(themeID int64, offset, limit int) ([]*models.ThemeVersion, error)

	// GetVersionByNumber retrieves a specific version by number
	GetVersionByNumber(themeID int64, versionNumber int) (*models.ThemeVersion, error)

	// GetVersionCount returns total number of versions for a theme
	GetVersionCount(themeID int64) (int, error)

	// GetNextVersionNumber returns the next version number for a theme
	GetNextVersionNumber(themeID int64) (int, error)

	// MarkAsNonCurrent marks a version as not current
	MarkAsNonCurrent(themeID int64) error

	// MarkAsCurrent marks a version as the current one
	MarkAsCurrent(themeID int64, versionNumber int) error

	// DeleteOldVersions deletes versions older than specified number of days
	DeleteOldVersions(themeID int64, daysOld int) error

	// GetVersionsByChangeType retrieves versions of specific change type
	GetVersionsByChangeType(themeID int64, changeType string) ([]*models.ThemeVersion, error)

	// GetVersionsByAuthor retrieves versions created by specific author
	GetVersionsByAuthor(themeID int64, authorEmail string) ([]*models.ThemeVersion, error)
}
