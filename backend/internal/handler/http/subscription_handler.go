package http

import (
	"log"
	"net/http"
	"strconv"

	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

// SubscriptionHandler handles subscription-related HTTP requests
type SubscriptionHandler struct {
	useCase *usecase.SubscriptionUseCase
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(useCase *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{useCase: useCase}
}

// GetTenantStatus handles GET /api/v1/tenant/status
// Returns current tenant's subscription status
func (h *SubscriptionHandler) GetTenantStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SubscriptionHandler] GetTenantStatus - path: %s\n", r.URL.Path)

	// Get tenant from context (set by tenant middleware)
	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		log.Printf("[SubscriptionHandler] ERROR: No tenant in context\n")
		BadRequestResponse(w, "tenant context not found")
		return
	}

	info, err := h.useCase.GetTenantSubscription(r.Context(), tenantID)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to get tenant status - error: %v\n", err)
		InternalErrorResponse(w, "failed to get subscription status")
		return
	}

	log.Printf("[SubscriptionHandler] GetTenantStatus success - status: %s\n", info.SubscriptionStatus)
	OKResponse(w, info)
}

// GetTenantStatusBySlug handles GET /api/v1/public/tenant-status/{slug}
// Public endpoint for restaurant website to check tenant status
func (h *SubscriptionHandler) GetTenantStatusBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	log.Printf("[SubscriptionHandler] GetTenantStatusBySlug - slug: %s\n", slug)

	if slug == "" {
		log.Printf("[SubscriptionHandler] ERROR: No slug provided\n")
		BadRequestResponse(w, "slug is required")
		return
	}

	info, err := h.useCase.GetTenantSubscriptionBySlug(r.Context(), slug)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to get status by slug - error: %v\n", err)
		NotFoundResponse(w, "tenant not found")
		return
	}

	// Return simplified response for public use
	response := map[string]interface{}{
		"can_access":     info.CanAccess,
		"block_reason":   info.BlockReason,
		"is_maintenance": !info.CanAccess,
	}

	log.Printf("[SubscriptionHandler] GetTenantStatusBySlug success - can_access: %v\n", info.CanAccess)
	OKResponse(w, response)
}

// ExtendTenantTrial handles POST /api/v1/super-admin/tenants/{id}/extend-trial
// Extends the trial period for a tenant (super admin only)
func (h *SubscriptionHandler) ExtendTenantTrial(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.PathValue("id")
	log.Printf("[SubscriptionHandler] ExtendTenantTrial - tenant_id: %s\n", tenantIDStr)

	// Get user from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == 0 {
		log.Printf("[SubscriptionHandler] ERROR: No user in context\n")
		UnauthorizedResponse(w, "user not authenticated")
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Invalid tenant ID - error: %v\n", err)
		BadRequestResponse(w, "invalid tenant ID")
		return
	}

	// Parse request body
	var req struct {
		Days   int    `json:"days"`
		Reason string `json:"reason"`
	}

	if !DecodeJSONBody(w, r, &req) {
		return
	}

	if req.Days <= 0 {
		log.Printf("[SubscriptionHandler] ERROR: Invalid days - days: %d\n", req.Days)
		BadRequestResponse(w, "days must be greater than 0")
		return
	}

	// Call usecase
	reasonPtr := &req.Reason
	if req.Reason == "" {
		reasonPtr = nil
	}

	err = h.useCase.ExtendTrial(
		r.Context(),
		tenantID,
		req.Days,
		userID,
		reasonPtr,
		getIPAddress(r),
		r.UserAgent(),
	)

	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to extend trial - error: %v\n", err)
		InternalErrorResponse(w, "failed to extend trial")
		return
	}

	// Get updated subscription info
	info, _ := h.useCase.GetTenantSubscription(r.Context(), tenantID)

	response := map[string]interface{}{
		"new_trial_end_date": info.TrialEndDate,
		"days_added":         req.Days,
		"total_trial_days":   info.TrialDays,
	}

	log.Printf("[SubscriptionHandler] ExtendTenantTrial success\n")
	OKResponse(w, response)
}

// DisableTenant handles POST /api/v1/super-admin/tenants/{id}/disable
// Disables a tenant's access (super admin only)
func (h *SubscriptionHandler) DisableTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.PathValue("id")
	log.Printf("[SubscriptionHandler] DisableTenant - tenant_id: %s\n", tenantIDStr)

	// Get user from context
	userID := middleware.GetUserID(r)
	if userID == 0 {
		log.Printf("[SubscriptionHandler] ERROR: No user in context\n")
		UnauthorizedResponse(w, "user not authenticated")
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Invalid tenant ID - error: %v\n", err)
		BadRequestResponse(w, "invalid tenant ID")
		return
	}

	// Parse request body
	var req struct {
		Reason string `json:"reason"`
	}

	if !DecodeJSONBody(w, r, &req) {
		return
	}

	if req.Reason == "" {
		log.Printf("[SubscriptionHandler] ERROR: Reason is required\n")
		BadRequestResponse(w, "reason is required")
		return
	}

	// Call usecase
	err = h.useCase.DisableTenant(
		r.Context(),
		tenantID,
		userID,
		req.Reason,
		getIPAddress(r),
		r.UserAgent(),
	)

	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to disable tenant - error: %v\n", err)
		InternalErrorResponse(w, "failed to disable tenant")
		return
	}

	info, _ := h.useCase.GetTenantSubscription(r.Context(), tenantID)

	response := map[string]interface{}{
		"disabled_at":     info.DisabledAt,
		"previous_status": "disabled",
	}

	log.Printf("[SubscriptionHandler] DisableTenant success\n")
	OKResponse(w, response)
}

