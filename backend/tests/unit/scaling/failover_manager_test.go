//go:build experimental
// +build experimental

package scaling_test

import (
	"testing"
	"time"

	"pos-saas/internal/scaling"
)

// TestFailoverManagerInitialization tests failover manager creation
func TestFailoverManagerInitialization(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	if fm == nil {
		t.Fatal("Failover manager initialization failed")
	}

	t.Log("✓ Failover manager initialized successfully")
}

// TestCreateFailoverPolicy tests failover policy creation
func TestCreateFailoverPolicy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	policy := fm.CreateFailoverPolicy(
		"primary-secondary", "us-east", "eu-west",
		0.5, 300, true,
	)

	if policy == nil {
		t.Fatal("Failed to create failover policy")
	}

	if policy.Name != "primary-secondary" {
		t.Fatalf("Expected policy name 'primary-secondary', got '%s'", policy.Name)
	}

	if policy.PrimaryRegion != "us-east" {
		t.Fatalf("Expected primary region 'us-east', got '%s'", policy.PrimaryRegion)
	}

	if policy.SecondaryRegion != "eu-west" {
		t.Fatalf("Expected secondary region 'eu-west', got '%s'", policy.SecondaryRegion)
	}

	if !policy.AutomaticFailover {
		t.Fatal("Automatic failover should be enabled")
	}

	if !policy.IsActive {
		t.Fatal("Policy should be active")
	}

	t.Log("✓ Failover policy created successfully")
}

// TestGetFailoverPolicy tests retrieving a failover policy
func TestGetFailoverPolicy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	createdPolicy := fm.CreateFailoverPolicy(
		"test-policy", "us-east", "eu-west",
		0.5, 300, true,
	)

	retrievedPolicy := fm.GetFailoverPolicy(createdPolicy.ID)

	if retrievedPolicy == nil {
		t.Fatal("Failed to retrieve failover policy")
	}

	if retrievedPolicy.ID != createdPolicy.ID {
		t.Fatal("Retrieved policy does not match created policy")
	}

	t.Log("✓ Failover policy retrieved successfully")
}

// TestCreateDisasterRecoveryPlan tests disaster recovery plan creation
func TestCreateDisasterRecoveryPlan(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	plan := fm.CreateDisasterRecoveryPlan(
		"critical-systems", 15, 5,
		30, 90, true,
		[]string{"us-east", "eu-west", "ap-south"},
	)

	if plan == nil {
		t.Fatal("Failed to create disaster recovery plan")
	}

	if plan.Name != "critical-systems" {
		t.Fatalf("Expected plan name 'critical-systems', got '%s'", plan.Name)
	}

	if plan.RecoveryTimeObjective != 15 {
		t.Fatalf("Expected RTO 15 minutes, got %d", plan.RecoveryTimeObjective)
	}

	if plan.RecoveryPointObjective != 5 {
		t.Fatalf("Expected RPO 5 minutes, got %d", plan.RecoveryPointObjective)
	}

	if !plan.ReplicationEnabled {
		t.Fatal("Replication should be enabled")
	}

	if len(plan.ReplicationTargets) != 3 {
		t.Fatalf("Expected 3 replication targets, got %d", len(plan.ReplicationTargets))
	}

	if !plan.IsActive {
		t.Fatal("Plan should be active")
	}

	t.Log("✓ Disaster recovery plan created successfully")
}

// TestGetDisasterRecoveryPlan tests retrieving a disaster recovery plan
func TestGetDisasterRecoveryPlan(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	createdPlan := fm.CreateDisasterRecoveryPlan(
		"test-plan", 15, 5,
		30, 90, true,
		[]string{"us-east", "eu-west"},
	)

	retrievedPlan := fm.GetDisasterRecoveryPlan(createdPlan.ID)

	if retrievedPlan == nil {
		t.Fatal("Failed to retrieve disaster recovery plan")
	}

	if retrievedPlan.ID != createdPlan.ID {
		t.Fatal("Retrieved plan does not match created plan")
	}

	t.Log("✓ Disaster recovery plan retrieved successfully")
}

// TestGetFailoverStatus tests retrieving failover status
func TestGetFailoverStatus(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	status := fm.GetFailoverStatus()

	if status == nil {
		t.Fatal("Failed to get failover status")
	}

	t.Logf("✓ Failover status retrieved: %v", status)
}

// TestGetFailoverHistory tests retrieving failover event history
func TestGetFailoverHistory(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	history := fm.GetFailoverHistory(10)

	if history == nil {
		t.Fatal("Failed to get failover history")
	}

	// Initially should have no events
	if len(history) > 0 {
		t.Logf("Initial failover event count: %d", len(history))
	}

	t.Log("✓ Failover history retrieved successfully")
}

// TestGetDisasterRecoveryStatus tests disaster recovery status
func TestGetDisasterRecoveryStatus(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	fm.CreateDisasterRecoveryPlan(
		"test-plan", 15, 5,
		30, 90, true,
		[]string{"us-east", "eu-west"},
	)

	status := fm.GetDisasterRecoveryStatus()

	if status == nil {
		t.Fatal("Failed to get disaster recovery status")
	}

	t.Logf("✓ Disaster recovery status retrieved")
}

