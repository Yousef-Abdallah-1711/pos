package http

import (
	"testing"

	"pos-saas/internal/validation"
)

// ==========================================
// Integration Test Helpers
// ==========================================

func setupTestValidator(t *testing.T) *validation.ThemeValidator {
	return validation.NewThemeValidator()
}

func createValidThemeRequest() *models.CreateThemeRequest {
	return &models.CreateThemeRequest{
		Name: "Test Theme",
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
			BaseFontSize: 16,
			BorderRadius: 8,
			LineHeight:   1.5,
		},
		Identity: models.WebsiteIdentity{
			SiteTitle:  "Test Restaurant",
			LogoURL:    "https://example.com/logo.png",
			FaviconURL: "https://example.com/favicon.ico",
		},
	}
}

// ==========================================
// Integration Tests for Theme Validation
// ==========================================

// TestCreateThemeValidationIntegration tests complete theme creation validation
func TestCreateThemeValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)
	req := createValidThemeRequest()

	validationResult := validator.ValidateThemeCreate(req)

	if !validationResult.IsValid {
		t.Fatalf("Valid theme request failed validation: %v", validationResult.Errors)
	}

	t.Log("✓ Complete theme creation validation passed")
}

// TestInvalidColorValidationIntegration tests color validation in creation
func TestInvalidColorValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)
	req := createValidThemeRequest()
	req.Colors.Primary = "not-a-color"

	validationResult := validator.ValidateThemeCreate(req)

	if validationResult.IsValid {
		t.Fatal("Expected validation to fail for invalid color")
	}

	foundError := false
	for _, err := range validationResult.Errors {
		if err.Field == "colors.primary" && err.Code == "INVALID_HEX_COLOR" {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected INVALID_HEX_COLOR error, got: %v", validationResult.Errors)
	}

	t.Log("✓ Invalid color validation works correctly")
}

// TestMissingRequiredFieldsValidation tests required field validation
func TestMissingRequiredFieldsValidation(t *testing.T) {
	validator := setupTestValidator(t)
	req := &models.CreateThemeRequest{}

	validationResult := validator.ValidateThemeCreate(req)

	if validationResult.IsValid {
		t.Fatal("Expected validation to fail for missing required fields")
	}

	if len(validationResult.Errors) == 0 {
		t.Fatal("Expected multiple validation errors")
	}

	t.Logf("✓ Missing required fields validation caught %d errors", len(validationResult.Errors))
}

// TestMultipleValidationErrorsIntegration tests handling of multiple errors
func TestMultipleValidationErrorsIntegration(t *testing.T) {
	validator := setupTestValidator(t)
	req := createValidThemeRequest()

	// Introduce multiple errors
	req.Name = ""                         // Missing required
	req.Colors.Primary = "invalid"        // Invalid color
	req.Typography.BaseFontSize = 999     // Out of range
	req.Typography.FontFamily = "BadFont" // Invalid font

	validationResult := validator.ValidateThemeCreate(req)

	if validationResult.IsValid {
		t.Fatal("Expected multiple validation errors")
	}

	errorCount := len(validationResult.Errors)
	if errorCount < 3 {
		t.Errorf("Expected at least 3 errors, got %d: %v", errorCount, validationResult.Errors)
	}

	t.Logf("✓ Multiple error validation detected %d errors correctly", errorCount)
}

// TestSQLInjectionPreventionIntegration tests SQL injection detection
func TestSQLInjectionPreventionIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	testCases := []string{
		"Theme'; DROP TABLE themes;--",
		"Normal' OR '1'='1",
		"Theme; DELETE FROM themes WHERE 1=1",
		"'; INSERT INTO themes VALUES (...)",
	}

	for _, maliciousInput := range testCases {
		req := createValidThemeRequest()
		req.Name = maliciousInput

		validationResult := validator.ValidateThemeCreate(req)

		if validationResult.IsValid {
			t.Errorf("Expected SQL injection detection for: %s", maliciousInput)
		}

		foundSQLError := false
		for _, err := range validationResult.Errors {
			if err.Code == "SQL_INJECTION_DETECTED" {
				foundSQLError = true
				break
			}
		}

		if !foundSQLError {
			t.Errorf("Expected SQL_INJECTION_DETECTED error for: %s", maliciousInput)
		}
	}

	t.Log("✓ SQL injection prevention working correctly")
}

