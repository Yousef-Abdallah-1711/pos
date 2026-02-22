//go:build legacy
// +build legacy

package service

import (
	"fmt"
	"pos-saas/internal/repository"
	"time"
)

// ComponentVersionService handles component versioning logic
type ComponentVersionService struct {
	versionRepo repository.ComponentVersionRepository
	componentRepo repository.ComponentRepository
}

// NewComponentVersionService creates a new component version service
func NewComponentVersionService(
	versionRepo repository.ComponentVersionRepository,
	componentRepo repository.ComponentRepository,
) *ComponentVersionService {
	return &ComponentVersionService{
		versionRepo: versionRepo,
		componentRepo: componentRepo,
	}
}

// CreateVersion creates a new version for a component
func (s *ComponentVersionService) CreateVersion(
	componentID int64,
	versionStr string,
	description string,
	releaseNotes string,
) (*models.ComponentVersion, error) {
	fmt.Printf("DEBUG: Creating version %s for component %d\n", versionStr, componentID)

	// Validate component exists
	_, err := s.componentRepo.GetByID(componentID)
	if err != nil {
		return nil, fmt.Errorf("component not found: %w", err)
	}

	// Parse semantic version
	parsed, err := models.ParseSemver(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid semantic version: %w", err)
	}

	// Set component ID and metadata
	parsed.ComponentID = componentID
	parsed.Description = description
	parsed.ReleaseNotes = releaseNotes

	// First version is always the latest
	existingVersions, _ := s.versionRepo.ListVersions(componentID)
	parsed.IsLatest = (len(existingVersions) == 0)

	// Initialize empty breaking changes if not set
	if parsed.BreakingChanges.Changes == nil {
		parsed.BreakingChanges.Changes = []string{}
	}
	if parsed.MigrationGuides == nil {
		parsed.MigrationGuides = []models.MigrationGuide{}
	}

	// Create version in database
	created, err := s.versionRepo.Create(parsed)
	if err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG: Created version %s (ID: %d) for component %d\n", versionStr, created.ID, componentID)
	return created, nil
}

// GetVersion gets a specific version of a component
func (s *ComponentVersionService) GetVersion(componentID int64, versionStr string) (*models.ComponentVersion, error) {
	return s.versionRepo.GetByComponentAndVersion(componentID, versionStr)
}

// GetLatestVersion gets the latest version of a component
func (s *ComponentVersionService) GetLatestVersion(componentID int64) (*models.ComponentVersion, error) {
	return s.versionRepo.GetLatestVersion(componentID)
}

// ListVersions lists all versions of a component
func (s *ComponentVersionService) ListVersions(componentID int64) ([]*models.ComponentVersion, error) {
	return s.versionRepo.ListVersions(componentID)
}

// SetLatestVersion marks a specific version as the latest
func (s *ComponentVersionService) SetLatestVersion(componentID int64, versionStr string) error {
	fmt.Printf("DEBUG: Setting latest version for component %d to %s\n", componentID, versionStr)

	// Verify version exists
	version, err := s.versionRepo.GetByComponentAndVersion(componentID, versionStr)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	if version.IsDeprecated {
		return fmt.Errorf("cannot set deprecated version as latest")
	}

	return s.versionRepo.SetLatestVersion(componentID, versionStr)
}

// DeprecateVersion marks a version as deprecated
func (s *ComponentVersionService) DeprecateVersion(
	componentID int64,
	versionStr string,
	message string,
	suggestedVersion string,
) error {
	fmt.Printf("DEBUG: Deprecating version %s for component %d\n", versionStr, componentID)

	// Verify version exists
	version, err := s.versionRepo.GetByComponentAndVersion(componentID, versionStr)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	if version.IsLatest {
		return fmt.Errorf("cannot deprecate the latest version")
	}

	fullMessage := message
	if suggestedVersion != "" {
		fullMessage = fmt.Sprintf("%s (upgrade to %s)", message, suggestedVersion)
	}

	return s.versionRepo.DeprecateVersion(componentID, versionStr, fullMessage)
}

