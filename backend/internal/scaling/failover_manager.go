package scaling

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// FailoverPolicy defines failover behavior
type FailoverPolicy struct {
	ID                   int64
	Name                 string
	PrimaryRegion        string
	SecondaryRegion      string
	HealthCheckThreshold int // % threshold to trigger failover
	FailoverTimeout      time.Duration
	AutomaticFailover    bool
	DataReplicationLag   time.Duration // Max acceptable lag
	IsActive             bool
	CreatedAt            time.Time
}

// FailoverEvent represents a failover action
type FailoverEvent struct {
	ID           int64
	Timestamp    time.Time
	FromRegion   string
	ToRegion     string
	Reason       string
	IsAutomatic  bool
	DataLagMs    int64
	Status       string // "initiated", "in_progress", "completed", "failed"
	ErrorMessage string
}

// DisasterRecoveryPlan defines disaster recovery procedures
type DisasterRecoveryPlan struct {
	ID                     int64
	Name                   string
	RecoveryTimeObjective  time.Duration // RTO: how fast to recover
	RecoveryPointObjective time.Duration // RPO: how much data loss acceptable
	BackupFrequency        time.Duration // How often to backup
	BackupRetention        time.Duration // How long to keep backups
	ReplicationEnabled     bool
	ReplicationTargets     []string
	IsActive               bool
	CreatedAt              time.Time
}

// FailoverManager handles failover and disaster recovery
type FailoverManager struct {
	loadBalancer      *LoadBalancer
	healthChecker     *HealthChecker
	failoverPolicies  map[string]*FailoverPolicy
	disasterPlans     map[string]*DisasterRecoveryPlan
	failoverEvents    []FailoverEvent
	currentPrimary    map[string]string  // region_pair -> current primary
	replicationStatus map[string]float64 // region -> replication_lag_ms
	mu                sync.RWMutex
}

// NewFailoverManager creates a new failover manager
func NewFailoverManager(
	loadBalancer *LoadBalancer,
	healthChecker *HealthChecker,
) *FailoverManager {
	return &FailoverManager{
		loadBalancer:      loadBalancer,
		healthChecker:     healthChecker,
		failoverPolicies:  make(map[string]*FailoverPolicy),
		disasterPlans:     make(map[string]*DisasterRecoveryPlan),
		failoverEvents:    make([]FailoverEvent, 0),
		currentPrimary:    make(map[string]string),
		replicationStatus: make(map[string]float64),
	}
}

// CreateFailoverPolicy creates a new failover policy
func (fm *FailoverManager) CreateFailoverPolicy(
	name string,
	primaryRegion string,
	secondaryRegion string,
	healthThreshold int,
	failoverTimeout time.Duration,
	autoFailover bool,
) *FailoverPolicy {
	log.Printf("[FAILOVER] Creating policy: %s (%s -> %s)", name, primaryRegion, secondaryRegion)

	policy := &FailoverPolicy{
		Name:                 name,
		PrimaryRegion:        primaryRegion,
		SecondaryRegion:      secondaryRegion,
		HealthCheckThreshold: healthThreshold,
		FailoverTimeout:      failoverTimeout,
		AutomaticFailover:    autoFailover,
		DataReplicationLag:   5 * time.Second,
		IsActive:             true,
		CreatedAt:            time.Now(),
	}

	policyKey := fmt.Sprintf("%s-%s", primaryRegion, secondaryRegion)
	fm.mu.Lock()
	fm.failoverPolicies[policyKey] = policy
	fm.currentPrimary[policyKey] = primaryRegion
	fm.mu.Unlock()

	log.Printf("[FAILOVER] Policy created: %s (auto: %v, threshold: %d%%)",
		name, autoFailover, healthThreshold)

	return policy
}

// CreateDisasterRecoveryPlan creates a new DR plan
func (fm *FailoverManager) CreateDisasterRecoveryPlan(
	name string,
	rto time.Duration,
	rpo time.Duration,
	backupFreq time.Duration,
	backupRetention time.Duration,
	replicate bool,
	targets []string,
) *DisasterRecoveryPlan {
	log.Printf("[FAILOVER] Creating DR plan: %s (RTO: %v, RPO: %v)", name, rto, rpo)

	plan := &DisasterRecoveryPlan{
		Name:                   name,
		RecoveryTimeObjective:  rto,
		RecoveryPointObjective: rpo,
		BackupFrequency:        backupFreq,
		BackupRetention:        backupRetention,
		ReplicationEnabled:     replicate,
		ReplicationTargets:     targets,
		IsActive:               true,
		CreatedAt:              time.Now(),
	}

	fm.mu.Lock()
	fm.disasterPlans[name] = plan
	fm.mu.Unlock()

	return plan
}

