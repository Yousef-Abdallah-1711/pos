//go:build legacy
// +build legacy

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MonitoringService handles real-time system monitoring
type MonitoringService struct {
	metrics         *SystemMetrics
	subscribers     map[int64]*MetricsSubscriber
	subscribersMu   sync.RWMutex
	alertRules      []*AlertRule
	alertRulesMu    sync.RWMutex
	collectionTick  *time.Ticker
	done            chan bool
	metricsHistory  *MetricsHistory
}

// SystemMetrics represents current system metrics
type SystemMetrics struct {
	Timestamp         time.Time                 `json:"timestamp"`
	CPUUsage          float64                   `json:"cpu_usage"`
	MemoryUsage       float64                   `json:"memory_usage"`
	DiskUsage         float64                   `json:"disk_usage"`
	ActiveConnections int64                     `json:"active_connections"`
	RequestsPerSecond float64                   `json:"requests_per_second"`
	ErrorRate         float64                   `json:"error_rate"`
	AvgResponseTime   int64                     `json:"avg_response_time"` // milliseconds
	DatabaseLatency   int64                     `json:"database_latency"`  // milliseconds
	CacheHitRate      float64                   `json:"cache_hit_rate"`
	QueueLength       int64                     `json:"queue_length"`
	ActiveProcesses   int64                     `json:"active_processes"`
	ServiceStatus     map[string]ServiceStatus  `json:"service_status"`
	CustomMetrics     map[string]float64        `json:"custom_metrics"`
}

// ServiceStatus represents status of a service
type ServiceStatus struct {
	Name       string    `json:"name"`
	Status     string    `json:"status"` // "healthy", "degraded", "unhealthy"
	Uptime     int64     `json:"uptime"` // seconds
	LastCheck  time.Time `json:"last_check"`
	Message    string    `json:"message"`
}

// MetricsSubscriber represents a subscriber to metrics updates
type MetricsSubscriber struct {
	ID            int64
	RestaurantID  int64
	UpdateChan    chan *SystemMetrics
	AlertChan     chan *Alert
	Done          chan bool
	LastHeartbeat time.Time
}

// Alert represents a triggered alert
type Alert struct {
	ID          string            `json:"id"`
	AlertID     string            `json:"alert_id"`
	Severity    string            `json:"severity"` // "info", "warning", "critical"
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Metric      string            `json:"metric"`
	Threshold   float64           `json:"threshold"`
	CurrentValue float64          `json:"current_value"`
	Timestamp   time.Time         `json:"timestamp"`
	Resolved    bool              `json:"resolved"`
	ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
}

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	ID           string                 `json:"id"`
	RestaurantID int64                  `json:"restaurant_id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Metric       string                 `json:"metric"`
	Condition    string                 `json:"condition"` // ">", "<", ">=", "<=", "==", "!="
	Threshold    float64                `json:"threshold"`
	Duration     int                    `json:"duration"` // seconds
	Severity     string                 `json:"severity"` // "info", "warning", "critical"
	Enabled      bool                   `json:"enabled"`
	Notification map[string]interface{} `json:"notification"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// MetricsHistory stores historical metrics
type MetricsHistory struct {
	metrics []*HistoricalMetric
	mu      sync.RWMutex
	maxSize int
}

// HistoricalMetric represents a single historical data point
type HistoricalMetric struct {
	Timestamp time.Time                `json:"timestamp"`
	Metrics   *SystemMetrics           `json:"metrics"`
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		metrics:        &SystemMetrics{ServiceStatus: make(map[string]ServiceStatus), CustomMetrics: make(map[string]float64)},
		subscribers:    make(map[int64]*MetricsSubscriber),
		alertRules:     make([]*AlertRule, 0),
		collectionTick: time.NewTicker(5 * time.Second),
		done:           make(chan bool),
		metricsHistory: &MetricsHistory{
			metrics: make([]*HistoricalMetric, 0),
			maxSize: 1000,
		},
	}
}

// Start starts the monitoring service
func (ms *MonitoringService) Start() {
	fmt.Println("DEBUG: Starting monitoring service")

	go func() {
		for {
			select {
			case <-ms.collectionTick.C:
				ms.collectMetrics()
				ms.checkAlertRules()
				ms.broadcastMetrics()
			case <-ms.done:
				ms.collectionTick.Stop()
				return
			}
		}
	}()
}

// Subscribe subscribes to metrics updates
func (ms *MonitoringService) Subscribe(restaurantID int64) *MetricsSubscriber {
	ms.subscribersMu.Lock()
	defer ms.subscribersMu.Unlock()

	subscriber := &MetricsSubscriber{
		ID:            int64(len(ms.subscribers)),
		RestaurantID:  restaurantID,
		UpdateChan:    make(chan *SystemMetrics, 10),
		AlertChan:     make(chan *Alert, 10),
		Done:          make(chan bool),
		LastHeartbeat: time.Now(),
	}

	ms.subscribers[subscriber.ID] = subscriber
	fmt.Printf("DEBUG: Subscriber %d added (total: %d)\n", subscriber.ID, len(ms.subscribers))

	return subscriber
}

