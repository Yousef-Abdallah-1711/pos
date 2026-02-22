//go:build legacy
// +build legacy

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"pos-saas/internal/repository"
)

// ElasticsearchService handles search operations using Elasticsearch
type ElasticsearchService struct {
	// In production, would use official Elasticsearch Go client
	// github.com/elastic/go-elasticsearch/v8
	searchRepo   repository.SearchRepository
	cacheService *CacheService
	config       *models.SearchConfig
}

// NewElasticsearchService creates a new search service
func NewElasticsearchService(
	searchRepo repository.SearchRepository,
	cacheService *CacheService,
) *ElasticsearchService {
	return &ElasticsearchService{
		searchRepo:   searchRepo,
		cacheService: cacheService,
	}
}

// Search performs a full-text search with optional filters and facets
func (es *ElasticsearchService) Search(
	ctx context.Context,
	restaurantID int64,
	req *models.SearchRequest,
) (*models.SearchResponse, error) {
	start := time.Now()

	// Check cache
	cacheKey := fmt.Sprintf("search:%d:%s:%v", restaurantID, req.Query, req.Filters)
	if cached, ok := es.cacheService.Get(cacheKey); ok {
		if resp, ok := cached.(*models.SearchResponse); ok {
			return resp, nil
		}
	}

	// Build Elasticsearch query
	query := es.buildElasticsearchQuery(req)

	// Execute search
	results, facets, aggs, totalHits, err := es.executeSearch(
		ctx,
		restaurantID,
		req.QueryType,
		query,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Get facets
	facetData := es.buildFacets(ctx, restaurantID, req.QueryType, facets)

	// Get suggestions
	suggestions := es.getSuggestions(ctx, req.Query)

	// Calculate pagination
	totalPages := int((totalHits + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := &models.SearchResponse{
		Query:         req.Query,
		TotalHits:     totalHits,
		Results:       results,
		Facets:        facetData,
		Aggregations:  aggs,
		Suggestions:   suggestions,
		ExecutionTime: time.Since(start).Milliseconds(),
		Page:          req.Page,
		PageSize:      req.PageSize,
		TotalPages:    totalPages,
	}

	// Cache for 1 hour
	es.cacheService.Set(cacheKey, response, 3600, []string{"search", fmt.Sprintf("restaurant:%d", restaurantID)})

	// Record search analytics
	_ = es.recordSearch(ctx, restaurantID, req, len(results))

	return response, nil
}

// buildElasticsearchQuery builds the Elasticsearch query DSL
func (es *ElasticsearchService) buildElasticsearchQuery(req *models.SearchRequest) map[string]interface{} {
	// Multi-match query for better relevance
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"multi_match": map[string]interface{}{
						"query":  req.Query,
						"fields": []string{"title^2", "description", "tags", "category"},
						"type":   "best_fields",
						"operator": "or",
						"fuzziness": "AUTO",
					},
				},
			},
		},
	}

	// Add filters
	if len(req.Filters) > 0 {
		filterClauses := make([]map[string]interface{}, len(req.Filters))
		for i, f := range req.Filters {
			filterClauses[i] = es.buildFilterClause(f)
		}
		query["bool"].(map[string]interface{})["filter"] = filterClauses
	}

	return query
}

// buildFilterClause builds a single filter clause
func (es *ElasticsearchService) buildFilterClause(filter *models.SearchFilter) map[string]interface{} {
	switch filter.Operator {
	case "equals":
		return map[string]interface{}{"term": map[string]interface{}{filter.Field: filter.Value}}
	case "contains":
		return map[string]interface{}{"match": map[string]interface{}{filter.Field: filter.Value}}
	case "in":
		return map[string]interface{}{"terms": map[string]interface{}{filter.Field: filter.Value}}
	case "range":
		return map[string]interface{}{"range": map[string]interface{}{filter.Field: filter.Value}}
	case "gt":
		return map[string]interface{}{"range": map[string]interface{}{filter.Field: map[string]interface{}{"gte": filter.Value}}}
	case "lt":
		return map[string]interface{}{"range": map[string]interface{}{filter.Field: map[string]interface{}{"lte": filter.Value}}}
	default:
		return map[string]interface{}{"term": map[string]interface{}{filter.Field: filter.Value}}
	}
}

