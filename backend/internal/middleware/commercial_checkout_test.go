package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommercialCheckoutDoesNotReachProvider(t *testing.T) {
	called := false
	handler := CommercialCheckoutGate(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/subscription/initiate", nil))
	if called || response.Code != http.StatusServiceUnavailable {
		t.Fatalf("checkout reached payment handler=%v status=%d", called, response.Code)
	}
}
