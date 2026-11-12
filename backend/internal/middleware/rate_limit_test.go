package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"pos-saas/internal/pkg/jwt"
)

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
