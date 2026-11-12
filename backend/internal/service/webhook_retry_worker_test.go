package service

import (
	"strings"
	"testing"
)

func TestWebhookRetryWorkerResolveEndpointStaysOnConfiguredOrigin(t *testing.T) {
	worker := &WebhookRetryWorker{baseURL: "https://hooks.example.test"}
	tests := []struct {
		name       string
		endpoint   string
		want       string
		wantReject bool
	}{
		{name: "relative webhook path", endpoint: "/api/v1/webhook/paymob?hmac=signature", want: "https://hooks.example.test/api/v1/webhook/paymob?hmac=signature"},
		{name: "same-origin absolute callback", endpoint: "https://hooks.example.test/api/v1/webhook/paymob", want: "https://hooks.example.test/api/v1/webhook/paymob"},
		{name: "empty endpoint default", endpoint: "", want: "https://hooks.example.test/api/v1/webhook/paymob"},
		{name: "different host", endpoint: "https://attacker.example.test/callback", wantReject: true},
		{name: "different scheme", endpoint: "http://hooks.example.test/callback", wantReject: true},
		{name: "scheme-relative host", endpoint: "//attacker.example.test/callback", wantReject: true},
		{name: "embedded credentials", endpoint: "https://user:pass@hooks.example.test/callback", wantReject: true},
		{name: "fragment", endpoint: "/callback#fragment", wantReject: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := worker.resolveEndpoint(test.endpoint)
			if test.wantReject {
				if err == nil {
					t.Fatalf("resolveEndpoint(%q) = %q, want rejection", test.endpoint, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveEndpoint(%q): %v", test.endpoint, err)
			}
			if got != test.want {
				t.Fatalf("resolveEndpoint(%q) = %q, want %q", test.endpoint, got, test.want)
			}
		})
	}
}

func TestWebhookRetryWorkerRejectsInvalidConfiguredOrigin(t *testing.T) {
	worker := &WebhookRetryWorker{baseURL: "file:///tmp"}
	if _, err := worker.resolveEndpoint("/api/v1/webhook/paymob"); err == nil || !strings.Contains(err.Error(), "base URL") {
		t.Fatalf("invalid configured base URL error = %v, want a base URL validation error", err)
	}
}
