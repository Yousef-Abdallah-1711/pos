package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestReadinessRouteReflectsDatabaseHealth(t *testing.T) {
	tests := []struct {
		name            string
		databaseHealthy bool
		wantStatus      int
		wantBody        string
	}{
		{name: "ready", databaseHealthy: true, wantStatus: http.StatusOK, wantBody: `"status":"ready"`},
		{name: "database unavailable", databaseHealthy: false, wantStatus: http.StatusServiceUnavailable, wantBody: `"status":"not_ready"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			registerHealthRoutes(mux, func(context.Context) (bool, string) {
				if test.databaseHealthy {
					return true, "ok"
				}
				return false, "unavailable"
			}, time.Now())

			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil))
			if response.Code != test.wantStatus {
				t.Fatalf("readiness status = %d, want %d", response.Code, test.wantStatus)
			}
			if body := response.Body.String(); !strings.Contains(body, test.wantBody) {
				t.Fatalf("readiness body %q does not contain %q", body, test.wantBody)
			}
		})
	}
}

func TestMetricsRouteRetainsInternalTokenMiddleware(t *testing.T) {
	t.Setenv("METRICS_TOKEN", "health-route-test-token")
	mux := http.NewServeMux()
	registerHealthRoutes(mux, func(context.Context) (bool, string) { return true, "ok" }, time.Now())

	unauthorized := httptest.NewRecorder()
	mux.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/api/v1/health/metrics", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("metrics without internal token status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/health/metrics", nil)
	request.Header.Set("X-Internal-Token", os.Getenv("METRICS_TOKEN"))
	authorized := httptest.NewRecorder()
	mux.ServeHTTP(authorized, request)
	if authorized.Code != http.StatusOK {
		t.Fatalf("metrics with internal token status = %d, want %d", authorized.Code, http.StatusOK)
	}
}
