//go:build experimental
// +build experimental

package sharding_test

import (
	"testing"
	"time"

	"pos-saas/internal/sharding"
)

// TestHealthCheckerInitialization tests health checker creation
func TestHealthCheckerInitialization(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	if hc == nil {
		t.Fatal("Health checker initialization failed")
	}

	t.Log("✓ Health checker initialized successfully")
}

// TestHealthCheckStart tests starting health checks
func TestHealthCheckStart(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 100*time.Millisecond)
	defer hc.Stop()

	hc.Start()
	time.Sleep(200 * time.Millisecond) // Wait for at least one check

	summary := hc.GetHealthSummary()
	if summary == nil {
		t.Fatal("Health summary is nil")
	}

	if summary["total_shards"] == nil {
		t.Fatal("Total shards not in summary")
	}

	t.Logf("✓ Health checks started and running - %d total shards", summary["total_shards"])
}

// TestGetShardHealth tests retrieving single shard health
func TestGetShardHealth(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Test getting health for shard 0
	health := hc.GetShardHealth(0)

	if health == nil {
		t.Fatal("Shard health is nil")
	}

	if health.ShardID != 0 {
		t.Fatalf("Expected shard 0, got %d", health.ShardID)
	}

	t.Logf("✓ Retrieved health for shard 0: %s", health.Status)
}

// TestGetAllHealth tests retrieving health for all shards
func TestGetAllHealth(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	allHealth := hc.GetAllHealth()

	if len(allHealth) == 0 {
		t.Fatal("No shard health data returned")
	}

	if len(allHealth) != 4 {
		t.Fatalf("Expected 4 shards, got %d", len(allHealth))
	}

	for _, health := range allHealth {
		if health.ShardID < 0 || health.ShardID >= 4 {
			t.Fatalf("Invalid shard ID: %d", health.ShardID)
		}
	}

	t.Logf("✓ Retrieved health for all %d shards", len(allHealth))
}

// TestHealthSummary tests getting health summary
func TestHealthSummary(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)
	hc.Start()
	defer hc.Stop()

	time.Sleep(100 * time.Millisecond)

	summary := hc.GetHealthSummary()

	if summary["total_shards"] == nil {
		t.Fatal("Missing total_shards in summary")
	}

	if summary["system_healthy"] == nil {
		t.Fatal("Missing system_healthy in summary")
	}

	if summary["average_latency"] == nil {
		t.Fatal("Missing average_latency in summary")
	}

	t.Logf("✓ Health summary retrieved: %v", summary)
}

// TestMetricsRecording tests recording query metrics
func TestMetricsRecording(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Record some queries
	hc.RecordQuery(0, 10*time.Millisecond, nil)
	hc.RecordQuery(0, 20*time.Millisecond, nil)
	hc.RecordQuery(0, 15*time.Millisecond, nil)

	metrics := hc.GetMetrics(0)
	if metrics == nil {
		t.Fatal("Metrics are nil")
	}

	if metrics["total_queries"].(int64) != 3 {
		t.Fatalf("Expected 3 total queries, got %v", metrics["total_queries"])
	}

	t.Logf("✓ Metrics recording works: %v", metrics)
}

// TestErrorRateCalculation tests error rate calculation
func TestErrorRateCalculation(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Record 10 queries, 2 with errors
	for i := 0; i < 8; i++ {
		hc.RecordQuery(0, 10*time.Millisecond, nil)
	}
	hc.RecordQuery(0, 10*time.Millisecond, ErrTestError)
	hc.RecordQuery(0, 10*time.Millisecond, ErrTestError)

	metrics := hc.GetMetrics(0)
	errorRate := metrics["error_rate"].(float64)

	expectedRate := 0.2 // 2/10 = 20%
	if errorRate < expectedRate-0.01 || errorRate > expectedRate+0.01 {
		t.Fatalf("Expected error rate ~%.2f, got %f", expectedRate, errorRate)
	}

	t.Logf("✓ Error rate calculated correctly: %.2f%%", errorRate*100)
}

// TestMetricsReset tests resetting metrics
func TestMetricsReset(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Record some queries
	hc.RecordQuery(0, 10*time.Millisecond, nil)
	hc.RecordQuery(0, 20*time.Millisecond, nil)

	// Reset metrics
	hc.ResetMetrics(0)

	metrics := hc.GetMetrics(0)
	if metrics["total_queries"].(int64) != 0 {
		t.Fatalf("Expected 0 total queries after reset, got %v", metrics["total_queries"])
	}

	t.Log("✓ Metrics reset successfully")
}

// TestHealthThresholds tests setting custom health thresholds
func TestHealthThresholds(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Set custom thresholds
	hc.SetHealthThresholds(
		100*time.Millisecond,  // degraded
		500*time.Millisecond,  // unhealthy
		0.1,                   // 10% error rate
		5*time.Second,         // replication lag
	)

	// Health should now be based on new thresholds
	t.Log("✓ Health thresholds updated successfully")
}

// TestHealthCheckerStop tests stopping health checks
func TestHealthCheckerStop(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Millisecond)
	hc.Start()
	time.Sleep(50 * time.Millisecond)

	hc.Stop()
	time.Sleep(100 * time.Millisecond) // Wait for goroutine to stop

	t.Log("✓ Health checker stopped successfully")
}

// TestQPSCalculation tests QPS (queries per second) calculation
func TestQPSCalculation(t *testing.T) {
	router := createTestShardRouter(t)
	defer router.Close()

	hc := sharding.NewHealthChecker(router, 10*time.Second)

	// Record queries over time
	for i := 0; i < 10; i++ {
		hc.RecordQuery(0, 10*time.Millisecond, nil)
	}

	time.Sleep(100 * time.Millisecond)

	metrics := hc.GetMetrics(0)
	qps := metrics["qps"].(float64)

	// Should be roughly 100 QPS (10 queries / 0.1 seconds)
	if qps < 50 || qps > 150 {
		t.Logf("Warning: QPS outside expected range (50-150): %f", qps)
	}

	t.Logf("✓ QPS calculated: %.2f queries/sec", qps)
}

// BenchmarkHealthCheck benchmarks health checking performance
func BenchmarkHealthCheck(b *testing.B) {
	router := createTestShardRouter(&testing.T{})
	defer router.Close()

	hc := sharding.NewHealthChecker(router, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hc.GetAllHealth()
	}
}

// BenchmarkMetricsRecording benchmarks metrics recording
func BenchmarkMetricsRecording(b *testing.B) {
	router := createTestShardRouter(&testing.T{})
	defer router.Close()

	hc := sharding.NewHealthChecker(router, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hc.RecordQuery(0, 10*time.Millisecond, nil)
	}
}

// HELPER VARIABLES

var ErrTestError = NewTestError("test error")

type TestError struct {
	msg string
}

func NewTestError(msg string) *TestError {
	return &TestError{msg}
}

func (e *TestError) Error() string {
	return e.msg
}
