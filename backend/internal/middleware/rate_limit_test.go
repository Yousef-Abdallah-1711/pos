package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/service"
)

func TestSharedIPRateLimitIsSharedAcrossMiddlewareInstances(t *testing.T) {
	store := service.NewInMemoryRateLimitStore(32)
	first := RateLimitByIPWithStore(store, "auth-login", 1, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	second := RateLimitByIPWithStore(store, "auth-login", 1, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))

	doRequest := func(handler http.Handler, remoteAddr string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = remoteAddr
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}
	if response := doRequest(first, "203.0.113.10:1234"); response.Code != http.StatusOK {
		t.Fatalf("first request = %d, want 200", response.Code)
	}
	blocked := doRequest(second, "203.0.113.10:5678")
	if blocked.Code != http.StatusTooManyRequests {
		t.Fatalf("second instance request = %d, want 429", blocked.Code)
	}
	if blocked.Header().Get("Retry-After") == "" || blocked.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Fatalf("rate-limit response headers missing: %v", blocked.Header())
	}
	if response := doRequest(second, "203.0.113.11:1234"); response.Code != http.StatusOK {
		t.Fatalf("different IP = %d, want 200", response.Code)
	}
}

type failingSharedRateLimitStore struct{}

func (failingSharedRateLimitStore) IncrementWindow(context.Context, string, time.Duration) (int64, time.Duration, error) {
	return 0, 0, errors.New("store unavailable")
}

func TestSharedIPRateLimitFailsClosedWhenStoreIsUnavailable(t *testing.T) {
	handlerCalled := false
	wrapped := RateLimitByIPWithStore(failingSharedRateLimitStore{}, "auth-login", 20, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	request.RemoteAddr = "203.0.113.10:1234"
	response := httptest.NewRecorder()
	wrapped.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable || handlerCalled {
		t.Fatalf("store failure status=%d handlerCalled=%t, want 503 and false", response.Code, handlerCalled)
	}
}

func withTestUser(req *http.Request, userID int) *http.Request {
	claims := &jwt.Claims{UserID: userID}
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	return req.WithContext(ctx)
}

func TestRateLimitByUserWindowEnforcesPerUser(t *testing.T) {
	t.Parallel()

	wrapped := RateLimitByUserWindow(2, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	doRequest := func(userID int) int {
		req := httptest.NewRequest(http.MethodGet, "/kds/test", nil)
		req = withTestUser(req, userID)
		req.RemoteAddr = "127.0.0.1:3000"
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		return rr.Code
	}

	if code := doRequest(1001); code != http.StatusOK {
		t.Fatalf("first request should pass, got %d", code)
	}
	if code := doRequest(1001); code != http.StatusOK {
		t.Fatalf("second request should pass, got %d", code)
	}
	if code := doRequest(1001); code != http.StatusTooManyRequests {
		t.Fatalf("third request should be rate-limited, got %d", code)
	}

	// Different user should have an independent counter.
	if code := doRequest(2002); code != http.StatusOK {
		t.Fatalf("different user should not be rate-limited, got %d", code)
	}
}

func TestRateLimitByUserWindowFallsBackToIP(t *testing.T) {
	t.Parallel()

	wrapped := RateLimitByUserWindow(1, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	doRequest := func(remoteAddr string) int {
		req := httptest.NewRequest(http.MethodGet, "/kds/test", nil)
		req.RemoteAddr = remoteAddr
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		return rr.Code
	}

	if code := doRequest("10.1.1.1:4444"); code != http.StatusOK {
		t.Fatalf("first request should pass, got %d", code)
	}
	if code := doRequest("10.1.1.1:5555"); code != http.StatusTooManyRequests {
		t.Fatalf("same IP should be rate-limited without user claims, got %d", code)
	}
	if code := doRequest("10.1.1.2:5555"); code != http.StatusOK {
		t.Fatalf("different IP should not be rate-limited, got %d", code)
	}
}

func TestRateLimitSuperAdminConcurrentRequests(t *testing.T) {
	t.Parallel()

	const requestCount = 64
	wrapped := RateLimitSuperAdmin(requestCount)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	var wait sync.WaitGroup
	statuses := make(chan int, requestCount)
	for range requestCount {
		wait.Add(1)
		go func() {
			defer wait.Done()
			req := httptest.NewRequest(http.MethodGet, "/super-admin/tenants", nil)
			req.RemoteAddr = "127.0.0.1:3000"
			response := httptest.NewRecorder()
			wrapped.ServeHTTP(response, req)
			statuses <- response.Code
		}()
	}
	wait.Wait()
	close(statuses)
	for status := range statuses {
		if status != http.StatusOK {
			t.Fatalf("request below configured limit returned %d", status)
		}
	}
}
