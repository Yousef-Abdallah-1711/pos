package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
)

func TestAccountingEventAdminRejectsMissingScope(t *testing.T) {
	handler := NewAccountingEventAdminHandler(nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounting/paid-order-events/failed", nil)
	response := httptest.NewRecorder()

	handler.ListFailed(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", response.Code, http.StatusBadRequest)
	}
	if response.Header().Get("Cache-Control") != "no-store, private" {
		t.Fatalf("cache control=%q, want no-store, private", response.Header().Get("Cache-Control"))
	}
}

func TestAccountingEventRedriveRequiresReasonBeforeRepository(t *testing.T) {
	handler := NewAccountingEventAdminHandler(nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/accounting/paid-order-events/1/redrive", strings.NewReader(`{"reason":"  "}`))
	request.SetPathValue("id", "1")
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(1))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(2))
	ctx = context.WithValue(ctx, middleware.UserContextKey, &jwt.Claims{UserID: 3})
	response := httptest.NewRecorder()

	handler.RedriveFailed(response, request.WithContext(ctx))

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", response.Code, http.StatusBadRequest)
	}
}
