//go:build experimental
// +build experimental

package auth_test

import (
	"testing"
	"time"

	"pos-saas/internal/auth"
)

// TestJWTManagerInitialization tests JWT manager creation
func TestJWTManagerInitialization(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	if jm == nil {
		t.Fatal("JWT manager initialization failed")
	}

	t.Log("✓ JWT manager initialized successfully")
}

// TestGenerateAccessToken tests access token generation
func TestGenerateAccessToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	token, err := jm.GenerateAccessToken(
		1, 1, "user@example.com", "testuser",
		[]string{"admin", "manager"}, []string{"read", "write", "delete"},
		"192.168.1.1", "Mozilla/5.0", "session_123",
	)

	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	if token == "" {
		t.Fatal("Access token is empty")
	}

	t.Logf("✓ Access token generated: %s...", token[:20])
}

// TestValidateAccessToken tests access token validation
func TestValidateAccessToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate token
	token, err := jm.GenerateAccessToken(
		1, 1, "user@example.com", "testuser",
		[]string{"admin"}, []string{"read", "write"},
		"192.168.1.1", "Mozilla/5.0", "session_123",
	)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := jm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 1 {
		t.Fatalf("Expected user ID 1, got %d", claims.UserID)
	}

	if claims.TenantID != 1 {
		t.Fatalf("Expected tenant ID 1, got %d", claims.TenantID)
	}

	if claims.Email != "user@example.com" {
		t.Fatalf("Expected email 'user@example.com', got '%s'", claims.Email)
	}

	t.Log("✓ Access token validated successfully")
}

// TestExpiredAccessToken tests that expired tokens are rejected
func TestExpiredAccessToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, -1*time.Second, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate token (with negative expiry, so it's already expired)
	token, err := jm.GenerateAccessToken(
		1, 1, "user@example.com", "testuser",
		[]string{"admin"}, []string{"read"},
		"192.168.1.1", "Mozilla/5.0", "session_123",
	)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate (should fail)
	_, err = jm.ValidateAccessToken(token)
	if err == nil {
		t.Fatal("Expected validation to fail for expired token")
	}

	t.Log("✓ Expired token correctly rejected")
}

// TestGenerateRefreshToken tests refresh token generation
func TestGenerateRefreshToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	token, err := jm.GenerateRefreshToken(1, 1, "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	if token == "" {
		t.Fatal("Refresh token is empty")
	}

	t.Logf("✓ Refresh token generated: %s...", token[:20])
}

// TestValidateRefreshToken tests refresh token validation
func TestValidateRefreshToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate refresh token
	token, err := jm.GenerateRefreshToken(1, 1, "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Validate refresh token
	refreshToken, err := jm.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if refreshToken.UserID != 1 {
		t.Fatalf("Expected user ID 1, got %d", refreshToken.UserID)
	}

	if refreshToken.TenantID != 1 {
		t.Fatalf("Expected tenant ID 1, got %d", refreshToken.TenantID)
	}

	t.Log("✓ Refresh token validated successfully")
}

// TestRevokeRefreshToken tests refresh token revocation
func TestRevokeRefreshToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate refresh token
	token, err := jm.GenerateRefreshToken(1, 1, "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Revoke token
	err = jm.RevokeRefreshToken(token)
	if err != nil {
		t.Fatalf("Failed to revoke token: %v", err)
	}

	// Try to validate revoked token (should fail)
	_, err = jm.ValidateRefreshToken(token)
	if err == nil {
		t.Fatal("Expected validation to fail for revoked token")
	}

	t.Log("✓ Refresh token revoked successfully")
}

