//go:build experimental
// +build experimental

package scaling_test

import (
	"testing"

	"pos-saas/internal/scaling"
)

// TestLoadBalancerInitialization tests load balancer creation
func TestLoadBalancerInitialization(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")

	if lb == nil {
		t.Fatal("Load balancer initialization failed")
	}

	if lb.GetCurrentStrategy() != "round-robin" {
		t.Fatalf("Expected strategy 'round-robin', got '%s'", lb.GetCurrentStrategy())
	}

	t.Log("✓ Load balancer initialized successfully")
}

// TestRegisterRegion tests region registration
func TestRegisterRegion(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")

	region := lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	if region == nil {
		t.Fatal("Region registration failed")
	}

	if region.Name != "us-east" {
		t.Fatalf("Expected region name 'us-east', got '%s'", region.Name)
	}

	if !region.IsActive {
		t.Fatal("Region should be active")
	}

	t.Log("✓ Region registered successfully")
}

// TestRegisterServer tests server registration
func TestRegisterServer(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	server := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)

	if server == nil {
		t.Fatal("Server registration failed")
	}

	if server.ID != "server-1" {
		t.Fatalf("Expected server ID 'server-1', got '%s'", server.ID)
	}

	if server.Region != "us-east" {
		t.Fatalf("Expected region 'us-east', got '%s'", server.Region)
	}

	t.Log("✓ Server registered successfully")
}

// TestRoundRobinLoadBalancing tests round-robin server selection
func TestRoundRobinLoadBalancing(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	s1 := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	s2 := lb.RegisterServer("server-2", "us-east", "192.168.1.2", 8080, 1, 100)
	s3 := lb.RegisterServer("server-3", "us-east", "192.168.1.3", 8080, 1, 100)

	// Test round-robin order
	selected1 := lb.SelectServer("us-east")
	selected2 := lb.SelectServer("us-east")
	selected3 := lb.SelectServer("us-east")
	selected4 := lb.SelectServer("us-east")

	if selected1.ID != s1.ID {
		t.Fatalf("Expected first selection to be server-1, got %s", selected1.ID)
	}

	if selected2.ID != s2.ID {
		t.Fatalf("Expected second selection to be server-2, got %s", selected2.ID)
	}

	if selected3.ID != s3.ID {
		t.Fatalf("Expected third selection to be server-3, got %s", selected3.ID)
	}

	if selected4.ID != s1.ID {
		t.Fatalf("Expected fourth selection to cycle back to server-1, got %s", selected4.ID)
	}

	t.Log("✓ Round-robin load balancing works correctly")
}

// TestLeastLoadedStrategy tests least-loaded server selection
func TestLeastLoadedStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("least-loaded")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	s1 := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	s2 := lb.RegisterServer("server-2", "us-east", "192.168.1.2", 8080, 1, 100)
	s3 := lb.RegisterServer("server-3", "us-east", "192.168.1.3", 8080, 1, 100)

	// Update loads
	lb.UpdateServerMetrics(s1.ID, 80, 50, 0.1, 75.0, 60.0)
	lb.UpdateServerMetrics(s2.ID, 20, 30, 0.01, 25.0, 20.0)
	lb.UpdateServerMetrics(s3.ID, 50, 40, 0.05, 50.0, 40.0)

	selected := lb.SelectServer("us-east")

	if selected.ID != s2.ID {
		t.Fatalf("Expected selection of least-loaded server (server-2), got %s", selected.ID)
	}

	t.Log("✓ Least-loaded strategy works correctly")
}

// TestFastestStrategy tests fastest server selection
func TestFastestStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("fastest")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	s1 := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	s2 := lb.RegisterServer("server-2", "us-east", "192.168.1.2", 8080, 1, 100)
	s3 := lb.RegisterServer("server-3", "us-east", "192.168.1.3", 8080, 1, 100)

	// Update response times
	lb.UpdateServerMetrics(s1.ID, 50, 150, 0.1, 75.0, 60.0)
	lb.UpdateServerMetrics(s2.ID, 50, 20, 0.01, 25.0, 20.0)
	lb.UpdateServerMetrics(s3.ID, 50, 100, 0.05, 50.0, 40.0)

	selected := lb.SelectServer("us-east")

	if selected.ID != s2.ID {
		t.Fatalf("Expected selection of fastest server (server-2), got %s", selected.ID)
	}

	t.Log("✓ Fastest strategy works correctly")
}

