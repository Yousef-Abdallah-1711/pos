package scaling

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	ServerID         string
	Timestamp        time.Time
	IsHealthy        bool
	ResponseTime     time.Duration
	StatusCode       int
	ErrorMessage     string
	CPUUsage         float64
	MemoryUsage      float64
	DiskUsagePercent float64
	DatabaseLatency  time.Duration
	CacheLatency     time.Duration
}

// HealthChecker continuously monitors server health
type HealthChecker struct {
	loadBalancer           *LoadBalancer
	healthProbes           map[string]HealthProbe
	checkInterval          time.Duration
	checkTimeout           time.Duration
	unhealthyCount         map[string]int // consecutive failures
	maxConsecutiveFailures int
	results                []HealthCheckResult
	mu                     sync.RWMutex
	stopChannel            chan bool
	isRunning              bool
}

// HealthProbe defines how to check health of a server
type HealthProbe struct {
	ServerID       string
	Endpoint       string
	Port           int
	Path           string // Usually "/health"
	Method         string
	ExpectedStatus int
	Timeout        time.Duration
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(
	loadBalancer *LoadBalancer,
	checkInterval time.Duration,
	checkTimeout time.Duration,
	maxFailures int,
) *HealthChecker {
	return &HealthChecker{
		loadBalancer:           loadBalancer,
		healthProbes:           make(map[string]HealthProbe),
		checkInterval:          checkInterval,
		checkTimeout:           checkTimeout,
		unhealthyCount:         make(map[string]int),
		maxConsecutiveFailures: maxFailures,
		results:                make([]HealthCheckResult, 0),
		stopChannel:            make(chan bool),
	}
}

// RegisterHealthProbe registers a health check probe for a server
func (hc *HealthChecker) RegisterHealthProbe(
	serverID string,
	endpoint string,
	port int,
	path string,
) {
	log.Printf("[HEALTH] Registering probe: %s -> %s:%d%s", serverID, endpoint, port, path)

	probe := HealthProbe{
		ServerID:       serverID,
		Endpoint:       endpoint,
		Port:           port,
		Path:           path,
		Method:         "GET",
		ExpectedStatus: 200,
		Timeout:        hc.checkTimeout,
	}

	hc.mu.Lock()
	hc.healthProbes[serverID] = probe
	hc.unhealthyCount[serverID] = 0
	hc.mu.Unlock()
}

// Start begins health checking
func (hc *HealthChecker) Start(ctx context.Context) {
	hc.mu.Lock()
	if hc.isRunning {
		hc.mu.Unlock()
		return
	}
	hc.isRunning = true
	hc.mu.Unlock()

	log.Printf("[HEALTH] Starting health checker (interval: %v)", hc.checkInterval)

	go func() {
		ticker := time.NewTicker(hc.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				hc.checkAllServers(ctx)
			case <-hc.stopChannel:
				log.Printf("[HEALTH] Stopping health checker")
				hc.mu.Lock()
				hc.isRunning = false
				hc.mu.Unlock()
				return
			case <-ctx.Done():
				hc.mu.Lock()
				hc.isRunning = false
				hc.mu.Unlock()
				return
			}
		}
	}()
}

// Stop stops the health checker
func (hc *HealthChecker) Stop() {
	hc.stopChannel <- true
}

// checkAllServers checks health of all registered servers
func (hc *HealthChecker) checkAllServers(ctx context.Context) {
	hc.mu.RLock()
	probes := make([]HealthProbe, 0)
	for _, probe := range hc.healthProbes {
		probes = append(probes, probe)
	}
	hc.mu.RUnlock()

	for _, probe := range probes {
		go hc.checkServer(ctx, probe)
	}
}

// checkServer performs a health check on a single server
func (hc *HealthChecker) checkServer(ctx context.Context, probe HealthProbe) {
	result := HealthCheckResult{
		ServerID:  probe.ServerID,
		Timestamp: time.Now(),
	}

	start := time.Now()

	// Create HTTP request
	url := fmt.Sprintf("http://%s:%d%s", probe.Endpoint, probe.Port, probe.Path)
	req, err := http.NewRequestWithContext(ctx, probe.Method, url, nil)
	if err != nil {
		result.IsHealthy = false
		result.ErrorMessage = fmt.Sprintf("Request creation failed: %v", err)
		result.ResponseTime = time.Since(start)
		hc.recordHealthCheck(result)
		hc.handleUnhealthyServer(probe.ServerID)
		return
	}

	// Perform request
	client := &http.Client{
		Timeout: probe.Timeout,
	}

	resp, err := client.Do(req)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.IsHealthy = false
		result.ErrorMessage = fmt.Sprintf("Request failed: %v", err)
		hc.recordHealthCheck(result)
		hc.handleUnhealthyServer(probe.ServerID)
		return
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.IsHealthy = resp.StatusCode == probe.ExpectedStatus

	if !result.IsHealthy {
		result.ErrorMessage = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
		hc.handleUnhealthyServer(probe.ServerID)
	} else {
		hc.handleHealthyServer(probe.ServerID)
	}

	hc.recordHealthCheck(result)
}

