package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
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
	audit    authAuditLogger
}

const maxAuthRequestBytes = 16 << 10

type authErrorResponse struct {
	Error     string            `json:"error"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
}

type authAuditLogger interface {
	RecordAuthEvent(context.Context, int64, int64, string, string, string, string) error
}

func NewAuthHandler(useCase *usecase.AuthUseCase, sessions *service.AuthSessionStore) *AuthHandler {
	return &AuthHandler{useCase: useCase, sessions: sessions}
}

func (h *AuthHandler) SetAuditLogger(audit authAuditLogger) { h.audit = audit }

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := decodeAuthJSON(w, r, &req); err != nil {
		respondAuthDecodeError(w, r, err)
		return
	}

	response, err := h.useCase.RegisterContext(r.Context(), &req)
	if err != nil {
		if response != nil && response.User != nil {
			log.Printf("registration committed without token request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			respondJSON(w, http.StatusCreated, map[string]interface{}{
				"user": response.User, "session_established": false,
				"message": "Your account was created. Please sign in to continue.",
			})
			return
		}
		respondRegistrationError(w, r, err)
		return
	}
	if err := h.establishSession(w, r, response.Token); err != nil {
		log.Printf("registration session creation failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"user": response.User, "session_established": false,
			"message": "Your account was created. Please sign in to continue.",
		})
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"user": response.User, "session_established": true})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := decodeAuthJSON(w, r, &req); err != nil {
		respondAuthDecodeError(w, r, err)
		return
	}

	response, err := h.useCase.Login(&req)
	if err != nil {
		var validation *domain.ValidationErrors
		if errors.As(err, &validation) {
			writeAuthError(w, r, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Please review the highlighted fields.", validation.Fields)
			return
		}
		if errors.Is(err, domain.ErrAuthDependencyUnavailable) {
			log.Printf("login dependency unavailable request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			writeAuthError(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Sign in is temporarily unavailable. Please try again.", nil)
			return
		}
		if errors.Is(err, domain.ErrTokenGeneration) {
			log.Printf("login token creation failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			writeAuthError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "We could not sign you in. Please try again.", nil)
			return
		}
		writeAuthError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email or password is incorrect.", nil)
		return
	}
	if authenticated, ok := response.(*domain.AuthResponse); ok {
		if err := h.establishSession(w, r, authenticated.Token); err != nil {
			log.Printf("login session creation failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			writeAuthError(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Sign in is temporarily unavailable. Please try again.", nil)
			return
		}
		h.recordAuthEvent(r, int64(authenticated.User.TenantID), int64(authenticated.User.ID), "login")
		respondJSON(w, http.StatusOK, map[string]interface{}{"user": authenticated.User})
		return
	}
	respondJSON(w, http.StatusOK, response)
}

// LoginConfirm handles tenant selection after multiple tenants are detected
func (h *AuthHandler) LoginConfirm(w http.ResponseWriter, r *http.Request) {
	var req domain.TenantSelectionRequest
	if err := decodeAuthJSON(w, r, &req); err != nil {
		respondAuthDecodeError(w, r, err)
		return
	}

	response, err := h.useCase.LoginConfirm(&req)
	if err != nil {
		if errors.Is(err, domain.ErrAuthDependencyUnavailable) {
			log.Printf("login-confirm dependency unavailable request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			writeAuthError(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Sign in is temporarily unavailable. Please try again.", nil)
			return
		}
		if errors.Is(err, domain.ErrTokenGeneration) {
			log.Printf("login-confirm token creation failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
			writeAuthError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "We could not sign you in. Please try again.", nil)
			return
		}
		writeAuthError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Please sign in again to choose this organization.", nil)
		return
	}
	if err := h.establishSession(w, r, response.Token); err != nil {
		log.Printf("login-confirm session creation failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
		writeAuthError(w, r, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Sign in is temporarily unavailable. Please try again.", nil)
		return
	}
	h.recordAuthEvent(r, int64(response.User.TenantID), int64(response.User.ID), "login")
	respondJSON(w, http.StatusOK, map[string]interface{}{"user": response.User})
}

func decodeAuthJSON(w http.ResponseWriter, r *http.Request, target interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func respondAuthDecodeError(w http.ResponseWriter, r *http.Request, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeAuthError(w, r, http.StatusRequestEntityTooLarge, "REQUEST_TOO_LARGE", "Request body is too large.", nil)
		return
	}
	writeAuthError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", nil)
}

func respondRegistrationError(w http.ResponseWriter, r *http.Request, err error) {
	var validation *domain.ValidationErrors
	if errors.As(err, &validation) {
		writeAuthError(w, r, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Please review the highlighted fields.", validation.Fields)
		return
	}
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		writeAuthError(w, r, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "An account with this email already exists.", map[string]string{"email": "An account with this email already exists."})
	case errors.Is(err, domain.ErrSubdomainAlreadyExists):
		writeAuthError(w, r, http.StatusConflict, "SUBDOMAIN_TAKEN", "This website address is already in use.", map[string]string{"restaurant_slug": "This website address is already in use."})
	case errors.Is(err, domain.ErrRestaurantNameConflict):
		writeAuthError(w, r, http.StatusConflict, "RESTAURANT_NAME_CONFLICT", "A restaurant with this name already exists.", map[string]string{"restaurant_name": "A restaurant with this name already exists."})
	default:
		log.Printf("registration failed request_id=%s error=%v", middleware.RequestIDFromContext(r.Context()), err)
		writeAuthError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "We could not create your account. Please try again.", nil)
	}
}

func writeAuthError(w http.ResponseWriter, r *http.Request, status int, code, message string, fields map[string]string) {
	requestID := middleware.RequestIDFromContext(r.Context())
	w.Header().Set("Cache-Control", "no-store")
	respondJSON(w, status, authErrorResponse{Error: message, Code: code, Message: message, Fields: fields, RequestID: requestID})
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
	var tenantID, userID int64
	if cookie, err := r.Cookie(service.AuthSessionCookieName); err == nil && cookie.Value != "" && h.sessions != nil {
		if credential, resolveErr := h.sessions.Resolve(r.Context(), cookie.Value); resolveErr == nil {
			if claims, validateErr := h.useCase.ValidateCredential(credential); validateErr == nil {
				tenantID, userID = int64(claims.TenantID), int64(claims.UserID)
			}
		}
	}
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
	if tenantID > 0 {
		h.recordAuthEvent(r, tenantID, userID, "logout")
	}
	for _, name := range []string{service.AuthSessionCookieName, service.ImpersonationSessionCookieName} {
		http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", HttpOnly: true, Secure: r.TLS != nil || config.IsProductionMode(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})
	}
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) recordAuthEvent(r *http.Request, tenantID, userID int64, action string) {
	if h.audit == nil {
		return
	}
	if err := h.audit.RecordAuthEvent(r.Context(), tenantID, userID, action, "success", middleware.GetClientIP(r), r.UserAgent()); err != nil {
		log.Printf("auth audit write failed for %s: %v", action, err)
	}
}

// Refresh rotates the opaque browser session ID and replaces its JWT. The
// route must be wrapped in SessionAuthMiddleware so only an active session can
// reach this handler; the route registration also requires a trusted Origin.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(service.AuthSessionCookieName)
	if err != nil || cookie.Value == "" || h.sessions == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	claims := middleware.GetUserClaims(r)
	credential, err := h.useCase.RefreshToken(claims)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	sessionID, err := h.sessions.Rotate(r.Context(), cookie.Value, credential)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			respondError(w, http.StatusUnauthorized, "Session expired")
			return
		}
		respondError(w, http.StatusInternalServerError, "Unable to refresh session")
		return
	}
	h.setSessionCookie(w, r, sessionID)
	h.recordAuthEvent(r, int64(claims.TenantID), int64(claims.UserID), "refresh")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, sessionID string) {
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