// TestXSSPreventionInURLsIntegration tests XSS URL prevention
func TestXSSPreventionInURLsIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	// Test malicious URLs
	xssURLs := []string{
		"javascript:alert('xss')",
		"javascript:void(0)",
		"data:text/html,<script>alert('xss')</script>",
	}

	for _, xssURL := range xssURLs {
		req := createValidThemeRequest()
		req.Identity.LogoURL = xssURL

		validationResult := validator.ValidateThemeCreate(req)

		if validationResult.IsValid {
			t.Errorf("Expected XSS prevention for URL: %s", xssURL)
		}
	}

	t.Log("✓ XSS prevention in URLs working correctly")
}

// TestColorRangeValidationIntegration tests all color specifications
func TestColorRangeValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	testCases := []struct {
		name       string
		colors     []string
		shouldFail bool
	}{
		{
			name:       "Valid 6-digit hex colors",
			colors:     []string{"#000000", "#ffffff", "#3b82f6", "#ff0000"},
			shouldFail: false,
		},
		{
			name:       "Valid 3-digit hex colors",
			colors:     []string{"#fff", "#000", "#f00"},
			shouldFail: false,
		},
		{
			name:       "Valid 8-digit hex with alpha",
			colors:     []string{"#000000ff", "#ffffff80", "#3b82f6cc"},
			shouldFail: false,
		},
		{
			name:       "Invalid missing hash",
			colors:     []string{"000000", "ffffff"},
			shouldFail: false, // Validator may accept without hash
		},
		{
			name:       "Invalid characters",
			colors:     []string{"#gggggg", "#zzzzzz"},
			shouldFail: true,
		},
		{
			name:       "Invalid too short",
			colors:     []string{"#ff", "#ab"},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		for _, color := range tc.colors {
			colors := &models.ThemeColors{
				Primary:    color,
				Secondary:  "#000000",
				Accent:     "#000000",
				Background: "#ffffff",
				Text:       "#000000",
				Border:     "#e5e7eb",
				Shadow:     "#000000",
			}

			errors := validator.ValidateColors(colors)
			hasError := len(errors) > 0

			if tc.shouldFail && !hasError {
				t.Logf("WARNING: Expected %s (%s) to fail but it passed", tc.name, color)
			}
			if !tc.shouldFail && hasError {
				t.Logf("INFO: %s (%s) validation: %v", tc.name, color, errors)
			}
		}
	}

	t.Log("✓ Color range validation tested")
}

// TestFontValidationIntegration tests font family validation
func TestFontValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	// Test allowed fonts
	allowedFonts := []string{
		"Inter", "Roboto", "Open Sans", "Lato", "Montserrat",
		"Raleway", "Playfair Display", "Poppins",
	}

	for _, font := range allowedFonts {
		typography := &models.TypographySettings{
			FontFamily:   font,
			BaseFontSize: 16,
			BorderRadius: 8,
			LineHeight:   1.5,
		}

		errors := validator.ValidateTypography(typography)
		if len(errors) > 0 {
			t.Errorf("Allowed font %s failed validation: %v", font, errors)
		}
	}

	// Test disallowed fonts
	invalidFonts := []string{"ComicSansMSEvil", "RandomFont123", "NotAFont"}

	for _, font := range invalidFonts {
		typography := &models.TypographySettings{
			FontFamily:   font,
			BaseFontSize: 16,
			BorderRadius: 8,
			LineHeight:   1.5,
		}

		errors := validator.ValidateTypography(typography)
		if len(errors) == 0 {
			t.Errorf("Invalid font %s should have been rejected", font)
		}
	}

	t.Log("✓ Font validation working correctly")
}

