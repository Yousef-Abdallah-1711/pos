package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID       int    `json:"user_id"`
	Email        string `json:"email"`
	TenantID     int    `json:"tenant_id"`
	RestaurantID int    `json:"restaurant_id,omitempty"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
	expiry time.Duration
}

func NewTokenService(secret string, expiry string) (*TokenService, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return nil, err
	}

	return &TokenService{
		secret: []byte(secret),
		expiry: duration,
	}, nil
}

func (s *TokenService) GenerateToken(userID, tenantID int, restaurantID *int, email, role string) (string, error) {
	restID := 0
	if restaurantID != nil {
		restID = *restaurantID
	}

	claims := Claims{
		UserID:       userID,
		Email:        email,
		TenantID:     tenantID,
		RestaurantID: restID,
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Login-selection tokens are proofs for one specific step, never sessions.
		for _, aud := range claims.Audience {
			if aud == selectionAudience {
				return nil, errors.New("invalid token")
			}
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

const selectionAudience = "pos-saas:login-selection"

// SelectionClaims is a short-lived proof that the holder already passed password
// verification for Email and may choose one of TenantIDs.
type SelectionClaims struct {
	Email     string `json:"email"`
	TenantIDs []int  `json:"tenant_ids"`
	jwt.RegisteredClaims
}

// AllowsTenant reports whether the proof covers the tenant.
func (c *SelectionClaims) AllowsTenant(tenantID int) bool {
	for _, id := range c.TenantIDs {
		if id == tenantID {
			return true
		}
	}
	return false
}

// GenerateSelectionToken issues the proof returned by a successful multi-tenant login.
func (s *TokenService) GenerateSelectionToken(email string, tenantIDs []int, ttl time.Duration) (string, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}
	now := time.Now()
	claims := SelectionClaims{
		Email:     strings.ToLower(strings.TrimSpace(email)),
		TenantIDs: tenantIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        hex.EncodeToString(idBytes),
			Audience:  jwt.ClaimStrings{selectionAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

// ValidateSelectionToken verifies signature, expiry and audience.
func (s *TokenService) ValidateSelectionToken(tokenString string) (*SelectionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SelectionClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithAudience(selectionAudience))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*SelectionClaims)
	if !ok || !token.Valid || claims.ID == "" {
		return nil, errors.New("invalid selection token")
	}
	return claims, nil
}
