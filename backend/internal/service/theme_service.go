package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// Cache provides a simple in-memory cache implementation (can be replaced with Redis)
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// simpleCache is a basic in-memory cache implementation
type simpleCache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
}

// NewSimpleCache creates a new in-memory cache
func NewSimpleCache() Cache {
	cache := &simpleCache{
		entries: make(map[string]CacheEntry),
	}
	// Start cleanup goroutine to remove expired entries
	go cache.cleanupExpired()
	return cache
}

// Get retrieves a value from the cache
func (c *simpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

// Set stores a value in the cache with TTL
func (c *simpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *simpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

// Clear removes all entries from the cache
func (c *simpleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]CacheEntry)
}

// cleanupExpired periodically removes expired entries
func (c *simpleCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// ThemeService provides business logic for theme operations
type ThemeService interface {
	// CreateTheme creates a new theme
	CreateTheme(ctx context.Context, req *domain.ThemeCreateRequest, tenantID, restaurantID, createdBy int64) (*domain.Theme, error)

	// GetTheme retrieves a theme by ID with caching
	GetTheme(ctx context.Context, tenantID, restaurantID, id int64) (*domain.Theme, error)

	// UpdateTheme updates an existing theme
	UpdateTheme(ctx context.Context, tenantID, restaurantID, id int64, req *domain.ThemeUpdateRequest, updatedBy int64) (*domain.Theme, error)

	// DeleteTheme deletes a theme
	DeleteTheme(ctx context.Context, tenantID, restaurantID, id int64) error

	// ListThemes retrieves a paginated list of themes
	ListThemes(ctx context.Context, tenantID, restaurantID int64, page, pageSize int) (*domain.ThemeListResponse, error)

	// ListThemePresets retrieves pre-built theme templates
	ListThemePresets(ctx context.Context, page, pageSize int) (map[string]interface{}, error)

	// GetPreset retrieves a single preset by ID
	GetPreset(ctx context.Context, id int64) (*domain.Theme, error)

	// ActivateTheme activates a theme for a tenant
	ActivateTheme(ctx context.Context, tenantID, restaurantID, id int64) error
}

// themeServiceImpl implements ThemeService
type themeServiceImpl struct {
	repo  repository.ThemeRepository
	cache Cache
}

// NewThemeService creates a new instance of ThemeService
func NewThemeService(repo repository.ThemeRepository) ThemeService {
	return &themeServiceImpl{
		repo:  repo,
		cache: NewSimpleCache(),
	}
}

// NewThemeServiceWithCache creates a new ThemeService with a provided cache
func NewThemeServiceWithCache(repo repository.ThemeRepository, cache Cache) ThemeService {
	return &themeServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

// CreateTheme creates a new theme
func (s *themeServiceImpl) CreateTheme(ctx context.Context, req *domain.ThemeCreateRequest, tenantID, restaurantID, createdBy int64) (*domain.Theme, error) {
	// Log the operation
	log.Printf("[ThemeService] Creating theme: %s (slug: %s) for tenant: %d, restaurant: %d", req.Name, req.Slug, tenantID, restaurantID)

	// Validate request
	if req.Name == "" || req.Slug == "" {
		return nil, errors.New("theme name and slug are required")
	}

	if tenantID == 0 {
		return nil, errors.New("tenant ID is required")
	}

	if restaurantID == 0 {
		return nil, errors.New("restaurant ID is required")
	}

	// Create theme domain object
	createdByPtr := &createdBy
	theme := &domain.Theme{
		TenantID:     tenantID,
		RestaurantID: restaurantID,
		Name:         req.Name,
		Slug:         req.Slug,
		Config:       req.Config,
		Description:  req.Description,
		IsActive:     false,
		IsPublished:  false,
		Version:      1,
		CreatedBy:    createdByPtr,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert into database
	if err := s.repo.Create(ctx, theme); err != nil {
		log.Printf("[ThemeService] Error creating theme: %v", err)
		return nil, err
	}

	// Invalidate list caches
	s.invalidateListCaches(tenantID, restaurantID)

	log.Printf("[ThemeService] Theme created successfully: ID=%d", theme.ID)
	return theme, nil
}

// GetTheme retrieves a theme by ID with caching
func (s *themeServiceImpl) GetTheme(ctx context.Context, tenantID, restaurantID, id int64) (*domain.Theme, error) {
	cacheKey := fmt.Sprintf("theme:%d:%d:%d", tenantID, restaurantID, id)

	// Check cache first
	if cached, found := s.cache.Get(cacheKey); found {
		log.Printf("[ThemeService] Cache HIT for theme: %d", id)
		return cached.(*domain.Theme), nil
	}

	log.Printf("[ThemeService] Cache MISS for theme: %d, fetching from database", id)

	// Fetch from database
	theme, err := s.repo.GetByID(ctx, tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[ThemeService] Error fetching theme: %v", err)
		return nil, err
	}

	// Cache the result for 10 minutes
	s.cache.Set(cacheKey, theme, 10*time.Minute)

	return theme, nil
}

// UpdateTheme updates an existing theme
func (s *themeServiceImpl) UpdateTheme(ctx context.Context, tenantID, restaurantID, id int64, req *domain.ThemeUpdateRequest, updatedBy int64) (*domain.Theme, error) {
	log.Printf("[ThemeService] Updating theme: %d", id)

	// Fetch current theme
	theme, err := s.repo.GetByID(ctx, tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[ThemeService] Error fetching theme for update: %v", err)
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		theme.Name = req.Name
	}
	if req.Description != "" {
		theme.Description = req.Description
	}
	if req.Config != nil && len(req.Config) > 0 {
		theme.Config = req.Config
	}
	if req.IsActive != nil {
		theme.IsActive = *req.IsActive
	}
	if req.IsPublished != nil {
		theme.IsPublished = *req.IsPublished
		if *req.IsPublished {
			now := time.Now()
			theme.PublishedAt = &now
		}
	}

	updatedByPtr := &updatedBy
	theme.UpdatedBy = updatedByPtr
	theme.UpdatedAt = time.Now()

	// Update in database
	if err := s.repo.Update(ctx, theme); err != nil {
		log.Printf("[ThemeService] Error updating theme: %v", err)
		return nil, err
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("theme:%d:%d:%d", tenantID, restaurantID, id)
	s.cache.Delete(cacheKey)
	cacheKeyActive := fmt.Sprintf("theme:active:%d:%d", tenantID, restaurantID)
	s.cache.Delete(cacheKeyActive)
	s.invalidateListCaches(tenantID, restaurantID)

	log.Printf("[ThemeService] Theme updated successfully: ID=%d, version=%d", id, theme.Version)
	return theme, nil
}

// DeleteTheme deletes a theme
func (s *themeServiceImpl) DeleteTheme(ctx context.Context, tenantID, restaurantID, id int64) error {
	log.Printf("[ThemeService] Deleting theme: %d", id)

	// Fetch theme to get tenant_id before deletion
	_, err := s.repo.GetByID(ctx, tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[ThemeService] Error fetching theme for deletion: %v", err)
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tenantID, restaurantID, id); err != nil {
		log.Printf("[ThemeService] Error deleting theme: %v", err)
		return err
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("theme:%d:%d:%d", tenantID, restaurantID, id)
	s.cache.Delete(cacheKey)
	cacheKeyActive := fmt.Sprintf("theme:active:%d:%d", tenantID, restaurantID)
	s.cache.Delete(cacheKeyActive)
	s.invalidateListCaches(tenantID, restaurantID)

	log.Printf("[ThemeService] Theme deleted successfully: ID=%d", id)
	return nil
}

// ListThemes retrieves a paginated list of themes
func (s *themeServiceImpl) ListThemes(ctx context.Context, tenantID, restaurantID int64, page, pageSize int) (*domain.ThemeListResponse, error) {
	log.Printf("[ThemeService] Listing themes for tenant: %d, page: %d, pageSize: %d", tenantID, page, pageSize)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	cacheKey := fmt.Sprintf("theme:list:%d:%d:%d:%d", tenantID, restaurantID, page, pageSize)

	// Check cache first
	if cached, found := s.cache.Get(cacheKey); found {
		log.Printf("[ThemeService] Cache HIT for theme list: tenant=%d, page=%d", tenantID, page)
		return cached.(*domain.ThemeListResponse), nil
	}

	log.Printf("[ThemeService] Cache MISS for theme list, fetching from database")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch from database
	themes, total, err := s.repo.List(ctx, tenantID, restaurantID, offset, pageSize)
	if err != nil {
		log.Printf("[ThemeService] Error listing themes: %v", err)
		return nil, err
	}

	// Build response
	totalPages := (int(total) + pageSize - 1) / pageSize
	response := &domain.ThemeListResponse{
		Themes:     themes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// Cache the result for 5 minutes
	s.cache.Set(cacheKey, response, 5*time.Minute)

	return response, nil
}

// ListThemePresets retrieves pre-built theme templates from theme_presets table
func (s *themeServiceImpl) ListThemePresets(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	log.Printf("[ThemeService] Listing theme presets, page: %d, pageSize: %d", page, pageSize)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	cacheKey := fmt.Sprintf("theme:presets:%d:%d", page, pageSize)

	// Check cache first
	if cached, found := s.cache.Get(cacheKey); found {
		log.Printf("[ThemeService] Cache HIT for theme presets")
		return cached.(map[string]interface{}), nil
	}

	log.Printf("[ThemeService] Cache MISS for theme presets, fetching from database")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch presets from repository
	presets, total, err := s.repo.ListPresets(ctx, offset, pageSize)
	if err != nil {
		log.Printf("[ThemeService] Error listing presets: %v", err)
		return nil, err
	}

	// Build response
	totalPages := (int(total) + pageSize - 1) / pageSize
	response := map[string]interface{}{
		"data": presets,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      pageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	log.Printf("[ThemeService] Fetched %d presets out of %d total", len(presets), total)

	// Cache for 30 minutes (presets don't change often)
	s.cache.Set(cacheKey, response, 30*time.Minute)

	return response, nil
}

// GetPreset retrieves a single preset by ID
func (s *themeServiceImpl) GetPreset(ctx context.Context, id int64) (*domain.Theme, error) {
	log.Printf("[ThemeService] Getting preset: %d", id)

	cacheKey := fmt.Sprintf("theme:preset:%d", id)

	// Check cache first
	if cached, found := s.cache.Get(cacheKey); found {
		log.Printf("[ThemeService] Cache HIT for preset %d", id)
		return cached.(*domain.Theme), nil
	}

	log.Printf("[ThemeService] Cache MISS for preset %d, fetching from database", id)

	// Fetch preset from repository
	preset, err := s.repo.GetPreset(ctx, id)
	if err != nil {
		log.Printf("[ThemeService] Error getting preset: %v", err)
		return nil, err
	}

	// Cache for 1 hour (presets don't change often)
	s.cache.Set(cacheKey, preset, 1*time.Hour)

	return preset, nil
}

// ActivateTheme activates a theme for a tenant
func (s *themeServiceImpl) ActivateTheme(ctx context.Context, tenantID, restaurantID, id int64) error {
	log.Printf("[ThemeService] Activating theme: %d", id)

	// Fetch theme to get tenant_id
	_, err := s.repo.GetByID(ctx, tenantID, restaurantID, id)
	if err != nil {
		log.Printf("[ThemeService] Error fetching theme for activation: %v", err)
		return err
	}

	// Activation is isolated to the selected restaurant or tenant-wide theme.
	if err := s.repo.SetActive(ctx, tenantID, restaurantID, id); err != nil {
		log.Printf("[ThemeService] Error activating theme: %v", err)
		return err
	}

	// Invalidate caches
	cacheKeyActive := fmt.Sprintf("theme:active:%d:%d", tenantID, restaurantID)
	s.cache.Delete(cacheKeyActive)
	s.invalidateListCaches(tenantID, restaurantID)

	log.Printf("[ThemeService] Theme activated successfully: ID=%d", id)
	return nil
}

// invalidateListCaches removes all cached list entries for a tenant
func (s *themeServiceImpl) invalidateListCaches(tenantID, restaurantID int64) {
	// Clear all list cache entries for this tenant
	// We need to invalidate all pagination variations
	log.Printf("[ThemeService] Invalidating list caches for tenant: %d", tenantID)

	// Since we don't know all the pagination variations, we'll clear the main ones
	// Common page sizes: 10, 20, 50, 100
	// Common pages: 1-10
	pageSizes := []int{10, 20, 50, 100}
	for page := 1; page <= 10; page++ {
		for _, pageSize := range pageSizes {
			cacheKey := fmt.Sprintf("theme:list:%d:%d:%d:%d", tenantID, restaurantID, page, pageSize)
			s.cache.Delete(cacheKey)
		}
	}

	// Also delete any cached active theme for this tenant
	cacheKeyActive := fmt.Sprintf("theme:active:%d:%d", tenantID, restaurantID)
	s.cache.Delete(cacheKeyActive)

	// For production, consider:
	// 1. Using Redis with pattern matching (KEYS pattern)
	// 2. Keeping a set of all cache keys per tenant
	// 3. Using TTL-based expiration with event-driven cache invalidation
}
