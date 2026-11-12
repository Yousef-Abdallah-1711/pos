package http

import (
	"log"
	"net/http"
	appconfig "pos-saas/internal/config"
	"strconv"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
	"time"
)

// SuperAdminHandler handles super admin HTTP requests
type SuperAdminHandler struct {
	authUC   *usecase.SuperAdminAuthUseCase
	tenantUC *usecase.SuperAdminTenantUseCase
	sessions *service.AuthSessionStore
}

// NewSuperAdminHandler creates a new super admin handler
func NewSuperAdminHandler(
	authUC *usecase.SuperAdminAuthUseCase,
	tenantUC *usecase.SuperAdminTenantUseCase,
	sessions *service.AuthSessionStore,
) *SuperAdminHandler {
	return &SuperAdminHandler{
		authUC:   authUC,
		tenantUC: tenantUC,
		sessions: sessions,
	}
}

// Login handles super admin login
func (h *SuperAdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsTrustedOrigin(r.Header.Get("Origin")) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{"success": false, "error": "invalid request origin"})
		return
	}

	var req domain.SuperAdminLoginRequest
	if err := decodeJSONStrict(r, &req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	response, err := h.authUC.Login(r.Context(), &req, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] Login failed: %v\n", err)
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if h.sessions == nil || response.Token == "" {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "session service unavailable"})
		return
	}
	if previous, cookieErr := r.Cookie(service.SuperAdminSessionCookieName); cookieErr == nil && previous.Value != "" {
		if revokeErr := h.sessions.Revoke(r.Context(), previous.Value); revokeErr != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "unable to rotate session"})
			return
		}
	}
	sessionID, err := h.sessions.Create(r.Context(), response.Token)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "unable to establish session"})
		return
	}
	secure := r.TLS != nil || appconfig.IsProductionMode()
	http.SetCookie(w, &http.Cookie{Name: service.SuperAdminSessionCookieName, Value: sessionID, Path: "/", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: int(h.sessions.TTL().Seconds()), Expires: time.Now().Add(h.sessions.TTL())})
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    map[string]interface{}{"super_admin": response.SuperAdmin, "user": response.User, "expires_at": response.ExpiresAt},
	})
}

// Session returns the authenticated super-admin identity without exposing its credential.
func (h *SuperAdminHandler) Session(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetSuperAdminClaims(r)
	if claims == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{"success": false, "error": "unauthorized"})
		return
	}
	superAdmin, err := h.authUC.ValidateSession(r.Context(), claims.SuperAdminID, middleware.GetClientIP(r))
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{"success": false, "error": "unauthorized"})
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"super_admin":    superAdmin,
			"email":          claims.Email,
			"is_super_admin": claims.IsSuperAdmin,
		},
	})
}

// Logout handles super admin logout
func (h *SuperAdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] Logout request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err := h.authUC.Logout(r.Context(), superAdminID, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] Logout failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if cookie, err := r.Cookie(service.SuperAdminSessionCookieName); err == nil && h.sessions != nil {
		if err := h.sessions.Revoke(r.Context(), cookie.Value); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "unable to end session"})
			return
		}
	}
	http.SetCookie(w, &http.Cookie{Name: service.SuperAdminSessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: r.TLS != nil || appconfig.IsProductionMode(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})
	w.Header().Set("Cache-Control", "no-store")
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "logged out successfully",
	})
}

// GetDashboardStats returns statistics for the super admin dashboard
func (h *SuperAdminHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] GetDashboardStats request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	stats, err := h.tenantUC.GetDashboardStats(r.Context(), superAdminID, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] GetDashboardStats failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// ListTenants returns a paginated list of all tenants
func (h *SuperAdminHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] ListTenants request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Parse query parameters
	pagination := parsePagination(r)
	filter := parseTenantFilter(r)

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	result, err := h.tenantUC.ListTenants(r.Context(), superAdminID, filter, pagination, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] ListTenants failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result.Data,
		"meta": map[string]interface{}{
			"total":       result.Total,
			"page":        result.Page,
			"page_size":   result.PageSize,
			"total_pages": result.TotalPages,
		},
	})
}

// GetTenant returns details for a specific tenant
func (h *SuperAdminHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] GetTenant request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get tenant ID from path
	tenantIDStr := r.PathValue("id")
	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid tenant ID",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	tenant, err := h.tenantUC.GetTenantDetails(r.Context(), superAdminID, tenantID, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] GetTenant failed: %v\n", err)
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tenant,
	})
}

// GetTenantUsers returns users for a specific tenant
func (h *SuperAdminHandler) GetTenantUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] GetTenantUsers request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get tenant ID from path
	tenantIDStr := r.PathValue("id")
	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid tenant ID",
		})
		return
	}

	pagination := parsePagination(r)
	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	result, err := h.tenantUC.GetTenantUsers(r.Context(), superAdminID, tenantID, pagination, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] GetTenantUsers failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result.Data,
		"meta": map[string]interface{}{
			"total":       result.Total,
			"page":        result.Page,
			"page_size":   result.PageSize,
			"total_pages": result.TotalPages,
		},
	})
}

