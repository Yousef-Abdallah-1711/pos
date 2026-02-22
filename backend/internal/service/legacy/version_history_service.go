//go:build legacy
// +build legacy

/**
 * VersionHistoryService
 * Manages theme version history, snapshots, rollbacks, and comparisons
 */

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"pos-saas/internal/repository"
)

// VersionListOptions defines pagination options for listing versions
type VersionListOptions struct {
	Page  int
	Limit int
}

// VersionHistoryService manages theme version operations
type VersionHistoryService struct {
	historyRepo   repository.ThemeHistoryRepository
	snapshotRepo  repository.SnapshotRepository
	themeRepo     repository.ThemeRepositoryV2
	componentRepo repository.ComponentRepositoryV2
}

// NewVersionHistoryService creates a new version history service
func NewVersionHistoryService(
	historyRepo repository.ThemeHistoryRepository,
	snapshotRepo repository.SnapshotRepository,
	themeRepo repository.ThemeRepositoryV2,
	componentRepo repository.ComponentRepositoryV2,
) *VersionHistoryService {
	return &VersionHistoryService{
		historyRepo:   historyRepo,
		snapshotRepo:  snapshotRepo,
		themeRepo:     themeRepo,
		componentRepo: componentRepo,
	}
}

// ListVersions returns paginated version history for a theme
func (s *VersionHistoryService) ListVersions(
	themeID int64,
	opts *VersionListOptions,
) ([]*models.ThemeVersion, int, error) {
	if opts == nil {
		opts = &VersionListOptions{Page: 1, Limit: 20}
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 100 {
		opts.Limit = 20
	}

	offset := (opts.Page - 1) * opts.Limit

	// Get versions from history
	versions, err := s.historyRepo.GetVersionHistory(themeID, offset, opts.Limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get version history: %w", err)
	}

	// Get total count
	totalCount, err := s.historyRepo.GetVersionCount(themeID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get version count: %w", err)
	}

	return versions, totalCount, nil
}

// GetVersionWithSnapshot retrieves a specific version with its snapshot data
func (s *VersionHistoryService) GetVersionWithSnapshot(
	themeID int64,
	versionNumber int,
) (*models.VersionHistoryEntry, error) {
	// Get version record
	version, err := s.historyRepo.GetVersionByNumber(themeID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	// Get snapshot
	snapshot, err := s.snapshotRepo.GetSnapshot(themeID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	// Combine into VersionHistoryEntry
	entry := &models.VersionHistoryEntry{
		Version:  version,
		Snapshot: snapshot,
	}

	return entry, nil
}

// CreateSnapshot creates a new version snapshot of current theme state
func (s *VersionHistoryService) CreateSnapshot(
	themeID int64,
	changeType string,
	summary string,
	userID int64,
	userEmail string,
) (*models.ThemeVersion, error) {
	// Get current theme
	theme, err := s.themeRepo.GetByID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get theme: %w", err)
	}

	// Get all components for this theme
	components, err := s.componentRepo.GetByTheme(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	// Get next version number
	nextVersion, err := s.historyRepo.GetNextVersionNumber(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next version number: %w", err)
	}

	// Create theme snapshot JSON
	themeSnapshotJSON, err := json.Marshal(map[string]interface{}{
		"id":                theme.ID,
		"name":              theme.Name,
		"colors":            theme.Colors,
		"typography":        theme.Typography,
		"identity":          theme.Identity,
		"header":            theme.Header,
		"footer":            theme.Footer,
		"lastModifiedAt":    theme.UpdatedAt,
		"snapshotTimestamp": time.Now().UTC(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal theme snapshot: %w", err)
	}

	// Create components snapshot JSON
	componentSnapshotJSON, err := json.Marshal(components)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal components snapshot: %w", err)
	}

	// Create snapshot record
	snapshot := &models.ThemeSnapshot{
		ThemeID:              themeID,
		VersionNumber:        nextVersion,
		ThemeSnapshot:        string(themeSnapshotJSON),
		ComponentsSnapshot:   string(componentSnapshotJSON),
		SnapshotSizeBytes:    int64(len(themeSnapshotJSON) + len(componentSnapshotJSON)),
		CreatedAt:            time.Now().UTC(),
	}

	_, err = s.snapshotRepo.Create(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Create history record
	history := &models.ThemeHistory{
		ThemeID:             themeID,
		VersionNumber:       nextVersion,
		Changes:             string(themeSnapshotJSON),
		ChangeSummary:       summary,
		ChangeType:          changeType,
		AuthorName:          userEmail, // Can be enhanced to get actual name
		AuthorEmail:         userEmail,
		PreviousVersionID:   nil,
		IsCurrent:           true,
		CreatedAt:           time.Now().UTC(),
	}

	createdHistory, err := s.historyRepo.Create(history)
	if err != nil {
		return nil, fmt.Errorf("failed to create history record: %w", err)
	}

	return s.historyToVersion(createdHistory), nil
}

// RollbackToVersion restores theme to a specific version
func (s *VersionHistoryService) RollbackToVersion(
	themeID int64,
	targetVersion int,
	userID int64,
	userEmail string,
	reason string,
) (*models.Theme, error) {
	// Version rollback not yet fully implemented for v2 themes
	return nil, errors.New("version rollback not yet implemented for v2 themes")

	// Get target version snapshot
	/*targetSnapshot, err := s.snapshotRepo.GetSnapshot(themeID, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get target version snapshot: %w", err)
	}

	if targetSnapshot == nil {
		return nil, errors.New("target version not found")
	}

	// Parse theme snapshot
	var themeData map[string]interface{}
	err = json.Unmarshal([]byte(targetSnapshot.ThemeSnapshot), &themeData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse theme snapshot: %w", err)
	}

	// Parse components snapshot
	var componentsData []*models.ThemeComponent
	err = json.Unmarshal([]byte(targetSnapshot.ComponentsSnapshot), &componentsData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse components snapshot: %w", err)
	}

	// Get current theme
	currentTheme, err := s.themeRepo.GetByID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current theme: %w", err)
	}

	// Update theme with snapshot data
	currentTheme.Colors = extractMapField(themeData, "colors")
	currentTheme.Typography = extractMapField(themeData, "typography")
	currentTheme.Identity = extractMapField(themeData, "identity")
	currentTheme.HeaderConfig = extractMapField(themeData, "headerConfig")
	currentTheme.FooterConfig = extractMapField(themeData, "footerConfig")
	currentTheme.CustomCSS = extractStringField(themeData, "customCSS")
	currentTheme.Responsive = extractBoolField(themeData, "responsive")
	currentTheme.ThemeJSON = extractMapField(themeData, "themeJSON")
	currentTheme.UpdatedAt = time.Now().UTC()

	// Save updated theme
	updatedTheme, err := s.themeRepo.Update(currentTheme)
	if err != nil {
		return nil, fmt.Errorf("failed to update theme: %w", err)
	}

	// Delete existing components for theme
	err = s.componentRepo.DeleteByThemeID(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing components: %w", err)
	}

	// Restore components from snapshot
	for _, component := range componentsData {
		component.ThemeID = themeID
		_, err := s.componentRepo.Create(component)
		if err != nil {
			return nil, fmt.Errorf("failed to restore component: %w", err)
		}
	}

	// Get next version number
	nextVersion, err := s.historyRepo.GetNextVersionNumber(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next version number: %w", err)
	}

	// Create rollback snapshot
	themeSnapshotJSON, _ := json.Marshal(updatedTheme)

	rollbackSnapshot := &models.ThemeSnapshot{
		ThemeID:            themeID,
		VersionNumber:      nextVersion,
		ThemeSnapshot:      string(themeSnapshotJSON),
		ComponentsSnapshot: targetSnapshot.ComponentsSnapshot,
		SnapshotSizeBytes:  int64(len(themeSnapshotJSON) + len(targetSnapshot.ComponentsSnapshot)),
		CreatedAt:          time.Now().UTC(),
	}

	_, err = s.snapshotRepo.Create(rollbackSnapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to create rollback snapshot: %w", err)
	}

	// Create history record for rollback
	rollbackHistory := &models.ThemeHistory{
		ThemeID:           themeID,
		VersionNumber:     nextVersion,
		Changes:           string(themeSnapshotJSON),
		ChangeSummary:     fmt.Sprintf("Rolled back to version %d", targetVersion),
		ChangeType:        "rollback",
		AuthorName:        userEmail,
		AuthorEmail:       userEmail,
		PreviousVersionID: &targetSnapshot.ID,
		IsCurrent:         true,
		RollbackReason:    &reason,
		CreatedAt:         time.Now().UTC(),
	}

	_, err = s.historyRepo.Create(rollbackHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to create rollback history: %w", err)
	}

	return updatedTheme, nil
	*/
}

// CompareVersions returns differences between two versions
func (s *VersionHistoryService) CompareVersions(
	themeID int64,
	fromVersion int,
	toVersion int,
) (*models.VersionComparison, error) {
	// Get both snapshots
	fromSnapshot, err := s.snapshotRepo.GetSnapshot(themeID, fromVersion)
	if err != nil || fromSnapshot == nil {
		return nil, fmt.Errorf("from version not found")
	}

	toSnapshot, err := s.snapshotRepo.GetSnapshot(themeID, toVersion)
	if err != nil || toSnapshot == nil {
		return nil, fmt.Errorf("to version not found")
	}

	// Parse snapshots
	var fromData, toData map[string]interface{}
	json.Unmarshal([]byte(fromSnapshot.ThemeSnapshot), &fromData)
	json.Unmarshal([]byte(toSnapshot.ThemeSnapshot), &toData)

	// Calculate differences
	comparison := &models.VersionComparison{
		FromVersionNumber: fromVersion,
		ToVersionNumber:   toVersion,
		FromSnapshot:      fromSnapshot,
		ToSnapshot:        toSnapshot,
		Changes:           calculateChanges(fromData, toData),
		ChangedFieldCount: 0,
		DifferencePercent: 0,
	}

	comparison.ChangedFieldCount = len(comparison.Changes)
	if comparison.ChangedFieldCount > 0 {
		comparison.DifferencePercent = (float64(comparison.ChangedFieldCount) / float64(countFields(toData))) * 100
	}

	return comparison, nil
}

// GetVersionDiff returns detailed diff of a version compared to previous
func (s *VersionHistoryService) GetVersionDiff(
	themeID int64,
	versionNumber int,
) (map[string]interface{}, error) {
	if versionNumber < 2 {
		return nil, errors.New("cannot diff first version")
	}

	// Get current and previous snapshots
	currentSnapshot, err := s.snapshotRepo.GetSnapshot(themeID, versionNumber)
	if err != nil || currentSnapshot == nil {
		return nil, fmt.Errorf("version not found")
	}

	previousSnapshot, err := s.snapshotRepo.GetSnapshot(themeID, versionNumber-1)
	if err != nil || previousSnapshot == nil {
		return nil, fmt.Errorf("previous version not found")
	}

	// Parse and compare
	var current, previous map[string]interface{}
	json.Unmarshal([]byte(currentSnapshot.ThemeSnapshot), &current)
	json.Unmarshal([]byte(previousSnapshot.ThemeSnapshot), &previous)

	changes := calculateChanges(previous, current)

	return map[string]interface{}{
		"versionNumber": versionNumber,
		"previousVersion": versionNumber - 1,
		"changes": changes,
		"changeCount": len(changes),
	}, nil
}

// GetVersionStats returns statistics about version history
func (s *VersionHistoryService) GetVersionStats(
	themeID int64,
) (map[string]interface{}, error) {
	totalVersions, err := s.historyRepo.GetVersionCount(themeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version count: %w", err)
	}

	// Get first and last versions
	versions, err := s.historyRepo.GetVersionHistory(themeID, 0, 1)
	var oldestVersion *models.ThemeVersion
	if err == nil && len(versions) > 0 {
		oldestVersion = versions[0]
	}

	stats := map[string]interface{}{
		"totalVersions": totalVersions,
		"oldestVersion": oldestVersion,
		"createdAt":     time.Now(),
	}

	// Calculate total snapshot size
	snapshots, err := s.snapshotRepo.GetSnapshotsByTheme(themeID)
	if err == nil {
		totalSize := int64(0)
		for _, snap := range snapshots {
			totalSize += snap.SnapshotSizeBytes
		}
		stats["totalSnapshotSize"] = totalSize
		stats["averageSnapshotSize"] = totalSize / int64(len(snapshots))
	}

	return stats, nil
}

// Helper functions

func (s *VersionHistoryService) historyToVersion(history *models.ThemeHistory) *models.ThemeVersion {
	if history == nil {
		return nil
	}

	return &models.ThemeVersion{
		ID:              history.ID,
		ThemeID:         history.ThemeID,
		VersionNumber:   history.VersionNumber,
		Changes:         history.Changes,
		ChangeSummary:   history.ChangeSummary,
		ChangeType:      history.ChangeType,
		AuthorName:      history.AuthorName,
		AuthorEmail:     history.AuthorEmail,
		CreatedBy:       nil,
		CreatedAt:       history.CreatedAt,
		IsCurrent:       history.IsCurrent,
		RollbackReason:  history.RollbackReason,
		PreviousVersionID: history.PreviousVersionID,
	}
}

func calculateChanges(from, to map[string]interface{}) []map[string]interface{} {
	changes := []map[string]interface{}{}

	// Check for modifications and additions
	for key, toValue := range to {
		if fromValue, exists := from[key]; !exists {
			// Added
			changes = append(changes, map[string]interface{}{
				"type":      "added",
				"field":     key,
				"oldValue":  nil,
				"newValue":  toValue,
				"icon":      "✨",
			})
		} else if !deepEqual(fromValue, toValue) {
			// Modified
			changes = append(changes, map[string]interface{}{
				"type":      "modified",
				"field":     key,
				"oldValue":  fromValue,
				"newValue":  toValue,
				"icon":      "🔄",
			})
		}
	}

	// Check for removals
	for key, fromValue := range from {
		if _, exists := to[key]; !exists {
			changes = append(changes, map[string]interface{}{
				"type":      "removed",
				"field":     key,
				"oldValue":  fromValue,
				"newValue":  nil,
				"icon":      "🗑️",
			})
		}
	}

	return changes
}

func deepEqual(a, b interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

func countFields(data map[string]interface{}) int {
	count := len(data)
	for _, v := range data {
		if m, ok := v.(map[string]interface{}); ok {
			count += countFields(m)
		}
	}
	return count
}

func extractMapField(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}

func extractStringField(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func extractBoolField(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}
