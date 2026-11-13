package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFrameAncestorsFromConfig(t *testing.T) {
	tests := []struct {
		name        string
		configured  string
		development bool
		want        string
	}{
		{name: "missing configuration denies framing", want: "'none'"},
		{name: "HTTPS origins are accepted", configured: "https://dashboard.example.test,https://admin.example.test", want: "https://dashboard.example.test https://admin.example.test"},
		{name: "loopback HTTP is development-only", configured: "http://127.0.0.1:3002", development: true, want: "http://127.0.0.1:3002"},
		{name: "HTTP origin is rejected outside development", configured: "http://127.0.0.1:3002", want: "'none'"},
		{name: "remote HTTP origin is rejected", configured: "http://dashboard.example.test", development: true, want: "'none'"},
		{name: "paths are rejected", configured: "https://dashboard.example.test/checkout", want: "'none'"},
		{name: "wildcards are rejected", configured: "https://*.example.test", want: "'none'"},
		{name: "invalid member rejects entire policy", configured: "https://dashboard.example.test,not-a-url", want: "'none'"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := frameAncestorsFromConfig(test.configured, test.development); got != test.want {
				t.Fatalf("frameAncestorsFromConfig() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSecurityHeadersMiddlewareUsesConfiguredCheckoutAncestors(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv(mockCheckoutFrameAncestorsEnv, "https://dashboard.example.test")

	handler := SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/v1/mock/paymob/iframe", nil)
	handler.ServeHTTP(recorder, request)

	policy := recorder.Header().Get("Content-Security-Policy")
	if !strings.Contains(policy, "frame-ancestors https://dashboard.example.test") {
		t.Fatalf("CSP does not include configured ancestor: %q", policy)
	}
	if got := recorder.Header().Get("X-Frame-Options"); got != "" {
		t.Fatalf("X-Frame-Options = %q, want unset for explicit CSP allowlist", got)
	}
}

func TestSecurityHeadersMiddlewareDeniesCheckoutFramingWithoutConfig(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv(mockCheckoutFrameAncestorsEnv, "")

	handler := SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/v1/mock/paymob/iframe", nil)
	handler.ServeHTTP(recorder, request)

	if policy := recorder.Header().Get("Content-Security-Policy"); !strings.Contains(policy, "frame-ancestors 'none'") {
		t.Fatalf("CSP = %q, want frame-ancestors 'none'", policy)
	}
	if got := recorder.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("X-Frame-Options = %q, want DENY", got)
	}
}