// executeSearch executes the search against Elasticsearch
func (es *ElasticsearchService) executeSearch(
	ctx context.Context,
	restaurantID int64,
	queryType string,
	query map[string]interface{},
	page int,
	pageSize int,
) ([]*models.SearchResult, map[string]interface{}, *models.SearchAggregation, int64, error) {
	// Get index name for query type
	indexName := es.getIndexName(queryType)

	// Build search request
	from := (page - 1) * pageSize
	searchRequest := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  pageSize,
		"sort": []map[string]interface{}{
			{"_score": "desc"},
			{"created_at": "desc"},
		},
		"aggs": es.buildAggregations(),
	}

	// Execute query (in production, would call Elasticsearch)
	// For now, simulating with mock data
	results, totalHits, facets, aggs := es.mockSearchResults(restaurantID, queryType, pageSize)

	return results, facets, aggs, totalHits, nil
}

// mockSearchResults simulates search results
func (es *ElasticsearchService) mockSearchResults(
	restaurantID int64,
	queryType string,
	pageSize int,
) ([]*models.SearchResult, int64, map[string]interface{}, *models.SearchAggregation) {
	results := make([]*models.SearchResult, 0)

	// Generate mock results
	for i := 0; i < pageSize; i++ {
		result := &models.SearchResult{
			ID:           int64(i + 1),
			ResourceID:   int64(100 + i),
			ResourceType: queryType,
			Title:        fmt.Sprintf("Result %d", i+1),
			Description:  fmt.Sprintf("Description for result %d", i+1),
			Score:        0.95 - float64(i)*0.05,
			Rank:         i + 1,
			Rating:       4.5 - float64(i%2)*0.5,
			Views:        int64(1000 - i*50),
			CreatedAt:    time.Now().AddDate(0, 0, -i),
			UpdatedAt:    time.Now().AddDate(0, 0, -i),
		}
		results = append(results, result)
	}

	// Calculate facets
	facets := map[string]interface{}{
		"category": map[string]int{
			"component": 150,
			"menu":      85,
			"layout":    42,
		},
		"rating": map[string]int{
			"5": 45,
			"4": 120,
			"3": 67,
		},
	}

	// Build aggregations
	aggs := &models.SearchAggregation{
		ResourceType: map[string]int{
			"component": 200,
			"menu":      150,
		},
		Categories: map[string]int{
			"UI Component": 85,
			"Form":         65,
			"Navigation":   50,
		},
	}

	return results, 245, facets, aggs
}

// buildAggregations builds Elasticsearch aggregations for faceting
func (es *ElasticsearchService) buildAggregations() map[string]interface{} {
	return map[string]interface{}{
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category",
				"size":  20,
			},
		},
		"rating": map[string]interface{}{
			"histogram": map[string]interface{}{
				"field":    "rating",
				"interval": 1,
			},
		},
		"date_range": map[string]interface{}{
			"date_histogram": map[string]interface{}{
				"field":    "created_at",
				"interval": "month",
			},
		},
	}
}

