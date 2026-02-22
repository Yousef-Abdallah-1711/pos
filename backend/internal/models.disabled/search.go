//go:build legacy
// +build legacy

package models

import (
	"time"
)

// SearchQuery represents a user search query
type SearchQuery struct {
	ID             int64                  `gorm:"primaryKey" json:"id"`
	RestaurantID   int64                  `json:"restaurant_id"`
	UserID         *int64                 `json:"user_id"`
	Query          string                 `json:"query"`
	QueryType      string                 `json:"query_type"` // "component", "menu", "layout", "user", "all"
	Filters        map[string]interface{} `json:"filters"`
	SortBy         string                 `json:"sort_by"` // "relevance", "popularity", "recent", "rating"
	ResultCount    int                    `json:"result_count"`
	ExecutionTime  int64                  `json:"execution_time"` // milliseconds
	ClickedResult  *int64                 `json:"clicked_result"` // Which result was clicked
	CreatedAt      time.Time              `json:"created_at"`
}

// SearchResult represents a single search result
type SearchResult struct {
	ID              int64           `json:"id"`
	ResourceID      int64           `json:"resource_id"`
	ResourceType    string          `json:"resource_type"` // "component", "menu", "layout", "user"
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Snippet         string          `json:"snippet"` // Highlighted excerpt
	Score           float64         `json:"score"` // Relevance score (0-1)
	Rank            int             `json:"rank"` // Position in results
	URL             string          `json:"url"`
	Thumbnail       *string         `json:"thumbnail"`
	Metadata        map[string]interface{} `json:"metadata"`
	Tags            []string        `json:"tags"`
	Category        string          `json:"category"`
	Rating          float64         `json:"rating"`
	Views           int64           `json:"views"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// FacetValue represents a single facet value
type FacetValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"` // Number of results with this value
}

// Facet represents a search filter dimension
type Facet struct {
	Name   string        `json:"name"`
	Label  string        `json:"label"`
	Type   string        `json:"type"` // "string", "numeric", "date", "boolean"
	Values []*FacetValue `json:"values"`
}

// SearchAggregation represents aggregated search statistics
type SearchAggregation struct {
	ResourceType map[string]int    `json:"resource_type"` // Type counts
	Categories   map[string]int    `json:"categories"`
	DateRange    map[string]int    `json:"date_range"` // By time period
	RatingBuckets map[string]int   `json:"rating_buckets"` // By rating range
	PopularityBuckets map[string]int `json:"popularity_buckets"`
}

// TrendingComponent represents a trending search item
type TrendingComponent struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	RestaurantID   int64           `json:"restaurant_id"`
	ComponentID    int64           `json:"component_id"`
	Title          string          `json:"title"`
	Category       string          `json:"category"`
	SearchVolume   int64           `json:"search_volume"` // Times searched
	ViewCount      int64           `json:"view_count"` // Times viewed from search
	ClickCount     int64           `json:"click_count"` // Times clicked
	CTR            float64         `json:"ctr"` // Click-through rate
	Trend          float64         `json:"trend"` // Change in popularity (1.0 = stable)
	TrendDirection string          `json:"trend_direction"` // "rising", "falling", "stable"
	Rank           int             `json:"rank"` // 1-based ranking
	Period         string          `json:"period"` // "daily", "weekly", "monthly"
	PeriodStart    time.Time       `json:"period_start"`
	PeriodEnd      time.Time       `json:"period_end"`
	CreatedAt      time.Time       `json:"created_at"`
}

// SearchAnalytics tracks search behavior for insights
type SearchAnalytics struct {
	ID                int64                  `gorm:"primaryKey" json:"id"`
	RestaurantID      int64                  `json:"restaurant_id"`
	Period            string                 `json:"period"` // "daily", "weekly", "monthly"
	PeriodStart       time.Time              `json:"period_start"`
	PeriodEnd         time.Time              `json:"period_end"`
	TotalSearches     int64                  `json:"total_searches"`
	UniqueUsers       int64                  `json:"unique_users"`
	AvgResultsPerQuery int                   `json:"avg_results_per_query"`
	AvgQueryTime      int64                  `json:"avg_query_time"` // ms
	ZeroResultsRate   float64                `json:"zero_results_rate"` // Percentage
	MostSearchedTerms []string               `json:"most_searched_terms"`
	TopClicks         map[string]int64       `json:"top_clicks"` // Component ID -> clicks
	QueryByType       map[string]int64       `json:"query_by_type"` // Type -> count
	CreatedAt         time.Time              `json:"created_at"`
}

