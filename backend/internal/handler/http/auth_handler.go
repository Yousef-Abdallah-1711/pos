package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"pos-saas/internal/config"
	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

type AuthHandler struct {
	useCase  *usecase.AuthUseCase
	sessions *service.AuthSessionStore
}

func NewAuthHandler(useCase *usecase.AuthUseCase, sessions *service.AuthSessionStore) *AuthHandler {
	return &AuthHandler{useCase: useCase, sessions: sessions}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.useCase.Register(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "already taken") {
			statusCode = http.StatusConflict
		}
		respondError(w, statusCode, err.Error())
		return
	}
	if err := h.establishSession(w, r, response.Token); err != nil {
		respondError(w, http.StatusInternalServerError, "Unable to establish session")
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"user": response.User})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.useCase.Login(&req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if authenticated, ok := response.(*domain.AuthResponse); ok {
		if err := h.establishSession(w, r, authenticated.Token); err != nil {
			respondError(w, http.StatusInternalServerError, "Unable to establish session")
			return
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{"user": authenticated.User})
		return
	}
	respondJSON(w, http.StatusOK, response)
}

// LoginConfirm handles tenant selection after multiple tenants are detected
func (h *AuthHandler) LoginConfirm(w http.ResponseWriter, r *http.Request) {
	var req domain.TenantSelectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.useCase.LoginConfirm(&req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := h.establishSession(w, r, response.Token); err != nil {
		respondError(w, http.StatusInternalServerError, "Unable to establish session")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"user": response.User})
}

func (h *AuthHandler) establishSession(w http.ResponseWriter, r *http.Request, credential string) error {
	if h.sessions == nil || credential == "" {
		return service.ErrSessionNotFound
	}
	var sessionID string
	var err error
	if previous, cookieErr := r.Cookie(service.AuthSessionCookieName); cookieErr == nil && previous.Value != "" {
		sessionID, err = h.sessions.Rotate(r.Context(), previous.Value, credential)
		if errors.Is(err, service.ErrSessionNotFound) {
			// A stale cookie must not block a fresh login. Create generates a new
			// unpredictable identifier; the stale identifier is never adopted.
			sessionID, err = h.sessions.Create(r.Context(), credential)
		}
	} else {
		sessionID, err = h.sessions.Create(r.Context(), credential)
	}
	if err != nil {
		return err
	}
	secure := r.TLS != nil || config.IsProductionMode()
	http.SetCookie(w, &http.Cookie{
		Name:     service.AuthSessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.sessions.TTL().Seconds()),
		Expires:  time.Now().Add(h.sessions.TTL()),
	})
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	return nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	for _, name := range []string{service.AuthSessionCookieName, service.ImpersonationSessionCookieName} {
		if cookie, err := r.Cookie(name); err == nil && cookie.Value != "" {
			if !middleware.IsTrustedOrigin(r.Header.Get("Origin")) {
				respondError(w, http.StatusForbidden, "Invalid request origin")
				return
			}
			if h.sessions == nil {
				respondError(w, http.StatusInternalServerError, "Session service unavailable")
				return
			}
			if err := h.sessions.Revoke(r.Context(), cookie.Value); err != nil {
				respondError(w, http.StatusInternalServerError, "Unable to end session")
				return
			}
		}
	}
	for _, name := range []string{service.AuthSessionCookieName, service.ImpersonationSessionCookieName} {
		http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", HttpOnly: true, Secure: r.TLS != nil || config.IsProductionMode(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})
	}
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id": claims.UserID, "user_id": claims.UserID, "email": claims.Email,
		"name": claims.Email, "tenant_id": claims.TenantID,
		"restaurant_id": claims.RestaurantID, "role": claims.Role,
		"is_super_admin": false,
	})
}

// CheckSubdomainAvailability checks if a subdomain is available
func (h *AuthHandler) CheckSubdomainAvailability(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		respondError(w, http.StatusBadRequest, "slug parameter is required")
		return
	}

	available, message := h.useCase.CheckSubdomainAvailability(slug)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"available": available,
		"message":   message,
		"slug":      slug,
	})
}

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.useCase.ForgotPassword(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Password reset email sent",
	})
}

// TokenExchange exchanges a super admin one-time token for a regular JWT
func (h *AuthHandler) TokenExchange(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusGone, "Browser token exchange is no longer supported")
}

// Helper functions
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
