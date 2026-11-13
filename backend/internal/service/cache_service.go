package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheService provides multi-layer caching strategies
type CacheService struct {
	memoryCache     map[string]*CacheEntry
	invalidations   map[string]time.Time
	metrics         *CacheMetrics
	mu              sync.RWMutex
	ttlTicker       *time.Ticker
	maxSize         int
	cleanupInterval time.Duration
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
	Hits      int64
	Size      int64
	Tags      []string
	Priority  int // Higher = keep longer
	Version   string
}

// CacheMetrics tracks cache performance
type CacheMetrics struct {
	Hits        int64
	Misses      int64
	Evictions   int64
	Size        int64
	ItemCount   int64
	LastUpdated time.Time
}

// CacheInvalidationRule defines when cache should be invalidated
type CacheInvalidationRule struct {
	ID        string
	Pattern   string // Glob pattern (e.g., "components/*")
	Trigger   string // "time_based", "event_based", "manual"
	Interval  time.Duration
	LastFired time.Time
}

// CacheWarmingStrategy defines how to pre-populate cache
type CacheWarmingStrategy struct {
	ID          string
	Name        string
	Keys        []string
	Schedule    string // "on_startup", "hourly", "daily"
	Priority    int
	Description string
}

// NewCacheService creates a new cache service
func NewCacheService(maxSize int, cleanupInterval time.Duration) *CacheService {
	service := &CacheService{
		memoryCache:     make(map[string]*CacheEntry),
		invalidations:   make(map[string]time.Time),
		metrics:         &CacheMetrics{},
		ttlTicker:       time.NewTicker(cleanupInterval),
		maxSize:         maxSize,
		cleanupInterval: cleanupInterval,
	}

	// Start cleanup goroutine
	go service.cleanupExpired()

	return service
}

// Get retrieves a value from cache
func (cs *CacheService) Get(key string) (interface{}, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	entry, exists := cs.memoryCache[key]
	if !exists {
		cs.metrics.Misses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		go cs.Delete(key)
		cs.metrics.Misses++
		return nil, false
	}

	entry.Hits++
	cs.metrics.Hits++

	fmt.Printf("DEBUG: Cache HIT for key: %s (hits: %d)\n", key, entry.Hits)
	return entry.Value, true
}

// Set stores a value in cache
func (cs *CacheService) Set(key string, value interface{}, ttl time.Duration, tags ...string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Check if we need to evict
	if cs.getCurrentSize() >= int64(cs.maxSize) {
		cs.evictLRU()
	}

	hash := md5.Sum([]byte(key))
	size := int64(len(hex.EncodeToString(hash[:])))

	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		Hits:      0,
		Size:      size,
		Tags:      tags,
		Priority:  1,
		Version:   fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	cs.memoryCache[key] = entry
	cs.metrics.ItemCount = int64(len(cs.memoryCache))
	cs.metrics.Size = cs.getCurrentSize()

	fmt.Printf("DEBUG: Cache SET for key: %s (ttl: %v, tags: %v)\n", key, ttl, tags)
	return nil
}

// Delete removes a key from cache
func (cs *CacheService) Delete(key string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.memoryCache[key]; exists {
		delete(cs.memoryCache, key)
		cs.metrics.ItemCount = int64(len(cs.memoryCache))
		cs.metrics.Size = cs.getCurrentSize()
		fmt.Printf("DEBUG: Cache DELETE for key: %s\n", key)
	}

	return nil
}

// InvalidateByTag invalidates all entries with a specific tag
func (cs *CacheService) InvalidateByTag(tag string) int {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	count := 0
	for key, entry := range cs.memoryCache {
		for _, t := range entry.Tags {
			if t == tag {
				delete(cs.memoryCache, key)
				count++
				break
			}
		}
	}

	cs.metrics.ItemCount = int64(len(cs.memoryCache))
	cs.metrics.Size = cs.getCurrentSize()

	fmt.Printf("DEBUG: Cache invalidated %d entries with tag: %s\n", count, tag)
	return count
}

