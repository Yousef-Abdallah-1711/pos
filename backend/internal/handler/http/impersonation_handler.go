package http

import (
	"log"
	"net/http"
	"strconv"

	"pos-saas/internal/middleware"
	"pos-saas/internal/usecase"
)

// ImpersonationHandler handles impersonation HTTP requests
type ImpersonationHandler struct {
	impersonationUC *usecase.ImpersonationUseCase
}

// NewImpersonationHandler creates a new impersonation handler
func NewImpersonationHandler(impersonationUC *usecase.ImpersonationUseCase) *ImpersonationHandler {
	return &ImpersonationHandler{
		impersonationUC: impersonationUC,
	}
}

// StartImpersonation starts impersonating a user
func (h *ImpersonationHandler) StartImpersonation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] StartImpersonation request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Get user ID from path
	userIDStr := r.PathValue("userID")
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	response, err := h.impersonationUC.StartImpersonation(r.Context(), superAdminID, targetUserID, ip, userAgent)
	if err != nil {
		log.Printf("[ImpersonationHandler] StartImpersonation failed: %v\n", err)
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// EndImpersonation ends the current impersonation session
func (h *ImpersonationHandler) EndImpersonation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] EndImpersonation request\n")

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

	err := h.impersonationUC.EndImpersonation(r.Context(), superAdminID, ip, userAgent)
	if err != nil {
		log.Printf("[ImpersonationHandler] EndImpersonation failed: %v\n", err)
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "impersonation ended successfully",
	})
}

// EndImpersonationByToken ends impersonation using the token (called from tenant dashboard)
func (h *ImpersonationHandler) EndImpersonationByToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] EndImpersonationByToken request\n")

	// Get token from body or header
	var token string

	// Try to get from body
	var body struct {
		Token string `json:"token"`
	}
	if err := decodeJSONStrictAllowEmpty(r, &body); err == nil && body.Token != "" {
		token = body.Token
	}

	// Try header
	if token == "" {
		token = r.Header.Get("X-Impersonation-Token")
	}

	// Try query parameter
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "impersonation token required",
		})
		return
	}

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err := h.impersonationUC.EndImpersonationByToken(r.Context(), token, ip, userAgent)
	if err != nil {
		log.Printf("[ImpersonationHandler] EndImpersonationByToken failed: %v\n", err)
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "impersonation ended successfully",
	})
}

// GetImpersonationStatus returns the current impersonation status
func (h *ImpersonationHandler) GetImpersonationStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] GetImpersonationStatus request\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	session, err := h.impersonationUC.GetImpersonationStatus(r.Context(), superAdminID)
	if err != nil {
		log.Printf("[ImpersonationHandler] GetImpersonationStatus failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if session == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"is_impersonating": false,
			},
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"is_impersonating": true,
			"target_user_id":   session.TargetUserID,
			"target_tenant_id": session.TargetTenantID,
			"started_at":       session.StartedAt,
			"expires_at":       session.ExpiresAt,
		},
	})
}

// ValidateImpersonationToken validates an impersonation token (called from tenant dashboard)
func (h *ImpersonationHandler) ValidateImpersonationToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] ValidateImpersonationToken request\n")

	// Get token from query or header
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("X-Impersonation-Token")
	}

	if token == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "token required",
		})
		return
	}

	session, user, err := h.impersonationUC.ValidateImpersonation(r.Context(), token)
	if err != nil {
		log.Printf("[ImpersonationHandler] ValidateImpersonationToken failed: %v\n", err)
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"valid":   false,
			"error":   err.Error(),
		})
		return
	}

	// Clear sensitive data
	user.PasswordHash = ""

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"valid":   true,
		"data": map[string]interface{}{
			"user":       user,
			"tenant_id":  session.TargetTenantID,
			"expires_at": session.ExpiresAt,
			"session_id": session.ID,
		},
	})
}

// CheckImpersonation checks if the current request is impersonated (for tenant dashboard)
func (h *ImpersonationHandler) CheckImpersonation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] CheckImpersonation request\n")

	info := middleware.GetImpersonationInfo(r)
	if info == nil || !info.IsImpersonated {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"is_impersonated": false,
			},
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"is_impersonated":  true,
			"super_admin_id":   info.SuperAdminID,
			"target_user_id":   info.TargetUserID,
			"target_tenant_id": info.TargetTenantID,
		},
	})
}
