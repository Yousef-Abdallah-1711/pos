//go:build legacy
// +build legacy

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ComponentVersion represents a semantic version of a component
type ComponentVersion struct {
	ID              int64               `db:"id" json:"id"`
	ComponentID     int64               `db:"component_id" json:"component_id"`
	Version         string              `db:"version" json:"version"`
	Major           int                 `db:"major" json:"major"`
	Minor           int                 `db:"minor" json:"minor"`
	Patch           int                 `db:"patch" json:"patch"`
	Prerelease      string              `db:"prerelease" json:"prerelease"`
	Metadata        string              `db:"metadata" json:"metadata"`
	Description     string              `db:"description" json:"description"`
	ReleaseNotes    string              `db:"release_notes" json:"release_notes"`
	IsLatest        bool                `db:"is_latest" json:"is_latest"`
	IsDeprecated    bool                `db:"is_deprecated" json:"is_deprecated"`
	DeprecationMsg  string              `db:"deprecation_message" json:"deprecation_message"`
	BreakingChanges BreakingChanges     `db:"breaking_changes" json:"breaking_changes"`
	MigrationGuides []MigrationGuide    `db:"migration_guides" json:"migration_guides"`
	Changelog       string              `db:"changelog" json:"changelog"`
	CreatedAt       time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `db:"updated_at" json:"updated_at"`
}

// BreakingChanges represents breaking changes in this version
type BreakingChanges struct {
	HasBreakingChanges bool     `json:"has_breaking_changes"`
	Changes            []string `json:"changes"`
	MigrationRequired  bool     `json:"migration_required"`
	MigrationSteps     []string `json:"migration_steps"`
}

// Scan implements sql.Scanner for PostgreSQL JSON
func (bc *BreakingChanges) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), bc)
}

// Value implements driver.Valuer for PostgreSQL JSON
func (bc BreakingChanges) Value() (driver.Value, error) {
	return json.Marshal(bc)
}

// MigrationGuide represents a guide to migrate from one version to another
type MigrationGuide struct {
	FromVersion string            `json:"from_version"`
	ToVersion   string            `json:"to_version"`
	Steps       []string          `json:"steps"`
	Warnings    []string          `json:"warnings"`
	EstimatedTime int             `json:"estimated_time_minutes"`
	AutoMigrable bool             `json:"auto_migrable"`
	AutoMigrationScript string     `json:"auto_migration_script,omitempty"`
}

// Scan implements sql.Scanner for PostgreSQL JSON array
func (guides *[]MigrationGuide) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), guides)
}

// Value implements driver.Valuer for PostgreSQL JSON array
func (guides []MigrationGuide) Value() (driver.Value, error) {
	return json.Marshal(guides)
}

// VersionConstraint represents a version constraint like ^1.2.0 or ~1.2
type VersionConstraint struct {
	Constraint string // ^, ~, >=, <=, >, <, =, !=
	Version    string // 1.2.0
}

