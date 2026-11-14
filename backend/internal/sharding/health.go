package sharding

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// HealthStatus represents the health state of a shard
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// ShardHealth contains detailed health information about a shard
type ShardHealth struct {
	ShardID              int
	Status               HealthStatus
	LastCheckTime        time.Time
	ResponseTime         time.Duration
	ConnectionPoolSize   int
	AvailableConnections int
	QueriesPerSecond     float64
	ErrorRate            float64
	ReplicationLag       time.Duration
	Message              string
	IsReachable          bool
	DatabaseSize         int64 // in bytes
	RecordCount          int64
}

// HealthChecker monitors the health of all shards
type HealthChecker struct {
	router             *ShardRouter
	shards             []ShardInfo
	healthCache        map[int]*ShardHealth
	metrics            map[int]*ShardMetrics
	mu                 sync.RWMutex
	checkInterval      time.Duration
	degradedThreshold  time.Duration // Latency threshold for degraded status
	unhealthyThreshold time.Duration // Latency threshold for unhealthy status
	maxErrorRate       float64       // Error rate threshold (0.1 = 10%)
	maxReplicationLag  time.Duration // Max acceptable replication lag
	stopChan           chan bool
	isRunning          bool
}

// ShardMetrics tracks metrics for a shard
type ShardMetrics struct {
	TotalQueries  int64
	FailedQueries int64
	TotalLatency  time.Duration
	MaxLatency    time.Duration
	MinLatency    time.Duration
	LastResetTime time.Time
	mu            sync.RWMutex
}

// NewHealthChecker creates a new health checker instance
func NewHealthChecker(router *ShardRouter, checkInterval time.Duration) *HealthChecker {
	shards := router.GetAllShards()

	hc := &HealthChecker{
		router:             router,
		shards:             shards,
		healthCache:        make(map[int]*ShardHealth),
		metrics:            make(map[int]*ShardMetrics),
		checkInterval:      checkInterval,
		degradedThreshold:  500 * time.Millisecond,  // 500ms = degraded
		unhealthyThreshold: 2000 * time.Millisecond, // 2 seconds = unhealthy
		maxErrorRate:       0.05,                    // 5% error rate
		maxReplicationLag:  10 * time.Second,        // 10 seconds max lag
		stopChan:           make(chan bool),
	}

	// Initialize metrics for each shard
	for _, shard := range shards {
		hc.metrics[shard.ID] = &ShardMetrics{
			LastResetTime: time.Now(),
		}
		hc.healthCache[shard.ID] = &ShardHealth{
			ShardID: shard.ID,
			Status:  HealthStatusUnknown,
		}
	}

	return hc
}

// Start begins health checking in a background goroutine
func (hc *HealthChecker) Start() {
	if hc.isRunning {
		log.Printf("[HEALTH] Health checker already running")
		return
	}

	hc.isRunning = true
	log.Printf("[HEALTH] Starting health checker with interval: %v", hc.checkInterval)

	go hc.checkLoop()
}

// Stop stops the health checking loop
func (hc *HealthChecker) Stop() {
	if !hc.isRunning {
		return
	}

	log.Printf("[HEALTH] Stopping health checker")
	hc.stopChan <- true
	hc.isRunning = false
}

// checkLoop runs the continuous health checking loop
func (hc *HealthChecker) checkLoop() {
	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hc.stopChan:
			log.Printf("[HEALTH] Health check loop stopped")
			return
		case <-ticker.C:
			hc.checkAllShards()
		}
	}
}

// checkAllShards checks health of all shards
func (hc *HealthChecker) checkAllShards() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	log.Printf("[HEALTH] Running health checks for %d shards", len(hc.shards))

	for i, shard := range hc.shards {
		health := hc.checkShard(&shard)
		hc.healthCache[shard.ID] = health

		// Log health status
		statusEmoji := "✓"
		if health.Status == HealthStatusDegraded {
			statusEmoji = "⚠"
		} else if health.Status == HealthStatusUnhealthy {
			statusEmoji = "✗"
		}

		log.Printf("[HEALTH] Shard %d (%s) %s | Latency: %v | Error Rate: %.1f%% | Status: %s",
			shard.ID,
			shard.Host,
			statusEmoji,
			health.ResponseTime,
			health.ErrorRate*100,
			health.Status,
		)

		if i%5 == 4 { // Log separator every 5 shards
			log.Printf("[HEALTH] ---")
		}
	}
}

