//go:build legacy
// +build legacy

package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// ComponentVersionRepository handles database operations for component versions
type ComponentVersionRepository interface {
	// Create creates a new component version
	Create(version *models.ComponentVersion) (*models.ComponentVersion, error)

	// GetByID gets a component version by ID
	GetByID(id int64) (*models.ComponentVersion, error)

	// GetByComponentAndVersion gets a specific version of a component
	GetByComponentAndVersion(componentID int64, versionStr string) (*models.ComponentVersion, error)

	// GetLatestVersion gets the latest version of a component
	GetLatestVersion(componentID int64) (*models.ComponentVersion, error)

	// ListVersions lists all versions of a component
	ListVersions(componentID int64) ([]*models.ComponentVersion, error)

	// ListCompatibleVersions lists versions compatible with a constraint
	ListCompatibleVersions(componentID int64, constraint string) ([]*models.ComponentVersion, error)

	// Update updates a component version
	Update(version *models.ComponentVersion) (*models.ComponentVersion, error)

	// SetLatestVersion marks a version as latest
	SetLatestVersion(componentID int64, versionStr string) error

	// DeprecateVersion marks a version as deprecated
	DeprecateVersion(componentID int64, versionStr string, message string) error

	// Delete deletes a component version
	Delete(id int64) error

	// GetCompatibility gets compatibility info between two versions
	GetCompatibility(sourceVersionID, targetVersionID int64) (*models.ComponentCompatibility, error)

	// CreateCompatibility records compatibility between two versions
	CreateCompatibility(compat *models.ComponentCompatibility) error

	// UpdateCompatibility updates compatibility info
	UpdateCompatibility(compat *models.ComponentCompatibility) error
}

