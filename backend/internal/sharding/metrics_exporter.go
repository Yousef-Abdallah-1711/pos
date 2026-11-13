package sharding

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// PrometheusMetricsExporter exports shard metrics in Prometheus format
type PrometheusMetricsExporter struct {
	healthChecker *HealthChecker
	mu            sync.RWMutex
}

// NewPrometheusMetricsExporter creates a new Prometheus metrics exporter
func NewPrometheusMetricsExporter(hc *HealthChecker) *PrometheusMetricsExporter {
	return &PrometheusMetricsExporter{
		healthChecker: hc,
	}
}

// ExportMetrics exports all metrics in Prometheus format
func (e *PrometheusMetricsExporter) ExportMetrics() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var metrics strings.Builder

	// Header
	metrics.WriteString("# HELP shard_metrics POS Shard Metrics\n")
	metrics.WriteString("# TYPE shard_metrics gauge\n")
	metrics.WriteString("\n")

	// Get health status for all shards
	allHealth := e.healthChecker.GetAllHealth()

	// Export shard health status
	metrics.WriteString("# HELP shard_health_status Shard health status (0=unknown, 1=healthy, 2=degraded, 3=unhealthy)\n")
	metrics.WriteString("# TYPE shard_health_status gauge\n")

	statusMap := map[HealthStatus]int{
		HealthStatusUnknown:   0,
		HealthStatusHealthy:   1,
		HealthStatusDegraded:  2,
		HealthStatusUnhealthy: 3,
	}

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_health_status{shard_id=\"%d\",host=\"%s\"} %d\n",
			health.ShardID,
			"", // Host would be added from shard info
			statusMap[health.Status],
		))
	}
	metrics.WriteString("\n")

	// Export response times
	metrics.WriteString("# HELP shard_response_time_ms Shard response time in milliseconds\n")
	metrics.WriteString("# TYPE shard_response_time_ms gauge\n")

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_response_time_ms{shard_id=\"%d\"} %.2f\n",
			health.ShardID,
			float64(health.ResponseTime.Milliseconds()),
		))
	}
	metrics.WriteString("\n")

	// Export error rates
	metrics.WriteString("# HELP shard_error_rate_percentage Shard error rate in percentage\n")
	metrics.WriteString("# TYPE shard_error_rate_percentage gauge\n")

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_error_rate_percentage{shard_id=\"%d\"} %.2f\n",
			health.ShardID,
			health.ErrorRate*100,
		))
	}
	metrics.WriteString("\n")

	// Export queries per second
	metrics.WriteString("# HELP shard_queries_per_second Shard queries per second\n")
	metrics.WriteString("# TYPE shard_queries_per_second gauge\n")

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_queries_per_second{shard_id=\"%d\"} %.2f\n",
			health.ShardID,
			health.QueriesPerSecond,
		))
	}
	metrics.WriteString("\n")

	// Export database sizes
	metrics.WriteString("# HELP shard_database_size_bytes Shard database size in bytes\n")
	metrics.WriteString("# TYPE shard_database_size_bytes gauge\n")

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_database_size_bytes{shard_id=\"%d\"} %d\n",
			health.ShardID,
			health.DatabaseSize,
		))
	}
	metrics.WriteString("\n")

	// Export record counts
	metrics.WriteString("# HELP shard_record_count Total records in shard\n")
	metrics.WriteString("# TYPE shard_record_count gauge\n")

	for _, health := range allHealth {
		metrics.WriteString(fmt.Sprintf(
			"shard_record_count{shard_id=\"%d\"} %d\n",
			health.ShardID,
			health.RecordCount,
		))
	}
	metrics.WriteString("\n")

	// Export reachability
	metrics.WriteString("# HELP shard_reachable Shard reachability (1=reachable, 0=unreachable)\n")
	metrics.WriteString("# TYPE shard_reachable gauge\n")

	for _, health := range allHealth {
		reachable := 0
		if health.IsReachable {
			reachable = 1
		}
		metrics.WriteString(fmt.Sprintf(
			"shard_reachable{shard_id=\"%d\"} %d\n",
			health.ShardID,
			reachable,
		))
	}
	metrics.WriteString("\n")

	// System-level metrics
	summary := e.healthChecker.GetHealthSummary()

	metrics.WriteString("# HELP system_shards_total Total number of shards\n")
	metrics.WriteString("# TYPE system_shards_total gauge\n")
	metrics.WriteString(fmt.Sprintf("system_shards_total %d\n", summary["total_shards"]))
	metrics.WriteString("\n")

	metrics.WriteString("# HELP system_shards_healthy Number of healthy shards\n")
	metrics.WriteString("# TYPE system_shards_healthy gauge\n")
	metrics.WriteString(fmt.Sprintf("system_shards_healthy %d\n", summary["healthy_shards"]))
	metrics.WriteString("\n")

	metrics.WriteString("# HELP system_shards_degraded Number of degraded shards\n")
	metrics.WriteString("# TYPE system_shards_degraded gauge\n")
	metrics.WriteString(fmt.Sprintf("system_shards_degraded %d\n", summary["degraded_shards"]))
	metrics.WriteString("\n")

	metrics.WriteString("# HELP system_shards_unhealthy Number of unhealthy shards\n")
	metrics.WriteString("# TYPE system_shards_unhealthy gauge\n")
	metrics.WriteString(fmt.Sprintf("system_shards_unhealthy %d\n", summary["unhealthy_shards"]))
	metrics.WriteString("\n")

	metrics.WriteString("# HELP system_average_latency_ms Average latency across all shards\n")
	metrics.WriteString("# TYPE system_average_latency_ms gauge\n")
	latencyMs := summary["average_latency"].(time.Duration).Milliseconds()
	metrics.WriteString(fmt.Sprintf("system_average_latency_ms %d\n", latencyMs))
	metrics.WriteString("\n")

	metrics.WriteString("# HELP system_health System health status (1=healthy, 0=degraded)\n")
	metrics.WriteString("# TYPE system_health gauge\n")
	systemHealth := 0
	if summary["system_healthy"].(bool) {
		systemHealth = 1
	}
	metrics.WriteString(fmt.Sprintf("system_health %d\n", systemHealth))
	metrics.WriteString("\n")

	return metrics.String()
}