// ParseSemver parses a semantic version string into components
// Supports formats: 1.2.3, 1.2.3-alpha, 1.2.3-alpha.1, 1.2.3+build.1, 1.2.3-alpha.1+build.1
func ParseSemver(versionStr string) (*ComponentVersion, error) {
	versionStr = strings.TrimSpace(versionStr)
	if versionStr == "" {
		return nil, errors.New("version string is empty")
	}

	cv := &ComponentVersion{
		Version: versionStr,
	}

	// Parse prerelease and metadata
	// Format: major.minor.patch[-prerelease][+metadata]
	metadataIdx := strings.Index(versionStr, "+")
	prereleaseIdx := strings.Index(versionStr, "-")

	var baseVersion string
	if metadataIdx > 0 {
		cv.Metadata = versionStr[metadataIdx+1:]
		versionStr = versionStr[:metadataIdx]
	}

	if prereleaseIdx > 0 {
		cv.Prerelease = versionStr[prereleaseIdx+1:]
		baseVersion = versionStr[:prereleaseIdx]
	} else {
		baseVersion = versionStr
	}

	// Parse major.minor.patch
	parts := strings.Split(baseVersion, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("invalid semver format: %s (expected major.minor[.patch])", versionStr)
	}

	// Parse major
	major, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || major < 0 {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}
	cv.Major = major

	// Parse minor
	minor, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || minor < 0 {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	cv.Minor = minor

	// Parse patch (optional, defaults to 0)
	if len(parts) == 3 {
		patch, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil || patch < 0 {
			return nil, fmt.Errorf("invalid patch version: %s", parts[2])
		}
		cv.Patch = patch
	} else {
		cv.Patch = 0
	}

	return cv, nil
}

// Compare compares this version with another
// Returns: -1 if this < other, 0 if this == other, 1 if this > other
func (cv *ComponentVersion) Compare(other *ComponentVersion) int {
	if cv.Major != other.Major {
		if cv.Major < other.Major {
			return -1
		}
		return 1
	}

	if cv.Minor != other.Minor {
		if cv.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if cv.Patch != other.Patch {
		if cv.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Compare prerelease (no prerelease > prerelease)
	if (cv.Prerelease == "") != (other.Prerelease == "") {
		if cv.Prerelease == "" {
			return 1
		}
		return -1
	}

	if cv.Prerelease != other.Prerelease {
		if cv.Prerelease < other.Prerelease {
			return -1
		}
		return 1
	}

	return 0
}

// IsCompatibleWith checks if this version is compatible with a constraint
func (cv *ComponentVersion) IsCompatibleWith(constraint VersionConstraint) bool {
	targetVersion, err := ParseSemver(constraint.Version)
	if err != nil {
		return false
	}

	switch constraint.Constraint {
	case "^": // Caret: compatible with version (allows minor and patch updates)
		// ^1.2.3 := >=1.2.3 <2.0.0
		// ^0.2.3 := >=0.2.3 <0.3.0
		// ^0.0.3 := >=0.0.3 <0.0.4
		if cv.Major == 0 && targetVersion.Major == 0 {
			if cv.Minor == 0 && targetVersion.Minor == 0 {
				return cv.Patch == targetVersion.Patch
			}
			return cv.Major == targetVersion.Major && cv.Minor == targetVersion.Minor && cv.Patch >= targetVersion.Patch
		}
		return cv.Major == targetVersion.Major && cv.Compare(targetVersion) >= 0

	case "~": // Tilde: approximately equivalent version (allows patch updates)
		// ~1.2.3 := >=1.2.3 <1.3.0
		// ~1.2 := >=1.2.0 <1.3.0
		return cv.Major == targetVersion.Major && cv.Minor == targetVersion.Minor && cv.Patch >= targetVersion.Patch

	case ">=":
		return cv.Compare(targetVersion) >= 0
	case ">":
		return cv.Compare(targetVersion) > 0
	case "<=":
		return cv.Compare(targetVersion) <= 0
	case "<":
		return cv.Compare(targetVersion) < 0
	case "=":
		return cv.Compare(targetVersion) == 0
	case "!=":
		return cv.Compare(targetVersion) != 0
	case "": // No constraint
		return true
	default:
		return false
	}
}

// String returns the semantic version string
func (cv *ComponentVersion) String() string {
	version := fmt.Sprintf("%d.%d.%d", cv.Major, cv.Minor, cv.Patch)
	if cv.Prerelease != "" {
		version += "-" + cv.Prerelease
	}
	if cv.Metadata != "" {
		version += "+" + cv.Metadata
	}
	return version
}

// ParseVersionConstraint parses a version constraint string like ^1.2.0
func ParseVersionConstraint(constraintStr string) (*VersionConstraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)
	if constraintStr == "" {
		return nil, errors.New("constraint string is empty")
	}

	// Extract constraint operator
	operators := []string{"^", "~", ">=", "<=", ">", "<", "=", "!="}
	var constraint string
	var versionPart string

	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			constraint = op
			versionPart = strings.TrimPrefix(constraintStr, op)
			versionPart = strings.TrimSpace(versionPart)
			break
		}
	}

	if constraint == "" {
		// Default to = if no operator specified
		constraint = "="
		versionPart = constraintStr
	}

	if versionPart == "" {
		return nil, errors.New("version part is empty")
	}

	return &VersionConstraint{
		Constraint: constraint,
		Version:    versionPart,
	}, nil
}

// ValidateSemver validates a semantic version string format
func ValidateSemver(versionStr string) error {
	_, err := ParseSemver(versionStr)
	return err
}

// ExtractVersionFromString extracts semantic version from a string like "v1.2.3" or "1.2.3"
func ExtractVersionFromString(s string) (string, error) {
	// Match semantic version pattern with optional 'v' prefix
	pattern := `v?(\d+\.\d+(?:\.\d+)?)(?:-[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*)?(?:\+[a-zA-Z0-9]+(?:\.[a-zA-Z0-9]+)*)?`
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) == 0 {
		return "", fmt.Errorf("no semantic version found in string: %s", s)
	}

	return matches[0], nil
}
