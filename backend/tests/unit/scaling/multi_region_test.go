//go:build experimental
// +build experimental

package scaling_test

import (
	"testing"

	"pos-saas/internal/scaling"
)

// TestMultiRegionDeploymentInitialization tests multi-region deployment creation
func TestMultiRegionDeploymentInitialization(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	if mrd == nil {
		t.Fatal("Multi-region deployment initialization failed")
	}

	t.Log("✓ Multi-region deployment initialized successfully")
}

// TestRegisterDataCenter tests data center registration
func TestRegisterDataCenter(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	dc := mrd.RegisterDataCenter(
		"us-east-1", "US East 1", "Virginia, USA",
		"AWS", "us-east-1",
		38.9072, -77.0369,
		1000, "premium", true,
	)

	if dc == nil {
		t.Fatal("Failed to register data center")
	}

	if dc.ID != "us-east-1" {
		t.Fatalf("Expected data center ID 'us-east-1', got '%s'", dc.ID)
	}

	if dc.Provider != "AWS" {
		t.Fatalf("Expected provider 'AWS', got '%s'", dc.Provider)
	}

	if !dc.IsPrimary {
		t.Fatal("Data center should be marked as primary")
	}

	if !dc.IsActive {
		t.Fatal("Data center should be active")
	}

	if dc.Capacity != 1000 {
		t.Fatalf("Expected capacity 1000, got %d", dc.Capacity)
	}

	t.Log("✓ Data center registered successfully")
}

// TestCreateDeploymentConfig tests deployment configuration creation
func TestCreateDeploymentConfig(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"global-deployment",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1", // data residency
		"ap-south-1", // disaster recovery
		"active-active",
		true, // cost optimization
		100,  // latency target
	)

	if config == nil {
		t.Fatal("Failed to create deployment config")
	}

	if config.Name != "global-deployment" {
		t.Fatalf("Expected config name 'global-deployment', got '%s'", config.Name)
	}

	if len(config.ActiveRegions) != 3 {
		t.Fatalf("Expected 3 active regions, got %d", len(config.ActiveRegions))
	}

	if config.FailoverStrategy != "active-active" {
		t.Fatalf("Expected failover strategy 'active-active', got '%s'", config.FailoverStrategy)
	}

	if config.DataResidencyRegion != "eu-west-1" {
		t.Fatalf("Expected data residency region 'eu-west-1', got '%s'", config.DataResidencyRegion)
	}

	if !config.IsActive {
		t.Fatal("Config should be active")
	}

	t.Log("✓ Deployment config created successfully")
}

// TestUpdateLatencyMatrix tests latency matrix updates
func TestUpdateLatencyMatrix(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	mrd.UpdateLatencyMatrix("us-east-1", "eu-west-1", 85)
	mrd.UpdateLatencyMatrix("eu-west-1", "ap-south-1", 120)
	mrd.UpdateLatencyMatrix("ap-south-1", "us-east-1", 150)

	latency := mrd.GetLatency("us-east-1", "eu-west-1")

	if latency != 85 {
		t.Fatalf("Expected latency 85ms, got %dms", latency)
	}

	t.Log("✓ Latency matrix updated successfully")
}

// TestGetLatencyMatrix tests retrieving the full latency matrix
func TestGetLatencyMatrix(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	mrd.UpdateLatencyMatrix("us-east-1", "eu-west-1", 85)
	mrd.UpdateLatencyMatrix("eu-west-1", "ap-south-1", 120)

	matrix := mrd.GetLatencyMatrix()

	if matrix == nil {
		t.Fatal("Failed to get latency matrix")
	}

	if latency, exists := matrix["us-east-1"]["eu-west-1"]; !exists || latency != 85 {
		t.Fatal("Latency matrix does not contain expected entry")
	}

	t.Log("✓ Latency matrix retrieved successfully")
}

// TestOptimizeTrafficDistribution tests traffic distribution optimization
func TestOptimizeTrafficDistribution(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"test-deployment",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1",
		"ap-south-1",
		"active-active",
		true,
		100,
	)

	distribution := mrd.OptimizeTrafficDistribution("test-deployment")

	if distribution == nil {
		t.Fatal("Failed to optimize traffic distribution")
	}

	// For active-active with 3 regions, each should get ~33.33%
	totalDistribution := 0.0
	for _, percentage := range distribution {
		totalDistribution += percentage
	}

	if totalDistribution < 99.0 || totalDistribution > 101.0 {
		t.Fatalf("Expected total distribution ~100%%, got %f%%", totalDistribution)
	}

	t.Logf("✓ Traffic distribution optimized: %v", distribution)
}

