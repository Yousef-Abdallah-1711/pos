package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

const (
	kdsQueueTimeout     = 2500 * time.Millisecond
	kdsStatsTimeout     = 2500 * time.Millisecond
	kdsSettingsTimeout  = 2500 * time.Millisecond
	kdsAnalyticsTimeout = 3500 * time.Millisecond
	kdsMutationTimeout  = 4000 * time.Millisecond
)

// KDSHandler handles KDS API endpoints.
type KDSHandler struct {
	kdsUC *usecase.KDSUseCase
}

// NewKDSHandler creates a new KDS HTTP handler.
func NewKDSHandler(kdsUC *usecase.KDSUseCase) *KDSHandler {
	return &KDSHandler{kdsUC: kdsUC}
}

// GetQueue returns active KDS queue entries.
// GET /api/v1/kds/queue
func (h *KDSHandler) GetQueue(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	filters := domain.KDSFilters{
		IncludeCompleted: parseBoolQuery(r.URL.Query().Get("include_completed")),
	}

	for _, status := range collectStatuses(r) {
		filters.Statuses = append(filters.Statuses, domain.KDSStatus(status))
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsQueueTimeout)
	defer cancel()

	orders, err := h.kdsUC.GetQueue(ctx, tenantID, restaurantID, filters)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS queue request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data":    []domain.KDSOrder{},
			})
			return
		}
		if strings.Contains(err.Error(), "invalid kds status") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch KDS queue")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    orders,
	})
}

// BumpOrder advances a KDS order to next status.
// PUT /api/v1/kds/orders/{id}/bump
func (h *KDSHandler) BumpOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	userID := middleware.GetUserID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	queueID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || queueID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid KDS order ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsMutationTimeout)
	defer cancel()

	order, err := h.kdsUC.BumpOrder(ctx, tenantID, restaurantID, queueID, userID)
	if err != nil {
		log.Printf("[KDS] bump failed: tenant=%d restaurant=%d queue=%d user=%d err=%v", tenantID, restaurantID, queueID, userID, err)
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS bump request timed out")
			return
		}
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "not found"):
			respondError(w, http.StatusNotFound, "KDS order not found")
		case strings.Contains(errMsg, "invalid"), strings.Contains(errMsg, "already completed"), strings.Contains(errMsg, "retry"):
			respondError(w, http.StatusBadRequest, errMsg)
		default:
			respondError(w, http.StatusInternalServerError, "Failed to bump KDS order")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

// RecallOrder moves a completed KDS order back to ready.
// PUT /api/v1/kds/orders/{id}/recall
func (h *KDSHandler) RecallOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	queueID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || queueID <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid KDS order ID")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsMutationTimeout)
	defer cancel()

	order, err := h.kdsUC.RecallOrder(ctx, tenantID, restaurantID, queueID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS recall request timed out")
			return
		}
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "not found"):
			respondError(w, http.StatusNotFound, "KDS order not found")
		case strings.Contains(errMsg, "invalid"), strings.Contains(errMsg, "recallable"):
			respondError(w, http.StatusBadRequest, errMsg)
		default:
			respondError(w, http.StatusInternalServerError, "Failed to recall KDS order")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

// GetStats returns queue statistics.
// GET /api/v1/kds/stats
func (h *KDSHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsStatsTimeout)
	defer cancel()

	stats, err := h.kdsUC.GetStats(ctx, tenantID, restaurantID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS stats request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"active_orders":    int64(0),
					"avg_prep_minutes": 0.0,
					"completed_today":  int64(0),
					"new_orders":       int64(0),
					"in_progress":      int64(0),
					"ready":            int64(0),
				},
			})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch KDS stats")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// GetAnalytics returns KDS analytics.
// GET /api/v1/kds/analytics
func (h *KDSHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsAnalyticsTimeout)
	defer cancel()

	analytics, err := h.kdsUC.GetAnalytics(ctx, tenantID, restaurantID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS analytics request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"average_prep_minutes":   0.0,
					"total_completed_orders": int64(0),
					"prep_time_trend":        []domain.KDSPrepTrendPoint{},
					"busiest_hours":          []domain.KDSHourLoadPoint{},
				},
			})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch KDS analytics")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analytics,
	})
}

