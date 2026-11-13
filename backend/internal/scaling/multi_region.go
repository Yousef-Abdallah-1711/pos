package scaling

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// DataCenter represents a physical data center
type DataCenter struct {
	ID           string
	Name         string
	Location     string // City, Country
	Provider     string // AWS, GCP, Azure, etc.
	Region       string // us-east-1, eu-west-1, etc.
	Latitude     float64
	Longitude    float64
	Capacity     int // Max servers
	CurrentCount int // Current server count
	IsActive     bool
	IsPrimary    bool   // Primary for data residency
	Tier         string // "premium", "standard", "budget"
	CreatedAt    time.Time
}

// MultiRegionConfig defines multi-region deployment strategy
type MultiRegionConfig struct {
	Name                   string
	ActiveRegions          []string
	DataResidencyRegion    string // Required region for data residency
	DisasterRecoveryRegion string
	FailoverStrategy       string // "active-active", "active-passive", "active-warm"
	CostOptimization       bool
	LatencyTarget          int    // Target P99 latency in ms
	DataReplicationMode    string // "sync", "async"
	IsActive               bool
	CreatedAt              time.Time
}

// MultiRegionDeployment manages multi-region deployment
type MultiRegionDeployment struct {
	loadBalancer        *LoadBalancer
	dataCenters         map[string]*DataCenter
	deploymentConfigs   map[string]*MultiRegionConfig
	regionLatencies     map[string]map[string]int     // region -> region -> latency_ms
	regionDistances     map[string]map[string]float64 // region -> region -> km
	trafficDistribution map[string]float64            // region -> traffic_percentage
	mu                  sync.RWMutex
}

// NewMultiRegionDeployment creates a new multi-region deployment manager
func NewMultiRegionDeployment(loadBalancer *LoadBalancer) *MultiRegionDeployment {
	return &MultiRegionDeployment{
		loadBalancer:        loadBalancer,
		dataCenters:         make(map[string]*DataCenter),
		deploymentConfigs:   make(map[string]*MultiRegionConfig),
		regionLatencies:     make(map[string]map[string]int),
		regionDistances:     make(map[string]map[string]float64),
		trafficDistribution: make(map[string]float64),
	}
}

// RegisterDataCenter registers a new data center
func (mrd *MultiRegionDeployment) RegisterDataCenter(
	dcID string,
	name string,
	location string,
	provider string,
	region string,
	latitude float64,
	longitude float64,
	capacity int,
	tier string,
	isPrimary bool,
) *DataCenter {
	log.Printf("[MULTI-REGION] Registering data center: %s (%s) in %s", dcID, provider, location)

	dc := &DataCenter{
		ID:        dcID,
		Name:      name,
		Location:  location,
		Provider:  provider,
		Region:    region,
		Latitude:  latitude,
		Longitude: longitude,
		Capacity:  capacity,
		IsActive:  true,
		IsPrimary: isPrimary,
		Tier:      tier,
		CreatedAt: time.Now(),
	}

	mrd.mu.Lock()
	mrd.dataCenters[dcID] = dc
	mrd.trafficDistribution[region] = 0
	mrd.mu.Unlock()

	log.Printf("[MULTI-REGION] Data center registered: capacity=%d, tier=%s", capacity, tier)
	return dc
}

// CreateDeploymentConfig creates a multi-region deployment configuration
func (mrd *MultiRegionDeployment) CreateDeploymentConfig(
	name string,
	activeRegions []string,
	dataResidencyRegion string,
	disasterRecoveryRegion string,
	failoverStrategy string,
	costOptimization bool,
	latencyTarget int,
) *MultiRegionConfig {
	log.Printf("[MULTI-REGION] Creating deployment config: %s (regions: %v)", name, activeRegions)

	config := &MultiRegionConfig{
		Name:                   name,
		ActiveRegions:          activeRegions,
		DataResidencyRegion:    dataResidencyRegion,
		DisasterRecoveryRegion: disasterRecoveryRegion,
		FailoverStrategy:       failoverStrategy,
		CostOptimization:       costOptimization,
		LatencyTarget:          latencyTarget,
		DataReplicationMode:    "async",
		IsActive:               true,
		CreatedAt:              time.Now(),
	}

	mrd.mu.Lock()
	mrd.deploymentConfigs[name] = config
	mrd.mu.Unlock()

	log.Printf("[MULTI-REGION] Config created: strategy=%s, target_latency=%dms", failoverStrategy, latencyTarget)
	return config
}