// buildFacets builds facet data for response
func (es *ElasticsearchService) buildFacets(
	ctx context.Context,
	restaurantID int64,
	queryType string,
	facetData map[string]interface{},
) []*models.Facet {
	facets := make([]*models.Facet, 0)

	// Category facet
	facets = append(facets, &models.Facet{
		Name:  "category",
		Label: "Category",
		Type:  "string",
		Values: []*models.FacetValue{
			{Value: "component", Label: "Components", Count: 150},
			{Value: "menu", Label: "Menus", Count: 85},
			{Value: "layout", Label: "Layouts", Count: 42},
		},
	})

	// Rating facet
	facets = append(facets, &models.Facet{
		Name:  "rating",
		Label: "Rating",
		Type:  "numeric",
		Values: []*models.FacetValue{
			{Value: "5", Label: "5 Stars", Count: 120},
			{Value: "4", Label: "4+ Stars", Count: 85},
			{Value: "3", Label: "3+ Stars", Count: 40},
		},
	})

	// Date range facet
	facets = append(facets, &models.Facet{
		Name:  "date_range",
		Label: "Date",
		Type:  "date",
		Values: []*models.FacetValue{
			{Value: "last_week", Label: "Last Week", Count: 45},
			{Value: "last_month", Label: "Last Month", Count: 120},
			{Value: "last_year", Label: "Last Year", Count: 80},
		},
	})

	return facets
}

// GetAutocompleteSuggestions returns autocomplete suggestions
func (es *ElasticsearchService) GetAutocompleteSuggestions(
	ctx context.Context,
	restaurantID int64,
	prefix string,
	limit int,
) ([]string, error) {
	// Check cache
	cacheKey := fmt.Sprintf("autocomplete:%d:%s", restaurantID, prefix)
	if cached, ok := es.cacheService.Get(cacheKey); ok {
		if suggestions, ok := cached.([]string); ok {
			return suggestions, nil
		}
	}

	// Get suggestions from database
	suggestions, err := es.searchRepo.GetAutocompleteSuggestions(ctx, restaurantID, prefix, limit)
	if err != nil {
		return nil, err
	}

	// Cache for 24 hours
	es.cacheService.Set(cacheKey, suggestions, 86400, []string{"autocomplete"})

	return suggestions, nil
}

// GetTrendingItems returns trending search items
func (es *ElasticsearchService) GetTrendingItems(
	ctx context.Context,
	restaurantID int64,
	period string,
	limit int,
) ([]*models.TrendingComponent, error) {
	// Get trending items from database
	trending, err := es.searchRepo.GetTrendingComponents(ctx, restaurantID, period, limit)
	if err != nil {
		return nil, err
	}

	// Calculate trends
	for _, item := range trending {
		item.Trend = es.calculateTrend(item)
		item.TrendDirection = es.determineTrendDirection(item.Trend)
	}

	return trending, nil
}

// calculateTrend calculates trend score
func (es *ElasticsearchService) calculateTrend(component *models.TrendingComponent) float64 {
	// Simple trend: compare current vs previous period
	// In production, would fetch actual previous period data
	return 1.05 // 5% growth
}

// determineTrendDirection determines if trend is rising, falling, or stable
func (es *ElasticsearchService) determineTrendDirection(trend float64) string {
	if trend > 1.1 {
		return "rising"
	}
	if trend < 0.9 {
		return "falling"
	}
	return "stable"
}

// SearchRelated finds related items to a given component
func (es *ElasticsearchService) SearchRelated(
	ctx context.Context,
	restaurantID int64,
	componentID int64,
	limit int,
) ([]*models.SearchResult, error) {
	// Build more-like-this query based on component attributes
	// In production, would use Elasticsearch's more_like_this query

	// For now, return mock related items
	results := make([]*models.SearchResult, 0)
	for i := 0; i < limit; i++ {
		result := &models.SearchResult{
			ID:           int64(1000 + i),
			ResourceID:   int64(100 + i),
			ResourceType: "component",
			Title:        fmt.Sprintf("Related Component %d", i+1),
			Description:  fmt.Sprintf("Related to component %d", componentID),
			Score:        0.8 - float64(i)*0.05,
			Rating:       4.5,
		}
		results = append(results, result)
	}

	return results, nil
}

