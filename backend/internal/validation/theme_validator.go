//go:build legacy
// +build legacy

package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ThemeValidator handles all validation for theme-related operations
type ThemeValidator struct {
	// Allowed font families for security
	allowedFonts map[string]bool
	// SQL injection patterns to prevent
	sqlInjectionPatterns []*regexp.Regexp
	// Max sizes for fields
	maxNameLength    int
	maxDescLength    int
	maxCSSSize       int
	maxTitleLength   int
}

// NewThemeValidator creates a new theme validator instance
func NewThemeValidator() *ThemeValidator {
	return &ThemeValidator{
		allowedFonts: map[string]bool{
			"Inter":              true,
			"Roboto":             true,
			"Open Sans":          true,
			"Lato":               true,
			"Montserrat":         true,
			"Raleway":            true,
			"Playfair Display":   true,
			"Poppins":            true,
			"Nunito":             true,
			"Source Sans Pro":    true,
			"Ubuntu":             true,
			"Lora":               true,
			"Merriweather":       true,
			"Georgia":            true,
			"Times New Roman":    true,
			"Arial":              true,
			"Helvetica":          true,
			"Verdana":            true,
			"Courier New":        true,
			"Comic Sans MS":      true,
		},
		sqlInjectionPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(DROP|DELETE|INSERT|UPDATE|SELECT|EXEC|EXECUTE|UNION|ALTER)\s`),
			regexp.MustCompile(`(?i)(;|--|\*/|/\*|xp_|sp_)`),
			regexp.MustCompile(`(?i)('|")\s*(OR|AND)\s*('|")?.*=.*`),
		},
		maxNameLength:  100,
		maxDescLength:  500,
		maxCSSSize:     50 * 1024, // 50KB
		maxTitleLength: 255,
	}
}

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

// ValidationResult contains the validation result with all errors
type ValidationResult struct {
	IsValid bool
	Errors  []*ValidationError
}

// ValidateThemeCreate validates a CreateThemeRequest
func (v *ThemeValidator) ValidateThemeCreate(req *models.CreateThemeRequest) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []*ValidationError{},
	}

	if req == nil {
		result.IsValid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "request",
			Message: "Request cannot be nil",
			Code:    "INVALID_REQUEST",
		})
		return result
	}

	// Validate theme name
	if errs := v.validateRequired(req.Name, "name", "Theme name"); len(errs) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, errs...)
	}

	if len(req.Name) > v.maxNameLength {
		result.IsValid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("Theme name must not exceed %d characters", v.maxNameLength),
			Code:    "NAME_TOO_LONG",
		})
	}

	// Check for SQL injection patterns
	if errs := v.validateAgainstSQLInjection(req.Name, "name"); len(errs) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, errs...)
	}

	// Validate colors
	if errs := v.ValidateColors(&req.Colors); len(errs) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, errs...)
	}

	// Validate typography
	if errs := v.ValidateTypography(&req.Typography); len(errs) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, errs...)
	}

	// Validate identity
	if errs := v.ValidateIdentity(&req.Identity); len(errs) > 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, errs...)
	}

	return result
}

// ValidateThemeUpdate validates an UpdateThemeRequest
func (v *ThemeValidator) ValidateThemeUpdate(req *models.UpdateThemeRequest) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []*ValidationError{},
	}

	if req == nil {
		result.IsValid = false
		result.Errors = append(result.Errors, &ValidationError{
			Field:   "request",
			Message: "Request cannot be nil",
			Code:    "INVALID_REQUEST",
		})
		return result
	}

	// Validate name if provided
	if req.Name != "" {
		if len(req.Name) > v.maxNameLength {
			result.IsValid = false
			result.Errors = append(result.Errors, &ValidationError{
				Field:   "name",
				Message: fmt.Sprintf("Theme name must not exceed %d characters", v.maxNameLength),
				Code:    "NAME_TOO_LONG",
			})
		}

		if errs := v.validateAgainstSQLInjection(req.Name, "name"); len(errs) > 0 {
			result.IsValid = false
			result.Errors = append(result.Errors, errs...)
		}
	}

	// Validate colors if provided
	if req.Colors != nil {
		if errs := v.ValidateColors(req.Colors); len(errs) > 0 {
			result.IsValid = false
			result.Errors = append(result.Errors, errs...)
		}
	}

	// Validate typography if provided
	if req.Typography != nil {
		if errs := v.ValidateTypography(req.Typography); len(errs) > 0 {
			result.IsValid = false
			result.Errors = append(result.Errors, errs...)
		}
	}

	// Validate identity if provided
	if req.Identity != nil {
		if errs := v.ValidateIdentity(req.Identity); len(errs) > 0 {
			result.IsValid = false
			result.Errors = append(result.Errors, errs...)
		}
	}

	return result
}