// ComponentCompatibility represents compatibility between two versions
type ComponentCompatibility struct {
	ID                     int64
	SourceVersionID        int64
	TargetVersionID        int64
	IsCompatible           bool
	CompatibilityNotes     string
	AutoMigrable           bool
	MigrationScript        string
	EstimatedMigrationTime int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// MySQLComponentVersionRepository implements ComponentVersionRepository for MySQL
type MySQLComponentVersionRepository struct {
	db *sql.DB
}

// NewComponentVersionRepository creates a new component version repository
func NewComponentVersionRepository(db *sql.DB) ComponentVersionRepository {
	return &MySQLComponentVersionRepository{db: db}
}

// Create implements ComponentVersionRepository.Create
func (r *MySQLComponentVersionRepository) Create(version *models.ComponentVersion) (*models.ComponentVersion, error) {
	query := `
		INSERT INTO component_versions (
			component_id,
			version,
			major,
			minor,
			patch,
			prerelease,
			metadata,
			description,
			release_notes,
			changelog,
			is_latest,
			is_deprecated,
			deprecation_message,
			breaking_changes,
			migration_guides,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	breakingChangesJSON, _ := json.Marshal(version.BreakingChanges)
	migrationGuidesJSON, _ := json.Marshal(version.MigrationGuides)

	err := r.db.QueryRow(
		query,
		version.ComponentID,
		version.Version,
		version.Major,
		version.Minor,
		version.Patch,
		version.Prerelease,
		version.Metadata,
		version.Description,
		version.ReleaseNotes,
		version.Changelog,
		version.IsLatest,
		version.IsDeprecated,
		version.DeprecationMsg,
		breakingChangesJSON,
		migrationGuidesJSON,
		time.Now(),
		time.Now(),
	).Scan(&version.ID, &version.CreatedAt, &version.UpdatedAt)

	if err != nil {
		fmt.Printf("DEBUG: Failed to create component version: %v\n", err)
		return nil, fmt.Errorf("failed to create component version: %w", err)
	}

	fmt.Printf("DEBUG: Created component version %d: %s\n", version.ID, version.Version)
	return version, nil
}

// GetByID implements ComponentVersionRepository.GetByID
func (r *MySQLComponentVersionRepository) GetByID(id int64) (*models.ComponentVersion, error) {
	query := `
		SELECT
			id, component_id, version, major, minor, patch, prerelease, metadata,
			description, release_notes, changelog,
			is_latest, is_deprecated, deprecation_message,
			breaking_changes, migration_guides,
			created_at, updated_at
		FROM component_versions
		WHERE id = ?
	`

	version := &models.ComponentVersion{}
	var breakingChangesJSON, migrationGuidesJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&version.ID, &version.ComponentID, &version.Version,
		&version.Major, &version.Minor, &version.Patch,
		&version.Prerelease, &version.Metadata,
		&version.Description, &version.ReleaseNotes, &version.Changelog,
		&version.IsLatest, &version.IsDeprecated, &version.DeprecationMsg,
		&breakingChangesJSON, &migrationGuidesJSON,
		&version.CreatedAt, &version.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("component version not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get component version: %w", err)
	}

	// Parse JSON fields
	if breakingChangesJSON != nil {
		json.Unmarshal(breakingChangesJSON, &version.BreakingChanges)
	}
	if migrationGuidesJSON != nil {
		json.Unmarshal(migrationGuidesJSON, &version.MigrationGuides)
	}

	return version, nil
}

// GetByComponentAndVersion implements ComponentVersionRepository.GetByComponentAndVersion
func (r *MySQLComponentVersionRepository) GetByComponentAndVersion(componentID int64, versionStr string) (*models.ComponentVersion, error) {
	query := `
		SELECT
			id, component_id, version, major, minor, patch, prerelease, metadata,
			description, release_notes, changelog,
			is_latest, is_deprecated, deprecation_message,
			breaking_changes, migration_guides,
			created_at, updated_at
		FROM component_versions
		WHERE component_id = ? AND version = ?
	`

	version := &models.ComponentVersion{}
	var breakingChangesJSON, migrationGuidesJSON []byte

	err := r.db.QueryRow(query, componentID, versionStr).Scan(
		&version.ID, &version.ComponentID, &version.Version,
		&version.Major, &version.Minor, &version.Patch,
		&version.Prerelease, &version.Metadata,
		&version.Description, &version.ReleaseNotes, &version.Changelog,
		&version.IsLatest, &version.IsDeprecated, &version.DeprecationMsg,
		&breakingChangesJSON, &migrationGuidesJSON,
		&version.CreatedAt, &version.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("component version not found: %d@%s", componentID, versionStr)
		}
		return nil, fmt.Errorf("failed to get component version: %w", err)
	}

	// Parse JSON fields
	if breakingChangesJSON != nil {
		json.Unmarshal(breakingChangesJSON, &version.BreakingChanges)
	}
	if migrationGuidesJSON != nil {
		json.Unmarshal(migrationGuidesJSON, &version.MigrationGuides)
	}

	return version, nil
}

// GetLatestVersion implements ComponentVersionRepository.GetLatestVersion
func (r *MySQLComponentVersionRepository) GetLatestVersion(componentID int64) (*models.ComponentVersion, error) {
	query := `
		SELECT
			id, component_id, version, major, minor, patch, prerelease, metadata,
			description, release_notes, changelog,
			is_latest, is_deprecated, deprecation_message,
			breaking_changes, migration_guides,
			created_at, updated_at
		FROM component_versions
		WHERE component_id = ? AND is_latest = true
		ORDER BY created_at DESC
		LIMIT 1
	`

	version := &models.ComponentVersion{}
	var breakingChangesJSON, migrationGuidesJSON []byte

	err := r.db.QueryRow(query, componentID).Scan(
		&version.ID, &version.ComponentID, &version.Version,
		&version.Major, &version.Minor, &version.Patch,
		&version.Prerelease, &version.Metadata,
		&version.Description, &version.ReleaseNotes, &version.Changelog,
		&version.IsLatest, &version.IsDeprecated, &version.DeprecationMsg,
		&breakingChangesJSON, &migrationGuidesJSON,
		&version.CreatedAt, &version.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no latest version found for component: %d", componentID)
		}
		return nil, fmt.Errorf("failed to get latest component version: %w", err)
	}

	// Parse JSON fields
	if breakingChangesJSON != nil {
		json.Unmarshal(breakingChangesJSON, &version.BreakingChanges)
	}
	if migrationGuidesJSON != nil {
		json.Unmarshal(migrationGuidesJSON, &version.MigrationGuides)
	}

	return version, nil
}

// ListVersions implements ComponentVersionRepository.ListVersions
func (r *MySQLComponentVersionRepository) ListVersions(componentID int64) ([]*models.ComponentVersion, error) {
	query := `
		SELECT
			id, component_id, version, major, minor, patch, prerelease, metadata,
			description, release_notes, changelog,
			is_latest, is_deprecated, deprecation_message,
			breaking_changes, migration_guides,
			created_at, updated_at
		FROM component_versions
		WHERE component_id = ?
		ORDER BY major DESC, minor DESC, patch DESC, created_at DESC
	`

	rows, err := r.db.Query(query, componentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list component versions: %w", err)
	}
	defer rows.Close()

	var versions []*models.ComponentVersion

	for rows.Next() {
		version := &models.ComponentVersion{}
		var breakingChangesJSON, migrationGuidesJSON []byte

		err := rows.Scan(
			&version.ID, &version.ComponentID, &version.Version,
			&version.Major, &version.Minor, &version.Patch,
			&version.Prerelease, &version.Metadata,
			&version.Description, &version.ReleaseNotes, &version.Changelog,
			&version.IsLatest, &version.IsDeprecated, &version.DeprecationMsg,
			&breakingChangesJSON, &migrationGuidesJSON,
			&version.CreatedAt, &version.UpdatedAt,
		)

		if err != nil {
			fmt.Printf("DEBUG: Error scanning version row: %v\n", err)
			continue
		}

		// Parse JSON fields
		if breakingChangesJSON != nil {
			json.Unmarshal(breakingChangesJSON, &version.BreakingChanges)
		}
		if migrationGuidesJSON != nil {
			json.Unmarshal(migrationGuidesJSON, &version.MigrationGuides)
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// ListCompatibleVersions implements ComponentVersionRepository.ListCompatibleVersions
func (r *MySQLComponentVersionRepository) ListCompatibleVersions(componentID int64, constraintStr string) ([]*models.ComponentVersion, error) {
	versions, err := r.ListVersions(componentID)
	if err != nil {
		return nil, err
	}

	constraint, err := models.ParseVersionConstraint(constraintStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint: %w", err)
	}

	var compatible []*models.ComponentVersion
	for _, v := range versions {
		if v.IsCompatibleWith(*constraint) && !v.IsDeprecated {
			compatible = append(compatible, v)
		}
	}

	return compatible, nil
}

// Update implements ComponentVersionRepository.Update
func (r *MySQLComponentVersionRepository) Update(version *models.ComponentVersion) (*models.ComponentVersion, error) {
	query := `
		UPDATE component_versions
		SET
			description = ?,
			release_notes = ?,
			changelog = ?,
			is_latest = ?,
			is_deprecated = ?,
			deprecation_message = ?,
			breaking_changes = ?,
			migration_guides = ?,
			updated_at = ?
		WHERE id = ?
	`

	breakingChangesJSON, _ := json.Marshal(version.BreakingChanges)
	migrationGuidesJSON, _ := json.Marshal(version.MigrationGuides)

	_, err := r.db.Exec(
		query,
		version.Description,
		version.ReleaseNotes,
		version.Changelog,
		version.IsLatest,
		version.IsDeprecated,
		version.DeprecationMsg,
		breakingChangesJSON,
		migrationGuidesJSON,
		time.Now(),
		version.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update component version: %w", err)
	}

	version.UpdatedAt = time.Now()
	fmt.Printf("DEBUG: Updated component version %d: %s\n", version.ID, version.Version)
	return version, nil
}

// SetLatestVersion implements ComponentVersionRepository.SetLatestVersion
func (r *MySQLComponentVersionRepository) SetLatestVersion(componentID int64, versionStr string) error {
	// First, unset all other versions as latest
	_, err := r.db.Exec(
		"UPDATE component_versions SET is_latest = false WHERE component_id = ?",
		componentID,
	)
	if err != nil {
		return fmt.Errorf("failed to unset previous latest version: %w", err)
	}

	// Then set the specified version as latest
	result, err := r.db.Exec(
		"UPDATE component_versions SET is_latest = true WHERE component_id = ? AND version = ?",
		componentID,
		versionStr,
	)
	if err != nil {
		return fmt.Errorf("failed to set latest version: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("version not found: %d@%s", componentID, versionStr)
	}

	fmt.Printf("DEBUG: Set latest version for component %d to %s\n", componentID, versionStr)
	return nil
}

// DeprecateVersion implements ComponentVersionRepository.DeprecateVersion
func (r *MySQLComponentVersionRepository) DeprecateVersion(componentID int64, versionStr string, message string) error {
	_, err := r.db.Exec(
		"UPDATE component_versions SET is_deprecated = true, deprecation_message = ? WHERE component_id = ? AND version = ?",
		message,
		componentID,
		versionStr,
	)

	if err != nil {
		return fmt.Errorf("failed to deprecate version: %w", err)
	}

	fmt.Printf("DEBUG: Deprecated component version %d@%s\n", componentID, versionStr)
	return nil
}

// Delete implements ComponentVersionRepository.Delete
func (r *MySQLComponentVersionRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM component_versions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete component version: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("component version not found: %d", id)
	}

	fmt.Printf("DEBUG: Deleted component version %d\n", id)
	return nil
}

// GetCompatibility implements ComponentVersionRepository.GetCompatibility
func (r *MySQLComponentVersionRepository) GetCompatibility(sourceVersionID, targetVersionID int64) (*models.ComponentCompatibility, error) {
	query := `
		SELECT id, source_version_id, target_version_id, is_compatible, compatibility_notes,
		       auto_migrable, migration_script, estimated_migration_minutes, created_at, updated_at
		FROM component_compatibility
		WHERE source_version_id = ? AND target_version_id = ?
	`

	compat := &models.ComponentCompatibility{}
	err := r.db.QueryRow(query, sourceVersionID, targetVersionID).Scan(
		&compat.ID, &compat.SourceVersionID, &compat.TargetVersionID,
		&compat.IsCompatible, &compat.CompatibilityNotes,
		&compat.AutoMigrable, &compat.MigrationScript,
		&compat.EstimatedMigrationTime, &compat.CreatedAt, &compat.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("compatibility not found")
		}
		return nil, fmt.Errorf("failed to get compatibility: %w", err)
	}

	return compat, nil
}

// CreateCompatibility implements ComponentVersionRepository.CreateCompatibility
func (r *MySQLComponentVersionRepository) CreateCompatibility(compat *models.ComponentCompatibility) error {
	query := `
		INSERT INTO component_compatibility (
			source_version_id, target_version_id, is_compatible,
			compatibility_notes, auto_migrable, migration_script,
			estimated_migration_minutes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		compat.SourceVersionID, compat.TargetVersionID, compat.IsCompatible,
		compat.CompatibilityNotes, compat.AutoMigrable, compat.MigrationScript,
		compat.EstimatedMigrationTime, time.Now(), time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create compatibility record: %w", err)
	}

	id, _ := result.LastInsertId()
	compat.ID = id
	compat.CreatedAt = time.Now()
	compat.UpdatedAt = time.Now()

	fmt.Printf("DEBUG: Created compatibility record %d\n", id)
	return nil
}

// UpdateCompatibility implements ComponentVersionRepository.UpdateCompatibility
func (r *MySQLComponentVersionRepository) UpdateCompatibility(compat *models.ComponentCompatibility) error {
	query := `
		UPDATE component_compatibility
		SET is_compatible = ?, compatibility_notes = ?,
		    auto_migrable = ?, migration_script = ?,
		    estimated_migration_minutes = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		compat.IsCompatible, compat.CompatibilityNotes,
		compat.AutoMigrable, compat.MigrationScript,
		compat.EstimatedMigrationTime, time.Now(),
		compat.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update compatibility: %w", err)
	}

	fmt.Printf("DEBUG: Updated compatibility record %d\n", compat.ID)
	return nil
}
