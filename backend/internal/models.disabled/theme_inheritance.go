//go:build legacy
// +build legacy

package models

import (
	"database/sql"
	"time"
)

// ThemeInheritance represents the parent-child relationship between themes
type ThemeInheritance struct {
	ID            string         `json:"id" db:"id"`
	ThemeID       string         `json:"themeId" db:"theme_id"`           // Child theme
	ParentThemeID string         `json:"parentThemeId" db:"parent_theme_id"` // Parent theme
	Overrides     sql.NullString `json:"overrides" db:"overrides"`        // JSON object of overridden properties
	InheritLevel  int            `json:"inheritLevel" db:"inherit_level"` // 0=direct, 1=grandparent, etc.
	IsActive      bool           `json:"isActive" db:"is_active"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
}

// ThemePropertyOverride tracks individual property overrides
type ThemePropertyOverride struct {
	ID            string         `json:"id" db:"id"`
	ThemeID       string         `json:"themeId" db:"theme_id"`
	PropertyPath  string         `json:"propertyPath" db:"property_path"` // e.g., "colors.primary", "typography.fontSize"
	OriginalValue sql.NullString `json:"originalValue" db:"original_value"` // Value from parent
	OverrideValue sql.NullString `json:"overrideValue" db:"override_value"` // Overridden value
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
}

// ResolvedTheme represents a theme with all inherited properties resolved
type ResolvedTheme struct {
	ID              string                      `json:"id"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description"`
	Slug            string                      `json:"slug"`
	SiteTitle       string                      `json:"siteTitle"`
	PrimaryColor    string                      `json:"primaryColor"`
	SecondaryColor  string                      `json:"secondaryColor"`
	AccentColor     string                      `json:"accentColor"`
	BackgroundColor string                      `json:"backgroundColor"`
	TextColor       string                      `json:"textColor"`
	BorderColor     string                      `json:"borderColor"`
	ShadowColor     string                      `json:"shadowColor"`
	FontFamily      string                      `json:"fontFamily"`
	BaseFontSize    int                         `json:"baseFontSize"`
	BorderRadius    int                         `json:"borderRadius"`
	LineHeight      float64                     `json:"lineHeight"`
	LogoURL         string                      `json:"logoUrl"`
	FaviconURL      string                      `json:"faviconUrl"`
	IsPublished     bool                        `json:"isPublished"`
	IsActive        bool                        `json:"isActive"`
	InheritanceChain []*ThemeInheritanceNode    `json:"inheritanceChain"` // Chain of parent themes
	CreatedAt       time.Time                   `json:"createdAt"`
	UpdatedAt       time.Time                   `json:"updatedAt"`
}

// ThemeInheritanceNode represents a node in the inheritance chain
type ThemeInheritanceNode struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Level           int                    `json:"level"`
	Overrides       map[string]interface{} `json:"overrides"`
}
