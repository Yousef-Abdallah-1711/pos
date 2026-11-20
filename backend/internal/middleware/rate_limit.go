package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"pos-saas/internal/observability"
)

type rateLimitClient struct {
	count   int
	resetAt time.Time
}

type rateLimitKeyFunc func(*http.Request) string

// SharedRateLimitStore atomically increments an expiring counter in a shared
// store such as Redis.
type SharedRateLimitStore interface {
	IncrementWindow(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error)
}

// RateLimitByIP applies an in-memory per-IP rate limit windowed by one minute.
func RateLimitByIP(requestsPerMinute int) func(http.Handler) http.Handler {
	return rateLimitWithKey(requestsPerMinute, time.Minute, func(r *http.Request) string {
		return "ip:" + GetClientIP(r)
	})
}

// RateLimitByIPWithStore applies an atomic shared per-IP window. The scope and
// IP are hashed before persistence so Redis keys contain no client addresses.
func RateLimitByIPWithStore(store SharedRateLimitStore, scope string, limit int, window time.Duration) func(http.Handler) http.Handler {
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}
	scope = strings.TrimSpace(scope)
	if scope == "" {
		scope = "auth"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if store == nil {
				recordRateLimitStoreError()
				writeRateLimitStoreUnavailable(w)
				return
			}
			key := sharedRateLimitKey(scope, GetClientIP(r))
			count, ttl, err := store.IncrementWindow(r.Context(), key, window)
			if err != nil || count < 1 {
				recordRateLimitStoreError()
				writeRateLimitStoreUnavailable(w)
				return
			}
			if ttl <= 0 {
				ttl = time.Second
			}
			resetAt := time.Now().Add(ttl)
			remaining := int64(limit) - count
			if remaining < 0 {
				remaining = 0
			}
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))
			if count > int64(limit) {
				recordRateLimitRejected()
				retryAfter := int(ttl / time.Second)
				if ttl%time.Second != 0 {
					retryAfter++
				}
				if retryAfter < 1 {
					retryAfter = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.Header().Set("Cache-Control", "no-store")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "rate limit exceeded"})
				return
			}
			recordRateLimitAllowed()
			next.ServeHTTP(w, r)
		})
	}
}

func sharedRateLimitKey(scope, clientIP string) string {
	sum := sha256.Sum256([]byte(scope + "\x00" + clientIP))
	return "pos-saas:auth:rate-limit:v1:" + hex.EncodeToString(sum[:])
}

func writeRateLimitStoreUnavailable(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Retry-After", "1")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "service temporarily unavailable"})
}

func recordRateLimitAllowed() {
	observability.RecordAuthSecurity(func(metrics *observability.AuthSecurityMetrics) { metrics.RateLimitAllowed++ })
}

func recordRateLimitRejected() {
	observability.RecordAuthSecurity(func(metrics *observability.AuthSecurityMetrics) { metrics.RateLimitRejected++ })
}

func recordRateLimitStoreError() {
	observability.RecordAuthSecurity(func(metrics *observability.AuthSecurityMetrics) { metrics.RateLimitStoreErrors++ })
}

// RateLimitByUser applies an in-memory per-user rate limit windowed by one minute.
// If user claims are unavailable, it falls back to IP address.
func RateLimitByUser(requestsPerMinute int) func(http.Handler) http.Handler {
	return RateLimitByUserWindow(requestsPerMinute, time.Minute)
}

// RateLimitByUserWindow applies an in-memory per-user limit for a custom window.
// If user claims are unavailable, it falls back to IP address.
func RateLimitByUserWindow(limit int, window time.Duration) func(http.Handler) http.Handler {
	return rateLimitWithKey(limit, window, getUserRateLimitKey)
}

func getUserRateLimitKey(r *http.Request) string {
	userID := GetUserID(r)
	if userID > 0 {
		return "user:" + strconv.FormatInt(userID, 10)
	}

	return "ip:" + GetClientIP(r)
}

func rateLimitWithKey(limit int, window time.Duration, keyFn rateLimitKeyFunc) func(http.Handler) http.Handler {
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	clients := make(map[string]rateLimitClient)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFn(r)
			if key == "" {
				key = "global"
			}

			now := time.Now()

			mu.Lock()
			client, exists := clients[key]
			if !exists || now.After(client.resetAt) {
				client = rateLimitClient{
					count:   0,
					resetAt: now.Add(window),
				}
			}

			if client.count >= limit {
				retryAfter := int(time.Until(client.resetAt).Seconds())
				if retryAfter < 1 {
					retryAfter = 1
				}
				mu.Unlock()

				w.Header().Set("Retry-After", intToString(retryAfter))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "rate limit exceeded",
				})
				return
			}

			client.count++
			clients[key] = client
			remaining := limit - client.count
			resetAt := client.resetAt
			mu.Unlock()

			w.Header().Set("X-RateLimit-Limit", intToString(limit))
			w.Header().Set("X-RateLimit-Remaining", intToString(remaining))
			w.Header().Set("X-RateLimit-Reset", intToString(int(resetAt.Unix())))

			next.ServeHTTP(w, r)
		})
	}
}

func intToString(value int) string {
	return strconv.Itoa(value)
}