// InvalidateByPattern invalidates entries matching pattern
func (cs *CacheService) InvalidateByPattern(pattern string) int {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	count := 0
	for key := range cs.memoryCache {
		if matchesPattern(key, pattern) {
			delete(cs.memoryCache, key)
			count++
		}
	}

	cs.metrics.ItemCount = int64(len(cs.memoryCache))
	cs.metrics.Size = cs.getCurrentSize()

	fmt.Printf("DEBUG: Cache invalidated %d entries matching pattern: %s\n", count, pattern)
	return count
}

// Clear clears all cache
func (cs *CacheService) Clear() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.memoryCache = make(map[string]*CacheEntry)
	cs.metrics = &CacheMetrics{}

	fmt.Println("DEBUG: Cache cleared")
	return nil
}

// GetMetrics returns cache metrics
func (cs *CacheService) GetMetrics() *CacheMetrics {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	metrics := &CacheMetrics{
		Hits:        cs.metrics.Hits,
		Misses:      cs.metrics.Misses,
		Evictions:   cs.metrics.Evictions,
		Size:        cs.metrics.Size,
		ItemCount:   cs.metrics.ItemCount,
		LastUpdated: time.Now(),
	}

	return metrics
}

// GetHitRate returns cache hit rate
func (cs *CacheService) GetHitRate() float64 {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	total := cs.metrics.Hits + cs.metrics.Misses
	if total == 0 {
		return 0
	}

	return float64(cs.metrics.Hits) / float64(total) * 100
}

// WarmCache pre-populates cache with provided data
func (cs *CacheService) WarmCache(key string, value interface{}, ttl time.Duration, priority int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		Hits:      0,
		Priority:  priority,
		Version:   fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	cs.memoryCache[key] = entry
	cs.metrics.ItemCount = int64(len(cs.memoryCache))

	fmt.Printf("DEBUG: Cache warmed for key: %s (priority: %d)\n", key, priority)
	return nil
}

// GetEntry returns cache entry metadata
func (cs *CacheService) GetEntry(key string) (*CacheEntry, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	entry, exists := cs.memoryCache[key]
	return entry, exists
}

// Helper functions

func (cs *CacheService) getCurrentSize() int64 {
	size := int64(0)
	for _, entry := range cs.memoryCache {
		size += entry.Size
	}
	return size
}

func (cs *CacheService) cleanupExpired() {
	for range cs.ttlTicker.C {
		cs.mu.Lock()

		now := time.Now()
		for key, entry := range cs.memoryCache {
			if now.After(entry.ExpiresAt) {
				delete(cs.memoryCache, key)
				fmt.Printf("DEBUG: Cache entry expired and removed: %s\n", key)
			}
		}

		cs.metrics.ItemCount = int64(len(cs.memoryCache))
		cs.metrics.Size = cs.getCurrentSize()
		cs.metrics.LastUpdated = now

		cs.mu.Unlock()
	}
}

func (cs *CacheService) evictLRU() {
	// Find entry with lowest priority and oldest access
	var lruKey string
	var lruEntry *CacheEntry

	for key, entry := range cs.memoryCache {
		if lruEntry == nil || entry.Priority < lruEntry.Priority ||
			(entry.Priority == lruEntry.Priority && entry.Hits < lruEntry.Hits) {
			lruKey = key
			lruEntry = entry
		}
	}

	if lruKey != "" {
		delete(cs.memoryCache, lruKey)
		cs.metrics.Evictions++
		fmt.Printf("DEBUG: Cache entry evicted (LRU): %s\n", lruKey)
	}
}

func matchesPattern(key string, pattern string) bool {
	// Simple pattern matching (in production, use regex)
	// For now, supports * wildcard
	if pattern == "*" {
		return true
	}

	// Convert simple glob to comparison
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}

	return key == pattern
}

// Shutdown gracefully shuts down the cache service
func (cs *CacheService) Shutdown() error {
	cs.ttlTicker.Stop()
	fmt.Println("DEBUG: Cache service shut down")
	return nil
}