// checkShard checks the health of a single shard
func (hc *HealthChecker) checkShard(shard *ShardInfo) *ShardHealth {
	health := &ShardHealth{
		ShardID:       shard.ID,
		LastCheckTime: time.Now(),
		IsReachable:   false,
	}

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()

	// Try to get connection
	conn, err := hc.router.GetConnection(fmt.Sprintf("%d", shard.ID))
	health.ResponseTime = time.Since(start)

	if err != nil {
		health.Status = HealthStatusUnhealthy
		health.Message = fmt.Sprintf("Connection failed: %v", err)
		health.IsReachable = false
		return health
	}

	if conn == nil {
		health.Status = HealthStatusUnhealthy
		health.Message = "No connection available"
		health.IsReachable = false
		return health
	}

	health.IsReachable = true

	// Test database connectivity
	start = time.Now()
	err = conn.PingContext(ctx)
	health.ResponseTime = time.Since(start)

	if err != nil {
		health.Status = HealthStatusUnhealthy
		health.Message = fmt.Sprintf("Ping failed: %v", err)
		return health
	}

	// Get database stats
	hc.getShardStats(ctx, conn, health)

	// Get metrics for this shard
	metrics := hc.metrics[shard.ID]
	health.ErrorRate = hc.calculateErrorRate(metrics)
	health.QueriesPerSecond = hc.calculateQPS(metrics)

	// Determine health status based on thresholds
	if health.ResponseTime > hc.unhealthyThreshold {
		health.Status = HealthStatusUnhealthy
		health.Message = fmt.Sprintf("Response time too high: %v", health.ResponseTime)
	} else if health.ResponseTime > hc.degradedThreshold {
		health.Status = HealthStatusDegraded
		health.Message = fmt.Sprintf("Response time degraded: %v", health.ResponseTime)
	} else if health.ErrorRate > hc.maxErrorRate {
		health.Status = HealthStatusDegraded
		health.Message = fmt.Sprintf("Error rate high: %.2f%%", health.ErrorRate*100)
	} else {
		health.Status = HealthStatusHealthy
		health.Message = "Shard is healthy"
	}

	return health
}

// getShardStats retrieves statistics about a shard
func (hc *HealthChecker) getShardStats(ctx context.Context, conn *sql.DB, health *ShardHealth) {
	// Get database size
	var dbSize int64
	query := "SELECT pg_database_size(current_database())"
	err := conn.QueryRowContext(ctx, query).Scan(&dbSize)
	if err == nil {
		health.DatabaseSize = dbSize
	}

	// Get record count across all tables
	var recordCount int64
	tables := []string{"orders", "order_items", "payments", "customers", "restaurants"}
	for _, table := range tables {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := conn.QueryRowContext(ctx, query).Scan(&count)
		if err == nil {
			recordCount += count
		}
	}
	health.RecordCount = recordCount
}

// RecordQuery records metrics for a query
func (hc *HealthChecker) RecordQuery(shardID int, latency time.Duration, err error) {
	hc.mu.RLock()
	metrics, exists := hc.metrics[shardID]
	hc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.TotalQueries++
	metrics.TotalLatency += latency

	if latency > metrics.MaxLatency {
		metrics.MaxLatency = latency
	}
	if latency < metrics.MinLatency || metrics.MinLatency == 0 {
		metrics.MinLatency = latency
	}

	if err != nil {
		metrics.FailedQueries++
	}
}

// GetShardHealth returns the current health status of a shard
func (hc *HealthChecker) GetShardHealth(shardID int) *ShardHealth {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	health, exists := hc.healthCache[shardID]
	if !exists {
		return &ShardHealth{
			ShardID: shardID,
			Status:  HealthStatusUnknown,
			Message: "Shard not found",
		}
	}

	return health
}

