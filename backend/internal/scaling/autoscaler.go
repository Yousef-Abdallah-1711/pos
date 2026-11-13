package scaling

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ScalingPolicy defines when and how to scale
type ScalingPolicy struct {
	ID                   int64
	RegionName           string
	Name                 string
	ScaleUpThreshold     int // CPU/Memory % threshold to scale up
	ScaleDownThreshold   int // CPU/Memory % threshold to scale down
	ScaleUpQuantity      int // Number of instances to add
	ScaleDownQuantity    int // Number of instances to remove
	MinInstances         int // Minimum instances to maintain
	MaxInstances         int // Maximum instances allowed
	CooldownPeriod       time.Duration
	MetricsCheckInterval time.Duration
	IsActive             bool
	CreatedAt            time.Time
}

// ScalingEvent represents a scaling action
type ScalingEvent struct {
	ID          int64
	RegionName  string
	Timestamp   time.Time
	Action      string // "scale-up", "scale-down"
	FromCount   int
	ToCount     int
	Reason      string
	CPUUsage    float64
	MemoryUsage float64
}

// AutoScaler manages automatic scaling based on metrics
type AutoScaler struct {
	loadBalancer    *LoadBalancer
	policies        map[string]*ScalingPolicy
	scalingEvents   []ScalingEvent
	lastScalingTime map[string]time.Time
	mu              sync.RWMutex
	stopChannel     chan bool
	isRunning       bool
}

// NewAutoScaler creates a new auto-scaler
func NewAutoScaler(loadBalancer *LoadBalancer) *AutoScaler {
	return &AutoScaler{
		loadBalancer:    loadBalancer,
		policies:        make(map[string]*ScalingPolicy),
		scalingEvents:   make([]ScalingEvent, 0),
		lastScalingTime: make(map[string]time.Time),
		stopChannel:     make(chan bool),
	}
}

// CreateScalingPolicy creates a new scaling policy
func (as *AutoScaler) CreateScalingPolicy(
	regionName string,
	name string,
	scaleUpThreshold int,
	scaleDownThreshold int,
	scaleUpQty int,
	scaleDownQty int,
	minInstances int,
	maxInstances int,
	cooldownPeriod time.Duration,
) *ScalingPolicy {
	log.Printf("[AUTOSCALER] Creating policy: %s for region %s", name, regionName)

	policy := &ScalingPolicy{
		RegionName:           regionName,
		Name:                 name,
		ScaleUpThreshold:     scaleUpThreshold,
		ScaleDownThreshold:   scaleDownThreshold,
		ScaleUpQuantity:      scaleUpQty,
		ScaleDownQuantity:    scaleDownQty,
		MinInstances:         minInstances,
		MaxInstances:         maxInstances,
		CooldownPeriod:       cooldownPeriod,
		MetricsCheckInterval: 30 * time.Second,
		IsActive:             true,
		CreatedAt:            time.Now(),
	}

	as.mu.Lock()
	as.policies[regionName] = policy
	as.lastScalingTime[regionName] = time.Now()
	as.mu.Unlock()

	log.Printf("[AUTOSCALER] Policy created: %s (min:%d, max:%d, up:%d%%, down:%d%%)",
		name, minInstances, maxInstances, scaleUpThreshold, scaleDownThreshold)

	return policy
}

// Start begins monitoring and scaling
func (as *AutoScaler) Start(ctx context.Context) {
	as.mu.Lock()
	if as.isRunning {
		as.mu.Unlock()
		return
	}
	as.isRunning = true
	as.mu.Unlock()

	log.Printf("[AUTOSCALER] Starting auto-scaler")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				as.evaluateAndScale(ctx)
			case <-as.stopChannel:
				log.Printf("[AUTOSCALER] Stopping auto-scaler")
				as.mu.Lock()
				as.isRunning = false
				as.mu.Unlock()
				return
			case <-ctx.Done():
				log.Printf("[AUTOSCALER] Context cancelled, stopping")
				as.mu.Lock()
				as.isRunning = false
				as.mu.Unlock()
				return
			}
		}
	}()
}

