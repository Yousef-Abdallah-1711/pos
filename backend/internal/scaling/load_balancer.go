package scaling

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Region represents a geographic region with servers
type Region struct {
	Name      string
	Continent string
	Latitude  float64
	Longitude float64
	Servers   []*Server
	Priority  int // Lower number = higher priority
	IsActive  bool
}

// Server represents a single application server
type Server struct {
	ID              string
	Region          *Region
	Endpoint        string // hostname or IP
	Port            int
	Weight          int // 1-100 for weighted round-robin
	Capacity        int // Max requests per second
	CurrentLoad     int // Current requests per second
	IsHealthy       bool
	ResponseTime    time.Duration
	ErrorRate       float64 // 0-1
	CPUUsage        float64 // 0-100
	MemoryUsage     float64 // 0-100
	LastHealthCheck time.Time
	CreatedAt       time.Time
}

// LoadBalancer implements global load balancing
type LoadBalancer struct {
	regions         map[string]*Region
	servers         map[string]*Server
	strategyName    string // "round-robin", "least-loaded", "fastest", "geo"
	currentIndex    int    // For round-robin
	mu              sync.RWMutex
	preferredRegion string            // User's closest region
	sessionAffinity bool              // Keep user on same server
	sessions        map[string]string // session_id -> server_id
	sessionTTL      time.Duration
	lastCleanup     time.Time
}

// LoadBalancingStrategy defines how requests are routed
type LoadBalancingStrategy string

const (
	StrategyRoundRobin  LoadBalancingStrategy = "round-robin"
	StrategyLeastLoaded LoadBalancingStrategy = "least-loaded"
	StrategyFastest     LoadBalancingStrategy = "fastest"
	StrategyGeo         LoadBalancingStrategy = "geo"
	StrategyWeighted    LoadBalancingStrategy = "weighted"
)

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(strategy LoadBalancingStrategy) *LoadBalancer {
	return &LoadBalancer{
		regions:         make(map[string]*Region),
		servers:         make(map[string]*Server),
		strategyName:    string(strategy),
		sessions:        make(map[string]string),
		sessionTTL:      30 * time.Minute,
		lastCleanup:     time.Now(),
		sessionAffinity: true,
	}
}

// RegisterRegion registers a new geographic region
func (lb *LoadBalancer) RegisterRegion(
	name string,
	continent string,
	latitude float64,
	longitude float64,
	priority int,
) *Region {
	log.Printf("[LB] Registering region: %s (priority: %d)", name, priority)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	region := &Region{
		Name:      name,
		Continent: continent,
		Latitude:  latitude,
		Longitude: longitude,
		Priority:  priority,
		IsActive:  true,
		Servers:   make([]*Server, 0),
	}

	lb.regions[name] = region
	log.Printf("[LB] Region registered: %s", name)
	return region
}

// RegisterServer registers a server in a region
func (lb *LoadBalancer) RegisterServer(
	serverID string,
	regionName string,
	endpoint string,
	port int,
	weight int,
	capacity int,
) (*Server, error) {
	log.Printf("[LB] Registering server: %s in region %s", serverID, regionName)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	region, exists := lb.regions[regionName]
	if !exists {
		return nil, fmt.Errorf("region %s not found", regionName)
	}

	server := &Server{
		ID:              serverID,
		Region:          region,
		Endpoint:        endpoint,
		Port:            port,
		Weight:          weight,
		Capacity:        capacity,
		IsHealthy:       true,
		CreatedAt:       time.Now(),
		LastHealthCheck: time.Now(),
	}

	lb.servers[serverID] = server
	region.Servers = append(region.Servers, server)

	log.Printf("[LB] Server registered: %s (capacity: %d req/s)", serverID, capacity)
	return server, nil
}

// SelectServer selects a server for a request
func (lb *LoadBalancer) SelectServer(
	clientIP string,
	sessionID string,
	preferredRegion string,
) (*Server, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Check session affinity
	if lb.sessionAffinity && sessionID != "" {
		if serverID, exists := lb.sessions[sessionID]; exists {
			if server, exists := lb.servers[serverID]; exists && server.IsHealthy {
				log.Printf("[LB] Using session affinity: server=%s", serverID)
				return server, nil
			}
			// Session expired or server unhealthy
			delete(lb.sessions, sessionID)
		}
	}

	// Get healthy servers
	healthyServers := lb.getHealthyServers(preferredRegion)
	if len(healthyServers) == 0 {
		return nil, fmt.Errorf("no healthy servers available")
	}

	var selected *Server

	// Select based on strategy
	switch lb.strategyName {
	case string(StrategyRoundRobin):
		selected = lb.selectRoundRobin(healthyServers)
	case string(StrategyLeastLoaded):
		selected = lb.selectLeastLoaded(healthyServers)
	case string(StrategyFastest):
		selected = lb.selectFastest(healthyServers)
	case string(StrategyGeo):
		selected = lb.selectGeo(healthyServers, clientIP)
	case string(StrategyWeighted):
		selected = lb.selectWeighted(healthyServers)
	default:
		selected = lb.selectRoundRobin(healthyServers)
	}

	// Store session affinity
	if lb.sessionAffinity && sessionID != "" {
		lb.sessions[sessionID] = selected.ID
	}

	log.Printf("[LB] Selected server: %s (strategy: %s)", selected.ID, lb.strategyName)
	return selected, nil
}

