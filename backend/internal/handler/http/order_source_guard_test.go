package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
)

func TestOrderSourceComesFromProtectedRoute(t *testing.T) {
	for _, tc := range []struct {
		pattern, source, want string
		allowed               bool
	}{
		{"POST /api/v1/public/orders", "pos", "pos", false},
		{"POST /api/v1/public/orders", " PoS ", " PoS ", false},
		{"POST /api/v1/public/orders", "", "website", true},
		{"POST /api/v1/public/orders", "website", "website", true},
		{"POST /api/v1/public/orders", "dine_in", "dine_in", true},
		{"POST /api/v1/pos/orders", "", "pos", true},
		{"POST /api/v1/pos/orders", "website", "pos", true},
	} {
		t.Run(tc.pattern+"/"+tc.source, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/public/orders", nil)
			request.Pattern = tc.pattern
			body := domain.CreateOrderRequest{OrderSource: tc.source}
			if allowed := applyTrustedOrderSource(request, &body); allowed != tc.allowed || body.OrderSource != tc.want {
				t.Fatalf("allowed=%t source=%q; want %t/%q", allowed, body.OrderSource, tc.allowed, tc.want)
			}
		})
	}
}

func TestPublicOrderHandlerRejectsPOSSourceBeforeAnyBusinessWrite(t *testing.T) {
	handler := &PublicOrderHandler{} // Any usecase call would panic: rejection must precede it.
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/public/orders", handler.CreateOrder)
	for _, source := range []string{"pos", "POS", " pos "} {
		request := httptest.NewRequest(http.MethodPost, "/api/v1/public/orders", strings.NewReader(`{"order_source":"`+source+`"}`))
		ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(42))
		ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(3))
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request.WithContext(ctx))
		if response.Code != http.StatusForbidden {
			t.Fatalf("source=%q status=%d; want 403 before business writes", source, response.Code)
		}
	}
}