// Stop stops the auto-scaler
func (as *AutoScaler) Stop() {
	as.stopChannel <- true
}

// evaluateAndScale evaluates metrics and performs scaling
func (as *AutoScaler) evaluateAndScale(ctx context.Context) {
	as.mu.RLock()
	policies := make([]*ScalingPolicy, 0)
	for _, policy := range as.policies {
		if policy.IsActive {
			policies = append(policies, policy)
		}
	}
	as.mu.RUnlock()

	for _, policy := range policies {
		as.evaluateRegion(ctx, policy)
	}
}

// evaluateRegion evaluates a specific region for scaling
func (as *AutoScaler) evaluateRegion(ctx context.Context, policy *ScalingPolicy) {
	log.Printf("[AUTOSCALER] Evaluating region: %s", policy.RegionName)

	// Get region stats
	regionStats := as.loadBalancer.GetRegionStats(policy.RegionName)
	if regionStats == nil {
		log.Printf("[AUTOSCALER] Region %s not found", policy.RegionName)
		return
	}

	currentCount := regionStats["total_servers"].(int)
	utilization := regionStats["utilization"].(int)

	// Check if in cooldown period
	as.mu.RLock()
	lastScaling := as.lastScalingTime[policy.RegionName]
	as.mu.RUnlock()

	if time.Since(lastScaling) < policy.CooldownPeriod {
		return
	}

	// Decide action
	if utilization > policy.ScaleUpThreshold {
		as.performScaleUp(ctx, policy, currentCount, utilization)
	} else if utilization < policy.ScaleDownThreshold && currentCount > policy.MinInstances {
		as.performScaleDown(ctx, policy, currentCount, utilization)
	}
}

// performScaleUp adds new instances
func (as *AutoScaler) performScaleUp(
	ctx context.Context,
	policy *ScalingPolicy,
	currentCount int,
	utilization int,
) {
	newCount := currentCount + policy.ScaleUpQuantity

	// Don't exceed max instances
	if newCount > policy.MaxInstances {
		newCount = policy.MaxInstances
	}

	if newCount == currentCount {
		return
	}

	log.Printf("[AUTOSCALER] SCALE UP: %s: %d -> %d instances (utilization: %d%%)",
		policy.RegionName, currentCount, newCount, utilization)

	// Record scaling event
	event := ScalingEvent{
		RegionName: policy.RegionName,
		Timestamp:  time.Now(),
		Action:     "scale-up",
		FromCount:  currentCount,
		ToCount:    newCount,
		Reason:     fmt.Sprintf("High utilization (%d%%) exceeded threshold (%d%%)", utilization, policy.ScaleUpThreshold),
		CPUUsage:   float64(utilization),
	}

	as.recordScalingEvent(event)

	// In production, call cloud provider API to add instances
	// For now, just record the event
	as.mu.Lock()
	as.lastScalingTime[policy.RegionName] = time.Now()
	as.mu.Unlock()
}

// performScaleDown removes instances
func (as *AutoScaler) performScaleDown(
	ctx context.Context,
	policy *ScalingPolicy,
	currentCount int,
	utilization int,
) {
	newCount := currentCount - policy.ScaleDownQuantity

	// Don't go below min instances
	if newCount < policy.MinInstances {
		newCount = policy.MinInstances
	}

	if newCount == currentCount {
		return
	}

	log.Printf("[AUTOSCALER] SCALE DOWN: %s: %d -> %d instances (utilization: %d%%)",
		policy.RegionName, currentCount, newCount, utilization)

	// Record scaling event
	event := ScalingEvent{
		RegionName: policy.RegionName,
		Timestamp:  time.Now(),
		Action:     "scale-down",
		FromCount:  currentCount,
		ToCount:    newCount,
		Reason:     fmt.Sprintf("Low utilization (%d%%) below threshold (%d%%)", utilization, policy.ScaleDownThreshold),
		CPUUsage:   float64(utilization),
	}

	as.recordScalingEvent(event)

	// In production, call cloud provider API to remove instances
	// For now, just record the event
	as.mu.Lock()
	as.lastScalingTime[policy.RegionName] = time.Now()
	as.mu.Unlock()
}