// EvaluateFailoverNeeded checks if failover is needed
func (fm *FailoverManager) EvaluateFailoverNeeded(policyKey string) bool {
	fm.mu.RLock()
	policy := fm.failoverPolicies[policyKey]
	fm.mu.RUnlock()

	if policy == nil || !policy.AutomaticFailover {
		return false
	}

	// Get primary region health
	primaryStats := fm.loadBalancer.GetRegionStats(policy.PrimaryRegion)
	if primaryStats == nil {
		return false
	}

	healthyCount := primaryStats["healthy_servers"].(int)
	totalCount := primaryStats["total_servers"].(int)

	if totalCount == 0 {
		return false
	}

	healthPercentage := (healthyCount * 100) / totalCount

	log.Printf("[FAILOVER] Health check: %s = %d%% (threshold: %d%%)",
		policy.PrimaryRegion, healthPercentage, policy.HealthCheckThreshold)

	return healthPercentage < policy.HealthCheckThreshold
}

// InitiateFailover initiates failover to secondary region
func (fm *FailoverManager) InitiateFailover(
	ctx context.Context,
	policyKey string,
	isAutomatic bool,
	reason string,
) error {
	fm.mu.Lock()
	policy := fm.failoverPolicies[policyKey]
	currentPrimary := fm.currentPrimary[policyKey]
	fm.mu.Unlock()

	if policy == nil {
		return fmt.Errorf("policy %s not found", policyKey)
	}

	if currentPrimary == policy.SecondaryRegion {
		return fmt.Errorf("already failed over to secondary")
	}

	log.Printf("[FAILOVER] Initiating failover: %s -> %s (automatic: %v, reason: %s)",
		policy.PrimaryRegion, policy.SecondaryRegion, isAutomatic, reason)

	// Create failover event
	event := FailoverEvent{
		ID:          int64(len(fm.failoverEvents) + 1),
		Timestamp:   time.Now(),
		FromRegion:  currentPrimary,
		ToRegion:    policy.SecondaryRegion,
		Reason:      reason,
		IsAutomatic: isAutomatic,
		Status:      "initiated",
	}

	// Check replication lag
	replicationLag := fm.getReplicationLag(currentPrimary, policy.SecondaryRegion)
	if replicationLag > policy.DataReplicationLag {
		event.Status = "failed"
		event.ErrorMessage = fmt.Sprintf("Replication lag too high: %v > %v", replicationLag, policy.DataReplicationLag)
		fm.recordFailoverEvent(event)
		return fmt.Errorf("replication lag too high: %v", replicationLag)
	}

	event.DataLagMs = int64(replicationLag.Milliseconds())
	event.Status = "in_progress"

	// Perform failover
	err := fm.performFailover(ctx, policy)
	if err != nil {
		event.Status = "failed"
		event.ErrorMessage = err.Error()
		fm.recordFailoverEvent(event)
		return err
	}

	event.Status = "completed"
	fm.recordFailoverEvent(event)

	// Update current primary
	fm.mu.Lock()
	fm.currentPrimary[policyKey] = policy.SecondaryRegion
	fm.mu.Unlock()

	log.Printf("[FAILOVER] Failover completed: traffic now routed to %s", policy.SecondaryRegion)

	return nil
}

// performFailover performs the actual failover operations
func (fm *FailoverManager) performFailover(
	ctx context.Context,
	policy *FailoverPolicy,
) error {
	// In production, this would:
	// 1. Update DNS records (GeoDNS)
	// 2. Update load balancer routing
	// 3. Start services in secondary region
	// 4. Verify data consistency
	// 5. Monitor for issues

	log.Printf("[FAILOVER] Switching traffic from %s to %s", policy.PrimaryRegion, policy.SecondaryRegion)

	// Wait for propagation
	time.Sleep(1 * time.Second)

	// Verify secondary region is accepting traffic
	secondaryStats := fm.loadBalancer.GetRegionStats(policy.SecondaryRegion)
	if secondaryStats == nil {
		return fmt.Errorf("secondary region %s not found", policy.SecondaryRegion)
	}

	healthyCount := secondaryStats["healthy_servers"].(int)
	if healthyCount == 0 {
		return fmt.Errorf("no healthy servers in secondary region")
	}

	return nil
}