// ExportMetricsJSON exports metrics in JSON format
func (e *PrometheusMetricsExporter) ExportMetricsJSON() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	allHealth := e.healthChecker.GetAllHealth()
	summary := e.healthChecker.GetHealthSummary()

	shards := make([]map[string]interface{}, 0)
	for _, health := range allHealth {
		shard := map[string]interface{}{
			"shard_id":      health.ShardID,
			"status":        health.Status,
			"response_time": health.ResponseTime.String(),
			"error_rate":    health.ErrorRate,
			"qps":           health.QueriesPerSecond,
			"database_size": health.DatabaseSize,
			"record_count":  health.RecordCount,
			"is_reachable":  health.IsReachable,
			"message":       health.Message,
		}
		shards = append(shards, shard)
	}

	return map[string]interface{}{
		"timestamp": time.Now(),
		"system":    summary,
		"shards":    shards,
	}
}

// GenerateAlerts generates alerts based on health status
func (e *PrometheusMetricsExporter) GenerateAlerts() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var alerts []string

	allHealth := e.healthChecker.GetAllHealth()

	for _, health := range allHealth {
		if health.Status == HealthStatusUnhealthy {
			alerts = append(alerts, fmt.Sprintf(
				"CRITICAL: Shard %d is unhealthy - %s",
				health.ShardID,
				health.Message,
			))
		} else if health.Status == HealthStatusDegraded {
			alerts = append(alerts, fmt.Sprintf(
				"WARNING: Shard %d is degraded - %s",
				health.ShardID,
				health.Message,
			))
		}

		if !health.IsReachable {
			alerts = append(alerts, fmt.Sprintf(
				"CRITICAL: Shard %d is unreachable",
				health.ShardID,
			))
		}
	}

	summary := e.healthChecker.GetHealthSummary()
	if !summary["system_healthy"].(bool) {
		alerts = append(alerts, fmt.Sprintf(
			"CRITICAL: System is not healthy - %d healthy, %d degraded, %d unhealthy",
			summary["healthy_shards"],
			summary["degraded_shards"],
			summary["unhealthy_shards"],
		))
	}

	return alerts
}

// LogMetrics logs current metrics to application logs
func (e *PrometheusMetricsExporter) LogMetrics() {
	e.mu.RLock()
	defer e.mu.RUnlock()

	allHealth := e.healthChecker.GetAllHealth()
	summary := e.healthChecker.GetHealthSummary()

	log.Printf("[METRICS] System Health Summary:")
	log.Printf("[METRICS]   Total Shards: %d", summary["total_shards"])
	log.Printf("[METRICS]   Healthy: %d", summary["healthy_shards"])
	log.Printf("[METRICS]   Degraded: %d", summary["degraded_shards"])
	log.Printf("[METRICS]   Unhealthy: %d", summary["unhealthy_shards"])
	log.Printf("[METRICS]   Average Latency: %v", summary["average_latency"])

	for _, health := range allHealth {
		log.Printf("[METRICS] Shard %d: %s | Latency: %v | Errors: %.2f%% | QPS: %.2f",
			health.ShardID,
			health.Status,
			health.ResponseTime,
			health.ErrorRate*100,
			health.QueriesPerSecond,
		)
	}
}