// ValidateColors validates all theme colors
func (v *ThemeValidator) ValidateColors(colors *models.ThemeColors) []*ValidationError {
	var errors []*ValidationError

	if colors == nil {
		errors = append(errors, &ValidationError{
			Field:   "colors",
			Message: "Colors configuration is required",
			Code:    "COLORS_REQUIRED",
		})
		return errors
	}

	colorFields := map[string]string{
		"primary":     colors.Primary,
		"secondary":   colors.Secondary,
		"accent":      colors.Accent,
		"background":  colors.Background,
		"text":        colors.Text,
		"border":      colors.Border,
		"shadow":      colors.Shadow,
	}

	for fieldName, colorValue := range colorFields {
		if errs := v.validateHexColor(colorValue, "colors."+fieldName); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	return errors
}

// ValidateTypography validates typography settings
func (v *ThemeValidator) ValidateTypography(typography *models.TypographySettings) []*ValidationError {
	var errors []*ValidationError

	if typography == nil {
		errors = append(errors, &ValidationError{
			Field:   "typography",
			Message: "Typography settings are required",
			Code:    "TYPOGRAPHY_REQUIRED",
		})
		return errors
	}

	// Validate font family
	if errs := v.validateFontFamily(typography.FontFamily, "typography.fontFamily"); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	// Validate base font size: 10-24px
	if errs := v.validateNumberRange(typography.BaseFontSize, 10, 24, "typography.baseFontSize", "Base font size"); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	// Validate border radius: 0-50px
	if errs := v.validateNumberRange(typography.BorderRadius, 0, 50, "typography.borderRadius", "Border radius"); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	// Validate line height: 1.0-3.0
	if errs := v.validateFloatRange(typography.LineHeight, 1.0, 3.0, "typography.lineHeight", "Line height"); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	return errors
}

// ValidateIdentity validates website identity settings
func (v *ThemeValidator) ValidateIdentity(identity *models.WebsiteIdentity) []*ValidationError {
	var errors []*ValidationError

	if identity == nil {
		errors = append(errors, &ValidationError{
			Field:   "identity",
			Message: "Identity settings are required",
			Code:    "IDENTITY_REQUIRED",
		})
		return errors
	}

	// Validate site title
	if errs := v.validateRequired(identity.SiteTitle, "identity.siteTitle", "Site title"); len(errs) > 0 {
		errors = append(errors, errs...)
	}

	if len(identity.SiteTitle) > v.maxTitleLength {
		errors = append(errors, &ValidationError{
			Field:   "identity.siteTitle",
			Message: fmt.Sprintf("Site title must not exceed %d characters", v.maxTitleLength),
			Code:    "TITLE_TOO_LONG",
		})
	}

	// Validate logo URL if provided
	if identity.LogoURL != "" {
		if errs := v.validateURL(identity.LogoURL, "identity.logoUrl", true); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate favicon URL if provided
	if identity.FaviconURL != "" {
		if errs := v.validateURL(identity.FaviconURL, "identity.faviconUrl", true); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	return errors
}

// ValidateCustomCSS validates custom CSS string
func (v *ThemeValidator) ValidateCustomCSS(css string) []*ValidationError {
	var errors []*ValidationError

	if css == "" {
		return errors // Custom CSS is optional
	}

	// Check size limit
	if len(css) > v.maxCSSSize {
		errors = append(errors, &ValidationError{
			Field:   "customCSS",
			Message: fmt.Sprintf("Custom CSS must not exceed %d KB", v.maxCSSSize/1024),
			Code:    "CSS_TOO_LARGE",
		})
	}

	// Check for dangerous CSS patterns
	dangerousPatterns := []string{
		"@import",
		"behavior:",
		"expression(",
		"-moz-binding",
		"javascript:",
		"vbscript:",
	}

	cssLower := strings.ToLower(css)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cssLower, pattern) {
			errors = append(errors, &ValidationError{
				Field:   "customCSS",
				Message: fmt.Sprintf("Dangerous CSS pattern '%s' is not allowed", pattern),
				Code:    "DANGEROUS_CSS",
			})
		}
	}

	return errors
}

// ValidateComponentConfig validates a component's configuration
func (v *ThemeValidator) ValidateComponentConfig(componentType string, config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if componentType == "" {
		errors = append(errors, &ValidationError{
			Field:   "componentType",
			Message: "Component type is required",
			Code:    "COMPONENT_TYPE_REQUIRED",
		})
		return errors
	}

	// Type-specific validation
	switch componentType {
	case "hero":
		errors = append(errors, v.validateHeroConfig(config)...)
	case "products":
		errors = append(errors, v.validateProductsConfig(config)...)
	case "why-us":
		errors = append(errors, v.validateWhyUsConfig(config)...)
	case "contact":
		errors = append(errors, v.validateContactConfig(config)...)
	case "testimonials":
		errors = append(errors, v.validateTestimonialsConfig(config)...)
	case "cta":
		errors = append(errors, v.validateCTAConfig(config)...)
	}

	return errors
}

// Helper validation functions

// validateRequired checks if a string is not empty
func (v *ThemeValidator) validateRequired(value, field, label string) []*ValidationError {
	var errors []*ValidationError

	if strings.TrimSpace(value) == "" {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s is required", label),
			Code:    "REQUIRED",
		})
	}

	return errors
}

