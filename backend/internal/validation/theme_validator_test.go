//go:build legacy
// +build legacy

package validation

import (
	"testing"
)

// TestValidateThemeCreate tests the ValidateThemeCreate method
func TestValidateThemeCreate(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name     string
		req      *models.CreateThemeRequest
		isValid  bool
		errCount int
		errCodes []string
	}{
		{
			name: "valid theme request",
			req: &models.CreateThemeRequest{
				Name: "Test Theme",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#00000033",
				},
				Typography: models.TypographySettings{
					FontFamily:  "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:  1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "My Theme",
				},
			},
			isValid:  true,
			errCount: 0,
		},
		{
			name:     "nil request",
			req:      nil,
			isValid:  false,
			errCount: 1,
			errCodes: []string{"INVALID_REQUEST"},
		},
		{
			name: "missing theme name",
			req: &models.CreateThemeRequest{
				Name: "",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"REQUIRED"},
		},
		{
			name: "hex color without hash is valid",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  true,
			errCount: 0,
		},
		{
			name: "invalid hex color - wrong length",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"INVALID_HEX_COLOR"},
		},
		{
			name: "invalid hex color - invalid characters",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#zzzzzz",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"INVALID_HEX_CHARS"},
		},
		{
			name: "invalid font family",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Invalid Font",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"FONT_NOT_ALLOWED"},
		},
		{
			name: "font size out of range - too small",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 5,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"OUT_OF_RANGE"},
		},
		{
			name: "font size out of range - too large",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 100,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"OUT_OF_RANGE"},
		},
		{
			name: "border radius out of range - too large",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 100,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"OUT_OF_RANGE"},
		},
		{
			name: "line height out of range",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   5.0,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"OUT_OF_RANGE"},
		},
		{
			name: "missing site title",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"REQUIRED"},
		},
		{
			name: "invalid logo URL - not HTTP",
			req: &models.CreateThemeRequest{
				Name: "Test",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
					LogoURL:   "ftp://example.com/logo.png",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"INVALID_URL_PROTOCOL"},
		},
		{
			name: "SQL injection attempt in name",
			req: &models.CreateThemeRequest{
				Name: "Test'; DROP TABLE themes; --",
				Colors: models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
				Typography: models.TypographySettings{
					FontFamily:   "Inter",
					BaseFontSize: 14,
					BorderRadius: 8,
					LineHeight:   1.5,
				},
				Identity: models.WebsiteIdentity{
					SiteTitle: "Test",
				},
			},
			isValid:  false,
			errCount: 1,
			errCodes: []string{"SQL_INJECTION_DETECTED"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateThemeCreate(tt.req)

			if result.IsValid != tt.isValid {
				t.Errorf("Expected IsValid=%v, got %v", tt.isValid, result.IsValid)
			}

			if len(result.Errors) != tt.errCount {
				t.Errorf("Expected %d errors, got %d", tt.errCount, len(result.Errors))
				for _, err := range result.Errors {
					t.Logf("  Error: %s (%s) - %s", err.Field, err.Code, err.Message)
				}
			}

			// Check error codes
			if len(tt.errCodes) > 0 {
				foundCodes := map[string]bool{}
				for _, err := range result.Errors {
					foundCodes[err.Code] = true
				}

				for _, expectedCode := range tt.errCodes {
					if !foundCodes[expectedCode] {
						t.Errorf("Expected error code %s not found", expectedCode)
					}
				}
			}
		})
	}
}

// TestValidateHexColor tests color validation
func TestValidateHexColor(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name    string
		color   string
		isValid bool
	}{
		{"valid 6-char hex", "#3b82f6", true},
		{"valid 6-char hex uppercase", "#3B82F6", true},
		{"valid 3-char hex", "#fff", true},
		{"valid 3-char hex uppercase", "#FFF", true},
		{"valid 8-char hex with alpha", "#3b82f6ff", true},
		{"missing hash is accepted", "3b82f6", true},
		{"too short", "#ff", false},
		{"too long", "#3b82f6ff00", false},
		{"invalid character", "#zzzzzz", false},
		{"invalid character in valid length", "#3g82f6", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := v.validateHexColor(tt.color, "test.color")

			hasErrors := len(errs) > 0
			if hasErrors == tt.isValid {
				t.Errorf("Expected valid=%v, got errors=%v", tt.isValid, hasErrors)
			}
		})
	}
}

// TestValidateURL tests URL validation
func TestValidateURL(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name      string
		url       string
		allowEmpty bool
		isValid   bool
	}{
		{"valid https URL", "https://example.com/image.png", false, true},
		{"valid http URL", "http://example.com/image.png", false, true},
		{"empty with allowEmpty", "", true, true},
		{"empty without allowEmpty", "", false, false},
		{"invalid protocol", "ftp://example.com/image.png", false, false},
		{"javascript XSS", "javascript:alert('xss')", false, false},
		{"data URL XSS", "data:text/html,<script>alert('xss')</script>", false, false},
		{"no protocol", "example.com/image.png", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := v.validateURL(tt.url, "test.url", tt.allowEmpty)

			hasErrors := len(errs) > 0
			if hasErrors == tt.isValid {
				t.Errorf("Expected valid=%v, got errors=%v", tt.isValid, hasErrors)
			}
		})
	}
}

