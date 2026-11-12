package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

const kdsMaxConnectionsPerRoom = 10

type kdsConnectionLimiter struct {
	mu         sync.Mutex
	maxPerRoom int
	counts     map[string]int
}

func newKDSConnectionLimiter(maxPerRoom int) *kdsConnectionLimiter {
	if maxPerRoom <= 0 {
		maxPerRoom = kdsMaxConnectionsPerRoom
	}

	return &kdsConnectionLimiter{
		maxPerRoom: maxPerRoom,
		counts:     make(map[string]int),
	}
}

func buildKDSConnectionKey(tenantID, restaurantID int64) string {
	return "tenant:" + strconv.FormatInt(tenantID, 10) + ":restaurant:" + strconv.FormatInt(restaurantID, 10)
}

func (l *kdsConnectionLimiter) Acquire(tenantID, restaurantID int64) bool {
	key := buildKDSConnectionKey(tenantID, restaurantID)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.counts[key] >= l.maxPerRoom {
		return false
	}

	l.counts[key]++
	return true
}

func (l *kdsConnectionLimiter) Release(tenantID, restaurantID int64) {
	key := buildKDSConnectionKey(tenantID, restaurantID)

	l.mu.Lock()
	defer l.mu.Unlock()

	current := l.counts[key]
	if current <= 1 {
		delete(l.counts, key)
		return
	}

	l.counts[key] = current - 1
}

// KDSWSHandler upgrades authenticated requests to websocket clients for KDS rooms.
type KDSWSHandler struct {
	hub          *KDSHub
	tokenService *jwt.TokenService
	sessions     interface {
		Resolve(context.Context, string) (string, error)
	}
	limiter  *kdsConnectionLimiter
	upgrader websocket.Upgrader
}

func NewKDSWSHandler(hub *KDSHub, tokenService *jwt.TokenService, sessions ...interface {
	Resolve(context.Context, string) (string, error)
}) *KDSWSHandler {
	var sessionResolver interface {
		Resolve(context.Context, string) (string, error)
	}
	if len(sessions) > 0 {
		sessionResolver = sessions[0]
	}
	return &KDSWSHandler{
		hub:          hub,
		tokenService: tokenService,
		sessions:     sessionResolver,
		limiter:      newKDSConnectionLimiter(kdsMaxConnectionsPerRoom),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return middleware.IsTrustedOrigin(r.Header.Get("Origin"))
			},
		},
	}
}

// ServeWS accepts only an API-origin browser session cookie. Credentials are never read from URLs.
func (h *KDSWSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	if h.hub == nil || h.tokenService == nil {
		respondJSONError(w, http.StatusServiceUnavailable, "KDS websocket is unavailable")
		return
	}

	cookie, err := r.Cookie(service.ImpersonationSessionCookieName)
	isImpersonation := err == nil && cookie.Value != ""
	if !isImpersonation {
		cookie, err = r.Cookie(service.AuthSessionCookieName)
	}
	if err != nil || cookie.Value == "" || h.sessions == nil {
		respondJSONError(w, http.StatusUnauthorized, "missing websocket session")
		return
	}
	if !middleware.IsTrustedOrigin(r.Header.Get("Origin")) {
		respondJSONError(w, http.StatusForbidden, "invalid websocket origin")
		return
	}
	token, err := h.sessions.Resolve(r.Context(), cookie.Value)
	if err != nil || token == "" {
		respondJSONError(w, http.StatusUnauthorized, "invalid websocket session")
		return
	}
	if isImpersonation && service.IsImpersonationCredential(token) {
		token = service.UnwrapImpersonationCredential(token)
	}

	claims, err := h.tokenService.ValidateToken(token)
	if err != nil {
		respondJSONError(w, http.StatusUnauthorized, "invalid websocket token")
		return
	}

	tenantID := int64(claims.TenantID)
	restaurantID := int64(claims.RestaurantID)
	userID := int64(claims.UserID)
	if tenantID <= 0 || restaurantID <= 0 {
		respondJSONError(w, http.StatusBadRequest, "invalid tenant context")
		return
	}
	if h.limiter != nil && !h.limiter.Acquire(tenantID, restaurantID) {
		respondJSONError(w, http.StatusTooManyRequests, "too many websocket connections for this restaurant")
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if h.limiter != nil {
			h.limiter.Release(tenantID, restaurantID)
		}
		return
	}

	client := newKDSClient(h.hub, conn, tenantID, restaurantID, userID, func() {
		if h.limiter != nil {
			h.limiter.Release(tenantID, restaurantID)
		}
	})
	h.hub.Register(client)

	// Best effort handshake payload to help clients verify room assignment.
	_ = conn.SetWriteDeadline(time.Now().Add(kdsWriteWait))
	_ = conn.WriteJSON(map[string]interface{}{
		"type":          "connected",
		"tenant_id":     tenantID,
		"restaurant_id": restaurantID,
		"user_id":       userID,
		"timestamp":     time.Now().UTC(),
	})
	_ = conn.SetWriteDeadline(time.Time{})

	go client.WritePump()
	go client.ReadPump()
}

func respondJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
