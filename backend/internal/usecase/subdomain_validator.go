package usecase

import (
	"fmt"
	"regexp"
	"strings"
)

// SubdomainValidator provides subdomain validation and availability checking
type SubdomainValidator struct {
	// Reserved subdomains that cannot be used
	reserved map[string]bool
}

// NewSubdomainValidator creates a new subdomain validator
func NewSubdomainValidator() *SubdomainValidator {
	return &SubdomainValidator{
		reserved: map[string]bool{
			// System reserved
			"api":       true,
			"admin":     true,
			"dashboard": true,
			"www":       true,
			"mail":      true,
			"ftp":       true,
			"cdn":       true,
			"auth":      true,
			"app":       true,
			"docs":      true,
			"status":    true,
			"help":      true,
			"support":   true,
			"blog":      true,
			"shop":      true,
			"store":     true,
			"dev":       true,
			"staging":   true,
			"test":      true,
			"localhost": true,
			"127":       true,
			// Common services
			"git":        true,
			"gitlab":     true,
			"github":     true,
			"jenkins":    true,
			"docker":     true,
			"kubernetes": true,
		},
	}
}

// ValidateSlug checks if a subdomain slug is valid
// Returns error if invalid
func (sv *SubdomainValidator) ValidateSlug(slug string) error {
	// Check if empty
	if slug == "" {
		return fmt.Errorf("subdomain cannot be empty")
	}

	// Check length (min 2, max 63 characters)
	if len(slug) < 2 {
		return fmt.Errorf("subdomain must be at least 2 characters long")
	}
	if len(slug) > 63 {
		return fmt.Errorf("subdomain must not exceed 63 characters")
	}

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Check for reserved names
	if sv.reserved[slug] {
		return fmt.Errorf("subdomain '%s' is reserved and cannot be used", slug)
	}

	// Validate format: alphanumeric and hyphens only
	// Must start and end with alphanumeric
	slugRegex := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
	if !slugRegex.MatchString(slug) {
		return fmt.Errorf("subdomain can only contain lowercase letters, numbers, and hyphens (must start and end with alphanumeric)")
	}

	// Check for consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("subdomain cannot contain consecutive hyphens")
	}

	return nil
}

// ToSlug converts a restaurant name to a valid subdomain slug
// Example: "My Cool Restaurant" -> "my-cool-restaurant"
func (sv *SubdomainValidator) ToSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple hyphens with single hyphen
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// If empty after processing, return a default
	if slug == "" {
		slug = "restaurant"
	}

	return slug
}

// GenerateUniqueSlug generates a unique slug from a restaurant name
// This should be called to auto-generate a slug suggestion
// The actual uniqueness check should happen at the database level
func (sv *SubdomainValidator) GenerateUniqueSlug(name string) string {
	baseSlug := sv.ToSlug(name)

	// Check if base slug is reserved
	if sv.reserved[baseSlug] {
		// Add restaurant prefix to make it unique
		baseSlug = "restaurant-" + baseSlug
	}

	return baseSlug
}

// ValidateAndFormat validates a slug and returns formatted version
// Returns the formatted slug and error if validation fails
func (sv *SubdomainValidator) ValidateAndFormat(slug string) (string, error) {
	// Normalize
	formatted := strings.ToLower(strings.TrimSpace(slug))

	// Validate
	if err := sv.ValidateSlug(formatted); err != nil {
		return "", err
	}

	return formatted, nil
}
