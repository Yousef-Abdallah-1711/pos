package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rateLimitClient struct {
	count   int
	resetAt time.Time
}

type rateLimitKeyFunc func(*http.Request) string

// RateLimitByIP applies an in-memory per-IP rate limit windowed by one minute.
func RateLimitByIP(requestsPerMinute int) func(http.Handler) http.Handler {
	return rateLimitWithKey(requestsPerMinute, time.Minute, func(r *http.Request) string {
		return "ip:" + GetClientIP(r)
	})
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