// TestActiveActiveStrategy tests active-active deployment strategy
func TestActiveActiveStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"active-active-test",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1",
		"ap-south-1",
		"active-active",
		true,
		100,
	)

	distribution := mrd.OptimizeTrafficDistribution("active-active-test")

	// Each region should get equal distribution
	expectedPerRegion := 100.0 / 3.0

	for region, percentage := range distribution {
		if percentage < expectedPerRegion-1 || percentage > expectedPerRegion+1 {
			t.Fatalf("Expected %f%% for region %s, got %f%%", expectedPerRegion, region, percentage)
		}
	}

	t.Log("✓ Active-active strategy distributes traffic equally")
}

// TestActivePassiveStrategy tests active-passive deployment strategy
func TestActivePassiveStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"active-passive-test",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1",
		"ap-south-1",
		"active-passive",
		true,
		100,
	)

	distribution := mrd.OptimizeTrafficDistribution("active-passive-test")

	// All traffic should go to primary (first region)
	if distribution["us-east-1"] != 100.0 {
		t.Fatalf("Expected 100%% for primary region, got %f%%", distribution["us-east-1"])
	}

	if distribution["eu-west-1"] != 0 {
		t.Fatalf("Expected 0%% for secondary region, got %f%%", distribution["eu-west-1"])
	}

	t.Log("✓ Active-passive strategy routes all traffic to primary")
}

// TestActiveWarmStrategy tests active-warm deployment strategy
func TestActiveWarmStrategy(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"active-warm-test",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1",
		"ap-south-1",
		"active-warm",
		true,
		100,
	)

	distribution := mrd.OptimizeTrafficDistribution("active-warm-test")

	// Primary should get 80%, secondary should get 20%
	if distribution["us-east-1"] != 80.0 {
		t.Fatalf("Expected 80%% for primary region, got %f%%", distribution["us-east-1"])
	}

	if distribution["eu-west-1"] != 20.0 {
		t.Fatalf("Expected 20%% for secondary region, got %f%%", distribution["eu-west-1"])
	}

	t.Log("✓ Active-warm strategy distributes 80-20 traffic")
}

// TestSelectRegionForRequest tests region selection for client requests
func TestSelectRegionForRequest(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config := mrd.CreateDeploymentConfig(
		"region-selection-test",
		[]string{"us-east-1", "eu-west-1"},
		"eu-west-1",
		"us-east-1",
		"active-active",
		true,
		100,
	)

	mrd.UpdateLatencyMatrix("us-east-1", "eu-west-1", 85)
	mrd.UpdateLatencyMatrix("eu-west-1", "us-east-1", 85)

	// Client in us-east-1 should be routed to us-east-1
	selectedRegion := mrd.SelectRegionForRequest("us-east-1", "region-selection-test")

	if selectedRegion != "us-east-1" {
		t.Fatalf("Expected client in us-east-1 to be routed to us-east-1, got %s", selectedRegion)
	}

	t.Log("✓ Region selection for client requests works correctly")
}

// TestGetMultiRegionStatus tests retrieving multi-region status
func TestGetMultiRegionStatus(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	mrd.RegisterDataCenter(
		"us-east-1", "US East 1", "Virginia, USA",
		"AWS", "us-east-1",
		38.9072, -77.0369,
		1000, "premium", true,
	)

	mrd.RegisterDataCenter(
		"eu-west-1", "EU West 1", "Ireland",
		"AWS", "eu-west-1",
		53.4129, -8.2439,
		800, "standard", false,
	)

	config := mrd.CreateDeploymentConfig(
		"status-test",
		[]string{"us-east-1", "eu-west-1"},
		"eu-west-1",
		"us-east-1",
		"active-active",
		true,
		100,
	)

	status := mrd.GetMultiRegionStatus("status-test")

	if status == nil {
		t.Fatal("Failed to get multi-region status")
	}

	if status["config_name"] != "status-test" {
		t.Fatalf("Expected config name 'status-test', got %v", status["config_name"])
	}

	t.Logf("✓ Multi-region status retrieved: %v", status)
}

// TestGetDataCenterStatus tests retrieving data center status
func TestGetDataCenterStatus(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	dc := mrd.RegisterDataCenter(
		"us-east-1", "US East 1", "Virginia, USA",
		"AWS", "us-east-1",
		38.9072, -77.0369,
		1000, "premium", true,
	)

	status := mrd.GetDataCenterStatus("us-east-1")

	if status == nil {
		t.Fatal("Failed to get data center status")
	}

	if status["id"] != dc.ID {
		t.Fatalf("Expected data center ID 'us-east-1', got %v", status["id"])
	}

	if status["is_primary"] != true {
		t.Fatalf("Expected primary data center, got %v", status["is_primary"])
	}

	t.Logf("✓ Data center status retrieved")
}

