package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"time"
)

func (h *PublicOrderHandler) QuoteOrder(w http.ResponseWriter, r *http.Request) {
	setPublicOrderNoCacheHeaders(w)
	tenant, restaurant := middleware.GetTenantID(r), middleware.GetRestaurantID(r)
	if tenant <= 0 || restaurant <= 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant context")
		return
	}
	var request domain.OrderQuoteRequest
	if err := decodeJSONStrict(r, &request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid quote request")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	quote, err := h.orderUC.QuoteOrder(ctx, tenant, restaurant, &request)
	if err != nil {
		status := http.StatusBadRequest
		if !errors.Is(err, domain.ErrInvalidOrderQuote) {
			status = http.StatusServiceUnavailable
			log.Printf("[ORDER_QUOTE] unavailable tenant=%d restaurant=%d error_type=%T", tenant, restaurant, err)
		}
		respondError(w, status, "Unable to quote the requested items")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": quote})
}

// The legacy validation endpoint must return actual scoped prices, never placeholders.
func (h *PublicOrderHandler) writeValidatedOrderQuote(w http.ResponseWriter, r *http.Request, request *domain.CreateOrderRequest) {
	setPublicOrderNoCacheHeaders(w)
	tenant, restaurant := middleware.GetTenantID(r), middleware.GetRestaurantID(r)
	if tenant <= 0 || restaurant <= 0 {
		respondError(w, http.StatusUnauthorized, "Missing restaurant context")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	quote, err := h.orderUC.QuoteOrder(ctx, tenant, restaurant, &domain.OrderQuoteRequest{Items: request.Items})
	if err != nil {
		respondError(w, http.StatusBadRequest, "Unable to validate the requested items")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Order is valid", "data": map[string]interface{}{"estimated_subtotal": quote.Subtotal, "subtotal": quote.Subtotal, "tax_amount": quote.TaxAmount, "delivery_fee": quote.DeliveryFee, "total_amount": quote.TotalAmount, "currency": quote.Currency, "item_count": len(request.Items)}})
}