// validateHexColor checks if a string is a valid hex color
func (v *ThemeValidator) validateHexColor(value, field string) []*ValidationError {
	var errors []*ValidationError

	if value == "" {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("Color value is required"),
			Code:    "COLOR_REQUIRED",
		})
		return errors
	}

	// Remove # prefix if present
	color := strings.TrimPrefix(value, "#")

	// Check length - allow #RGB or #RRGGBB or #RRGGBBAA
	if len(color) != 3 && len(color) != 6 && len(color) != 8 {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: "Color must be in hex format (#RGB, #RRGGBB, or #RRGGBBAA)",
			Code:    "INVALID_HEX_COLOR",
		})
		return errors
	}

	// Validate all characters are hexadecimal
	for _, c := range color {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			errors = append(errors, &ValidationError{
				Field:   field,
				Message: "Color must contain only valid hexadecimal characters (0-9, a-f, A-F)",
				Code:    "INVALID_HEX_CHARS",
			})
			return errors
		}
	}

	return errors
}

// validateURL checks if a string is a valid URL
func (v *ThemeValidator) validateURL(value, field string, allowEmpty bool) []*ValidationError {
	var errors []*ValidationError

	if value == "" {
		if !allowEmpty {
			errors = append(errors, &ValidationError{
				Field:   field,
				Message: "URL is required",
				Code:    "URL_REQUIRED",
			})
		}
		return errors
	}

	// Must start with http:// or https://
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: "URL must start with http:// or https://",
			Code:    "INVALID_URL_PROTOCOL",
		})
		return errors
	}

	// Basic URL validation using regex
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(value) {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: "Invalid URL format",
			Code:    "INVALID_URL_FORMAT",
		})
		return errors
	}

	// Check for XSS attempts in URL
	if strings.Contains(strings.ToLower(value), "javascript:") || strings.Contains(strings.ToLower(value), "data:") {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: "URL contains potentially dangerous content",
			Code:    "DANGEROUS_URL",
		})
		return errors
	}

	return errors
}

