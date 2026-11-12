package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SuperAdminClaims represents the claims in a super admin JWT
// This is the canonical definition - middleware should import this
type SuperAdminClaims struct {
	UserID        int64  `json:"user_id"`
	Email         string `json:"email"`
	SuperAdminID  int64  `json:"super_admin_id"`
	IsSuperAdmin  bool   `json:"is_super_admin"`
	ExpiresAtUnix int64  `json:"-"`
	jwt.RegisteredClaims
}

// SuperAdminTokenService handles super admin JWT tokens
type SuperAdminTokenService struct {
	secret []byte
	expiry time.Duration
}

// NewSuperAdminTokenService creates a new super admin token service
// Uses a separate secret for super admin tokens for additional security
func NewSuperAdminTokenService(secret string, expiry string) (*SuperAdminTokenService, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		// Default to 15 minutes if not specified
		duration = 15 * time.Minute
	}

	return &SuperAdminTokenService{
		secret: []byte(secret),
		expiry: duration,
	}, nil
}

// GenerateSuperAdminToken generates a JWT token for super admin
func (s *SuperAdminTokenService) GenerateSuperAdminToken(userID int64, email string, superAdminID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.expiry)

	claims := SuperAdminClaims{
		UserID:       userID,
		Email:        email,
		SuperAdminID: superAdminID,
		IsSuperAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pos-saas-super-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// ValidateSuperAdminToken validates a super admin JWT token
func (s *SuperAdminTokenService) ValidateSuperAdminToken(tokenString string) (*SuperAdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SuperAdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SuperAdminClaims); ok && token.Valid {
		// Check if it's actually a super admin token
		if !claims.IsSuperAdmin {
			return nil, errors.New("not a super admin token")
		}

		// Check issuer
		if claims.Issuer != "pos-saas-super-admin" {
			return nil, errors.New("invalid token issuer")
		}

		// Set ExpiresAt from RegisteredClaims
		if claims.RegisteredClaims.ExpiresAt != nil {
			claims.ExpiresAtUnix = claims.RegisteredClaims.ExpiresAt.Unix()
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshSuperAdminToken refreshes a super admin token
func (s *SuperAdminTokenService) RefreshSuperAdminToken(tokenString string) (string, time.Time, error) {
	// Validate the existing token first
	claims, err := s.ValidateSuperAdminToken(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate a new token
	return s.GenerateSuperAdminToken(claims.UserID, claims.Email, claims.SuperAdminID)
}
