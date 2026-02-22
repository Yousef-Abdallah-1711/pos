//go:build legacy
// +build legacy

package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pos-saas/internal/service"
)

// SearchHandler handles search-related endpoints
type SearchHandler struct {
	searchService *service.ElasticsearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *service.ElasticsearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// RegisterRoutes registers all search routes
func (h *SearchHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/search")
	{
		// Global search
		group.POST("/", h.Search)
		group.GET("/", h.SearchQuery)

		// Autocomplete
		group.GET("/autocomplete", h.GetAutocomplete)

		// Facets
		group.GET("/facets/:facetName", h.GetFacetValues)

		// Trending
		group.GET("/trending", h.GetTrending)

		// Related
		group.GET("/:componentID/related", h.GetRelated)

		// Saved searches
		group.POST("/saved", h.SaveSearch)
		group.GET("/saved", h.GetSavedSearches)
		group.GET("/saved/:savedID", h.ExecuteSavedSearch)
		group.DELETE("/saved/:savedID", h.DeleteSavedSearch)

		// Analytics
		group.GET("/analytics", h.GetSearchAnalytics)

		// Query analysis
		group.POST("/analyze", h.AnalyzeQuery)
	}
}

// Search performs a full-text search
// POST /api/v1/search
// Body: { "query": "...", "queryType": "...", "filters": [], "page": 1, "pageSize": 20 }
func (h *SearchHandler) Search(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.QueryType == "" {
		req.QueryType = "all"
	}

	// Perform search
	resp, err := h.searchService.Search(c.Request.Context(), restaurantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchQuery performs a search via query parameters
// GET /api/v1/search?q=button&type=component&page=1&size=20
func (h *SearchHandler) SearchQuery(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	queryType := c.DefaultQuery("type", "all")
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if sizeStr := c.Query("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	// Parse sort
	sortBy := c.DefaultQuery("sort", "relevance") // relevance, recent, popular, rating

	req := &models.SearchRequest{
		Query:    query,
		QueryType: queryType,
		SortBy:   sortBy,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.searchService.Search(c.Request.Context(), restaurantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAutocomplete returns autocomplete suggestions
// GET /api/v1/search/autocomplete?prefix=but&limit=10
func (h *SearchHandler) GetAutocomplete(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	prefix := c.Query("prefix")
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prefix parameter is required"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	suggestions, err := h.searchService.GetAutocompleteSuggestions(
		c.Request.Context(),
		restaurantID,
		prefix,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prefix":       prefix,
		"suggestions":  suggestions,
		"count":        len(suggestions),
	})
}

// GetFacetValues returns values for a specific facet
// GET /api/v1/search/facets/category?query=&limit=20
func (h *SearchHandler) GetFacetValues(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	facetName := c.Param("facetName")
	query := c.DefaultQuery("query", "")
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	values, err := h.searchService.GetFacetValues(c.Request.Context(), restaurantID, facetName, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"facet":  facetName,
		"values": values,
		"count":  len(values),
	})
}

// GetTrending returns trending search items
// GET /api/v1/search/trending?period=weekly&limit=10
func (h *SearchHandler) GetTrending(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	period := c.DefaultQuery("period", "daily")
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	trending, err := h.searchService.GetTrendingItems(c.Request.Context(), restaurantID, period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period":   period,
		"trending": trending,
		"count":    len(trending),
	})
}

// GetRelated returns components related to a given component
// GET /api/v1/search/:componentID/related?limit=10
func (h *SearchHandler) GetRelated(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	componentID, err := strconv.ParseInt(c.Param("componentID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid component_id"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	related, err := h.searchService.SearchRelated(c.Request.Context(), restaurantID, componentID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"component_id": componentID,
		"related":      related,
		"count":        len(related),
	})
}

// SaveSearch saves a search query
// POST /api/v1/search/saved
// Body: { "name": "...", "query": "...", "filters": [], "sortBy": "..." }
func (h *SearchHandler) SaveSearch(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	var req struct {
		Name   string                  `json:"name" binding:"required"`
		Query  string                  `json:"query" binding:"required"`
		Filters []*models.SearchFilter  `json:"filters"`
		SortBy string                  `json:"sort_by"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := h.searchService.SaveSearch(
		c.Request.Context(),
		restaurantID,
		userID,
		req.Name,
		req.Query,
		req.Filters,
		req.SortBy,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, saved)
}

// GetSavedSearches returns saved searches for the current user
// GET /api/v1/search/saved
func (h *SearchHandler) GetSavedSearches(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")
	userID := c.GetInt64("user_id")

	saved, err := h.searchService.GetSavedSearches(c.Request.Context(), restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"saved": saved,
		"count": len(saved),
	})
}

// ExecuteSavedSearch executes a saved search
// GET /api/v1/search/saved/:savedID
func (h *SearchHandler) ExecuteSavedSearch(c *gin.Context) {
	savedID, err := strconv.ParseInt(c.Param("savedID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid saved_id"})
		return
	}

	resp, err := h.searchService.ExecuteSavedSearch(c.Request.Context(), savedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteSavedSearch deletes a saved search
// DELETE /api/v1/search/saved/:savedID
func (h *SearchHandler) DeleteSavedSearch(c *gin.Context) {
	savedID, err := strconv.ParseInt(c.Param("savedID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid saved_id"})
		return
	}

	// In production, would delete from database
	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
		"saved_id": savedID,
	})
}

// GetSearchAnalytics returns search analytics for a period
// GET /api/v1/search/analytics?period=daily&start=2024-01-01&end=2024-01-31
func (h *SearchHandler) GetSearchAnalytics(c *gin.Context) {
	restaurantID := c.GetInt64("restaurant_id")

	period := c.DefaultQuery("period", "daily")
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	if startStr := c.Query("start"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			start = t
		}
	}

	if endStr := c.Query("end"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			end = t
		}
	}

	analytics, err := h.searchService.GetSearchAnalytics(c.Request.Context(), restaurantID, period, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// AnalyzeQuery analyzes a search query
// POST /api/v1/search/analyze
// Body: { "query": "..." }
func (h *SearchHandler) AnalyzeQuery(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis := h.searchService.AnalyzeQuery(req.Query)

	c.JSON(http.StatusOK, analysis)
}