// validateFontFamily checks if font is in allowed list
func (v *ThemeValidator) validateFontFamily(value, field string) []*ValidationError {
	var errors []*ValidationError

	if value == "" {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: "Font family is required",
			Code:    "FONT_REQUIRED",
		})
		return errors
	}

	if !v.allowedFonts[value] {
		allowedList := []string{}
		for font := range v.allowedFonts {
			allowedList = append(allowedList, font)
		}
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("Font '%s' is not in the allowed list", value),
			Code:    "FONT_NOT_ALLOWED",
		})
		return errors
	}

	return errors
}

// validateNumberRange checks if a number is within a valid range
func (v *ThemeValidator) validateNumberRange(value, min, max int, field, label string) []*ValidationError {
	var errors []*ValidationError

	if value < min || value > max {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be between %d and %d", label, min, max),
			Code:    "OUT_OF_RANGE",
		})
	}

	return errors
}

// validateFloatRange checks if a float is within a valid range
func (v *ThemeValidator) validateFloatRange(value float64, min, max float64, field, label string) []*ValidationError {
	var errors []*ValidationError

	if value < min || value > max {
		errors = append(errors, &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be between %.1f and %.1f", label, min, max),
			Code:    "OUT_OF_RANGE",
		})
	}

	return errors
}

// validateAgainstSQLInjection checks for SQL injection patterns
func (v *ThemeValidator) validateAgainstSQLInjection(value, field string) []*ValidationError {
	var errors []*ValidationError

	for _, pattern := range v.sqlInjectionPatterns {
		if pattern.MatchString(value) {
			errors = append(errors, &ValidationError{
				Field:   field,
				Message: "Input contains potentially dangerous SQL patterns",
				Code:    "SQL_INJECTION_DETECTED",
			})
			break
		}
	}

	return errors
}

// Component-specific validators

// validateHeroConfig validates hero component settings
func (v *ThemeValidator) validateHeroConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	// Hero component may have backgroundImage, overlayColor, overlayOpacity, etc.
	if config == nil {
		return errors
	}

	// Check for background image URL if present
	if bgImage, ok := config["backgroundImage"].(string); ok && bgImage != "" {
		if errs := v.validateURL(bgImage, "components.hero.backgroundImage", true); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Check overlay color if present
	if overlayColor, ok := config["overlayColor"].(string); ok && overlayColor != "" {
		if errs := v.validateHexColor(overlayColor, "components.hero.overlayColor"); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	return errors
}

// validateProductsConfig validates products component settings
func (v *ThemeValidator) validateProductsConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if config == nil {
		return errors
	}

	// Check layout if present
	if layout, ok := config["layout"].(string); ok {
		validLayouts := map[string]bool{"2": true, "3": true, "4": true}
		if !validLayouts[layout] {
			errors = append(errors, &ValidationError{
				Field:   "components.products.layout",
				Message: "Layout must be 2, 3, or 4 columns",
				Code:    "INVALID_LAYOUT",
			})
		}
	}

	return errors
}

// validateWhyUsConfig validates why-us component settings
func (v *ThemeValidator) validateWhyUsConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if config == nil {
		return errors
	}

	// Feature list validation would go here
	return errors
}

// validateContactConfig validates contact component settings
func (v *ThemeValidator) validateContactConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if config == nil {
		return errors
	}

	// Email and phone validation could go here
	return errors
}

// validateTestimonialsConfig validates testimonials component settings
func (v *ThemeValidator) validateTestimonialsConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if config == nil {
		return errors
	}

	// Auto-rotate and max items validation
	if maxItems, ok := config["maxItems"].(float64); ok {
		if maxItems < 1 || maxItems > 20 {
			errors = append(errors, &ValidationError{
				Field:   "components.testimonials.maxItems",
				Message: "Max items must be between 1 and 20",
				Code:    "INVALID_MAX_ITEMS",
			})
		}
	}

	return errors
}

