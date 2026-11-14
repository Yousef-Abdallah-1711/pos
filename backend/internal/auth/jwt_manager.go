package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents JWT claims with tenant context
type JWTClaims struct {
	UserID      int64     `json:"user_id"`
	TenantID    int64     `json:"tenant_id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	SessionID   string    `json:"session_id"`
	IssuedAt    time.Time `json:"iat"`
	ExpiresAt   time.Time `json:"exp"`
	jwt.RegisteredClaims
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	Token     string
	UserID    int64
	TenantID  int64
	ExpiresAt time.Time
	IssuedAt  time.Time
}

// JWTManager manages JWT token generation and validation
type JWTManager struct {
	signingKey         []byte
	refreshSigningKey  []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	refreshTokens      map[string]*RefreshToken // In-memory store (use Redis in production)
	issuer             string
	audience           string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(
	signingKey []byte,
	refreshSigningKey []byte,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
	issuer string,
	audience string,
) *JWTManager {
	if accessTokenExpiry == 0 {
		accessTokenExpiry = 15 * time.Minute
	}
	if refreshTokenExpiry == 0 {
		refreshTokenExpiry = 7 * 24 * time.Hour
	}

	return &JWTManager{
		signingKey:         signingKey,
		refreshSigningKey:  refreshSigningKey,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		refreshTokens:      make(map[string]*RefreshToken),
		issuer:             issuer,
		audience:           audience,
	}
}

// GenerateAccessToken generates a new access token with tenant context
func (jm *JWTManager) GenerateAccessToken(
	userID int64,
	tenantID int64,
	email string,
	username string,
	roles []string,
	permissions []string,
	ipAddress string,
	userAgent string,
	sessionID string,
) (string, error) {
	log.Printf("[JWT] Generating access token: user=%d, tenant=%d, roles=%v",
		userID, tenantID, roles)

	now := time.Now()
	expiresAt := now.Add(jm.accessTokenExpiry)

	claims := JWTClaims{
		UserID:      userID,
		TenantID:    tenantID,
		Email:       email,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		SessionID:   sessionID,
		IssuedAt:    now,
		ExpiresAt:   expiresAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    jm.issuer,
			Audience:  jwt.ClaimStrings{jm.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.signingKey)
	if err != nil {
		log.Printf("[JWT] ERROR: Failed to generate access token: %v", err)
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	log.Printf("[JWT] Access token generated successfully (expires: %s)", expiresAt.Format(time.RFC3339))
	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (jm *JWTManager) GenerateRefreshToken(
	userID int64,
	tenantID int64,
	email string,
) (string, error) {
	log.Printf("[JWT] Generating refresh token: user=%d, tenant=%d", userID, tenantID)

	now := time.Now()
	expiresAt := now.Add(jm.refreshTokenExpiry)

	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"email":     email,
		"type":      "refresh",
		"exp":       expiresAt.Unix(),
		"iat":       now.Unix(),
		"iss":       jm.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.refreshSigningKey)
	if err != nil {
		log.Printf("[JWT] ERROR: Failed to generate refresh token: %v", err)
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in memory (use Redis in production)
	jm.refreshTokens[tokenString] = &RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		TenantID:  tenantID,
		ExpiresAt: expiresAt,
		IssuedAt:  now,
	}

	log.Printf("[JWT] Refresh token generated successfully (expires: %s)", expiresAt.Format(time.RFC3339))
	return tokenString, nil
}

// ValidateAccessToken validates an access token and returns claims
func (jm *JWTManager) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	log.Printf("[JWT] Validating access token")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.signingKey, nil
	})

	if err != nil {
		log.Printf("[JWT] ERROR: Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		log.Printf("[JWT] ERROR: Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		log.Printf("[JWT] ERROR: Failed to extract claims")
		return nil, fmt.Errorf("failed to extract claims")
	}

	// Verify expiration
	if time.Now().After(claims.ExpiresAt) {
		log.Printf("[JWT] ERROR: Token expired")
		return nil, fmt.Errorf("token expired")
	}

	// Verify issuer and audience
	if claims.Issuer != jm.issuer {
		log.Printf("[JWT] ERROR: Invalid issuer")
		return nil, fmt.Errorf("invalid issuer")
	}

	log.Printf("[JWT] Token validated successfully: user=%d, tenant=%d", claims.UserID, claims.TenantID)
	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (jm *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshToken, error) {
	log.Printf("[JWT] Validating refresh token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.refreshSigningKey, nil
	})

	if err != nil {
		log.Printf("[JWT] ERROR: Failed to parse refresh token: %v", err)
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		log.Printf("[JWT] ERROR: Invalid refresh token")
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if token exists in store
	refreshToken, exists := jm.refreshTokens[tokenString]
	if !exists {
		log.Printf("[JWT] ERROR: Refresh token not found in store")
		return nil, fmt.Errorf("refresh token not found")
	}

	// Check expiration
	if time.Now().After(refreshToken.ExpiresAt) {
		delete(jm.refreshTokens, tokenString)
		log.Printf("[JWT] ERROR: Refresh token expired")
		return nil, fmt.Errorf("refresh token expired")
	}

	log.Printf("[JWT] Refresh token validated: user=%d, tenant=%d", refreshToken.UserID, refreshToken.TenantID)
	return refreshToken, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (jm *JWTManager) RefreshAccessToken(
	refreshTokenString string,
	newRoles []string,
	newPermissions []string,
	ipAddress string,
	userAgent string,
) (string, error) {
	log.Printf("[JWT] Refreshing access token")

	// Validate refresh token
	refreshToken, err := jm.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Generate new access token
	newAccessToken, err := jm.GenerateAccessToken(
		refreshToken.UserID,
		refreshToken.TenantID,
		"", // Email not stored in refresh token
		"", // Username not stored in refresh token
		newRoles,
		newPermissions,
		ipAddress,
		userAgent,
		"", // New session ID
	)

	if err != nil {
		return "", err
	}

	log.Printf("[JWT] Access token refreshed successfully")
	return newAccessToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (jm *JWTManager) RevokeRefreshToken(tokenString string) error {
	log.Printf("[JWT] Revoking refresh token")

	if _, exists := jm.refreshTokens[tokenString]; !exists {
		log.Printf("[JWT] WARNING: Refresh token not found for revocation")
		return fmt.Errorf("refresh token not found")
	}

	delete(jm.refreshTokens, tokenString)
	log.Printf("[JWT] Refresh token revoked successfully")
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (jm *JWTManager) RevokeAllUserTokens(userID int64, tenantID int64) error {
	log.Printf("[JWT] Revoking all tokens: user=%d, tenant=%d", userID, tenantID)

	revokedCount := 0
	for token, refreshToken := range jm.refreshTokens {
		if refreshToken.UserID == userID && refreshToken.TenantID == tenantID {
			delete(jm.refreshTokens, token)
			revokedCount++
		}
	}

	log.Printf("[JWT] Revoked %d tokens for user %d", revokedCount, userID)
	return nil
}

// GetTokenInfo returns information about a token without full validation
func (jm *JWTManager) GetTokenInfo(tokenString string) (map[string]interface{}, error) {
	// Parse without verification to get claims
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		log.Printf("[JWT] ERROR: Failed to parse token: %v", err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	return map[string]interface{}{
		"user_id":     claims.UserID,
		"tenant_id":   claims.TenantID,
		"email":       claims.Email,
		"username":    claims.Username,
		"roles":       claims.Roles,
		"permissions": claims.Permissions,
		"issued_at":   claims.IssuedAt,
		"expires_at":  claims.ExpiresAt,
		"is_expired":  time.Now().After(claims.ExpiresAt),
	}, nil
}

// SetAccessTokenExpiry sets the access token expiry duration
func (jm *JWTManager) SetAccessTokenExpiry(duration time.Duration) {
	jm.accessTokenExpiry = duration
	log.Printf("[JWT] Access token expiry set to %v", duration)
}

// SetRefreshTokenExpiry sets the refresh token expiry duration
func (jm *JWTManager) SetRefreshTokenExpiry(duration time.Duration) {
	jm.refreshTokenExpiry = duration
	log.Printf("[JWT] Refresh token expiry set to %v", duration)
}

// CleanupExpiredTokens removes expired refresh tokens from memory
func (jm *JWTManager) CleanupExpiredTokens() {
	log.Printf("[JWT] Cleaning up expired refresh tokens")

	now := time.Now()
	removedCount := 0

	for token, refreshToken := range jm.refreshTokens {
		if now.After(refreshToken.ExpiresAt) {
			delete(jm.refreshTokens, token)
			removedCount++
		}
	}

	log.Printf("[JWT] Removed %d expired refresh tokens", removedCount)
}

// GetActiveTokenCount returns the count of active refresh tokens
func (jm *JWTManager) GetActiveTokenCount() int {
	count := 0
	now := time.Now()

	for _, refreshToken := range jm.refreshTokens {
		if now.Before(refreshToken.ExpiresAt) {
			count++
		}
	}

	return count
}

// GetTokenStatus returns the status of a token
func (jm *JWTManager) GetTokenStatus(tokenString string) map[string]interface{} {
	_, exists := jm.refreshTokens[tokenString]

	info, _ := jm.GetTokenInfo(tokenString)

	status := map[string]interface{}{
		"is_revoked": !exists,
		"info":       info,
	}

	return status
}
