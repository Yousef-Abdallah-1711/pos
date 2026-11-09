package middleware

import (
	"context"
	"net/http"
	"strings"

	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

type contextKey string

const UserContextKey contextKey = "user"

// ImpersonationAuthenticator validates a presented impersonation token against
// server-side session state and returns the identity of the impersonated user.
// A nil error means the token is genuine, active and unexpired.
type ImpersonationAuthenticator interface {
	AuthenticateImpersonation(ctx context.Context, token string) (*jwt.Claims, *ImpersonationInfo, error)
}

type authSessionResolver interface {
	Resolve(context.Context, string) (string, error)
}

// SessionAuthMiddleware resolves the host-only HttpOnly session cookie through
// shared server-side storage, then validates the stored JWT with the existing
// authorization middleware. State-changing cookie requests require a trusted
// Origin to prevent cross-site request forgery.
func SessionAuthMiddleware(tokenService *jwt.TokenService, sessions authSessionResolver, impersonation ...ImpersonationAuthenticator) func(http.Handler) http.Handler {
	legacyAuth := AuthMiddleware(tokenService, impersonation...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(service.ImpersonationSessionCookieName)
			isImpersonation := err == nil && cookie.Value != ""
			if !isImpersonation {
				cookie, err = r.Cookie(service.AuthSessionCookieName)
			}
			if err != nil || cookie.Value == "" {
				legacyAuth(next).ServeHTTP(w, r)
				return
			}
			if isUnsafeMethod(r.Method) && !IsTrustedOrigin(r.Header.Get("Origin")) {
				http.Error(w, "Forbidden - invalid request origin", http.StatusForbidden)
				return
			}
			credential, err := sessions.Resolve(r.Context(), cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			request := r.Clone(r.Context())
			request.Header = r.Header.Clone()
			if isImpersonation && service.IsImpersonationCredential(credential) {
				request.Header.Set("X-Impersonation-Token", service.UnwrapImpersonationCredential(credential))
			} else {
				request.Header.Set("Authorization", "Bearer "+credential)
			}
			legacyAuth(next).ServeHTTP(w, request)
		})
	}
}

func isUnsafeMethod(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}

// AuthMiddleware authenticates requests with a bearer JWT. When an
// ImpersonationAuthenticator is supplied, an X-Impersonation-Token header is
// accepted only if that authenticator validates it. The mere presence of the
// header never authenticates a request, and an invalid impersonation token is
// never silently replaced by another credential.
func AuthMiddleware(tokenService *jwt.TokenService, impersonation ...ImpersonationAuthenticator) func(http.Handler) http.Handler {
	var impersonationAuth ImpersonationAuthenticator
	if len(impersonation) > 0 {
		impersonationAuth = impersonation[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if impersonationToken := r.Header.Get("X-Impersonation-Token"); impersonationToken != "" {
				if impersonationAuth == nil {
					http.Error(w, "Unauthorized - impersonation not supported", http.StatusUnauthorized)
					return
				}
				claims, info, err := impersonationAuth.AuthenticateImpersonation(r.Context(), impersonationToken)
				if err != nil || claims == nil {
					http.Error(w, "Unauthorized - invalid impersonation token", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				if info != nil {
					ctx = context.WithValue(ctx, ImpersonationContextKey, info)
					ctx = context.WithValue(ctx, IsImpersonatedContextKey, true)
					ctx = context.WithValue(ctx, OriginalSuperAdminIDKey, info.SuperAdminID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authHeader := r.Header.Get("Authorization")
			headerToken := ""
			if authHeader != "" {
				// Extract token from "Bearer <token>"
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					headerToken = strings.TrimSpace(parts[1])
				}
			}
			cookieToken := ""
			if cookie, err := r.Cookie("auth_token"); err == nil {
				cookieToken = strings.TrimSpace(cookie.Value)
			}

			if headerToken == "" && cookieToken == "" {
				http.Error(w, "Unauthorized - missing auth token", http.StatusUnauthorized)
				return
			}

			claims, err := tokenService.ValidateToken(headerToken)
			if err != nil && cookieToken != "" && cookieToken != headerToken {
				claims, err = tokenService.ValidateToken(cookieToken)
			}
			if err != nil || claims == nil {
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
