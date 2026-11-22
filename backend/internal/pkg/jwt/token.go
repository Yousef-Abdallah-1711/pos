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
	if s == nil {
		return "", errors.New("token service is not configured")
	}
	return s.GenerateTokenWithExpiry(userID, tenantID, restaurantID, email, role, s.expiry)
}

// GenerateTokenWithExpiry issues an access credential with a caller-selected
// maximum lifetime. Session-backed browser credentials use the server session
// absolute deadline; bearer and other token callers retain the configured TTL.
func (s *TokenService) GenerateTokenWithExpiry(userID, tenantID int, restaurantID *int, email, role string, expiry time.Duration) (string, error) {
	if s == nil || expiry <= 0 {
		return "", errors.New("token service or expiry is invalid")
	}
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
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

// SelectionAccount identifies the exact legacy account whose password matched.
type SelectionAccount struct {
	UserID   int `json:"user_id"`
	TenantID int `json:"tenant_id"`
}

// SelectionClaims is a short-lived proof that the holder already passed password
// verification for Email and may choose one of the exact account/tenant pairs.
type SelectionClaims struct {
	Email    string             `json:"email"`
	Accounts []SelectionAccount `json:"accounts"`
	jwt.RegisteredClaims
}

// AllowsTenant reports whether the proof covers the tenant.
func (c *SelectionClaims) AllowsTenant(tenantID int) bool {
	for _, account := range c.Accounts {
		if account.TenantID == tenantID {
			return true
		}
	}
	return false
}

// AllowsAccount reports whether the proof covers this exact legacy account.
func (c *SelectionClaims) AllowsAccount(tenantID, userID int) bool {
	for _, account := range c.Accounts {
		if account.TenantID == tenantID && account.UserID == userID {
			return true
		}
	}
	return false
}

// GenerateSelectionToken issues the proof returned by a successful multi-tenant login.
func (s *TokenService) GenerateSelectionToken(email string, accounts []SelectionAccount, ttl time.Duration) (string, error) {
	if s == nil || len(accounts) == 0 || strings.TrimSpace(email) == "" {
		return "", errors.New("selection token requires a configured service, email, and account")
	}
	seenUsers := make(map[int]struct{}, len(accounts))
	for _, account := range accounts {
		if account.UserID <= 0 || account.TenantID <= 0 {
			return "", errors.New("selection token account identifiers must be positive")
		}
		if _, exists := seenUsers[account.UserID]; exists {
			return "", errors.New("selection token accounts must not repeat a user")
		}
		seenUsers[account.UserID] = struct{}{}
	}
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}
	now := time.Now()
	claims := SelectionClaims{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Accounts: append([]SelectionAccount(nil), accounts...),
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
	if !ok || !token.Valid || claims.ID == "" || len(claims.Accounts) == 0 || claims.Email == "" {
		return nil, errors.New("invalid selection token")
	}
	for _, account := range claims.Accounts {
		if account.UserID <= 0 || account.TenantID <= 0 {
			return nil, errors.New("invalid selection token account")
		}
	}
	return claims, nil
}