// validateCTAConfig validates CTA component settings
func (v *ThemeValidator) validateCTAConfig(config map[string]interface{}) []*ValidationError {
	var errors []*ValidationError

	if config == nil {
		return errors
	}

	// Button color and text validation
	if btnColor, ok := config["buttonColor"].(string); ok && btnColor != "" {
		if errs := v.validateHexColor(btnColor, "components.cta.buttonColor"); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	return errors
}

// FormatValidationErrors converts validation errors to user-friendly message
func FormatValidationErrors(errors []*ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		if err.Field != "" {
			messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
		} else {
			messages = append(messages, err.Message)
		}
	}

	return strings.Join(messages, "; ")
}

// ValidateColorContrast checks if two colors have sufficient contrast (WCAG AA 4.5:1)
func (v *ThemeValidator) ValidateColorContrast(foreground, background string) bool {
	fg, err1 := hexToRGB(foreground)
	bg, err2 := hexToRGB(background)

	if err1 != nil || err2 != nil {
		return false // Invalid color, fail validation
	}

	ratio := calculateContrastRatio(fg, bg)
	return ratio >= 4.5 // WCAG AA standard
}

// Helper function to convert hex to RGB
func hexToRGB(hex string) ([]float64, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		// Convert #RGB to #RRGGBB
		hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
	}

	if len(hex) != 6 {
		return nil, errors.New("invalid hex color")
	}

	r, g, b := 0, 0, 0
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return nil, err
	}

	return []float64{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}, nil
}

// Helper function to calculate relative luminance
func getRelativeLuminance(rgb []float64) float64 {
	var r, g, b float64

	for i, val := range rgb {
		if val <= 0.03928 {
			switch i {
			case 0:
				r = val / 12.92
			case 1:
				g = val / 12.92
			case 2:
				b = val / 12.92
			}
		} else {
			switch i {
			case 0:
				r = ((val + 0.055) / 1.055) * ((val + 0.055) / 1.055)
			case 1:
				g = ((val + 0.055) / 1.055) * ((val + 0.055) / 1.055)
			case 2:
				b = ((val + 0.055) / 1.055) * ((val + 0.055) / 1.055)
			}
		}
	}

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// Helper function to calculate contrast ratio
func calculateContrastRatio(fg, bg []float64) float64 {
	l1 := getRelativeLuminance(fg)
	l2 := getRelativeLuminance(bg)

	lighter := l1
	darker := l2
	if l2 > l1 {
		lighter = l2
		darker = l1
	}

	return (lighter + 0.05) / (darker + 0.05)
}

// Additional validation utility functions

// ValidateThemeSlug validates a theme slug
func (v *ThemeValidator) ValidateThemeSlug(slug string) []*ValidationError {
	var errors []*ValidationError

	if slug == "" {
		errors = append(errors, &ValidationError{
			Field:   "slug",
			Message: "Theme slug is required",
			Code:    "SLUG_REQUIRED",
		})
		return errors
	}

	// Slug must be lowercase alphanumeric with hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(slug) {
		errors = append(errors, &ValidationError{
			Field:   "slug",
			Message: "Slug must contain only lowercase letters, numbers, and hyphens",
			Code:    "INVALID_SLUG_FORMAT",
		})
	}

	// Max length
	if len(slug) > 100 {
		errors = append(errors, &ValidationError{
			Field:   "slug",
			Message: "Slug must not exceed 100 characters",
			Code:    "SLUG_TOO_LONG",
		})
	}

	return errors
}

// ValidateComponentTitle validates a component title
func (v *ThemeValidator) ValidateComponentTitle(title string) []*ValidationError {
	var errors []*ValidationError

	if title == "" {
		errors = append(errors, &ValidationError{
			Field:   "title",
			Message: "Component title is required",
			Code:    "TITLE_REQUIRED",
		})
		return errors
	}

	if utf8.RuneCountInString(title) > 255 {
		errors = append(errors, &ValidationError{
			Field:   "title",
			Message: "Component title must not exceed 255 characters",
			Code:    "TITLE_TOO_LONG",
		})
	}

	return errors
}