// GetSettings returns KDS settings.
// GET /api/v1/kds/settings
func (h *KDSHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsSettingsTimeout)
	defer cancel()

	settings, err := h.kdsUC.GetSettings(ctx, tenantID, restaurantID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS settings request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data":    defaultKDSSettingsResponse(tenantID, restaurantID),
			})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch KDS settings")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    settings,
	})
}

// UpdateSettings updates KDS settings.
// PUT /api/v1/kds/settings
func (h *KDSHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r)
	restaurantID := middleware.GetRestaurantID(r)
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}
	if h.kdsUC == nil {
		respondError(w, http.StatusServiceUnavailable, "KDS service unavailable")
		return
	}

	var req struct {
		AutoBumpMinutes     *int    `json:"auto_bump_minutes"`
		AlertSoundEnabled   *bool   `json:"alert_sound_enabled"`
		ShowAnalyticsPanel  *bool   `json:"show_analytics_panel"`
		DisplayDensity      *string `json:"display_density"`
		ColorWarningMinutes *int    `json:"color_warning_minutes"`
		ColorDangerMinutes  *int    `json:"color_danger_minutes"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AutoBumpMinutes == nil &&
		req.AlertSoundEnabled == nil &&
		req.ShowAnalyticsPanel == nil &&
		req.DisplayDensity == nil &&
		req.ColorWarningMinutes == nil &&
		req.ColorDangerMinutes == nil {
		respondError(w, http.StatusBadRequest, "At least one settings field is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), kdsMutationTimeout)
	defer cancel()

	current, err := h.kdsUC.GetSettings(ctx, tenantID, restaurantID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS settings request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			defaults := defaultKDSSettingsResponse(tenantID, restaurantID)
			current = &defaults
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to load current KDS settings")
			return
		}
	}

	settings := *current
	if req.AutoBumpMinutes != nil {
		settings.AutoBumpMinutes = *req.AutoBumpMinutes
	}
	if req.AlertSoundEnabled != nil {
		settings.AlertSoundEnabled = *req.AlertSoundEnabled
	}
	if req.ShowAnalyticsPanel != nil {
		settings.ShowAnalyticsPanel = *req.ShowAnalyticsPanel
	}
	if req.DisplayDensity != nil {
		settings.DisplayDensity = strings.TrimSpace(*req.DisplayDensity)
	}
	if req.ColorWarningMinutes != nil {
		settings.ColorWarningMinutes = *req.ColorWarningMinutes
	}
	if req.ColorDangerMinutes != nil {
		settings.ColorDangerMinutes = *req.ColorDangerMinutes
	}
	settings.TenantID = tenantID
	settings.RestaurantID = restaurantID

	if err := h.kdsUC.UpdateSettings(ctx, &settings); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS update request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data":    settings,
			})
			return
		}
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid") ||
			strings.Contains(errMsg, "must be") ||
			strings.Contains(errMsg, "required") ||
			strings.Contains(errMsg, "greater than") {
			respondError(w, http.StatusBadRequest, errMsg)
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update KDS settings")
		return
	}

	updated, err := h.kdsUC.GetSettings(ctx, tenantID, restaurantID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "KDS settings request timed out")
			return
		}
		if isKDSStorageUnavailable(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"success": true,
				"data":    settings,
			})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to load updated KDS settings")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func collectStatuses(r *http.Request) []string {
	values := make([]string, 0)
	addList := func(raw string) {
		for _, part := range strings.Split(raw, ",") {
			status := strings.TrimSpace(part)
			if status == "" {
				continue
			}
			values = append(values, status)
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		addList(status)
	}

	if statuses := r.URL.Query().Get("statuses"); statuses != "" {
		addList(statuses)
	}

	if repeated, ok := r.URL.Query()["status"]; ok && len(repeated) > 1 {
		for _, value := range repeated {
			addList(value)
		}
	}

	return values
}

func parseBoolQuery(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func isKDSStorageUnavailable(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "database is not initialized") ||
		strings.Contains(msg, "kds repository is not configured") ||
		strings.Contains(msg, "relation \"kds_")
}

func defaultKDSSettingsResponse(tenantID, restaurantID int64) domain.KDSSettings {
	return domain.KDSSettings{
		TenantID:            tenantID,
		RestaurantID:        restaurantID,
		AutoBumpMinutes:     0,
		AlertSoundEnabled:   true,
		ShowAnalyticsPanel:  false,
		DisplayDensity:      "comfortable",
		ColorWarningMinutes: 10,
		ColorDangerMinutes:  20,
	}
}
