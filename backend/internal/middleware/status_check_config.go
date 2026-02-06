package middleware

// StatusCheckConfig holds configuration for tenant status checks
type StatusCheckConfig struct {
	// Enabled determines if status checks are active
	Enabled bool

	// SkipPaths are exact paths that should skip the status check
	SkipPaths []string

	// SkipPrefixes are path prefixes that should skip the status check
	SkipPrefixes []string

	// AllowedWhenBlocked are paths that should be accessible even when tenant is blocked
	AllowedWhenBlocked []string

	// BlockedStatusCode is the HTTP status code to return when blocked
	BlockedStatusCode int
}

// DefaultStatusCheckConfig returns the default configuration for status checks
func DefaultStatusCheckConfig() *StatusCheckConfig {
	return &StatusCheckConfig{
		Enabled: true,

		// Exact paths to skip
		SkipPaths: []string{
			"/api/v1/tenant/status",
			"/api/v1/auth/logout",
			"/api/v1/auth/refresh",
			"/api/v1/auth/validate-token",
		},

		// Path prefixes to skip
		SkipPrefixes: []string{
			"/api/v1/public/",
			"/api/v1/super-admin/",
			"/api/v1/auth/",
			"/health",
			"/metrics",
			"/docs",
		},

		// Paths allowed even when tenant is blocked
		AllowedWhenBlocked: []string{
			"/api/v1/tenant/status",
			"/api/v1/auth/logout",
		},

		// Use 403 Forbidden for blocked tenants
		BlockedStatusCode: 403,
	}
}

// MergeConfig merges a custom config with defaults
// Custom config values override defaults
func MergeConfig(custom *StatusCheckConfig) *StatusCheckConfig {
	if custom == nil {
		return DefaultStatusCheckConfig()
	}

	defaults := DefaultStatusCheckConfig()

	// If custom has any skip paths, use them (don't merge)
	if len(custom.SkipPaths) > 0 {
		defaults.SkipPaths = custom.SkipPaths
	}

	// If custom has any skip prefixes, use them (don't merge)
	if len(custom.SkipPrefixes) > 0 {
		defaults.SkipPrefixes = custom.SkipPrefixes
	}

	// Use custom enabled flag
	defaults.Enabled = custom.Enabled

	// Use custom status code if provided
	if custom.BlockedStatusCode > 0 {
		defaults.BlockedStatusCode = custom.BlockedStatusCode
	}

	return defaults
}