// TestGetDeploymentComparison tests deployment strategy comparison
func TestGetDeploymentComparison(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	comparison := mrd.GetDeploymentComparison()

	if comparison == nil {
		t.Fatal("Failed to get deployment comparison")
	}

	strategies, exists := comparison["strategies"]
	if !exists {
		t.Fatal("Strategies not found in comparison")
	}

	strategyMap, ok := strategies.(map[string]map[string]string)
	if !ok {
		t.Fatal("Invalid strategies structure")
	}

	if _, exists := strategyMap["active-active"]; !exists {
		t.Fatal("Active-active strategy not found in comparison")
	}

	if _, exists := strategyMap["active-passive"]; !exists {
		t.Fatal("Active-passive strategy not found in comparison")
	}

	if _, exists := strategyMap["active-warm"]; !exists {
		t.Fatal("Active-warm strategy not found in comparison")
	}

	t.Log("✓ Deployment comparison retrieved successfully")
}

// TestCalculateGeometricDistance tests geographic distance calculation
func TestCalculateGeometricDistance(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	// US East (Washington DC): 38.9072, -77.0369
	// EU West (Dublin): 53.4129, -8.2439
	distance := mrd.CalculateGeometricDistance(38.9072, -77.0369, 53.4129, -8.2439)

	// Should be approximately 5,200 km
	if distance < 5000 || distance > 5500 {
		t.Fatalf("Expected distance ~5200km, got %f km", distance)
	}

	t.Logf("✓ Geographic distance calculated: %f km", distance)
}

// TestMultipleDataCenters tests managing multiple data centers
func TestMultipleDataCenters(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	mrd.RegisterDataCenter(
		"us-east-1", "US East 1", "Virginia, USA",
		"AWS", "us-east-1",
		38.9072, -77.0369,
		1000, "premium", true,
	)

	mrd.RegisterDataCenter(
		"eu-west-1", "EU West 1", "Ireland",
		"AWS", "eu-west-1",
		53.4129, -8.2439,
		800, "standard", false,
	)

	mrd.RegisterDataCenter(
		"ap-south-1", "AP South 1", "Mumbai, India",
		"AWS", "ap-south-1",
		19.0760, 72.8777,
		600, "standard", false,
	)

	dataCenters := mrd.GetAllDataCenters()

	if len(dataCenters) != 3 {
		t.Fatalf("Expected 3 data centers, got %d", len(dataCenters))
	}

	t.Logf("✓ Multiple data centers registered: %d", len(dataCenters))
}

// TestMultipleDeploymentConfigs tests managing multiple deployment configurations
func TestMultipleDeploymentConfigs(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	config1 := mrd.CreateDeploymentConfig(
		"global-config",
		[]string{"us-east-1", "eu-west-1", "ap-south-1"},
		"eu-west-1",
		"ap-south-1",
		"active-active",
		true,
		100,
	)

	config2 := mrd.CreateDeploymentConfig(
		"regional-config",
		[]string{"us-east-1", "us-west-2"},
		"us-east-1",
		"us-west-2",
		"active-passive",
		false,
		150,
	)

	if config1 == nil || config2 == nil {
		t.Fatal("Failed to create both configs")
	}

	configs := mrd.GetAllDeploymentConfigs()

	if len(configs) < 2 {
		t.Fatalf("Expected at least 2 configs, got %d", len(configs))
	}

	t.Logf("✓ Multiple deployment configs created: %d", len(configs))
}

// TestDataResidencyCompliance tests data residency region enforcement
func TestDataResidencyCompliance(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	// For GDPR compliance, data must reside in eu-west-1
	config := mrd.CreateDeploymentConfig(
		"gdpr-compliant",
		[]string{"us-east-1", "eu-west-1", "eu-central-1"},
		"eu-west-1", // Data residency requirement
		"eu-central-1",
		"active-active",
		false,
		100,
	)

	if config.DataResidencyRegion != "eu-west-1" {
		t.Fatal("Data residency region not properly set")
	}

	t.Log("✓ Data residency compliance validated")
}

// TestTierConfiguration tests data center tier configuration
func TestTierConfiguration(t *testing.T) {
	lb := scaling.NewLoadBalancer("round-robin")
	mrd := scaling.NewMultiRegionDeployment(lb)

	premium := mrd.RegisterDataCenter(
		"premium-dc", "Premium DC", "Singapore",
		"GCP", "asia-southeast1",
		1.3521, 103.8198,
		2000, "premium", true,
	)

	standard := mrd.RegisterDataCenter(
		"standard-dc", "Standard DC", "Tokyo",
		"AWS", "ap-northeast-1",
		35.6762, 139.6503,
		1000, "standard", false,
	)

	budget := mrd.RegisterDataCenter(
		"budget-dc", "Budget DC", "Mumbai",
		"Azure", "southindia",
		19.0760, 72.8777,
		500, "budget", false,
	)

	if premium.Tier != "premium" || standard.Tier != "standard" || budget.Tier != "budget" {
		t.Fatal("Data center tiers not properly configured")
	}

	t.Log("✓ Data center tier configuration validated")
}