// recordHealthCheck records a health check result
func (hc *HealthChecker) recordHealthCheck(result HealthCheckResult) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.results = append(hc.results, result)

	// Keep only last 10000 results
	if len(hc.results) > 10000 {
		hc.results = hc.results[1:]
	}

	status := "healthy"
	if !result.IsHealthy {
		status = "unhealthy"
	}

	log.Printf("[HEALTH] Check result: %s is %s (response: %v)", result.ServerID, status, result.ResponseTime)
}

// handleHealthyServer handles a healthy server status
func (hc *HealthChecker) handleHealthyServer(serverID string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Reset failure count
	if hc.unhealthyCount[serverID] > 0 {
		hc.unhealthyCount[serverID] = 0

		// Mark server as healthy
		_ = hc.loadBalancer.SetServerHealth(serverID, true)
		log.Printf("[HEALTH] Server %s recovered", serverID)
	}
}

// handleUnhealthyServer handles an unhealthy server status
func (hc *HealthChecker) handleUnhealthyServer(serverID string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.unhealthyCount[serverID]++

	// Mark unhealthy after max consecutive failures
	if hc.unhealthyCount[serverID] >= hc.maxConsecutiveFailures {
		_ = hc.loadBalancer.SetServerHealth(serverID, false)
		log.Printf("[HEALTH] Server %s marked unhealthy (%d failures)",
			serverID, hc.unhealthyCount[serverID])
	}
}

// GetHealthStatus returns current health status of a server
func (hc *HealthChecker) GetHealthStatus(serverID string) map[string]interface{} {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	var lastResult *HealthCheckResult
	for i := len(hc.results) - 1; i >= 0; i-- {
		if hc.results[i].ServerID == serverID {
			lastResult = &hc.results[i]
			break
		}
	}

	if lastResult == nil {
		return map[string]interface{}{
			"server_id": serverID,
			"status":    "unknown",
		}
	}

	status := "healthy"
	if !lastResult.IsHealthy {
		status = "unhealthy"
	}

	return map[string]interface{}{
		"server_id":            serverID,
		"status":               status,
		"response_time_ms":     lastResult.ResponseTime.Milliseconds(),
		"status_code":          lastResult.StatusCode,
		"error_message":        lastResult.ErrorMessage,
		"cpu_usage":            fmt.Sprintf("%.1f%%", lastResult.CPUUsage),
		"memory_usage":         fmt.Sprintf("%.1f%%", lastResult.MemoryUsage),
		"disk_usage":           fmt.Sprintf("%.1f%%", lastResult.DiskUsagePercent),
		"database_latency_ms":  lastResult.DatabaseLatency.Milliseconds(),
		"cache_latency_ms":     lastResult.CacheLatency.Milliseconds(),
		"last_check":           lastResult.Timestamp,
		"consecutive_failures": hc.unhealthyCount[serverID],
	}
}

// GetHealthHistory returns health check history for a server
func (hc *HealthChecker) GetHealthHistory(serverID string, limit int) []HealthCheckResult {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	var results []HealthCheckResult

	for _, result := range hc.results {
		if result.ServerID == serverID {
			results = append(results, result)
		}
	}

	if len(results) > limit {
		results = results[len(results)-limit:]
	}

	return results
}

// GetGlobalHealthStatus returns global health status of all servers
func (hc *HealthChecker) GetGlobalHealthStatus() map[string]interface{} {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	healthyCount := 0
	unhealthyCount := 0
	var avgResponseTime time.Duration
	checkCount := 0

	for _, result := range hc.results {
		if len(hc.results) > 0 && result.Timestamp.After(hc.results[len(hc.results)-1].Timestamp.Add(-hc.checkInterval)) {
			if result.IsHealthy {
				healthyCount++
			} else {
				unhealthyCount++
			}
			avgResponseTime += result.ResponseTime
			checkCount++
		}
	}

	if checkCount > 0 {
		avgResponseTime /= time.Duration(checkCount)
	}

	return map[string]interface{}{
		"total_servers":        len(hc.healthProbes),
		"healthy_servers":      healthyCount,
		"unhealthy_servers":    unhealthyCount,
		"avg_response_time_ms": avgResponseTime.Milliseconds(),
		"total_checks":         len(hc.results),
		"is_running":           hc.isRunning,
	}
}

// CalculateHealthPercentage calculates server health percentage
func (hc *HealthChecker) CalculateHealthPercentage(serverID string) int {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	totalChecks := 0
	healthyChecks := 0

	for _, result := range hc.results {
		if result.ServerID == serverID {
			totalChecks++
			if result.IsHealthy {
				healthyChecks++
			}
		}
	}

	if totalChecks == 0 {
		return 100
	}

	return (healthyChecks * 100) / totalChecks
}
