package middleware

import (
	"net/http"
)

// CachingMiddlewareStack returns the recommended middleware chain for caching
// Usage in main.go:
//
//	router := http.NewServeMux()
//	router = middleware.CachingMiddlewareStack(router)
//	http.ListenAndServe(":8080", router)
func CachingMiddlewareStack(next http.Handler) http.Handler {
	// Apply middleware in order (executes in reverse order)
	// 1. Cache headers based on content type
	next = CacheMiddleware(next)

	// 2. Handle conditional requests (ETag, If-None-Match)
	next = ConditionalRequestMiddleware(next)

	// 3. Add security headers
	next = SecurityHeadersMiddleware(next)

	// 4. Add performance headers
	next = PerformanceHeadersMiddleware(next)

	// 5. Compression support
	next = CompressThenCacheMiddleware(next)

	return next
}

// APICacheConfig returns cache configuration for a specific API endpoint
func APICacheConfig(pattern string) CacheConfig {
	// Override default API caching for specific patterns
	switch pattern {
	// User-specific data - no cache
	case "/api/v1/users/me":
		return CacheConfig{
			CacheControl:    "private, no-cache, no-store, must-revalidate",
			EnableETag:      true,
			SetLastModified: false,
		}

	// Public list endpoints - short cache (5 min)
	case "/api/v1/products":
		fallthrough
	case "/api/v1/categories":
		fallthrough
	case "/api/v1/components":
		return CacheConfig{
			CacheControl:    "public, max-age=300, must-revalidate",
			EnableETag:      true,
			SetLastModified: true,
		}

	// Theme data - medium cache (30 min)
	case "/api/v1/themes":
		return CacheConfig{
			CacheControl:    "public, max-age=1800, must-revalidate",
			EnableETag:      true,
			SetLastModified: true,
		}

	// Analytics data - short cache (1 min)
	case "/api/v1/analytics":
		return CacheConfig{
			CacheControl:    "private, max-age=60, must-revalidate",
			EnableETag:      true,
			SetLastModified: false,
		}

	// Notifications - no cache
	case "/api/v1/notifications":
		return CacheConfig{
			CacheControl:    "private, no-cache, no-store, must-revalidate",
			EnableETag:      true,
			SetLastModified: false,
		}

	// Default API response - no cache
	default:
		return CacheStrategies["api"]
	}
}

// CacheControlDirectives provides helper functions for building Cache-Control headers
type CacheControlDirectives struct {
	// Public or private
	Visibility string // "public" or "private"

	// How long to cache in seconds
	MaxAge int

	// How long can be used after expiration while validating
	SMaxAge int

	// Must revalidate with origin
	MustRevalidate bool

	// Can only be used if fresh
	NoCache bool

	// Cannot be stored
	NoStore bool

	// Can be stored but not used without validation
	ProxyRevalidate bool

	// Immutable (never changes)
	Immutable bool
}

// String converts directives to Cache-Control header value
func (d CacheControlDirectives) String() string {
	result := d.Visibility

	if d.MaxAge > 0 {
		result += ", max-age=" + string(rune(d.MaxAge))
	}

	if d.SMaxAge > 0 {
		result += ", s-maxage=" + string(rune(d.SMaxAge))
	}

	if d.NoCache {
		result += ", no-cache"
	}

	if d.NoStore {
		result += ", no-store"
	}

	if d.MustRevalidate {
		result += ", must-revalidate"
	}

	if d.ProxyRevalidate {
		result += ", proxy-revalidate"
	}

	if d.Immutable {
		result += ", immutable"
	}

	return result
}

// CommonCachePatterns provides pre-configured cache patterns
var CommonCachePatterns = map[string]string{
	// Immutable content (versioned/hashed filenames)
	"immutable": "public, max-age=31536000, immutable",

	// Long-term cache with revalidation
	"longterm": "public, max-age=31536000, must-revalidate",

	// Medium-term cache
	"mediumterm": "public, max-age=604800, must-revalidate", // 1 week

	// Short-term cache
	"shortterm": "public, max-age=3600, must-revalidate", // 1 hour

	// Very short-term cache
	"veryshort": "public, max-age=300, must-revalidate", // 5 minutes

	// Private, non-cached
	"nocache": "private, no-cache, no-store, must-revalidate",

	// Stale while revalidate (useful for service workers)
	"stalerevalidate": "public, max-age=300, stale-while-revalidate=86400", // 5 min cache, 1 day stale
}

// HTTPCacheHeaders contains common HTTP cache-related headers
type HTTPCacheHeaders struct {
	CacheControl    string
	ETag            string
	LastModified    string
	Expires         string
	Vary            string
	Age             string
	Pragma          string
	IfNoneMatch     string
	IfModifiedSince string
	Expires_RFC1123 string
	Expires_RFC2822 string
}

// GetCacheHeaders returns recommended headers for different scenarios
func GetCacheHeaders(scenario string) HTTPCacheHeaders {
	switch scenario {
	case "static":
		return HTTPCacheHeaders{
			CacheControl:    "public, max-age=31536000, immutable",
			Vary:            "Accept-Encoding",
			Expires_RFC1123: "Thu, 31 Dec 2099 23:59:59 GMT",
		}

	case "html":
		return HTTPCacheHeaders{
			CacheControl: "public, max-age=3600, must-revalidate",
			Vary:         "Accept-Encoding, Accept",
		}

	case "api":
		return HTTPCacheHeaders{
			CacheControl: "private, no-cache, no-store, must-revalidate",
			Pragma:       "no-cache",
			Expires:      "-1",
		}

	case "manifest":
		return HTTPCacheHeaders{
			CacheControl: "public, max-age=3600, must-revalidate",
		}

	default:
		return HTTPCacheHeaders{
			CacheControl: "private, no-cache, no-store, must-revalidate",
		}
	}
}
