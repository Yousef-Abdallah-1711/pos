//go:build legacy
// +build legacy

/**
 * Version History Models
 * Defines structures for version control and rollback functionality
 */

package models

import (
	"time"
)

// ThemeVersion represents a version in the version history
type ThemeVersion struct {
	ID              int64      `json:"id" db:"id"`
	ThemeID         int64      `json:"themeId" db:"theme_id"`
	VersionNumber   int        `json:"versionNumber" db:"version_number"`
	Changes         string     `json:"changes" db:"changes"`
	ChangeSummary   string     `json:"changeSummary" db:"change_summary"`
	ChangeType      string     `json:"changeType" db:"change_type"` // create, update, rollback, activate, publish, duplicate
	AuthorName      string     `json:"authorName" db:"author_name"`
	AuthorEmail     string     `json:"authorEmail" db:"author_email"`
	CreatedBy       *int64     `json:"createdBy" db:"created_by"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	IsCurrent       bool       `json:"isCurrent" db:"is_current"`
	RollbackReason  *string    `json:"rollbackReason" db:"rollback_reason"`
	PreviousVersionID *int64   `json:"previousVersionId" db:"previous_version_id"`
}

// ThemeSnapshot stores a complete snapshot of theme state at a specific version
type ThemeSnapshot struct {
	ID                 int64     `json:"id" db:"id"`
	ThemeID            int64     `json:"themeId" db:"theme_id"`
	VersionNumber      int       `json:"versionNumber" db:"version_number"`
	ThemeSnapshot      string    `json:"themeSnapshot" db:"theme_snapshot"` // JSONB stored as string
	ComponentsSnapshot string    `json:"componentsSnapshot" db:"components_snapshot"` // JSONB stored as string
	SnapshotSizeBytes  int64     `json:"snapshotSizeBytes" db:"snapshot_size_bytes"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

// ThemeHistory represents a record in the theme_history table
type ThemeHistory struct {
	ID                int64      `json:"id" db:"id"`
	ThemeID           int64      `json:"themeId" db:"theme_id"`
	VersionNumber     int        `json:"versionNumber" db:"version_number"`
	Changes           string     `json:"changes" db:"changes"`
	ChangeSummary     string     `json:"changeSummary" db:"change_summary"`
	ChangeType        string     `json:"changeType" db:"change_type"`
	AuthorName        string     `json:"authorName" db:"author_name"`
	AuthorEmail       string     `json:"authorEmail" db:"author_email"`
	CreatedBy         *int64     `json:"createdBy" db:"created_by"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	IsCurrent         bool       `json:"isCurrent" db:"is_current"`
	RollbackReason    *string    `json:"rollbackReason" db:"rollback_reason"`
	PreviousVersionID *int64     `json:"previousVersionId" db:"previous_version_id"`
}

// VersionHistoryEntry combines version info with snapshot data
type VersionHistoryEntry struct {
	Version  *ThemeVersion   `json:"version"`
	Snapshot *ThemeSnapshot  `json:"snapshot"`
}

// PropertyChange represents a single property change
type PropertyChange struct {
	Type      string      `json:"type"`      // added, modified, removed
	Field     string      `json:"field"`
	FieldPath string      `json:"fieldPath,omitempty"`
	OldValue  interface{} `json:"oldValue"`
	NewValue  interface{} `json:"newValue"`
	Icon      string      `json:"icon"`
}

// ChangeType enumerates possible change types
type ChangeType string

const (
	ChangeTypeCreate     ChangeType = "create"
	ChangeTypeUpdate     ChangeType = "update"
	ChangeTypeRollback   ChangeType = "rollback"
	ChangeTypeActivate   ChangeType = "activate"
	ChangeTypePublish    ChangeType = "publish"
	ChangeTypeDuplicate  ChangeType = "duplicate"
)

// VersionComparison represents a comparison between two versions
type VersionComparison struct {
	FromVersionNumber  int                 `json:"fromVersionNumber"`
	ToVersionNumber    int                 `json:"toVersionNumber"`
	FromSnapshot       *ThemeSnapshot      `json:"fromSnapshot"`
	ToSnapshot         *ThemeSnapshot      `json:"toSnapshot"`
	Changes            []map[string]interface{} `json:"changes"`
	ChangedFieldCount  int                 `json:"changedFieldCount"`
	DifferencePercent  float64             `json:"differencePercent"`
}

// VersionStats contains statistics about version history
type VersionStats struct {
	TotalVersions       int       `json:"totalVersions"`
	OldestVersion       *ThemeVersion `json:"oldestVersion"`
	LatestVersion       *ThemeVersion `json:"latestVersion"`
	TotalSnapshotSize   int64     `json:"totalSnapshotSize"`
	AverageSnapshotSize int64     `json:"averageSnapshotSize"`
	CreatedAt           time.Time `json:"createdAt"`
}

// VersionDiff represents detailed differences for a version
type VersionDiff struct {
	VersionNumber     int                       `json:"versionNumber"`
	PreviousVersion   int                       `json:"previousVersion"`
	Changes           []map[string]interface{} `json:"changes"`
	ChangeCount       int                       `json:"changeCount"`
	Timestamp         time.Time                 `json:"timestamp"`
}

// RollbackOptions defines options for rollback operation
type RollbackOptions struct {
	TargetVersion int    `json:"targetVersion"`
	Reason        string `json:"reason,omitempty"`
	PreviewFirst  bool   `json:"previewFirst,omitempty"`
}

// CreateSnapshotRequest defines request for creating a manual snapshot
type CreateSnapshotRequest struct {
	Summary string `json:"summary,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// CompareVersionsRequest defines request for comparing versions
type CompareVersionsRequest struct {
	FromVersion int `json:"fromVersion"`
	ToVersion   int `json:"toVersion"`
}

// ChangeTypeConfig defines display configuration for each change type
type ChangeTypeConfig struct {
	Icon  string `json:"icon"`
	Label string `json:"label"`
	Color string `json:"color"`
}

// ChangeTypeConfigs maps change types to their display config
var ChangeTypeConfigs = map[string]ChangeTypeConfig{
	"create": {
		Icon:  "✨",
		Label: "Created",
		Color: "green",
	},
	"update": {
		Icon:  "✏️",
		Label: "Updated",
		Color: "blue",
	},
	"rollback": {
		Icon:  "↩️",
		Label: "Rolled Back",
		Color: "orange",
	},
	"activate": {
		Icon:  "🔓",
		Label: "Activated",
		Color: "purple",
	},
	"publish": {
		Icon:  "📤",
		Label: "Published",
		Color: "indigo",
	},
	"duplicate": {
		Icon:  "📋",
		Label: "Duplicated",
		Color: "cyan",
	},
}
