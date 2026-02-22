//go:build legacy
// +build legacy

package repository

import (
	"context"
	"time"
)

// SearchRepository handles search-related data persistence
type SearchRepository interface {
	// Autocomplete suggestions
	GetAutocompleteSuggestions(ctx context.Context, restaurantID int64, prefix string, limit int) ([]string, error)
	CreateSuggestion(ctx context.Context, suggestion *models.SearchSuggestion) error
	UpdateSuggestion(ctx context.Context, suggestion *models.SearchSuggestion) error

	// Trending items
	GetTrendingComponents(ctx context.Context, restaurantID int64, period string, limit int) ([]*models.TrendingComponent, error)
	CreateTrendingComponent(ctx context.Context, component *models.TrendingComponent) error
	UpdateTrendingComponent(ctx context.Context, component *models.TrendingComponent) error

	// Search queries
	RecordSearch(ctx context.Context, search *models.SearchQuery) error
	GetSearchHistory(ctx context.Context, restaurantID int64, userID *int64, limit int) ([]*models.SearchQuery, error)

	// Saved searches
	SaveSearch(ctx context.Context, search *models.SavedSearch) (*models.SavedSearch, error)
	GetSavedSearch(ctx context.Context, id int64) (*models.SavedSearch, error)
	GetSavedSearches(ctx context.Context, restaurantID int64, userID int64) ([]*models.SavedSearch, error)
	UpdateSavedSearch(ctx context.Context, search *models.SavedSearch) error
	DeleteSavedSearch(ctx context.Context, id int64) error

	// Search analytics
	GetSearchAnalytics(ctx context.Context, restaurantID int64, period string, start time.Time, end time.Time) (*models.SearchAnalytics, error)
	RecordSearchAnalytics(ctx context.Context, analytics *models.SearchAnalytics) error

	// Index management
	GetSearchIndex(ctx context.Context, restaurantID int64, resourceType string) (*models.SearchIndex, error)
	UpdateSearchIndex(ctx context.Context, index *models.SearchIndex) error
	RefreshIndex(ctx context.Context, restaurantID int64, indexName string) error

	// Configuration
	GetSearchConfig(ctx context.Context, restaurantID int64) (*models.SearchConfig, error)
	UpdateSearchConfig(ctx context.Context, config *models.SearchConfig) error

	// Related searches
	GetRelatedSearches(ctx context.Context, query string, limit int) ([]string, error)
	RecordRelatedSearch(ctx context.Context, original string, related string) error
}