// GetAllHealth returns health status of all shards
func (hc *HealthChecker) GetAllHealth() map[int]*ShardHealth {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	healthMap := make(map[int]*ShardHealth)
	for id, health := range hc.healthCache {
		healthMap[id] = health
	}

	return healthMap
}

// GetHealthSummary returns overall health summary
func (hc *HealthChecker) GetHealthSummary() map[string]interface{} {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	var totalShards, healthyShards, degradedShards, unhealthyShards int
	var avgLatency time.Duration

	for _, health := range hc.healthCache {
		totalShards++
		avgLatency += health.ResponseTime

		switch health.Status {
		case HealthStatusHealthy:
			healthyShards++
		case HealthStatusDegraded:
			degradedShards++
		case HealthStatusUnhealthy:
			unhealthyShards++
		}
	}

	if totalShards > 0 {
		avgLatency = avgLatency / time.Duration(totalShards)
	}

	return map[string]interface{}{
		"total_shards":     totalShards,
		"healthy_shards":   healthyShards,
		"degraded_shards":  degradedShards,
		"unhealthy_shards": unhealthyShards,
		"average_latency":  avgLatency,
		"system_healthy":   unhealthyShards == 0,
		"last_check_time":  time.Now(),
	}
}

// calculateErrorRate calculates error rate from metrics
func (hc *HealthChecker) calculateErrorRate(metrics *ShardMetrics) float64 {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	if metrics.TotalQueries == 0 {
		return 0
	}

	return float64(metrics.FailedQueries) / float64(metrics.TotalQueries)
}

// calculateQPS calculates queries per second
func (hc *HealthChecker) calculateQPS(metrics *ShardMetrics) float64 {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	elapsed := time.Since(metrics.LastResetTime).Seconds()
	if elapsed == 0 {
		return 0
	}

	return float64(metrics.TotalQueries) / elapsed
}

// ResetMetrics resets metrics for a shard
func (hc *HealthChecker) ResetMetrics(shardID int) {
	hc.mu.RLock()
	metrics, exists := hc.metrics[shardID]
	hc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	metrics.TotalQueries = 0
	metrics.FailedQueries = 0
	metrics.TotalLatency = 0
	metrics.MaxLatency = 0
	metrics.MinLatency = 0
	metrics.LastResetTime = time.Now()
}

// GetMetrics returns metrics for a shard
func (hc *HealthChecker) GetMetrics(shardID int) map[string]interface{} {
	hc.mu.RLock()
	metrics, exists := hc.metrics[shardID]
	hc.mu.RUnlock()

	if !exists {
		return nil
	}

	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	var avgLatency time.Duration
	if metrics.TotalQueries > 0 {
		avgLatency = metrics.TotalLatency / time.Duration(metrics.TotalQueries)
	}

	return map[string]interface{}{
		"total_queries":   metrics.TotalQueries,
		"failed_queries":  metrics.FailedQueries,
		"error_rate":      float64(metrics.FailedQueries) / float64(metrics.TotalQueries),
		"avg_latency":     avgLatency,
		"max_latency":     metrics.MaxLatency,
		"min_latency":     metrics.MinLatency,
		"qps":             hc.calculateQPS(metrics),
		"last_reset_time": metrics.LastResetTime,
	}
}

// SetHealthThresholds allows customization of health thresholds
func (hc *HealthChecker) SetHealthThresholds(
	degraded, unhealthy time.Duration,
	maxErrorRate float64,
	maxReplicationLag time.Duration) {

	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.degradedThreshold = degraded
	hc.unhealthyThreshold = unhealthy
	hc.maxErrorRate = maxErrorRate
	hc.maxReplicationLag = maxReplicationLag

	log.Printf("[HEALTH] Thresholds updated - Degraded: %v, Unhealthy: %v, Error Rate: %.2f%%",
		degraded, unhealthy, maxErrorRate*100)
}
