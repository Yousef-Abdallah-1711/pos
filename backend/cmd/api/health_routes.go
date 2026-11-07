package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"pos-saas/internal/config"
	"pos-saas/internal/middleware"
	"pos-saas/internal/observability"
)

type databaseHealthCheck func(context.Context) (bool, string)

func registerHealthRoutes(mux *http.ServeMux, check databaseHealthCheck, startedAt time.Time) {
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		dbHealthy, dbStatus := check(r.Context())
		status := "ok"
		if !dbHealthy {
			status = "degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     status,
			"uptime_sec": int64(time.Since(startedAt).Seconds()),
			"db_healthy": dbHealthy,
			"db_status":  dbStatus,
		})
	})

	mux.HandleFunc("GET /api/v1/health/ready", func(w http.ResponseWriter, r *http.Request) {
		dbHealthy, dbStatus := check(r.Context())
		statusCode := http.StatusOK
		status := "ready"
		if !dbHealthy {
			statusCode = http.StatusServiceUnavailable
			status = "not_ready"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     status,
			"uptime_sec": int64(time.Since(startedAt).Seconds()),
			"db_healthy": dbHealthy,
			"db_status":  dbStatus,
		})
	})

	mux.Handle("GET /api/v1/health/metrics", middleware.RequireInternalToken("METRICS_TOKEN", config.IsProductionMode())(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    observability.Snapshot(),
		})
	})))
}