// Unsubscribe unsubscribes from metrics updates
func (ms *MonitoringService) Unsubscribe(subscriberID int64) {
	ms.subscribersMu.Lock()
	defer ms.subscribersMu.Unlock()

	if subscriber, exists := ms.subscribers[subscriberID]; exists {
		subscriber.Done <- true
		delete(ms.subscribers, subscriberID)
		fmt.Printf("DEBUG: Subscriber %d removed (total: %d)\n", subscriberID, len(ms.subscribers))
	}
}

// GetCurrentMetrics returns current metrics
func (ms *MonitoringService) GetCurrentMetrics() *SystemMetrics {
	return ms.metrics
}

// GetMetricsHistory returns historical metrics
func (ms *MonitoringService) GetMetricsHistory(duration time.Duration) []*HistoricalMetric {
	ms.metricsHistory.mu.RLock()
	defer ms.metricsHistory.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	history := make([]*HistoricalMetric, 0)

	for _, metric := range ms.metricsHistory.metrics {
		if metric.Timestamp.After(cutoff) {
			history = append(history, metric)
		}
	}

	return history
}

// AddAlertRule adds an alert rule
func (ms *MonitoringService) AddAlertRule(rule *AlertRule) error {
	ms.alertRulesMu.Lock()
	defer ms.alertRulesMu.Unlock()

	if rule.ID == "" {
		rule.ID = fmt.Sprintf("rule_%d", time.Now().UnixNano())
	}
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	ms.alertRules = append(ms.alertRules, rule)
	fmt.Printf("DEBUG: Alert rule added: %s\n", rule.Name)

	return nil
}

// UpdateAlertRule updates an alert rule
func (ms *MonitoringService) UpdateAlertRule(rule *AlertRule) error {
	ms.alertRulesMu.Lock()
	defer ms.alertRulesMu.Unlock()

	for i, r := range ms.alertRules {
		if r.ID == rule.ID {
			rule.UpdatedAt = time.Now()
			ms.alertRules[i] = rule
			fmt.Printf("DEBUG: Alert rule updated: %s\n", rule.Name)
			return nil
		}
	}

	return fmt.Errorf("alert rule not found")
}

// DeleteAlertRule deletes an alert rule
func (ms *MonitoringService) DeleteAlertRule(ruleID string) error {
	ms.alertRulesMu.Lock()
	defer ms.alertRulesMu.Unlock()

	for i, rule := range ms.alertRules {
		if rule.ID == ruleID {
			ms.alertRules = append(ms.alertRules[:i], ms.alertRules[i+1:]...)
			fmt.Printf("DEBUG: Alert rule deleted: %s\n", rule.Name)
			return nil
		}
	}

	return fmt.Errorf("alert rule not found")
}

// GetAlertRules returns all alert rules
func (ms *MonitoringService) GetAlertRules() []*AlertRule {
	ms.alertRulesMu.RLock()
	defer ms.alertRulesMu.RUnlock()

	rules := make([]*AlertRule, len(ms.alertRules))
	copy(rules, ms.alertRules)
	return rules
}

// SetCustomMetric sets a custom metric value
func (ms *MonitoringService) SetCustomMetric(name string, value float64) {
	ms.metrics.CustomMetrics[name] = value
	fmt.Printf("DEBUG: Custom metric set: %s = %.2f\n", name, value)
}

// SetServiceStatus sets service status
func (ms *MonitoringService) SetServiceStatus(name string, status ServiceStatus) {
	status.LastCheck = time.Now()
	ms.metrics.ServiceStatus[name] = status
	fmt.Printf("DEBUG: Service status updated: %s = %s\n", name, status.Status)
}

// Helper functions

func (ms *MonitoringService) collectMetrics() {
	ms.metrics.Timestamp = time.Now()

	// Simulate metric collection (in production, would gather real metrics)
	ms.metrics.CPUUsage = 45.3 + (time.Now().Unix()%10 - 5) // Simulate variation
	ms.metrics.MemoryUsage = 62.5 + (time.Now().Unix()%8 - 4)
	ms.metrics.DiskUsage = 78.2
	ms.metrics.ActiveConnections = 150 + (time.Now().Unix() % 50)
	ms.metrics.RequestsPerSecond = 1234.5 + (time.Now().Unix()%500 - 250)
	ms.metrics.ErrorRate = 0.5
	ms.metrics.AvgResponseTime = 145 + (time.Now().Unix() % 100)
	ms.metrics.DatabaseLatency = 45
	ms.metrics.CacheHitRate = 92.3
	ms.metrics.QueueLength = 25

	// Update service statuses
	if ms.metrics.ErrorRate > 1.0 {
		ms.metrics.ServiceStatus["api"] = ServiceStatus{
			Name:      "api",
			Status:    "degraded",
			Message:   "High error rate detected",
			LastCheck: time.Now(),
		}
	} else {
		ms.metrics.ServiceStatus["api"] = ServiceStatus{
			Name:      "api",
			Status:    "healthy",
			Message:   "API operating normally",
			LastCheck: time.Now(),
		}
	}

	if ms.metrics.DatabaseLatency > 100 {
		ms.metrics.ServiceStatus["database"] = ServiceStatus{
			Name:      "database",
			Status:    "degraded",
			Message:   "High database latency",
			LastCheck: time.Now(),
		}
	} else {
		ms.metrics.ServiceStatus["database"] = ServiceStatus{
			Name:      "database",
			Status:    "healthy",
			Message:   "Database operating normally",
			LastCheck: time.Now(),
		}
	}

	// Store in history
	ms.metricsHistory.addMetric(&HistoricalMetric{
		Timestamp: ms.metrics.Timestamp,
		Metrics: &SystemMetrics{
			Timestamp:         ms.metrics.Timestamp,
			CPUUsage:          ms.metrics.CPUUsage,
			MemoryUsage:       ms.metrics.MemoryUsage,
			DiskUsage:         ms.metrics.DiskUsage,
			ActiveConnections: ms.metrics.ActiveConnections,
			RequestsPerSecond: ms.metrics.RequestsPerSecond,
			ErrorRate:         ms.metrics.ErrorRate,
			AvgResponseTime:   ms.metrics.AvgResponseTime,
			DatabaseLatency:   ms.metrics.DatabaseLatency,
			CacheHitRate:      ms.metrics.CacheHitRate,
			QueueLength:       ms.metrics.QueueLength,
		},
	})
}

