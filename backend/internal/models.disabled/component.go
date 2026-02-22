//go:build legacy
// +build legacy

package models

import (
	"database/sql"
	"time"
)

// ComponentRegistry represents a registered component that can be used in themes
type ComponentRegistry struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	ComponentID string         `json:"componentId" db:"component_id"`     // e.g., "hero", "products", "cta"
	Version     string         `json:"version" db:"version"`               // Semantic version (1.0.0)
	Aliases     sql.NullString `json:"aliases" db:"aliases"`               // JSON array of aliases
	Category    string         `json:"category" db:"category"`             // e.g., "sections", "headers", "footers"
	Tags        sql.NullString `json:"tags" db:"tags"`                     // JSON array of tags
	Schema      sql.NullString `json:"schema" db:"schema"`                 // JSON schema for component props
	ThumbnailURL sql.NullString `json:"thumbnailUrl" db:"thumbnail_url"`  // Component preview thumbnail
	IsActive    bool           `json:"isActive" db:"is_active"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

// ComponentThemeMapping represents which components are used in a specific theme
type ComponentThemeMapping struct {
	ID          string    `json:"id" db:"id"`
	ThemeID     string    `json:"themeId" db:"theme_id"`
	ComponentID string    `json:"componentId" db:"component_id"`
	Position    int       `json:"position" db:"position"`           // Display order
	IsEnabled   bool      `json:"isEnabled" db:"is_enabled"`
	Config      sql.NullString `json:"config" db:"config"`          // JSON component configuration
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// ComponentVersion represents a version history entry for a component
type ComponentVersion struct {
	ID          string         `json:"id" db:"id"`
	ComponentID string         `json:"componentId" db:"component_id"`
	Version     string         `json:"version" db:"version"`
	Changes     sql.NullString `json:"changes" db:"changes"`  // JSON changelog
	IsActive    bool           `json:"isActive" db:"is_active"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
}
