package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsOriginAllowedSupportsSingleRestaurantSubdomain(t *testing.T) {
	corsOrigins = map[string]struct{}{
		"https://pos.rakmyat.online":               {},
		"https://*.restaurants.pos.rakmyat.online": {},
	}
	tests := []struct {
		origin string
		want   bool
	}{
		{origin: "https://pos.rakmyat.online", want: true},
		{origin: "https://cafe.restaurants.pos.rakmyat.online", want: true},
		{origin: "https://a-b.restaurants.pos.rakmyat.online", want: true},
		{origin: "https://restaurants.pos.rakmyat.online", want: false},
		{origin: "https://a.b.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online.evil.test", want: false},
		{origin: "https://cafe-restaurants.pos.rakmyat.online", want: false},
		{origin: "http://cafe.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online:444", want: false},
		{origin: "https://user@cafe.restaurants.pos.rakmyat.online", want: false},
		{origin: "https://cafe.restaurants.pos.rakmyat.online/path", want: false},
	}
	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			if got := isOriginAllowed(test.origin); got != test.want {
				t.Fatalf("isOriginAllowed(%q) = %t, want %t", test.origin, got, test.want)
			}
		})
	}
}

func TestCORSPreflightAllowsOrderCancellationIdempotencyHeader(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/admin/orders/1/items/1/cancel", nil)
	request.Header.Set("Access-Control-Request-Headers", "Content-Type, Idempotency-Key, X-Order-Token, X-Customer-Phone")
	response := httptest.NewRecorder()
	CORSMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("preflight must not mutate order") })).ServeHTTP(response, request)
	allowed := response.Header().Get("Access-Control-Allow-Headers")
	if response.Code != http.StatusOK || !strings.Contains(allowed, "Idempotency-Key") || !strings.Contains(allowed, "X-Order-Token") || !strings.Contains(allowed, "X-Customer-Phone") {
		t.Fatal("browser cancellation retries require the idempotency header in CORS preflight")
	}
}
