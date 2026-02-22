//go:build legacy
// +build legacy

/**
 * SnapshotRepository Interface
 * Defines operations for theme snapshot storage
 */

package repository

import (
)

// SnapshotRepository defines operations for theme snapshots
type SnapshotRepository interface {
	// Create saves a new snapshot
	Create(snapshot *models.ThemeSnapshot) (*models.ThemeSnapshot, error)

	// GetSnapshot retrieves a specific snapshot by theme and version
	GetSnapshot(themeID int64, versionNumber int) (*models.ThemeSnapshot, error)

	// GetSnapshotsByTheme retrieves all snapshots for a theme
	GetSnapshotsByTheme(themeID int64) ([]*models.ThemeSnapshot, error)

	// GetSnapshotPaginated retrieves snapshots with pagination
	GetSnapshotPaginated(themeID int64, offset, limit int) ([]*models.ThemeSnapshot, error)

	// UpdateSnapshotSize updates the size metadata of a snapshot
	UpdateSnapshotSize(themeID int64, versionNumber int, sizeBytes int64) error

	// DeleteSnapshot removes a specific snapshot
	DeleteSnapshot(themeID int64, versionNumber int) error

	// DeleteOldSnapshots removes snapshots older than specified days
	DeleteOldSnapshots(themeID int64, daysOld int) error

	// GetTotalSnapshotSize returns total size of all snapshots for a theme
	GetTotalSnapshotSize(themeID int64) (int64, error)

	// VerifySnapshotIntegrity checks if snapshot data is complete
	VerifySnapshotIntegrity(themeID int64, versionNumber int) (bool, error)
}
