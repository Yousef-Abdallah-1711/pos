package middleware

import (
	"net/http"
	"strings"
)

func MaxBodyBytesMiddleware(limit int64) func(http.Handler) http.Handler {
	if limit <= 0 {
		limit = 1 << 20 // 1MB
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				r.Body = http.MaxBytesReader(w, r.Body, limit)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// MaxJSONBodyBytesMiddleware limits only JSON request bodies.
// It avoids breaking multipart/form-data endpoints such as image uploads.
func MaxJSONBodyBytesMiddleware(limit int64) func(http.Handler) http.Handler {
	if limit <= 0 {
		limit = 1 << 20 // 1MB
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
				if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "+json") {
					r.Body = http.MaxBytesReader(w, r.Body, limit)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