// TestRevokeAllUserTokens tests revoking all tokens for a user
func TestRevokeAllUserTokens(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate multiple tokens for same user
	token1, _ := jm.GenerateRefreshToken(1, 1, "user@example.com")
	token2, _ := jm.GenerateRefreshToken(1, 1, "user@example.com")
	token3, _ := jm.GenerateRefreshToken(2, 1, "other@example.com") // Different user

	// Revoke all tokens for user 1
	err := jm.RevokeAllUserTokens(1, 1)
	if err != nil {
		t.Fatalf("Failed to revoke all tokens: %v", err)
	}

	// Verify token1 and token2 are revoked
	_, err1 := jm.ValidateRefreshToken(token1)
	_, err2 := jm.ValidateRefreshToken(token2)

	if err1 == nil || err2 == nil {
		t.Fatal("Expected tokens to be revoked")
	}

	// Verify token3 is still valid (different user)
	_, err3 := jm.ValidateRefreshToken(token3)
	if err3 != nil {
		t.Fatal("Token for other user should still be valid")
	}

	t.Log("✓ All user tokens revoked successfully")
}

// TestRefreshAccessToken tests generating new access token from refresh token
func TestRefreshAccessToken(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate refresh token
	refreshToken, err := jm.GenerateRefreshToken(1, 1, "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Use refresh token to get new access token
	newAccessToken, err := jm.RefreshAccessToken(
		refreshToken,
		[]string{"admin"},
		[]string{"read", "write"},
		"192.168.1.1",
		"Mozilla/5.0",
	)
	if err != nil {
		t.Fatalf("Failed to refresh access token: %v", err)
	}

	// Validate new access token
	claims, err := jm.ValidateAccessToken(newAccessToken)
	if err != nil {
		t.Fatalf("Failed to validate refreshed token: %v", err)
	}

	if claims.UserID != 1 || claims.TenantID != 1 {
		t.Fatal("Refreshed token has incorrect claims")
	}

	t.Log("✓ Access token refreshed successfully")
}

// TestGetTokenInfo tests getting info without full validation
func TestGetTokenInfo(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate token
	token, err := jm.GenerateAccessToken(
		1, 1, "user@example.com", "testuser",
		[]string{"admin"}, []string{"read"},
		"192.168.1.1", "Mozilla/5.0", "session_123",
	)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Get token info
	info, err := jm.GetTokenInfo(token)
	if err != nil {
		t.Fatalf("Failed to get token info: %v", err)
	}

	if info["user_id"] != int64(1) {
		t.Fatal("Token info has incorrect user ID")
	}

	t.Log("✓ Token info retrieved successfully")
}

// TestCleanupExpiredTokens tests cleanup of expired tokens
func TestCleanupExpiredTokens(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, -1*time.Second, 7*24*time.Hour, "test_issuer", "test_audience")

	// Generate expired token
	expiredToken, _ := jm.GenerateRefreshToken(1, 1, "user@example.com")

	// Generate valid token
	jm.SetRefreshTokenExpiry(7 * 24 * time.Hour)
	validToken, _ := jm.GenerateRefreshToken(2, 1, "user2@example.com")

	if expiredToken == "" || validToken == "" {
		t.Fatal("token generation returned an empty token")
	}

	initialCount := jm.GetActiveTokenCount()

	// Clean up expired tokens
	jm.CleanupExpiredTokens()

	finalCount := jm.GetActiveTokenCount()

	if finalCount >= initialCount {
		t.Fatal("Cleanup did not remove expired tokens")
	}

	t.Logf("✓ Expired tokens cleaned up (before: %d, after: %d)", initialCount, finalCount)
}

// TestTokenExpiry tests setting token expiry durations
func TestTokenExpiry(t *testing.T) {
	signingKey := []byte("test_signing_key_32_chars_minimum")
	refreshKey := []byte("test_refresh_key_32_chars_minimum")

	jm := auth.NewJWTManager(signingKey, refreshKey, 15*time.Minute, 7*24*time.Hour, "test_issuer", "test_audience")

	// Change expiry
	jm.SetAccessTokenExpiry(30 * time.Minute)
	jm.SetRefreshTokenExpiry(14 * 24 * time.Hour)

	token, err := jm.GenerateAccessToken(
		1, 1, "user@example.com", "testuser",
		[]string{"admin"}, []string{"read"},
		"192.168.1.1", "Mozilla/5.0", "session_123",
	)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := jm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check that token won't expire soon
	if time.Until(claims.ExpiresAt) < 25*time.Minute {
		t.Fatal("Token expiry was not updated")
	}

	t.Log("✓ Token expiry updated successfully")
}
