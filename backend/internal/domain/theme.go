package domain

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

// Theme represents a complete theme configuration for a restaurant/tenant
type Theme struct {
	ID           int64           `json:"id"`
	TenantID     int64           `json:"tenant_id"`
	RestaurantID int64           `json:"restaurant_id,omitempty"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	Config       json.RawMessage `json:"config"` // Entire theme as JSON - JSONB from DB
	Description  string          `json:"description,omitempty"`
	IsActive     bool            `json:"is_active"`
	IsPublished  bool            `json:"is_published"`
	Version      int             `json:"version"`
	CreatedBy    *int64          `json:"created_by,omitempty"`
	UpdatedBy    *int64          `json:"updated_by,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	PublishedAt  *time.Time      `json:"published_at,omitempty"`
}

// ThemeConfig is the structure of the config JSONB field
// This defines the complete theme configuration with all styling options
type ThemeConfig struct {
	Styling struct {
		Colors struct {
			Primary    string `json:"primary"`    // e.g. #3b82f6
			Secondary  string `json:"secondary"`  // e.g. #1e40af
			Accent     string `json:"accent"`     // e.g. #0ea5e9
			Background string `json:"background"` // e.g. #ffffff
			Text       string `json:"text"`       // e.g. #1f2937
			Border     string `json:"border"`     // e.g. #e5e7eb
			Shadow     string `json:"shadow"`     // e.g. #000000
		} `json:"colors"`
		Typography struct {
			FontSans     string  `json:"font_sans"`     // e.g. Inter, Arial
			FontSerif    string  `json:"font_serif"`    // e.g. Georgia
			FontMono     string  `json:"font_mono"`     // e.g. Courier
			BaseSize     int     `json:"base_size"`     // e.g. 16
			LineHeight   float64 `json:"line_height"`   // e.g. 1.5
			BorderRadius int     `json:"border_radius"` // e.g. 8
		} `json:"typography"`
	} `json:"styling"`

	Header struct {
		LogoURL   string `json:"logo_url"`
		NavBg     string `json:"nav_bg"`
		NavText   string `json:"nav_text"`
		NavHeight int    `json:"nav_height"`
		Sticky    bool   `json:"sticky"`
	} `json:"header"`

	Footer struct {
		Text      string `json:"text"`
		BgColor   string `json:"bg_color"`
		TextColor string `json:"text_color"`
		Links     []struct {
			Label string `json:"label"`
			URL   string `json:"url"`
		} `json:"links"`
	} `json:"footer"`

	HomePage struct {
		Sections []struct {
			ID      string                 `json:"id"`
			Type    string                 `json:"type"` // hero, products, testimonials, etc
			Order   int                    `json:"order"`
			Enabled bool                   `json:"enabled"`
			Config  map[string]interface{} `json:"config"`
		} `json:"sections"`
	} `json:"homepage"`

	DarkMode struct {
		Enabled        bool   `json:"enabled"`
		PrimaryDark    string `json:"primary_dark"`
		SecondaryDark  string `json:"secondary_dark"`
		BackgroundDark string `json:"background_dark"`
		TextDark       string `json:"text_dark"`
	} `json:"dark_mode"`

	RTL struct {
		Enabled bool `json:"enabled"`
	} `json:"rtl"`

	Custom map[string]interface{} `json:"custom"` // For extensibility
}

// Validation error type
type ThemeValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ThemeValidationError) Error() string {
	if e.Field != "" {
		return "theme validation error: " + e.Field + " - " + e.Message
	}
	return "theme validation error: " + e.Message
}