// AddBreakingChanges updates breaking changes for a version
func (s *ComponentVersionService) AddBreakingChanges(
	componentID int64,
	versionStr string,
	changes []string,
	migrationRequired bool,
	migrationSteps []string,
) error {
	fmt.Printf("DEBUG: Adding breaking changes to version %s\n", versionStr)

	version, err := s.versionRepo.GetByComponentAndVersion(componentID, versionStr)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	version.BreakingChanges = models.BreakingChanges{
		HasBreakingChanges: len(changes) > 0,
		Changes:           changes,
		MigrationRequired: migrationRequired,
		MigrationSteps:    migrationSteps,
	}

	_, err = s.versionRepo.Update(version)
	return err
}

// AddMigrationGuide adds a migration guide between two versions
func (s *ComponentVersionService) AddMigrationGuide(
	componentID int64,
	fromVersion string,
	toVersion string,
	steps []string,
	autoMigrable bool,
	autoScript string,
) error {
	fmt.Printf("DEBUG: Adding migration guide from %s to %s\n", fromVersion, toVersion)

	version, err := s.versionRepo.GetByComponentAndVersion(componentID, toVersion)
	if err != nil {
		return fmt.Errorf("target version not found: %w", err)
	}

	guide := models.MigrationGuide{
		FromVersion:        fromVersion,
		ToVersion:          toVersion,
		Steps:              steps,
		AutoMigrable:       autoMigrable,
		AutoMigrationScript: autoScript,
	}

	version.MigrationGuides = append(version.MigrationGuides, guide)

	_, err = s.versionRepo.Update(version)
	return err
}

// CheckCompatibility checks if source version can upgrade to target version
func (s *ComponentVersionService) CheckCompatibility(
	sourceVersionID int64,
	targetVersionID int64,
) (compatible bool, guide *models.MigrationGuide, err error) {
	fmt.Printf("DEBUG: Checking compatibility between versions %d and %d\n", sourceVersionID, targetVersionID)

	sourceVersion, err := s.versionRepo.GetByID(sourceVersionID)
	if err != nil {
		return false, nil, fmt.Errorf("source version not found: %w", err)
	}

	targetVersion, err := s.versionRepo.GetByID(targetVersionID)
	if err != nil {
		return false, nil, fmt.Errorf("target version not found: %w", err)
	}

	// If major version changed, it's a breaking change
	if sourceVersion.Major != targetVersion.Major {
		return false, nil, fmt.Errorf("major version change - manual migration required")
	}

	// Check for migration guide in target version
	for _, mg := range targetVersion.MigrationGuides {
		if mg.FromVersion == sourceVersion.Version {
			return true, &mg, nil
		}
	}

	// If no guide found but target version is compatible, allow it
	return true, nil, nil
}

// GetCompatibleVersions gets all versions compatible with a constraint
func (s *ComponentVersionService) GetCompatibleVersions(
	componentID int64,
	constraintStr string,
) ([]*models.ComponentVersion, error) {
	fmt.Printf("DEBUG: Getting versions compatible with constraint %s for component %d\n", constraintStr, componentID)

	// Parse constraint
	constraint, err := models.ParseVersionConstraint(constraintStr)
	if err != nil {
		return nil, fmt.Errorf("invalid constraint: %w", err)
	}

	versions, err := s.versionRepo.ListVersions(componentID)
	if err != nil {
		return nil, err
	}

	var compatible []*models.ComponentVersion
	for _, v := range versions {
		if v.IsCompatibleWith(*constraint) && !v.IsDeprecated {
			compatible = append(compatible, v)
		}
	}

	return compatible, nil
}

