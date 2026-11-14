package middleware

import (
	"net/http"
	"net/url"
	"os"
	"pos-saas/internal/config"
	"strings"
	"sync"
)

var (
	corsOriginsOnce sync.Once
	corsOrigins     map[string]struct{}
)

func CORSMiddleware(next http.Handler) http.Handler {
	corsOriginsOnce.Do(loadAllowedOrigins)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := normalizeOrigin(r.Header.Get("Origin"))
		originAllowed := isOriginAllowed(origin)

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-Restaurant-ID, X-Impersonation-Token, Idempotency-Key, X-Order-Token, X-Customer-Phone")
		w.Header().Set("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, Retry-After")
		w.Header().Add("Vary", "Origin")

		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			if origin != "" && !originAllowed {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loadAllowedOrigins() {
	corsOrigins = make(map[string]struct{})

	raw := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	if raw == "" && config.IsDevelopmentMode() {
		raw = "http://localhost:3001,http://localhost:3002,http://localhost:3003,http://localhost:3004"
	}

	for _, item := range strings.Split(raw, ",") {
		origin := normalizeOrigin(item)
		if origin == "" {
			continue
		}
		corsOrigins[origin] = struct{}{}
	}

	// Phase 8: allow Paymob hosted origin for iframe/payment communication.
	corsOrigins["https://accept.paymob.com"] = struct{}{}
}

func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	if _, ok := corsOrigins[origin]; ok {
		return true
	}
	parsedOrigin, err := url.Parse(origin)
	if err != nil || parsedOrigin.Host == "" || parsedOrigin.User != nil || parsedOrigin.Path != "" || parsedOrigin.RawQuery != "" || parsedOrigin.Fragment != "" {
		return false
	}
	originHost := strings.ToLower(parsedOrigin.Hostname())
	for configuredOrigin := range corsOrigins {
		pattern, err := url.Parse(configuredOrigin)
		if err != nil || !strings.HasPrefix(pattern.Hostname(), "*.") {
			continue
		}
		if parsedOrigin.Scheme != pattern.Scheme || parsedOrigin.Port() != pattern.Port() {
			continue
		}
		suffix := strings.TrimPrefix(strings.ToLower(pattern.Hostname()), "*")
		if !strings.HasSuffix(originHost, suffix) {
			continue
		}
		subdomain := strings.TrimSuffix(originHost, suffix)
		if subdomain != "" && !strings.Contains(subdomain, ".") {
			return true
		}
	}
	// Keep local development flexible for alternate ports only outside production.
	return config.IsDevelopmentMode() && isLocalhostOrigin(origin)
}

// IsTrustedOrigin reports whether an exact configured browser origin is allowed.
// Cookie-authenticated writes use it as an additional CSRF check.
func IsTrustedOrigin(origin string) bool {
	corsOriginsOnce.Do(loadAllowedOrigins)
	return isOriginAllowed(normalizeOrigin(origin))
}

// RequireTrustedOrigin protects browser endpoints that create or rotate sessions.
// Browsers include Origin on these POST requests; reject missing or unlisted origins.
func RequireTrustedOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsTrustedOrigin(r.Header.Get("Origin")) {
			http.Error(w, "Forbidden - invalid request origin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isLocalhostOrigin(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || host == "127.0.0.1" {
		return true
	}

	// Support local tenant subdomains in development (e.g. demo.localhost:3003).
	return strings.HasSuffix(host, ".localhost")
}

func normalizeOrigin(origin string) string {
	trimmed := strings.TrimSpace(origin)
	if trimmed == "" {
		return ""
	}
	return strings.TrimSuffix(trimmed, "/")
}