// selectRoundRobin selects next server in round-robin fashion
func (lb *LoadBalancer) selectRoundRobin(servers []*Server) *Server {
	server := servers[lb.currentIndex%len(servers)]
	lb.currentIndex++
	return server
}

// selectLeastLoaded selects server with least current load
func (lb *LoadBalancer) selectLeastLoaded(servers []*Server) *Server {
	var selected *Server
	minLoad := math.MaxInt32

	for _, server := range servers {
		utilization := (server.CurrentLoad * 100) / server.Capacity
		if utilization < minLoad {
			minLoad = utilization
			selected = server
		}
	}

	return selected
}

// selectFastest selects server with fastest response time
func (lb *LoadBalancer) selectFastest(servers []*Server) *Server {
	var selected *Server
	minTime := time.Hour

	for _, server := range servers {
		if server.ResponseTime < minTime {
			minTime = server.ResponseTime
			selected = server
		}
	}

	return selected
}

// selectGeo selects server closest to client (by geographic distance)
func (lb *LoadBalancer) selectGeo(servers []*Server, clientIP string) *Server {
	// In production, use real GeoIP database to get client location
	// For now, return first server or use random

	if len(servers) == 0 {
		return nil
	}

	// Group by region and prefer by priority
	regionMap := make(map[*Region][]*Server)
	for _, server := range servers {
		regionMap[server.Region] = append(regionMap[server.Region], server)
	}

	// Find region with highest priority (lowest number)
	var bestRegion *Region
	minPriority := math.MaxInt32

	for region := range regionMap {
		if region.Priority < minPriority {
			minPriority = region.Priority
			bestRegion = region
		}
	}

	if bestRegion != nil {
		regionServers := regionMap[bestRegion]
		return regionServers[rand.Intn(len(regionServers))]
	}

	return servers[0]
}

// selectWeighted selects server based on weight (weighted round-robin)
func (lb *LoadBalancer) selectWeighted(servers []*Server) *Server {
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	if totalWeight == 0 {
		return servers[rand.Intn(len(servers))]
	}

	random := rand.Intn(totalWeight)
	current := 0

	for _, server := range servers {
		current += server.Weight
		if random < current {
			return server
		}
	}

	return servers[0]
}

// getHealthyServers returns list of healthy servers in preferred region
func (lb *LoadBalancer) getHealthyServers(preferredRegion string) []*Server {
	var result []*Server

	// First try preferred region
	if preferredRegion != "" {
		if region, exists := lb.regions[preferredRegion]; exists {
			for _, server := range region.Servers {
				if server.IsHealthy {
					result = append(result, server)
				}
			}
		}
	}

	// If no healthy servers in preferred region, get from all regions
	if len(result) == 0 {
		for _, server := range lb.servers {
			if server.IsHealthy {
				result = append(result, server)
			}
		}
	}

	return result
}

// UpdateServerMetrics updates server performance metrics
func (lb *LoadBalancer) UpdateServerMetrics(
	serverID string,
	currentLoad int,
	responseTime time.Duration,
	errorRate float64,
	cpuUsage float64,
	memoryUsage float64,
) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server, exists := lb.servers[serverID]
	if !exists {
		return fmt.Errorf("server %s not found", serverID)
	}

	server.CurrentLoad = currentLoad
	server.ResponseTime = responseTime
	server.ErrorRate = errorRate
	server.CPUUsage = cpuUsage
	server.MemoryUsage = memoryUsage
	server.LastHealthCheck = time.Now()

	log.Printf("[LB] Updated metrics: %s (load: %d/%d, response: %v, error: %.2f%%)",
		serverID, currentLoad, server.Capacity, responseTime, errorRate*100)

	return nil
}

// SetServerHealth sets the health status of a server
func (lb *LoadBalancer) SetServerHealth(serverID string, isHealthy bool) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server, exists := lb.servers[serverID]
	if !exists {
		return fmt.Errorf("server %s not found", serverID)
	}

	server.IsHealthy = isHealthy
	status := "healthy"
	if !isHealthy {
		status = "unhealthy"
	}

	log.Printf("[LB] Server %s marked as %s", serverID, status)
	return nil
}

