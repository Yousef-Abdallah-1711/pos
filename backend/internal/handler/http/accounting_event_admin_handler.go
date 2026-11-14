package http

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/middleware"
	"pos-saas/internal/repository"
)

const maxAccountingEventPage = 1_000_000

type AccountingEventAdminHandler struct {
	events *repository.AccountingEventRepository
}

func NewAccountingEventAdminHandler(events *repository.AccountingEventRepository) *AccountingEventAdminHandler {
	return &AccountingEventAdminHandler{events: events}
}

// ListFailed returns terminal accounting failures scoped to the authenticated restaurant.
func (h *AccountingEventAdminHandler) ListFailed(w http.ResponseWriter, r *http.Request) {
	setAccountingEventNoCache(w)
	tenantID, restaurantID := middleware.GetTenantID(r), middleware.GetRestaurantID(r)
	if tenantID <= 0 || restaurantID <= 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	page, limit, err := accountingEventPagination(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	items, total, err := h.events.ListFailed(r.Context(), tenantID, restaurantID, limit, (page-1)*limit)
	if err != nil {
		log.Printf("[AccountingEventAdmin] list_failed tenant_id=%d restaurant_id=%d error=%v", tenantID, restaurantID, err)
		respondError(w, http.StatusInternalServerError, "Failed to list failed accounting events")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": items,
		"pagination": map[string]int{
			"page": page, "limit": limit, "total": int(total),
			"total_pages": (int(total) + limit - 1) / limit,
		},
	})
}

// RedriveFailed requeues a terminal event only while its paid order remains paid.
func (h *AccountingEventAdminHandler) RedriveFailed(w http.ResponseWriter, r *http.Request) {
	setAccountingEventNoCache(w)
	tenantID, restaurantID := middleware.GetTenantID(r), middleware.GetRestaurantID(r)
	actorUserID := middleware.GetUserID(r)
	if tenantID <= 0 || restaurantID <= 0 || actorUserID <= 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant, restaurant, or user context")
		return
	}
	eventDBID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || eventDBID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid accounting event ID")
		return
	}
	var request struct {
		Reason string `json:"reason"`
	}
	if err := decodeJSONStrict(r, &request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	request.Reason = strings.TrimSpace(request.Reason)
	if request.Reason == "" || len(request.Reason) > 1000 {
		respondError(w, http.StatusBadRequest, "A redrive reason of 1 to 1000 bytes is required")
		return
	}
	eventID, err := h.events.RedriveFailed(r.Context(), tenantID, restaurantID, eventDBID, actorUserID, request.Reason)
	switch {
	case errors.Is(err, repository.ErrAccountingEventNotFound):
		respondError(w, http.StatusNotFound, "Accounting event not found")
		return
	case errors.Is(err, repository.ErrAccountingEventNotRedrivable):
		respondError(w, http.StatusConflict, "Accounting event is not failed or its order is no longer paid")
		return
	case err != nil:
		log.Printf("[AccountingEventAdmin] redrive_failed tenant_id=%d restaurant_id=%d event_id=%d actor_user_id=%d error=%v",
			tenantID, restaurantID, eventDBID, actorUserID, err)
		respondError(w, http.StatusInternalServerError, "Failed to requeue accounting event")
		return
	}
	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"success": true, "event_id": eventID, "status": repository.AccountingEventRetryable,
	})
}

func accountingEventPagination(r *http.Request) (page, limit int, err error) {
	page, limit = 1, 25
	if raw := r.URL.Query().Get("page"); raw != "" {
		page, err = strconv.Atoi(raw)
		if err != nil || page < 1 || page > maxAccountingEventPage {
			return 0, 0, errors.New("page must be between 1 and 1000000")
		}
	}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > 100 {
			return 0, 0, errors.New("limit must be between 1 and 100")
		}
	}
	return page, limit, nil
}

func setAccountingEventNoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, private")
}
