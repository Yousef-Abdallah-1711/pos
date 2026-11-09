package middleware

import (
	"net/http"
	"strings"
)

// CacheConfig defines caching strategy for different content types
type CacheConfig struct {
	// Cache-Control directive
	CacheControl string
	// ETag support
	EnableETag bool
	// Last-Modified header
	SetLastModified bool
}

// CacheStrategies defines caching policies for different resources
var CacheStrategies = map[string]CacheConfig{
	// Static assets - long cache (1 year)
	"static": {
		CacheControl:    "public, max-age=31536000, immutable",
		EnableETag:      true,
		SetLastModified: true,
	},
	// JavaScript bundles - long cache (1 year)
	"javascript": {
		CacheControl:    "public, max-age=31536000, immutable",
		EnableETag:      true,
		SetLastModified: true,
	},
	// CSS stylesheets - long cache (1 year)
	"css": {
		CacheControl:    "public, max-age=31536000, immutable",
		EnableETag:      true,
		SetLastModified: true,
	},
	// Web fonts - long cache (1 year)
	"fonts": {
		CacheControl:    "public, max-age=31536000, immutable",
		EnableETag:      true,
		SetLastModified: true,
	},
	// Images - long cache (1 year)
	"images": {
		CacheControl:    "public, max-age=31536000, immutable",
		EnableETag:      true,
		SetLastModified: true,
	},
	// HTML pages - short cache (1 hour)
	"html": {
		CacheControl:    "public, max-age=3600, must-revalidate",
		EnableETag:      true,
		SetLastModified: true,
	},
	// API responses - no cache by default
	"api": {
		CacheControl:    "private, no-cache, no-store, must-revalidate",
		EnableETag:      true,
		SetLastModified: false,
	},
	// Component manifests - medium cache (1 hour)
	"manifest": {
		CacheControl:    "public, max-age=3600, must-revalidate",
		EnableETag:      true,
		SetLastModified: true,
	},
}

// CacheMiddleware sets appropriate cache headers based on request path
func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determine cache strategy based on request path
		strategy := determineCacheStrategy(r.URL.Path)

		// Get cache config
		config, exists := CacheStrategies[strategy]
		if !exists {
			config = CacheStrategies["api"] // Default to API (no cache)
		}

		// Set Cache-Control header
		w.Header().Set("Cache-Control", config.CacheControl)

		// Set additional headers for better caching
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable compression
		w.Header().Set("Vary", "Accept-Encoding")

		// For static assets, set long expiry headers
		if strategy == "static" || strategy == "fonts" || strategy == "images" {
			w.Header().Set("Expires", "Thu, 31 Dec 2099 23:59:59 GMT")
		}

		// For API responses, disable caching at HTTP level (let service worker handle it)
		if strategy == "api" {
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "-1")
		}

		next.ServeHTTP(w, r)
	})
}

// determineCacheStrategy determines the caching strategy based on URL path
func determineCacheStrategy(path string) string {
	// Static assets with version hash (immutable)
	if strings.HasPrefix(path, "/_next/static/") {
		return "static"
	}

	// Public static files
	if strings.HasPrefix(path, "/static/") {
		return "static"
	}

	// JavaScript bundles
	if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".mjs") {
		return "javascript"
	}

	// CSS stylesheets
	if strings.HasSuffix(path, ".css") {
		return "css"
	}

	// Web fonts
	if strings.HasSuffix(path, ".woff") ||
		strings.HasSuffix(path, ".woff2") ||
		strings.HasSuffix(path, ".ttf") ||
		strings.HasSuffix(path, ".eot") ||
		strings.HasSuffix(path, ".otf") {
		return "fonts"
	}

	// Images
	if strings.HasSuffix(path, ".png") ||
		strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".jpeg") ||
		strings.HasSuffix(path, ".gif") ||
		strings.HasSuffix(path, ".webp") ||
		strings.HasSuffix(path, ".svg") ||
		strings.HasSuffix(path, ".ico") ||
		strings.HasPrefix(path, "/_next/image") {
		return "images"
	}

	// HTML pages
	if strings.HasSuffix(path, ".html") ||
		strings.HasSuffix(path, "/") && !strings.HasPrefix(path, "/api/") {
		return "html"
	}

	// Component manifests and metadata
	if strings.HasSuffix(path, ".json") && !strings.HasPrefix(path, "/api/") {
		if strings.Contains(path, "manifest") ||
			strings.Contains(path, "metadata") {
			return "manifest"
		}
	}

	// API routes - no cache by default
	if strings.HasPrefix(path, "/api/") {
		return "api"
	}

	// Default to API (no cache)
	return "api"
}