// SearchSuggestion represents autocomplete suggestions
type SearchSuggestion struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	RestaurantID   int64     `json:"restaurant_id"`
	Prefix         string    `json:"prefix"`
	Suggestion     string    `json:"suggestion"`
	Type           string    `json:"type"` // "component", "menu", "user", "saved_search"
	Popularity     int64     `json:"popularity"` // Based on search frequency
	CTR            float64   `json:"ctr"` // Historical click-through rate
	LastSearched   time.Time `json:"last_searched"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SavedSearch represents a saved search query
type SavedSearch struct {
	ID             int64                  `gorm:"primaryKey" json:"id"`
	RestaurantID   int64                  `json:"restaurant_id"`
	UserID         int64                  `json:"user_id"`
	Name           string                 `json:"name"` // User-friendly name
	Query          string                 `json:"query"`
	Filters        map[string]interface{} `json:"filters"`
	SortBy         string                 `json:"sort_by"`
	IsPublic       bool                   `json:"is_public"` // Shared with others
	ExecutionCount int                    `json:"execution_count"` // Times run
	LastRun        *time.Time             `json:"last_run"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// SearchFilter represents a filter applied to search
type SearchFilter struct {
	Field    string        `json:"field"`
	Operator string        `json:"operator"` // "equals", "contains", "range", "in", "gt", "lt"
	Value    interface{}   `json:"value"`
}

// SearchRequest represents a complete search request
type SearchRequest struct {
	Query          string           `json:"query"`
	QueryType      string           `json:"query_type"` // "component", "menu", "layout", "user", "all"
	Filters        []*SearchFilter  `json:"filters"`
	Facets         []string         `json:"facets"` // Which facets to return
	SortBy         string           `json:"sort_by"`
	SortOrder      string           `json:"sort_order"` // "asc", "desc"
	Page           int              `json:"page"`
	PageSize       int              `json:"page_size"`
	HighlightFields []string        `json:"highlight_fields"`
	Explain        bool             `json:"explain"` // Return scoring explanation
}

// SearchResponse represents a complete search response
type SearchResponse struct {
	Query        string              `json:"query"`
	TotalHits    int64               `json:"total_hits"`
	Results      []*SearchResult     `json:"results"`
	Facets       []*Facet            `json:"facets"`
	Aggregations *SearchAggregation  `json:"aggregations"`
	Suggestions  []string            `json:"suggestions"`
	ExecutionTime int64              `json:"execution_time"` // ms
	Page         int                 `json:"page"`
	PageSize     int                 `json:"page_size"`
	TotalPages   int                 `json:"total_pages"`
}

// SearchIndex represents Elasticsearch index configuration
type SearchIndex struct {
	ID            int64                  `gorm:"primaryKey" json:"id"`
	IndexName     string                 `json:"index_name"`
	ResourceType  string                 `json:"resource_type"` // "component", "menu", etc.
	Mapping       map[string]interface{} `json:"mapping"` // Elasticsearch mapping
	DocCount      int64                  `json:"doc_count"`
	IndexSize     int64                  `json:"index_size"` // bytes
	LastRefresh   time.Time              `json:"last_refresh"`
	Status        string                 `json:"status"` // "healthy", "degraded", "unhealthy"
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// SearchConfig stores search configuration
type SearchConfig struct {
	ID                    int64           `gorm:"primaryKey" json:"id"`
	RestaurantID          int64           `json:"restaurant_id"`
	EnableFullText        bool            `json:"enable_full_text"`
	EnableAutoComplete    bool            `json:"enable_auto_complete"`
	EnableFacets          bool            `json:"enable_facets"`
	EnableTrendingItems   bool            `json:"enable_trending_items"`
	EnableAnalytics       bool            `json:"enable_analytics"`
	SearchTimeout         int             `json:"search_timeout"` // seconds
	MaxResults            int             `json:"max_results"`
	MinScoreThreshold     float64         `json:"min_score_threshold"`
	FuzzyMatching         bool            `json:"fuzzy_matching"`
	FuzzyMaxEdits         int             `json:"fuzzy_max_edits"` // Levenshtein distance
	BoostRecent           bool            `json:"boost_recent"`
	BoostPopular          bool            `json:"boost_popular"`
	EnableSpellCheck      bool            `json:"enable_spell_check"`
	IndexRefreshInterval  int             `json:"index_refresh_interval"` // seconds
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}
