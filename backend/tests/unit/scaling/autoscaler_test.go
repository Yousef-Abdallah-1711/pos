//go:build experimental
// +build experimental

package scaling_test

import (
	"testing"
	"time"

	"pos-saas/internal/scaling"
)

// TestAutoScalerInitialization tests auto-scaler creation
func TestAutoScalerInitialization(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	if as == nil {
		t.Fatal("Auto-scaler initialization failed")
	}

	t.Log("✓ Auto-scaler initialized successfully")
}

// TestCreateScalingPolicy tests scaling policy creation
func TestCreateScalingPolicy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	policy := as.CreateScalingPolicy(
		"us-east", "standard-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	if policy == nil {
		t.Fatal("Failed to create scaling policy")
	}

	if policy.RegionName != "us-east" {
		t.Fatalf("Expected region 'us-east', got '%s'", policy.RegionName)
	}

	if policy.ScaleUpThreshold != 75 {
		t.Fatalf("Expected scale up threshold 75, got %d", policy.ScaleUpThreshold)
	}

	if policy.MinInstances != 2 {
		t.Fatalf("Expected min instances 2, got %d", policy.MinInstances)
	}

	if policy.MaxInstances != 10 {
		t.Fatalf("Expected max instances 10, got %d", policy.MaxInstances)
	}

	if !policy.IsActive {
		t.Fatal("Policy should be active")
	}

	t.Log("✓ Scaling policy created successfully")
}

// TestGetScalingPolicy tests retrieving a scaling policy
func TestGetScalingPolicy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	createdPolicy := as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	retrievedPolicy := as.GetScalingPolicy("us-east")

	if retrievedPolicy == nil {
		t.Fatal("Failed to retrieve scaling policy")
	}

	if retrievedPolicy.RegionName != createdPolicy.RegionName {
		t.Fatal("Retrieved policy does not match created policy")
	}

	t.Log("✓ Scaling policy retrieved successfully")
}

// TestUpdateScalingPolicy tests updating a scaling policy
func TestUpdateScalingPolicy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	newThreshold := 80
	updatedPolicy := as.UpdateScalingPolicy(
		"us-east",
		&newThreshold,
		nil,
		nil,
		nil,
		nil,
	)

	if updatedPolicy == nil {
		t.Fatal("Failed to update scaling policy")
	}

	if updatedPolicy.ScaleUpThreshold != 80 {
		t.Fatalf("Expected scale up threshold 80, got %d", updatedPolicy.ScaleUpThreshold)
	}

	t.Log("✓ Scaling policy updated successfully")
}

// TestScalingPolicyConstraints tests policy constraint validation
func TestScalingPolicyConstraints(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	policy := as.CreateScalingPolicy(
		"us-east", "constraint-policy",
		75, 25, 2, 1,
		5, 20, 300, 30,
	)

	// Verify constraints
	if policy.MinInstances > policy.MaxInstances {
		t.Fatal("Min instances should not exceed max instances")
	}

	if policy.ScaleDownThreshold >= policy.ScaleUpThreshold {
		t.Fatal("Scale down threshold should be less than scale up threshold")
	}

	t.Log("✓ Scaling policy constraints validated")
}

// TestAutoScalerStart tests starting the auto-scaler
func TestAutoScalerStart(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	as.Start()
	defer as.Stop()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	if !as.IsRunning() {
		t.Fatal("Auto-scaler should be running")
	}

	t.Log("✓ Auto-scaler started successfully")
}

// TestAutoScalerStop tests stopping the auto-scaler
func TestAutoScalerStop(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	as.Start()
	time.Sleep(100 * time.Millisecond)

	as.Stop()
	time.Sleep(100 * time.Millisecond)

	if as.IsRunning() {
		t.Fatal("Auto-scaler should be stopped")
	}

	t.Log("✓ Auto-scaler stopped successfully")
}

// TestGetScalingEvents tests retrieving scaling events
func TestGetScalingEvents(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	events := as.GetScalingEvents(10)

	if events == nil {
		t.Fatal("Failed to get scaling events")
	}

	// Initially there should be no events
	if len(events) > 0 {
		t.Logf("Initial event count: %d", len(events))
	}

	t.Log("✓ Scaling events retrieved successfully")
}

// TestGetScalingMetrics tests retrieving scaling metrics
func TestGetScalingMetrics(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	metrics := as.GetScalingMetrics()

	if metrics == nil {
		t.Fatal("Failed to get scaling metrics")
	}

	t.Logf("✓ Scaling metrics retrieved: %v", len(metrics))
}

// TestGetRegionScalingMetrics tests retrieving region-specific metrics
func TestGetRegionScalingMetrics(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	metrics := as.GetRegionScalingMetrics("us-east")

	if metrics == nil {
		t.Fatal("Failed to get region scaling metrics")
	}

	t.Logf("✓ Region scaling metrics retrieved")
}

// TestMultipleScalingPolicies tests managing multiple policies
func TestMultipleScalingPolicies(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	lb.RegisterRegion("us-east", "North America", 40.7128, -74.0060, 1)
	lb.RegisterRegion("eu-west", "Europe", 51.5074, -0.1278, 2)

	policy1 := as.CreateScalingPolicy(
		"us-east", "east-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	policy2 := as.CreateScalingPolicy(
		"eu-west", "west-policy",
		80, 20, 3, 2,
		3, 15, 300, 30,
	)

	policies := as.GetAllScalingPolicies()

	if len(policies) < 2 {
		t.Fatalf("Expected at least 2 policies, got %d", len(policies))
	}

	if policy1 == nil || policy2 == nil {
		t.Fatal("Failed to create both policies")
	}

	t.Logf("✓ Multiple scaling policies created: %d", len(policies))
}

// TestScalingPolicyCooldown tests cooldown period functionality
func TestScalingPolicyCooldown(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	policy := as.CreateScalingPolicy(
		"us-east", "cooldown-policy",
		75, 25, 2, 1,
		2, 10, 60, 30, // 60 second cooldown
	)

	if policy.CooldownPeriod.Seconds() != 60 {
		t.Fatalf("Expected 60 second cooldown, got %v", policy.CooldownPeriod.Seconds())
	}

	t.Log("✓ Scaling policy cooldown validated")
}

// TestScalingPolicyDisabling tests disabling a scaling policy
func TestScalingPolicyDisabling(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	as := scaling.NewAutoScaler(lb)

	policy := as.CreateScalingPolicy(
		"us-east", "test-policy",
		75, 25, 2, 1,
		2, 10, 300, 30,
	)

	if !policy.IsActive {
		t.Fatal("Policy should be active initially")
	}

	isActive := false
	updatedPolicy := as.UpdateScalingPolicy(
		"us-east",
		nil, nil, nil, nil,
		&isActive,
	)

	if updatedPolicy == nil {
		t.Fatal("Failed to update policy")
	}

	if updatedPolicy.IsActive {
		t.Fatal("Policy should be disabled")
	}

	t.Log("✓ Scaling policy disabling works correctly")
}