// GetFacetValues returns values for a specific facet
func (es *ElasticsearchService) GetFacetValues(
	ctx context.Context,
	restaurantID int64,
	facetName string,
	query string,
	limit int,
) ([]*models.FacetValue, error) {
	// Build aggregation query for specific facet
	// In production, would use Elasticsearch aggregations

	values := make([]*models.FacetValue, 0)

	// Mock response
	switch facetName {
	case "category":
		values = append(values,
			&models.FacetValue{Value: "button", Label: "Button", Count: 50},
			&models.FacetValue{Value: "form", Label: "Form", Count: 45},
			&models.FacetValue{Value: "navigation", Label: "Navigation", Count: 38},
		)
	case "rating":
		values = append(values,
			&models.FacetValue{Value: "5", Label: "5 Stars", Count: 100},
			&models.FacetValue{Value: "4", Label: "4+ Stars", Count: 80},
		)
	}

	return values, nil
}

// getSuggestions generates search suggestions
func (es *ElasticsearchService) getSuggestions(ctx context.Context, query string) []string {
	// In production, would use spell checking and related searches
	suggestions := []string{
		query + "s",           // Plural
		"related " + query,    // Related
		query + " component",  // Add type
	}

	return suggestions
}

// recordSearch records search analytics
func (es *ElasticsearchService) recordSearch(
	ctx context.Context,
	restaurantID int64,
	req *models.SearchRequest,
	resultCount int,
) error {
	searchQuery := &models.SearchQuery{
		RestaurantID: restaurantID,
		Query:        req.Query,
		QueryType:    req.QueryType,
		ResultCount:  resultCount,
		CreatedAt:    time.Now(),
	}

	return es.searchRepo.RecordSearch(ctx, searchQuery)
}

// IndexDocument adds a document to the search index
func (es *ElasticsearchService) IndexDocument(
	ctx context.Context,
	restaurantID int64,
	doc map[string]interface{},
) error {
	// In production, would call Elasticsearch API
	// For now, just simulate
	return nil
}

// DeleteDocument removes a document from the search index
func (es *ElasticsearchService) DeleteDocument(
	ctx context.Context,
	restaurantID int64,
	documentID string,
) error {
	// In production, would call Elasticsearch API
	return nil
}

// UpdateDocument updates a document in the search index
func (es *ElasticsearchService) UpdateDocument(
	ctx context.Context,
	restaurantID int64,
	documentID string,
	doc map[string]interface{},
) error {
	// In production, would call Elasticsearch API
	return nil
}

// BulkIndex indexes multiple documents
func (es *ElasticsearchService) BulkIndex(
	ctx context.Context,
	restaurantID int64,
	documents []map[string]interface{},
) error {
	// In production, would use bulk API for efficiency
	for _, doc := range documents {
		_ = es.IndexDocument(ctx, restaurantID, doc)
	}
	return nil
}

// RefreshIndex refreshes the search index
func (es *ElasticsearchService) RefreshIndex(
	ctx context.Context,
	restaurantID int64,
	queryType string,
) error {
	// In production, would call Elasticsearch refresh API
	indexName := es.getIndexName(queryType)
	return es.searchRepo.RefreshIndex(ctx, restaurantID, indexName)
}

// getIndexName returns the Elasticsearch index name
func (es *ElasticsearchService) getIndexName(queryType string) string {
	if queryType == "" {
		queryType = "all"
	}
	return fmt.Sprintf("pos-search-%s", strings.ToLower(queryType))
}