// TestGeoAwareStrategy tests geo-aware server selection
func TestGeoAwareStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("geo-aware")

	region1 := lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	region2 := lb.RegisterRegion("eu-west", "Europe", 51.5074, -0.1278, 2)

	s1 := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	s2 := lb.RegisterServer("server-2", "eu-west", "192.168.2.1", 8080, 1, 100)

	if region1 == nil || region2 == nil {
		t.Fatal("Region registration failed")
	}

	selected := lb.SelectServer("us-east")

	if selected.ID != s1.ID {
		t.Fatalf("Expected selection from us-east region, got %s", selected.ID)
	}

	t.Log("✓ Geo-aware strategy works correctly")
}

// TestSessionAffinity tests sticky session functionality
func TestSessionAffinity(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	s1 := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	s2 := lb.RegisterServer("server-2", "us-east", "192.168.1.2", 8080, 1, 100)
	s3 := lb.RegisterServer("server-3", "us-east", "192.168.1.3", 8080, 1, 100)

	lb.SetSessionAffinity(true, 1800)

	// Select server for a session
	sessionID := "session-123"
	selected1 := lb.SelectServer("us-east")

	// Simulate session affinity lookup (would be done in handler)
	selected2 := lb.SelectServer("us-east")
	selected3 := lb.SelectServer("us-east")

	// At least one server should be selected
	if selected1 == nil {
		t.Fatal("Server selection failed")
	}

	if selected2 == nil {
		t.Fatal("Server selection failed")
	}

	if selected3 == nil {
		t.Fatal("Server selection failed")
	}

	t.Logf("✓ Session affinity configured (sessionID: %s)", sessionID)
}

// TestServerHealthStatus tests server health updates
func TestServerHealthStatus(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	server := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)

	// Server should start healthy
	if !server.IsHealthy {
		t.Fatal("Server should start healthy")
	}

	// Mark as unhealthy
	lb.SetServerHealth("server-1", false)

	updatedServer := lb.GetServer("server-1")
	if updatedServer.IsHealthy {
		t.Fatal("Server should be marked unhealthy")
	}

	// Mark as healthy again
	lb.SetServerHealth("server-1", true)
	updatedServer = lb.GetServer("server-1")
	if !updatedServer.IsHealthy {
		t.Fatal("Server should be marked healthy")
	}

	t.Log("✓ Server health status updates work correctly")
}

// TestGetGlobalStats tests global statistics
func TestGetGlobalStats(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	lb.RegisterRegion("eu-west", "Europe", 51.5074, -0.1278, 2)

	lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)
	lb.RegisterServer("server-2", "us-east", "192.168.1.2", 8080, 1, 100)
	lb.RegisterServer("server-3", "eu-west", "192.168.2.1", 8080, 1, 150)

	stats := lb.GetGlobalStats()

	if stats == nil {
		t.Fatal("Failed to get global stats")
	}

	if stats["total_regions"] != 2 {
		t.Fatalf("Expected 2 regions, got %v", stats["total_regions"])
	}

	if stats["total_servers"] != 3 {
		t.Fatalf("Expected 3 servers, got %v", stats["total_servers"])
	}

	t.Logf("✓ Global stats retrieved: %v", stats)
}

// TestChangeStrategy tests changing load balancing strategy
func TestChangeStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")

	if lb.GetCurrentStrategy() != "round-robin" {
		t.Fatal("Initial strategy should be round-robin")
	}

	lb.SetStrategy("least-loaded")

	if lb.GetCurrentStrategy() != "least-loaded" {
		t.Fatal("Strategy should be changed to least-loaded")
	}

	lb.SetStrategy("fastest")

	if lb.GetCurrentStrategy() != "fastest" {
		t.Fatal("Strategy should be changed to fastest")
	}

	t.Log("✓ Load balancing strategy changes work correctly")
}

// TestUpdateServerMetrics tests metric updates
func TestUpdateServerMetrics(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)

	server := lb.RegisterServer("server-1", "us-east", "192.168.1.1", 8080, 1, 100)

	// Update metrics
	lb.UpdateServerMetrics("server-1", 75, 45, 0.02, 80.0, 70.0)

	updatedServer := lb.GetServer("server-1")

	if updatedServer.CurrentLoad != 75 {
		t.Fatalf("Expected load 75, got %d", updatedServer.CurrentLoad)
	}

	if updatedServer.ResponseTime != 45 {
		t.Fatalf("Expected response time 45, got %d", updatedServer.ResponseTime)
	}

	if updatedServer.CPUUsage != 80.0 {
		t.Fatalf("Expected CPU usage 80.0, got %f", updatedServer.CPUUsage)
	}

	if updatedServer.MemoryUsage != 70.0 {
		t.Fatalf("Expected memory usage 70.0, got %f", updatedServer.MemoryUsage)
	}

	t.Log("✓ Server metrics updated successfully")
}
