package middleware

import (
	"net/http/httptest"
	"testing"
)

func TestGetClientIPWalksTrustedForwardedChainFromRight(t *testing.T) {
	t.Setenv("TRUSTED_PROXY_CIDRS", "10.0.0.0/8,192.0.2.0/24")
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "10.0.0.10:8080"
	request.Header.Set("X-Forwarded-For", "198.51.100.99, 203.0.113.8, 192.0.2.7")
	if got := GetClientIP(request); got != "203.0.113.8" {
		t.Fatalf("GetClientIP() = %q, want nearest untrusted client 203.0.113.8", got)
	}
}

func TestGetClientIPIgnoresForwardedHeadersFromUntrustedPeer(t *testing.T) {
	t.Setenv("TRUSTED_PROXY_CIDRS", "10.0.0.0/8")
	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "198.51.100.44:8080"
	request.Header.Set("X-Forwarded-For", "203.0.113.8")
	request.Header.Set("X-Real-IP", "203.0.113.9")
	if got := GetClientIP(request); got != "198.51.100.44" {
		t.Fatalf("GetClientIP() = %q, want remote peer 198.51.100.44", got)
	}
}
