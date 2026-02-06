package domain

import (
	"time"
)

// RestaurantTheme represents the theme configuration for a restaurant
type RestaurantTheme struct {
	ID             int64          `json:"id"`
	RestaurantID   int64          `json:"restaurant_id"`
	TenantID       int64          `json:"tenant_id"`
	PrimaryColor   string         `json:"primary_color"`   // Hex color code
	SecondaryColor string         `json:"secondary_color"` // Hex color code
	AccentColor    string         `json:"accent_color"`    // Hex color code
	LogoURL        string         `json:"logo_url"`        // URL to restaurant logo
	FontFamily     string         `json:"font_family"`     // Font family name
	Sections       []ThemeSection `json:"sections"`        // Ordered sections
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// ThemeSection represents a section in the home page
type ThemeSection struct {
	ID              int64                  `json:"id"`
	ThemeID         int64                  `json:"theme_id"`
	SectionType     string                 `json:"section_type"`     // hero, featured_items, why_choose_us, info_cards, cta, testimonials
	Order           int                    `json:"order"`            // Display order (1, 2, 3, ...)
	IsVisible       bool                   `json:"is_visible"`       // Show/hide section
	Title           string                 `json:"title"`            // Section title
	Subtitle        string                 `json:"subtitle"`         // Section subtitle
	Description     string                 `json:"description"`      // Section description
	BackgroundImage string                 `json:"background_image"` // Background image URL
	ButtonText      string                 `json:"button_text"`      // CTA button text
	ButtonLink      string                 `json:"button_link"`      // CTA button link
	Content         map[string]interface{} `json:"content"`          // JSON content specific to section type
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// ComponentType defines available component types for home page builder
type ComponentType struct {
	ID             string                 `json:"id"`              // hero, featured_items, etc.
	Name           string                 `json:"name"`            // Display name
	Description    string                 `json:"description"`     // What this component does
	Icon           string                 `json:"icon"`            // Icon name
	DefaultData    map[string]interface{} `json:"default_data"`    // Default values
	EditableFields []ComponentField       `json:"editable_fields"` // Fields user can edit
}

// ComponentField defines an editable field in a component
type ComponentField struct {
	Key             string   `json:"key"`               // Field key (title, subtitle, image, etc.)
	Label           string   `json:"label"`             // Display label
	Type            string   `json:"type"`              // text, textarea, image, color, select, etc.
	Required        bool     `json:"required"`          // Is required
	Placeholder     string   `json:"placeholder"`       // Placeholder text
	Options         []string `json:"options,omitempty"` // For select fields
	MaxLength       int      `json:"max_length,omitempty"`
	ValidationRules []string `json:"validation_rules,omitempty"` // email, url, etc.
}

// RestaurantThemeResponse is the API response for theme
type RestaurantThemeResponse struct {
	Success bool             `json:"success"`
	Data    *RestaurantTheme `json:"data,omitempty"`
	Message string           `json:"message,omitempty"`
}

// ComponentLibraryResponse is the API response for available components
type ComponentLibraryResponse struct {
	Success bool            `json:"success"`
	Data    []ComponentType `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}

// UpdateThemeRequest is the request for updating restaurant theme
type UpdateThemeRequest struct {
	PrimaryColor   string         `json:"primary_color"`
	SecondaryColor string         `json:"secondary_color"`
	AccentColor    string         `json:"accent_color"`
	LogoURL        string         `json:"logo_url"`
	FontFamily     string         `json:"font_family"`
	Sections       []ThemeSection `json:"sections"`
}

// SectionOrderRequest is the request for reordering sections
type SectionOrderRequest struct {
	SectionID int64 `json:"section_id"`
	Order     int   `json:"order"`
}
