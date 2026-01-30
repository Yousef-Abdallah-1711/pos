package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

// SuperAdminContextKey is the context key for super admin data
type superAdminContextKeyType string

const (
	SuperAdminContextKey     superAdminContextKeyType = "superAdmin"
	SuperAdminIDContextKey   superAdminContextKeyType = "superAdminID"
	SuperAdminUserContextKey superAdminContextKeyType = "superAdminUser"
)

// SuperAdminClaims is an alias to jwt.SuperAdminClaims for convenience
type SuperAdminClaims = jwt.SuperAdminClaims

// SuperAdminTokenValidator interface for validating super admin tokens
type SuperAdminTokenValidator interface {
	ValidateSuperAdminToken(token string) (*jwt.SuperAdminClaims, error)
}

// SuperAdminAuthMiddleware validates super admin JWT tokens
func SuperAdminAuthMiddleware(validator SuperAdminTokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[SuperAdminAuthMiddleware] Processing request: %s %s\n", r.Method, r.URL.Path)

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("[SuperAdminAuthMiddleware] No Authorization header\n")
				respondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Printf("[SuperAdminAuthMiddleware] Invalid Authorization header format\n")
				respondError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			claims, err := validator.ValidateSuperAdminToken(token)
			if err != nil {
				log.Printf("[SuperAdminAuthMiddleware] Token validation failed: %v\n", err)
				respondError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Check if token is expired
			if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
				log.Printf("[SuperAdminAuthMiddleware] Token expired\n")
				respondError(w, http.StatusUnauthorized, "token expired")
				return
			}

			// Check if this is a super admin token
			if !claims.IsSuperAdmin {
				log.Printf("[SuperAdminAuthMiddleware] Token is not for super admin\n")
				respondError(w, http.StatusForbidden, "access denied")
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), SuperAdminContextKey, claims)
			ctx = context.WithValue(ctx, SuperAdminIDContextKey, claims.SuperAdminID)

			log.Printf("[SuperAdminAuthMiddleware] Authenticated super admin: %d\n", claims.SuperAdminID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireSuperAdmin ensures the request is from an active super admin
func RequireSuperAdmin(repo *repository.SuperAdminRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[RequireSuperAdmin] Validating super admin status\n")

			// Get super admin ID from context
			superAdminID, ok := r.Context().Value(SuperAdminIDContextKey).(int64)
			if !ok || superAdminID == 0 {
				log.Printf("[RequireSuperAdmin] No super admin ID in context\n")
				respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// Verify super admin is active
			superAdmin, err := repo.GetByID(r.Context(), superAdminID)
			if err != nil {
				log.Printf("[RequireSuperAdmin] Super admin not found: %d\n", superAdminID)
				respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if !superAdmin.IsActive {
				log.Printf("[RequireSuperAdmin] Super admin is inactive: %d\n", superAdminID)
				respondError(w, http.StatusForbidden, "account is inactive")
				return
			}

			if superAdmin.IsLocked() {
				log.Printf("[RequireSuperAdmin] Super admin is locked: %d\n", superAdminID)
				respondError(w, http.StatusForbidden, "account is locked")
				return
			}

			// Add super admin to context
			ctx := context.WithValue(r.Context(), SuperAdminUserContextKey, superAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateIPWhitelist checks if the request IP is in the super admin's whitelist
func ValidateIPWhitelist(repo *repository.SuperAdminRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[ValidateIPWhitelist] Validating IP whitelist\n")

			// Get super admin ID from context
			superAdminID, ok := r.Context().Value(SuperAdminIDContextKey).(int64)
			if !ok || superAdminID == 0 {
				next.ServeHTTP(w, r) // Let other middleware handle auth
				return
			}

			// Get client IP
			ip := GetClientIP(r)

			// Get super admin
			superAdmin, err := repo.GetByID(r.Context(), superAdminID)
			if err != nil {
				log.Printf("[ValidateIPWhitelist] Super admin not found: %d\n", superAdminID)
				respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// Check IP whitelist
			if !superAdmin.IsIPAllowed(ip) {
				log.Printf("[ValidateIPWhitelist] IP not allowed: %s for super admin %d\n", ip, superAdminID)
				respondError(w, http.StatusForbidden, "access denied from this IP address")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitSuperAdmin applies stricter rate limiting for super admin endpoints
func RateLimitSuperAdmin(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter
	// In production, use Redis or similar
	type clientInfo struct {
		count    int
		resetAt  time.Time
	}
	clients := make(map[string]*clientInfo)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)

			// Get or create client info
			client, exists := clients[ip]
			if !exists || time.Now().After(client.resetAt) {
				clients[ip] = &clientInfo{
					count:   1,
					resetAt: time.Now().Add(time.Minute),
				}
			} else {
				client.count++
				if client.count > requestsPerMinute {
					log.Printf("[RateLimitSuperAdmin] Rate limit exceeded for IP: %s\n", ip)
					w.Header().Set("Retry-After", "60")
					respondError(w, http.StatusTooManyRequests, "rate limit exceeded")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClientIP extracts the real client IP from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP in the list
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetSuperAdminClaims retrieves super admin claims from context
func GetSuperAdminClaims(r *http.Request) *SuperAdminClaims {
	claims, ok := r.Context().Value(SuperAdminContextKey).(*SuperAdminClaims)
	if !ok {
		return nil
	}
	return claims
}

// GetSuperAdminID retrieves super admin ID from context
func GetSuperAdminID(r *http.Request) int64 {
	id, ok := r.Context().Value(SuperAdminIDContextKey).(int64)
	if !ok {
		return 0
	}
	return id
}

// Helper function to respond with error
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   message,
		"success": false,
	})
}