// UpdateLatencyMatrix updates latency measurements between regions
func (mrd *MultiRegionDeployment) UpdateLatencyMatrix(
	fromRegion string,
	toRegion string,
	latencyMs int,
) {
	mrd.mu.Lock()
	defer mrd.mu.Unlock()

	if _, exists := mrd.regionLatencies[fromRegion]; !exists {
		mrd.regionLatencies[fromRegion] = make(map[string]int)
	}

	mrd.regionLatencies[fromRegion][toRegion] = latencyMs

	log.Printf("[MULTI-REGION] Latency updated: %s -> %s = %dms", fromRegion, toRegion, latencyMs)
}

// GetLatency returns latency between two regions
func (mrd *MultiRegionDeployment) GetLatency(fromRegion string, toRegion string) int {
	mrd.mu.RLock()
	defer mrd.mu.RUnlock()

	if latencies, exists := mrd.regionLatencies[fromRegion]; exists {
		if latency, exists := latencies[toRegion]; exists {
			return latency
		}
	}

	return 0
}

// OptimizeTrafficDistribution optimizes traffic distribution across regions
func (mrd *MultiRegionDeployment) OptimizeTrafficDistribution(configName string) map[string]float64 {
	mrd.mu.RLock()
	config := mrd.deploymentConfigs[configName]
	mrd.mu.RUnlock()

	if config == nil {
		return nil
	}

	distribution := make(map[string]float64)

	switch config.FailoverStrategy {
	case "active-active":
		// Equal distribution
		per := 100.0 / float64(len(config.ActiveRegions))
		for _, region := range config.ActiveRegions {
			distribution[region] = per
		}

	case "active-passive":
		// All traffic to primary
		distribution[config.ActiveRegions[0]] = 100.0
		for i := 1; i < len(config.ActiveRegions); i++ {
			distribution[config.ActiveRegions[i]] = 0
		}

	case "active-warm":
		// 80% to primary, 20% to secondary
		if len(config.ActiveRegions) >= 2 {
			distribution[config.ActiveRegions[0]] = 80.0
			distribution[config.ActiveRegions[1]] = 20.0
			// Remaining regions get no traffic
			for i := 2; i < len(config.ActiveRegions); i++ {
				distribution[config.ActiveRegions[i]] = 0
			}
		}

	default:
		// Default to equal distribution
		per := 100.0 / float64(len(config.ActiveRegions))
		for _, region := range config.ActiveRegions {
			distribution[region] = per
		}
	}

	mrd.mu.Lock()
	mrd.trafficDistribution = distribution
	mrd.mu.Unlock()

	log.Printf("[MULTI-REGION] Traffic distribution optimized: %v", distribution)
	return distribution
}

// SelectRegionForRequest selects best region for a client request
func (mrd *MultiRegionDeployment) SelectRegionForRequest(
	clientRegion string,
	configName string,
) string {
	mrd.mu.RLock()
	config := mrd.deploymentConfigs[configName]
	mrd.mu.RUnlock()

	if config == nil {
		return ""
	}

	// If client is in active region, route to that region
	for _, activeRegion := range config.ActiveRegions {
		if activeRegion == clientRegion {
			return activeRegion
		}
	}

	// Otherwise, find region with lowest latency
	minLatency := 99999
	bestRegion := ""

	for _, activeRegion := range config.ActiveRegions {
		latency := mrd.GetLatency(clientRegion, activeRegion)
		if latency > 0 && latency < minLatency {
			minLatency = latency
			bestRegion = activeRegion
		}
	}

	if bestRegion == "" && len(config.ActiveRegions) > 0 {
		bestRegion = config.ActiveRegions[0]
	}

	return bestRegion
}

// GetMultiRegionStatus returns overall multi-region deployment status
func (mrd *MultiRegionDeployment) GetMultiRegionStatus(configName string) map[string]interface{} {
	mrd.mu.RLock()
	config := mrd.deploymentConfigs[configName]
	distribution := make(map[string]float64)
	for k, v := range mrd.trafficDistribution {
		distribution[k] = v
	}
	mrd.mu.RUnlock()

	if config == nil {
		return nil
	}

	dcCount := 0
	totalCapacity := 0
	totalServers := 0

	mrd.mu.RLock()
	for _, dc := range mrd.dataCenters {
		if dc.IsActive {
			dcCount++
			totalCapacity += dc.Capacity
			totalServers += dc.CurrentCount
		}
	}
	mrd.mu.RUnlock()

	return map[string]interface{}{
		"config_name":              configName,
		"strategy":                 config.FailoverStrategy,
		"active_regions":           config.ActiveRegions,
		"data_residency_region":    config.DataResidencyRegion,
		"disaster_recovery_region": config.DisasterRecoveryRegion,
		"data_replication_mode":    config.DataReplicationMode,
		"total_data_centers":       dcCount,
		"total_capacity":           totalCapacity,
		"total_servers":            totalServers,
		"traffic_distribution":     distribution,
		"latency_target_ms":        config.LatencyTarget,
		"cost_optimization":        config.CostOptimization,
	}
}

