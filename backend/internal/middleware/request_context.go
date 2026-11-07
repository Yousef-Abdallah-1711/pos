package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

const RequestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

var requestIDFallbackCounter atomic.Uint64

// RequestIDMiddleware assigns a server-generated ID to every request. It deliberately
// ignores caller-supplied IDs so untrusted values cannot forge log correlation fields.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		ctx := context.WithValue(r.Context(), requestIDContextKey{}, requestID)
		w.Header().Set(RequestIDHeader, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the ID assigned by RequestIDMiddleware, or an empty
// string when the request did not pass through that middleware.
func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

// RecoveryMiddleware prevents handler panics from terminating the server and keeps
// panic values and stack traces out of both the response and logs.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() == nil {
				return
			}

			slog.Error("request panic recovered",
				"request_id", RequestIDFromContext(r.Context()),
				"method", r.Method,
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	var id [16]byte
	if _, err := rand.Read(id[:]); err == nil {
		id[6] = (id[6] & 0x0f) | 0x40
		id[8] = (id[8] & 0x3f) | 0x80
		encoded := hex.EncodeToString(id[:])
		return fmt.Sprintf("%s-%s-%s-%s-%s", encoded[:8], encoded[8:12], encoded[12:16], encoded[16:20], encoded[20:])
	}

	return fmt.Sprintf("fallback-%d", requestIDFallbackCounter.Add(1))
}
