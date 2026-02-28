package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// CurrentEnv returns the active runtime environment using common env keys.
func CurrentEnv() string {
	for _, key := range []string{"ENV", "APP_ENV", "GO_ENV", "NODE_ENV"} {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return strings.ToLower(value)
		}
	}
	return "development"
}

func IsProductionMode() bool {
	env := CurrentEnv()
	return env == "production" || env == "prod"
}

func IsDevelopmentMode() bool {
	env := CurrentEnv()
	return env == "development" || env == "dev" || env == "local"
}

func ResolveRuntimeSetting(keys []string, devFallback string, label string) (string, error) {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value, nil
		}
	}

	if !IsProductionMode() {
		return strings.TrimSpace(devFallback), nil
	}

	return "", fmt.Errorf("%s is required in production (set %s)", label, strings.Join(keys, " or "))
}

func ResolveRuntimeURL(keys []string, devFallback string, label string) (string, error) {
	value, err := ResolveRuntimeSetting(keys, devFallback, label)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(value, "/"), nil
}

// ValidateProductionURLConfig ensures required public URL/domain settings are present in production.
func ValidateProductionURLConfig() error {
	if !IsProductionMode() {
		return nil
	}

	type urlRequirement struct {
		label string
		key   string
	}

	requiredURLs := []urlRequirement{
		{label: "APP_BASE_URL", key: "APP_BASE_URL"},
		{label: "UPLOAD_BASE_URL", key: "UPLOAD_BASE_URL"},
		{label: "WEBHOOK_BASE_URL", key: "WEBHOOK_BASE_URL"},
		{label: "DASHBOARD_URL", key: "DASHBOARD_URL"},
	}

	for _, req := range requiredURLs {
		raw := strings.TrimSpace(os.Getenv(req.key))
		if raw == "" {
			return fmt.Errorf("%s is required in production", req.label)
		}

		parsed, err := url.Parse(raw)
		if err != nil || !parsed.IsAbs() || parsed.Host == "" {
			return fmt.Errorf("%s must be a valid absolute URL (got %q)", req.label, raw)
		}
	}

	restaurantDomain := strings.TrimSpace(os.Getenv("RESTAURANT_DOMAIN"))
	if restaurantDomain == "" {
		return fmt.Errorf("RESTAURANT_DOMAIN is required in production")
	}
	if strings.Contains(restaurantDomain, "://") {
		return fmt.Errorf("RESTAURANT_DOMAIN must be a host/domain, not a full URL (got %q)", restaurantDomain)
	}

	corsOrigins := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	if corsOrigins == "" {
		return fmt.Errorf("CORS_ORIGINS is required in production")
	}

	for _, item := range strings.Split(corsOrigins, ",") {
		origin := strings.TrimSpace(item)
		if origin == "" {
			continue
		}
		parsed, err := url.Parse(origin)
		if err != nil || !parsed.IsAbs() || parsed.Host == "" {
			return fmt.Errorf("CORS_ORIGINS contains invalid origin %q", origin)
		}
	}

	return nil
}