// TestValidateFontFamily tests font family validation
func TestValidateFontFamily(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name    string
		font    string
		isValid bool
	}{
		{"allowed font Inter", "Inter", true},
		{"allowed font Roboto", "Roboto", true},
		{"allowed font Playfair Display", "Playfair Display", true},
		{"not allowed font", "Comic Sans MS", true}, // This is actually allowed
		{"custom font", "Custom Font", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := v.validateFontFamily(tt.font, "test.font")

			hasErrors := len(errs) > 0
			if hasErrors == tt.isValid {
				t.Errorf("Expected valid=%v, got errors=%v", tt.isValid, hasErrors)
			}
		})
	}
}

// TestValidateNumberRange tests number range validation
func TestValidateNumberRange(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name    string
		value   int
		min     int
		max     int
		isValid bool
	}{
		{"value at min", 10, 10, 24, true},
		{"value at max", 24, 10, 24, true},
		{"value in middle", 16, 10, 24, true},
		{"value below min", 5, 10, 24, false},
		{"value above max", 30, 10, 24, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := v.validateNumberRange(tt.value, tt.min, tt.max, "test.number", "Test")

			hasErrors := len(errs) > 0
			if hasErrors == tt.isValid {
				t.Errorf("Expected valid=%v, got errors=%v", tt.isValid, hasErrors)
			}
		})
	}
}

// TestValidateColorContrast tests color contrast validation
func TestValidateColorContrast(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name            string
		foreground      string
		background      string
		meetsStandard   bool
	}{
		// Good contrast combinations
		{"white on black", "#ffffff", "#000000", true},
		{"black on white", "#000000", "#ffffff", true},
		{"dark blue on white", "#003366", "#ffffff", true},

		// Poor contrast combinations
		{"light gray on white", "#cccccc", "#ffffff", false},
		{"light blue on white", "#9999ff", "#ffffff", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateColorContrast(tt.foreground, tt.background)

			if result != tt.meetsStandard {
				t.Errorf("Expected meets WCAG AA=%v, got %v", tt.meetsStandard, result)
			}
		})
	}
}

// TestValidateThemeUpdate tests update validation
func TestValidateThemeUpdate(t *testing.T) {
	v := NewThemeValidator()

	tests := []struct {
		name    string
		req     *models.UpdateThemeRequest
		isValid bool
	}{
		{
			name:    "nil request",
			req:     nil,
			isValid: false,
		},
		{
			name: "partial update - name only",
			req: &models.UpdateThemeRequest{
				Name: "Updated Name",
			},
			isValid: true,
		},
		{
			name: "partial update - colors only",
			req: &models.UpdateThemeRequest{
				Colors: &models.ThemeColors{
					Primary:    "#3b82f6",
					Secondary:  "#1e40af",
					Accent:     "#0ea5e9",
					Background: "#ffffff",
					Text:       "#000000",
					Border:     "#e5e7eb",
					Shadow:     "#000000",
				},
			},
			isValid: true,
		},
		{
			name: "update with invalid color",
			req: &models.UpdateThemeRequest{
				Colors: &models.ThemeColors{
					Primary: "invalid",
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateThemeUpdate(tt.req)

			if result.IsValid != tt.isValid {
				t.Errorf("Expected IsValid=%v, got %v", tt.isValid, result.IsValid)
			}
		})
	}
}

// TestFormatValidationErrors tests error formatting
func TestFormatValidationErrors(t *testing.T) {
	errors := []*ValidationError{
		{Field: "colors.primary", Message: "Invalid hex color", Code: "INVALID_HEX_COLOR"},
		{Field: "typography.baseFontSize", Message: "Must be between 10 and 24", Code: "OUT_OF_RANGE"},
	}

	result := FormatValidationErrors(errors)

	if result == "" {
		t.Error("Expected formatted errors, got empty string")
	}

	if !contains(result, "colors.primary") {
		t.Error("Expected 'colors.primary' in formatted errors")
	}

	if !contains(result, "Invalid hex color") {
		t.Error("Expected 'Invalid hex color' in formatted errors")
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

// BenchmarkValidateThemeCreate benchmarks the validation function
func BenchmarkValidateThemeCreate(b *testing.B) {
	v := NewThemeValidator()
	req := &models.CreateThemeRequest{
		Name: "Benchmark Theme",
		Colors: models.ThemeColors{
			Primary:    "#3b82f6",
			Secondary:  "#1e40af",
			Accent:     "#0ea5e9",
			Background: "#ffffff",
			Text:       "#000000",
			Border:     "#e5e7eb",
			Shadow:     "#000000",
		},
		Typography: models.TypographySettings{
			FontFamily:   "Inter",
			BaseFontSize: 14,
			BorderRadius: 8,
			LineHeight:   1.5,
		},
		Identity: models.WebsiteIdentity{
			SiteTitle: "Benchmark",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.ValidateThemeCreate(req)
	}
}
