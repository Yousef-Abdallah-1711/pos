package domain

import (
	"encoding/json"
	"testing"
	"time"
)

// TestThemeValidation tests basic theme validation
func TestThemeValidation(t *testing.T) {
	tests := []struct {
		name    string
		theme   *Theme
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid theme",
			theme: &Theme{
				Name:      "Modern Theme",
				Slug:      "modern-theme",
				TenantID:  1,
				Config:    json.RawMessage("{}"),
				IsActive:  false,
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing name",
			theme: &Theme{
				Name:      "",
				Slug:      "modern-theme",
				TenantID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too short",
			theme: &Theme{
				Name:      "Mo",
				Slug:      "modern-theme",
				TenantID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name must be at least 3 characters",
		},
		{
			name: "name too long",
			theme: &Theme{
				Name:      "A" + string(make([]byte, 255)),
				Slug:      "modern-theme",
				TenantID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name must not exceed 255 characters",
		},
		{
			name: "missing slug",
			theme: &Theme{
				Name:      "Modern Theme",
				Slug:      "",
				TenantID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "slug is required",
		},
		{
			name: "invalid slug format",
			theme: &Theme{
				Name:      "Modern Theme",
				Slug:      "Modern Theme",
				TenantID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "slug must be lowercase alphanumeric",
		},
		{
			name: "invalid tenant id",
			theme: &Theme{
				Name:      "Modern Theme",
				Slug:      "modern-theme",
				TenantID:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "tenant_id must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.theme.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				if errStr := err.Error(); tt.errMsg != "" && !contains(errStr, tt.errMsg) {
					t.Errorf("Validate() error = %v, want message containing %v", errStr, tt.errMsg)
				}
			}
		})
	}
}

// TestSlugValidation tests slug format validation
func TestSlugValidation(t *testing.T) {
	tests := []struct {
		name  string
		slug  string
		valid bool
	}{
		{"valid slug", "modern-theme", true},
		{"valid slug with numbers", "theme-v2-premium", true},
		{"single char", "a", true},
		{"slug too short", "ab", false},
		{"slug too long", "a" + string(make([]byte, 49)), false},
		{"uppercase", "Modern-Theme", false},
		{"spaces", "modern theme", false},
		{"special chars", "modern@theme", false},
		{"ends with hyphen", "modern-", false},
		{"starts with hyphen", "-modern", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidSlug(tt.slug)
			if valid != tt.valid {
				t.Errorf("isValidSlug(%q) = %v, want %v", tt.slug, valid, tt.valid)
			}
		})
	}
}

// TestColorValidation tests hex color validation
func TestColorValidation(t *testing.T) {
	// Create valid config
	validCfg := &ThemeConfig{}
	validCfg.Styling.Colors.Primary = "#3b82f6"
	validCfg.Styling.Colors.Secondary = "#1e40af"
	validCfg.Styling.Colors.Accent = "#0ea5e9"
	validCfg.Styling.Colors.Background = "#ffffff"
	validCfg.Styling.Colors.Text = "#1f2937"
	validCfg.Styling.Colors.Border = "#e5e7eb"
	validCfg.Styling.Colors.Shadow = "#000000"

	// Create invalid config
	invalidCfg := &ThemeConfig{}
	invalidCfg.Styling.Colors.Primary = "rgb(255,0,0)"

	tests := []struct {
		name  string
		cfg   *ThemeConfig
		valid bool
	}{
		{"valid colors", validCfg, true},
		{"invalid color format", invalidCfg, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateThemeColors(tt.cfg)
			if (err != nil) != !tt.valid {
				t.Errorf("validateThemeColors() error = %v, want valid = %v", err, tt.valid)
			}
		})
	}
}

// TestParseConfig tests config parsing
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  json.RawMessage
		wantErr bool
	}{
		{
			name:    "empty config",
			config:  json.RawMessage(""),
			wantErr: false,
		},
		{
			name:    "valid config",
			config:  json.RawMessage(`{"styling":{"colors":{"primary":"#3b82f6"}}}`),
			wantErr: false,
		},
		{
			name:    "invalid json",
			config:  json.RawMessage(`{invalid}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := &Theme{Config: tt.config}
			cfg, err := theme.ParseConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && cfg == nil {
				t.Errorf("ParseConfig() returned nil config for valid input")
			}
		})
	}
}

// TestSetConfig tests config setting
func TestSetConfig(t *testing.T) {
	cfg := &ThemeConfig{}
	theme := &Theme{}

	err := theme.SetConfig(cfg)
	if err != nil {
		t.Errorf("SetConfig() error = %v, want nil", err)
	}

	if len(theme.Config) == 0 {
		t.Errorf("SetConfig() did not set config")
	}

	// Verify we can parse it back
	parsedCfg, err := theme.ParseConfig()
	if err != nil {
		t.Errorf("ParseConfig() after SetConfig() error = %v", err)
	}
	if parsedCfg == nil {
		t.Errorf("ParseConfig() after SetConfig() returned nil")
	}
}

// TestNewTheme tests theme constructor
func TestNewTheme(t *testing.T) {
	tests := []struct {
		name      string
		themeName string
		slug      string
		tenantID  int64
		wantErr   bool
	}{
		{"valid", "Modern Theme", "modern-theme", 1, false},
		{"missing name", "", "modern-theme", 1, true},
		{"invalid slug", "Modern Theme", "Modern Theme", 1, true},
		{"missing tenant", "Modern Theme", "modern-theme", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := NewTheme(tt.themeName, tt.slug, tt.tenantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTheme() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if theme == nil {
					t.Errorf("NewTheme() returned nil for valid input")
				}
				if theme.Name != tt.themeName || theme.Slug != tt.slug {
					t.Errorf("NewTheme() fields not set correctly")
				}
				if theme.IsActive != false || theme.IsPublished != false {
					t.Errorf("NewTheme() initial flags incorrect")
				}
			}
		})
	}
}

// TestThemeListResponse tests response structure
func TestThemeListResponse(t *testing.T) {
	response := &ThemeListResponse{
		Themes:     make([]*Theme, 0),
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	}

	if response.TotalPages != (int(response.Total)+response.PageSize-1)/response.PageSize {
		t.Errorf("TotalPages calculation incorrect")
	}
}

// TestThemeCreateRequest tests create request
func TestThemeCreateRequest(t *testing.T) {
	req := &ThemeCreateRequest{
		Name:        "Test Theme",
		Slug:        "test-theme",
		Description: "A test theme",
		Config:      json.RawMessage(`{}`),
	}

	if req.Name == "" || req.Slug == "" {
		t.Errorf("ThemeCreateRequest fields empty")
	}
}

// TestThemeUpdateRequest tests update request
func TestThemeUpdateRequest(t *testing.T) {
	isActive := true
	req := &ThemeUpdateRequest{
		Name:     "Updated Theme",
		IsActive: &isActive,
	}

	if req.Name != "Updated Theme" || *req.IsActive != true {
		t.Errorf("ThemeUpdateRequest fields incorrect")
	}
}

// TestThemeValidationError tests error type
func TestThemeValidationError(t *testing.T) {
	err := &ThemeValidationError{
		Field:   "name",
		Message: "name is required",
	}

	errStr := err.Error()
	if !contains(errStr, "name") || !contains(errStr, "required") {
		t.Errorf("ThemeValidationError.Error() = %v, want to contain 'name' and 'required'", errStr)
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