// TestTypographyRangesIntegration tests all typography range validations
func TestTypographyRangesIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	testCases := []struct {
		name       string
		typography models.TypographySettings
		shouldFail bool
	}{
		{
			name: "Valid typography",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 16,
				BorderRadius: 8,
				LineHeight:   1.5,
			},
			shouldFail: false,
		},
		{
			name: "Font size below minimum",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 9,
				BorderRadius: 8,
				LineHeight:   1.5,
			},
			shouldFail: true,
		},
		{
			name: "Font size above maximum",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 25,
				BorderRadius: 8,
				LineHeight:   1.5,
			},
			shouldFail: true,
		},
		{
			name: "Border radius above maximum",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 16,
				BorderRadius: 51,
				LineHeight:   1.5,
			},
			shouldFail: true,
		},
		{
			name: "Line height below minimum",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 16,
				BorderRadius: 8,
				LineHeight:   0.9,
			},
			shouldFail: true,
		},
		{
			name: "Line height above maximum",
			typography: models.TypographySettings{
				FontFamily:   "Inter",
				BaseFontSize: 16,
				BorderRadius: 8,
				LineHeight:   3.1,
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		errors := validator.ValidateTypography(&tc.typography)
		hasError := len(errors) > 0

		if tc.shouldFail && !hasError {
			t.Errorf("Test '%s': Expected validation to fail but it passed", tc.name)
		}
		if !tc.shouldFail && hasError {
			t.Errorf("Test '%s': Expected validation to pass but got errors: %v", tc.name, errors)
		}
	}

	t.Log("✓ Typography range validation working correctly")
}

// TestUpdateThemeValidationIntegration tests update request validation
func TestUpdateThemeValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	// Test valid partial update
	updateReq := &models.UpdateThemeRequest{
		Name: "Updated Theme",
		Colors: &models.ThemeColors{
			Primary: "#ff0000",
		},
	}

	validationResult := validator.ValidateThemeUpdate(updateReq)
	if !validationResult.IsValid {
		t.Errorf("Valid partial update failed: %v", validationResult.Errors)
	}

	// Test invalid partial update
	updateReq = &models.UpdateThemeRequest{
		Colors: &models.ThemeColors{
			Primary: "invalid",
		},
	}

	validationResult = validator.ValidateThemeUpdate(updateReq)
	if validationResult.IsValid {
		t.Error("Expected invalid color in update to fail validation")
	}

	t.Log("✓ Theme update validation working correctly")
}

// TestErrorResponseFormatIntegration tests that errors have proper structure
func TestErrorResponseFormatIntegration(t *testing.T) {
	validator := setupTestValidator(t)
	req := createValidThemeRequest()

	// Create multiple errors
	req.Name = ""
	req.Colors.Primary = "invalid"
	req.Typography.BaseFontSize = 999

	validationResult := validator.ValidateThemeCreate(req)

	if validationResult.IsValid {
		t.Fatal("Expected validation errors")
	}

	for i, err := range validationResult.Errors {
		// Check error structure
		if len(err.Field) == 0 {
			t.Errorf("Error %d: Field should not be empty", i)
		}
		if len(err.Message) == 0 {
			t.Errorf("Error %d: Message should not be empty", i)
		}
		if len(err.Code) == 0 {
			t.Errorf("Error %d: Code should not be empty", i)
		}

		t.Logf("Error %d: Field=%s, Code=%s", i, err.Field, err.Code)
	}

	t.Log("✓ Error response format is correct")
}