func (ms *MonitoringService) checkAlertRules() {
	ms.alertRulesMu.RLock()
	rules := make([]*AlertRule, len(ms.alertRules))
	copy(rules, ms.alertRules)
	ms.alertRulesMu.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		var metricValue float64
		var triggered bool

		// Evaluate metric
		switch rule.Metric {
		case "cpu_usage":
			metricValue = ms.metrics.CPUUsage
		case "memory_usage":
			metricValue = ms.metrics.MemoryUsage
		case "disk_usage":
			metricValue = ms.metrics.DiskUsage
		case "error_rate":
			metricValue = ms.metrics.ErrorRate
		case "avg_response_time":
			metricValue = float64(ms.metrics.AvgResponseTime)
		case "database_latency":
			metricValue = float64(ms.metrics.DatabaseLatency)
		default:
			if value, ok := ms.metrics.CustomMetrics[rule.Metric]; ok {
				metricValue = value
			}
		}

		// Check condition
		switch rule.Condition {
		case ">":
			triggered = metricValue > rule.Threshold
		case "<":
			triggered = metricValue < rule.Threshold
		case ">=":
			triggered = metricValue >= rule.Threshold
		case "<=":
			triggered = metricValue <= rule.Threshold
		case "==":
			triggered = metricValue == rule.Threshold
		case "!=":
			triggered = metricValue != rule.Threshold
		}

		if triggered {
			alert := &Alert{
				ID:           fmt.Sprintf("alert_%d", time.Now().UnixNano()),
				AlertID:      rule.ID,
				Severity:     rule.Severity,
				Title:        rule.Name,
				Description:  rule.Description,
				Metric:       rule.Metric,
				Threshold:    rule.Threshold,
				CurrentValue: metricValue,
				Timestamp:    time.Now(),
				Resolved:     false,
			}

			ms.broadcastAlert(alert)
			fmt.Printf("DEBUG: Alert triggered: %s (metric: %s, value: %.2f, threshold: %.2f)\n",
				alert.Title, alert.Metric, metricValue, rule.Threshold)
		}
	}
}

func (ms *MonitoringService) broadcastMetrics() {
	ms.subscribersMu.RLock()
	subscribers := make([]*MetricsSubscriber, 0, len(ms.subscribers))
	for _, sub := range ms.subscribers {
		subscribers = append(subscribers, sub)
	}
	ms.subscribersMu.RUnlock()

	for _, sub := range subscribers {
		select {
		case sub.UpdateChan <- ms.metrics:
		default:
			// Channel full, skip this update
		}
	}
}

func (ms *MonitoringService) broadcastAlert(alert *Alert) {
	ms.subscribersMu.RLock()
	subscribers := make([]*MetricsSubscriber, 0, len(ms.subscribers))
	for _, sub := range ms.subscribers {
		subscribers = append(subscribers, sub)
	}
	ms.subscribersMu.RUnlock()

	for _, sub := range subscribers {
		select {
		case sub.AlertChan <- alert:
		default:
			// Channel full, skip this alert
		}
	}
}

// MetricsHistory methods

func (mh *MetricsHistory) addMetric(metric *HistoricalMetric) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	mh.metrics = append(mh.metrics, metric)

	// Keep only recent metrics
	if len(mh.metrics) > mh.maxSize {
		mh.metrics = mh.metrics[len(mh.metrics)-mh.maxSize:]
	}
}

func (mh *MetricsHistory) getMetrics() []*HistoricalMetric {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	metrics := make([]*HistoricalMetric, len(mh.metrics))
	copy(metrics, mh.metrics)
	return metrics
}

// Shutdown stops the monitoring service
func (ms *MonitoringService) Shutdown(ctx context.Context) error {
	fmt.Println("DEBUG: Shutting down monitoring service")
	ms.done <- true
	<-ctx.Done()
	return nil
}