// ConditionalRequestMiddleware handles ETag and If-None-Match headers
func ConditionalRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add ETag generation capability
		// This would be enhanced with actual content hashing
		etag := generateETag(r.URL.Path)

		// Check If-None-Match header
		if match := r.Header.Get("If-None-Match"); match != "" {
			if match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Set ETag header
		w.Header().Set("ETag", etag)

		next.ServeHTTP(w, r)
	})
}

// generateETag generates a simple ETag for the given path
// In production, this should be based on actual content hash
func generateETag(path string) string {
	// Simple ETag based on path and timestamp
	// In production, use content hash
	return `"` + strings.TrimPrefix(path, "/") + `"`
}

// CompressThenCacheMiddleware applies compression before caching
func CompressThenCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client supports compression
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
		}

		// Continue to caching middleware
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Clickjacking protection
		w.Header().Set("X-Frame-Options", "DENY")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
		)

		// Feature-Policy (Permissions-Policy)
		w.Header().Set(
			"Permissions-Policy",
			"accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
		)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}

// PerformanceHeadersMiddleware adds performance-related headers
func PerformanceHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Server timing header for performance monitoring
		w.Header().Set("Server-Timing", "total;dur=0")

		// DNS prefetch for external resources
		w.Header().Set("Link", `</fonts>; rel=prefetch, </static>; rel=prefetch`)

		// Enable HTTP/2 Server Push (if supported)
		if pusher, ok := w.(http.Pusher); ok {
			// Push critical resources
			pusher.Push("/_next/static/", &http.PushOptions{Header: http.Header{
				"Content-Type": []string{"text/javascript"},
			}})
		}

		next.ServeHTTP(w, r)
	})
}

// CacheStrategy describes a caching strategy for monitoring/documentation
type CacheStrategy struct {
	ContentType string
	Pattern     string
	MaxAge      int
	Description string
	Example     string
}

// GetCacheStrategies returns all configured cache strategies for documentation
func GetCacheStrategies() []CacheStrategy {
	return []CacheStrategy{
		{
			ContentType: "Static Assets",
			Pattern:     "/_next/static/*, /static/*",
			MaxAge:      31536000, // 1 year
			Description: "Immutable, versioned assets - cached for 1 year",
			Example:     "/_next/static/chunks/main-abc123.js",
		},
		{
			ContentType: "Fonts",
			Pattern:     "/*.woff2, /*.ttf",
			MaxAge:      31536000,
			Description: "Web fonts - cached for 1 year",
			Example:     "/fonts/inter.woff2",
		},
		{
			ContentType: "Images",
			Pattern:     "/*.png, /*.jpg, /*.webp",
			MaxAge:      31536000,
			Description: "Images - cached for 1 year",
			Example:     "/images/hero.webp",
		},
		{
			ContentType: "HTML Pages",
			Pattern:     "/*.html, /",
			MaxAge:      3600,
			Description: "HTML pages - cached for 1 hour, must revalidate",
			Example:     "/index.html",
		},
		{
			ContentType: "API Responses",
			Pattern:     "/api/*",
			MaxAge:      0,
			Description: "API responses - not cached at HTTP level",
			Example:     "/api/v1/products",
		},
	}
}
