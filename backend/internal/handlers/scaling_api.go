//go:build legacy
// +build legacy

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"pos-saas/internal/scaling"
)

// LoadBalancerHandler handles load balancer API requests
type LoadBalancerHandler struct {
	lb *scaling.LoadBalancer
}

// NewLoadBalancerHandler creates a new load balancer handler
func NewLoadBalancerHandler(lb *scaling.LoadBalancer) *LoadBalancerHandler {
	return &LoadBalancerHandler{lb: lb}
}

// RegisterLoadBalancerRoutes registers load balancer API routes
func (h *LoadBalancerHandler) RegisterLoadBalancerRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /scaling/load-balancer/select", h.SelectServer)
	router.HandleFunc("GET /scaling/load-balancer/strategy", h.GetCurrentStrategy)
	router.HandleFunc("POST /scaling/load-balancer/strategy", h.SetStrategy)
	router.HandleFunc("GET /scaling/load-balancer/regions", h.ListRegions)
	router.HandleFunc("POST /scaling/load-balancer/regions", h.RegisterRegion)
	router.HandleFunc("GET /scaling/load-balancer/servers", h.ListServers)
	router.HandleFunc("POST /scaling/load-balancer/servers", h.RegisterServer)
	router.HandleFunc("POST /scaling/load-balancer/servers/:id/metrics", h.UpdateServerMetrics)
	router.HandleFunc("POST /scaling/load-balancer/servers/:id/health", h.SetServerHealth)
	router.HandleFunc("GET /scaling/load-balancer/stats/global", h.GetGlobalStats)
	router.HandleFunc("GET /scaling/load-balancer/stats/regions", h.GetRegionStats)
	router.HandleFunc("GET /scaling/load-balancer/stats/servers", h.GetServerStats)
	router.HandleFunc("POST /scaling/load-balancer/session-affinity", h.ConfigureSessionAffinity)
}

