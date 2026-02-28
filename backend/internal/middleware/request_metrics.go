package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/observability"
)

type statusCaptureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCaptureResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func RequestMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Websocket upgrades require the original response writer for hijacking.
		// Wrapping can break the handshake if Hijacker is not preserved.
		if isUpgradeRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		startedAt := time.Now()
		statusWriter := &statusCaptureResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(statusWriter, r)

		pattern := r.Pattern
		if pattern == "" {
			pattern = r.URL.Path
		}
		metric := r.Method + " " + pattern + " [" + strconv.Itoa(statusWriter.statusCode) + "]"
		observability.RecordAPI(metric, time.Since(startedAt))
	})
}

func isUpgradeRequest(r *http.Request) bool {
	if r == nil {
		return false
	}

	connection := strings.ToLower(strings.TrimSpace(r.Header.Get("Connection")))
	upgrade := strings.ToLower(strings.TrimSpace(r.Header.Get("Upgrade")))
	return strings.Contains(connection, "upgrade") || upgrade != ""
}
