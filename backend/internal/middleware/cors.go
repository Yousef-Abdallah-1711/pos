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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-Restaurant-ID, X-Impersonation-Token")
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
	// Keep local development flexible for alternate ports only outside production.
	return config.IsDevelopmentMode() && isLocalhostOrigin(origin)
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