// SelectServer selects the next server for a request
func (h *LoadBalancerHandler) SelectServer(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		http.Error(w, "region parameter required", http.StatusBadRequest)
		return
	}

	server := h.lb.SelectServer(region)
	if server == nil {
		http.Error(w, "no available servers", http.StatusServiceUnavailable)
		return
	}

	response := map[string]interface{}{
		"server_id":   server.ID,
		"region":      server.Region,
		"endpoint":    server.Endpoint,
		"port":        server.Port,
		"weight":      server.Weight,
		"current_load": server.CurrentLoad,
		"is_healthy":  server.IsHealthy,
		"response_time": server.ResponseTime,
		"cpu_usage":   server.CPUUsage,
		"memory_usage": server.MemoryUsage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentStrategy returns the current load balancing strategy
func (h *LoadBalancerHandler) GetCurrentStrategy(w http.ResponseWriter, r *http.Request) {
	strategy := h.lb.GetCurrentStrategy()
	response := map[string]string{
		"strategy": strategy,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetStrategyRequest represents a strategy change request
type SetStrategyRequest struct {
	Strategy string `json:"strategy"`
}

// SetStrategy changes the load balancing strategy
func (h *LoadBalancerHandler) SetStrategy(w http.ResponseWriter, r *http.Request) {
	var req SetStrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.lb.SetStrategy(req.Strategy)
	response := map[string]string{
		"strategy": req.Strategy,
		"message":  "strategy changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListRegions returns all registered regions
func (h *LoadBalancerHandler) ListRegions(w http.ResponseWriter, r *http.Request) {
	regions := h.lb.GetAllRegions()
	response := map[string]interface{}{
		"count":   len(regions),
		"regions": regions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRegionRequest represents a region registration request
type RegisterRegionRequest struct {
	Name      string  `json:"name"`
	Continent string  `json:"continent"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Priority  int     `json:"priority"`
}

// RegisterRegion registers a new region
func (h *LoadBalancerHandler) RegisterRegion(w http.ResponseWriter, r *http.Request) {
	var req RegisterRegionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	region := h.lb.RegisterRegion(req.Name, req.Continent, req.Latitude, req.Longitude, req.Priority)
	response := map[string]interface{}{
		"name":      region.Name,
		"continent": region.Continent,
		"latitude":  region.Latitude,
		"longitude": region.Longitude,
		"priority":  region.Priority,
		"is_active": region.IsActive,
		"message":   "region registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListServers returns all registered servers
func (h *LoadBalancerHandler) ListServers(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")

	var servers []*scaling.Server
	if region != "" {
		servers = h.lb.GetServersByRegion(region)
	} else {
		servers = h.lb.GetAllServers()
	}

	response := map[string]interface{}{
		"count":   len(servers),
		"servers": servers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterServerRequest represents a server registration request
type RegisterServerRequest struct {
	ID       string  `json:"id"`
	Region   string  `json:"region"`
	Endpoint string  `json:"endpoint"`
	Port     int     `json:"port"`
	Weight   int     `json:"weight"`
	Capacity int     `json:"capacity"`
}

// RegisterServer registers a new server
func (h *LoadBalancerHandler) RegisterServer(w http.ResponseWriter, r *http.Request) {
	var req RegisterServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	server := h.lb.RegisterServer(req.ID, req.Region, req.Endpoint, req.Port, req.Weight, req.Capacity)
	response := map[string]interface{}{
		"id":       server.ID,
		"region":   server.Region,
		"endpoint": server.Endpoint,
		"port":     server.Port,
		"weight":   server.Weight,
		"capacity": server.Capacity,
		"message":  "server registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateServerMetricsRequest represents a server metrics update request
type UpdateServerMetricsRequest struct {
	CurrentLoad  int     `json:"current_load"`
	ResponseTime int     `json:"response_time"`
	ErrorRate    float64 `json:"error_rate"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
}

// UpdateServerMetrics updates server performance metrics
func (h *LoadBalancerHandler) UpdateServerMetrics(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/scaling/load-balancer/servers/")
	serverID = strings.TrimSuffix(serverID, "/metrics")

	var req UpdateServerMetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.lb.UpdateServerMetrics(serverID, req.CurrentLoad, req.ResponseTime, req.ErrorRate, req.CPUUsage, req.MemoryUsage)
	response := map[string]interface{}{
		"server_id":    serverID,
		"current_load": req.CurrentLoad,
		"response_time": req.ResponseTime,
		"cpu_usage":    req.CPUUsage,
		"memory_usage": req.MemoryUsage,
		"message":      "metrics updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetServerHealthRequest represents a server health update request
type SetServerHealthRequest struct {
	IsHealthy bool   `json:"is_healthy"`
	Reason    string `json:"reason"`
}

// SetServerHealth sets server health status
func (h *LoadBalancerHandler) SetServerHealth(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/scaling/load-balancer/servers/")
	serverID = strings.TrimSuffix(serverID, "/health")

	var req SetServerHealthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.lb.SetServerHealth(serverID, req.IsHealthy)
	response := map[string]interface{}{
		"server_id":  serverID,
		"is_healthy": req.IsHealthy,
		"reason":     req.Reason,
		"message":    "health status updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGlobalStats returns global load balancing statistics
func (h *LoadBalancerHandler) GetGlobalStats(w http.ResponseWriter, r *http.Request) {
	stats := h.lb.GetGlobalStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetRegionStats returns statistics for all regions
func (h *LoadBalancerHandler) GetRegionStats(w http.ResponseWriter, r *http.Request) {
	stats := h.lb.GetRegionStats()
	response := map[string]interface{}{
		"regions": stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetServerStats returns statistics for all servers
func (h *LoadBalancerHandler) GetServerStats(w http.ResponseWriter, r *http.Request) {
	stats := h.lb.GetServerStats()
	response := map[string]interface{}{
		"servers": stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConfigureSessionAffinityRequest represents a session affinity configuration request
type ConfigureSessionAffinityRequest struct {
	Enabled    bool `json:"enabled"`
	TTLSeconds int  `json:"ttl_seconds"`
}

// ConfigureSessionAffinity configures session affinity (sticky sessions)
func (h *LoadBalancerHandler) ConfigureSessionAffinity(w http.ResponseWriter, r *http.Request) {
	var req ConfigureSessionAffinityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.lb.SetSessionAffinity(req.Enabled, req.TTLSeconds)
	response := map[string]interface{}{
		"enabled":     req.Enabled,
		"ttl_seconds": req.TTLSeconds,
		"message":     "session affinity configured successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AutoScalerHandler handles auto-scaler API requests
type AutoScalerHandler struct {
	autoscaler *scaling.AutoScaler
}

// NewAutoScalerHandler creates a new auto-scaler handler
func NewAutoScalerHandler(autoscaler *scaling.AutoScaler) *AutoScalerHandler {
	return &AutoScalerHandler{autoscaler: autoscaler}
}

// RegisterAutoScalerRoutes registers auto-scaler API routes
func (h *AutoScalerHandler) RegisterAutoScalerRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /scaling/autoscaler/policies", h.CreateScalingPolicy)
	router.HandleFunc("GET /scaling/autoscaler/policies", h.ListScalingPolicies)
	router.HandleFunc("GET /scaling/autoscaler/policies/:region", h.GetScalingPolicy)
	router.HandleFunc("PUT /scaling/autoscaler/policies/:region", h.UpdateScalingPolicy)
	router.HandleFunc("GET /scaling/autoscaler/events", h.GetScalingEvents)
	router.HandleFunc("GET /scaling/autoscaler/metrics", h.GetScalingMetrics)
	router.HandleFunc("GET /scaling/autoscaler/metrics/:region", h.GetRegionScalingMetrics)
}

// CreateScalingPolicyRequest represents a scaling policy creation request
type CreateScalingPolicyRequest struct {
	RegionName          string `json:"region_name"`
	Name                string `json:"name"`
	ScaleUpThreshold    int    `json:"scale_up_threshold"`
	ScaleDownThreshold  int    `json:"scale_down_threshold"`
	ScaleUpQuantity     int    `json:"scale_up_quantity"`
	ScaleDownQuantity   int    `json:"scale_down_quantity"`
	MinInstances        int    `json:"min_instances"`
	MaxInstances        int    `json:"max_instances"`
	CooldownPeriodSecs  int    `json:"cooldown_period_secs"`
	MetricsCheckIntervalSecs int `json:"metrics_check_interval_secs"`
}

// CreateScalingPolicy creates a new scaling policy
func (h *AutoScalerHandler) CreateScalingPolicy(w http.ResponseWriter, r *http.Request) {
	var req CreateScalingPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	policy := h.autoscaler.CreateScalingPolicy(
		req.RegionName, req.Name,
		req.ScaleUpThreshold, req.ScaleDownThreshold,
		req.ScaleUpQuantity, req.ScaleDownQuantity,
		req.MinInstances, req.MaxInstances,
		req.CooldownPeriodSecs, req.MetricsCheckIntervalSecs,
	)

	response := map[string]interface{}{
		"region_name":         policy.RegionName,
		"name":                policy.Name,
		"scale_up_threshold":  policy.ScaleUpThreshold,
		"scale_down_threshold": policy.ScaleDownThreshold,
		"min_instances":       policy.MinInstances,
		"max_instances":       policy.MaxInstances,
		"is_active":           policy.IsActive,
		"message":             "scaling policy created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListScalingPolicies returns all scaling policies
func (h *AutoScalerHandler) ListScalingPolicies(w http.ResponseWriter, r *http.Request) {
	policies := h.autoscaler.GetAllScalingPolicies()
	response := map[string]interface{}{
		"count":    len(policies),
		"policies": policies,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetScalingPolicy returns a specific scaling policy
func (h *AutoScalerHandler) GetScalingPolicy(w http.ResponseWriter, r *http.Request) {
	region := strings.TrimPrefix(r.URL.Path, "/scaling/autoscaler/policies/")
	policy := h.autoscaler.GetScalingPolicy(region)

	if policy == nil {
		http.Error(w, "policy not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

// UpdateScalingPolicyRequest represents a scaling policy update request
type UpdateScalingPolicyRequest struct {
	ScaleUpThreshold    *int `json:"scale_up_threshold,omitempty"`
	ScaleDownThreshold  *int `json:"scale_down_threshold,omitempty"`
	MinInstances        *int `json:"min_instances,omitempty"`
	MaxInstances        *int `json:"max_instances,omitempty"`
	IsActive            *bool `json:"is_active,omitempty"`
}

// UpdateScalingPolicy updates a scaling policy
func (h *AutoScalerHandler) UpdateScalingPolicy(w http.ResponseWriter, r *http.Request) {
	region := strings.TrimPrefix(r.URL.Path, "/scaling/autoscaler/policies/")
	var req UpdateScalingPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	policy := h.autoscaler.UpdateScalingPolicy(region, req.ScaleUpThreshold, req.ScaleDownThreshold, req.MinInstances, req.MaxInstances, req.IsActive)

	if policy == nil {
		http.Error(w, "policy not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"policy":  policy,
		"message": "scaling policy updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetScalingEvents returns recent scaling events
func (h *AutoScalerHandler) GetScalingEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	events := h.autoscaler.GetScalingEvents(limit)
	response := map[string]interface{}{
		"count":  len(events),
		"events": events,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetScalingMetrics returns current scaling metrics
func (h *AutoScalerHandler) GetScalingMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.autoscaler.GetScalingMetrics()
	response := map[string]interface{}{
		"metrics": metrics,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRegionScalingMetrics returns scaling metrics for a specific region
func (h *AutoScalerHandler) GetRegionScalingMetrics(w http.ResponseWriter, r *http.Request) {
	region := strings.TrimPrefix(r.URL.Path, "/scaling/autoscaler/metrics/")
	metrics := h.autoscaler.GetRegionScalingMetrics(region)

	if metrics == nil {
		http.Error(w, "region not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// HealthCheckerHandler handles health checker API requests
type HealthCheckerHandler struct {
	healthChecker *scaling.HealthChecker
}

// NewHealthCheckerHandler creates a new health checker handler
func NewHealthCheckerHandler(hc *scaling.HealthChecker) *HealthCheckerHandler {
	return &HealthCheckerHandler{healthChecker: hc}
}

// RegisterHealthCheckerRoutes registers health checker API routes
func (h *HealthCheckerHandler) RegisterHealthCheckerRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /scaling/health-checks/probes", h.RegisterHealthProbe)
	router.HandleFunc("GET /scaling/health-checks/probes", h.ListHealthProbes)
	router.HandleFunc("GET /scaling/health-checks/status", h.GetHealthStatus)
	router.HandleFunc("GET /scaling/health-checks/status/:server_id", h.GetServerHealthStatus)
	router.HandleFunc("GET /scaling/health-checks/history/:server_id", h.GetHealthHistory)
	router.HandleFunc("GET /scaling/health-checks/global", h.GetGlobalHealthStatus)
}

// RegisterHealthProbeRequest represents a health probe registration request
type RegisterHealthProbeRequest struct {
	ServerID       string `json:"server_id"`
	Endpoint       string `json:"endpoint"`
	Port           int    `json:"port"`
	Path           string `json:"path"`
	Method         string `json:"method"`
	ExpectedStatus int    `json:"expected_status"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

// RegisterHealthProbe registers a new health probe
func (h *HealthCheckerHandler) RegisterHealthProbe(w http.ResponseWriter, r *http.Request) {
	var req RegisterHealthProbeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}
	if req.ExpectedStatus == 0 {
		req.ExpectedStatus = 200
	}
	if req.TimeoutSeconds == 0 {
		req.TimeoutSeconds = 5
	}

	h.healthChecker.RegisterHealthProbe(
		req.ServerID, req.Endpoint, req.Port,
		req.Path, req.Method, req.ExpectedStatus,
	)

	response := map[string]interface{}{
		"server_id":       req.ServerID,
		"endpoint":        req.Endpoint,
		"path":            req.Path,
		"expected_status": req.ExpectedStatus,
		"message":         "health probe registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListHealthProbes returns all registered health probes
func (h *HealthCheckerHandler) ListHealthProbes(w http.ResponseWriter, r *http.Request) {
	probes := h.healthChecker.GetAllHealthProbes()
	response := map[string]interface{}{
		"count":  len(probes),
		"probes": probes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetHealthStatus returns current health status of all servers
func (h *HealthCheckerHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	status := h.healthChecker.GetHealthStatus()
	response := map[string]interface{}{
		"status": status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetServerHealthStatus returns health status of a specific server
func (h *HealthCheckerHandler) GetServerHealthStatus(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/scaling/health-checks/status/")
	status := h.healthChecker.GetServerHealthStatus(serverID)

	if status == nil {
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetHealthHistory returns recent health check history for a server
func (h *HealthCheckerHandler) GetHealthHistory(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/scaling/health-checks/history/")
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history := h.healthChecker.GetHealthHistory(serverID, limit)
	response := map[string]interface{}{
		"server_id": serverID,
		"count":     len(history),
		"history":   history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGlobalHealthStatus returns global health status summary
func (h *HealthCheckerHandler) GetGlobalHealthStatus(w http.ResponseWriter, r *http.Request) {
	status := h.healthChecker.GetGlobalHealthStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// FailoverHandler handles failover API requests
type FailoverHandler struct {
	failoverManager *scaling.FailoverManager
}

// NewFailoverHandler creates a new failover handler
func NewFailoverHandler(fm *scaling.FailoverManager) *FailoverHandler {
	return &FailoverHandler{failoverManager: fm}
}

// RegisterFailoverRoutes registers failover API routes
func (h *FailoverHandler) RegisterFailoverRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /scaling/failover/policies", h.CreateFailoverPolicy)
	router.HandleFunc("GET /scaling/failover/policies", h.ListFailoverPolicies)
	router.HandleFunc("GET /scaling/failover/status", h.GetFailoverStatus)
	router.HandleFunc("POST /scaling/failover/initiate", h.InitiateFailover)
	router.HandleFunc("POST /scaling/failover/failback", h.InitiateFailback)
	router.HandleFunc("GET /scaling/failover/events", h.GetFailoverEvents)
	router.HandleFunc("POST /scaling/disaster-recovery/plans", h.CreateDisasterRecoveryPlan)
	router.HandleFunc("GET /scaling/disaster-recovery/status", h.GetDisasterRecoveryStatus)
}

// CreateFailoverPolicyRequest represents a failover policy creation request
type CreateFailoverPolicyRequest struct {
	Name                    string `json:"name"`
	PrimaryRegion           string `json:"primary_region"`
	SecondaryRegion         string `json:"secondary_region"`
	HealthCheckThreshold    float64 `json:"health_check_threshold"`
	FailoverTimeoutSeconds  int    `json:"failover_timeout_seconds"`
	AutomaticFailover       bool   `json:"automatic_failover"`
}

// CreateFailoverPolicy creates a new failover policy
func (h *FailoverHandler) CreateFailoverPolicy(w http.ResponseWriter, r *http.Request) {
	var req CreateFailoverPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	policy := h.failoverManager.CreateFailoverPolicy(
		req.Name, req.PrimaryRegion, req.SecondaryRegion,
		req.HealthCheckThreshold, req.FailoverTimeoutSeconds,
		req.AutomaticFailover,
	)

	response := map[string]interface{}{
		"id":                   policy.ID,
		"name":                 policy.Name,
		"primary_region":       policy.PrimaryRegion,
		"secondary_region":     policy.SecondaryRegion,
		"automatic_failover":   policy.AutomaticFailover,
		"is_active":            policy.IsActive,
		"message":              "failover policy created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListFailoverPolicies returns all failover policies
func (h *FailoverHandler) ListFailoverPolicies(w http.ResponseWriter, r *http.Request) {
	policies := h.failoverManager.GetAllFailoverPolicies()
	response := map[string]interface{}{
		"count":    len(policies),
		"policies": policies,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFailoverStatus returns current failover status
func (h *FailoverHandler) GetFailoverStatus(w http.ResponseWriter, r *http.Request) {
	status := h.failoverManager.GetFailoverStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// InitiateFailoverRequest represents a manual failover request
type InitiateFailoverRequest struct {
	FromRegion string `json:"from_region"`
	ToRegion   string `json:"to_region"`
	Reason     string `json:"reason"`
}

// InitiateFailover initiates a manual failover
func (h *FailoverHandler) InitiateFailover(w http.ResponseWriter, r *http.Request) {
	var req InitiateFailoverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	event, err := h.failoverManager.InitiateFailover(req.FromRegion, req.ToRegion, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"event_id":    event.ID,
		"from_region": event.FromRegion,
		"to_region":   event.ToRegion,
		"status":      event.Status,
		"message":     "failover initiated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// InitiateFailbackRequest represents a failback request
type InitiateFailbackRequest struct {
	FromRegion string `json:"from_region"`
	ToRegion   string `json:"to_region"`
	Reason     string `json:"reason"`
}

// InitiateFailback initiates a failback to primary region
func (h *FailoverHandler) InitiateFailback(w http.ResponseWriter, r *http.Request) {
	var req InitiateFailbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	event, err := h.failoverManager.InitiateFailback(req.FromRegion, req.ToRegion, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"event_id":    event.ID,
		"from_region": event.FromRegion,
		"to_region":   event.ToRegion,
		"status":      event.Status,
		"message":     "failback initiated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFailoverEvents returns recent failover events
func (h *FailoverHandler) GetFailoverEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	events := h.failoverManager.GetFailoverHistory(limit)
	response := map[string]interface{}{
		"count":  len(events),
		"events": events,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateDisasterRecoveryPlanRequest represents a DR plan creation request
type CreateDisasterRecoveryPlanRequest struct {
	Name                    string   `json:"name"`
	RecoveryTimeObjective   int      `json:"recovery_time_objective_minutes"`
	RecoveryPointObjective  int      `json:"recovery_point_objective_minutes"`
	BackupFrequencyMinutes  int      `json:"backup_frequency_minutes"`
	BackupRetentionDays     int      `json:"backup_retention_days"`
	ReplicationEnabled      bool     `json:"replication_enabled"`
	ReplicationTargets      []string `json:"replication_targets"`
}

// CreateDisasterRecoveryPlan creates a new disaster recovery plan
func (h *FailoverHandler) CreateDisasterRecoveryPlan(w http.ResponseWriter, r *http.Request) {
	var req CreateDisasterRecoveryPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan := h.failoverManager.CreateDisasterRecoveryPlan(
		req.Name, req.RecoveryTimeObjective, req.RecoveryPointObjective,
		req.BackupFrequencyMinutes, req.BackupRetentionDays,
		req.ReplicationEnabled, req.ReplicationTargets,
	)

	response := map[string]interface{}{
		"id":                          plan.ID,
		"name":                        plan.Name,
		"recovery_time_objective":     plan.RecoveryTimeObjective,
		"recovery_point_objective":    plan.RecoveryPointObjective,
		"replication_enabled":         plan.ReplicationEnabled,
		"is_active":                   plan.IsActive,
		"message":                     "disaster recovery plan created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetDisasterRecoveryStatus returns current disaster recovery status
func (h *FailoverHandler) GetDisasterRecoveryStatus(w http.ResponseWriter, r *http.Request) {
	status := h.failoverManager.GetDisasterRecoveryStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// MultiRegionHandler handles multi-region deployment API requests
type MultiRegionHandler struct {
	multiRegion *scaling.MultiRegionDeployment
}

// NewMultiRegionHandler creates a new multi-region handler
func NewMultiRegionHandler(mrd *scaling.MultiRegionDeployment) *MultiRegionHandler {
	return &MultiRegionHandler{multiRegion: mrd}
}

// RegisterMultiRegionRoutes registers multi-region API routes
func (h *MultiRegionHandler) RegisterMultiRegionRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /scaling/multi-region/data-centers", h.RegisterDataCenter)
	router.HandleFunc("GET /scaling/multi-region/data-centers", h.ListDataCenters)
	router.HandleFunc("GET /scaling/multi-region/data-centers/:id", h.GetDataCenterStatus)
	router.HandleFunc("POST /scaling/multi-region/configs", h.CreateDeploymentConfig)
	router.HandleFunc("GET /scaling/multi-region/configs", h.ListDeploymentConfigs)
	router.HandleFunc("GET /scaling/multi-region/status/:config_name", h.GetMultiRegionStatus)
	router.HandleFunc("POST /scaling/multi-region/latency", h.UpdateLatencyMatrix)
	router.HandleFunc("GET /scaling/multi-region/latency-matrix", h.GetLatencyMatrix)
	router.HandleFunc("GET /scaling/multi-region/comparison", h.GetDeploymentComparison)
}

// RegisterDataCenterRequest represents a data center registration request
type RegisterDataCenterRequest struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Provider  string  `json:"provider"`
	Region    string  `json:"region"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Capacity  int     `json:"capacity"`
	Tier      string  `json:"tier"`
	IsPrimary bool    `json:"is_primary"`
}

// RegisterDataCenter registers a new data center
func (h *MultiRegionHandler) RegisterDataCenter(w http.ResponseWriter, r *http.Request) {
	var req RegisterDataCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dc := h.multiRegion.RegisterDataCenter(
		req.ID, req.Name, req.Location, req.Provider, req.Region,
		req.Latitude, req.Longitude, req.Capacity, req.Tier, req.IsPrimary,
	)

	response := map[string]interface{}{
		"id":         dc.ID,
		"name":       dc.Name,
		"location":   dc.Location,
		"provider":   dc.Provider,
		"region":     dc.Region,
		"capacity":   dc.Capacity,
		"is_primary": dc.IsPrimary,
		"tier":       dc.Tier,
		"message":    "data center registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListDataCenters returns all registered data centers
func (h *MultiRegionHandler) ListDataCenters(w http.ResponseWriter, r *http.Request) {
	dataCenters := h.multiRegion.GetAllDataCenters()
	response := map[string]interface{}{
		"count":         len(dataCenters),
		"data_centers": dataCenters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDataCenterStatus returns status of a specific data center
func (h *MultiRegionHandler) GetDataCenterStatus(w http.ResponseWriter, r *http.Request) {
	dcID := strings.TrimPrefix(r.URL.Path, "/scaling/multi-region/data-centers/")
	status := h.multiRegion.GetDataCenterStatus(dcID)

	if status == nil {
		http.Error(w, "data center not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// CreateDeploymentConfigRequest represents a deployment config creation request
type CreateDeploymentConfigRequest struct {
	Name                   string   `json:"name"`
	ActiveRegions          []string `json:"active_regions"`
	DataResidencyRegion    string   `json:"data_residency_region"`
	DisasterRecoveryRegion string   `json:"disaster_recovery_region"`
	FailoverStrategy       string   `json:"failover_strategy"`
	CostOptimization       bool     `json:"cost_optimization"`
	LatencyTarget          int      `json:"latency_target_ms"`
}

// CreateDeploymentConfig creates a new deployment configuration
func (h *MultiRegionHandler) CreateDeploymentConfig(w http.ResponseWriter, r *http.Request) {
	var req CreateDeploymentConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	config := h.multiRegion.CreateDeploymentConfig(
		req.Name, req.ActiveRegions, req.DataResidencyRegion,
		req.DisasterRecoveryRegion, req.FailoverStrategy,
		req.CostOptimization, req.LatencyTarget,
	)

	response := map[string]interface{}{
		"name":                    config.Name,
		"active_regions":          config.ActiveRegions,
		"failover_strategy":       config.FailoverStrategy,
		"data_residency_region":   config.DataResidencyRegion,
		"disaster_recovery_region": config.DisasterRecoveryRegion,
		"is_active":               config.IsActive,
		"message":                 "deployment config created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListDeploymentConfigs returns all deployment configurations
func (h *MultiRegionHandler) ListDeploymentConfigs(w http.ResponseWriter, r *http.Request) {
	configs := h.multiRegion.GetAllDeploymentConfigs()
	response := map[string]interface{}{
		"count":   len(configs),
		"configs": configs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMultiRegionStatus returns overall multi-region deployment status
func (h *MultiRegionHandler) GetMultiRegionStatus(w http.ResponseWriter, r *http.Request) {
	configName := strings.TrimPrefix(r.URL.Path, "/scaling/multi-region/status/")
	status := h.multiRegion.GetMultiRegionStatus(configName)

	if status == nil {
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// UpdateLatencyMatrixRequest represents a latency matrix update request
type UpdateLatencyMatrixRequest struct {
	FromRegion string `json:"from_region"`
	ToRegion   string `json:"to_region"`
	LatencyMs  int    `json:"latency_ms"`
}

// UpdateLatencyMatrix updates latency between regions
func (h *MultiRegionHandler) UpdateLatencyMatrix(w http.ResponseWriter, r *http.Request) {
	var req UpdateLatencyMatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.multiRegion.UpdateLatencyMatrix(req.FromRegion, req.ToRegion, req.LatencyMs)

	response := map[string]interface{}{
		"from_region": req.FromRegion,
		"to_region":   req.ToRegion,
		"latency_ms":  req.LatencyMs,
		"message":     "latency matrix updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetLatencyMatrix returns the latency matrix for all regions
func (h *MultiRegionHandler) GetLatencyMatrix(w http.ResponseWriter, r *http.Request) {
	matrix := h.multiRegion.GetLatencyMatrix()
	response := map[string]interface{}{
		"latency_matrix": matrix,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDeploymentComparison returns comparison of deployment strategies
func (h *MultiRegionHandler) GetDeploymentComparison(w http.ResponseWriter, r *http.Request) {
	comparison := h.multiRegion.GetDeploymentComparison()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparison)
}