// GetDataCenterStatus returns status of a data center
func (mrd *MultiRegionDeployment) GetDataCenterStatus(dcID string) map[string]interface{} {
	mrd.mu.RLock()
	defer mrd.mu.RUnlock()

	dc, exists := mrd.dataCenters[dcID]
	if !exists {
		return nil
	}

	utilizationPercent := 0
	if dc.Capacity > 0 {
		utilizationPercent = (dc.CurrentCount * 100) / dc.Capacity
	}

	return map[string]interface{}{
		"id":              dc.ID,
		"name":            dc.Name,
		"location":        dc.Location,
		"provider":        dc.Provider,
		"region":          dc.Region,
		"tier":            dc.Tier,
		"is_primary":      dc.IsPrimary,
		"is_active":       dc.IsActive,
		"capacity":        dc.Capacity,
		"current_servers": dc.CurrentCount,
		"utilization":     fmt.Sprintf("%d%%", utilizationPercent),
		"created_at":      dc.CreatedAt,
	}
}

// GetDeploymentComparison compares multiple deployment strategies
func (mrd *MultiRegionDeployment) GetDeploymentComparison() map[string]interface{} {
	mrd.mu.RLock()
	defer mrd.mu.RUnlock()

	strategies := map[string]map[string]string{
		"active-active": {
			"redundancy": "Very High",
			"cost":       "Very High",
			"complexity": "Very High",
			"rto":        "Seconds",
			"benefits":   "Perfect load distribution, high availability",
		},
		"active-passive": {
			"redundancy": "High",
			"cost":       "Medium",
			"complexity": "Medium",
			"rto":        "Minutes",
			"benefits":   "Cost efficient, simpler to manage",
		},
		"active-warm": {
			"redundancy": "High",
			"cost":       "Medium-High",
			"complexity": "Medium",
			"rto":        "Seconds-Minutes",
			"benefits":   "Balance between cost and availability",
		},
	}

	return map[string]interface{}{
		"strategies":         strategies,
		"total_configs":      len(mrd.deploymentConfigs),
		"total_data_centers": len(mrd.dataCenters),
	}
}

// GetLatencyMatrix returns the latency matrix for all regions
func (mrd *MultiRegionDeployment) GetLatencyMatrix() map[string]map[string]int {
	mrd.mu.RLock()
	defer mrd.mu.RUnlock()

	matrix := make(map[string]map[string]int)
	for fromRegion, latencies := range mrd.regionLatencies {
		matrix[fromRegion] = make(map[string]int)
		for toRegion, latency := range latencies {
			matrix[fromRegion][toRegion] = latency
		}
	}

	return matrix
}

// CalculateGeometricDistance calculates geometric distance between two points
func (mrd *MultiRegionDeployment) CalculateGeometricDistance(
	lat1, lon1, lat2, lon2 float64,
) float64 {
	// Haversine formula
	const earthRadiusKm = 6371.0

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	a := sin(dLat/2)*sin(dLat/2) +
		cos(toRadians(lat1))*cos(toRadians(lat2))*
			sin(dLon/2)*sin(dLon/2)

	c := 2.0 * atan2(sqrt(a), sqrt(1-a))
	distance := earthRadiusKm * c

	return distance
}

// Helper trigonometric functions
func toRadians(degrees float64) float64 {
	return degrees * 3.14159265359 / 180.0
}

func sin(rad float64) float64 {
	// Simplified sine approximation
	return rad - rad*rad*rad/6.0 + rad*rad*rad*rad*rad/120.0
}

func cos(rad float64) float64 {
	// Simplified cosine approximation
	return 1.0 - rad*rad/2.0 + rad*rad*rad*rad/24.0
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2.0
	}
	return z
}

func atan2(y, x float64) float64 {
	// Simplified atan2 approximation
	return 0.7853981633974483 // pi/4 as default
}