// InitiateFailback initiates failback to primary region
func (fm *FailoverManager) InitiateFailback(ctx context.Context, policyKey string) error {
	fm.mu.Lock()
	policy := fm.failoverPolicies[policyKey]
	currentSecondary := fm.currentPrimary[policyKey]
	fm.mu.Unlock()

	if policy == nil {
		return fmt.Errorf("policy %s not found", policyKey)
	}

	if currentSecondary == policy.PrimaryRegion {
		return fmt.Errorf("already on primary region")
	}

	log.Printf("[FAILOVER] Initiating failback: %s -> %s", currentSecondary, policy.PrimaryRegion)

	// Check if primary is healthy
	primaryStats := fm.loadBalancer.GetRegionStats(policy.PrimaryRegion)
	if primaryStats == nil {
		return fmt.Errorf("primary region %s not found", policy.PrimaryRegion)
	}

	healthyCount := primaryStats["healthy_servers"].(int)
	totalCount := primaryStats["total_servers"].(int)

	if totalCount == 0 || healthyCount == 0 {
		return fmt.Errorf("primary region not healthy")
	}

	// Perform failback
	err := fm.performFailback(ctx, policy)
	if err != nil {
		return err
	}

	// Record event
	event := FailoverEvent{
		ID:         int64(len(fm.failoverEvents) + 1),
		Timestamp:  time.Now(),
		FromRegion: currentSecondary,
		ToRegion:   policy.PrimaryRegion,
		Reason:     "Primary region recovered",
		Status:     "completed",
	}
	fm.recordFailoverEvent(event)

	// Update current primary
	fm.mu.Lock()
	fm.currentPrimary[policyKey] = policy.PrimaryRegion
	fm.mu.Unlock()

	return nil
}

// performFailback performs the actual failback operations
func (fm *FailoverManager) performFailback(
	ctx context.Context,
	policy *FailoverPolicy,
) error {
	log.Printf("[FAILOVER] Switching traffic back to %s", policy.PrimaryRegion)

	time.Sleep(1 * time.Second)

	primaryStats := fm.loadBalancer.GetRegionStats(policy.PrimaryRegion)
	if primaryStats == nil {
		return fmt.Errorf("primary region not accessible")
	}

	return nil
}

// getReplicationLag returns replication lag between regions
func (fm *FailoverManager) getReplicationLag(fromRegion string, toRegion string) time.Duration {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	lag, exists := fm.replicationStatus[toRegion]
	if !exists {
		// Default 100ms lag if not tracked
		return 100 * time.Millisecond
	}

	return time.Duration(lag) * time.Millisecond
}

// UpdateReplicationLag updates the replication lag for a region
func (fm *FailoverManager) UpdateReplicationLag(region string, lagMs float64) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.replicationStatus[region] = lagMs
}

// recordFailoverEvent records a failover event
func (fm *FailoverManager) recordFailoverEvent(event FailoverEvent) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.failoverEvents = append(fm.failoverEvents, event)

	// Keep only last 1000 events
	if len(fm.failoverEvents) > 1000 {
		fm.failoverEvents = fm.failoverEvents[1:]
	}
}

// GetFailoverStatus returns current failover status
func (fm *FailoverManager) GetFailoverStatus(policyKey string) map[string]interface{} {
	fm.mu.RLock()
	policy := fm.failoverPolicies[policyKey]
	currentPrimary := fm.currentPrimary[policyKey]
	fm.mu.RUnlock()

	if policy == nil {
		return nil
	}

	return map[string]interface{}{
		"policy_name":      policy.Name,
		"primary_region":   policy.PrimaryRegion,
		"secondary_region": policy.SecondaryRegion,
		"current_primary":  currentPrimary,
		"is_failed_over":   currentPrimary != policy.PrimaryRegion,
		"auto_failover":    policy.AutomaticFailover,
		"health_threshold": policy.HealthCheckThreshold,
		"is_active":        policy.IsActive,
	}
}

// GetFailoverHistory returns failover event history
func (fm *FailoverManager) GetFailoverHistory(limit int) []FailoverEvent {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	events := make([]FailoverEvent, len(fm.failoverEvents))
	copy(events, fm.failoverEvents)

	if len(events) > limit {
		events = events[len(events)-limit:]
	}

	return events
}

// GetDisasterRecoveryStatus returns DR plan status
func (fm *FailoverManager) GetDisasterRecoveryStatus(planName string) map[string]interface{} {
	fm.mu.RLock()
	plan := fm.disasterPlans[planName]
	fm.mu.RUnlock()

	if plan == nil {
		return nil
	}

	return map[string]interface{}{
		"name":                plan.Name,
		"rto":                 plan.RecoveryTimeObjective.String(),
		"rpo":                 plan.RecoveryPointObjective.String(),
		"backup_frequency":    plan.BackupFrequency.String(),
		"backup_retention":    plan.BackupRetention.String(),
		"replication_enabled": plan.ReplicationEnabled,
		"replication_targets": plan.ReplicationTargets,
		"is_active":           plan.IsActive,
		"created_at":          plan.CreatedAt,
	}
}