// GetSearchAnalytics returns search analytics for a period
func (es *ElasticsearchService) GetSearchAnalytics(
	ctx context.Context,
	restaurantID int64,
	period string,
	periodStart time.Time,
	periodEnd time.Time,
) (*models.SearchAnalytics, error) {
	analytics, err := es.searchRepo.GetSearchAnalytics(ctx, restaurantID, period, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	// Calculate derived metrics
	if analytics.TotalSearches > 0 {
		analytics.AvgResultsPerQuery = int(100 / analytics.TotalSearches) // Mock
	}

	return analytics, nil
}

// SaveSearch saves a search query for later
func (es *ElasticsearchService) SaveSearch(
	ctx context.Context,
	restaurantID int64,
	userID int64,
	name string,
	query string,
	filters []*models.SearchFilter,
	sortBy string,
) (*models.SavedSearch, error) {
	filterMap := make(map[string]interface{})
	for _, f := range filters {
		filterMap[f.Field] = f.Value
	}

	saved := &models.SavedSearch{
		RestaurantID:   restaurantID,
		UserID:         userID,
		Name:           name,
		Query:          query,
		Filters:        filterMap,
		SortBy:         sortBy,
		IsPublic:       false,
		ExecutionCount: 0,
		CreatedAt:      time.Now(),
	}

	return es.searchRepo.SaveSearch(ctx, saved)
}

// GetSavedSearches returns saved searches for a user
func (es *ElasticsearchService) GetSavedSearches(
	ctx context.Context,
	restaurantID int64,
	userID int64,
) ([]*models.SavedSearch, error) {
	return es.searchRepo.GetSavedSearches(ctx, restaurantID, userID)
}

// ExecuteSavedSearch executes a saved search
func (es *ElasticsearchService) ExecuteSavedSearch(
	ctx context.Context,
	savedSearchID int64,
) (*models.SearchResponse, error) {
	// Get saved search
	saved, err := es.searchRepo.GetSavedSearch(ctx, savedSearchID)
	if err != nil {
		return nil, err
	}

	// Build request from saved search
	filters := make([]*models.SearchFilter, 0)
	for field, value := range saved.Filters {
		filters = append(filters, &models.SearchFilter{
			Field:    field,
			Value:    value,
			Operator: "equals",
		})
	}

	req := &models.SearchRequest{
		Query:    saved.Query,
		Filters:  filters,
		SortBy:   saved.SortBy,
		Page:     1,
		PageSize: 20,
	}

	// Execute search
	resp, err := es.Search(ctx, saved.RestaurantID, req)
	if err != nil {
		return nil, err
	}

	// Increment execution count
	saved.ExecutionCount++
	now := time.Now()
	saved.LastRun = &now
	_ = es.searchRepo.UpdateSavedSearch(ctx, saved)

	return resp, nil
}

// AnalyzeQuery analyzes search query for insights
func (es *ElasticsearchService) AnalyzeQuery(query string) map[string]interface{} {
	// Tokenize query
	terms := strings.Fields(strings.ToLower(query))

	// Count tokens
	analysis := map[string]interface{}{
		"query":          query,
		"tokens":         terms,
		"token_count":    len(terms),
		"estimated_type": es.estimateQueryType(terms),
	}

	return analysis
}

// estimateQueryType estimates what the user is searching for
func (es *ElasticsearchService) estimateQueryType(terms []string) string {
	for _, term := range terms {
		switch term {
		case "component", "button", "form", "input":
			return "component"
		case "menu", "navigation", "sidebar":
			return "menu"
		case "layout", "grid", "flexbox":
			return "layout"
		}
	}
	return "all"
}

// CalculateRelevance calculates relevance score for search results
func (es *ElasticsearchService) CalculateRelevance(
	query string,
	result *models.SearchResult,
) float64 {
	// TF-IDF based relevance scoring
	score := 0.0

	// Title match has higher weight
	if strings.Contains(strings.ToLower(result.Title), strings.ToLower(query)) {
		score += 0.5
	}

	// Description match
	if strings.Contains(strings.ToLower(result.Description), strings.ToLower(query)) {
		score += 0.3
	}

	// Tag match
	for _, tag := range result.Tags {
		if strings.EqualFold(tag, query) {
			score += 0.2
			break
		}
	}

	// Normalize
	return math.Min(score, 1.0)
}
