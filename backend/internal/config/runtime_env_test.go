package config

import "testing"

var runtimeEnvKeys = []string{
	"ENV",
	"APP_ENV",
	"GO_ENV",
	"NODE_ENV",
	"APP_BASE_URL",
	"UPLOAD_BASE_URL",
	"WEBHOOK_BASE_URL",
	"DASHBOARD_URL",
	"LANDING_URL",
	"RESTAURANT_DOMAIN",
	"CORS_ORIGINS",
	"PRIMARY_URL",
	"SECONDARY_URL",
}

func clearRuntimeEnv(t *testing.T) {
	t.Helper()
	for _, key := range runtimeEnvKeys {
		t.Setenv(key, "")
	}
}

func seedValidProductionURLConfig(t *testing.T) {
	t.Helper()
	t.Setenv("ENV", "production")
	t.Setenv("APP_BASE_URL", "https://api.example.com")
	t.Setenv("UPLOAD_BASE_URL", "https://api.example.com/uploads")
	t.Setenv("WEBHOOK_BASE_URL", "https://api.example.com")
	t.Setenv("DASHBOARD_URL", "https://dashboard.example.com")
	t.Setenv("RESTAURANT_DOMAIN", "example.com")
	t.Setenv("CORS_ORIGINS", "https://www.example.com,https://dashboard.example.com")
}

func TestCurrentEnvPrefersConfiguredKeysInOrder(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("APP_ENV", "staging")
	t.Setenv("GO_ENV", "production")

	if got := CurrentEnv(); got != "staging" {
		t.Fatalf("CurrentEnv() = %q, want %q", got, "staging")
	}

	t.Setenv("ENV", "PROD")
	if got := CurrentEnv(); got != "prod" {
		t.Fatalf("CurrentEnv() with ENV set = %q, want %q", got, "prod")
	}
}

func TestResolveRuntimeSettingUsesDevFallbackOutsideProduction(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("ENV", "development")

	got, err := ResolveRuntimeSetting([]string{"PRIMARY_URL", "SECONDARY_URL"}, "http://localhost:8080", "primary URL")
	if err != nil {
		t.Fatalf("ResolveRuntimeSetting() unexpected error: %v", err)
	}
	if got != "http://localhost:8080" {
		t.Fatalf("ResolveRuntimeSetting() = %q, want dev fallback", got)
	}
}

func TestResolveRuntimeSettingFailsInProductionWhenMissing(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("ENV", "production")

	_, err := ResolveRuntimeSetting([]string{"PRIMARY_URL"}, "http://localhost:8080", "primary URL")
	if err == nil {
		t.Fatal("expected error in production when setting is missing")
	}
}

func TestResolveRuntimeURLTrimsTrailingSlash(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("ENV", "production")
	t.Setenv("PRIMARY_URL", "https://api.example.com///")

	got, err := ResolveRuntimeURL([]string{"PRIMARY_URL"}, "http://localhost:8080", "primary URL")
	if err != nil {
		t.Fatalf("ResolveRuntimeURL() unexpected error: %v", err)
	}
	if got != "https://api.example.com" {
		t.Fatalf("ResolveRuntimeURL() = %q, want %q", got, "https://api.example.com")
	}
}

func TestValidateProductionURLConfigAcceptsValidProductionConfig(t *testing.T) {
	clearRuntimeEnv(t)
	seedValidProductionURLConfig(t)

	if err := ValidateProductionURLConfig(); err != nil {
		t.Fatalf("ValidateProductionURLConfig() unexpected error: %v", err)
	}
}

func TestValidateProductionURLConfigRejectsInvalidRestaurantDomain(t *testing.T) {
	clearRuntimeEnv(t)
	seedValidProductionURLConfig(t)
	t.Setenv("RESTAURANT_DOMAIN", "https://example.com")

	err := ValidateProductionURLConfig()
	if err == nil {
		t.Fatal("expected error for invalid RESTAURANT_DOMAIN")
	}
}

func TestValidateProductionURLConfigRejectsInvalidCORSOrigin(t *testing.T) {
	clearRuntimeEnv(t)
	seedValidProductionURLConfig(t)
	t.Setenv("CORS_ORIGINS", "not-a-url")

	err := ValidateProductionURLConfig()
	if err == nil {
		t.Fatal("expected error for invalid CORS origin")
	}
}

func TestValidateProductionURLConfigNoOpOutsideProduction(t *testing.T) {
	clearRuntimeEnv(t)
	t.Setenv("ENV", "development")

	if err := ValidateProductionURLConfig(); err != nil {
		t.Fatalf("expected no error outside production, got %v", err)
	}
}