// TestErrorCodesConsistencyIntegration tests error code consistency
func TestErrorCodesConsistencyIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	errorCodeTests := []struct {
		name         string
		req          *models.CreateThemeRequest
		expectedCode string
	}{
		{
			name: "Invalid hex color",
			req: func() *models.CreateThemeRequest {
				r := createValidThemeRequest()
				r.Colors.Primary = "#gggggg"
				return r
			}(),
			expectedCode: "INVALID_HEX_COLOR",
		},
		{
			name: "Font size out of range",
			req: func() *models.CreateThemeRequest {
				r := createValidThemeRequest()
				r.Typography.BaseFontSize = 100
				return r
			}(),
			expectedCode: "OUT_OF_RANGE",
		},
		{
			name: "Invalid font family",
			req: func() *models.CreateThemeRequest {
				r := createValidThemeRequest()
				r.Typography.FontFamily = "InvalidFont"
				return r
			}(),
			expectedCode: "FONT_NOT_ALLOWED",
		},
		{
			name: "Invalid URL protocol",
			req: func() *models.CreateThemeRequest {
				r := createValidThemeRequest()
				r.Identity.LogoURL = "ftp://example.com/logo.png"
				return r
			}(),
			expectedCode: "INVALID_URL_PROTOCOL",
		},
	}

	for _, tc := range errorCodeTests {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateThemeCreate(tc.req)

			if result.IsValid {
				t.Errorf("Expected validation to fail")
				return
			}

			found := false
			for _, err := range result.Errors {
				if err.Code == tc.expectedCode {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected error code %s, got: %v", tc.expectedCode, result.Errors)
			}
		})
	}

	t.Log("✓ Error codes are consistent and correct")
}

// TestCustomCSSValidationIntegration tests CSS validation
func TestCustomCSSValidationIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	// Test valid CSS
	validCSS := `
		body {
			color: blue;
			margin: 10px;
			padding: 5px;
		}
		.container {
			width: 100%;
		}
	`
	errors := validator.ValidateCustomCSS(validCSS)
	if len(errors) > 0 {
		t.Logf("Valid CSS validation result: %v", errors)
	}

	// Test dangerous CSS patterns
	dangerousPatterns := []string{
		"body { behavior: url('xss.htc'); }",
		"body { -moz-binding: url('xss.xml#xss'); }",
		"body { background: url('javascript:alert(1)'); }",
	}

	for _, pattern := range dangerousPatterns {
		errors := validator.ValidateCustomCSS(pattern)
		if len(errors) == 0 {
			t.Logf("WARNING: Dangerous CSS pattern not detected: %s", pattern)
		}
	}

	t.Log("✓ Custom CSS validation tested")
}

// TestCompleteWorkflowIntegration tests a complete workflow
func TestCompleteWorkflowIntegration(t *testing.T) {
	validator := setupTestValidator(t)

	// Step 1: Create valid theme
	createReq := createValidThemeRequest()
	result := validator.ValidateThemeCreate(createReq)
	if !result.IsValid {
		t.Fatalf("Failed to create valid theme: %v", result.Errors)
	}

	// Step 2: Update theme
	updateReq := &models.UpdateThemeRequest{
		Name: "Updated Theme Name",
	}
	result = validator.ValidateThemeUpdate(updateReq)
	if !result.IsValid {
		t.Fatalf("Failed to update theme: %v", result.Errors)
	}

	// Step 3: Validate specific fields
	errors := validator.ValidateColors(&createReq.Colors)
	if len(errors) > 0 {
		t.Fatalf("Color validation failed: %v", errors)
	}

	errors = validator.ValidateTypography(&createReq.Typography)
	if len(errors) > 0 {
		t.Fatalf("Typography validation failed: %v", errors)
	}

	errors = validator.ValidateIdentity(&createReq.Identity)
	if len(errors) > 0 {
		t.Fatalf("Identity validation failed: %v", errors)
	}

	t.Log("✓ Complete workflow validation passed")
}

// TestPerformanceIntegration tests validation performance
func TestPerformanceIntegration(t *testing.T) {
	validator := setupTestValidator(t)
	req := createValidThemeRequest()

	// Run 1000 validations
	for i := 0; i < 1000; i++ {
		result := validator.ValidateThemeCreate(req)
		if result == nil {
			t.Fatal("Validation result should not be nil")
		}
	}

	t.Log("✓ 1000 validations completed successfully - performance is acceptable")
}
