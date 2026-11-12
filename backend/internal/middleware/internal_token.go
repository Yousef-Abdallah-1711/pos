package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

// RequireInternalToken protects operational endpoints (metrics, diagnostics) with a
// shared secret read from the environment variable envName.
//
//   - secret set: requests must present it as "Authorization: Bearer <secret>" or
//     "X-Internal-Token: <secret>" (constant-time comparison), otherwise 401;
//   - secret unset in production: the endpoint is disabled (404), never open;
//   - secret unset elsewhere: open, so local development keeps working.
func RequireInternalToken(envName string, production bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		secret := strings.TrimSpace(os.Getenv(envName))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				if production {
					http.NotFound(w, r)
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			presented := strings.TrimSpace(r.Header.Get("X-Internal-Token"))
			if presented == "" {
				if parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					presented = strings.TrimSpace(parts[1])
				}
			}
			if subtle.ConstantTimeCompare([]byte(presented), []byte(secret)) != 1 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