// ValidateMigrationPath checks if a migration path exists from source to target
func (s *ComponentVersionService) ValidateMigrationPath(
	componentID int64,
	fromVersionStr string,
	toVersionStr string,
) (exists bool, steps []string, err error) {
	fmt.Printf("DEBUG: Validating migration path from %s to %s\n", fromVersionStr, toVersionStr)

	fromVersion, err := s.versionRepo.GetByComponentAndVersion(componentID, fromVersionStr)
	if err != nil {
		return false, nil, fmt.Errorf("source version not found: %w", err)
	}

	toVersion, err := s.versionRepo.GetByComponentAndVersion(componentID, toVersionStr)
	if err != nil {
		return false, nil, fmt.Errorf("target version not found: %w", err)
	}

	// If same version, no migration needed
	if fromVersion.ID == toVersion.ID {
		return true, []string{"No migration needed - versions are the same"}, nil
	}

	// Check if direct migration guide exists
	for _, guide := range toVersion.MigrationGuides {
		if guide.FromVersion == fromVersionStr {
			return true, guide.Steps, nil
		}
	}

	// Check for indirect paths (e.g., 1.0 -> 1.1 -> 1.2)
	if s.canMigrateIndirectly(componentID, fromVersion, toVersion) {
		steps, _ := s.findMigrationPath(componentID, fromVersion, toVersion)
		return true, steps, nil
	}

	return false, nil, nil
}

// canMigrateIndirectly checks if migration is possible through intermediate versions
func (s *ComponentVersionService) canMigrateIndirectly(
	componentID int64,
	from *models.ComponentVersion,
	to *models.ComponentVersion,
) bool {
	// Only allow indirect migration if major version is same
	if from.Major != to.Major {
		return false
	}

	// Allow if target is newer
	return to.Compare(from) > 0
}

// findMigrationPath finds the path to migrate from source to target version
func (s *ComponentVersionService) findMigrationPath(
	componentID int64,
	from *models.ComponentVersion,
	to *models.ComponentVersion,
) ([]string, error) {
	var steps []string

	versions, _ := s.versionRepo.ListVersions(componentID)

	// Sort versions between from and to
	var intermediate []*models.ComponentVersion
	for _, v := range versions {
		if v.Compare(from) > 0 && v.Compare(to) <= 0 && v.Major == from.Major {
			intermediate = append(intermediate, v)
		}
	}

	// Sort by version
	for _, v := range intermediate {
		steps = append(steps, fmt.Sprintf("Upgrade to %s", v.Version))
	}

	return steps, nil
}

// UpdateReleaseNotes updates release notes for a version
func (s *ComponentVersionService) UpdateReleaseNotes(
	componentID int64,
	versionStr string,
	releaseNotes string,
	changelog string,
) error {
	fmt.Printf("DEBUG: Updating release notes for version %s\n", versionStr)

	version, err := s.versionRepo.GetByComponentAndVersion(componentID, versionStr)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	version.ReleaseNotes = releaseNotes
	version.Changelog = changelog
	version.UpdatedAt = time.Now()

	_, err = s.versionRepo.Update(version)
	return err
}

// GetVersionStatistics returns statistics about versions for a component
func (s *ComponentVersionService) GetVersionStatistics(componentID int64) (map[string]interface{}, error) {
	versions, err := s.versionRepo.ListVersions(componentID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_versions": len(versions),
		"deprecated":     0,
		"with_breaking_changes": 0,
		"latest_version": "",
		"oldest_version": "",
	}

	for _, v := range versions {
		if v.IsDeprecated {
			stats["deprecated"] = stats["deprecated"].(int) + 1
		}
		if v.BreakingChanges.HasBreakingChanges {
			stats["with_breaking_changes"] = stats["with_breaking_changes"].(int) + 1
		}
		if v.IsLatest {
			stats["latest_version"] = v.Version
		}
	}

	if len(versions) > 0 {
		stats["oldest_version"] = versions[len(versions)-1].Version
	}

	return stats, nil
}
