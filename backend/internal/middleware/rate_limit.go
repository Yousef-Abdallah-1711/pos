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

// RateLimitByIP applies an in-memory per-IP rate limit windowed by one minute.
func RateLimitByIP(requestsPerMinute int) func(http.Handler) http.Handler {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60
	}

	clients := make(map[string]rateLimitClient)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)
			now := time.Now()

			mu.Lock()
			client, exists := clients[ip]
			if !exists || now.After(client.resetAt) {
				client = rateLimitClient{
					count:   0,
					resetAt: now.Add(time.Minute),
				}
			}

			if client.count >= requestsPerMinute {
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
			clients[ip] = client
			remaining := requestsPerMinute - client.count
			resetAt := client.resetAt
			mu.Unlock()

			w.Header().Set("X-RateLimit-Limit", intToString(requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", intToString(remaining))
			w.Header().Set("X-RateLimit-Reset", intToString(int(resetAt.Unix())))

			next.ServeHTTP(w, r)
		})
	}
}

func intToString(value int) string {
	return strconv.Itoa(value)
}
