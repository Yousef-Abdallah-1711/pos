package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CSRFCookieName  = "csrf_token"
	CSRFHeaderName  = "X-CSRF-Token"
	CSRFTokenLength = 32
)

func GenerateCSRFToken() (string, error) {
	buf := make([]byte, CSRFTokenLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func SetCSRFCookie(w http.ResponseWriter, r *http.Request, token string) {
	isSecure := r.TLS != nil || strings.EqualFold(os.Getenv("ENV"), "production")
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // frontend reads it to set X-CSRF-Token header
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(CSRFCookieName)
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			http.Error(w, "Forbidden - missing CSRF cookie", http.StatusForbidden)
			return
		}
		headerToken := strings.TrimSpace(r.Header.Get(CSRFHeaderName))
		if headerToken == "" || headerToken != cookie.Value {
			http.Error(w, "Forbidden - invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