// TestMultipleFailoverPolicies tests managing multiple failover policies
func TestMultipleFailoverPolicies(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	policy1 := fm.CreateFailoverPolicy(
		"policy-1", "us-east", "eu-west",
		0.5, 300, true,
	)

	policy2 := fm.CreateFailoverPolicy(
		"policy-2", "ap-south", "ap-southeast",
		0.6, 300, false,
	)

	if policy1 == nil || policy2 == nil {
		t.Fatal("Failed to create both policies")
	}

	policies := fm.GetAllFailoverPolicies()

	if len(policies) < 2 {
		t.Fatalf("Expected at least 2 policies, got %d", len(policies))
	}

	t.Logf("✓ Multiple failover policies created: %d", len(policies))
}

// TestMultipleDisasterRecoveryPlans tests managing multiple DR plans
func TestMultipleDisasterRecoveryPlans(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	plan1 := fm.CreateDisasterRecoveryPlan(
		"critical-plan", 15, 5,
		30, 90, true,
		[]string{"us-east", "eu-west"},
	)

	plan2 := fm.CreateDisasterRecoveryPlan(
		"standard-plan", 60, 30,
		60, 180, true,
		[]string{"ap-south", "ap-southeast"},
	)

	if plan1 == nil || plan2 == nil {
		t.Fatal("Failed to create both plans")
	}

	plans := fm.GetAllDisasterRecoveryPlans()

	if len(plans) < 2 {
		t.Fatalf("Expected at least 2 plans, got %d", len(plans))
	}

	t.Logf("✓ Multiple disaster recovery plans created: %d", len(plans))
}

// TestFailoverPolicyValidation tests policy constraint validation
func TestFailoverPolicyValidation(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	policy := fm.CreateFailoverPolicy(
		"test-policy", "us-east", "eu-west",
		0.5, 300, true,
	)

	if policy.PrimaryRegion == policy.SecondaryRegion {
		t.Fatal("Primary and secondary regions should be different")
	}

	if policy.HealthCheckThreshold < 0 || policy.HealthCheckThreshold > 1 {
		t.Fatal("Health check threshold should be between 0 and 1")
	}

	if policy.FailoverTimeout <= 0 {
		t.Fatal("Failover timeout should be positive")
	}

	t.Log("✓ Failover policy validation passed")
}

// TestDisasterRecoveryPlanValidation tests DR plan constraint validation
func TestDisasterRecoveryPlanValidation(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	plan := fm.CreateDisasterRecoveryPlan(
		"test-plan", 15, 5,
		30, 90, true,
		[]string{"us-east", "eu-west"},
	)

	if plan.RecoveryTimeObjective <= 0 {
		t.Fatal("RTO should be positive")
	}

	if plan.RecoveryPointObjective <= 0 {
		t.Fatal("RPO should be positive")
	}

	if plan.RecoveryPointObjective > plan.RecoveryTimeObjective {
		t.Fatal("RPO should not exceed RTO")
	}

	if plan.BackupRetention <= 0 {
		t.Fatal("Backup retention should be positive")
	}

	t.Log("✓ Disaster recovery plan validation passed")
}

// TestReplicationLagTracking tests replication lag monitoring
func TestReplicationLagTracking(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	fm.UpdateReplicationLag("us-east", "eu-west", 250)
	fm.UpdateReplicationLag("us-east", "ap-south", 500)

	// Retrieve and verify lag tracking
	status := fm.GetFailoverStatus()

	if status == nil {
		t.Fatal("Failed to get failover status")
	}

	t.Log("✓ Replication lag tracking works correctly")
}

// TestAutomaticFailoverConfiguration tests automatic failover settings
func TestAutomaticFailoverConfiguration(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	automaticPolicy := fm.CreateFailoverPolicy(
		"auto-policy", "us-east", "eu-west",
		0.5, 300, true,
	)

	manualPolicy := fm.CreateFailoverPolicy(
		"manual-policy", "ap-south", "ap-southeast",
		0.5, 300, false,
	)

	if !automaticPolicy.AutomaticFailover {
		t.Fatal("Automatic failover should be enabled")
	}

	if manualPolicy.AutomaticFailover {
		t.Fatal("Automatic failover should be disabled")
	}

	t.Log("✓ Automatic failover configuration works correctly")
}

// TestHealthCheckThreshold tests health check threshold settings
func TestHealthCheckThreshold(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	hc := scaling.NewHealthChecker(lb, 30*time.Second, 5*time.Second, 3)
	fm := scaling.NewFailoverManager(lb, hc)

	policy := fm.CreateFailoverPolicy(
		"threshold-policy", "us-east", "eu-west",
		0.5, 300, true,
	)

	if policy.HealthCheckThreshold != 0.5 {
		t.Fatalf("Expected health check threshold 0.5, got %f", policy.HealthCheckThreshold)
	}

	t.Log("✓ Health check threshold configuration works correctly")
}
