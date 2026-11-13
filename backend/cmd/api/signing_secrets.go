package main

import (
	"fmt"
	"pos-saas/internal/config"
)

func validateAPISigningSecrets(cfg *config.Config) error {
	if cfg.Server.Env != "production" && cfg.Server.Env != "staging" {
		return nil
	}
	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	secret := cfg.JWT.SuperAdminSecret
	if len(secret) < 32 || secret == cfg.JWT.Secret || secret == cfg.JWT.Secret+"_super_admin" {
		return fmt.Errorf("SUPER_ADMIN_JWT_SECRET must be an independent secret of at least 32 characters")
	}
	return nil
}