// GenerateDashboardLoginToken fails closed because the legacy exchange creates a tenant-owner session
// without retaining the Super Admin as a distinct, restricted actor.
func (h *SuperAdminHandler) GenerateDashboardLoginToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
		"success": false,
		"error":   "tenant dashboard access is paused until restricted, actor-attributed support sessions are available",
	})
}

// DisableTenant disables a tenant
func (h *SuperAdminHandler) DisableTenant(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] DisableTenant request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get tenant ID from path
	tenantIDStr := r.PathValue("id")
	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid tenant ID",
		})
		return
	}

	// Parse reason from body
	var body struct {
		Reason string `json:"reason"`
	}
	_ = decodeJSONStrictAllowEmpty(r, &body)

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err = h.tenantUC.DisableTenant(r.Context(), superAdminID, tenantID, body.Reason, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] DisableTenant failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "tenant disabled successfully",
	})
}

// EnableTenant enables a tenant
func (h *SuperAdminHandler) EnableTenant(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] EnableTenant request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get tenant ID from path
	tenantIDStr := r.PathValue("id")
	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid tenant ID",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err = h.tenantUC.EnableTenant(r.Context(), superAdminID, tenantID, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] EnableTenant failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "tenant enabled successfully",
	})
}

// ListUsers returns a paginated list of all users across tenants
func (h *SuperAdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] ListUsers request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	pagination := parsePagination(r)
	filter := parseUserFilter(r)

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	result, err := h.tenantUC.ListAllUsers(r.Context(), superAdminID, filter, pagination, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] ListUsers failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result.Data,
		"meta": map[string]interface{}{
			"total":       result.Total,
			"page":        result.Page,
			"page_size":   result.PageSize,
			"total_pages": result.TotalPages,
		},
	})
}

// GetUser returns details for a specific user
func (h *SuperAdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] GetUser request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get user ID from path
	userIDStr := r.PathValue("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	user, err := h.tenantUC.GetUserDetails(r.Context(), superAdminID, userID, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] GetUser failed: %v\n", err)
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// GetAuditLogs returns audit logs with filtering
func (h *SuperAdminHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] GetAuditLogs request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	pagination := parsePagination(r)
	filter := parseAuditLogFilter(r)

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	result, err := h.tenantUC.GetAuditLogs(r.Context(), superAdminID, filter, pagination, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] GetAuditLogs failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result.Data,
		"meta": map[string]interface{}{
			"total":       result.Total,
			"page":        result.Page,
			"page_size":   result.PageSize,
			"total_pages": result.TotalPages,
		},
	})
}

// UpdateIPWhitelist updates the IP whitelist for the super admin
func (h *SuperAdminHandler) UpdateIPWhitelist(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SuperAdminHandler] UpdateIPWhitelist request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var body struct {
		IPWhitelist []string `json:"ip_whitelist"`
	}
	if err := decodeJSONStrict(r, &body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err := h.authUC.UpdateIPWhitelist(r.Context(), superAdminID, body.IPWhitelist, ip, userAgent)
	if err != nil {
		log.Printf("[SuperAdminHandler] UpdateIPWhitelist failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "IP whitelist updated successfully",
	})
}

// Helper functions

func parsePagination(r *http.Request) domain.Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func parseTenantFilter(r *http.Request) domain.TenantFilter {
	filter := domain.TenantFilter{}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	return filter
}

func parseUserFilter(r *http.Request) domain.UserFilter {
	filter := domain.UserFilter{}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
	}

	if tenantIDStr := r.URL.Query().Get("tenant_id"); tenantIDStr != "" {
		tenantID, _ := strconv.ParseInt(tenantIDStr, 10, 64)
		if tenantID > 0 {
			filter.TenantID = &tenantID
		}
	}

	if role := r.URL.Query().Get("role"); role != "" {
		filter.Role = &role
	}

	return filter
}

func parseAuditLogFilter(r *http.Request) domain.AuditLogFilter {
	filter := domain.AuditLogFilter{}

	if action := r.URL.Query().Get("action"); action != "" {
		filter.Action = &action
	}

	if targetType := r.URL.Query().Get("target_type"); targetType != "" {
		filter.TargetType = &targetType
	}

	if targetIDStr := r.URL.Query().Get("target_id"); targetIDStr != "" {
		targetID, _ := strconv.ParseInt(targetIDStr, 10, 64)
		if targetID > 0 {
			filter.TargetID = &targetID
		}
	}

	return filter
}

// Note: respondJSON is defined in auth_handler.go and reused here
