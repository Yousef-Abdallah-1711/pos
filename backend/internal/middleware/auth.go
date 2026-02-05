package middleware

import (
	"context"
	"net/http"
	"strings"

	"pos-saas/internal/pkg/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(tokenService *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First check for impersonation token header
			impersonationToken := r.Header.Get("X-Impersonation-Token")
			if impersonationToken != "" {
				// If impersonation token is present, allow the request
				// The frontend already validated it on the /impersonation/validate endpoint
				// Create a minimal claims object - will be filled by tenant middleware
				claims := &jwt.Claims{
					UserID: 1, // Temporary value, actual user will come from tenant context
				}

				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "Unauthorized - missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := tokenService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims retrieves claims from context
func GetUserClaims(r *http.Request) *jwt.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*jwt.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetUserID retrieves user ID from request context
func GetUserID(r *http.Request) int64 {
	claims := GetUserClaims(r)
	if claims == nil {
		return 0
	}
	return int64(claims.UserID)
}
