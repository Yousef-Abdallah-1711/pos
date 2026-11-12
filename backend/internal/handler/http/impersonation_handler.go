package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/config"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

// ImpersonationHandler handles impersonation HTTP requests
type ImpersonationHandler struct {
	impersonationUC *usecase.ImpersonationUseCase
	sessions        *service.AuthSessionStore
}

// NewImpersonationHandler creates a new impersonation handler
func NewImpersonationHandler(impersonationUC *usecase.ImpersonationUseCase, sessions ...*service.AuthSessionStore) *ImpersonationHandler {
	var sessionStore *service.AuthSessionStore
	if len(sessions) > 0 {
		sessionStore = sessions[0]
	}
	return &ImpersonationHandler{
		impersonationUC: impersonationUC,
		sessions:        sessionStore,
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
	if h.sessions == nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "session service unavailable"})
		return
	}
	ttl := time.Until(response.ExpiresAt)
	sessionID, err := h.sessions.CreateImpersonation(r.Context(), response.Token, ttl)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "unable to establish impersonation session"})
		return
	}
	secure := r.TLS != nil || config.IsProductionMode()
	http.SetCookie(w, &http.Cookie{Name: service.ImpersonationSessionCookieName, Value: sessionID, Path: "/", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: int(ttl.Seconds()), Expires: response.ExpiresAt})
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    map[string]interface{}{"expires_at": response.ExpiresAt, "target_user": response.TargetUser, "target_tenant_id": response.TargetTenantID},
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
	if cookie, cookieErr := r.Cookie(service.ImpersonationSessionCookieName); cookieErr == nil && h.sessions != nil {
		_ = h.sessions.Revoke(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: service.ImpersonationSessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: r.TLS != nil || config.IsProductionMode(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "impersonation ended successfully",
	})
}

// EndImpersonationByToken ends impersonation using the token (called from tenant dashboard)
func (h *ImpersonationHandler) EndImpersonationByToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ImpersonationHandler] EndImpersonationByToken request\n")

	cookie, cookieErr := r.Cookie(service.ImpersonationSessionCookieName)
	if cookieErr != nil || cookie.Value == "" || h.sessions == nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "impersonation session required",
		})
		return
	}
	credential, err := h.sessions.Resolve(r.Context(), cookie.Value)
	if err != nil || !service.IsImpersonationCredential(credential) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{"success": false, "error": "invalid impersonation session"})
		return
	}
	token := service.UnwrapImpersonationCredential(credential)

	ip := middleware.GetClientIP(r)
	userAgent := r.UserAgent()

	err = h.impersonationUC.EndImpersonationByToken(r.Context(), token, ip, userAgent)
	if err != nil {
		log.Printf("[ImpersonationHandler] EndImpersonationByToken failed: %v\n", err)
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if err := h.sessions.Revoke(r.Context(), cookie.Value); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "unable to end impersonation session"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: service.ImpersonationSessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: r.TLS != nil || config.IsProductionMode(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})

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

	cookie, cookieErr := r.Cookie(service.ImpersonationSessionCookieName)
	if cookieErr != nil || cookie.Value == "" || h.sessions == nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "impersonation session required",
		})
		return
	}
	credential, err := h.sessions.Resolve(r.Context(), cookie.Value)
	if err != nil || !service.IsImpersonationCredential(credential) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{"success": false, "valid": false, "error": "invalid impersonation session"})
		return
	}
	token := service.UnwrapImpersonationCredential(credential)

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