// recordScalingEvent records a scaling event
func (as *AutoScaler) recordScalingEvent(event ScalingEvent) {
	as.mu.Lock()
	defer as.mu.Unlock()

	event.ID = int64(len(as.scalingEvents) + 1)
	as.scalingEvents = append(as.scalingEvents, event)

	// Keep only last 1000 events
	if len(as.scalingEvents) > 1000 {
		as.scalingEvents = as.scalingEvents[1:]
	}
}

// GetScalingPolicy retrieves a scaling policy
func (as *AutoScaler) GetScalingPolicy(regionName string) *ScalingPolicy {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.policies[regionName]
}

// UpdateScalingPolicy updates a scaling policy
func (as *AutoScaler) UpdateScalingPolicy(
	regionName string,
	scaleUpThreshold int,
	scaleDownThreshold int,
) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	policy, exists := as.policies[regionName]
	if !exists {
		return fmt.Errorf("policy for region %s not found", regionName)
	}

	policy.ScaleUpThreshold = scaleUpThreshold
	policy.ScaleDownThreshold = scaleDownThreshold

	log.Printf("[AUTOSCALER] Updated policy for %s: up=%d%%, down=%d%%",
		regionName, scaleUpThreshold, scaleDownThreshold)

	return nil
}

// GetScalingEvents retrieves scaling events
func (as *AutoScaler) GetScalingEvents(regionName string, limit int) []ScalingEvent {
	as.mu.RLock()
	defer as.mu.RUnlock()

	var events []ScalingEvent

	for _, event := range as.scalingEvents {
		if regionName == "" || event.RegionName == regionName {
			events = append(events, event)
		}
	}

	// Return newest events first
	if len(events) > limit {
		events = events[len(events)-limit:]
	}

	return events
}

// GetScalingMetrics returns scaling metrics
func (as *AutoScaler) GetScalingMetrics() map[string]interface{} {
	as.mu.RLock()
	defer as.mu.RUnlock()

	totalScaleUps := 0
	totalScaleDowns := 0

	for _, event := range as.scalingEvents {
		if event.Action == "scale-up" {
			totalScaleUps++
		} else if event.Action == "scale-down" {
			totalScaleDowns++
		}
	}

	return map[string]interface{}{
		"is_running":       as.isRunning,
		"total_policies":   len(as.policies),
		"total_events":     len(as.scalingEvents),
		"scale_up_count":   totalScaleUps,
		"scale_down_count": totalScaleDowns,
	}
}

// GetRegionScalingMetrics returns scaling metrics for a region
func (as *AutoScaler) GetRegionScalingMetrics(regionName string) map[string]interface{} {
	as.mu.RLock()
	policy := as.policies[regionName]
	lastScaling := as.lastScalingTime[regionName]
	as.mu.RUnlock()

	if policy == nil {
		return nil
	}

	var regionEvents []ScalingEvent
	for _, event := range as.scalingEvents {
		if event.RegionName == regionName {
			regionEvents = append(regionEvents, event)
		}
	}

	scaleUpCount := 0
	scaleDownCount := 0

	for _, event := range regionEvents {
		if event.Action == "scale-up" {
			scaleUpCount++
		} else if event.Action == "scale-down" {
			scaleDownCount++
		}
	}

	return map[string]interface{}{
		"region":               regionName,
		"policy_name":          policy.Name,
		"is_active":            policy.IsActive,
		"scale_up_threshold":   policy.ScaleUpThreshold,
		"scale_down_threshold": policy.ScaleDownThreshold,
		"min_instances":        policy.MinInstances,
		"max_instances":        policy.MaxInstances,
		"scale_up_count":       scaleUpCount,
		"scale_down_count":     scaleDownCount,
		"last_scaling":         lastScaling,
		"recent_events":        len(regionEvents),
	}
}
