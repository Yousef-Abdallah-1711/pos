package jwt

import (
	"testing"
)

func TestSuperAdminTokenIncludesAndValidatesExpiration(t *testing.T) {
	service, err := NewSuperAdminTokenService("test-secret-at-least-32-characters-long", "2m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	token, expiresAt, err := service.GenerateSuperAdminToken(7, "admin@example.test", 3)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := service.ValidateSuperAdminToken(token)
	if err != nil {
		t.Fatalf("validate generated token: %v", err)
	}
	if claims.ExpiresAtUnix != expiresAt.Unix() {
		t.Fatalf("expiration claim = %d, want %d", claims.ExpiresAtUnix, expiresAt.Unix())
	}
}

func TestSuperAdminTokenRejectsExpiredToken(t *testing.T) {
	service, err := NewSuperAdminTokenService("test-secret-at-least-32-characters-long", "-1m")
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	token, _, err := service.GenerateSuperAdminToken(7, "admin@example.test", 3)
	if err != nil {
		t.Fatalf("generate expired token: %v", err)
	}
	if _, err := service.ValidateSuperAdminToken(token); err == nil {
		t.Fatal("expired super-admin token unexpectedly validated")
	}
}