// Validate checks if theme data is valid
func (t *Theme) Validate() error {
	// Check name
	if t.Name == "" {
		return &ThemeValidationError{Field: "name", Message: "name is required"}
	}
	if len(t.Name) < 3 {
		return &ThemeValidationError{Field: "name", Message: "name must be at least 3 characters"}
	}
	if len(t.Name) > 255 {
		return &ThemeValidationError{Field: "name", Message: "name must not exceed 255 characters"}
	}

	// Check slug
	if t.Slug == "" {
		return &ThemeValidationError{Field: "slug", Message: "slug is required"}
	}
	if !isValidSlug(t.Slug) {
		return &ThemeValidationError{Field: "slug", Message: "slug must be lowercase alphanumeric with hyphens, between 3-50 characters"}
	}

	// Check tenant_id
	if t.TenantID <= 0 {
		return &ThemeValidationError{Field: "tenant_id", Message: "tenant_id must be positive"}
	}

	// Check config (if provided, should be valid JSON)
	if len(t.Config) > 0 {
		// First, check if it's valid JSON
		var rawConfig map[string]interface{}
		if err := json.Unmarshal(t.Config, &rawConfig); err != nil {
			return &ThemeValidationError{Field: "config", Message: "config must be valid JSON: " + err.Error()}
		}

		// Try to parse as structured ThemeConfig for validation
		// but don't fail if it's a flat structure (frontend sends flat)
		var cfg ThemeConfig
		if err := json.Unmarshal(t.Config, &cfg); err == nil {
			// Only validate colors if config was parsed successfully
			if err := validateThemeColors(&cfg); err != nil {
				return err
			}
		}
		// If config doesn't match ThemeConfig structure, that's OK -
		// it means frontend sent a flat structure which is also valid
	}

	return nil
}

// isValidSlug validates slug format
func isValidSlug(slug string) bool {
	if len(slug) < 3 || len(slug) > 50 {
		return false
	}
	// Match: lowercase letters, numbers, hyphens only
	match, err := regexp.MatchString(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`, slug)
	if err != nil {
		return false
	}
	return match
}

// validateThemeColors validates color hex codes
func validateThemeColors(cfg *ThemeConfig) error {
	hexRegex := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	colors := []struct {
		name  string
		value string
	}{
		{"primary", cfg.Styling.Colors.Primary},
		{"secondary", cfg.Styling.Colors.Secondary},
		{"accent", cfg.Styling.Colors.Accent},
		{"background", cfg.Styling.Colors.Background},
		{"text", cfg.Styling.Colors.Text},
		{"border", cfg.Styling.Colors.Border},
		{"shadow", cfg.Styling.Colors.Shadow},
	}

	for _, color := range colors {
		if color.value != "" && !hexRegex.MatchString(color.value) {
			return &ThemeValidationError{Field: "config.styling.colors." + color.name, Message: "invalid hex color format"}
		}
	}

	return nil
}

// ParseConfig parses the JSON config into ThemeConfig struct
func (t *Theme) ParseConfig() (*ThemeConfig, error) {
	if len(t.Config) == 0 {
		return &ThemeConfig{}, nil
	}

	var cfg ThemeConfig
	if err := json.Unmarshal(t.Config, &cfg); err != nil {
		return nil, errors.New("failed to parse theme config: " + err.Error())
	}

	return &cfg, nil
}

// SetConfig sets the config from a ThemeConfig struct
func (t *Theme) SetConfig(cfg *ThemeConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.New("failed to marshal theme config: " + err.Error())
	}

	t.Config = data
	return nil
}

// NewTheme creates a new theme with validation
func NewTheme(name, slug string, tenantID int64) (*Theme, error) {
	theme := &Theme{
		Name:        name,
		Slug:        slug,
		TenantID:    tenantID,
		IsActive:    false,
		IsPublished: false,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := theme.Validate(); err != nil {
		return nil, err
	}

	return theme, nil
}

// ThemeCreateRequest represents a request to create a new theme
type ThemeCreateRequest struct {
	Name        string          `json:"name" binding:"required"`
	Slug        string          `json:"slug" binding:"required"`
	Description string          `json:"description"`
	Config      json.RawMessage `json:"config"`
}

// ThemeUpdateRequest represents a request to update a theme
type ThemeUpdateRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Config      json.RawMessage `json:"config"`
	IsActive    *bool           `json:"is_active"`
	IsPublished *bool           `json:"is_published"`
}

// ThemeListResponse represents a paginated list of themes
type ThemeListResponse struct {
	Themes     []*Theme `json:"themes"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}