// EnableTenant handles POST /api/v1/super-admin/tenants/{id}/enable
// Re-enables a disabled tenant (super admin only)
func (h *SubscriptionHandler) EnableTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.PathValue("id")
	log.Printf("[SubscriptionHandler] EnableTenant - tenant_id: %s\n", tenantIDStr)

	// Get user from context
	userID := middleware.GetUserID(r)
	if userID == 0 {
		log.Printf("[SubscriptionHandler] ERROR: No user in context\n")
		UnauthorizedResponse(w, "user not authenticated")
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Invalid tenant ID - error: %v\n", err)
		BadRequestResponse(w, "invalid tenant ID")
		return
	}

	// Call usecase
	err = h.useCase.EnableTenant(
		r.Context(),
		tenantID,
		userID,
		getIPAddress(r),
		r.UserAgent(),
	)

	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to enable tenant - error: %v\n", err)
		InternalErrorResponse(w, "failed to enable tenant")
		return
	}

	info, _ := h.useCase.GetTenantSubscription(r.Context(), tenantID)

	response := map[string]interface{}{
		"is_disabled":     info.IsDisabled,
		"subscription_status": info.SubscriptionStatus,
	}

	log.Printf("[SubscriptionHandler] EnableTenant success\n")
	OKResponse(w, response)
}

// DeleteTenant handles DELETE /api/v1/super-admin/tenants/{id}
// Deletes a tenant permanently (super admin only)
func (h *SubscriptionHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.PathValue("id")
	log.Printf("[SubscriptionHandler] DeleteTenant - tenant_id: %s\n", tenantIDStr)

	// Get user from context
	userID := middleware.GetUserID(r)
	if userID == 0 {
		log.Printf("[SubscriptionHandler] ERROR: No user in context\n")
		UnauthorizedResponse(w, "user not authenticated")
		return
	}

	// Check for confirmation parameter
	confirm := r.URL.Query().Get("confirm")
	if confirm != "true" {
		log.Printf("[SubscriptionHandler] ERROR: Deletion not confirmed\n")
		BadRequestResponse(w, "deletion must be confirmed with ?confirm=true")
		return
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Invalid tenant ID - error: %v\n", err)
		BadRequestResponse(w, "invalid tenant ID")
		return
	}

	// Parse request body
	var req struct {
		Reason string `json:"reason"`
	}

	if !DecodeJSONBody(w, r, &req) {
		return
	}

	// Call usecase
	reasonPtr := &req.Reason
	if req.Reason == "" {
		reasonPtr = nil
	}

	err = h.useCase.DeleteTenant(
		r.Context(),
		tenantID,
		userID,
		reasonPtr,
		getIPAddress(r),
		r.UserAgent(),
	)

	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to delete tenant - error: %v\n", err)
		InternalErrorResponse(w, "failed to delete tenant")
		return
	}

	response := map[string]interface{}{
		"message":    "tenant scheduled for deletion",
		"deleted_at": "NOW()",
	}

	log.Printf("[SubscriptionHandler] DeleteTenant success\n")
	OKResponse(w, response)
}

// GetSubscriptionHistory handles GET /api/v1/super-admin/tenants/{id}/subscription-history
// Gets paginated subscription change history (super admin only)
func (h *SubscriptionHandler) GetSubscriptionHistory(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.PathValue("id")
	log.Printf("[SubscriptionHandler] GetSubscriptionHistory - tenant_id: %s\n", tenantIDStr)

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Invalid tenant ID - error: %v\n", err)
		BadRequestResponse(w, "invalid tenant ID")
		return
	}

	// Parse pagination parameters
	page := GetQueryParamInt(r, "page", 1)
	pageSize := GetQueryParamInt(r, "page_size", 20)

	// Call usecase
	history, total, err := h.useCase.GetSubscriptionHistory(r.Context(), tenantID, page, pageSize)
	if err != nil {
		log.Printf("[SubscriptionHandler] ERROR: Failed to get history - error: %v\n", err)
		InternalErrorResponse(w, "failed to get history")
		return
	}

	response := map[string]interface{}{
		"items":     history,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	log.Printf("[SubscriptionHandler] GetSubscriptionHistory success - count: %d\n", len(history))
	OKResponse(w, response)
}

// Helper function to get IP address from request
func getIPAddress(r *http.Request) string {
	// Try X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Try X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
