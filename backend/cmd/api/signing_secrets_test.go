package main

import (
	"pos-saas/internal/config"
	"strings"
	"testing"
)

func TestStagingAndProductionRejectMissingSharedAndDerivedSigningSecrets(t *testing.T) {
	tenant := strings.Repeat("t", 48)
	for _, environment := range []string{"staging", "production"} {
		for _, secret := range []string{"", "short", tenant, tenant + "_super_admin"} {
			cfg := &config.Config{Server: config.ServerConfig{Env: environment}, JWT: config.JWTConfig{Secret: tenant, SuperAdminSecret: secret}}
			if err := validateAPISigningSecrets(cfg); err == nil {
				t.Errorf("%s accepted missing/shared/derived admin signing material", environment)
			}
		}
	}
}

func TestAPIAcceptsIndependentSigningSecretsAndRejectsWeakTenantSecret(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{Env: "staging"}, JWT: config.JWTConfig{Secret: strings.Repeat("t", 48), SuperAdminSecret: strings.Repeat("s", 48)}}
	if err := validateAPISigningSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	cfg.JWT.Secret = "short"
	if err := validateAPISigningSecrets(cfg); err == nil {
		t.Fatal("weak tenant signing material accepted")
	}
}
