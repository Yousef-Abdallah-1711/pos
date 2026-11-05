package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

func TestAuthRefreshRotatesCookieSessionExactlyOnceConcurrently(t *testing.T) {
	tokenService, err := jwt.NewTokenService("test-secret-that-is-long-enough-for-hmac", "15m")
	if err != nil {
		t.Fatal(err)
	}
	restaurantID := 17
	credential, err := tokenService.GenerateToken(23, 41, &restaurantID, "user@example.test", "manager")
	if err != nil {
		t.Fatal(err)
	}
	sessions := service.NewAuthSessionStore(service.NewInMemoryCache(), time.Minute)
	oldSessionID, err := sessions.Create(context.Background(), credential)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewAuthHandler(usecase.NewAuthUseCase(nil, tokenService), sessions)
	claims, err := tokenService.ValidateToken(credential)
	if err != nil {
		t.Fatal(err)
	}

	const attempts = 12
	var wait sync.WaitGroup
	var mu sync.Mutex
	successes, rejected := 0, 0
	var rotatedID string
	for range attempts {
		wait.Add(1)
		go func() {
			defer wait.Done()
			request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
			request.AddCookie(&http.Cookie{Name: service.AuthSessionCookieName, Value: oldSessionID})
			request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, claims))
			response := httptest.NewRecorder()
			handler.Refresh(response, request)
			mu.Lock()
			defer mu.Unlock()
			switch response.Code {
			case http.StatusNoContent:
				successes++
				cookies := response.Result().Cookies()
				if len(cookies) == 1 {
					rotatedID = cookies[0].Value
					if !cookies[0].HttpOnly || cookies[0].Value == oldSessionID {
						t.Errorf("refresh cookie must be HttpOnly and rotate the opaque ID")
					}
				} else {
					t.Errorf("refresh returned %d cookies, want exactly one", len(cookies))
				}
			case http.StatusUnauthorized:
				rejected++
			default:
				t.Errorf("refresh status = %d, want 204 or 401; body=%s", response.Code, response.Body.String())
			}
		}()
	}
	wait.Wait()
	if successes != 1 || rejected != attempts-1 {
		t.Fatalf("refresh results: successes=%d rejected=%d, want 1 and %d", successes, rejected, attempts-1)
	}
	if _, err := sessions.Resolve(context.Background(), oldSessionID); err != service.ErrSessionNotFound {
		t.Fatalf("old session remains valid after refresh: %v", err)
	}
	if resolved, err := sessions.Resolve(context.Background(), rotatedID); err != nil {
		t.Fatalf("rotated session cannot be resolved: %v", err)
	} else if parsed, parseErr := tokenService.ValidateToken(resolved); parseErr != nil || parsed.UserID != claims.UserID || parsed.TenantID != claims.TenantID {
		t.Fatalf("rotated session credential claims = %#v, %v", parsed, parseErr)
	}
}
