package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"pos-saas/internal/middleware"
	"strings"
	"testing"
)

func TestPublicOrderQuoteRequiresScopeAndRejectsClientPrices(t *testing.T) {
	handler := &PublicOrderHandler{}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/public/orders/quote", strings.NewReader(`{"items":[]}`))
	response := httptest.NewRecorder()
	handler.QuoteOrder(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatal("missing scope was accepted")
	}
	request = httptest.NewRequest(http.MethodPost, "/api/v1/public/orders/quote", strings.NewReader(`{"items":[{"product_id":1,"quantity":1,"price":0.01}]}`))
	ctx := context.WithValue(request.Context(), middleware.TenantContextKey, int64(1))
	ctx = context.WithValue(ctx, middleware.RestaurantContextKey, int64(1))
	response = httptest.NewRecorder()
	handler.QuoteOrder(response, request.WithContext(ctx))
	if response.Code != http.StatusBadRequest {
		t.Fatal("client-supplied price was accepted")
	}
	if response.Header().Get("Cache-Control") != "no-cache, no-store, private" {
		t.Fatal("quote response must not be shared-cacheable")
	}
}