// GetServerStats returns statistics for a server
func (lb *LoadBalancer) GetServerStats(serverID string) map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	server, exists := lb.servers[serverID]
	if !exists {
		return nil
	}

	utilization := 0
	if server.Capacity > 0 {
		utilization = (server.CurrentLoad * 100) / server.Capacity
	}

	return map[string]interface{}{
		"id":                server.ID,
		"region":            server.Region.Name,
		"endpoint":          server.Endpoint,
		"port":              server.Port,
		"is_healthy":        server.IsHealthy,
		"current_load":      server.CurrentLoad,
		"capacity":          server.Capacity,
		"utilization":       utilization,
		"response_time_ms":  server.ResponseTime.Milliseconds(),
		"error_rate":        fmt.Sprintf("%.2f%%", server.ErrorRate*100),
		"cpu_usage":         fmt.Sprintf("%.1f%%", server.CPUUsage),
		"memory_usage":      fmt.Sprintf("%.1f%%", server.MemoryUsage),
		"last_health_check": server.LastHealthCheck,
		"created_at":        server.CreatedAt,
	}
}

// GetRegionStats returns statistics for a region
func (lb *LoadBalancer) GetRegionStats(regionName string) map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	region, exists := lb.regions[regionName]
	if !exists {
		return nil
	}

	healthyCount := 0
	totalLoad := 0
	totalCapacity := 0
	avgResponseTime := time.Duration(0)

	for _, server := range region.Servers {
		if server.IsHealthy {
			healthyCount++
		}
		totalLoad += server.CurrentLoad
		totalCapacity += server.Capacity
		avgResponseTime += server.ResponseTime
	}

	if len(region.Servers) > 0 {
		avgResponseTime /= time.Duration(len(region.Servers))
	}

	utilization := 0
	if totalCapacity > 0 {
		utilization = (totalLoad * 100) / totalCapacity
	}

	return map[string]interface{}{
		"name":            region.Name,
		"continent":       region.Continent,
		"is_active":       region.IsActive,
		"total_servers":   len(region.Servers),
		"healthy_servers": healthyCount,
		"total_load":      totalLoad,
		"total_capacity":  totalCapacity,
		"utilization":     utilization,
		"avg_response_ms": avgResponseTime.Milliseconds(),
	}
}

// GetGlobalStats returns global load balancing statistics
func (lb *LoadBalancer) GetGlobalStats() map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	healthyServers := 0
	totalLoad := 0
	totalCapacity := 0
	var avgResponseTime time.Duration

	for _, server := range lb.servers {
		if server.IsHealthy {
			healthyServers++
		}
		totalLoad += server.CurrentLoad
		totalCapacity += server.Capacity
		avgResponseTime += server.ResponseTime
	}

	if len(lb.servers) > 0 {
		avgResponseTime /= time.Duration(len(lb.servers))
	}

	utilization := 0
	if totalCapacity > 0 {
		utilization = (totalLoad * 100) / totalCapacity
	}

	return map[string]interface{}{
		"strategy":         lb.strategyName,
		"total_regions":    len(lb.regions),
		"total_servers":    len(lb.servers),
		"healthy_servers":  healthyServers,
		"total_load":       totalLoad,
		"total_capacity":   totalCapacity,
		"utilization":      utilization,
		"avg_response_ms":  avgResponseTime.Milliseconds(),
		"active_sessions":  len(lb.sessions),
		"session_affinity": lb.sessionAffinity,
	}
}

// SetStrategy changes the load balancing strategy
func (lb *LoadBalancer) SetStrategy(strategy LoadBalancingStrategy) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.strategyName = string(strategy)
	log.Printf("[LB] Load balancing strategy changed to: %s", strategy)
}

// CleanupExpiredSessions removes expired sessions
func (lb *LoadBalancer) CleanupExpiredSessions() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	removed := 0

	for sessionID := range lb.sessions {
		// In production, track session creation time
		// For now, remove periodically
		delete(lb.sessions, sessionID)
		removed++
	}

	lb.lastCleanup = now
	log.Printf("[LB] Cleaned up %d expired sessions", removed)
}

// GetEndpoint returns the full endpoint URL for a server
func (s *Server) GetEndpoint() string {
	return fmt.Sprintf("http://%s:%d", s.Endpoint, s.Port)
}

// GetHealthPercentage returns health percentage for a server
func (s *Server) GetHealthPercentage() int {
	health := 100

	// Reduce health based on error rate
	health -= int(s.ErrorRate * 100)

	// Reduce health based on high CPU usage
	if s.CPUUsage > 90 {
		health -= 10
	} else if s.CPUUsage > 80 {
		health -= 5
	}

	// Reduce health based on high memory usage
	if s.MemoryUsage > 90 {
		health -= 10
	} else if s.MemoryUsage > 80 {
		health -= 5
	}

	if health < 0 {
		health = 0
	}

	return health
}
