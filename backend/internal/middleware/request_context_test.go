package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddlewareGeneratesServerOwnedID(t *testing.T) {
	const forgedID = "attacker-controlled-id"
	var contextID string
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextID = RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set(RequestIDHeader, forgedID)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	responseID := response.Header().Get(RequestIDHeader)
	if responseID == "" || responseID == forgedID {
		t.Fatalf("expected a generated request ID, got %q", responseID)
	}
	if contextID != responseID {
		t.Fatalf("context request ID %q does not match response request ID %q", contextID, responseID)
	}
}

func TestRecoveryMiddlewareHidesPanicDetails(t *testing.T) {
	const secretPanic = "sensitive-panic-payload"
	handler := RequestIDMiddleware(RecoveryMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(secretPanic)
	})))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/private?token=secret", nil))

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
	if strings.Contains(response.Body.String(), secretPanic) || strings.Contains(response.Body.String(), "token=secret") {
		t.Fatalf("panic or request details leaked in response: %q", response.Body.String())
	}
	if response.Header().Get(RequestIDHeader) == "" {
		t.Fatal("recovered request is missing its correlation ID")
	}
}
