package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
)

func (h *AdminOrderHandler) CancelOrderItem(w http.ResponseWriter, r *http.Request) {
	request := domain.OrderItemCancellation{TenantID: middleware.GetTenantID(r), RestaurantID: middleware.GetRestaurantID(r), ActorID: middleware.GetUserID(r), RequestKey: r.Header.Get("Idempotency-Key")}
	var orderError, itemError error
	request.OrderID, orderError = strconv.ParseInt(r.PathValue("id"), 10, 64)
	request.ItemID, itemError = strconv.ParseInt(r.PathValue("item_id"), 10, 64)
	var input struct {
		Quantity int    `json:"quantity"`
		Reason   string `json:"reason"`
	}
	if request.ActorID <= 0 || request.TenantID <= 0 || request.RestaurantID <= 0 {
		respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	if decodeJSONStrict(r, &input) != nil || orderError != nil || itemError != nil || request.OrderID <= 0 || request.ItemID <= 0 || input.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "Valid order, item and positive quantity are required")
		return
	}
	request.Quantity = input.Quantity
	request.Reason = strings.TrimSpace(input.Reason)
	if _, err := uuid.Parse(request.RequestKey); err != nil || utf8.RuneCountInString(request.Reason) < 1 || utf8.RuneCountInString(request.Reason) > 500 {
		respondError(w, http.StatusBadRequest, "A UUID Idempotency-Key and a reason of 1–500 characters are required")
		return
	}
	result, err := h.orderUC.CancelOrderItem(r.Context(), request)
	if err != nil {
		respondCancellationError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": result})
}

func respondCancellationError(w http.ResponseWriter, err error) {
	var databaseError *pq.Error
	if errors.As(err, &databaseError) {
		switch string(databaseError.Code) {
		case "P0002":
			respondError(w, http.StatusNotFound, "Order or item not found")
			return
		case "42501":
			respondError(w, http.StatusForbidden, "Cancellation is outside authorized scope")
			return
		case "23514":
			respondError(w, http.StatusConflict, "Cancellation conflicts with order state, quantity or idempotency key")
			return
		}
	}
	respondError(w, http.StatusInternalServerError, "Failed to cancel order item")
}
