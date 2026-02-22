//go:build legacy
// +build legacy

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"pos-saas/internal/sharding"
)

// HealthMonitoringHandler manages HTTP API for health monitoring
type HealthMonitoringHandler struct {
	healthChecker *sharding.HealthChecker
}

// NewHealthMonitoringHandler creates a new health monitoring handler
func NewHealthMonitoringHandler(hc *sharding.HealthChecker) *HealthMonitoringHandler {
	return &HealthMonitoringHandler{
		healthChecker: hc,
	}
}

// GetSystemHealth returns overall system health status
// GET /api/v1/admin/health/system
func (h *HealthMonitoringHandler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] GetSystemHealth called")

	summary := h.healthChecker.GetHealthSummary()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-System-Healthy", fmt.Sprintf("%v", summary["system_healthy"]))

	json.NewEncoder(w).Encode(summary)
}

// GetShardHealth returns health status of a specific shard
// GET /api/v1/admin/health/shards/{shard_id}
func (h *HealthMonitoringHandler) GetShardHealth(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.PathValue("shard_id")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		http.Error(w, "Invalid shard ID", http.StatusBadRequest)
		return
	}

	log.Printf("[MONITORING] GetShardHealth for shard %d", shardID)

	health := h.healthChecker.GetShardHealth(shardID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Shard-Status", string(health.Status))
	w.Header().Set("X-Shard-Reachable", fmt.Sprintf("%v", health.IsReachable))

	json.NewEncoder(w).Encode(health)
}

// GetAllShardHealth returns health status of all shards
// GET /api/v1/admin/health/shards
func (h *HealthMonitoringHandler) GetAllShardHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] GetAllShardHealth called")

	health := h.healthChecker.GetAllHealth()

	type ShardHealthResponse struct {
		ShardID              int    `json:"shard_id"`
		Status               string `json:"status"`
		ResponseTime         string `json:"response_time"`
		ErrorRate            string `json:"error_rate"`
		IsReachable          bool   `json:"is_reachable"`
		Message              string `json:"message"`
		QueriesPerSecond     string `json:"queries_per_second"`
		DatabaseSize         int64  `json:"database_size"`
		RecordCount          int64  `json:"record_count"`
	}

	var shards []ShardHealthResponse
	for _, h := range health {
		shard := ShardHealthResponse{
			ShardID:      h.ShardID,
			Status:       string(h.Status),
			ResponseTime: fmt.Sprintf("%.1fms", float64(h.ResponseTime.Milliseconds())),
			ErrorRate:    fmt.Sprintf("%.2f%%", h.ErrorRate*100),
			IsReachable:  h.IsReachable,
			Message:      h.Message,
			QueriesPerSecond: fmt.Sprintf("%.2f", h.QueriesPerSecond),
			DatabaseSize: h.DatabaseSize,
			RecordCount:  h.RecordCount,
		}
		shards = append(shards, shard)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(shards),
		"shards": shards,
	})
}

// GetShardMetrics returns detailed metrics for a shard
// GET /api/v1/admin/metrics/shards/{shard_id}
func (h *HealthMonitoringHandler) GetShardMetrics(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.PathValue("shard_id")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		http.Error(w, "Invalid shard ID", http.StatusBadRequest)
		return
	}

	log.Printf("[MONITORING] GetShardMetrics for shard %d", shardID)

	metrics := h.healthChecker.GetMetrics(shardID)
	if metrics == nil {
		http.Error(w, "Shard not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shard_id": shardID,
		"metrics":  metrics,
	})
}

// ResetShardMetrics resets metrics for a shard
// POST /api/v1/admin/metrics/shards/{shard_id}/reset
func (h *HealthMonitoringHandler) ResetShardMetrics(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.PathValue("shard_id")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		http.Error(w, "Invalid shard ID", http.StatusBadRequest)
		return
	}

	log.Printf("[MONITORING] ResetShardMetrics for shard %d", shardID)

	h.healthChecker.ResetMetrics(shardID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shard_id": shardID,
		"message":  "Metrics reset successfully",
	})
}

// GetSystemMetrics returns system-wide metrics
// GET /api/v1/admin/metrics/system
func (h *HealthMonitoringHandler) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] GetSystemMetrics called")

	summary := h.healthChecker.GetHealthSummary()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"system_health": summary,
		"timestamp":     summary["last_check_time"],
	})
}

// StartHealthChecks starts the background health checking
// POST /api/v1/admin/health/start
func (h *HealthMonitoringHandler) StartHealthChecks(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] StartHealthChecks called")

	h.healthChecker.Start()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Health checking started",
		"status":  "started",
	})
}

// StopHealthChecks stops the background health checking
// POST /api/v1/admin/health/stop
func (h *HealthMonitoringHandler) StopHealthChecks(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] StopHealthChecks called")

	h.healthChecker.Stop()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Health checking stopped",
		"status":  "stopped",
	})
}

// UpdateHealthThresholds updates health check thresholds
// POST /api/v1/admin/health/thresholds
func (h *HealthMonitoringHandler) UpdateHealthThresholds(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MONITORING] UpdateHealthThresholds called")

	var req struct {
		DegradedThreshold    string  `json:"degraded_threshold_ms"`
		UnhealthyThreshold   string  `json:"unhealthy_threshold_ms"`
		MaxErrorRate         float64 `json:"max_error_rate"`
		MaxReplicationLag    string  `json:"max_replication_lag_ms"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse durations
	degradedMs, _ := strconv.Atoi(req.DegradedThreshold)
	unhealthyMs, _ := strconv.Atoi(req.UnhealthyThreshold)
	lagMs, _ := strconv.Atoi(req.MaxReplicationLag)

	h.healthChecker.SetHealthThresholds(
		GetDurationFromMs(degradedMs),
		GetDurationFromMs(unhealthyMs),
		req.MaxErrorRate,
		GetDurationFromMs(lagMs),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":          "Thresholds updated successfully",
		"degraded_ms":      degradedMs,
		"unhealthy_ms":     unhealthyMs,
		"max_error_rate":   req.MaxErrorRate,
		"max_replication_lag_ms": lagMs,
	})
}

// HealthResponse is a standard health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SimpleHealthCheck is a lightweight health endpoint
// GET /health
func (h *HealthMonitoringHandler) SimpleHealthCheck(w http.ResponseWriter, r *http.Request) {
	summary := h.healthChecker.GetHealthSummary()

	status := "UP"
	if !summary["system_healthy"].(bool) {
		status = "DEGRADED"
	}

	response := HealthResponse{
		Status:  status,
		Message: fmt.Sprintf("System health: %d healthy, %d degraded, %d unhealthy",
			summary["healthy_shards"],
			summary["degraded_shards"],
			summary["unhealthy_shards"],
		),
	}

	w.Header().Set("Content-Type", "application/json")

	if status == "UP" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// Helper function
func GetDurationFromMs(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
